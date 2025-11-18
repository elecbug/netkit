package standard_graph

import (
	"github.com/elecbug/netkit/graph"
)

// BarabasiAlbertGraph generates a graph based on the Barabási–Albert preferential attachment model.
func (sg *StandardGraph) BarabasiAlbertGraph(n int, m int, isUndirected bool) *graph.Graph {
	if m < 1 || n <= m {
		return nil
	}

	ra := sg.genRand()
	g := graph.New(isUndirected)

	// --- 1. initialize ---
	for i := 0; i < m; i++ {
		g.AddNode(graph.NodeID(toString(i)))
	}
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			g.AddEdge(graph.NodeID(toString(i)), graph.NodeID(toString(j)))
		}
	}

	// --- 2. preferential attachment ---
	for i := m; i < n; i++ {
		newNode := graph.NodeID(toString(i))
		g.AddNode(newNode)

		// calculate current node degrees
		degrees := make(map[graph.NodeID]int)
		totalDegree := 0
		for _, id := range g.Nodes() {
			d := len(g.Neighbors(id)) // each node degree
			degrees[id] = d
			totalDegree += d
		}

		// degree based sampling
		chosen := make(map[graph.NodeID]bool)
		for len(chosen) < m {
			r := ra.Intn(totalDegree)
			accum := 0
			var target graph.NodeID
			for id, d := range degrees {
				accum += d
				if r < accum {
					target = id
					break
				}
			}
			// self-loop and duplicate edges are not allowed
			if target != newNode && !chosen[target] {
				g.AddEdge(newNode, target)
				chosen[target] = true
			}
		}
	}

	return g
}
