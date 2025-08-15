package g_algorithm

import (
	"runtime"
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

// ShortestPath computes a shortest path between start and end using BFS.
// It returns an empty path when no path exists.
func ShortestPath(g *graph.Graph, start, end node.ID) path.Path {
	if v, ok := cachedPaths[g.Hash()]; ok {
		return v[start][end]
	}

	if start == end {
		return *path.NewPath(start)
	}

	queue := []node.ID{start}
	visited := make(map[node.ID]bool)
	parent := make(map[node.ID]node.ID)
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		neighbors := g.GetNeighbors(current)

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)

				if neighbor == end {
					// Reconstruct path
					p := []node.ID{}

					for n := end; n != start; n = parent[n] {
						p = append([]node.ID{n}, p...)
					}

					p = append([]node.ID{start}, p...)

					return *path.NewPath(p...)
				}
			}
		}
	}

	return *path.NewPath() // No path found
}

// AllShortestPathsConcurrent computes all-pairs shortest paths using a worker pool.
func AllShortestPaths(g *graph.Graph, config *Config) map[node.ID]map[node.ID]path.Path {
	if v, ok := cachedPaths[g.Hash()]; ok {
		return v
	}

	if config == nil {
		config = &Config{}
	}

	workers := config.Workers

	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	nodes := g.GetNodes()

	// Prepare per-start row buckets with independent locks
	rows := make(map[node.ID]*row, len(nodes))

	for _, s := range nodes {
		rows[s] = &row{m: make(map[node.ID]path.Path, len(nodes)-1)}
	}

	// Create job queue
	jobs := make(chan pair, workers*2)

	var wg sync.WaitGroup

	// Worker goroutines
	workerFn := func() {
		defer wg.Done()
		for job := range jobs {
			// Compute the shortest path for (start, end)
			p := ShortestPath(g, job.start, job.end)

			// Write into the per-start row with minimal lock scope
			rS := rows[job.start]
			rS.mu.Lock()
			rS.m[job.end] = p
			rS.mu.Unlock()

			if g.IsBidirectional() {
				rE := rows[job.end]
				rE.mu.Lock()
				rE.m[job.start] = p
				rE.mu.Unlock()
			}
		}
	}

	// Start workers
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go workerFn()
	}

	// Enqueue all (start, end) pairs where start != end
	for i, s := range nodes {
		for j, e := range nodes {
			if (g.IsBidirectional() && i > j) || i == j {
				continue // Skip if bidirectional and i < j
			}

			jobs <- pair{start: s, end: e}
		}
	}

	close(jobs)

	// Wait for all jobs to finish
	wg.Wait()

	// Assemble final output (rows[*].m is already the desired inner map)
	out := make(map[node.ID]map[node.ID]path.Path, len(rows))

	for s, r := range rows {
		out[s] = r.m
	}

	cachedPaths[g.Hash()] = out

	return out
}

// row is a per-start-node bucket with its own mutex.
// It prevents concurrent writes to the inner map.
type row struct {
	mu sync.Mutex
	m  map[node.ID]path.Path
}

// pair is a single job unit (start, end) to compute.
type pair struct {
	start node.ID
	end   node.ID
}
