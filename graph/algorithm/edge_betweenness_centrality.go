package algorithm

import (
	"runtime"
	"sync"

	"github.com/elecbug/netkit/graph"
)

func makeEdgeKey(u, v graph.NodeID, undirected bool) (graph.NodeID, graph.NodeID) {
	if undirected && v < u {
		u, v = v, u
	}

	return u, v
}

// EdgeBetweennessCentrality computes edge betweenness centrality (unweighted)
// compatible with NetworkX's nx.edge_betweenness_centrality.
//
// Algorithm:
//   - Brandes (2001) single-source shortest paths with dependency accumulation,
//     extended for edges.
//
// Parallelization:
//   - Sources (s) are split across a worker pool (cfg.Workers).
//
// Normalization (normalized=true):
//   - Directed:   multiply by 1 / ((n-1)*(n-2))
//   - Undirected: multiply by 2 / ((n-1)*(n-2))
//
// Additionally, undirected results are divided by 2 to correct for double counting
// in Brandes accumulation (same practice as NetworkX).
//
// Returns:
//   - map[graph.NodeID]map[graph.NodeID]float64 where:
//   - Undirected: key is canonical [min(u,v), max(u,v)]
//   - Directed:   key is (u,v) ordered
func EdgeBetweennessCentrality(g *graph.Graph, cfg *Config) map[graph.NodeID]map[graph.NodeID]float64 {
	out := make(map[graph.NodeID]map[graph.NodeID]float64)
	if g == nil {
		return out
	}

	// ----- config defaults (match NetworkX) -----
	workers := runtime.NumCPU()
	normalized := true
	if cfg != nil {
		if cfg.Workers > 0 {
			workers = cfg.Workers
		}
		if cfg.EdgeBetweenness != nil {
			normalized = cfg.EdgeBetweenness.Normalized
		}
	}
	if workers < 1 {
		workers = 1
	}

	ids := g.Nodes()
	n := len(ids)
	if n == 0 {
		return out
	}
	isUndirected := g.IsBidirectional()

	// Pre-initialize all existing edges to 0 in the output map,
	// so the result contains every graph edge (like NetworkX).
	for _, u := range ids {
		for _, v := range g.Neighbors(u) {
			if u == v {
				continue // skip self-loops
			}
			if isUndirected && v < u {
				// store only once for undirected
				continue
			}
			u, v := makeEdgeKey(u, v, isUndirected)

			if out[u] == nil {
				out[u] = make(map[graph.NodeID]float64)
			}

			out[u][v] = 0.0
		}
	}

	// ----- worker pool over source nodes -----
	type job struct{ s graph.NodeID }
	jobs := make(chan job, n)

	var mu sync.Mutex
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		// Local accumulator to reduce lock contention
		local := make(map[graph.NodeID]map[graph.NodeID]float64, 64)

		for jb := range jobs {
			s := jb.s

			// Brandes data structures
			stack := make([]graph.NodeID, 0, n)
			preds := make(map[graph.NodeID][]graph.NodeID, n)
			sigma := make(map[graph.NodeID]float64, n)
			dist := make(map[graph.NodeID]int, n)

			for _, v := range ids {
				dist[v] = -1
			}
			sigma[s] = 1.0
			dist[s] = 0

			// BFS (unweighted)
			q := []graph.NodeID{s}
			for len(q) > 0 {
				v := q[0]
				q = q[1:]
				stack = append(stack, v)

				for _, w := range g.Neighbors(v) {
					// Discover w?
					if dist[w] < 0 {
						dist[w] = dist[v] + 1
						q = append(q, w)
					}
					// Is (v,w) on a shortest path from s?
					if dist[w] == dist[v]+1 {
						sigma[w] += sigma[v]
						preds[w] = append(preds[w], v)
					}
				}
			}

			// Dependency accumulation
			delta := make(map[graph.NodeID]float64, n)
			for len(stack) > 0 {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				for _, v := range preds[w] {
					if sigma[w] == 0 {
						continue
					}
					c := (sigma[v] / sigma[w]) * (1.0 + delta[w])
					// Edge (v,w) gets c
					eu, ev := makeEdgeKey(v, w, isUndirected)

					if local[eu] == nil {
						local[eu] = make(map[graph.NodeID]float64)
					}

					local[eu][ev] += c
					// Propagate to node dependency
					delta[v] += c
				}
			}
		}

		// Merge local into global
		if len(local) > 0 {
			mu.Lock()
			for e, val := range local {
				for k, v := range val {
					if e == k {
						continue // skip self-loops
					}

					out[e][k] += v
				}
			}
			mu.Unlock()
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}
	for _, s := range ids {
		jobs <- job{s: s}
	}
	close(jobs)
	wg.Wait()

	// ----- Undirected: correct double counting -----
	if isUndirected {
		for e := range out {
			for k := range out[e] {
				out[e][k] *= 0.5
			}
		}
	}

	// ----- Normalization (match NetworkX semantics) -----
	if normalized && n > 2 {
		var scale float64
		if isUndirected {
			// 2 / ((n-1)(n-2))
			scale = 2.0 / (float64(n-1) * float64(n-2))
		} else {
			// 1 / ((n-1)(n-2))
			scale = 1.0 / (float64(n-1) * float64(n-2))
		}
		for e := range out {
			for k := range out[e] {
				out[e][k] *= scale
			}
		}
	}

	return out
}
