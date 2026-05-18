package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// RandomRegularGraph generates a random k-regular graph with n nodes.
// Each node has exactly degree k. Returns nil if impossible.
// This implementation requires n*k to be even.
func RandomRegularGraph(seed int, directed bool, weightFunc WeightedFunc, n, k int) (*graph.Graph, error) {
	if k < 0 || k >= n {
		// degree must be between 0 and n-1
		return nil, fmt.Errorf("invalid degree: k must be between 0 and n-1")
	}
	if n < 0 {
		return nil, fmt.Errorf("invalid number of nodes: n must be non-negative")
	}

	if (n*k)%2 != 0 {
		// n*k must be even for this implementation
		return nil, fmt.Errorf("invalid parameters: n*k must be even for undirected graphs")
	}

	r := generateRand(seed)
	g := graph.New(directed, weightFunc != nil)

	if weightFunc == nil {
		weightFunc = func(from, to *graph.Node) *graph.Weight {
			return nil
		}
	}

	// add nodes
	for i := 0; i < n; i++ {
		if err := g.AddNode(graph.NodeID(fmt.Sprintf("%d", i))); err != nil {
			return nil, fmt.Errorf("failed to add node: %w", err)
		}
	}

	// duplicate each node k times as stubs
	stubs := make([]graph.NodeID, 0, n*k)
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			stubs = append(stubs, graph.NodeID(fmt.Sprintf("%d", i)))
		}
	}

	// shuffle
	r.Shuffle(len(stubs), func(i, j int) { stubs[i], stubs[j] = stubs[j], stubs[i] })

	// attempt to create edges (self-loop / duplicate prevention)
	maxTries := n * k * 10
	edges := make(map[[2]string]bool)
	try := 0

	for len(stubs) > 1 && try < maxTries {
		a := stubs[len(stubs)-1]
		b := stubs[len(stubs)-2]
		stubs = stubs[:len(stubs)-2]

		// self-loop is not allowed
		if a == b {
			// put them back and shuffle
			stubs = append(stubs, a, b)
			r.Shuffle(len(stubs), func(i, j int) { stubs[i], stubs[j] = stubs[j], stubs[i] })
			try++
			continue
		}

		// if undirected, edge key must be sorted
		key := [2]string{string(a), string(b)}
		if !directed && key[0] > key[1] {
			key[0], key[1] = key[1], key[0]
		}

		if edges[key] {
			// edge already exists, put them back and shuffle
			stubs = append(stubs, a, b)
			r.Shuffle(len(stubs), func(i, j int) { stubs[i], stubs[j] = stubs[j], stubs[i] })
			try++
			continue
		}

		// add edge
		fromNode, err := g.Node(a)
		if err != nil {
			return nil, fmt.Errorf("failed to get node: %w", err)
		}
		toNode, err := g.Node(b)
		if err != nil {
			return nil, fmt.Errorf("failed to get node: %w", err)
		}
		if err := g.AddEdge(a, b, weightFunc(fromNode, toNode)); err != nil {
			return nil, fmt.Errorf("failed to add edge: %w", err)
		}
		edges[key] = true
	}

	if len(stubs) > 0 {
		return RandomRegularGraph(seed*seed, directed, weightFunc, n, k)
	}

	return g, nil
}
