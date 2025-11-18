package standard_graph

import (
	"github.com/elecbug/netkit/graph"
)

// ErdosRenyiGraph generates a random graph based on the Erdős-Rényi model.
func (sg *StandardGraph) ErdosRenyiGraph(n int, p float64, isUndirected bool) *graph.Graph {
	ra := sg.genRand()

	g := graph.New(isUndirected)

	for i := 0; i < n; i++ {
		g.AddNode(graph.NodeID(toString(i)))
	}

	if isUndirected {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if ra.Float64() < p {
					g.AddEdge(graph.NodeID(toString(i)), graph.NodeID(toString(j)))
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && ra.Float64() < p {
					g.AddEdge(graph.NodeID(toString(i)), graph.NodeID(toString(j)))
				}
			}
		}
	}

	return g
}
