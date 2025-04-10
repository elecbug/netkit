package graph_algorithm

import (
	"sort"
	"sync"

	"github.com/elecbug/go-dspkg/graph"
	"github.com/elecbug/go-dspkg/graph/graph_type"
)

// `computePaths` calculates all shortest paths between every pair of nodes in the graph for a `Unit`.
func (u *Unit) computePaths() {
	u.shortestPaths = []Path{}
	list := u.graph.AliveNodes()
	matrix := u.graph.ToMatrix()
	weighted := u.graph.Type() == graph_type.DIRECTED_WEIGHTED || u.graph.Type() == graph_type.UNDIRECTED_WEIGHTED

	for _, start := range list {
		allPaths := computeAllShortestPathsFrom(matrix, start, weighted)
		for _, path := range allPaths {
			u.shortestPaths = append(u.shortestPaths, *path)
		}
	}

	sort.Slice(u.shortestPaths, func(i, j int) bool {
		return u.shortestPaths[i].Distance() < u.shortestPaths[j].Distance()
	})

	u.updateVersion = u.graph.Version()
}

// `computePaths` calculates all shortest paths in parallel for a `ParallelUnit`.
func (pu *ParallelUnit) computePaths() {
	pu.shortestPaths = []Path{}
	list := pu.graph.AliveNodes()
	matrix := pu.graph.ToMatrix()
	weighted := pu.graph.Type() == graph_type.DIRECTED_WEIGHTED || pu.graph.Type() == graph_type.UNDIRECTED_WEIGHTED

	type job struct {
		start graph.NodeID
	}

	jobChan := make(chan job, 100)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Worker goroutines
	for i := uint(0); i < pu.maxCore; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localResults := []Path{}

			for j := range jobChan {
				allPaths := computeAllShortestPathsFrom(matrix, j.start, weighted)
				for _, path := range allPaths {
					localResults = append(localResults, *path)
				}
				// Batch flush
				if len(localResults) >= 100 {
					mu.Lock()
					pu.shortestPaths = append(pu.shortestPaths, localResults...)
					mu.Unlock()
					localResults = localResults[:0]
				}
			}

			// Flush remaining
			if len(localResults) > 0 {
				mu.Lock()
				pu.shortestPaths = append(pu.shortestPaths, localResults...)
				mu.Unlock()
			}
		}()
	}

	// Producer
	go func() {
		for _, start := range list {
			jobChan <- job{start: start}
		}
		close(jobChan)
	}()

	wg.Wait()

	sort.Slice(pu.shortestPaths, func(i, j int) bool {
		return pu.shortestPaths[i].Distance() < pu.shortestPaths[j].Distance()
	})

	pu.updateVersion = pu.graph.Version()
}

// `computeAllShortestPathsFrom` computes all shortest paths from a single `start` node to all other nodes.
func computeAllShortestPathsFrom(matrix graph.Matrix, start graph.NodeID, weighted bool) map[graph.NodeID]*Path {
	n := len(matrix)
	dist := make([]graph.Distance, n)
	prev := make([]int, n)
	visited := make([]bool, n)

	for i := range dist {
		dist[i] = graph.INF_DISTANCE
		prev[i] = -1
	}
	dist[start] = 0

	if weighted {
		// Dijkstra's algorithm
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
	} else {
		// BFS for unweighted graph
		queue := []int{int(start)}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]

			for v := 0; v < n; v++ {
				if matrix[u][v] > 0 && dist[v] == graph.INF_DISTANCE {
					dist[v] = dist[u] + 1
					prev[v] = u
					queue = append(queue, v)
				}
			}
		}
	}

	paths := make(map[graph.NodeID]*Path)

	for end := 0; end < n; end++ {
		if graph.NodeID(end) == start || dist[end] == graph.INF_DISTANCE {
			continue
		}

		path := []graph.NodeID{}
		for at := end; at != -1; at = prev[at] {
			path = append(path, graph.NodeID(at))
		}

		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}

		paths[graph.NodeID(end)] = newPath(dist[end], path)
	}

	return paths
}
