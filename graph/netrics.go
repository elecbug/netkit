package graph

import (
	"github.com/elecbug/go-dspkg/graph/internal/algorithm"
	"github.com/elecbug/go-dspkg/graph/internal/graph"
)

// Type aliases for commonly used graph-related types from the internal packages.
type GraphType = graph.GraphType // Represents the type of graph (directed/undirected, weighted/unweighted).
type Distance = graph.Distance   // Represents the weight or distance between nodes.
type Node = graph.Node           // Represents a node in the graph.
type NodeID = graph.NodeID       // Represents the unique identifier of a node.
type Matrix = graph.Matrix       // Represents the adjacency matrix of the graph.

// Type aliases for algorithm-related structures from the internal packages.
type Unit = algorithm.Unit                 // Represents a computation unit for sequential graph algorithms.
type ParallelUnit = algorithm.ParallelUnit // Represents a computation unit for parallel graph algorithms.

// GraphParams wraps the internal graph.Graph to implement the Graph interface.
type GraphParams struct{ *graph.Graph }

// PathParams wraps the internal graph.Path to implement the Path interface.
type PathParams struct{ *graph.Path }

// Graph defines the interface for interacting with graph structures.
// It includes methods for managing nodes and edges, retrieving graph properties, and converting to computation units.
type Graph interface {
	AddNode(name string) (*Node, error)                     // Adds a new node to the graph.
	RemoveNode(identifier NodeID) error                     // Removes a node from the graph.
	FindNode(identifier NodeID) (*Node, error)              // Finds a node by its identifier.
	FindNodesByName(name string) ([]*Node, error)           // Finds all nodes with the given name.
	AddEdge(from, to NodeID) error                          // Adds an unweighted edge between two nodes.
	AddWeightEdge(from, to NodeID, distance Distance) error // Adds a weighted edge between two nodes.
	RemoveEdge(from, to NodeID) error                       // Removes an edge between two nodes.
	FindEdge(from, to NodeID) (*Distance, error)            // Finds the distance of an edge between two nodes.
	Matrix() Matrix                                         // Returns the adjacency matrix of the graph.
	String() string                                         // Returns a string representation of the graph.
	NodeCount() int                                         // Returns the number of nodes in the graph.
	EdgeCount() int                                         // Returns the number of edges in the graph.
	Type() GraphType                                        // Returns the type of the graph.
	IsUpdated() bool                                        // Checks if the graph has been updated since the last computation.
	ToUnit() *Unit                                          // Converts the graph to a Unit for sequential computation.
	ToParallelUnit(core uint) *ParallelUnit                 // Converts the graph to a ParallelUnit for parallel computation.
}

// Path defines the interface for interacting with paths in a graph.
// It includes methods to retrieve the distance and the nodes in the path.
type Path interface {
	Distance() Distance // Returns the total distance of the path.
	Nodes() []NodeID    // Returns the sequence of nodes in the path.
}

// Ensure GraphParams implements the Graph interface.
var _ Graph = (*GraphParams)(nil)

// Ensure PathParams implements the Path interface.
var _ Path = (*PathParams)(nil)

// NewGraph creates a new graph instance with the specified type and capacity.
//
// Parameters:
//   - graphType: The type of the graph (directed/undirected, weighted/unweighted).
//   - capacity: The initial capacity for nodes and edges.
//
// Returns:
//   - A Graph interface representing the new graph.
func NewGraph(graphType GraphType, capacity int) Graph {
	return &GraphParams{graph.NewGraph(graphType, capacity)}
}

// ToUnit converts the GraphParams to a sequential computation unit (Unit).
//
// Returns:
//   - A pointer to the Unit for sequential computation.
func (g *GraphParams) ToUnit() *Unit {
	return algorithm.NewUnit(g.Graph)
}

// ToParallelUnit converts the GraphParams to a parallel computation unit (ParallelUnit).
//
// Parameters:
//   - core: The number of CPU cores to use for parallel computation.
//
// Returns:
//   - A pointer to the ParallelUnit for parallel computation.
func (g *GraphParams) ToParallelUnit(core uint) *ParallelUnit {
	return algorithm.NewParallelUnit(g.Graph, core)
}

// Constants representing infinity for distances.
const INF = Distance(graph.INF)

// Constants representing graph types.
const (
	DIRECTED_UNWEIGHTED   = GraphType(graph.DIRECTED_UNWEIGHTED)   // Directed unweighted graph.
	DIRECTED_WEIGHTED     = GraphType(graph.DIRECTED_WEIGHTED)     // Directed weighted graph.
	UNDIRECTED_UNWEIGHTED = GraphType(graph.UNDIRECTED_UNWEIGHTED) // Undirected unweighted graph.
	UNDIRECTED_WEIGHTED   = GraphType(graph.UNDIRECTED_WEIGHTED)   // Undirected weighted graph.
)
