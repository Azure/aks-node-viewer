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
	"fmt"
	"math"
	"regexp"
	"sync"
	"time"

	"github.com/Azure/aks-node-viewer/pkg/pricing"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

var (
	// nodeLabelPrefix matches both the aks-node-viewer and eks-node-viewer label prefixes so
	// that configs written for either tool keep working.
	nodeLabelPrefix     = "(?:aks|eks)-node-viewer"
	resourceUsageLabels = regexp.MustCompile(nodeLabelPrefix + `/node-(.*?)-usage`)
)

type objectKey struct {
	namespace string
	name      string
}
type Node struct {
	mu      sync.RWMutex
	visible bool
	node    v1.Node
	pods    map[objectKey]*Pod
	used    v1.ResourceList
	price   float64
	// onDemandPrice is the on-demand price for the node's instance type,
	// regardless of the capacity type actually used. It is math.NaN() when the
	// on-demand price is unknown and is used to compute spot savings.
	onDemandPrice float64
}

func NewNode(n *v1.Node) *Node {
	node := &Node{
		node: *n,
		pods: map[objectKey]*Pod{},
		used: v1.ResourceList{},
	}

	return node
}
func (n *Node) IsOnDemand() bool {
	return n.node.Labels["karpenter.sh/capacity-type"] == "on-demand"
}

func (n *Node) IsSpot() bool {
	return n.node.Labels["karpenter.sh/capacity-type"] == "spot"
}

func (n *Node) Update(node *v1.Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.node = *node
}

func (n *Node) Name() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.node.Name
}

func (n *Node) BindPod(pod *Pod) {
	n.mu.Lock()
	defer n.mu.Unlock()
	key := objectKey{
		namespace: pod.Namespace(),
		name:      pod.Name(),
	}
	_, alreadyBound := n.pods[key]
	n.pods[key] = pod

	if !alreadyBound {
		for rn, q := range pod.Requested() {
			existing := n.used[rn]
			existing.Add(q)
			n.used[rn] = existing
		}
	}
}

func (n *Node) DeletePod(namespace string, name string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	key := objectKey{namespace: namespace, name: name}
	if p, ok := n.pods[key]; ok {
		// subtract the pod requests
		for rn, q := range p.Requested() {
			existing := n.used[rn]
			existing.Sub(q)
			n.used[rn] = existing
		}
		delete(n.pods, key)
	}
}

func (n *Node) Allocatable() v1.ResourceList {
	n.mu.RLock()
	defer n.mu.RUnlock()
	// shouldn't be modified so it's safe to return
	return n.node.Status.Allocatable
}

func (n *Node) Used() v1.ResourceList {
	n.mu.RLock()
	defer n.mu.RUnlock()
	used := v1.ResourceList{}
	for rn, q := range n.used {
		used[rn] = q.DeepCopy()
	}
	return used
}

func (n *Node) Cordoned() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.node.Spec.Unschedulable {
		return true
	}
	for _, taint := range n.node.Spec.Taints {
		if taint.Key == "karpenter.sh/disruption" && taint.Effect == v1.TaintEffectNoSchedule {
			return true
		}
	}
	return false
}

func (n *Node) Ready() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for _, c := range n.node.Status.Conditions {
		if c.Status == v1.ConditionTrue && c.Type == v1.NodeReady {
			return true
		}
	}
	return false
}

func (n *Node) Created() time.Time {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.node.CreationTimestamp.Time
}

func (n *Node) InstanceType() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.node.Labels[v1.LabelInstanceTypeStable]
}

func (n *Node) Zone() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if v := n.node.Labels[v1.LabelTopologyZone]; v != "" {
		return v
	}
	return n.node.Labels["failure-domain.beta.kubernetes.io/zone"]
}

// Region returns the Azure region the node runs in, from the standard
// topology.kubernetes.io/region label, or "" when the label is not set.
func (n *Node) Region() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if v := n.node.Labels[v1.LabelTopologyRegion]; v != "" {
		return v
	}
	return n.node.Labels["failure-domain.beta.kubernetes.io/region"]
}

func (n *Node) NumPods() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.pods)
}

func (n *Node) Hide() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.visible = false
}
func (n *Node) Visible() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.visible
}

func (n *Node) Show() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.visible = true
}

func (n *Node) Deleting() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return !n.node.DeletionTimestamp.IsZero()
}

func (n *Node) HasPrice() bool {
	// we use NaN for an unknown price, so if this is true the price is known
	return n.Price() == n.Price()
}

// SpotSavings returns the hourly saving of running this node on spot capacity
// compared to on-demand. It is zero unless the node is a spot node and both the
// spot and on-demand prices are known and the spot price is lower.
func (n *Node) SpotSavings() float64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if !n.isSpotLocked() || math.IsNaN(n.price) || math.IsNaN(n.onDemandPrice) {
		return 0
	}
	if savings := n.onDemandPrice - n.price; savings > 0 {
		return savings
	}
	return 0
}

// SpotSavingsReference returns the on-demand price of a spot node (rendered for
// display) and the percentage saved by running it on spot instead. The
// on-demand price is only a reference: it is not what the node costs. ok is
// false for nodes that are not on spot or whose prices are unknown.
func (n *Node) SpotSavingsReference() (reference string, pct int, ok bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if !n.isSpotLocked() || math.IsNaN(n.price) || math.IsNaN(n.onDemandPrice) || n.onDemandPrice <= 0 {
		return "", 0, false
	}
	pct = int(math.Round(100 * (n.onDemandPrice - n.price) / n.onDemandPrice))
	return fmt.Sprintf("$%0.3f", n.onDemandPrice), pct, true
}

// Price returns the last known price for the node. math.NaN() is used to
// indicate that the price is unknown.
func (n *Node) Price() float64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.price
}

func (n *Node) UpdatePrice(pricing *pricing.Provider) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.lookupPriceLocked(pricing)
}

// lookupPriceLocked must be called with n.mu held
func (n *Node) lookupPriceLocked(pprov *pricing.Provider) {
	// lookup our n price
	n.price = math.NaN()
	n.onDemandPrice = math.NaN()
	instanceType := n.node.Labels[v1.LabelInstanceTypeStable]

	// always resolve the on-demand price so spot savings can be computed
	if price, ok := pprov.OnDemandPrice(instanceType); ok {
		n.onDemandPrice = price
	}

	if n.isOnDemandLocked() {
		n.price = n.onDemandPrice
	} else if n.isSpotLocked() {
		if price, ok := pprov.SpotPrice(instanceType); ok {
			n.price = price
		}
		// Use on-demand prices if no Karpenter label is found
	} else {
		n.price = n.onDemandPrice
	}
}

// isOnDemandLocked and isSpotLocked must be called with n.mu held
func (n *Node) isOnDemandLocked() bool {
	return n.node.Labels["karpenter.sh/capacity-type"] == "on-demand"
}

func (n *Node) isSpotLocked() bool {
	return n.node.Labels["karpenter.sh/capacity-type"] == "spot"
}

// ComputeLabel computes dynamic labels
func (n *Node) ComputeLabel(labelName string) string {
	switch labelName {
	case "aks-node-viewer/node-age", "eks-node-viewer/node-age":
		return duration.HumanDuration(time.Since(n.Created()))
	case "aks-node-viewer/node-spot-savings", "eks-node-viewer/node-spot-savings":
		if reference, pct, ok := n.SpotSavingsReference(); ok {
			return fmt.Sprintf("%s %d%%", reference, pct)
		}
		return "-"
	}
	// resource based custom labels, e.g. eks-node-viewer/node-cpu-usage
	if match := resourceUsageLabels.FindStringSubmatch(labelName); len(match) > 0 {
		return pctUsage(n.Allocatable(), n.Used(), match[1])
	}
	return "-"
}

// NotReadyTime is the time that the node went NotReady, or when it was created if it hasn't been marked as NotReady.
func (n *Node) NotReadyTime() time.Time {
	n.mu.RLock()
	var notReadyTransitionTime time.Time
	for _, c := range n.node.Status.Conditions {
		if c.Type == v1.NodeReady && (c.Status == v1.ConditionFalse || c.Status == v1.ConditionUnknown) {
			notReadyTransitionTime = c.LastTransitionTime.Time
			break
		}
	}
	n.mu.RUnlock()
	if !notReadyTransitionTime.IsZero() {
		return notReadyTransitionTime
	}
	return n.Created()
}

func pctUsage(allocatable v1.ResourceList, used v1.ResourceList, resource string) string {
	allocRes, hasAlloc := allocatable[v1.ResourceName(resource)]
	if !hasAlloc {
		return "N/A"
	}
	usedRes, hasUsed := used[v1.ResourceName(resource)]
	if !hasUsed || usedRes.AsApproximateFloat64() == 0 {
		return "0%"
	}
	pctUsed := 0.0
	if allocRes.AsApproximateFloat64() != 0 {
		pctUsed = 100 * (usedRes.AsApproximateFloat64() / allocRes.AsApproximateFloat64())
	}
	return fmt.Sprintf("%.0f%%", pctUsed)
}
