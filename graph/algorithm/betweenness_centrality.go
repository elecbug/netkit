package algorithm

import (
	"runtime"
	"sync"

	"github.com/elecbug/netkit/graph"
)

// BetweennessCentrality computes betweenness centrality using cached all shortest paths.
// - Uses AllShortestPaths(g, cfg) (assumed cached/fast) to enumerate all shortest paths.
// - For each pair (s,t), each interior node on a shortest path gets 1/|SP(s,t)| credit.
// - Undirected graphs enqueue only i<j pairs (no double counting).
// - Normalization matches NetworkX: undirected => 2/((n-1)(n-2)), directed => 1/((n-1)(n-2)).
func BetweennessCentrality(g *graph.Graph, cfg *Config) map[graph.NodeID]float64 {
	res := make(map[graph.NodeID]float64)
	if g == nil {
		return res
	}

	ids := g.Nodes()
	n := len(ids)
	if n < 3 {
		// With fewer than 3 nodes, BC is zero for every node.
		for _, u := range ids {
			res[u] = 0
		}
		return res
	}

	// Read/normalize worker count.
	workers := 1
	normalized := true
	if cfg != nil && cfg.Workers > 0 {
		workers = cfg.Workers

		if cfg.Betweenness != nil {
			normalized = cfg.Betweenness.Normalized
		}
	} else {
		workers = runtime.NumCPU()
	}
	if workers > n {
		workers = n
	}

	// Use cached all-pairs shortest paths.
	// Type: map[graph.NodeID]map[graph.NodeID][]path.Path
	all := AllShortestPaths(g, cfg)

	// Build an index for stable iteration and pair generation.
	idxOf := make(map[graph.NodeID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	type pair struct{ s, t graph.NodeID }
	isUndirected := g.IsBidirectional()

	// Generate all (s,t) jobs.
	jobs := make(chan pair, n)
	var wg sync.WaitGroup

	// Global accumulator with lock; each worker keeps a local map to minimize contention.
	global := make(map[graph.NodeID]float64, n)
	var mu sync.Mutex

	// Worker: consume pairs and accumulate contributions into a local map, then merge.
	workerFn := func() {
		defer wg.Done()
		local := make(map[graph.NodeID]float64, n)

		for job := range jobs {
			row, ok := all[job.s]
			if !ok {
				continue
			}
			pathsST, ok := row[job.t]
			if !ok || len(pathsST) == 0 {
				continue
			}
			den := float64(len(pathsST))

			// For each shortest path s->...->t, every interior node gets 1/den.
			for _, pth := range pathsST {
				seq := pth.Nodes() // []graph.NodeID; interior nodes are [1 : len-1)
				if len(seq) <= 2 {
					continue // no interior node
				}
				for k := 1; k < len(seq)-1; k++ {
					v := seq[k]
					local[v] += 1.0 / den
				}
			}
		}

		if len(local) > 0 {
			mu.Lock()
			for v, val := range local {
				global[v] += val
			}
			mu.Unlock()
		}
	}

	// Start workers.
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go workerFn()
	}

	// Enqueue jobs.
	if isUndirected {
		// Only upper triangle (i<j) to avoid double counting undirected pairs.
		for i := 0; i < n; i++ {
			s := ids[i]
			for j := i + 1; j < n; j++ {
				t := ids[j]
				jobs <- pair{s: s, t: t}
			}
		}
	} else {
		// Directed: all ordered pairs (i!=j).
		for i := 0; i < n; i++ {
			s := ids[i]
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				t := ids[j]
				jobs <- pair{s: s, t: t}
			}
		}
	}

	close(jobs)
	wg.Wait()

	if normalized {
		// Normalization matching NetworkX
		// Undirected: factor = 2/((n-1)(n-2))
		// Directed:   factor = 1/((n-1)(n-2))
		var norm float64
		if isUndirected {
			norm = 2.0 / float64((n-1)*(n-2))
		} else {
			norm = 1.0 / float64((n-1)*(n-2))
		}

		// Write normalized results; include nodes that never appeared (zero).
		for _, u := range ids {
			res[u] = global[u] * norm
		}

	}

	return res
}
