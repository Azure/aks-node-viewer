[![GitHub License](https://img.shields.io/badge/License-Apache%202.0-ff69b4.svg)](https://github.com/Azure/aks-node-viewer/blob/main/LICENSE.txt)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/Azure/aks-node-viewer/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/Azure/aks-node-viewer)](https://goreportcard.com/report/github.com/Azure/aks-node-viewer)

## Usage

`aks-node-viewer` is a tool for visualizing dynamic node usage within a cluster.  It was originally developed as an internal tool at AWS for demonstrating consolidation with [Karpenter](https://karpenter.sh/).  It displays the scheduled pod resource requests vs the allocatable capacity on the node.  It *does not* look at the actual pod resource usage.

![](./.static/screenshot.png)

### Installation

Please either fetch the latest [release](https://github.com/Azure/aks-node-viewer/releases) or install manually using:
```shell
go install github.com/Azure/aks-node-viewer/cmd/aks-node-viewer@latest
```

Note: This will install it to your `GOBIN` directory, typically `~/go/bin` if it is unconfigured.

## Usage
```shell
Usage of ./aks-node-viewer:
  -attribution
    	Show the Open Source Attribution
  -context string
    	Name of the kubernetes context to use
  -disable-pricing
    	Disable pricing lookups
  -extra-labels string
    	A comma separated set of extra node labels to display
  -kubeconfig string
    	Absolute path to the kubeconfig file (default "~/.kube/config")
  -node-selector string
    	Node label selector used to filter nodes, if empty all nodes are selected
  -resources string
    	List of comma separated resources to monitor (default "cpu")
  -v	Display aks-node-viewer version
  -version
    	Display aks-node-viewer version
```

### Examples
```shell
# Standard usage
aks-node-viewer
# Karenter nodes only
aks-node-viewer --node-selector "karpenter.sh/provisioner-name"
# Display both CPU and Memory Usage
aks-node-viewer --resources cpu,memory
# Display extra labels, i.e. AZ
aks-node-viewer --extra-labels topology.kubernetes.io/zone
```

### Default Options
You can supply default options to `aks-node-viewer` by creating a file named `.aks-node-viewer` in your home directory and specifying
options there. The format is `option-name=value` where the option names are the command line flags:
```text
# select only Karpenter managed nodes
node-selector=karpenter.sh/provisioner-name

# display both CPU and memory
resources=cpu,memory
```

### Building

```shell
$ make build
```

Or local execution of GoReleaser build:
```shell
$ make goreleaser
```

### Source Attribution

Notice: Files in this source code originated from a fork of https://github.com/aws/eks-node-viewer
which is under an Apache 2.0 license. Those files have been modified to reflect environmental requirements in AKS and Azure.

Many thanks to @ellistarn, @jonathan-innis, @tzneal, @bwagner5, @njtran, and many other developers active in the Karpenter community for laying the foundations of a Karpenter provider ecosystem!
