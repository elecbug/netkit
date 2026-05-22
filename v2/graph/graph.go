package graph

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
)

// Graph maintains nodes and adjacency edges.
type Graph struct {
	nodes    map[NodeID]*Node // nodes maps NodeID to the corresponding Node struct.
	directed bool             // Directed indicates whether the graph is directed (true) or undirected (false).
	weighted bool             // Weighted indicates whether the graph is weighted (true) or unweighted (false).
}

// graphSerialization is a helper struct for JSON serialization of the Graph.
type graphSerialization struct {
	Nodes    map[NodeID]map[NodeID]Weight `json:"nodes"`
	Directed bool                         `json:"directed"`
	Weighted bool                         `json:"weighted"`
}

// Matrix represents the adjacency matrix of the graph, where the value at matrix[i][j] is
// the weight of the edge from node i to node j, or 0 if no edge exists.
type Matrix [][]Weight

// New creates and returns an empty Graph.
func New(directed bool, weighted bool) *Graph {
	return &Graph{
		nodes:    make(map[NodeID]*Node),
		directed: directed,
		weighted: weighted,
	}
}

// Free clears all nodes and edges from the graph, effectively resetting it to an empty state.
func (g *Graph) Free() {
	for id := range g.nodes {
		delete(g.nodes, id)
	}
}

/* Node */

// AddNode adds a node to the graph.
func (g *Graph) AddNode(id NodeID) error {
	if _, ok := g.nodes[id]; !ok {
		g.nodes[id] = NewNode(id)
		return nil
	} else {
		return fmt.Errorf("node %s already exists", id)
	}
}

// RemoveNode removes a node from the graph, along with all edges to and from that node.
func (g *Graph) RemoveNode(id NodeID) error {
	if _, ok := g.nodes[id]; !ok {
		return fmt.Errorf("node %s does not exist", id)
	}

	delete(g.nodes, id)

	// Remove edges to this node from other nodes
	for _, node := range g.nodes {
		node.removeEdge(id)
	}

	return nil
}

// HasNode checks if a node with the given ID exists in the graph.
func (g *Graph) HasNode(id NodeID) bool {
	_, exists := g.nodes[id]
	return exists
}

// Node returns the node with the given ID, or an error if it does not exist.
func (g *Graph) Node(id NodeID) (*Node, error) {
	if node, ok := g.nodes[id]; ok {
		return node, nil
	} else {
		return nil, fmt.Errorf("node %s does not exist", id)
	}
}

// Nodes returns a slice of all node IDs in the graph.
func (g *Graph) Nodes() []NodeID {
	var nodes []NodeID

	for id := range g.nodes {
		nodes = append(nodes, id)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i] < nodes[j]
	})

	return nodes
}

// Size returns the number of nodes in the graph.
func (g *Graph) Size() int {
	return len(g.nodes)
}

/* Edge */

// AddEdge adds an edge from one node to another with the specified weight.
func (g *Graph) AddEdge(from NodeID, to NodeID, weight *Weight) error {
	fromNode, fromExists := g.nodes[from]
	toNode, toExists := g.nodes[to]

	if !fromExists {
		return fmt.Errorf("node %s does not exist", from)
	}

	if !toExists {
		return fmt.Errorf("node %s does not exist", to)
	}

	if g.weighted {
		if weight == nil || *weight <= 0 {
			return fmt.Errorf("weight must be positive for weighted graphs")
		}

		err := fromNode.addEdge(to, *weight)
		if err != nil {
			return err
		}
		if !g.directed {
			err := toNode.addEdge(from, *weight)
			if err != nil {
				return err
			}
		}
	} else {
		if weight != nil {
			return fmt.Errorf("weight should be nil for unweighted graphs")
		}

		err := fromNode.addEdge(to, 1)
		if err != nil {
			return err
		}
		if !g.directed {
			err := toNode.addEdge(from, 1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveEdge removes the edge from one node to another.
func (g *Graph) RemoveEdge(from NodeID, to NodeID) error {
	fromNode, fromExists := g.nodes[from]
	toNode, toExists := g.nodes[to]

	if !fromExists {
		return fmt.Errorf("node %s does not exist", from)
	}

	if !toExists {
		return fmt.Errorf("node %s does not exist", to)
	}

	err := fromNode.removeEdge(to)
	if err != nil {
		return err
	}
	if !g.directed {
		err := toNode.removeEdge(from)
		if err != nil {
			return err
		}
	}

	return nil
}

// HasEdge checks if there is an edge from one node to another.
func (g *Graph) HasEdge(from NodeID, to NodeID) bool {
	fromNode, fromExists := g.nodes[from]
	_, toExists := g.nodes[to]

	if !fromExists || !toExists {
		return false
	}

	return fromNode.hasEdge(to)
}

// EdgeWeight returns the weight of the edge from one node to another, or an error
// if the edge does not exist.
func (g *Graph) EdgeWeight(from NodeID, to NodeID) (Weight, error) {
	fromNode, fromExists := g.nodes[from]
	_, toExists := g.nodes[to]

	if !fromExists {
		return 0, fmt.Errorf("node %s does not exist", from)
	}

	if !toExists {
		return 0, fmt.Errorf("node %s does not exist", to)
	}

	return fromNode.edgeWeight(to)
}

/* Formatting and Serialization */

// String returns a string representation of the graph, including its nodes and edges.
func (g *Graph) String() string {
	result := "Graph:\n"
	result += fmt.Sprintf("  Directed: %t\n", g.directed)
	result += fmt.Sprintf("  Weighted: %t\n", g.weighted)
	result += "  Nodes:\n"

	type pair struct {
		id   NodeID
		node *Node
	}

	pairs := make([]pair, 0, len(g.nodes))
	for id, node := range g.nodes {
		pairs = append(pairs, pair{id: id, node: node})
	}

	// Sort nodes by ID for consistent output
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].id < pairs[j].id
	})

	for _, p := range pairs {
		result += fmt.Sprintf("    %s: %v\n", p.id, p.node.edges)
	}

	return result
}

// Matrix returns the adjacency matrix representation of the graph, where the value at
// matrix[i][j] is the weight of the edge from node i to node j, or 0 if no edge exists.
func (g *Graph) Matrix() *Matrix {
	nodeIDs := g.Nodes()
	idToIndex := make(map[NodeID]int)
	for i, id := range nodeIDs {
		idToIndex[id] = i
	}

	matrix := make(Matrix, len(nodeIDs))
	for i := range matrix {
		matrix[i] = make([]Weight, len(nodeIDs))
	}

	mapping := make(map[NodeID]int)
	for i, id := range nodeIDs {
		mapping[id] = i
	}

	for fromID, node := range g.nodes {
		fromIndex := mapping[fromID]
		for toID, weight := range node.edges {
			toIndex := mapping[toID]
			matrix[fromIndex][toIndex] = weight
		}
	}

	return &matrix
}

// Serialize serializes the graph to a JSON string.
func (g *Graph) Serialize() (string, error) {
	serialization := graphSerialization{
		Nodes:    make(map[NodeID]map[NodeID]Weight),
		Directed: g.directed,
		Weighted: g.weighted,
	}

	for id, node := range g.nodes {
		serialization.Nodes[id] = node.edges
	}

	jsonBytes, err := json.Marshal(serialization)
	if err != nil {
		return "", fmt.Errorf("error serializing graph: %v", err)
	}

	return string(jsonBytes), nil
}

// Deserialize deserializes a JSON string into a Graph.
func Deserialize(jsonStr string) (*Graph, error) {
	var serialization graphSerialization
	err := json.Unmarshal([]byte(jsonStr), &serialization)
	if err != nil {
		return nil, fmt.Errorf("error deserializing graph: %v", err)
	}

	g := New(true, true)

	for id := range serialization.Nodes {
		if err := g.AddNode(id); err != nil {
			return nil, fmt.Errorf("error adding node %s: %v", id, err)
		}
	}

	for fromID, edges := range serialization.Nodes {
		for toID, weight := range edges {
			var weightPtr *Weight
			if g.weighted {
				weightCopy := weight
				weightPtr = &weightCopy
			}
			if err := g.AddEdge(fromID, toID, weightPtr); err != nil {
				return nil, fmt.Errorf("error adding edge from %s to %s: %v", fromID, toID, err)
			}
		}
	}

	g.directed = serialization.Directed
	g.weighted = serialization.Weighted

	if err := g.checkProperties(); err != nil {
		return nil, fmt.Errorf("graph properties check failed: %v", err)
	}

	return g, nil
}

// Hash returns the SHA-256 hash of the graph.
func (g *Graph) Hash() string {
	h := sha256.New()
	h.Write(fmt.Appendf(nil, "%s", g.String()))

	return fmt.Sprintf("%x", h.Sum(nil))
}

/* Utility Functions */

// checkProperties checks that the graph's properties (directed/undirected, weighted/unweighted) are consistent with its edges.
func (g *Graph) checkProperties() error {
	if !g.directed {
		for fromID, node := range g.nodes {
			for toID, weight := range node.edges {
				toNode, _ := g.Node(toID)
				if w, exists := toNode.edges[fromID]; !exists || w != weight {
					if !exists {
						return fmt.Errorf("undirected graph property violated: edge from %s to %s has weight %f but reverse edge does not exist", fromID, toID, weight)
					} else if w != weight {
						return fmt.Errorf("undirected graph property violated: edge from %s to %s has weight %f but reverse edge has weight %f", fromID, toID, weight, w)
					}
				}
			}
		}
	}

	if !g.weighted {
		for fromID, node := range g.nodes {
			for toID, weight := range node.edges {
				if weight != 1 {
					return fmt.Errorf("unweighted graph property violated: edge from %s to %s has weight %f", fromID, toID, weight)
				}
			}
		}
	}

	return nil
}

/* Properties */

// IsDirected returns true if the graph is directed, false otherwise.
func (g *Graph) IsDirected() bool {
	return g.directed
}

// IsWeighted returns true if the graph is weighted, false otherwise.
func (g *Graph) IsWeighted() bool {
	return g.weighted
}
