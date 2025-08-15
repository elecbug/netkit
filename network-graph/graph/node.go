package graph

import (
	"fmt"

	"github.com/elecbug/go-dspkg/network-graph/node"
)

// AddNode adds a node to the graph.
func (g *Graph) AddNode(id node.ID) error {
	if _, ok := g.nodes[id]; !ok {
		g.nodes[id] = true

		return nil
	} else {
		return fmt.Errorf("node %s already exists", id)
	}
}

// RemoveNode removes a node and its incident edges from the graph.
func (g *Graph) RemoveNode(id node.ID) error {
	if _, ok := g.nodes[id]; !ok {
		return fmt.Errorf("node %s does not exist", id)
	}

	delete(g.nodes, id)
	delete(g.edges, id)

	for from := range g.edges {
		delete(g.edges[from], id)
	}

	return nil
}

// FindNode reports whether a node with the given id exists.
func (g *Graph) FindNode(id node.ID) (bool, error) {
	if _, ok := g.nodes[id]; ok {
		return true, nil
	}

	return false, fmt.Errorf("node %s does not exist", id)
}
