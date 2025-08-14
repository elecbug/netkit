package graph_algorithm

// `Diameter` computes the diameter of the graph for a `Unit`.
// The diameter is defined as the longest shortest path between any two nodes in the graph.
func (u *Unit) Diameter() Path {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	// The diameter corresponds to the last (longest) path in the sorted shortestPaths slice.
	if len(u.shortestPaths) == 0 {
		return Path{} // Return an empty path if no paths are computed
	}

	// The diameter corresponds to the last (longest) path in the sorted shortestPaths slice.
	return u.shortestPaths[len(u.shortestPaths)-1]
}

// `Diameter` computes the diameter of the graph for a `ParallelUnit`.
func (pu *ParallelUnit) Diameter() Path {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	// The diameter corresponds to the last (longest) path in the sorted shortestPaths slice.
	if len(pu.shortestPaths) == 0 {
		return Path{} // Return an empty path if no paths are computed
	}

	return pu.shortestPaths[len(pu.shortestPaths)-1]
}
