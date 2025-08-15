# GO - Data Structure PacKaGe

Go (generic) data structures and graph algorithms with a focus on clarity and performance.

- Go 1.21+
- Module: `github.com/elecbug/go-dspkg`

## Install

```powershell
go get github.com/elecbug/go-dspkg@latest
```

## Packages

- `bimap`: Bidirectional map with O(1) lookups key->value and value->key.
- `slice`: Generic helpers: binary search, stable merge sort, parallel sort, and `IsSorted`.
- `graph`: Dense graph (adjacency-matrix) with directed/undirected and weighted/unweighted modes.
- `graph/graph_algorithm`: Shortest paths, diameter, centralities, and efficiencies. Sequential and parallel units.
- `network-graph`: A separate, lightweight graph with its own algorithms (legacy/alt API).

## Quick start

### Bidirectional map (bimap)

```go
package main

import (
	"fmt"
	"github.com/elecbug/go-dspkg/bimap"
)

func main() {
	m := bimap.New[string, int]()
	m.Set("alice", 1)
	m.Set("bob", 2)

	if id, ok := m.GetByKey("alice"); ok {
		fmt.Println("alice ->", id)
	}
	if name, ok := m.GetByValue(2); ok {
		fmt.Println("2 ->", name)
	}
}
```

### Slice utilities (binary search, sort, parallel sort)

```go
package main

import (
	"fmt"
	sl "github.com/elecbug/go-dspkg/slice"
)

func main() {
	// Binary search on a sorted slice
	arr := []int{1, 3, 5, 7, 9}
	target := 7
	idx := sl.Bsearch(arr, func(x int) sl.CompareType {
		if x == target {
			return sl.EQUAL
		}
		if x < target {
			return sl.TARGET_SMALL // search right
		}
		return sl.TARGET_BIG // search left
	})
	fmt.Println("index of 7:", idx)

	// In-place stable merge sort
	sl.Sort(arr, func(a, b int) bool { return a < b })

	// Parallel sort (limit depth with 'level')
	sl.ParallelSort(arr, func(a, b int) bool { return a < b }, 3)

	fmt.Println("sorted:", arr)
}
```

### Graph + algorithms

```go
package main

import (
	"fmt"
	"github.com/elecbug/go-dspkg/graph"
	ga "github.com/elecbug/go-dspkg/graph/graph_algorithm"
	gt "github.com/elecbug/go-dspkg/graph/graph_type"
)

func main() {
	// Undirected, unweighted graph with capacity for 10 nodes
	g := graph.NewGraph(gt.UNDIRECTED_UNWEIGHTED, 10)

	a, _ := g.AddNode()
	b, _ := g.AddNode()
	c, _ := g.AddNode()

	_ = g.AddEdge(a, b)
	_ = g.AddEdge(b, c)

	u := ga.NewUnit(g) // sequential unit

	// Shortest path
	sp := u.ShortestPath(a, c)
	fmt.Println("distance a->c:", sp.Distance(), "path:", sp.Nodes())

	// Diameter (longest shortest-path)
	d := u.Diameter()
	fmt.Println("diameter distance:", d.Distance(), "path:", d.Nodes())

	// Centralities
	deg := u.DegreeCentrality()
	btw := u.BetweennessCentrality()
	fmt.Println("degree:", deg)
	fmt.Println("betweenness:", btw)

	// Averages / percentiles
	avg := u.AverageShortestPathLength()
	p95 := u.PercentileShortestPathLength(0.95)
	fmt.Println("avg spl:", avg, "p95 spl:", p95)

	// Parallel unit (use multiple cores)
	pu := ga.NewParallelUnit(g, 4)
	_ = pu.Diameter()
}
```

Notes
- For weighted graphs, construct with `gt.DIRECTED_WEIGHTED` or `gt.UNDIRECTED_WEIGHTED` and add edges via `AddWeightEdge(from, to, distance)`.
- Graph uses an adjacency-matrix internally; `INF` distances (no edge) are handled for you.

### Network graph (alt API)

`network-graph` offers a separate graph and algorithms package with a simplified API. See `network-graph/` for details and tests.

## Development

- Run tests

```powershell
go test ./...
```

## License

MIT © 2025 elecbug. See `LICENSE`.