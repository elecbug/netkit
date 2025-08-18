package g_algorithm

import (
	"runtime"
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

// ShortestPaths finds all shortest paths between two nodes in a graph.
func ShortestPaths(g *graph.Graph, start, end node.ID) []path.Path {
	gh := g.Hash()

	cacheMu.RLock()

	if byStart, ok := cachedAllShortestPaths[gh]; ok {
		if byEnd, ok2 := byStart[start]; ok2 {
			if paths, ok3 := byEnd[end]; ok3 {
				cacheMu.RUnlock()
				return paths
			}
		}
	}

	cacheMu.RUnlock()

	res := allShortestPathsBFS(g, start, end)

	cacheMu.Lock()

	if _, ok := cachedAllShortestPaths[gh]; !ok {
		cachedAllShortestPaths[gh] = make(map[node.ID]map[node.ID][]path.Path)
	}

	if _, ok := cachedAllShortestPaths[gh][start]; !ok {
		cachedAllShortestPaths[gh][start] = make(map[node.ID][]path.Path)
	}

	if _, exists := cachedAllShortestPaths[gh][start][end]; !exists {
		cachedAllShortestPaths[gh][start][end] = res
	}

	cacheMu.Unlock()

	return res
}

// AllShortestPaths computes all-pairs shortest paths while keeping the same return structure (path.GraphPaths).
// Performance improvements over the (s,t)-pair BFS approach:
//   - Run exactly one BFS per source node (O(n*(m+n)) instead of O(n^2*(m+n)) in the worst case).
//   - Reconstruct all shortest paths to every target using predecessors (no repeated BFS).
//   - Use memoization to enumerate all s->u shortest paths once and reuse for all targets.
//
// Notes:
//   - For undirected graphs we fill symmetric entries [t][s] with the SAME slice reference as [s][t]
//     (matching prior behavior and saving work). If you need reversed node order per entry, that can be changed,
//     but be aware of the extra cost.
//   - Self paths [u][u] are set to a single self path.
func AllShortestPaths(g *graph.Graph, cfg *config.Config) path.GraphPaths {
	if g == nil {
		return path.GraphPaths{}
	}

	gh := g.Hash()

	// Check cache first
	cacheMu.RLock()
	if v, ok := cachedAllShortestPaths[gh]; ok {
		cacheMu.RUnlock()
		return v
	}
	cacheMu.RUnlock()

	// Workers
	workers := runtime.NumCPU()
	if cfg != nil && cfg.Workers > 0 {
		workers = cfg.Workers
	}
	if workers < 1 {
		workers = 1
	}

	ids := g.Nodes()
	n := len(ids)
	if n == 0 {
		return path.GraphPaths{}
	}

	// Precompute adjacency lists once to avoid per-step allocations in g.Neighbors.
	adj := make(map[node.ID][]node.ID, n)
	for _, u := range ids {
		adj[u] = g.Neighbors(u)
	}

	// Per-source row buckets with independent locks.
	type row struct {
		mu sync.Mutex
		m  map[node.ID][]path.Path
	}
	rows := make(map[node.ID]*row, n)
	for _, s := range ids {
		rows[s] = &row{m: make(map[node.ID][]path.Path, n-1)}
	}

	isUndirected := g.IsBidirectional()

	// Jobs are source nodes (one BFS per source).
	srcJobs := make(chan node.ID, workers*2)

	var wg sync.WaitGroup
	wg.Add(workers)

	// Worker: run BFS once from s, then reconstruct paths to every t using preds with memoization.
	goBFSWorker := func() {
		defer wg.Done()

		for s := range srcJobs {
			// --- BFS from s to get dist and preds on shortest-path DAG ---
			dist := make(map[node.ID]int, n)
			preds := make(map[node.ID][]node.ID, n)

			for _, u := range ids {
				dist[u] = -1
			}
			dist[s] = 0

			q := []node.ID{s}
			for len(q) > 0 {
				v := q[0]
				q = q[1:]
				dv := dist[v]

				for _, w := range adj[v] {
					if dist[w] < 0 {
						// First time discovered
						dist[w] = dv + 1
						preds[w] = append(preds[w], v)
						q = append(q, w)
					} else if dist[w] == dv+1 {
						// Also a predecessor on some shortest path
						preds[w] = append(preds[w], v)
					}
				}
			}

			// --- Memoized enumeration of ALL shortest paths s->u for all u reachable ---
			// returns list of sequences (each sequence is []node.ID from s to x).
			memo := make(map[node.ID][][]node.ID, n)
			var enum func(u node.ID) [][]node.ID
			enum = func(u node.ID) [][]node.ID {
				if paths, ok := memo[u]; ok {
					return paths
				}
				if u == s {
					out := [][]node.ID{{s}}
					memo[u] = out
					return out
				}
				var out [][]node.ID
				for _, p := range preds[u] {
					if dist[p] >= 0 && dist[p] == dist[u]-1 {
						pfxs := enum(p)
						for _, pfx := range pfxs {
							seq := make([]node.ID, len(pfx)+1)
							copy(seq, pfx)
							seq[len(pfx)] = u
							out = append(out, seq)
						}
					}
				}
				memo[u] = out
				return out
			}

			// Build result slice for this source s and write into rows.
			rS := rows[s]
			loc := make(map[node.ID][]path.Path, n-1)

			for _, t := range ids {
				if t == s {
					continue
				}
				if dist[t] <= 0 {
					// No path or same node: skip. Self path is handled after all workers finish.
					continue
				}
				seqs := enum(t)
				if len(seqs) == 0 {
					continue
				}
				pp := make([]path.Path, 0, len(seqs))
				for _, seq := range seqs {
					pp = append(pp, *path.New(seq...))
				}
				loc[t] = pp

				// For undirected graphs, mirror into [t][s] with the same slice reference (as previous code did).
				if isUndirected {
					rT := rows[t]
					rT.mu.Lock()
					// Share the same slice to save work and memory (matches previous semantics).
					rT.m[s] = pp
					rT.mu.Unlock()
				}
			}

			// Commit this source row in one lock.
			rS.mu.Lock()
			for t, pp := range loc {
				rS.m[t] = pp
			}
			rS.mu.Unlock()
		}
	}

	for i := 0; i < workers; i++ {
		go goBFSWorker()
	}
	for _, s := range ids {
		srcJobs <- s
	}
	close(srcJobs)
	wg.Wait()

	// Assemble final output
	out := make(path.GraphPaths, n)
	for s, r := range rows {
		out[s] = r.m
	}

	// Ensure self paths exist
	for _, u := range ids {
		out[u][u] = []path.Path{*path.NewSelf(u)}
	}

	// Cache the result
	cacheMu.Lock()
	cachedAllShortestPaths[gh] = out
	cacheMu.Unlock()

	return out
}

// allShortestPathsBFS finds all shortest paths between two nodes in a graph using BFS.
func allShortestPathsBFS(g *graph.Graph, start, end node.ID) []path.Path {
	if start == end {
		return []path.Path{*path.New(start)}
	}

	queue := []node.ID{start}
	dist := make(map[node.ID]int)
	dist[start] = 0
	preds := make(map[node.ID][]node.ID)
	targetDist := -1

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		if targetDist >= 0 && dist[v] >= targetDist {
			continue
		}

		for _, w := range g.Neighbors(v) {
			_, seen := dist[w]

			if !seen {
				dist[w] = dist[v] + 1
				preds[w] = append(preds[w], v)
				queue = append(queue, w)

				if w == end {
					targetDist = dist[w]
				}

				continue
			}

			if dist[w] == dist[v]+1 {
				preds[w] = append(preds[w], v)
			}
		}
	}

	if targetDist < 0 {
		return []path.Path{}
	}

	var all [][]node.ID
	cur := []node.ID{end}

	var dfs func(u node.ID)
	dfs = func(u node.ID) {
		if u == start {
			seq := make([]node.ID, len(cur))

			for i := range cur {
				seq[i] = cur[len(cur)-1-i]
			}

			all = append(all, seq)

			return
		}

		for _, p := range preds[u] {
			cur = append(cur, p)
			dfs(p)
			cur = cur[:len(cur)-1]
		}
	}

	dfs(end)

	res := make([]path.Path, 0, len(all))

	for _, seq := range all {
		res = append(res, *path.New(seq...))
	}
	return res
}
