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
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
