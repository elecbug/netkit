# Netkit

**Netkit** is a Go toolkit for graph algorithms, network analysis, and P2P network simulation.

It provides reusable graph data structures, common network analysis algorithms, and a programmable P2P layer for testing message propagation over graph-based network topologies.

Selected graph algorithm results are validated against [NetworkX](https://networkx.org/).

- Go 1.25+
- Module: `github.com/elecbug/netkit`
- License: MIT

---

## Features

Netkit provides graph utilities, network analysis algorithms, and P2P simulation tools for both directed and undirected graphs.

Current focus areas include:

- Graph data structures
- Shortest path analysis
- Centrality metrics
- Clustering metrics
- Graph diameter
- PageRank
- Modularity analysis
- P2P overlay construction
- Message propagation simulation
- NetworkX-compatible validation for selected algorithms

Implemented or planned algorithms include:

- Degree centrality
- Betweenness centrality
- Edge betweenness centrality
- Closeness centrality
- Eigenvector centrality
- PageRank
- Clustering coefficient
- Degree assortativity coefficient
- Diameter / weighted diameter
- Modularity

The P2P simulation layer supports:

- Creating peer networks from graph topologies
- Defining message propagation behavior
- Configuring processing latency
- Configuring network latency
- Testing broadcast and relay strategies
- Evaluating propagation behavior over generated or custom graphs

---

## Installation

```bash
go get github.com/elecbug/netkit@latest
```

---

## Usage

> [!NOTE]
> Netkit is under active development. Public APIs may change before the first stable release.

The same graph can be used both for graph analysis and as a P2P overlay topology.

```go
package main

import (
	"fmt"

	"github.com/elecbug/netkit/v2/analyzer"
	"github.com/elecbug/netkit/v2/graph"
	"github.com/elecbug/netkit/v2/p2p"
)

func main() {
	g := graph.NewGraph(false)

	g.AddNode("0")
	g.AddNode("1")
	g.AddNode("2")

	g.AddEdge("0", "1")
	g.AddEdge("1", "2")

	a := analyzer.NewAnalyzer(g, 4, analyzer.DefaultConfig())

	degree, err := a.DegreeCentrality()
	if err != nil {
		panic(err)
	}

	fmt.Println("degree centrality:", degree)

	cfg := &p2p.Config{
		ProcessingLatencyFunc: func(src p2p.PeerID) float64 {
			return 100
		},
		NetworkLatencyFunc: func(src, dst p2p.PeerID) float64 {
			return 1
		},
	}

	network, err := p2p.New(g, cfg)
	if err != nil {
		panic(err)
	}

	_ = network

	// The graph is now also available as a P2P overlay.
	// Users can define propagation behavior and run message dissemination tests
	// over the same topology used for graph analysis.
}
```

A propagation function can decide which peers should receive a message next.

```go
package main

import "github.com/elecbug/netkit/v2/p2p"

func ForwardToUnseenPeers(
	id p2p.PeerID,
	msg p2p.Message,
	known []p2p.PeerID,
	sent []p2p.PeerID,
	received []p2p.PeerID,
	params map[string]any,
) *[]p2p.PeerID {
	result := make([]p2p.PeerID, 0)

	for _, peerID := range known {
		alreadySent := false
		for _, s := range sent {
			if peerID == s {
				alreadySent = true
				break
			}
		}
		if alreadySent {
			continue
		}

		alreadyReceived := false
		for _, r := range received {
			if peerID == r {
				alreadyReceived = true
				break
			}
		}
		if alreadyReceived {
			continue
		}

		result = append(result, peerID)
	}

	return &result
}
```

> API details may vary by version. See package documentation and examples for the latest usage.

---

## P2P Simulation

Netkit can form a P2P network directly from a graph.

This allows users to generate or load a topology, analyze its structural properties, and then run message propagation experiments on the same network.

Typical use cases include:

* Broadcast protocol testing
* Gossip-style dissemination experiments
* Overlay topology evaluation
* Reachability analysis
* Duplicate message analysis
* Latency-sensitive propagation experiments
* Comparing propagation behavior across graph models

This makes Netkit useful not only as a graph algorithm library, but also as a lightweight experimental framework for P2P network behavior.

---

## Validation

Netkit reimplements common graph and network algorithms in Go.

Where possible, algorithm outputs are validated against NetworkX to ensure correctness and compatibility of definitions.

Validation currently covers metrics such as:

* Degree centrality
* Betweenness centrality
* Edge betweenness centrality
* Closeness centrality
* Clustering coefficient
* Degree assortativity coefficient
* Diameter
* Weighted diameter
* Eigenvector centrality
* PageRank
* Shortest paths
* Modularity

Some algorithms, especially greedy community detection for modularity, may not produce byte-for-byte identical results to NetworkX because heuristic merge order, tie-breaking, and implementation details can differ.

---

## Development

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Format code:

```bash
gofmt -w .
```

---

## Project Goals

Netkit is designed with the following goals:

1. **Clarity**
   Algorithms should be readable and easy to inspect.

2. **Extensibility**
   Graph structures, analysis components, and P2P protocol logic should be easy to extend.

3. **Practical performance**
   Implementations should be efficient enough for medium-to-large network analysis and simulation workloads.

4. **P2P experimentation**
   Users should be able to construct graph-based overlays and test message propagation behavior on them.

5. **Validation**
   Results should be checked against established libraries such as NetworkX when possible.

---

## License

MIT © 2025 elecbug. See [`LICENSE`](./LICENSE).

---

## Credits

This project reimplements common graph and network algorithms in Go with selected results validated against NetworkX.

NetworkX is © the NetworkX Developers and distributed under the BSD 3-Clause License.