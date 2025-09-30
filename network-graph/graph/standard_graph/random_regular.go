package standard_graph

import (
	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
)

// RandomRegularGraph generates a random k-regular graph with n nodes.
// Each node has exactly degree k. Returns nil if impossible.
func RandomRegularGraph(n, k int, isUndirected bool) *graph.Graph {
	if k < 0 || n < 1 || k >= n {
		return nil
	}
	if (n*k)%2 != 0 {
		// if undirected, n*k must be even
		return nil
	}

	ra := genRand()
	g := graph.New(isUndirected)

	// add nodes
	for i := 0; i < n; i++ {
		g.AddNode(node.ID(toString(i)))
	}

	// duplicate each node k times as stubs
	stubs := make([]node.ID, 0, n*k)
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			stubs = append(stubs, node.ID(toString(i)))
		}
	}

	// shuffle
	ra.Shuffle(len(stubs), func(i, j int) { stubs[i], stubs[j] = stubs[j], stubs[i] })

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
			ra.Shuffle(len(stubs), func(i, j int) { stubs[i], stubs[j] = stubs[j], stubs[i] })
			try++
			continue
		}

		// if undirected, edge key must be sorted
		key := [2]string{string(a), string(b)}
		if isUndirected && key[0] > key[1] {
			key[0], key[1] = key[1], key[0]
		}

		if edges[key] {
			// edge already exists, put them back and shuffle
			stubs = append(stubs, a, b)
			ra.Shuffle(len(stubs), func(i, j int) { stubs[i], stubs[j] = stubs[j], stubs[i] })
			try++
			continue
		}

		// add edge
		g.AddEdge(a, b)
		edges[key] = true
	}

	if len(stubs) > 0 {
		// if impossible to satisfy conditions
		return nil
	}

	return g
}
