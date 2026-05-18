package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// GridGraph generates a grid graph with the specified number of rows and columns.
// If torus is true, the graph will wrap around at the edges, creating a toroidal structure.
func GridGraph(seed int, directed bool, weightFunc WeightedFunc, rows, cols int, torus bool) (*graph.Graph, error) {
	if rows < 0 || cols < 0 {
		return nil, fmt.Errorf("rows and cols must be non-negative")
	}

	g := graph.New(directed, weightFunc != nil)

	if weightFunc == nil {
		weightFunc = func(from, to *graph.Node) *graph.Weight {
			return nil
		}
	}

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
			node, err := g.Node(id)
			if err != nil {
				return nil, fmt.Errorf("failed to get node: %w", err)
			}

			if torus {
				if rows > 2 {
					ni := (i + 1) % rows
					neighborNode, err := g.Node(nodeID(ni, j))
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					if err := g.AddEdge(id, nodeID(ni, j), weightFunc(node, neighborNode)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			} else {
				if i < rows-1 {
					neighborNode, err := g.Node(nodeID(i+1, j))
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					if err := g.AddEdge(id, nodeID(i+1, j), weightFunc(node, neighborNode)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}

			if torus {
				if cols > 2 {
					nj := (j + 1) % cols
					neighborNode, err := g.Node(nodeID(i, nj))
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					if err := g.AddEdge(id, nodeID(i, nj), weightFunc(node, neighborNode)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			} else {
				if j < cols-1 {
					neighborNode, err := g.Node(nodeID(i, j+1))
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					if err := g.AddEdge(id, nodeID(i, j+1), weightFunc(node, neighborNode)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}
		}
	}

	return g, nil
}
