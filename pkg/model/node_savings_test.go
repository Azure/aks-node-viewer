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
	"math"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func nodeWithCapacity(name string) *Node {
	return NewNode(&v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Status: v1.NodeStatus{
			Allocatable: v1.ResourceList{
				v1.ResourceCPU: resource.MustParse("2"),
			},
		},
	})
}

func setCapacity(n *Node, capacity string) {
	n.node.Labels = map[string]string{
		"karpenter.sh/capacity-type": capacity,
		v1.LabelInstanceTypeStable:   "Standard_D2s_v3",
	}
}

func TestSpotSavings(t *testing.T) {
	spot := nodeWithCapacity("spot-1")
	setCapacity(spot, "spot")
	spot.price = 0.050
	spot.onDemandPrice = 0.100

	if got, want := spot.SpotSavings(), 0.050; math.Abs(got-want) > 1e-9 {
		t.Errorf("expected savings %v, got %v", want, got)
	}
	for _, label := range []string{"aks-node-viewer/node-spot-savings", "eks-node-viewer/node-spot-savings"} {
		if got, want := spot.ComputeLabel(label), "$0.050"; got != want {
			t.Errorf("expected %s to be %q, got %q", label, want, got)
		}
	}

	// on-demand nodes have no spot savings
	onDemand := nodeWithCapacity("od-1")
	setCapacity(onDemand, "on-demand")
	onDemand.price = 0.100
	onDemand.onDemandPrice = 0.100
	if got := onDemand.SpotSavings(); got != 0 {
		t.Errorf("expected no savings for an on-demand node, got %v", got)
	}
	if got, want := onDemand.ComputeLabel("aks-node-viewer/node-spot-savings"), "-"; got != want {
		t.Errorf("expected %q for an on-demand node, got %q", want, got)
	}

	// unknown prices produce no savings
	unknown := nodeWithCapacity("spot-2")
	setCapacity(unknown, "spot")
	unknown.price = math.NaN()
	unknown.onDemandPrice = math.NaN()
	if got := unknown.SpotSavings(); got != 0 {
		t.Errorf("expected no savings for unknown prices, got %v", got)
	}
	if got, want := unknown.ComputeLabel("eks-node-viewer/node-spot-savings"), "-"; got != want {
		t.Errorf("expected %q for unknown prices, got %q", want, got)
	}
}

func TestClusterSpotSavings(t *testing.T) {
	style, err := ParseStyle("#04B575,#FFFF00,#FF0000")
	if err != nil {
		t.Fatalf("ParseStyle: %v", err)
	}
	u := NewUIModel(nil, "creation=dsc", style)
	u.SetResources([]string{"cpu"})

	spot := nodeWithCapacity("spot-1")
	setCapacity(spot, "spot")
	spot.price = 0.050
	spot.onDemandPrice = 0.100
	spot.Show()
	u.Cluster().AddNode(spot)

	onDemand := nodeWithCapacity("od-1")
	setCapacity(onDemand, "on-demand")
	onDemand.price = 0.100
	onDemand.onDemandPrice = 0.100
	onDemand.Show()
	u.Cluster().AddNode(onDemand)

	stats := u.Cluster().Stats()
	if got, want := stats.TotalPrice, 0.150; math.Abs(got-want) > 1e-9 {
		t.Errorf("expected total price %v, got %v", want, got)
	}
	if got, want := stats.SpotSavings, 0.050; math.Abs(got-want) > 1e-9 {
		t.Errorf("expected spot savings %v, got %v", want, got)
	}

	// baseline is the on-demand equivalent cost: 0.150 + 0.050 = 0.200, so the
	// savings represent 25% of it.
	out := u.View()
	if !strings.Contains(out, "saved $0.050/hour (25%)") {
		t.Errorf("expected the cluster savings in the summary, got:\n%s", out)
	}
}

func TestClusterNoSavingsWithoutSpot(t *testing.T) {
	style, err := ParseStyle("#04B575,#FFFF00,#FF0000")
	if err != nil {
		t.Fatalf("ParseStyle: %v", err)
	}
	u := NewUIModel(nil, "creation=dsc", style)
	u.SetResources([]string{"cpu"})

	onDemand := nodeWithCapacity("od-1")
	setCapacity(onDemand, "on-demand")
	onDemand.price = 0.100
	onDemand.onDemandPrice = 0.100
	onDemand.Show()
	u.Cluster().AddNode(onDemand)

	if out := u.View(); strings.Contains(out, "saved") {
		t.Errorf("expected no savings note without spot nodes, got:\n%s", out)
	}
}
