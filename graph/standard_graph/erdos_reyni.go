package standard_graph

import (
	"github.com/elecbug/netkit/graph"
	"github.com/elecbug/netkit/graph/node"
)

// ErdosRenyiGraph generates a random graph based on the Erdős-Rényi model.
func ErdosRenyiGraph(n int, p float64, isUndirected bool) *graph.Graph {
	ra := genRand()

	g := graph.New(isUndirected)

	for i := 0; i < n; i++ {
		g.AddNode(node.ID(toString(i)))
	}

	if isUndirected {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if ra.Float64() < p {
					g.AddEdge(node.ID(toString(i)), node.ID(toString(j)))
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && ra.Float64() < p {
					g.AddEdge(node.ID(toString(i)), node.ID(toString(j)))
				}
			}
		}
	}

	return g
}
