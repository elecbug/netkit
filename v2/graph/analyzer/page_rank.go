package analyzer

import (
	"math"
	"runtime"
	"sync"

	"github.com/elecbug/netkit/v2/graph"
)

// PageRank computes PageRank using parallel power iteration.
//
// The implementation follows NetworkX-style PageRank semantics:
//   - teleportation uses the personalization vector
//   - dangling-node mass is redistributed using the dangling vector
//   - convergence is checked by L1 error:
//     sum(|x_new - x_old|) < n * tol
//   - the final vector is L1-normalized
//
// Semantics:
//   - Directed graph + Reverse=false:
//     standard PageRank over outgoing edges
//   - Directed graph + Reverse=true:
//     PageRank over the reversed graph
//   - Undirected graph:
//     neighbors are treated as outgoing links
//
// Directed reverse mode uses a precomputed in-neighbor adjacency list rather than
// relying on Node.Neighbors() to expose reverse edges.
func (a *Analyzer) PageRank() (map[graph.NodeID]float64, error) {
	res := make(map[graph.NodeID]float64)

	if a == nil || a.baseGraph == nil {
		return res, nil
	}

	g := a.baseGraph

	workers := runtime.NumCPU()
	if a.parallelCoreCount > 0 {
		workers = a.parallelCoreCount
	}
	if workers < 1 {
		workers = 1
	}

	alpha := 0.85
	maxIter := 1000
	tol := 1e-6
	reverse := false

	var pers *map[graph.NodeID]float64
	var dang *map[graph.NodeID]float64

	if a.cfg != nil && a.cfg.PageRank != nil {
		pr := a.cfg.PageRank

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

	ids := g.Nodes()
	n := len(ids)

	if n == 0 {
		return res, nil
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

	outsForward := make([][]int, n)
	ins := make([][]int, n)

	for i, u := range ids {
		for _, v := range neighbors(u) {
			if v == u {
				continue
			}

			j, ok := idxOf[v]
			if !ok {
				continue
			}

			// u -> v
			outsForward[i] = append(outsForward[i], j)
			ins[j] = append(ins[j], i)
		}
	}

	outs := make([][]int, n)

	if isUndirected || !reverse {
		outs = outsForward
	} else {
		outs = ins
	}

	outdeg := make([]int, n)
	for i := 0; i < n; i++ {
		outdeg[i] = len(outs[i])
	}

	P := make([]float64, n)
	if pers == nil || len(*pers) == 0 {
		for i := range P {
			P[i] = 1.0 / float64(n)
		}
	} else {
		sum := 0.0
		for _, u := range ids {
			sum += (*pers)[u]
		}

		if sum == 0 {
			for i := range P {
				P[i] = 1.0 / float64(n)
			}
		} else {
			inv := 1.0 / sum
			for i, u := range ids {
				P[i] = (*pers)[u] * inv
			}
		}
	}

	D := make([]float64, n)
	if dang == nil || len(*dang) == 0 {
		copy(D, P)
	} else {
		sum := 0.0
		for _, u := range ids {
			sum += (*dang)[u]
		}

		if sum == 0 {
			for i := range D {
				D[i] = 1.0 / float64(n)
			}
		} else {
			inv := 1.0 / sum
			for i, u := range ids {
				D[i] = (*dang)[u] * inv
			}
		}
	}

	rank := make([]float64, n)
	tele := make([]float64, n)

	for i := 0; i < n; i++ {
		rank[i] = 1.0 / float64(n)
		tele[i] = (1.0 - alpha) * P[i]
	}

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

	for iter := 0; iter < maxIter; iter++ {
		dsum := 0.0

		var mu sync.Mutex
		var wg sync.WaitGroup

		wg.Add(len(ranges))
		for _, rg := range ranges {
			rg := rg

			go func() {
				defer wg.Done()

				local := 0.0
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

		newRank := make([]float64, n)

		for u := 0; u < n; u++ {
			val := tele[u] + dangTerm*D[u]

			for wi := 0; wi < len(shards); wi++ {
				val += shards[wi][u]
			}

			newRank[u] = val
		}

		errSum := 0.0
		for i := 0; i < n; i++ {
			errSum += math.Abs(newRank[i] - rank[i])
		}

		rank = newRank

		if errSum < tol*float64(n) {
			break
		}
	}

	sum := 0.0
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

	return res, nil
}
