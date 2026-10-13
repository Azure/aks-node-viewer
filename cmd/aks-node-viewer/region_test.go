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

import "testing"

func TestRegionTracker(t *testing.T) {
	var rt regionTracker

	if got := rt.region(); got != "" {
		t.Errorf("expected no region initially, got %q", got)
	}
	if got := rt.summary(); got != "" {
		t.Errorf("expected no summary initially, got %q", got)
	}

	// an empty region (node without the label) is ignored
	rt.observe("")
	if got := rt.summary(); got != "" {
		t.Errorf("expected an empty region to be ignored, got %q", got)
	}

	rt.observe("eastus")
	if got := rt.region(); got != "eastus" {
		t.Errorf("expected eastus, got %q", got)
	}

	// duplicates are not recorded twice
	rt.observe("eastus")
	if got := rt.summary(); got != "eastus" {
		t.Errorf("expected a single region, got %q", got)
	}

	// the first region wins, but a discrepancy stays visible in the summary
	rt.observe("westus2")
	if got := rt.region(); got != "eastus" {
		t.Errorf("expected the first region to win, got %q", got)
	}
	if got := rt.summary(); got != "eastus, westus2" {
		t.Errorf("expected both regions in the summary, got %q", got)
	}
}
