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
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	apiVersion  = "2021-10-01-preview"
	pricingPath = "/api/retail/prices?api-version=" + apiVersion

	// maxAttempts bounds how often a single page is retried when the API rate
	// limits us or returns a transient error. A full paginated fetch can take a
	// while (the API answers a rate limited request with 429 and asks us to wait
	// about a minute), but it must still finish rather than fail part way
	// through.
	maxAttempts = 5
	// defaultRetryAfter is used when the API rate limits us without telling us
	// how long to wait.
	defaultRetryAfter = 30 * time.Second

	userAgent = "aks-node-viewer"
)

// errRateLimited marks a response the API rejected because we exceeded its rate
// limit. It is retried after the interval the API asks us to wait.
var errRateLimited = errors.New("rate limited by the Azure retail prices API")

// statusError is an unexpected (non-rate-limit) HTTP response.
type statusError struct {
	code int
	body string
}

func (e *statusError) Error() string {
	msg := fmt.Sprintf("unexpected status code %d", e.code)
	if e.body != "" {
		msg += ": " + e.body
	}
	return msg
}

type PricingAPI interface {
	GetProductsPricePages(context.Context, []*Filter, func(output *ProductsPricePage)) error
}

type pricingAPI struct {
	baseURL string
	client  *http.Client
}

func New() PricingAPI {
	return &pricingAPI{
		baseURL: "https://prices.azure.com",
		client: &http.Client{
			// Pages of VM prices are large, so allow plenty of time for each
			// request without hanging forever on a stalled connection.
			Timeout: 60 * time.Second,
		},
	}
}

func (papi *pricingAPI) GetProductsPricePages(ctx context.Context, filters []*Filter, pageHandler func(output *ProductsPricePage)) error {
	nextURL := papi.baseURL + pricingPath

	if len(filters) > 0 {
		filterParams := []string{}
		for _, filter := range filters {
			filterParams = append(filterParams, filter.String())
		}

		filterParamsEscaped := url.QueryEscape(strings.Join(filterParams[:], " and "))

		nextURL += fmt.Sprintf("&$filter=%s", filterParamsEscaped)
	}

	for nextURL != "" {
		page, throttle, err := papi.fetchPage(ctx, nextURL)
		if err != nil {
			return err
		}

		pageHandler(page)
		nextURL = page.NextPageLink

		// The API tells us how much of our request budget is left; if it is
		// exhausted, wait before asking for the next page so we do not get rate
		// limited mid-fetch.
		if nextURL != "" {
			if err := sleepContext(ctx, throttle); err != nil {
				return err
			}
		}
	}

	return nil
}

// fetchPage retrieves a single page, retrying while the API rate limits us or
// returns a transient server error. The retail prices API answers a rate limited
// request with 429 and an x-ms-ratelimit-retry-after header that says how many
// seconds to wait; honouring it lets a full paginated fetch complete instead of
// failing part way through.
func (papi *pricingAPI) fetchPage(ctx context.Context, url string) (*ProductsPricePage, time.Duration, error) {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		page, throttle, err := papi.requestPage(ctx, url)
		if err == nil {
			return page, throttle, nil
		}
		lastErr = err

		var wait time.Duration
		switch {
		case errors.Is(err, errRateLimited):
			// Wait for as long as the API asked us to; this is what lets a
			// paginated fetch outlast the request budget and still complete.
			wait = throttle
		case isRetryable(err):
			wait = backoff(attempt + 1)
		default:
			return nil, 0, err
		}

		if err := sleepContext(ctx, wait); err != nil {
			return nil, 0, err
		}
	}
	return nil, 0, lastErr
}

// requestPage performs a single HTTP request. On a rate limit it returns
// errRateLimited along with the interval the API wants us to wait before
// retrying. For any other response it returns the page and how long to wait
// before the next request (zero unless our request budget is exhausted).
func (papi *pricingAPI) requestPage(ctx context.Context, url string) (*ProductsPricePage, time.Duration, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := papi.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		page := &ProductsPricePage{}
		if err := json.NewDecoder(res.Body).Decode(page); err != nil {
			return nil, 0, err
		}
		return page, throttleFromHeaders(res.Header), nil
	case http.StatusTooManyRequests:
		if wait, ok := retryAfterFromHeaders(res.Header); ok {
			return nil, wait, errRateLimited
		}
		return nil, defaultRetryAfter, errRateLimited
	default:
		body, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return nil, 0, &statusError{code: res.StatusCode, body: strings.TrimSpace(string(body))}
	}
}

// isRetryable reports whether a failed request is worth retrying. Network
// failures and 5xx responses are transient, everything else is not.
func isRetryable(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	var statusErr *statusError
	if errors.As(err, &statusErr) {
		return statusErr.code >= 500
	}
	return false
}

// retryAfterFromHeaders returns how long the API wants us to wait before
// retrying, and whether it said so. The retail prices API uses
// x-ms-ratelimit-retry-after; Retry-After is the standard header and is used as
// a fallback.
func retryAfterFromHeaders(h http.Header) (time.Duration, bool) {
	for _, name := range []string{"x-ms-ratelimit-retry-after", "Retry-After"} {
		value := strings.TrimSpace(h.Get(name))
		if value == "" {
			continue
		}
		// The headers carry a number of seconds. Retry-After may also be an
		// HTTP date, so accept that too.
		if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
			return time.Duration(seconds) * time.Second, true
		}
		if when, err := http.ParseTime(value); err == nil {
			if wait := time.Until(when); wait > 0 {
				return wait, true
			}
			return 0, true
		}
	}
	return 0, false
}

// throttleFromHeaders returns how long to wait before the next request when the
// API reports that our request budget is exhausted, and 0 otherwise.
func throttleFromHeaders(h http.Header) time.Duration {
	remaining := strings.TrimSpace(h.Get("x-ms-ratelimit-remaining-retailprices-requests"))
	if remaining == "" {
		return 0
	}
	left, err := strconv.Atoi(remaining)
	if err != nil || left > 0 {
		return 0
	}
	if wait, ok := retryAfterFromHeaders(h); ok {
		return wait
	}
	return defaultRetryAfter
}

// backoff returns how long to wait before the given retry attempt.
func backoff(attempt int) time.Duration {
	wait := time.Duration(attempt) * time.Second
	if wait > defaultRetryAfter {
		wait = defaultRetryAfter
	}
	return wait
}

// sleepContext waits for d, or until the context is done.
func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
