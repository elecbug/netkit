package graph

import (
	"fmt"
	"math"

	"github.com/elecbug/go-dspkg/graph/internal/graph/internal/graph_err" // Custom error package
)

// Graph represents the core structure of a graph.
// It manages nodes, tracks the current unique identifier (`nowID`), and defines the graph type (directed/undirected, weighted/unweighted).
// The `updated` field indicates whether the graph has been modified since the last algorithmic computation.
type Graph struct {
	nodes     *graphNodes // A collection of all nodes in the graph.
	nowID     NodeID      // The next unique identifier to be assigned to a new node.
	graphType GraphType   // The type of the graph (e.g., directed, undirected, weighted, unweighted).
	updated   bool        // Tracks if the graph has been modified since the last update.
	edgeCount int         // Number of edges in graph.
}

// NewGraph creates and initializes a new Graph instance.
//
// Parameters:
//   - graphType: The type of the graph (from the GraphType enumeration).
//   - capacity: The initial capacity for the node collection.
//
// Returns a pointer to the newly created Graph.
func NewGraph(graphType GraphType, capacity int) *Graph {
	return &Graph{
		nodes:     newNodes(capacity),
		nowID:     0,
		graphType: graphType,
		updated:   false,
		edgeCount: 0,
	}
}

// AddNode adds a new node to the graph with the given name.
//
// Parameters:
//   - name: The display name for the node.
//
// Returns the newly created Node and an error if insertion fails.
func (g *Graph) AddNode(name string) (*Node, error) {
	node := newNode(g.nowID, name)
	err := g.nodes.insert(node)

	if err != nil {
		return nil, err
	}

	// Increment the unique identifier for the next node.
	g.nowID++
	g.updated = false // Mark the graph as modified.

	return node, nil
}

// RemoveNode removes a node from the graph using its identifier.
//
// Parameters:
//   - identifier: The unique identifier of the node to remove.
//
// Returns an error if the node does not exist.
func (g *Graph) RemoveNode(identifier NodeID) error {
	g.updated = false // Mark the graph as modified.

	for _, edge := range g.nodes.find(identifier).edges {
		err := g.RemoveEdge(identifier, edge.to)

		if err != nil {
			return err
		}
	}

	return g.nodes.remove(identifier)
}

// FindNode retrieves a node from the graph by its identifier.
//
// Parameters:
//   - identifier: The unique identifier of the node.
//
// Returns the Node and an error if the node does not exist.
func (g *Graph) FindNode(identifier NodeID) (*Node, error) {
	result := g.nodes.find(identifier)

	if result != nil {
		return result, nil
	} else {
		return nil, graph_err.NotExistNode(identifier.String())
	}
}

// FindNodesByName retrieves all nodes with the given name.
//
// Parameters:
//   - name: The name of the nodes to find.
//
// Returns a slice of Nodes and an error if no nodes with the given name exist.
func (g *Graph) FindNodesByName(name string) ([]*Node, error) {
	result := g.nodes.findAll(name)

	if result != nil {
		return result, nil
	} else {
		return nil, graph_err.NotExistNode(name)
	}
}

// AddEdge adds an unweighted edge between two nodes in the graph.
//
// Parameters:
//   - from: The identifier of the source node.
//   - to: The identifier of the destination node.
//
// Returns an error if the edge cannot be added.
func (g *Graph) AddEdge(from, to NodeID) error {
	return g.AddWeightEdge(from, to, 1)
}

// AddWeightEdge adds a weighted edge between two nodes in the graph.
//
// Parameters:
//   - from: The identifier of the source node.
//   - to: The identifier of the destination node.
//   - distance: The weight of the edge.
//
// Returns an error if the edge cannot be added.
func (g *Graph) AddWeightEdge(from, to NodeID, distance Distance) error {
	// Check for invalid edge types and self-loops.
	if (g.graphType == DIRECTED_UNWEIGHTED || g.graphType == UNDIRECTED_UNWEIGHTED) && distance != 1 {
		return graph_err.InvalidEdge(g.graphType.String(), fmt.Sprintf("weight: %d", distance))
	}

	if from == to {
		return graph_err.SelfEdge(from.String())
	}

	f := g.nodes.find(from)
	t := g.nodes.find(to)

	// Ensure both nodes exist in the graph.
	if f == nil {
		return graph_err.NotExistNode(from.String())
	}
	if t == nil {
		return graph_err.NotExistNode(to.String())
	}

	// Add the edge to the source node.
	err := f.addEdge(to, distance)

	if err != nil {
		return err
	}

	// Add a reverse edge for undirected graphs.
	if g.graphType == UNDIRECTED_UNWEIGHTED || g.graphType == UNDIRECTED_WEIGHTED {
		err = t.addEdge(from, distance)

		if err != nil {
			return err
		}
	}

	g.updated = false // Mark the graph as modified.
	g.edgeCount++     // Update edge count

	return nil
}

// RemoveEdge removes an edge between two nodes in the graph.
// For undirected graphs, the reverse edge is also removed.
//
// Parameters:
//   - from: The identifier of the source node.
//   - to: The identifier of the destination node.
//
// Returns:
//   - nil if the edge is successfully removed.
//   - An error if the edge or nodes do not exist, or if the edge is invalid (e.g., a self-loop).
//
// Notes:
//   - The graph's `updated` flag is set to false to indicate that modifications have been made.
//   - For undirected graphs, the reverse edge (to -> from) is also removed.
func (g *Graph) RemoveEdge(from, to NodeID) error {
	if from == to {
		return graph_err.SelfEdge(from.String())
	}

	// Ensure both nodes exist in the graph.
	if g.nodes.find(from) == nil {
		return graph_err.NotExistNode(from.String())
	}
	if g.nodes.find(to) == nil {
		return graph_err.NotExistNode(to.String())
	}

	err := g.nodes.find(from).removeEdge(to)

	if err != nil {
		return err
	}

	if g.graphType == UNDIRECTED_UNWEIGHTED || g.graphType == UNDIRECTED_WEIGHTED {
		err = g.nodes.find(to).removeEdge(to)

		if err != nil {
			return err
		}
	}

	g.updated = false // Mark the graph as modified.
	g.edgeCount--     // Update edge count

	return nil
}

// FindEdge searches for an edge between two nodes in the graph and returns its distance.
//
// Parameters:
//   - from: The identifier of the source node.
//   - to: The identifier of the destination node.
//
// Returns:
//   - A pointer to the `Distance` of the edge if it exists.
//   - An error if the edge or either of the nodes does not exist, or if attempting to find a self-loop edge.
func (g *Graph) FindEdge(from, to NodeID) (*Distance, error) {
	if from == to {
		return nil, graph_err.SelfEdge(from.String())
	}

	f := g.nodes.find(from)
	t := g.nodes.find(to)

	// Ensure both nodes exist in the graph.
	if f == nil {
		return nil, graph_err.NotExistNode(from.String())
	}
	if t == nil {
		return nil, graph_err.NotExistNode(to.String())
	}

	for _, e := range f.edges {
		if e.to == to {
			return &e.distance, nil
		}
	}

	return nil, graph_err.NotExistEdge(from.String(), to.String())
}

// Matrix converts the graph to an adjacency matrix representation.
// Returns a Matrix where each element represents the distance between two nodes.
func (g *Graph) Matrix() Matrix {
	size := g.nowID
	matrix := make([][]Distance, size)

	// Initialize the matrix with infinity values.
	for i := range matrix {
		matrix[i] = make([]Distance, size)
		for j := range matrix[i] {
			matrix[i][j] = math.MaxUint
		}
	}

	// Populate the matrix with edge distances.
	for from_id, from := range g.nodes.nodes {
		for _, from_edge := range from.edges {
			matrix[from_id][from_edge.to] = from_edge.distance
		}
	}

	return matrix
}

// String returns a string representation of the Matrix.
// This method formats the matrix for easy readability:
//   - Each row of the matrix is printed on a new line.
//   - Values are separated by spaces, with "INF" used for unreachable nodes.
func (g Graph) String() string {
	result := ""

	matrix := g.Matrix()

	// Iterate over each row of the matrix.
	for _, arr := range [][]Distance(matrix) {
		// Iterate over each element in the row.
		for _, a := range arr {
			if a != INF {
				// Print the distance if it is not `INF`.
				result += fmt.Sprintf("%3d ", a)
			} else {
				// Use "INF" to represent unreachable nodes.
				result += "INF "
			}
		}

		// Add a newline at the end of each row.
		result += "\n"
	}

	return result
}

// NodeCount returns the number of nodes in the graph.
func (g Graph) NodeCount() int {
	return len(g.nodes.nodes)
}

// EdgeCount returns the number of edges in the graph.
func (g Graph) EdgeCount() int {
	return g.edgeCount
}

// Type returns the type of the graph (e.g., directed/undirected, weighted/unweighted).
func (g Graph) Type() GraphType {
	return g.graphType
}

// IsUpdated returns whether the graph has been updated since the last algorithmic computation.
func (g Graph) IsUpdated() bool {
	return g.updated
}

// Update sets the graph's updated status to true.
// This should be called after performing an algorithmic computation.
func (g *Graph) Update() {
	g.updated = true
}
