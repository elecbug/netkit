package graph

import "fmt"

// NodeID uniquely identifies a node in a network-graph.
type NodeID string

// Weight represents the weight of an edge in the graph.
type Weight float64

// Node represents a node in the graph, containing its ID and edges to other nodes.
type Node struct {
	ID    NodeID            // ID is the unique identifier for the node.
	edges map[NodeID]Weight // Edges maps the destination NodeID to the weight of the edge.
	tags  map[string]string // Tags can hold additional metadata about the node.
}

// NewNode creates a new node with the given ID.
func NewNode(id NodeID) *Node {
	return &Node{
		ID:    id,
		edges: make(map[NodeID]Weight),
		tags:  make(map[string]string),
	}
}

/* Edge */

// addEdge adds an edge from this node to another node with the specified weight.
func (n *Node) addEdge(to NodeID, weight Weight) error {
	if _, exists := n.edges[to]; exists {
		return fmt.Errorf("edge to node %s already exists", to)
	}

	n.edges[to] = weight
	return nil
}

// removeEdge removes the edge from this node to the specified destination node.
func (n *Node) removeEdge(to NodeID) error {
	if _, exists := n.edges[to]; !exists {
		return fmt.Errorf("edge to node %s does not exist", to)
	}

	delete(n.edges, to)
	return nil
}

// hasEdge checks if there is an edge from this node to the specified destination node.
func (n *Node) hasEdge(to NodeID) bool {
	_, exists := n.edges[to]
	return exists
}

// edgeWeight returns the weight of the edge from this node to the specified destination node, if it exists.
func (n *Node) edgeWeight(to NodeID) (Weight, error) {
	weight, exists := n.edges[to]
	if !exists {
		return 0, fmt.Errorf("edge to node %s does not exist", to)
	}

	return weight, nil
}

/* Tagging */

// AddTag adds a key-value pair as a tag to the node.
// It returns an error if a tag with the same key already exists.
func (n *Node) AddTag(key, value string) error {
	if n.HasTag(key) {
		return fmt.Errorf("tag with key %s already exists", key)
	}

	n.tags[key] = value
	return nil
}

// UpdateTag updates the value of an existing tag on the node.
// If the tag does not exist, it will be added.
func (n *Node) UpdateTag(key, value string) {
	n.tags[key] = value
}

// RemoveTag removes a tag from the node by its key. It returns an error if the tag does not exist.
func (n *Node) RemoveTag(key string) error {
	if _, exists := n.tags[key]; !exists {
		return fmt.Errorf("tag with key %s does not exist", key)
	}

	delete(n.tags, key)
	return nil
}

// HasTag checks if a tag with the specified key exists on the node.
func (n *Node) HasTag(key string) bool {
	_, exists := n.tags[key]
	return exists
}

// Tag retrieves the value of a tag by its key. It returns an error if the tag does not exist.
func (n *Node) Tag(key string) (string, bool) {
	value, exists := n.tags[key]
	if !exists {
		return "", false
	}
	return value, true
}

/* Connectivity */

// Neighbors returns a slice of NodeIDs representing the neighbors of this node.
func (n *Node) Neighbors() []NodeID {
	var neighbors []NodeID
	for to := range n.edges {
		neighbors = append(neighbors, to)
	}
	return neighbors
}

// Degree returns the number of edges connected to this node.
func (n *Node) Degree() int {
	return len(n.edges)
}

/* Formatting */

// String returns a string representation of the node, including its ID and edges.
func (n *Node) String() string {
	return fmt.Sprintf("Node(ID: %s, Edges: %v)", n.ID, n.edges)
}

/* Weight */

func NewWeight(value float64) *Weight {
	weight := Weight(value)
	return &weight
}
