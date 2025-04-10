package graph_algorithm

import (
	"sync"

	"github.com/elecbug/go-dspkg/graph"
)

// `GlobalEfficiency` computes the global efficiency of a graph using a Unit.
// Global efficiency is the average inverse shortest path length for all node pairs.
//
// Returns:
//   - The global efficiency as a float64.
func (u *Unit) GlobalEfficiency() float64 {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	var totalEfficiency float64
	var pairCount int

	// Iterate through all shortest paths to compute the efficiency
	for _, path := range u.shortestPaths {
		if path.Distance() != graph.INF_DISTANCE && path.Distance() > 0 {
			totalEfficiency += 1.0 / float64(path.Distance())
			pairCount++
		}
	}

	// Calculate average efficiency
	if pairCount == 0 {
		return 0.0
	}
	return totalEfficiency / float64(pairCount)
}

// `GlobalEfficiency` computes the global efficiency of a graph using a ParallelUnit.
// The computation is performed in parallel for better performance.
//
// Returns:
//   - The global efficiency as a float64.
func (pu *ParallelUnit) GlobalEfficiency() float64 {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	var totalEfficiency float64
	var pairCount int

	// Iterate through all shortest paths to compute the efficiency
	for _, path := range pu.shortestPaths {
		if path.Distance() != graph.INF_DISTANCE && path.Distance() > 0 {
			totalEfficiency += 1.0 / float64(path.Distance())
			pairCount++
		}
	}

	// Calculate average efficiency
	if pairCount == 0 {
		return 0.0
	}
	return totalEfficiency / float64(pairCount)
}

// `LocalEfficiency` computes the local efficiency of each node in the graph using a Unit.
// Local efficiency measures how well the neighbors of a node are connected.
//
// Returns:
//   - A map where the keys are node identifiers and the values are the local efficiency scores.
func (u *Unit) LocalEfficiency() map[graph.NodeID]float64 {
	g := u.graph

	if g.Version() != u.updateVersion {
		u.computePaths()
	}

	localEfficiency := make(map[graph.NodeID]float64)

	// Group paths by their starting node
	pathsBySource := make(map[graph.NodeID][]Path)
	for _, path := range u.shortestPaths {
		if len(path.Nodes()) > 0 {
			source := path.Nodes()[0]
			pathsBySource[source] = append(pathsBySource[source], path)
		}
	}

	// Compute local efficiency for each node
	for node, paths := range pathsBySource {
		neighbors := make(map[graph.NodeID]bool)

		// Identify neighbors of the current node
		for _, path := range paths {
			if len(path.Nodes()) > 1 {
				neighbors[path.Nodes()[1]] = true
			}
		}

		neighborList := make([]graph.NodeID, 0, len(neighbors))
		for neighbor := range neighbors {
			neighborList = append(neighborList, neighbor)
		}

		k := len(neighborList)
		if k < 2 {
			// Local efficiency is undefined for nodes with fewer than 2 neighbors.
			localEfficiency[node] = 0.0
			continue
		}

		// Calculate efficiency among neighbors
		totalEfficiency := 0.0
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				for _, path := range u.shortestPaths {
					if len(path.Nodes()) > 1 && path.Nodes()[0] == neighborList[i] && path.Nodes()[len(path.Nodes())-1] == neighborList[j] {
						if path.Distance() != graph.INF_DISTANCE && path.Distance() > 0 {
							totalEfficiency += 1.0 / float64(path.Distance())
						}
					}
				}
			}
		}

		// Normalize by the number of possible connections
		localEfficiency[node] = totalEfficiency / float64(k*(k-1))
	}

	return localEfficiency
}

// `LocalEfficiency` computes the local efficiency of each node in the graph using a ParallelUnit.
// The computation is performed in parallel for better performance.
//
// Returns:
//   - A map where the keys are node identifiers and the values are the local efficiency scores.
func (pu *ParallelUnit) LocalEfficiency() map[graph.NodeID]float64 {
	g := pu.graph

	if g.Version() != pu.updateVersion {
		pu.computePaths()
	}

	localEfficiency := make(map[graph.NodeID]float64)
	efficiencyChan := make(chan struct {
		node  graph.NodeID
		value float64
	}, len(pu.shortestPaths))
	var wg sync.WaitGroup

	// Group paths by their starting node
	pathsBySource := make(map[graph.NodeID][]Path)
	for _, path := range pu.shortestPaths {
		if len(path.Nodes()) > 0 {
			source := path.Nodes()[0]
			pathsBySource[source] = append(pathsBySource[source], path)
		}
	}

	// Compute local efficiency for each node in parallel
	for node, paths := range pathsBySource {
		wg.Add(1)
		go func(node graph.NodeID, paths []Path) {
			defer wg.Done()

			neighbors := make(map[graph.NodeID]bool)
			for _, path := range paths {
				if len(path.Nodes()) > 1 {
					neighbors[path.Nodes()[1]] = true
				}
			}

			neighborList := make([]graph.NodeID, 0, len(neighbors))
			for neighbor := range neighbors {
				neighborList = append(neighborList, neighbor)
			}

			k := len(neighborList)
			if k < 2 {
				efficiencyChan <- struct {
					node  graph.NodeID
					value float64
				}{node: node, value: 0.0}
				return
			}

			totalEfficiency := 0.0
			for i := 0; i < k; i++ {
				for j := i + 1; j < k; j++ {
					for _, path := range pu.shortestPaths {
						if len(path.Nodes()) > 1 && path.Nodes()[0] == neighborList[i] && path.Nodes()[len(path.Nodes())-1] == neighborList[j] {
							if path.Distance() != graph.INF_DISTANCE && path.Distance() > 0 {
								totalEfficiency += 1.0 / float64(path.Distance())
							}
						}
					}
				}
			}

			efficiencyChan <- struct {
				node  graph.NodeID
				value float64
			}{node: node, value: totalEfficiency / float64(k*(k-1))}
		}(node, paths)
	}

	// Close the channel after all goroutines complete
	go func() {
		wg.Wait()
		close(efficiencyChan)
	}()

	// Collect results
	for res := range efficiencyChan {
		localEfficiency[res.node] = res.value
	}

	return localEfficiency
}
