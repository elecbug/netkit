package analyzer

import (
	"runtime"
	"sync"

	"github.com/elecbug/netkit/v2/graph"
)

// EdgeBetweennessCentrality computes edge betweenness centrality.
//
// The implementation uses Brandes' unweighted shortest-path dependency
// accumulation algorithm extended to edges.
//
// For undirected graphs, edge keys are canonicalized as (min(u,v), max(u,v)).
// For directed graphs, edge keys preserve direction as (u,v).
//
// Self-loops are ignored. Existing graph edges are pre-initialized with score 0
// so that the returned map contains every edge.
//
// If normalization is enabled, the result follows NetworkX-compatible scaling:
//   - Undirected: 2 / ((n - 1)(n - 2))
//   - Directed:   1 / ((n - 1)(n - 2))
func (a *Analyzer) EdgeBetweennessCentrality() (map[graph.NodeID]map[graph.NodeID]float64, error) {
	out := make(map[graph.NodeID]map[graph.NodeID]float64)

	if a == nil || a.baseGraph == nil {
		return out, nil
	}

	g := a.baseGraph
	ids := g.Nodes()
	n := len(ids)

	if n == 0 {
		return out, nil
	}

	workers := runtime.NumCPU()
	if a.parallelCoreCount > 0 {
		workers = a.parallelCoreCount
	}
	if workers < 1 {
		workers = 1
	}
	if workers > n {
		workers = n
	}

	normalized := true
	if a.cfg != nil && a.cfg.EdgeBetweenness != nil {
		normalized = a.cfg.EdgeBetweenness.Normalized
	}

	isUndirected := !g.IsDirected()

	neighbors := func(u graph.NodeID) []graph.NodeID {
		uNode, err := g.Node(u)
		if err != nil {
			return nil
		}
		return uNode.Neighbors()
	}

	for _, u := range ids {
		for _, v := range neighbors(u) {
			if u == v {
				continue
			}

			if isUndirected && v < u {
				continue
			}

			eu, ev := makeEdgeKey(u, v, isUndirected)

			if out[eu] == nil {
				out[eu] = make(map[graph.NodeID]float64)
			}

			out[eu][ev] = 0
		}
	}

	type job struct {
		s graph.NodeID
	}

	jobs := make(chan job, n)

	var mu sync.Mutex
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		local := make(map[graph.NodeID]map[graph.NodeID]float64, 64)

		for jb := range jobs {
			s := jb.s

			stack := make([]graph.NodeID, 0, n)
			preds := make(map[graph.NodeID][]graph.NodeID, n)
			sigma := make(map[graph.NodeID]float64, n)
			dist := make(map[graph.NodeID]int, n)

			for _, v := range ids {
				dist[v] = -1
			}

			sigma[s] = 1
			dist[s] = 0

			q := []graph.NodeID{s}
			for len(q) > 0 {
				v := q[0]
				q = q[1:]
				stack = append(stack, v)

				for _, w := range neighbors(v) {
					if dist[w] < 0 {
						dist[w] = dist[v] + 1
						q = append(q, w)
					}

					if dist[w] == dist[v]+1 {
						sigma[w] += sigma[v]
						preds[w] = append(preds[w], v)
					}
				}
			}

			delta := make(map[graph.NodeID]float64, n)

			for len(stack) > 0 {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				for _, v := range preds[w] {
					if sigma[w] == 0 {
						continue
					}

					c := (sigma[v] / sigma[w]) * (1 + delta[w])
					eu, ev := makeEdgeKey(v, w, isUndirected)

					if local[eu] == nil {
						local[eu] = make(map[graph.NodeID]float64)
					}

					local[eu][ev] += c
					delta[v] += c
				}
			}
		}

		if len(local) > 0 {
			mu.Lock()
			for eu, row := range local {
				if out[eu] == nil {
					out[eu] = make(map[graph.NodeID]float64)
				}

				for ev, val := range row {
					if eu == ev {
						continue
					}

					out[eu][ev] += val
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

	if isUndirected {
		for eu := range out {
			for ev := range out[eu] {
				out[eu][ev] *= 0.5
			}
		}
	}

	if normalized && n > 2 {
		scale := 1.0 / (float64(n-1) * float64(n-2))
		if isUndirected {
			scale = 2.0 / (float64(n-1) * float64(n-2))
		}

		for eu := range out {
			for ev := range out[eu] {
				out[eu][ev] *= scale
			}
		}
	}

	return out, nil
}

// makeEdgeKey returns a canonical edge key (u,v) for the given nodes.
// For undirected graphs, it ensures u <= v to avoid double counting.
func makeEdgeKey(u, v graph.NodeID, undirected bool) (graph.NodeID, graph.NodeID) {
	if undirected && v < u {
		u, v = v, u
	}

	return u, v
}
