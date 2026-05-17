package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// GridGraph generates a grid graph with the specified number of rows and columns.
// If torus is true, the graph will wrap around at the edges, creating a toroidal structure.
func GridGraph(seed int, directed bool, weightFunc WeightedFunc, rows, cols int, torus bool) (*graph.Graph, error) {
	if weightFunc == nil {
		weightFunc = Unweighted
	}
	if rows < 0 || cols < 0 {
		return nil, fmt.Errorf("rows and cols must be non-negative")
	}

	g := graph.New(directed, true)

	nodeID := func(row, col int) graph.NodeID {
		return graph.NodeID(fmt.Sprintf("%d", row*cols+col))
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			g.AddNode(nodeID(i, j))
			if node, err := g.Node(nodeID(i, j)); err != nil {
				return nil, fmt.Errorf("failed to retrieve node: %w", err)
			} else {
				node.AddTag("x", fmt.Sprintf("%d", i))
				node.AddTag("y", fmt.Sprintf("%d", j))
			}
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			id := nodeID(i, j)

			if torus {
				if rows > 2 {
					ni := (i + 1) % rows
					if err := g.AddEdge(id, nodeID(ni, j), weightFunc(id, nodeID(ni, j))); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			} else {
				if i < rows-1 {
					if err := g.AddEdge(id, nodeID(i+1, j), weightFunc(id, nodeID(i+1, j))); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}

			if torus {
				if cols > 2 {
					nj := (j + 1) % cols
					if err := g.AddEdge(id, nodeID(i, nj), weightFunc(id, nodeID(i, nj))); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			} else {
				if j < cols-1 {
					if err := g.AddEdge(id, nodeID(i, j+1), weightFunc(id, nodeID(i, j+1))); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}
		}
	}

	return g, nil
}
