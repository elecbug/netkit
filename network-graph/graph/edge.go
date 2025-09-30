package graph

import (
	"fmt"

	"github.com/elecbug/netkit/network-graph/node"
)

// AddEdge adds an edge from -> to. If bidirectional is true, adds the reverse edge as well.
func (g *Graph) AddEdge(from, to node.ID) error {
	if _, ok := g.nodes[from]; !ok {
		return fmt.Errorf("from node %s does not exist", from)
	}
	if _, ok := g.nodes[to]; !ok {
		return fmt.Errorf("to node %s does not exist", to)
	}

	if _, ok := g.edges[from]; !ok {
		g.edges[from] = make(map[node.ID]bool)
	}

	if g.edges[from][to] {
		return fmt.Errorf("edge from %s to %s already exists", from, to)
	}

	g.edges[from][to] = true

	if g.isUndirected {
		if _, ok := g.edges[to]; !ok {
			g.edges[to] = make(map[node.ID]bool)
		}

		if g.edges[to][from] {
			return fmt.Errorf("edge from %s to %s already exists", to, from)
		}

		g.edges[to][from] = true
	}

	return nil
}

// RemoveEdge removes the edge from -> to. If bidirectional is true, removes the reverse edge as well.
func (g *Graph) RemoveEdge(from, to node.ID) error {
	if _, ok := g.edges[from]; !ok {
		return fmt.Errorf("no edges from node %s", from)
	}

	if _, ok := g.edges[from][to]; !ok {
		return fmt.Errorf("edge from %s to %s does not exist", from, to)
	}

	delete(g.edges[from], to)

	if g.isUndirected {
		if _, ok := g.edges[to]; !ok {
			g.edges[to] = make(map[node.ID]bool)
		}

		delete(g.edges[to], from)
	}

	return nil
}

// HasEdge checks if an edge exists from -> to.
func (g *Graph) HasEdge(from, to node.ID) bool {
	if edges, ok := g.edges[from]; ok {
		return edges[to]
	}

	return false
}

// EdgeCount returns the number of edges in the graph.
func (g *Graph) EdgeCount() int {
	count := 0

	for _, edges := range g.edges {
		count += len(edges)
	}

	return count
}
