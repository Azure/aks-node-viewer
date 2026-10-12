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
package model_test

import (
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Azure/aks-node-viewer/pkg/model"
)

func testNode(name string) *v1.Node {
	n := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Status: v1.NodeStatus{
			Phase: v1.NodePending,
		},
	}
	return n
}

func TestNewNode(t *testing.T) {
	n := testNode("mynode")
	node := model.NewNode(n)
	if exp, got := "mynode", node.Name(); exp != got {
		t.Errorf("expeted Name == %s, got %s", exp, got)
	}
}

func TestNodeTypeUnknown(t *testing.T) {
	n := testNode("mynode")
	node := model.NewNode(n)
	if node.IsOnDemand() {
		t.Errorf("exepcted to not be on-demand")
	}
	if node.IsSpot() {
		t.Errorf("exepcted to not be spot")
	}
}

func TestNodeTypeOnDemand(t *testing.T) {
	for label, value := range map[string]string{
		"karpenter.sh/capacity-type": "on-demand",
	} {
		n := testNode("mynode")
		n.Labels = map[string]string{
			label: value,
		}
		node := model.NewNode(n)
		if !node.IsOnDemand() {
			t.Errorf("exepcted on-demand")
		}
		if node.IsSpot() {
			t.Errorf("exepcted to not be spot")
		}
	}
}

func TestNodeTypeSpot(t *testing.T) {
	for label, value := range map[string]string{
		"karpenter.sh/capacity-type": "spot",
	} {
		n := testNode("mynode")
		n.Labels = map[string]string{
			label: value,
		}
		node := model.NewNode(n)
		if node.IsOnDemand() {
			t.Errorf("exepcted to not be on-demand")
		}
		if !node.IsSpot() {
			t.Errorf("exepcted to be spot")
		}
	}
}

func TestNodeRegion(t *testing.T) {
	n := testNode("mynode")
	if got := model.NewNode(n).Region(); got != "" {
		t.Errorf("expected no region without the label, got %q", got)
	}

	n.Labels = map[string]string{"topology.kubernetes.io/region": "eastus"}
	if got := model.NewNode(n).Region(); got != "eastus" {
		t.Errorf("expected the node region, got %q", got)
	}
}

func TestComputeLabel(t *testing.T) {
	n := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "mynode",
			CreationTimestamp: metav1.Time{Time: time.Now().Add(-time.Hour)},
		},
		Status: v1.NodeStatus{
			Allocatable: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("2"),
				v1.ResourceMemory: resource.MustParse("4Gi"),
				v1.ResourcePods:   resource.MustParse("110"),
			},
		},
	}
	node := model.NewNode(n)

	if age := node.ComputeLabel("eks-node-viewer/node-age"); age == "-" || age == "" {
		t.Errorf("expected eks-node-viewer/node-age to be a human duration, got %q", age)
	}
	if age := node.ComputeLabel("aks-node-viewer/node-age"); age == "-" || age == "" {
		t.Errorf("expected aks-node-viewer/node-age to be a human duration, got %q", age)
	}
	if got := node.ComputeLabel("eks-node-viewer/node-cpu-usage"); got != "0%" {
		t.Errorf("expected node-cpu-usage to be 0%%, got %q", got)
	}
	if got := node.ComputeLabel("eks-node-viewer/node-nope-usage"); got != "N/A" {
		t.Errorf("expected unknown resource usage to be N/A, got %q", got)
	}
	if got := node.ComputeLabel("some-other-label"); got != "-" {
		t.Errorf("expected unknown label to be '-', got %q", got)
	}
}

func TestComputeLabelResourceUsage(t *testing.T) {
	node := model.NewNode(&v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "mynode"},
		Status: v1.NodeStatus{
			Allocatable: v1.ResourceList{
				v1.ResourceCPU: resource.MustParse("2"),
			},
		},
	})
	pod := model.NewPod(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "mypod"},
		Spec: v1.PodSpec{
			NodeName: "mynode",
			Containers: []v1.Container{
				{
					Name: "c",
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("500m")},
					},
				},
			},
		},
	})
	node.BindPod(pod)
	if got := node.ComputeLabel("eks-node-viewer/node-cpu-usage"); got != "25%" {
		t.Errorf("expected node-cpu-usage to be 25%%, got %q", got)
	}
	if got := node.ComputeLabel("aks-node-viewer/node-memory-usage"); got != "N/A" {
		t.Errorf("expected aks-node-viewer/node-memory-usage to be N/A, got %q", got)
	}
}
