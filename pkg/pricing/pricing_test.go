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

	pprov.StartRefresh(context.Background())

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
	// without StartRefresh, live prices must never replace the static ones
	time.Sleep(200 * time.Millisecond)
	if _, ok := pprov.OnDemandPrice("Magic_Instance_9"); ok {
		t.Errorf("expected live prices to remain unloaded without StartRefresh")
	}
}
