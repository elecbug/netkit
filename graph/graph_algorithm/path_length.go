package graph_algorithm

import (
	"github.com/elecbug/go-dspkg/graph"
)

// `ShortestPath` finds and returns the shortest path between two nodes in the graph for a Unit.
// Returns a `graph.Path` representing the shortest path from `from` to `to`.
// If no path exists, returns a `graph.Path` with distance `INF` and the nodes `{from, to}`.
func (u *Unit) ShortestPath(from, to graph.NodeID) Path {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	for _, p := range u.shortestPaths {
		if p.Nodes()[0] == from && p.Nodes()[len(p.Nodes())-1] == to {
			return p
		}
	}

	return *newPath(graph.INF_DISTANCE, []graph.NodeID{from, to})
}

// `ShortestPath` finds and returns the shortest path between two nodes in the graph for a ParallelUnit.
// Returns a `graph.Path` representing the shortest path from `from` to `to`.
// If no path exists, returns a `graph.Path` with distance `INF` and the nodes `{from, to}`.
func (pu *ParallelUnit) ShortestPath(from, to graph.NodeID) Path {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	for _, p := range pu.shortestPaths {
		if p.Nodes()[0] == from && p.Nodes()[len(p.Nodes())-1] == to {
			return p
		}
	}

	return *newPath(graph.INF_DISTANCE, []graph.NodeID{from, to})
}

// `AverageShortestPathLength` computes the average shortest path length in the graph.
//
// Returns:
//   - The average shortest path length as a float64.
//
// Notes:
//   - If no shortest paths are found, the function returns 0.
func (u *Unit) AverageShortestPathLength() float64 {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	var totalDistance graph.Distance = 0
	var pairCount int

	// Sum up distances for all shortest paths.
	for _, path := range u.shortestPaths {
		totalDistance += path.Distance()
		pairCount++
	}

	if pairCount == 0 {
		return 0 // Avoid division by zero if no paths exist.
	}

	return float64(totalDistance) / float64(pairCount)
}

// ParallelUnit version of `AverageShortestPathLength`.
// Computes the average shortest path length using parallel computations.
//
// Returns:
//   - The average shortest path length as a float64.
//
// Notes:
//   - If no shortest paths are found, the function returns 0.
func (pu *ParallelUnit) AverageShortestPathLength() float64 {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	var totalDistance graph.Distance = 0
	var pairCount int

	// Sum up distances for all shortest paths.
	for _, path := range pu.shortestPaths {
		totalDistance += path.Distance()
		pairCount++
	}

	if pairCount == 0 {
		return 0 // Avoid division by zero if no paths exist.
	}

	return float64(totalDistance) / float64(pairCount)
}

// `PercentileShortestPathLength` returns the shortest path length at the specified percentile.
//
// Parameters:
//   - percentile: A float64 between 0 and 1 indicating the desired percentile.
//
// Returns:
//   - The shortest path length corresponding to the given percentile.
//
// Notes:
//   - The percentile is calculated based on the sorted list of shortest paths.
//   - If the percentile is out of range, it is clamped to valid indices.
func (u *Unit) PercentileShortestPathLength(percentile float64) graph.Distance {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	// Calculate the index for the desired percentile.
	index := int(percentile * float64(len(u.shortestPaths)))

	// Clamp the index to the valid range.
	if index >= len(u.shortestPaths) {
		index = len(u.shortestPaths) - 1
	} else if index < 0 {
		index = 0
	}

	return u.shortestPaths[index].Distance()
}

// ParallelUnit version of `PercentileShortestPathLength`.
// Computes the percentile shortest path length using parallel computations.
//
// Parameters:
//   - percentile: A float64 between 0 and 1 indicating the desired percentile.
//
// Returns:
//   - The shortest path length corresponding to the given percentile.
//
// Notes:
//   - The percentile is calculated based on the sorted list of shortest paths.
//   - If the percentile is out of range, it is clamped to valid indices.
func (pu *ParallelUnit) PercentileShortestPathLength(percentile float64) graph.Distance {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	// Calculate the index for the desired percentile.
	index := int(percentile * float64(len(pu.shortestPaths)))

	// Clamp the index to the valid range.
	if index >= len(pu.shortestPaths) {
		index = len(pu.shortestPaths) - 1
	} else if index < 0 {
		index = 0
	}

	return pu.shortestPaths[index].Distance()
}
