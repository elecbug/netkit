package algorithm

import (
	"github.com/elecbug/go-dspkg/graph/graph"
)

// `Unit` represents a computation unit for graph algorithms.
// It stores shortest paths within the graph and performs computations.
//
// Fields:
//   - shortestPaths: A slice of all shortest paths in the graph, sorted by distance in ascending order.
//   - graph: A reference to the graph on which computations are performed.
type Unit struct {
	shortestPaths []graph.Path // Stores the shortest paths for the graph, sorted by distance in ascending order.
	graph         *graph.Graph // A reference to the graph associated with this computation unit.
	updated       bool         // Update information for shortest paths
}

// `ParallelUnit` extends Unit for parallel computation of graph algorithms.
// It supports dividing tasks across multiple CPU cores for better performance.
//
// Fields:
//   - Unit: Embeds the Unit structure for shared functionality.
//   - maxCore: The maximum number of CPU cores to use for parallel computation.
type ParallelUnit struct {
	Unit         // Embeds the Unit structure for shared functionality.
	maxCore uint // Maximum number of CPU cores to use for parallel computation.
}

// `NewUnit` creates and initializes a new Unit instance.
//
// Parameters:
//   - g: The graph to associate with this computation unit.
//
// Returns:
//   - A pointer to the newly created Unit.
func NewUnit(g *graph.Graph) *Unit {
	return &Unit{
		shortestPaths: make([]graph.Path, g.EdgeCount()), // Initialize with an empty slice of paths.
		graph:         g,                                 // Associate the graph with this Unit.
		updated:       false,
	}
}

// `NewParallelUnit` creates and initializes a new ParallelUnit instance.
//
// Parameters:
//   - g: The graph to associate with this computation unit.
//   - core: The maximum number of CPU cores to use for parallel computations.
//
// Returns:
//   - A pointer to the newly created ParallelUnit.
func NewParallelUnit(g *graph.Graph, core uint) *ParallelUnit {
	return &ParallelUnit{
		Unit:    *NewUnit(g), // Initialize the embedded Unit structure.
		maxCore: core,        // Set the maximum number of cores for parallel processing.
	}
}
