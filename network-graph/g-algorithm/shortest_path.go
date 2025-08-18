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

// AllShortestPaths finds all shortest paths between all pairs of nodes in a graph.
func AllShortestPaths(g *graph.Graph, cfg *config.Config) path.GraphPaths {
	gh := g.Hash()

	cacheMu.RLock()

	if v, ok := cachedAllShortestPaths[gh]; ok {
		cacheMu.RUnlock()
		return v
	}

	cacheMu.RUnlock()

	if cfg == nil {
		cfg = &config.Config{}
	}

	workers := cfg.Workers

	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	nodes := g.Nodes()
	n := len(nodes)

	rows := make(map[node.ID]*row, n)

	for _, s := range nodes {
		rows[s] = &row{m: make(map[node.ID][]path.Path, n-1)}
	}

	jobs := make(chan pair, workers*2)

	var wg sync.WaitGroup
	isUndirected := g.IsBidirectional()

	workerFn := func() {
		defer wg.Done()
		for job := range jobs {

			p := allShortestPathsBFS(g, job.start, job.end)

			rS := rows[job.start]
			rS.mu.Lock()
			rS.m[job.end] = p
			rS.mu.Unlock()

			if isUndirected {
				rE := rows[job.end]
				rE.mu.Lock()
				rE.m[job.start] = p
				rE.mu.Unlock()
			}
		}
	}

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go workerFn()
	}

	for i, s := range nodes {
		for j, e := range nodes {
			if i == j {
				continue
			}

			if isUndirected && i > j {
				continue
			}

			jobs <- pair{start: s, end: e}
		}
	}

	close(jobs)

	wg.Wait()

	out := make(path.GraphPaths, len(rows))

	for s, r := range rows {
		out[s] = r.m
	}

	for _, id := range nodes {
		out[id][id] = []path.Path{*path.NewSelf(id)}
	}

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

type row struct {
	mu sync.Mutex
	m  map[node.ID][]path.Path
}

type pair struct {
	start node.ID
	end   node.ID
}
