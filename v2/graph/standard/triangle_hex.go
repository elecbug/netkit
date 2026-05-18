package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// TriangleHexGraph generates a hexagonal lattice graph with a specified edge length.
func TriangleHexGraph(seed int, directed bool, weightFunc WeightedFunc, edge int) (*graph.Graph, error) {
	if weightFunc == nil {
		weightFunc = Unweighted
	}
	if edge < 0 {
		return nil, fmt.Errorf("edge must be non-negative")
	}

	g := graph.New(directed, true)

	radius := edge - 1

	type coord struct {
		q int
		r int
	}

	nextID := 0
	centerCoord := coord{q: 0, r: 0}
	coordToID := make(map[coord]graph.NodeID)

	inHex := func(c coord) bool {
		return abs(c.q) <= radius && abs(c.r) <= radius && abs(c.q+c.r) <= radius
	}

	nextNumericID := func() graph.NodeID {
		id := graph.NodeID(fmt.Sprintf("%d", nextID))
		nextID++
		return id
	}

	for q := -radius; q <= radius; q++ {
		for r := -radius; r <= radius; r++ {
			c := coord{q: q, r: r}

			if !inHex(c) {
				continue
			}

			id := nextNumericID()
			if err := g.AddNode(id); err != nil {
				return nil, fmt.Errorf("failed to add node: %w", err)
			}
			if node, err := g.Node(id); err != nil {
				return nil, fmt.Errorf("failed to retrieve node: %w", err)
			} else {
				node.AddTag("q", fmt.Sprintf("%d", c.q))
				node.AddTag("r", fmt.Sprintf("%d", c.r))

				if c.q == centerCoord.q && c.r == centerCoord.r {
					node.AddTag("center", "true")
				}
			}

			coordToID[c] = id
		}
	}

	dirs := []coord{
		{q: 1, r: 0},
		{q: 0, r: 1},
		{q: 1, r: -1},
	}

	for q := -radius; q <= radius; q++ {
		for r := -radius; r <= radius; r++ {
			c := coord{q: q, r: r}
			if !inHex(c) {
				continue
			}

			id := coordToID[c]

			for _, d := range dirs {
				neighbor := coord{q: c.q + d.q, r: c.r + d.r}

				nid, ok := coordToID[neighbor]
				if ok {
					node, err := g.Node(id)
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					neighborNode, err := g.Node(nid)
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					if err := g.AddEdge(id, nid, weightFunc(node, neighborNode)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}
		}
	}

	return g, nil
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
