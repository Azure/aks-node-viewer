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

package model

import (
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Azure/aks-node-viewer/pkg/pricing"
)

func nodeWithCreation(name string, created time.Time) *Node {
	return NewNode(&v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: name, CreationTimestamp: metav1.Time{Time: created}},
	})
}

func TestMakeNodeSorterCreation(t *testing.T) {
	older := nodeWithCreation("old", time.Now().Add(-2*time.Hour))
	newer := nodeWithCreation("new", time.Now().Add(-time.Hour))

	// matches eks-node-viewer semantics: creation=asc sorts newest first,
	// creation=dsc (the default) sorts oldest first.
	asc := makeNodeSorter("creation=asc")
	if !asc(newer, older) {
		t.Errorf("expected 'creation=asc' to sort newer before older")
	}
	if asc(older, newer) {
		t.Errorf("expected 'creation=asc' to not sort older before newer")
	}

	dsc := makeNodeSorter("creation=dsc")
	if !dsc(older, newer) {
		t.Errorf("expected 'creation=dsc' to sort older before newer")
	}
	if dsc(newer, older) {
		t.Errorf("expected 'creation=dsc' to not sort newer before older")
	}
}

func TestNodeRowPriceIsSeparateColumn(t *testing.T) {
	style, err := ParseStyle("#04B575,#FFFF00,#FF0000")
	if err != nil {
		t.Fatalf("ParseStyle: %v", err)
	}
	u := NewUIModel(nil, "creation=dsc", style)
	u.SetResources([]string{"cpu"})

	n := NewNode(&v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node-1"},
		Status: v1.NodeStatus{
			Allocatable: v1.ResourceList{
				v1.ResourceCPU:  resource.MustParse("2"),
				v1.ResourcePods: resource.MustParse("110"),
			},
		},
	})
	n.node.Labels = map[string]string{v1.LabelInstanceTypeStable: "Standard_D2s_v3"}
	n.price = 0.096
	n.Show()
	u.Cluster().AddNode(n)

	out := u.View()

	// the price must be its own column, separated from the instance type by
	// padding rather than glued on with a slash
	if !strings.Contains(out, "Standard_D2s_v3 $0.096") {
		t.Errorf("expected the price to be its own column, got:\n%s", out)
	}
	if strings.Contains(out, "Standard_D2s_v3/") {
		t.Errorf("expected the price not to be glued to the instance type, got:\n%s", out)
	}
}

func TestPricingNote(t *testing.T) {
	style, err := ParseStyle("#04B575,#FFFF00,#FF0000")
	if err != nil {
		t.Fatalf("ParseStyle: %v", err)
	}

	newModel := func(st pricing.Status) *UIModel {
		u := NewUIModel(nil, "creation=dsc", style)
		u.PricingStatus = func() pricing.Status { return st }
		return u
	}

	// live prices show when they were updated
	live := time.Date(2026, 10, 12, 15, 4, 0, 0, time.Local)
	if got := newModel(pricing.Status{Source: pricing.SourceLive, Updated: live}).pricingNote(); !strings.Contains(got, "prices updated 2026-10-12 15:04") {
		t.Errorf("expected a live note, got %q", got)
	}

	// cached prices show the cache age
	cached := time.Now().Add(-2*time.Hour - time.Minute)
	if got := newModel(pricing.Status{Source: pricing.SourceCache, Updated: cached}).pricingNote(); !strings.Contains(got, "prices cached") || !strings.Contains(got, "2h ago") {
		t.Errorf("expected a cached note with its age, got %q", got)
	}

	// static prices are called out as not updated
	static := time.Date(2023, 3, 1, 19, 26, 21, 0, time.UTC)
	if got := newModel(pricing.Status{Source: pricing.SourceStatic, Updated: static}).pricingNote(); !strings.Contains(got, "prices not updated") || !strings.Contains(got, "2023-03-01") {
		t.Errorf("expected a not-updated note, got %q", got)
	}

	// an in-flight refresh is surfaced
	refresh := newModel(pricing.Status{Source: pricing.SourceStatic, Updated: static, Updating: true}).pricingNote()
	if !strings.Contains(refresh, "refreshing") {
		t.Errorf("expected a refreshing note, got %q", refresh)
	}

	// without a status callback there is nothing to show
	u := NewUIModel(nil, "creation=dsc", style)
	if got := u.pricingNote(); got != "" {
		t.Errorf("expected an empty note without PricingStatus, got %q", got)
	}
}

func TestMakeNodeSorterByLabel(t *testing.T) {
	az1 := nodeWithCreation("az1-node", time.Now().Add(-time.Hour))
	az2 := nodeWithCreation("az2-node", time.Now().Add(-time.Hour))
	az1.node.Labels = map[string]string{"topology.kubernetes.io/zone": "a"}
	az2.node.Labels = map[string]string{"topology.kubernetes.io/zone": "b"}

	asc := makeNodeSorter("topology.kubernetes.io/zone=asc")
	if !asc(az1, az2) {
		t.Errorf("expected 'a' zone node before 'b' zone node")
	}
	dsc := makeNodeSorter("topology.kubernetes.io/zone=dsc")
	if !dsc(az2, az1) {
		t.Errorf("expected 'b' zone node before 'a' zone node with dsc")
	}
}
