# Netkit

**Netkit** is a Go graph algorithm library focused on clarity, extensibility, and practical performance.

It provides reusable graph data structures and common network analysis algorithms, with selected results validated against [NetworkX](https://networkx.org/).

- Go 1.25+
- Module: `github.com/elecbug/netkit`
- License: MIT

---

## Features

Netkit provides graph utilities and network analysis algorithms for both directed and undirected graphs.

Current focus areas include:

- Graph data structures
- Shortest path analysis
- Centrality metrics
- Clustering metrics
- Graph diameter
- PageRank
- Modularity analysis
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

---

## Installation

```bash
go get github.com/elecbug/netkit@latest
````

---

## Usage

> [!NOTE]
> Netkit is under active development. Public APIs may change before the first stable release.

```go
package main

import (
	"fmt"

	"github.com/elecbug/netkit/v2/analyzer"
	"github.com/elecbug/netkit/v2/graph"
)

func main() {
	g := graph.NewGraph(false)

	g.AddNode("0")
	g.AddNode("1")
	g.AddNode("2")

	g.AddEdge("0", "1")
	g.AddEdge("1", "2")

	a := analyzer.NewAnalyzer(g, nil)

	degree, err := a.DegreeCentrality()
	if err != nil {
		panic(err)
	}

	fmt.Println(degree)
}
```

> API details may vary by version. See package documentation and examples for the latest usage.

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
   Graph structures and analysis components should be easy to extend.

3. **Practical performance**
   Implementations should be efficient enough for medium-to-large network analysis workloads.

4. **Validation**
   Results should be checked against established libraries such as NetworkX when possible.

---

## License

MIT © 2025 elecbug. See [`LICENSE`](./LICENSE).

---

## Credits

This project reimplements common graph and network algorithms in Go with selected results validated against NetworkX.

NetworkX is © the NetworkX Developers and distributed under the BSD 3-Clause License.
