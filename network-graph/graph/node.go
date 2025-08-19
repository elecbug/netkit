package graph

import (
	"fmt"

	"github.com/elecbug/netkit/network-graph/node"
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

// HasNode reports whether a node with the given id exists.
func (g *Graph) HasNode(id node.ID) bool {
	if _, ok := g.nodes[id]; ok {
		return true
	} else {
		return false
	}
}

// Nodes returns a slice of all node IDs in the graph.
func (g *Graph) Nodes() []node.ID {
	var nodes []node.ID

	for id := range g.nodes {
		nodes = append(nodes, id)
	}

	return nodes
}

// Neighbors returns the list of neighbors reachable from the given node id.
func (g *Graph) Neighbors(id node.ID) []node.ID {
	if edges, ok := g.edges[id]; ok {
		var result []node.ID

		for to, v := range edges {
			if v {
				result = append(result, to)
			}
		}

		return result
	}

	return nil
}
