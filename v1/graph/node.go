package graph

import (
	"fmt"
)

// NodeID uniquely identifies a node in a network-graph.
type NodeID string

func (id NodeID) String() string {
	return string(id)
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(id NodeID) error {
	if _, ok := g.nodes[id]; !ok {
		g.nodes[id] = true

		return nil
	} else {
		return fmt.Errorf("node %s already exists", id)
	}
}

// RemoveNode removes a node and its incident edges from the graph.
func (g *Graph) RemoveNode(id NodeID) error {
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
func (g *Graph) HasNode(id NodeID) bool {
	if _, ok := g.nodes[id]; ok {
		return true
	} else {
		return false
	}
}

// Nodes returns a slice of all node IDs in the graph.
func (g *Graph) Nodes() []NodeID {
	var nodes []NodeID

	for id := range g.nodes {
		nodes = append(nodes, id)
	}

	return nodes
}

// Neighbors returns the list of neighbors reachable from the given node id.
func (g *Graph) Neighbors(id NodeID) []NodeID {
	if edges, ok := g.edges[id]; ok {
		var result []NodeID

		for to, v := range edges {
			if v {
				result = append(result, to)
			}
		}

		return result
	}

	return nil
}
