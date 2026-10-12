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
	"reflect"
	"testing"
	"time"

	"github.com/Azure/aks-node-viewer/pkg/pricing/client"
)

func TestCacheRoundTrip(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	onDemandTime := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	spotTime := time.Now().Add(-90 * time.Minute).UTC().Truncate(time.Second)
	data := &cacheData{
		Region:             "westus2",
		OnDemandUpdateTime: onDemandTime,
		OnDemandPrices:     map[string]float64{"Standard_D2s_v3": 0.100},
		SpotUpdateTime:     spotTime,
		SpotPrices:         map[string]float64{"Standard_D2s_v3": 0.050},
	}
	if err := writeCache(data); err != nil {
		t.Fatalf("writeCache: %v", err)
	}

	got, ok := readCache("westus2")
	if !ok {
		t.Fatalf("expected the cache to be readable")
	}
	if !got.OnDemandUpdateTime.Equal(onDemandTime) {
		t.Errorf("expected on-demand update time %v, got %v", onDemandTime, got.OnDemandUpdateTime)
	}
	if !got.SpotUpdateTime.Equal(spotTime) {
		t.Errorf("expected spot update time %v, got %v", spotTime, got.SpotUpdateTime)
	}
	if !reflect.DeepEqual(got.OnDemandPrices, data.OnDemandPrices) {
		t.Errorf("expected on-demand prices %v, got %v", data.OnDemandPrices, got.OnDemandPrices)
	}
	if !reflect.DeepEqual(got.SpotPrices, data.SpotPrices) {
		t.Errorf("expected spot prices %v, got %v", data.SpotPrices, got.SpotPrices)
	}

	// a cache for a different region must not be returned
	if _, ok := readCache("eastus"); ok {
		t.Errorf("expected no cache for a different region")
	}
}

func TestCacheRejectsCorruptData(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	if _, ok := readCache("westus2"); ok {
		t.Errorf("expected no cache when the file is missing")
	}

	// an empty cache (no prices) is not usable
	if err := writeCache(&cacheData{Region: "westus2"}); err != nil {
		t.Fatalf("writeCache: %v", err)
	}
	if _, ok := readCache("westus2"); ok {
		t.Errorf("expected an empty cache to be rejected")
	}
}

func TestProviderSeedsFromCache(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	cachedTime := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	if err := writeCache(&cacheData{
		Region:             "westus2",
		OnDemandUpdateTime: cachedTime,
		OnDemandPrices:     map[string]float64{"Cached_Instance": 1.5},
		SpotUpdateTime:     cachedTime,
		SpotPrices:         map[string]float64{"Cached_Instance": 0.75},
	}); err != nil {
		t.Fatalf("writeCache: %v", err)
	}

	pprov := NewProvider(context.Background(), &fakePricingAPI{}, "westus2")

	if price, ok := pprov.OnDemandPrice("Cached_Instance"); !ok || price != 1.5 {
		t.Errorf("expected the cached on-demand price, got price=%v ok=%v", price, ok)
	}
	if price, ok := pprov.SpotPrice("Cached_Instance"); !ok || price != 0.75 {
		t.Errorf("expected the cached spot price, got price=%v ok=%v", price, ok)
	}

	status := pprov.Status()
	if status.Source != SourceCache {
		t.Errorf("expected the source to be SourceCache, got %v", status.Source)
	}
	if !status.Updated.Equal(cachedTime) {
		t.Errorf("expected the status update time to be %v, got %v", cachedTime, status.Updated)
	}
	if status.Updating {
		t.Errorf("expected Updating to be false before StartRefresh")
	}
}

func TestProviderDefaultsToStatic(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	pprov := NewProvider(context.Background(), &fakePricingAPI{}, "westus2")
	status := pprov.Status()
	if status.Source != SourceStatic {
		t.Errorf("expected the source to be SourceStatic, got %v", status.Source)
	}
	if !status.Updated.Equal(initialPriceUpdate) {
		t.Errorf("expected the static update time, got %v", status.Updated)
	}
}

func TestProviderPersistsCache(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	api := &fakePricingAPI{pages: []*client.ProductsPricePage{
		{Items: []client.Item{onDemandItem(), spotItem()}},
	}}
	pprov := NewProvider(context.Background(), api, "westus2")

	if err := pprov.UpdateOnDemandPricing(context.Background()); err != nil {
		t.Fatalf("UpdateOnDemandPricing: %v", err)
	}
	if err := pprov.UpdateSpotPricing(context.Background()); err != nil {
		t.Fatalf("UpdateSpotPricing: %v", err)
	}

	cached, ok := readCache("westus2")
	if !ok {
		t.Fatalf("expected a cache to have been written")
	}
	if price := cached.OnDemandPrices["Standard_D2s_v3"]; price != 0.100 {
		t.Errorf("expected the on-demand price to be cached, got %v", price)
	}
	if price := cached.SpotPrices["Standard_D2s_v3"]; price != 0.050 {
		t.Errorf("expected the spot price to be cached, got %v", price)
	}

	if status := pprov.Status(); status.Source != SourceLive {
		t.Errorf("expected the source to be SourceLive, got %v", status.Source)
	}
}
