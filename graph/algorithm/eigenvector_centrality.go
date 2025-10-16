package algorithm

import (
	"math"
	"runtime"
	"sync"

	"github.com/elecbug/netkit/graph"
)

// EigenvectorCentrality computes eigenvector centrality using parallel power iteration.
// Semantics match NetworkX by default:
//   - Undirected graphs: neighbors are used (symmetric).
//   - Directed graphs: by default (Reverse=false) we use predecessors/in-edges (left eigenvector).
//     Set Reverse=true to use successors/out-edges (right eigenvector).
//
// Unweighted edges are assumed. The result vector is L2-normalized (sum of squares == 1).
func EigenvectorCentrality(g *graph.Graph, cfg *Config) map[graph.NodeID]float64 {
	out := make(map[graph.NodeID]float64)
	if g == nil {
		return out
	}

	// --------- Read config & defaults ----------
	maxIter := 100
	tol := 1e-6
	reverse := false
	var nstart *map[graph.NodeID]float64
	workers := runtime.NumCPU()

	if cfg != nil {
		if cfg.Workers > 0 {
			workers = cfg.Workers
		}
		if ec := cfg.Eigenvector; ec != nil {
			if ec.MaxIter > 0 {
				maxIter = ec.MaxIter
			}
			if ec.Tol > 0 {
				tol = ec.Tol
			}
			reverse = ec.Reverse
			nstart = ec.NStart
		}
	}
	if workers < 1 {
		workers = 1
	}

	// --------- Indexing ----------
	ids := g.Nodes()
	n := len(ids)
	if n == 0 {
		return out
	}
	idxOf := make(map[graph.NodeID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	isUndirected := g.IsBidirectional()

	// --------- Build adjacency lists for multiplication y = A*x (or A^T*x) ----------
	// We will compute y[u] = sum_{v in N(u)} x[v],
	// where N(u) are:
	//   - Undirected: neighbors(u)
	//   - Directed:
	//       Reverse=false (NetworkX default): predecessors of u (in-edges v->u)
	//       Reverse=true: successors of u (out-edges u->v)
	//
	// Since Graph only exposes Neighbors (out-neighbors), we precompute both:
	//   outs[uIdx] = { indices of v | u->v }
	//   ins[uIdx]  = { indices of v | v->u }  (built by reversing outs)
	outs := make([][]int, n)
	for i, u := range ids {
		nbrs := g.Neighbors(u) // out-neighbors
		if len(nbrs) == 0 {
			continue
		}
		row := make([]int, 0, len(nbrs))
		for _, v := range nbrs {
			row = append(row, idxOf[v])
		}
		outs[i] = row
	}
	ins := make([][]int, n)
	if isUndirected {
		// For undirected storage, Neighbors(u) is symmetric, but we treat as generic neighbors.
		// ins won't be used; keep nil to avoid extra memory.
	} else {
		for u := 0; u < n; u++ {
			for _, v := range outs[u] {
				// edge u->v implies v has predecessor u
				ins[v] = append(ins[v], u)
			}
		}
	}

	// Choose which neighborhood to use per node according to direction policy.
	adj := make([][]int, n)
	if isUndirected {
		// neighbors (same as outs for undirected)
		for i := 0; i < n; i++ {
			adj[i] = outs[i]
		}
	} else {
		if !reverse {
			// NetworkX default on digraphs: use predecessors (in-edges)
			for i := 0; i < n; i++ {
				adj[i] = ins[i]
			}
		} else {
			// Out-edges (successors)
			for i := 0; i < n; i++ {
				adj[i] = outs[i]
			}
		}
	}

	// --------- Start vector x ----------
	x := make([]float64, n)
	if nstart == nil || len(*nstart) == 0 {
		// uniform start
		val := 1.0 / float64(n)
		for i := 0; i < n; i++ {
			x[i] = val
		}
	} else {
		// user-provided start (no need to normalize here; we'll normalize after the first multiply)
		for i, u := range ids {
			x[i] = (*nstart)[u]
		}
	}

	// --------- Parallel helpers ----------
	// Split [0..n) into ~workers contiguous chunks.
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

	multiply := func(dst, src []float64) {
		var wg sync.WaitGroup
		wg.Add(len(chunks))
		for _, rg := range chunks {
			rg := rg
			go func() {
				defer wg.Done()
				for u := rg[0]; u < rg[1]; u++ {
					sum := 0.0
					for _, v := range adj[u] {
						sum += src[v]
					}
					dst[u] = sum
				}
			}()
		}
		wg.Wait()
	}

	l2norm := func(vec []float64) float64 {
		// Parallel reduction of sum of squares
		var wg sync.WaitGroup
		var mu sync.Mutex
		total := 0.0
		wg.Add(len(chunks))
		for _, rg := range chunks {
			rg := rg
			go func() {
				defer wg.Done()
				local := 0.0
				for i := rg[0]; i < rg[1]; i++ {
					local += vec[i] * vec[i]
				}
				if local != 0 {
					mu.Lock()
					total += local
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
		return math.Sqrt(total)
	}

	// --------- Power iteration ----------
	y := make([]float64, n)
	for iter := 0; iter < maxIter; iter++ {
		// y = A* x  (or A^T * x) depending on adj we chose
		for i := 0; i < n; i++ {
			y[i] = 0
		}
		multiply(y, x)

		// Normalize by L2 norm (NetworkX does this each iteration)
		nrm := l2norm(y)
		if nrm == 0 {
			// Degenerate case: all-zeros (e.g., no in-edges when using predecessors).
			// Fall back to previous x and stop.
			break
		}
		inv := 1.0 / nrm
		for i := 0; i < n; i++ {
			y[i] *= inv
		}

		// Convergence check: L1 difference (sum of abs diffs) < tol * n  (NetworkX criterion)
		diff := 0.0
		for i := 0; i < n; i++ {
			diff += math.Abs(y[i] - x[i])
		}
		// swap: x <- y
		copy(x, y)

		if diff < tol*float64(n) {
			break
		}
	}

	// --------- Build result map (x is L2-normalized) ----------
	for i, u := range ids {
		out[u] = x[i]
	}
	return out
}
