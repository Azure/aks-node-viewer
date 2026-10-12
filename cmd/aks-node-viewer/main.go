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
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/Azure/aks-node-viewer/pkg/client"
	"github.com/Azure/aks-node-viewer/pkg/model"
	"github.com/Azure/aks-node-viewer/pkg/pricing"
)

var attribution string

func main() {
	flags, err := ParseFlags()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		log.Fatalf("cannot parse flags: %v", err)
	}

	if flags.ShowAttribution {
		fmt.Println(attribution)
		os.Exit(0)
	}

	if flags.Version {
		fmt.Printf("aks-node-viewer version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
		fmt.Printf("built by: %s\n", builtBy)
		os.Exit(0)
	}

	cs, err := client.Create(flags.Kubeconfig, flags.Context)
	if err != nil {
		log.Fatalf("creating client, %s", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	style, err := model.ParseStyle(flags.Style)
	if err != nil {
		log.Fatalf("parsing style, %s", err)
	}
	m := model.NewUIModel(strings.Split(flags.ExtraLabels, ","), flags.NodeSort, style)
	m.DisablePricing = flags.DisablePricing
	m.SetResources(strings.FieldsFunc(flags.Resources, func(r rune) bool { return r == ',' }))

	// Prices are looked up for a single Azure region. Unless the region is set
	// explicitly, it is taken from the cluster's nodes so a cluster outside the
	// default region is priced correctly; an explicitly configured region always
	// wins, and a mismatch with the cluster is surfaced in the UI.
	regionExplicit := flags.Region != ""
	region := flags.Region
	if region == "" {
		region = pricing.DefaultRegion
	}
	clusterRegions := &regionTracker{}

	// The pricing provider starts with cached (or embedded static) prices so the UI
	// renders immediately; live prices are fetched in the background and applied as
	// they arrive. The status is surfaced in the UI so stale prices are obvious.
	pprov := pricing.NewProvider(ctx, pricing.NewAPI(), region)
	m.PricingStatus = pprov.Status
	m.ClusterRegion = clusterRegions.summary
	if !flags.DisablePricing {
		pprov.Refresh(ctx)
		go func() {
			// re-apply prices to nodes whenever fresh pricing data arrives
			for range pprov.Watch() {
				m.Cluster().UpdatePrices(pprov)
			}
		}()
	}

	// observeNodeRegion keeps the pricing region in sync with the cluster. When
	// the region is not set explicitly it is taken from the first node that
	// reports one; an explicit region is never overridden.
	observeNodeRegion := func(nodeRegion string) {
		if nodeRegion == "" {
			return
		}
		clusterRegions.observe(nodeRegion)
		if flags.DisablePricing || regionExplicit {
			return
		}
		if pprov.SetRegion(clusterRegions.region()) {
			pprov.Refresh(ctx)
		}
	}

	var nodeSelector labels.Selector
	if ns, err := labels.Parse(flags.NodeSelector); err != nil {
		log.Fatalf("parsing node selector: %s", err)
	} else {
		nodeSelector = ns
	}

	monitorSettings := &monitorSettings{
		clientset:         cs,
		model:             m,
		nodeSelector:      nodeSelector,
		pricing:           pprov,
		observeNodeRegion: observeNodeRegion,
	}
	startMonitor(ctx, monitorSettings)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("error running tea: %s", err)
	}
	cancel()
}

type monitorSettings struct {
	clientset         *kubernetes.Clientset
	model             *model.UIModel
	nodeSelector      labels.Selector
	pricing           *pricing.Provider
	observeNodeRegion func(region string)
}

func startMonitor(ctx context.Context, settings *monitorSettings) {
	if settings.observeNodeRegion == nil {
		settings.observeNodeRegion = func(string) {}
	}

	podWatchList := cache.NewListWatchFromClient(settings.clientset.CoreV1().RESTClient(), "pods",
		v1.NamespaceAll, fields.Everything())

	cluster := settings.model.Cluster()
	_, podController := cache.NewInformer(
		podWatchList,
		&v1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				p := obj.(*v1.Pod)
				if !isTerminalPod(p) {
					cluster.AddPod(model.NewPod(p), settings.pricing)
				}
			},
			DeleteFunc: func(obj interface{}) {
				p := obj.(*v1.Pod)
				cluster.DeletePod(p.Namespace, p.Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				p := newObj.(*v1.Pod)
				if isTerminalPod(p) {
					cluster.DeletePod(p.Namespace, p.Name)
				} else {
					pod, ok := cluster.GetPod(p.Namespace, p.Name)
					if !ok {
						cluster.AddPod(model.NewPod(p), settings.pricing)
					} else {
						pod.Update(p)
						cluster.AddPod(pod, settings.pricing)
					}
				}
			},
		},
	)
	go podController.Run(ctx.Done())

	nodeWatchList := cache.NewFilteredListWatchFromClient(settings.clientset.CoreV1().RESTClient(), "nodes",
		v1.NamespaceAll, func(options *metav1.ListOptions) {
			options.LabelSelector = settings.nodeSelector.String()
		})
	_, controller := cache.NewInformer(
		nodeWatchList,
		&v1.Node{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				node := model.NewNode(obj.(*v1.Node))
				settings.observeNodeRegion(node.Region())
				node.UpdatePrice(settings.pricing)
				n := cluster.AddNode(node)
				n.Show()
			},
			DeleteFunc: func(obj interface{}) {
				cluster.DeleteNode(obj.(*v1.Node).Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				n := newObj.(*v1.Node)
				if !n.DeletionTimestamp.IsZero() && len(n.Finalizers) == 0 {
					cluster.DeleteNode(n.Name)
				} else {
					node, ok := cluster.GetNode(n.Name)
					if !ok {
						log.Println("unable to find node", n.Name)
					} else {
						node.Update(n)
						settings.observeNodeRegion(node.Region())
					}
					node.Show()
				}
			},
		},
	)
	go controller.Run(ctx.Done())

}

// isTerminalPod returns true if the pod is deleting or in a terminal state
func isTerminalPod(p *v1.Pod) bool {
	if !p.DeletionTimestamp.IsZero() {
		return true
	}
	switch p.Status.Phase {
	case v1.PodSucceeded, v1.PodFailed:
		return true
	}
	return false
}
