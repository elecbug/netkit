package g_algorithm

import (
	"math"
	"runtime"
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// PageRank computes PageRank using a parallel power-iteration that mirrors NetworkX semantics:
//   - Teleport and dangling redistribution match NetworkX (personalization used for both by default).
//   - Convergence check uses L1 error before any per-iteration normalization: sum(|x - x_last|) < n*tol.
//   - Result is L1-normalized once after convergence (or after max-iter).
//   - For directed graphs, Reverse=false is the standard PageRank (incoming influence).
//     Reverse=true flips direction (treat in-neighbors as outs).
//   - For undirected graphs, neighbors are used as outs (degree-based normalization).
func PageRank(g *graph.Graph, cfg *config.Config) map[node.ID]float64 {
	res := make(map[node.ID]float64)
	if g == nil {
		return res
	}

	// ----- read config & defaults -----
	workers := runtime.NumCPU()
	alpha := 0.85
	maxIter := 1000 // a bit larger than NX default to ensure tighter match
	tol := 1e-6
	reverse := false
	var pers *map[node.ID]float64
	var dang *map[node.ID]float64

	if cfg != nil {
		if cfg.Workers > 0 {
			workers = cfg.Workers
		}
		if pr := cfg.PageRank; pr != nil {
			if pr.Alpha > 0 && pr.Alpha < 1 {
				alpha = pr.Alpha
			}
			if pr.MaxIter > 0 {
				maxIter = pr.MaxIter
			}
			if pr.Tol > 0 {
				tol = pr.Tol
			}
			reverse = pr.Reverse
			pers = pr.Personalization
			dang = pr.Dangling
		}
	}
	if workers < 1 {
		workers = 1
	}

	// ----- index nodes -----
	ids := g.Nodes()
	n := len(ids)
	if n == 0 {
		return res
	}
	idxOf := make(map[node.ID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}
	bidir := g.IsBidirectional()

	// ----- build outs (scatter list) -----
	// We always scatter rank[v]/outdeg[v] to outs[v].
	outs := make([][]int, n)
	outdeg := make([]int, n)

	getOuts := func(u node.ID) []int {
		nbrs := g.Neighbors(u)
		if bidir {
			// Undirected: treat neighbors as outs directly.
			out := make([]int, 0, len(nbrs))
			for _, w := range nbrs {
				if w == u {
					continue
				}
				out = append(out, idxOf[w])
			}
			return out
		}
		// Directed:
		// reverse=false => (u -> w)
		// reverse=true  => (w -> u)
		out := make([]int, 0, len(nbrs))
		if !reverse {
			for _, w := range nbrs {
				if w == u {
					continue
				}
				if g.HasEdge(u, w) {
					out = append(out, idxOf[w])
				}
			}
		} else {
			for _, w := range nbrs {
				if w == u {
					continue
				}
				if g.HasEdge(w, u) {
					out = append(out, idxOf[w])
				}
			}
		}
		return out
	}

	for i, u := range ids {
		li := getOuts(u)
		outs[i] = li
		outdeg[i] = len(li)
	}

	// ----- personalization P (L1 normalized) -----
	P := make([]float64, n)
	if pers == nil || len(*pers) == 0 {
		for i := range P {
			P[i] = 1.0 / float64(n)
		}
	} else {
		var s float64
		for _, u := range ids {
			s += (*pers)[u]
		}
		if s == 0 {
			for i := range P {
				P[i] = 1.0 / float64(n)
			}
		} else {
			inv := 1.0 / s
			for i, u := range ids {
				P[i] = (*pers)[u] * inv
			}
		}
	}

	// ----- dangling distribution D (default = P), L1 normalized -----
	D := make([]float64, n)
	if dang == nil || len(*dang) == 0 {
		copy(D, P)
	} else {
		var s float64
		for _, u := range ids {
			s += (*dang)[u]
		}
		if s == 0 {
			for i := range D {
				D[i] = 1.0 / float64(n)
			}
		} else {
			inv := 1.0 / s
			for i, u := range ids {
				D[i] = (*dang)[u] * inv
			}
		}
	}

	// ----- init rank (uniform) and teleport term -----
	rank := make([]float64, n)
	tele := make([]float64, n) // (1 - alpha)*P
	for i := 0; i < n; i++ {
		rank[i] = 1.0 / float64(n)
		tele[i] = (1.0 - alpha) * P[i]
	}

	// helper: split range
	split := func(total, k int) [][2]int {
		if k > total {
			k = total
		}
		size := (total + k - 1) / k
		chunks := make([][2]int, 0, k)
		for start := 0; start < total; start += size {
			end := start + size
			if end > total {
				end = total
			}
			chunks = append(chunks, [2]int{start, end})
		}
		return chunks
	}
	ranges := split(n, workers)

	// ----- power iteration (no per-iteration normalization; L1 test like NetworkX) -----
	for iter := 0; iter < maxIter; iter++ {
		// 1) sum of ranks on dangling nodes
		var dsum float64
		var mu sync.Mutex
		var wg sync.WaitGroup

		wg.Add(len(ranges))
		for _, rg := range ranges {
			rg := rg
			go func() {
				defer wg.Done()
				var local float64
				for i := rg[0]; i < rg[1]; i++ {
					if outdeg[i] == 0 {
						local += rank[i]
					}
				}
				if local != 0 {
					mu.Lock()
					dsum += local
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
		dangTerm := alpha * dsum

		// 2) parallel scatter into private shards
		shards := make([][]float64, len(ranges))
		for i := range shards {
			shards[i] = make([]float64, n)
		}

		wg.Add(len(ranges))
		for wi, rg := range ranges {
			wi, rg := wi, rg
			go func() {
				defer wg.Done()
				buf := shards[wi]
				for v := rg[0]; v < rg[1]; v++ {
					od := outdeg[v]
					if od == 0 {
						continue
					}
					share := alpha * rank[v] / float64(od)
					for _, u := range outs[v] {
						buf[u] += share
					}
				}
			}()
		}
		wg.Wait()

		// 3) reduce shards + add teleport and dangling redistribution
		newRank := make([]float64, n)
		for u := 0; u < n; u++ {
			val := tele[u] + dangTerm*D[u]
			for wi := 0; wi < len(shards); wi++ {
				val += shards[wi][u]
			}
			newRank[u] = val
		}

		// 4) L1 error BEFORE any normalization (NetworkX stopping rule)
		var err float64
		for i := 0; i < n; i++ {
			err += math.Abs(newRank[i] - rank[i])
		}
		rank = newRank
		if err < tol*float64(n) {
			break
		}
	}

	// 5) single final normalization so sum(rank)=1 exactly
	var sum float64
	for i := 0; i < n; i++ {
		sum += rank[i]
	}
	if sum != 0 {
		inv := 1.0 / sum
		for i := 0; i < n; i++ {
			rank[i] *= inv
		}
	}

	for i, u := range ids {
		res[u] = rank[i]
	}
	return res
}
