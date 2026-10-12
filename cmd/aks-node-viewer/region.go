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

package main

import (
	"strings"
	"sync"
)

// regionTracker records the Azure regions reported by the cluster's nodes. AKS
// clusters are single-region, so in practice all nodes agree; if they don't, the
// first region observed is used for pricing and the discrepancy is surfaced in
// the UI.
type regionTracker struct {
	mu      sync.Mutex
	regions []string
}

// observe records a region reported by a node, ignoring the empty value used
// when a node has no region label.
func (t *regionTracker) observe(region string) {
	if region == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, r := range t.regions {
		if r == region {
			return
		}
	}
	t.regions = append(t.regions, region)
}

// region returns the region prices should be looked up for: the first region
// reported by the nodes, or "" when none is known yet.
func (t *regionTracker) region() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.regions) == 0 {
		return ""
	}
	return t.regions[0]
}

// summary describes the regions reported by the nodes for display, or "" when
// none is known yet.
func (t *regionTracker) summary() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return strings.Join(t.regions, ", ")
}
