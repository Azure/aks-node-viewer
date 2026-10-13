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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/aks-node-viewer/pkg/pricing/client"
	"github.com/samber/lo"
)

// Source identifies where the pricing data currently in use came from.
type Source int

const (
	// SourceStatic is the price list embedded at build time. It is the least
	// accurate source and is only used when no cache is available.
	SourceStatic Source = iota
	// SourceCache is pricing persisted by a previous run and loaded from the
	// XDG cache directory.
	SourceCache
	// SourceLive is pricing fetched from the Azure retail prices API during
	// this run.
	SourceLive
)

func (s Source) String() string {
	switch s {
	case SourceLive:
		return "live"
	case SourceCache:
		return "cache"
	default:
		return "static"
	}
}

// DefaultRegion is used for pricing lookups when no region is configured and
// the cluster's region cannot (yet) be determined from its nodes. The embedded
// static price snapshot only contains eastus prices.
const DefaultRegion = "eastus"

// Status describes the freshness of the pricing data currently in use.
type Status struct {
	// Source is where the current pricing data came from.
	Source Source
	// Updated is when the current pricing data was obtained.
	Updated time.Time
	// Updating is true while a live refresh is in flight.
	Updating bool
	// Region is the Azure region the prices were looked up for.
	Region string
	// Error describes why the most recent refresh failed, if it did. It is
	// empty when the last refresh succeeded.
	Error string
}

// Provider provides actual pricing data to the Azure cloud provider to allow it to make more informed decisions
// regarding which instances to launch.  This is initialized at startup with a periodically updated static price list to
// support running in locations where pricing data is unavailable.  In those cases the static pricing data provides a
// relative ordering that is still more accurate than our previous pricing model.  In the event that a pricing update
// fails, the previous pricing information is retained and used which may be the static initial pricing data if pricing
// updates never succeed.
type Provider struct {
	pricing client.PricingAPI
	region  string

	mu                 sync.RWMutex
	source             Source
	updating           bool
	refreshPending     bool
	statusUpdated      time.Time
	onDemandUpdateTime time.Time
	onDemandPrices     map[string]float64
	spotUpdateTime     time.Time
	spotPrices         map[string]float64

	// onDemandErr and spotErr record why the most recent update of each price
	// set failed, so the UI can explain why prices are stale instead of leaving
	// the failure invisible.
	onDemandErr error
	spotErr     error

	// watch is signalled (coalesced) whenever a pricing refresh completes
	watch chan struct{}

	// cacheMu serializes cache writes so concurrent on-demand and spot
	// refreshes cannot leave an older snapshot on disk.
	cacheMu sync.Mutex
}

type Err struct {
	error
	lastUpdateTime time.Time
}

// NewPricingAPI returns a pricing API
func NewAPI() client.PricingAPI {
	return client.New()
}

func NewProvider(ctx context.Context, pricing client.PricingAPI, region string) *Provider {
	p := &Provider{
		pricing: pricing,
		watch:   make(chan struct{}, 1),
	}
	p.seedLocked(region)
	return p
}

// seedLocked initialises the prices for a region from the cache written by a
// previous run, falling back to the embedded static snapshot. It must be called
// with p.mu held.
func (p *Provider) seedLocked(region string) {
	// see if we've got region specific pricing data
	staticPricing, ok := initialOnDemandPrices[region]
	if !ok {
		// and if not, fall back to the always available eastus
		staticPricing = initialOnDemandPrices["eastus"]
	}

	p.region = region
	p.source = SourceStatic
	p.statusUpdated = initialPriceUpdate
	p.onDemandUpdateTime = initialPriceUpdate
	p.onDemandPrices = staticPricing
	p.spotUpdateTime = initialPriceUpdate
	// default our spot pricing to the same as the on-demand pricing until a price update
	p.spotPrices = staticPricing

	// any error from a previous region no longer applies
	p.onDemandErr = nil
	p.spotErr = nil

	// Prefer pricing persisted by a previous run over the embedded static
	// snapshot: it gives a much more accurate baseline while the live refresh
	// runs in the background.
	if cached, ok := readCache(region); ok {
		p.source = SourceCache
		if len(cached.OnDemandPrices) > 0 {
			p.onDemandPrices = lo.Assign(cached.OnDemandPrices)
			p.onDemandUpdateTime = cached.OnDemandUpdateTime
			p.statusUpdated = cached.OnDemandUpdateTime
		}
		if len(cached.SpotPrices) > 0 {
			p.spotPrices = lo.Assign(cached.SpotPrices)
			p.spotUpdateTime = cached.SpotUpdateTime
			if p.statusUpdated.Before(cached.SpotUpdateTime) {
				p.statusUpdated = cached.SpotUpdateTime
			}
		}
	}
}

// SetRegion changes the region used for pricing lookups and reseeds the
// provider from that region's cache (or the embedded static snapshot). It
// reports whether the region changed; call Refresh to fetch live prices for the
// new region. Prices already fetched for the previous region stay in effect
// until the refresh completes.
func (p *Provider) SetRegion(region string) bool {
	if region == "" {
		return false
	}
	p.mu.Lock()
	if region == p.region {
		p.mu.Unlock()
		return false
	}
	p.seedLocked(region)
	p.mu.Unlock()
	p.signalUpdate()
	return true
}

// Region returns the region prices are currently looked up for.
func (p *Provider) Region() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.region
}

// Status reports the freshness of the pricing data currently in use.
func (p *Provider) Status() Status {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return Status{
		Source:   p.source,
		Updated:  p.statusUpdated,
		Updating: p.updating,
		Region:   p.region,
		Error:    p.pricingErrorLocked(),
	}
}

// pricingErrorLocked describes the most recent pricing failures. It must be
// called with p.mu held.
func (p *Provider) pricingErrorLocked() string {
	switch {
	case p.onDemandErr != nil && p.spotErr != nil:
		if p.onDemandErr.Error() == p.spotErr.Error() {
			return p.onDemandErr.Error()
		}
		return fmt.Sprintf("%s; %s", p.onDemandErr, p.spotErr)
	case p.onDemandErr != nil:
		return p.onDemandErr.Error()
	case p.spotErr != nil:
		return p.spotErr.Error()
	}
	return ""
}

// persistCache writes the current pricing data to the XDG cache directory so
// that the next run can start from it. Cache writes are best-effort.
func (p *Provider) persistCache() {
	p.cacheMu.Lock()
	defer p.cacheMu.Unlock()

	p.mu.RLock()
	data := &cacheData{
		Region:             p.region,
		OnDemandUpdateTime: p.onDemandUpdateTime,
		OnDemandPrices:     p.onDemandPrices,
		SpotUpdateTime:     p.spotUpdateTime,
		SpotPrices:         p.spotPrices,
	}
	p.mu.RUnlock()
	_ = writeCache(data)
}

// Refresh starts loading live pricing data for the current region in the
// background. It returns immediately: until the refresh completes the provider
// keeps serving whatever prices it already has, so the UI is never blocked.
// Consumers can use Watch to react to completed refreshes.
//
// At most one refresh runs at a time. A request that arrives while a refresh is
// in flight is coalesced into a follow-up refresh, so a region change made
// mid-refresh is still fetched (and the results of the in-flight refresh for the
// old region are discarded).
func (p *Provider) Refresh(ctx context.Context) {
	p.mu.Lock()
	if p.updating {
		p.refreshPending = true
		p.mu.Unlock()
		return
	}
	p.updating = true
	p.mu.Unlock()
	go p.refreshLoop(ctx)
}

func (p *Provider) refreshLoop(ctx context.Context) {
	for {
		p.updatePricing(ctx)

		p.mu.Lock()
		if p.refreshPending {
			p.refreshPending = false
			p.mu.Unlock()
			continue
		}
		p.updating = false
		p.mu.Unlock()
		return
	}
}

// Watch returns a channel that receives a signal whenever pricing data has been
// refreshed. Signals are coalesced: a refresh that completes while a previous signal
// is still pending may not produce an additional signal.
func (p *Provider) Watch() <-chan struct{} {
	return p.watch
}

func (p *Provider) signalUpdate() {
	select {
	case p.watch <- struct{}{}:
	default:
	}
}

// InstanceTypes returns the list of all instance types for which either a price is known.
func (p *Provider) InstanceTypes() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return lo.Union(lo.Keys(p.onDemandPrices), lo.Keys(p.spotPrices))
}

// OnDemandLastUpdated returns the time that the on-demand pricing was last updated
func (p *Provider) OnDemandLastUpdated() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.onDemandUpdateTime
}

// SpotLastUpdated returns the time that the spot pricing was last updated
func (p *Provider) SpotLastUpdated() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.spotUpdateTime
}

// OnDemandPrice returns the last known on-demand price for a given instance type, returning false if there is no
// known on-demand pricing for the instance type.
func (p *Provider) OnDemandPrice(instanceType string) (float64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	price, ok := p.onDemandPrices[instanceType]
	if !ok {
		return 0.0, false
	}
	return price, true
}

// SpotPrice returns the last known spot price for a given instance type, returning false
// if there is no known spot pricing for that instance type
func (p *Provider) SpotPrice(instanceType string) (float64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	price, ok := p.spotPrices[instanceType]
	if !ok {
		return 0.0, false
	}
	return price, true
}

func (p *Provider) updatePricing(ctx context.Context) {
	// Snapshot the region so a SetRegion during the fetch cannot mix prices from
	// two regions in a single update.
	region := p.Region()

	onDemandPrices := map[string]float64{}
	spotPrices := map[string]float64{}
	onDemandHandler := onDemandPage(onDemandPrices)
	spotHandler := spotPage(spotPrices)

	// Feed both price sets from a single paginated fetch: the API rate limits
	// requests, so fetching on-demand and spot separately doubles the number of
	// calls and makes a partial failure much more likely.
	err := p.fetchPricing(ctx, region, func(page *client.ProductsPricePage) {
		onDemandHandler(page)
		spotHandler(page)
	})

	_ = p.storeOnDemandPricing(region, onDemandPrices, err)
	_ = p.storeSpotPricing(region, spotPrices, err)
}

// UpdateOnDemandPricing fetches the on-demand prices for the current region.
func (p *Provider) UpdateOnDemandPricing(ctx context.Context) *Err {
	region := p.Region()
	onDemandPrices := map[string]float64{}
	err := p.fetchPricing(ctx, region, onDemandPage(onDemandPrices))
	return p.storeOnDemandPricing(region, onDemandPrices, err)
}

// storeOnDemandPricing applies a fetched set of on-demand prices. It reports why
// the update failed when the fetch failed or returned no prices, and drops the
// results when the region changed while the fetch was in flight.
func (p *Provider) storeOnDemandPricing(region string, prices map[string]float64, fetchErr error) *Err {
	p.mu.Lock()
	if region != p.region {
		// the region changed while we were fetching; the follow-up refresh will
		// fetch the new region, so drop these results
		p.mu.Unlock()
		return nil
	}
	lastUpdateTime := p.onDemandUpdateTime
	if fetchErr != nil {
		err := fmt.Errorf("pricing update failed: %w", fetchErr)
		p.onDemandErr = err
		p.mu.Unlock()
		return &Err{error: err, lastUpdateTime: lastUpdateTime}
	}
	if len(prices) == 0 {
		err := errors.New("no on-demand pricing found")
		p.onDemandErr = err
		p.mu.Unlock()
		return &Err{error: err, lastUpdateTime: lastUpdateTime}
	}

	p.onDemandPrices = lo.Assign(prices)
	p.onDemandUpdateTime = time.Now()
	p.onDemandErr = nil
	p.source = SourceLive
	p.statusUpdated = p.onDemandUpdateTime
	p.mu.Unlock()

	p.signalUpdate()
	p.persistCache()
	return nil
}

func (p *Provider) fetchPricing(ctx context.Context, region string, pageHandler func(output *client.ProductsPricePage)) error {
	filters := []*client.Filter{
		{
			Field:    "priceType",
			Operator: client.Equals,
			Value:    "Consumption",
		},
		{
			Field:    "currencyCode",
			Operator: client.Equals,
			Value:    "USD",
		},
		{
			Field:    "serviceFamily",
			Operator: client.Equals,
			Value:    "Compute",
		},
		{
			Field:    "serviceName",
			Operator: client.Equals,
			Value:    "Virtual Machines",
		},
		{
			Field:    "armRegionName",
			Operator: client.Equals,
			Value:    region,
		}}
	return p.pricing.GetProductsPricePages(ctx, filters, pageHandler)
}

func onDemandPage(prices map[string]float64) func(page *client.ProductsPricePage) {
	return func(page *client.ProductsPricePage) {
		for _, pItem := range page.Items {
			if strings.HasSuffix(pItem.ProductName, " Windows") {
				continue
			}
			if strings.HasSuffix(pItem.MeterName, " Low Priority") {
				// https://learn.microsoft.com/en-us/azure/batch/batch-spot-vms#differences-between-spot-and-low-priority-vms
				continue
			}
			if strings.HasSuffix(pItem.SkuName, " Spot") {
				continue
			}
			prices[pItem.ArmSkuName] = pItem.RetailPrice
		}
	}
}

// UpdateSpotPricing fetches the spot prices for the current region.
func (p *Provider) UpdateSpotPricing(ctx context.Context) *Err {
	region := p.Region()
	spotPrices := map[string]float64{}
	err := p.fetchPricing(ctx, region, spotPage(spotPrices))
	return p.storeSpotPricing(region, spotPrices, err)
}

// storeSpotPricing applies a fetched set of spot prices. It reports why the
// update failed when the fetch failed or returned no prices, and drops the
// results when the region changed while the fetch was in flight.
func (p *Provider) storeSpotPricing(region string, prices map[string]float64, fetchErr error) *Err {
	p.mu.Lock()
	if region != p.region {
		// the region changed while we were fetching; the follow-up refresh will
		// fetch the new region, so drop these results
		p.mu.Unlock()
		return nil
	}
	lastUpdateTime := p.spotUpdateTime
	if fetchErr != nil {
		err := fmt.Errorf("pricing update failed: %w", fetchErr)
		p.spotErr = err
		p.mu.Unlock()
		return &Err{error: err, lastUpdateTime: lastUpdateTime}
	}
	if len(prices) == 0 {
		err := errors.New("no spot pricing found")
		p.spotErr = err
		p.mu.Unlock()
		return &Err{error: err, lastUpdateTime: lastUpdateTime}
	}

	p.spotPrices = lo.Assign(prices)
	p.spotUpdateTime = time.Now()
	p.spotErr = nil
	p.source = SourceLive
	p.statusUpdated = p.spotUpdateTime
	p.mu.Unlock()

	p.signalUpdate()
	p.persistCache()
	return nil
}

func spotPage(prices map[string]float64) func(page *client.ProductsPricePage) {
	return func(page *client.ProductsPricePage) {
		for _, pItem := range page.Items {
			if strings.HasSuffix(pItem.ProductName, " Windows") {
				continue
			}
			if strings.HasSuffix(pItem.MeterName, " Low Priority") {
				// https://learn.microsoft.com/en-us/azure/batch/batch-spot-vms#differences-between-spot-and-low-priority-vms
				continue
			}
			if !strings.HasSuffix(pItem.SkuName, " Spot") {
				continue
			}
			prices[pItem.ArmSkuName] = pItem.RetailPrice
		}
	}
}

func (p *Provider) LivenessProbe(req *http.Request) error {
	// ensure we don't deadlock and nolint for the empty critical section
	p.mu.Lock()
	//nolint: staticcheck
	p.mu.Unlock()
	return nil
}
