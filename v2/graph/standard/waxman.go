package standard

import (
	"fmt"
	"math"

	"github.com/elecbug/netkit/v2/graph"
)

type waxmanPoint struct {
	x float64
	y float64
}

func waxmanDistance(a, b waxmanPoint) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return math.Sqrt(dx*dx + dy*dy)
}

// WaxmanGraph generates a Waxman random graph.
// Nodes are placed uniformly at random in a 2D unit square.
// For each node pair (u, v), an edge is added with probability:
//
//	P(u, v) = alpha * exp(-dist(u, v) / (beta * L))
//
// L is the maximum Euclidean distance among all node pairs.
// alpha controls edge density.
// beta controls locality sensitivity.
func WaxmanGraph(seed int, directed bool, weightFunc WeightedFunc, n int, alpha, beta float64) (*graph.Graph, error) {
	if n < 0 {
		return nil, fmt.Errorf("invalid number of nodes: n must be non-negative")
	}
	if alpha <= 0 || alpha > 1 {
		return nil, fmt.Errorf("invalid alpha: alpha must be in (0, 1]")
	}
	if beta <= 0 || beta > 1 {
		return nil, fmt.Errorf("invalid beta: beta must be in (0, 1]")
	}

	r := generateRand(seed)
	g := graph.New(directed, weightFunc != nil)

	if weightFunc == nil {
		weightFunc = func(from, to *graph.Node) *graph.Weight {
			return nil
		}
	}

	points := make([]waxmanPoint, n)

	// --- 1. Generate Nodes and Positions ---
	for i := 0; i < n; i++ {
		id := graph.NodeID(fmt.Sprintf("%d", i))

		if err := g.AddNode(id); err != nil {
			return nil, fmt.Errorf("failed to add node: %w", err)
		}

		points[i] = waxmanPoint{
			x: r.Float64(),
			y: r.Float64(),
		}
	}

	if n <= 1 {
		return g, nil
	}

	// --- 2. Compute L: maximum pairwise distance ---
	L := 0.0

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := waxmanDistance(points[i], points[j])
			if d > L {
				L = d
			}
		}
	}

	if L == 0 {
		return g, nil
	}

	// --- 3. Generate Edges ---
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := waxmanDistance(points[i], points[j])
			p := alpha * math.Exp(-d/(beta*L))

			if r.Float64() >= p {
				continue
			}

			from := graph.NodeID(fmt.Sprintf("%d", i))
			to := graph.NodeID(fmt.Sprintf("%d", j))

			fromNode, err := g.Node(from)
			if err != nil {
				return nil, fmt.Errorf("failed to get node: %w", err)
			}

			toNode, err := g.Node(to)
			if err != nil {
				return nil, fmt.Errorf("failed to get node: %w", err)
			}

			if err := g.AddEdge(from, to, weightFunc(fromNode, toNode)); err != nil {
				return nil, fmt.Errorf("failed to add edge: %w", err)
			}
		}
	}

	return g, nil
}
