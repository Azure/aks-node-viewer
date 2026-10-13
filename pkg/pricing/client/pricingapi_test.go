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

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testAPI(server *httptest.Server) *pricingAPI {
	return &pricingAPI{baseURL: server.URL, client: server.Client()}
}

func writePage(w http.ResponseWriter, page ProductsPricePage) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(page)
}

func TestGetProductsPricePagesFollowsPagination(t *testing.T) {
	var requests int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&requests, 1) == 1 {
			writePage(w, ProductsPricePage{
				Items:        []Item{{ArmSkuName: "Standard_D2s_v3", RetailPrice: 0.1}},
				NextPageLink: server.URL + "/next",
			})
			return
		}
		writePage(w, ProductsPricePage{
			Items: []Item{{ArmSkuName: "Standard_D4s_v3", RetailPrice: 0.2}},
		})
	}))
	defer server.Close()

	var pages, items int
	err := testAPI(server).GetProductsPricePages(context.Background(), nil, func(page *ProductsPricePage) {
		pages++
		items += len(page.Items)
	})
	if err != nil {
		t.Fatalf("GetProductsPricePages: %v", err)
	}
	if pages != 2 || items != 2 {
		t.Errorf("expected 2 pages with 2 items, got %d pages and %d items", pages, items)
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Errorf("expected 2 requests, got %d", got)
	}
}

func TestGetProductsPricePagesRetriesOnRateLimit(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&requests, 1) == 1 {
			w.Header().Set("x-ms-ratelimit-retry-after", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writePage(w, ProductsPricePage{Items: []Item{{ArmSkuName: "Standard_D2s_v3", RetailPrice: 0.1}}})
	}))
	defer server.Close()

	var pages int
	err := testAPI(server).GetProductsPricePages(context.Background(), nil, func(*ProductsPricePage) {
		pages++
	})
	if err != nil {
		t.Fatalf("expected the rate limited request to succeed on retry, got %v", err)
	}
	if pages != 1 {
		t.Errorf("expected 1 page, got %d", pages)
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Errorf("expected a single retry (2 requests), got %d", got)
	}
}

func TestGetProductsPricePagesGivesUpOnPersistentRateLimit(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.Header().Set("x-ms-ratelimit-retry-after", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	err := testAPI(server).GetProductsPricePages(context.Background(), nil, func(*ProductsPricePage) {})
	if err == nil {
		t.Fatal("expected an error after exhausting retries")
	}
	if !strings.Contains(err.Error(), "rate limited") {
		t.Errorf("expected a rate limit error, got %v", err)
	}
	if got := atomic.LoadInt32(&requests); got != maxAttempts {
		t.Errorf("expected %d attempts, got %d", maxAttempts, got)
	}
}

func TestGetProductsPricePagesRetriesServerError(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&requests, 1) == 1 {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}
		writePage(w, ProductsPricePage{})
	}))
	defer server.Close()

	if err := testAPI(server).GetProductsPricePages(context.Background(), nil, func(*ProductsPricePage) {}); err != nil {
		t.Fatalf("expected the server error to be retried, got %v", err)
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Errorf("expected a retry (2 requests), got %d", got)
	}
}

func TestGetProductsPricePagesFailsOnClientError(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		http.Error(w, "bad filter", http.StatusBadRequest)
	}))
	defer server.Close()

	err := testAPI(server).GetProductsPricePages(context.Background(), nil, func(*ProductsPricePage) {})
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "bad filter") {
		t.Errorf("expected the status and body in the error, got %v", err)
	}
	if got := atomic.LoadInt32(&requests); got != 1 {
		t.Errorf("expected no retry for a client error, got %d requests", got)
	}
}

func TestGetProductsPricePagesBuildsFilter(t *testing.T) {
	var query string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery
		writePage(w, ProductsPricePage{})
	}))
	defer server.Close()

	filters := []*Filter{
		{Field: "serviceName", Operator: Equals, Value: "Virtual Machines"},
		{Field: "armRegionName", Operator: Equals, Value: "eastus"},
	}
	if err := testAPI(server).GetProductsPricePages(context.Background(), filters, func(*ProductsPricePage) {}); err != nil {
		t.Fatalf("GetProductsPricePages: %v", err)
	}
	if !strings.Contains(query, "api-version="+apiVersion) {
		t.Errorf("expected the API version in the query, got %q", query)
	}
	want := url.QueryEscape("serviceName eq 'Virtual Machines' and armRegionName eq 'eastus'")
	if !strings.Contains(query, want) {
		t.Errorf("expected the escaped filter %q in the query, got %q", want, query)
	}
}

func TestRetryAfterFromHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("x-ms-ratelimit-retry-after", "42")
	if got, ok := retryAfterFromHeaders(headers); !ok || got != 42*time.Second {
		t.Errorf("expected 42s, got %v (ok=%v)", got, ok)
	}

	headers = http.Header{}
	headers.Set("Retry-After", "7")
	if got, ok := retryAfterFromHeaders(headers); !ok || got != 7*time.Second {
		t.Errorf("expected 7s, got %v (ok=%v)", got, ok)
	}

	// a zero value is a valid answer ("retry immediately") and must be honoured
	headers = http.Header{}
	headers.Set("x-ms-ratelimit-retry-after", "0")
	if got, ok := retryAfterFromHeaders(headers); !ok || got != 0 {
		t.Errorf("expected an immediate retry, got %v (ok=%v)", got, ok)
	}

	// Retry-After may be an HTTP date rather than a number of seconds
	headers = http.Header{}
	headers.Set("Retry-After", time.Now().Add(time.Minute).UTC().Format(http.TimeFormat))
	if got, ok := retryAfterFromHeaders(headers); !ok || got < 50*time.Second || got > time.Minute {
		t.Errorf("expected about a minute, got %v (ok=%v)", got, ok)
	}

	if got, ok := retryAfterFromHeaders(http.Header{}); ok || got != 0 {
		t.Errorf("expected no wait without headers, got %v (ok=%v)", got, ok)
	}
}

func TestThrottleFromHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("x-ms-ratelimit-remaining-retailprices-requests", "3")
	if got := throttleFromHeaders(headers); got != 0 {
		t.Errorf("expected no throttle while budget remains, got %v", got)
	}

	headers.Set("x-ms-ratelimit-remaining-retailprices-requests", "0")
	headers.Set("x-ms-ratelimit-retry-after", "15")
	if got := throttleFromHeaders(headers); got != 15*time.Second {
		t.Errorf("expected to wait for the retry-after, got %v", got)
	}

	if got := throttleFromHeaders(http.Header{}); got != 0 {
		t.Errorf("expected no throttle without headers, got %v", got)
	}
}

func TestSleepContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sleepContext(ctx, time.Hour); err == nil {
		t.Error("expected a cancelled context to interrupt the wait")
	}
	if err := sleepContext(context.Background(), 0); err != nil {
		t.Errorf("expected no error for a zero wait, got %v", err)
	}
}
