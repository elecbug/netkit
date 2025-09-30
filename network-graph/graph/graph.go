// Package graph provides a simple adjacency map graph for network-graph.
package graph

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elecbug/netkit/network-graph/node"
)

// Graph maintains nodes and adjacency edges.
type Graph struct {
	nodes        map[node.ID]bool
	edges        map[node.ID]map[node.ID]bool
	isUndirected bool
}

// New creates and returns an empty Graph.
func New(isUndirected bool) *Graph {
	return &Graph{
		nodes:        make(map[node.ID]bool),
		edges:        make(map[node.ID]map[node.ID]bool),
		isUndirected: isUndirected,
	}
}

// FromMatrix creates a new graph from a boolean adjacency matrix.
func FromMatrix(matrix [][]bool, bidirectional bool) *Graph {
	g := New(bidirectional)

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			if matrix[i][j] {
				g.AddEdge(node.ID(fmt.Sprintf("%d", i)), node.ID(fmt.Sprintf("%d", j)))
			}
		}
	}

	return g
}

// Save serializes the graph to a string.
func Save(g *Graph) (string, error) {
	if g == nil {
		return "", fmt.Errorf("graph is nil")
	}

	nodes, err := json.Marshal(g.nodes)

	if err != nil {
		return "", fmt.Errorf("failed to marshal nodes: %v", err)
	}

	edges, err := json.Marshal(g.edges)

	if err != nil {
		return "", fmt.Errorf("failed to marshal edges: %v", err)
	}

	bidirectional, err := json.Marshal(g.isUndirected)

	if err != nil {
		return "", fmt.Errorf("failed to marshal bidirectional: %v", err)
	}

	return fmt.Sprintf("%s\n%s\n%s", nodes, edges, bidirectional), nil
}

// Load deserializes the graph from a string.
func Load(data string) (*Graph, error) {
	lines := strings.Split(data, "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("invalid graph data")
	}

	nodes := make(map[node.ID]bool)
	edges := make(map[node.ID]map[node.ID]bool)
	var bidirectional bool

	if err := json.Unmarshal([]byte(lines[0]), &nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nodes: %v", err)
	}

	if err := json.Unmarshal([]byte(lines[1]), &edges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal edges: %v", err)
	}

	if err := json.Unmarshal([]byte(lines[2]), &bidirectional); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bidirectional: %v", err)
	}

	return &Graph{
		nodes:        nodes,
		edges:        edges,
		isUndirected: bidirectional,
	}, nil
}

// [deprecated] This function is deprecated. Use graph.IsUndirected() instead.
// IsBidirectional returns true if the graph is bidirectional.
func (g *Graph) IsBidirectional() bool {
	return g.isUndirected
}

// IsUndirected returns true if the graph is undirected.
func (g *Graph) IsUndirected() bool {
	return g.isUndirected
}

// Hash returns the SHA-256 hash of the graph.
func (g *Graph) Hash() string {
	h := sha256.New()

	h.Write(fmt.Appendf(nil, "%v", g.nodes))
	h.Write(fmt.Appendf(nil, "%v", g.edges))

	return fmt.Sprintf("%x", h.Sum(nil))
}
