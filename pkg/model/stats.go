// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package model

import v1 "k8s.io/api/core/v1"

type Stats struct {
	NumNodes             int
	AllocatableResources v1.ResourceList
	UsedResources        v1.ResourceList
	PercentUsedResoruces map[v1.ResourceName]float64
	Nodes                []*Node
	TotalPods            int
	PodsByPhase          map[v1.PodPhase]int
	BoundPodCount        int
	TotalPrice           float64
}
