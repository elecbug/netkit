package analyzer

import (
	"runtime"
	"sync"

	"github.com/elecbug/netkit/v2/graph"
)

// BetweennessCentrality computes node betweenness centrality.
//
// For each source-target pair, the score of an intermediate node v is increased
// by the fraction of shortest paths between the pair that pass through v.
//
// The implementation uses cached all-pairs shortest paths computed by Analyzer.
// For undirected graphs, unordered node pairs are counted once. For directed graphs,
// ordered source-target pairs are counted.
//
// If normalization is enabled, the result follows NetworkX-compatible scaling:
//   - Undirected: 2 / ((n - 1)(n - 2))
//   - Directed:   1 / ((n - 1)(n - 2))
//
// Nodes that do not lie on any shortest path are included with score 0.
func (a *Analyzer) BetweennessCentrality() (map[graph.NodeID]float64, error) {
	res := make(map[graph.NodeID]float64)

	if a == nil || a.baseGraph == nil {
		return res, nil
	}

	if err := a.computeAllShortestPaths(); err != nil {
		return nil, err
	}

	g := a.baseGraph
	ids := g.Nodes()
	n := len(ids)

	if n < 3 {
		for _, u := range ids {
			res[u] = 0
		}
		return res, nil
	}

	workers := runtime.NumCPU()
	if a.parallelCoreCount > 0 {
		workers = a.parallelCoreCount
	}
	if workers > n {
		workers = n
	}
	if workers < 1 {
		workers = 1
	}

	normalized := true
	if a.cfg != nil && a.cfg.Betweenness != nil {
		normalized = a.cfg.Betweenness.Normalized
	}

	all := a.allShortestPaths

	type pair struct {
		s graph.NodeID
		t graph.NodeID
	}

	isUndirected := !g.IsDirected()

	jobs := make(chan pair, n)
	var wg sync.WaitGroup

	global := make(map[graph.NodeID]float64, n)
	var mu sync.Mutex

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

			for _, pth := range pathsST {
				seq := pth.Nodes()
				if len(seq) <= 2 {
					continue
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

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go workerFn()
	}

	if isUndirected {
		for i := 0; i < n; i++ {
			s := ids[i]
			for j := i + 1; j < n; j++ {
				t := ids[j]
				jobs <- pair{s: s, t: t}
			}
		}
	} else {
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

	norm := 1.0
	if normalized {
		if isUndirected {
			norm = 2.0 / float64((n-1)*(n-2))
		} else {
			norm = 1.0 / float64((n-1)*(n-2))
		}
	}

	for _, u := range ids {
		res[u] = global[u] * norm
	}

	return res, nil
}
