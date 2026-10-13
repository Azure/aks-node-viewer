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

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFormatCPU(t *testing.T) {
	cases := []struct{ in, want string }{
		{"0", "0"},
		{"250m", "250m"},
		{"500m", "500m"},
		{"1", "1"},
		{"8", "8"},
		{"1500m", "1.5"},
		{"4250m", "4.25"},
		{"1333m", "1.333"},
		{"4250000m", "4,250"},
	}
	for _, c := range cases {
		if got := formatResource(v1.ResourceCPU, resource.MustParse(c.in)); got != c.want {
			t.Errorf("formatResource(cpu, %q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct{ in, want string }{
		{"100", "100"},
		{"512Mi", "512Mi"},
		{"1024Ki", "1Mi"},
		{"12Gi", "12Gi"},
		{"16384Mi", "16Gi"},
		{"12800Mi", "12.5Gi"},
		{"2Ti", "2Ti"},
	}
	for _, c := range cases {
		if got := formatResource(v1.ResourceMemory, resource.MustParse(c.in)); got != c.want {
			t.Errorf("formatResource(memory, %q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatCount(t *testing.T) {
	if got := formatCount(1); got != "1" {
		t.Errorf("expected 1, got %q", got)
	}
	if got := formatCount(1234); got != "1,234" {
		t.Errorf("expected 1,234, got %q", got)
	}
	if got := formatCount(1234567); got != "1,234,567" {
		t.Errorf("expected 1,234,567, got %q", got)
	}
}

func TestGroupThousands(t *testing.T) {
	cases := []struct{ in, want string }{
		{"0", "0"},
		{"999", "999"},
		{"1000", "1,000"},
		{"1234567", "1,234,567"},
		{"1234.5", "1,234.5"},
		{"-1234", "-1,234"},
	}
	for _, c := range cases {
		if got := groupThousands(c.in); got != c.want {
			t.Errorf("groupThousands(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFormatResourceRouting(t *testing.T) {
	if got := formatResource(v1.ResourceEphemeralStorage, resource.MustParse("20Gi")); got != "20Gi" {
		t.Errorf("expected ephemeral-storage to use binary units, got %q", got)
	}
	if got := formatResource(v1.ResourceStorage, resource.MustParse("1Ti")); got != "1Ti" {
		t.Errorf("expected storage to use binary units, got %q", got)
	}
	// unknown resources keep their canonical Kubernetes representation
	if got := formatResource(v1.ResourcePods, resource.MustParse("110")); got != "110" {
		t.Errorf("expected pods to keep its canonical form, got %q", got)
	}
}

func TestClusterSummaryFormatsResources(t *testing.T) {
	style, err := ParseStyle("#04B575,#FFFF00,#FF0000")
	if err != nil {
		t.Fatalf("ParseStyle: %v", err)
	}
	u := NewUIModel(nil, "creation=dsc", style)
	u.SetResources([]string{"cpu", "memory"})

	node := NewNode(&v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node-1"},
		Status: v1.NodeStatus{
			Allocatable: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("4"),
				v1.ResourceMemory: resource.MustParse("16Gi"),
				v1.ResourcePods:   resource.MustParse("110"),
			},
		},
	})
	node.Show()
	u.Cluster().AddNode(node)

	pod := NewPod(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "pod-1"},
		Spec: v1.PodSpec{
			NodeName: "node-1",
			Containers: []v1.Container{{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("1500m"),
						v1.ResourceMemory: resource.MustParse("12Gi"),
					},
				},
			}},
		},
	})
	u.Cluster().AddPod(pod, nil)

	out := u.View()
	if !strings.Contains(out, "1.5/4") {
		t.Errorf("expected cpu to be displayed in cores, got:\n%s", out)
	}
	if !strings.Contains(out, "12Gi/16Gi") {
		t.Errorf("expected memory to be displayed in Gi, got:\n%s", out)
	}
	if strings.Contains(out, "Ki") {
		t.Errorf("expected memory not to be displayed in Ki, got:\n%s", out)
	}
	if strings.Contains(out, "1500m") {
		t.Errorf("expected cpu not to be displayed in millicores, got:\n%s", out)
	}
}
