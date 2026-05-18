package analyzer

import (
	"math"
	"runtime"
	"sync"

	"github.com/elecbug/netkit/v2/graph"
)

// EigenvectorCentrality computes eigenvector centrality using power iteration.
//
// The result vector is L2-normalized.
//
// Semantics:
//   - Undirected graph:
//     centrality is propagated through undirected neighbors
//   - Directed graph + Reverse=false:
//     centrality is computed from predecessors / incoming edges
//   - Directed graph + Reverse=true:
//     centrality is computed from successors / outgoing edges
//
// Configurable parameters include max iteration count, tolerance, reverse mode,
// worker count, and optional initial vector NStart.
//
// The convergence check uses the L1 difference between consecutive vectors.
func (a *Analyzer) EigenvectorCentrality() (map[graph.NodeID]float64, error) {
	out := make(map[graph.NodeID]float64)

	if a == nil || a.baseGraph == nil {
		return out, nil
	}

	g := a.baseGraph

	maxIter := 100
	tol := 1e-6
	reverse := false
	var nstart *map[graph.NodeID]float64

	workers := runtime.NumCPU()
	if a.parallelCoreCount > 0 {
		workers = a.parallelCoreCount
	}

	if a.cfg != nil && a.cfg.Eigenvector != nil {
		ec := a.cfg.Eigenvector

		if ec.MaxIter > 0 {
			maxIter = ec.MaxIter
		}
		if ec.Tol > 0 {
			tol = ec.Tol
		}

		reverse = ec.Reverse
		nstart = ec.NStart
	}

	if workers < 1 {
		workers = 1
	}

	ids := g.Nodes()
	n := len(ids)

	if n == 0 {
		return out, nil
	}

	if workers > n {
		workers = n
	}

	idxOf := make(map[graph.NodeID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	neighbors := func(u graph.NodeID) []graph.NodeID {
		uNode, err := g.Node(u)
		if err != nil {
			return nil
		}

		return uNode.Neighbors()
	}

	isUndirected := !g.IsDirected()

	outs := make([][]int, n)
	for i, u := range ids {
		nbrs := neighbors(u)
		if len(nbrs) == 0 {
			continue
		}

		row := make([]int, 0, len(nbrs))
		for _, v := range nbrs {
			vIdx, ok := idxOf[v]
			if !ok {
				continue
			}

			row = append(row, vIdx)
		}

		outs[i] = row
	}

	ins := make([][]int, n)
	if !isUndirected {
		for u := 0; u < n; u++ {
			for _, v := range outs[u] {
				ins[v] = append(ins[v], u)
			}
		}
	}

	adj := make([][]int, n)
	if isUndirected {
		for i := 0; i < n; i++ {
			adj[i] = outs[i]
		}
	} else if !reverse {
		for i := 0; i < n; i++ {
			adj[i] = ins[i]
		}
	} else {
		for i := 0; i < n; i++ {
			adj[i] = outs[i]
		}
	}

	x := make([]float64, n)
	if nstart == nil || len(*nstart) == 0 {
		val := 1.0 / float64(n)
		for i := 0; i < n; i++ {
			x[i] = val
		}
	} else {
		for i, u := range ids {
			x[i] = (*nstart)[u]
		}
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

	y := make([]float64, n)

	for iter := 0; iter < maxIter; iter++ {
		for i := 0; i < n; i++ {
			y[i] = 0
		}

		multiply(y, x)

		nrm := l2norm(y)
		if nrm == 0 {
			break
		}

		inv := 1.0 / nrm
		for i := 0; i < n; i++ {
			y[i] *= inv
		}

		diff := 0.0
		for i := 0; i < n; i++ {
			diff += math.Abs(y[i] - x[i])
		}

		copy(x, y)

		if diff < tol*float64(n) {
			break
		}
	}

	for i, u := range ids {
		out[u] = x[i]
	}

	return out, nil
}
