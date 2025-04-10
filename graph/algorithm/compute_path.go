package algorithm

import (
	"sort"
	"sync"

	"github.com/elecbug/go-dspkg/graph/graph"
)

// `computePaths` calculates all shortest paths between every pair of nodes in the graph for a Unit.
// After computation, the `shortestPaths` field in the Unit is updated and sorted by path distance in ascending order.
//
// Parameters:
//   - g: The graph to perform the computation on.
func (u *Unit) computePaths() {
	g := u.graph
	u.shortestPaths = []graph.Path{}
	n := len(g.Matrix())

	for start := graph.NodeID(0); start < graph.NodeID(n); start++ {
		for end := graph.NodeID(0); end < graph.NodeID(n); end++ {
			if start == end {
				continue
			}

			path := shortestPath(g, start, end)

			if path.Distance() != graph.INF_DISTANCE {
				u.shortestPaths = append(u.shortestPaths, *path)
			}
		}
	}

	// Sort the paths by their total distance.
	sort.Slice(u.shortestPaths, func(i, j int) bool {
		return u.shortestPaths[i].Distance() < u.shortestPaths[j].Distance()
	})

	u.updated = true
	g.Update()
}

// `computePaths` calculates all shortest paths in parallel for a ParallelUnit.
// After computation, the `shortestPaths` field in the ParallelUnit is updated and sorted by path distance in ascending order.
//
// Parameters:
//   - g: The graph to perform the computation on.
func (pu *ParallelUnit) computePaths() {
	g := pu.graph
	pu.shortestPaths = []graph.Path{}

	type to struct {
		start graph.NodeID
		end   graph.NodeID
	}

	n := len(g.Matrix())

	jobChan := make(chan to)
	resultChan := make(chan graph.Path)
	workerCount := pu.maxCore

	var wg sync.WaitGroup
	wg.Add(int(workerCount))

	// Start worker goroutines to compute paths in parallel.
	for i := uint(0); i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for job := range jobChan {
				path := shortestPath(g, job.start, job.end)

				if path.Distance() != graph.INF_DISTANCE {
					resultChan <- *path
				}
			}
		}()
	}

	// Generate jobs for every pair of nodes.
	go func() {
		for start := 0; start < n; start++ {
			for end := 0; end < n; end++ {
				if start != end {
					jobChan <- to{graph.NodeID(start), graph.NodeID(end)}
				}
			}
		}
		close(jobChan)
	}()

	// Close the result channel after all workers finish.
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from workers.
	for result := range resultChan {
		pu.shortestPaths = append(pu.shortestPaths, result)
	}

	// Sort the paths by their total distance.
	sort.Slice(pu.shortestPaths, func(i, j int) bool {
		return pu.shortestPaths[i].Distance() < pu.shortestPaths[j].Distance()
	})

	pu.updated = true
	g.Update()
}

// `shortestPath` computes the shortest path between two nodes in a graph.
//
// Parameters:
//   - g: The graph to perform the computation on.
//   - start: The starting node identifier.
//   - end: The ending node identifier.
//
// Returns:
//   - A graph.Path containing the shortest path and its total distance.
//   - If no path exists, the returned Path has distance INF and an empty node sequence.
func shortestPath(g *graph.Graph, start, end graph.NodeID) *graph.Path {
	if g.Type() == graph.DIRECTED_WEIGHTED || g.Type() == graph.UNDIRECTED_WEIGHTED {
		return weightedShortestPath(g.Matrix(), start, end)
	} else if g.Type() == graph.DIRECTED_UNWEIGHTED || g.Type() == graph.UNDIRECTED_UNWEIGHTED {
		return unweightedShortestPath(g.Matrix(), start, end)
	} else {
		return graph.NewPath(graph.INF_DISTANCE, []graph.NodeID{})
	}
}

// `weightedShortestPath` computes the shortest path between two nodes in a weighted graph.
// Uses Dijkstra's algorithm to calculate the path.
//
// Parameters:
//   - matrix: The adjacency matrix representation of the graph.
//   - start: The starting node identifier.
//   - end: The ending node identifier.
//
// Returns:
//   - A graph.Path containing the shortest path and its total distance.
func weightedShortestPath(matrix graph.Matrix, start, end graph.NodeID) *graph.Path {
	n := len(matrix)

	if int(start) >= n || int(end) >= n {
		return graph.NewPath(graph.INF_DISTANCE, []graph.NodeID{})
	}

	dist := make([]graph.Distance, n)
	prev := make([]int, n)
	visited := make([]bool, n)

	for i := range dist {
		dist[i] = graph.INF_DISTANCE
		prev[i] = -1
	}

	dist[start] = 0

	for {
		minDist := graph.INF_DISTANCE
		u := -1
		for i := 0; i < n; i++ {
			if !visited[i] && dist[i] < minDist {
				minDist = dist[i]
				u = i
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		for v := 0; v < n; v++ {
			if matrix[u][v] > 0 && !visited[v] {
				alt := dist[u] + matrix[u][v]
				if alt < dist[v] {
					dist[v] = alt
					prev[v] = u
				}
			}
		}
	}

	path := []graph.NodeID{}

	for at := int(end); at != -1; at = prev[at] {
		path = append(path, graph.NodeID(at))
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	if dist[end] == graph.INF_DISTANCE {
		return graph.NewPath(graph.INF_DISTANCE, []graph.NodeID{})
	}

	return graph.NewPath(dist[end], path)
}

// `unweightedShortestPath` computes the shortest path between two nodes in an unweighted graph.
// Uses BFS to calculate the path.
//
// Parameters:
//   - matrix: The adjacency matrix representation of the graph.
//   - start: The starting node identifier.
//   - end: The ending node identifier.
//
// Returns:
//   - A graph.Path containing the shortest path and its total distance.
func unweightedShortestPath(matrix graph.Matrix, start, end graph.NodeID) *graph.Path {
	n := len(matrix)

	if int(start) >= n || int(end) >= n {
		return graph.NewPath(graph.INF_DISTANCE, []graph.NodeID{})
	}

	dist := make([]graph.Distance, n)
	prev := make([]int, n)

	for i := range dist {
		dist[i] = graph.INF_DISTANCE
		prev[i] = -1
	}

	queue := []int{int(start)}
	dist[start] = 0

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for v := 0; v < n; v++ {
			if matrix[u][v] == 1 && dist[v] == graph.INF_DISTANCE {
				dist[v] = dist[u] + 1
				prev[v] = u
				queue = append(queue, v)
			}
		}
	}

	path := []graph.NodeID{}

	for at := int(end); at != -1; at = prev[at] {
		path = append(path, graph.NodeID(at))
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	if dist[end] == graph.INF_DISTANCE {
		return graph.NewPath(graph.INF_DISTANCE, []graph.NodeID{})
	}

	return graph.NewPath(dist[end], path)
}
