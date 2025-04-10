package graph_algorithm

// `Diameter` computes the diameter of the graph for a Unit.
// The diameter is defined as the longest shortest path between any two nodes in the graph.
//
// Returns:
//   - A graph.Path representing the longest shortest path in the graph.
//
// Notes:
//   - If the graph or the Unit has been updated, shortest paths are recomputed.
func (u *Unit) Diameter() Path {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	// The diameter corresponds to the last (longest) path in the sorted shortestPaths slice.
	return u.shortestPaths[len(u.shortestPaths)-1]
}

// `Diameter` computes the diameter of the graph for a ParallelUnit.
//
// Returns:
//   - A graph.Path representing the longest shortest path in the graph.
//
// Notes:
//   - If the graph or the ParallelUnit has been updated, shortest paths are recomputed in parallel.
func (pu *ParallelUnit) Diameter() Path {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	// The diameter corresponds to the last (longest) path in the sorted shortestPaths slice.
	return pu.shortestPaths[len(pu.shortestPaths)-1]
}
