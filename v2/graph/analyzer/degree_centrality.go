package analyzer

import (
	"runtime"
	"sync"

	"github.com/elecbug/netkit/v2/graph"
)

// DegreeCentrality computes degree-based node centrality.
//
// Semantics follow NetworkX-style degree centrality:
//   - Undirected graph:
//     C_D(u) = deg(u) / (n - 1)
//   - Directed graph:
//     mode="total": (in(u) + out(u)) / (n - 1)
//     mode="in":    in(u) / (n - 1)
//     mode="out":   out(u) / (n - 1)
//
// Self-loops are ignored. For directed graphs, total degree centrality may exceed 1.0.
// The result contains one centrality value for every node in the base graph.
func (a *Analyzer) DegreeCentrality() (map[graph.NodeID]float64, error) {
	res := make(map[graph.NodeID]float64)

	if a == nil || a.baseGraph == nil {
		return res, nil
	}

	g := a.baseGraph
	ids := g.Nodes()
	n := len(ids)

	if n <= 1 {
		for _, u := range ids {
			res[u] = 0
		}
		return res, nil
	}

	denom := float64(n - 1)

	mode := DegreeCentralityTotal
	if a.cfg != nil && a.cfg.Degree != nil && a.cfg.Degree.Mode != "" {
		mode = a.cfg.Degree.Mode
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
			if v == u {
				continue
			}

			vIdx, ok := idxOf[v]
			if !ok {
				continue
			}

			row = append(row, vIdx)
		}

		outs[i] = row
	}

	var ins [][]int
	if !isUndirected {
		ins = make([][]int, n)

		for u := 0; u < n; u++ {
			for _, v := range outs[u] {
				ins[v] = append(ins[v], u)
			}
		}
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
	nums := make([]float64, n)

	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for _, rg := range chunks {
		rg := rg

		go func() {
			defer wg.Done()

			for i := rg[0]; i < rg[1]; i++ {
				if isUndirected {
					nums[i] = float64(len(outs[i]))
					continue
				}

				switch mode {
				case DegreeCentralityIn:
					nums[i] = float64(len(ins[i]))

				case DegreeCentralityOut:
					nums[i] = float64(len(outs[i]))

				default:
					nums[i] = float64(len(ins[i]) + len(outs[i]))
				}
			}
		}()
	}

	wg.Wait()

	for i, u := range ids {
		res[u] = nums[i] / denom
	}

	return res, nil
}
