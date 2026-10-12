/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pricing

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/aks-node-viewer/pkg/pricing/client"
)

type fakePricingAPI struct {
	pages []*client.ProductsPricePage
}

func (f *fakePricingAPI) GetProductsPricePages(ctx context.Context, filters []*client.Filter, pageHandler func(output *client.ProductsPricePage)) error {
	for _, p := range f.pages {
		pageHandler(p)
	}
	return nil
}

func onDemandItem() client.Item {
	return client.Item{
		ProductName: "Virtual Machines Dv3 Series",
		MeterName:   "D2s v3",
		SkuName:     "D2s_v3",
		ArmSkuName:  "Standard_D2s_v3",
		RetailPrice: 0.100,
	}
}

func spotItem() client.Item {
	return client.Item{
		ProductName: "Virtual Machines Dv3 Series",
		MeterName:   "D2s v3",
		SkuName:     "D2s_v3 Spot",
		ArmSkuName:  "Standard_D2s_v3",
		RetailPrice: 0.050,
	}
}

func liveOnlyItem() client.Item {
	return client.Item{
		ProductName: "Virtual Machines Dv3 Series",
		MeterName:   "Magic Instance",
		SkuName:     "Magic_Instance_9",
		ArmSkuName:  "Magic_Instance_9",
		RetailPrice: 0.100,
	}
}

func TestProviderLazyRefresh(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	api := &fakePricingAPI{pages: []*client.ProductsPricePage{
		{Items: []client.Item{liveOnlyItem(), onDemandItem(), spotItem()}},
	}}
	pprov := NewProvider(context.Background(), api, "westus2")

	// static prices must be available immediately, before any refresh completes
	if price, ok := pprov.OnDemandPrice("Basic_A0"); !ok || price <= 0 {
		t.Errorf("expected static on-demand pricing to be available immediately")
	}

	pprov.Refresh(context.Background())

	deadline := time.Now().Add(5 * time.Second)
	var onDemand, spot, replaced float64
	var okOD, okSpot, okRepl bool
	for time.Now().Before(deadline) {
		select {
		case <-pprov.Watch():
		case <-time.After(100 * time.Millisecond):
		}
		onDemand, okOD = pprov.OnDemandPrice("Magic_Instance_9")
		spot, okSpot = pprov.SpotPrice("Standard_D2s_v3")
		replaced, okRepl = pprov.OnDemandPrice("Standard_D2s_v3")
		if okOD && okRepl && okSpot && onDemand == 0.100 && spot == 0.050 && replaced == 0.100 {
			return
		}
	}
	t.Fatalf("timed out waiting for the pricing refresh to complete: onDemand=%v okOD=%v spot=%v okSpot=%v replaced=%v okRepl=%v",
		onDemand, okOD, spot, okSpot, replaced, okRepl)
}

func TestProviderStaticWithoutRefresh(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	api := &fakePricingAPI{pages: []*client.ProductsPricePage{
		{Items: []client.Item{liveOnlyItem(), onDemandItem(), spotItem()}},
	}}
	pprov := NewProvider(context.Background(), api, "not-a-region")

	// an unknown region falls back to the eastus static prices
	if price, ok := pprov.OnDemandPrice("Basic_A0"); !ok || price <= 0 {
		t.Errorf("expected static pricing fallback, got price=%v ok=%v", price, ok)
	}
	// without Refresh, live prices must never replace the static ones
	time.Sleep(200 * time.Millisecond)
	if _, ok := pprov.OnDemandPrice("Magic_Instance_9"); ok {
		t.Errorf("expected live prices to remain unloaded without Refresh")
	}
}

func TestProviderSetRegion(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	cachedTime := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	if err := writeCache(&cacheData{
		Region:             "eastus",
		OnDemandUpdateTime: cachedTime,
		OnDemandPrices:     map[string]float64{"Cached_Instance": 1.5},
		SpotUpdateTime:     cachedTime,
		SpotPrices:         map[string]float64{"Cached_Instance": 0.75},
	}); err != nil {
		t.Fatalf("writeCache: %v", err)
	}

	pprov := NewProvider(context.Background(), &fakePricingAPI{}, "westus2")
	if got := pprov.Region(); got != "westus2" {
		t.Fatalf("expected region westus2, got %q", got)
	}
	if got := pprov.Status().Region; got != "westus2" {
		t.Errorf("expected the status to report the region, got %q", got)
	}

	// switching to a region with a cache reseeds the prices from it
	if !pprov.SetRegion("eastus") {
		t.Fatalf("expected SetRegion to report a change")
	}
	if got := pprov.Region(); got != "eastus" {
		t.Errorf("expected region eastus, got %q", got)
	}
	if price, ok := pprov.OnDemandPrice("Cached_Instance"); !ok || price != 1.5 {
		t.Errorf("expected the cached on-demand price, got price=%v ok=%v", price, ok)
	}
	if price, ok := pprov.SpotPrice("Cached_Instance"); !ok || price != 0.75 {
		t.Errorf("expected the cached spot price, got price=%v ok=%v", price, ok)
	}
	if got := pprov.Status().Source; got != SourceCache {
		t.Errorf("expected the source to be SourceCache, got %v", got)
	}

	// no change is reported for the same region or an empty one
	if pprov.SetRegion("eastus") {
		t.Errorf("expected no change for the same region")
	}
	if pprov.SetRegion("") {
		t.Errorf("expected no change for an empty region")
	}
}

// gatedPricingAPI blocks fetches for one region until the gate is closed, so a
// test can change the provider's region while a refresh is in flight.
type gatedPricingAPI struct {
	blockRegion  string
	gate         chan struct{}
	started      chan string
	regionPrices map[string]float64
}

func (g *gatedPricingAPI) GetProductsPricePages(ctx context.Context, filters []*client.Filter, pageHandler func(output *client.ProductsPricePage)) error {
	region := ""
	for _, f := range filters {
		if f.Field == "armRegionName" {
			region = f.Value
		}
	}
	g.started <- region
	if region == g.blockRegion {
		<-g.gate
	}
	pageHandler(&client.ProductsPricePage{Items: []client.Item{{
		ProductName: "Virtual Machines Test",
		MeterName:   "Test",
		SkuName:     "Test",
		ArmSkuName:  "Instance",
		RetailPrice: g.regionPrices[region],
	}}})
	return nil
}

func TestProviderRegionChangeWhileRefreshing(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	api := &gatedPricingAPI{
		blockRegion:  "region-a",
		gate:         make(chan struct{}),
		started:      make(chan string, 16),
		regionPrices: map[string]float64{"region-a": 1.0, "region-b": 2.0},
	}
	pprov := NewProvider(context.Background(), api, "region-a")
	pprov.Refresh(context.Background())

	// wait until the region-a fetch is in flight and blocked
	waitForFetch(t, api.started, "region-a")

	// change the region while that refresh is still in flight
	if !pprov.SetRegion("region-b") {
		t.Fatalf("expected the region to change")
	}
	pprov.Refresh(context.Background())

	// release the in-flight region-a fetch: its results must be discarded and
	// the follow-up region-b refresh must win
	close(api.gate)

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if price, ok := pprov.OnDemandPrice("Instance"); ok && price == 2.0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	price, ok := pprov.OnDemandPrice("Instance")
	t.Fatalf("timed out waiting for the region-b price: price=%v ok=%v region=%q", price, ok, pprov.Region())
}

// waitForFetch waits until the pricing API is asked for the given region.
func waitForFetch(t *testing.T, started <-chan string, region string) {
	t.Helper()
	timeout := time.After(5 * time.Second)
	for {
		select {
		case got := <-started:
			if got == region {
				return
			}
		case <-timeout:
			t.Fatalf("timed out waiting for a %s fetch", region)
		}
	}
}
