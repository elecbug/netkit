package standard

import (
	"fmt"
	"math"

	"github.com/elecbug/netkit/v2/graph"
)

// RandomGeometricGraph generates a random geometric graph (RGG).
// Nodes are placed uniformly at random in the unit square, and edges
// are added between nodes that are within a specified radius r.
func RandomGeometricGraph(seed int, directed bool, weightFunc WeightedFunc, n int, r float64) (*graph.Graph, error) {
	if r < 0 || r > 1 {
		return nil, fmt.Errorf("invalid radius: r must be between 0 and 1")
	}
	if n < 0 {
		return nil, fmt.Errorf("invalid number of nodes: n must be non-negative")
	}
	if weightFunc == nil {
		weightFunc = Unweighted()
	}

	rr := generateRand(seed)
	g := graph.New(directed, true)

	// --- 1. Generate Nodes ---
	type point struct{ x, y float64 }
	positions := make(map[graph.NodeID]point)

	for i := 0; i < n; i++ {
		id := graph.NodeID(fmt.Sprintf("%d", i))
		if err := g.AddNode(id); err != nil {
			return nil, fmt.Errorf("failed to add node: %w", err)
		}
		if node, err := g.Node(id); err != nil {
			return nil, fmt.Errorf("failed to retrieve node: %w", err)
		} else {
			node.AddTag("x", fmt.Sprintf("%f", rr.Float64()))
			node.AddTag("y", fmt.Sprintf("%f", rr.Float64()))
		}
		positions[id] = point{
			x: rr.Float64(),
			y: rr.Float64(),
		}
	}

	// --- 2. Generate Edges ---
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			from := graph.NodeID(fmt.Sprintf("%d", i))
			to := graph.NodeID(fmt.Sprintf("%d", j))
			pF, pT := positions[from], positions[to]

			dx := pF.x - pT.x
			dy := pF.y - pT.y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= r {
				if err := g.AddEdge(from, to, weightFunc(from, to)); err != nil {
					return nil, fmt.Errorf("failed to add edge: %w", err)
				}
			}
		}
	}

	return g, nil
}

// RForRandomGeometricGraph calculates the radius r for a random geometric graph to achieve a target average degree.
func RForRandomGeometricGraph(targetDegree float64, n int) (float64, error) {
	if targetDegree < 0 {
		return 0, fmt.Errorf("invalid target degree: must be non-negative")
	}
	if n <= 0 {
		return 0, fmt.Errorf("invalid number of nodes: n must be positive")
	}

	r := math.Sqrt(targetDegree / (math.Pi * float64(n-1)))

	return r, nil
}
