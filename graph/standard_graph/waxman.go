package standard_graph

import (
	"math"

	"github.com/elecbug/netkit/graph"
	"github.com/elecbug/netkit/graph/node"
)

// WaxmanGraph generates a Waxman random graph.
// n = number of nodes
// alpha, beta = Waxman parameters (0<alpha<=1, 0<beta<=1)
// isUndirected = undirected or directed graph
func WaxmanGraph(n int, alpha, beta float64, isUndirected bool) *graph.Graph {
	if n < 1 || alpha <= 0 || alpha > 1 || beta <= 0 || beta > 1 {
		return nil
	}

	ra := genRand()
	g := graph.New(isUndirected)

	// --- 1. Generate Node Positions ---
	type point struct{ x, y float64 }
	positions := make(map[node.ID]point)

	for i := 0; i < n; i++ {
		id := node.ID(toString(i))
		g.AddNode(id)
		positions[id] = point{
			x: ra.Float64(),
			y: ra.Float64(),
		}
	}

	// Maximum distance (diagonal)
	L := math.Sqrt(2.0)

	// --- 2. Generate Edges Based on Distance ---
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			id1 := node.ID(toString(i))
			id2 := node.ID(toString(j))
			p1, p2 := positions[id1], positions[id2]

			dx := p1.x - p2.x
			dy := p1.y - p2.y
			dist := math.Sqrt(dx*dx + dy*dy)

			// Waxman probability
			prob := alpha * math.Exp(-dist/(beta*L))

			if ra.Float64() < prob {
				g.AddEdge(id1, id2)
			}
		}
	}

	return g
}
