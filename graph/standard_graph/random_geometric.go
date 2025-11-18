package standard_graph

import (
	"math"

	"github.com/elecbug/netkit/graph"
)

// RandomGeometricGraph generates a random geometric graph (RGG).
// n = number of nodes
// r = connection radius (0~1)
// isUndirected = undirected or directed graph
func (sg *StandardGraph) RandomGeometricGraph(n int, r float64, isUndirected bool) *graph.Graph {
	if n < 1 || r <= 0 {
		return nil
	}

	ra := sg.genRand()
	g := graph.New(isUndirected)

	// --- 1. Generate Nodes ---
	type point struct{ x, y float64 }
	positions := make(map[graph.NodeID]point)

	for i := 0; i < n; i++ {
		id := graph.NodeID(toString(i))
		g.AddNode(id)
		positions[id] = point{
			x: ra.Float64(),
			y: ra.Float64(),
		}
	}

	// --- 2. Generate Edges ---
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			id1 := graph.NodeID(toString(i))
			id2 := graph.NodeID(toString(j))
			p1, p2 := positions[id1], positions[id2]

			dx := p1.x - p2.x
			dy := p1.y - p2.y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= r {
				g.AddEdge(id1, id2)
			}
		}
	}

	return g
}
