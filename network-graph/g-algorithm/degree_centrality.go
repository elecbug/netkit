package g_algorithm

import (
	"runtime"
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// DegreeCentralityConfig (suggested to be added inside your config package)
// NetworkX parity:
//   - Undirected: degree_centrality = deg(u) / (n-1)
//   - Directed:
//       mode="total" (default) => (in(u)+out(u)) / (n-1)  (can exceed 1.0; same as nx.degree_centrality)
//       mode="in"               => in(u) / (n-1)          (nx.in_degree_centrality)
//       mode="out"              => out(u) / (n-1)         (nx.out_degree_centrality)
//
// type DegreeCentralityConfig struct {
//     // "total" | "in" | "out"
//     Mode string
// }
//
// In config.Config:
// type Config struct {
//     Workers int
//     Degree  *DegreeCentralityConfig
//     // ...
// }

// DegreeCentrality computes degree-based centrality with NetworkX-compatible semantics.
// - Undirected: deg(u)/(n-1).
// - Directed (default "total"): (in(u)+out(u))/(n-1). Use "in"/"out" for the specific variants.
// Self-loops are ignored for centrality.
func DegreeCentrality(g *graph.Graph, cfg *config.Config) map[node.ID]float64 {
	res := make(map[node.ID]float64)
	if g == nil {
		return res
	}

	ids := g.Nodes()
	n := len(ids)
	if n <= 1 {
		for _, u := range ids {
			res[u] = 0
		}
		return res
	}
	denom := float64(n - 1)

	// --- read config ---
	mode := "total"
	if cfg != nil && cfg.Degree != nil && cfg.Degree.Mode != "" {
		mode = cfg.Degree.Mode
	}

	// --- indexing ---
	idxOf := make(map[node.ID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	isUndirected := g.IsBidirectional()

	// --- build out-neighbors; ignore self-loops ---
	outs := make([][]int, n)
	for i, u := range ids {
		nbrs := g.Neighbors(u)
		if len(nbrs) == 0 {
			continue
		}
		row := make([]int, 0, len(nbrs))
		for _, v := range nbrs {
			if v == u { // ignore self-loop
				continue
			}
			row = append(row, idxOf[v])
		}
		outs[i] = row
	}

	// --- build in-neighbors for digraphs (using outs) ---
	var ins [][]int
	if isUndirected {
		// Not needed; degree uses outs directly.
	} else {
		ins = make([][]int, n)
		for u := 0; u < n; u++ {
			for _, v := range outs[u] {
				// edge u->v implies v has predecessor u
				ins[v] = append(ins[v], u)
			}
		}
	}

	// --- parallel degree counting ---
	workers := runtime.NumCPU()
	if cfg != nil && cfg.Workers > 0 {
		workers = cfg.Workers
	}
	if workers < 1 {
		workers = 1
	}
	if workers > n {
		workers = n
	}

	split := func(total, k int) [][2]int {
		if k > total {
			k = total
		}
		size := (total + k - 1) / k
		out := make([][2]int, 0, k)
		for s := 0; s < total; s += size {
			e := s + size
			if e > total {
				e = total
			}
			out = append(out, [2]int{s, e})
		}
		return out
	}
	chunks := split(n, workers)

	nums := make([]float64, n) // raw degree counts according to mode

	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for _, rg := range chunks {
		rg := rg
		go func() {
			defer wg.Done()
			for i := rg[0]; i < rg[1]; i++ {
				switch {
				case isUndirected:
					// Undirected: degree is just len(neighbors)
					nums[i] = float64(len(outs[i]))
				default:
					switch mode {
					case "in":
						nums[i] = float64(len(ins[i]))
					case "out":
						nums[i] = float64(len(outs[i]))
					default: // "total"
						nums[i] = float64(len(ins[i]) + len(outs[i]))
					}
				}
			}
		}()
	}
	wg.Wait()

	// --- normalize by (n-1) and fill map ---
	for i, u := range ids {
		res[u] = nums[i] / denom
	}
	return res
}
