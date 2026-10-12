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

#### Nix

A [Nix flake](./flake.nix) is included, so on NixOS (or any system with Nix installed) you can run
the viewer without installing anything permanently:

```shell
# run directly from the repository
nix run .
# ...or from GitHub
nix run github:Azure/aks-node-viewer

# build and install into your profile
nix profile install .
```

#### Development shell

The flake provides a development shell with Go, goreleaser, golangci-lint and go-licenses:

```shell
nix develop
# or, if you use direnv, .envrc already runs `use flake`
```

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
  -node-sort string
    	Sort order for the nodes, either 'creation' or a label name. The sort order defaults to ascending and can be controlled by appending =asc or =dsc to the value. (default "creation=dsc")
  -region string
    	Azure region to use for pricing lookups (default "westus2")
  -resources string
    	List of comma separated resources to monitor (default "cpu,memory")
  -style string
    	Three colors to use for styling 'good', 'ok' and 'bad' values, e.g. "#04B575,#FFFF00,#FF0000"
  -v	Display aks-node-viewer version
  -version
    	Display aks-node-viewer version
```

### Examples
```shell
# Standard usage
aks-node-viewer
# Karpenter nodes only
aks-node-viewer --node-selector "karpenter.sh/provisioner-name"
# Display both CPU and Memory Usage
aks-node-viewer --resources cpu,memory
# Display extra labels, i.e. AZ
aks-node-viewer --extra-labels topology.kubernetes.io/zone
# Sort nodes by CPU usage, ascending
aks-node-viewer --node-sort cpu
# Sort nodes newest first
aks-node-viewer --node-sort creation=asc
# Use a custom color scheme
aks-node-viewer --style "#04B575,#FFFF00,#FF0000"
```

### Computed Labels

`aks-node-viewer` supports some custom label names that can be passed to the `--extra-labels` to display additional node information. Both the `aks-node-viewer/` and `eks-node-viewer/` prefixes are accepted, so configs written for [eks-node-viewer](https://github.com/awslabs/eks-node-viewer) keep working.

- `<prefix>/node-age` - Age of the node
- `<prefix>/node-cpu-usage` - CPU requests as a percentage of allocatable CPU
- `<prefix>/node-memory-usage` - Memory requests as a percentage of allocatable memory
- `<prefix>/node-pods-usage` - Pods scheduled as a percentage of the pod capacity
- `<prefix>/node-ephemeral-storage-usage` - Ephemeral storage requests as a percentage of allocatable ephemeral storage

where `<prefix>` is either `aks-node-viewer` or `eks-node-viewer`.

### Default Options
You can supply default options to `aks-node-viewer` by creating a config file. The
file is read from the XDG config directory:

```text
${XDG_CONFIG_HOME:-~/.config}/aks-node-viewer/config
```

If that file does not exist, the legacy `~/.aks-node-viewer` file is read instead.
On first run the legacy file is copied into the XDG location automatically, so
existing configurations keep working. The legacy file is left in place and is
still used as a fallback if the XDG config is later removed.

The format is `option-name=value` where the option names are the command line flags:
```text
# select only Karpenter managed nodes
node-selector=karpenter.sh/provisioner-name

# display both CPU and memory
resources=cpu,memory
```

### Pricing

On startup the viewer immediately displays prices and refreshes them from the
Azure retail prices API in the background, so the UI is never blocked. Prices are
sourced, in order of preference, from:

1. the live Azure retail prices API (fetched during the current run),
2. the on-disk cache written by a previous run, and
3. the static price list embedded at build time.

Fetched prices are cached so the next run starts from a recent baseline instead
of the static snapshot. The cache is stored per region in the XDG cache
directory:

```text
${XDG_CACHE_HOME:-~/.cache}/aks-node-viewer/pricing-<region>.json
```

The cluster summary shows where the displayed prices came from and when they were
last updated, so it is obvious when they are stale, e.g.:

```text
10 nodes  24/96  25.0%  cpu  ...  12.345/hour $9,001.234/month (prices updated 2026-10-12 15:04)
```

While a background refresh is in flight the note ends with `, refreshing`. If no
live prices could be fetched the note makes that explicit, for example
`(prices not updated since 2023-03-01)` or
`(prices cached 2026-10-11 09:00 (1d ago))`.

### Building

```shell
$ make build
```

Or local execution of GoReleaser build:
```shell
$ make goreleaser
```

Or with Nix:
```shell
$ nix build
$ ./result/bin/aks-node-viewer
```

### Source Attribution

Notice: Files in this source code originated from a fork of https://github.com/aws/eks-node-viewer
which is under an Apache 2.0 license. Those files have been modified to reflect environmental requirements in AKS and Azure.

Many thanks to @ellistarn, @jonathan-innis, @tzneal, @bwagner5, @njtran, and many other developers active in the Karpenter community for laying the foundations of a Karpenter provider ecosystem!
