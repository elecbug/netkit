package analyzer

import (
	"runtime"
	"sync"

	"github.com/elecbug/netkit/v2/graph"
)

// ClusteringCoefficient computes local clustering coefficients for all nodes
// and returns their average as gcc.
//
// Returns:
//   - gcc: average local clustering coefficient over all nodes
//   - cc:  local clustering coefficient per node
//
// Semantics:
//   - Undirected graph:
//     standard local clustering coefficient:
//     links among neighbors / possible links among neighbors
//   - Directed graph:
//     Fagiolo-style directed clustering coefficient, considering both
//     incoming and outgoing neighbors and reciprocal edges.
//
// Self-loops are ignored. Nodes with degree less than 2 receive coefficient 0.
func (a *Analyzer) ClusteringCoefficient() (float64, map[graph.NodeID]float64, error) {
	cc := make(map[graph.NodeID]float64)

	if a == nil || a.baseGraph == nil {
		return 0, cc, nil
	}

	g := a.baseGraph
	nodes := g.Nodes()
	n := len(nodes)

	if n == 0 {
		return 0, cc, nil
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

	outNeighbors := make(map[graph.NodeID][]graph.NodeID, n)
	for _, v := range nodes {
		vNode, err := g.Node(v)
		if err != nil {
			continue
		}

		ns := vNode.Neighbors()

		buf := make([]graph.NodeID, 0, len(ns))
		for _, w := range ns {
			if w != v {
				buf = append(buf, w)
			}
		}

		outNeighbors[v] = buf
	}

	isDirected := g.IsDirected()

	var inNeighbors map[graph.NodeID][]graph.NodeID
	if isDirected {
		inNeighbors = make(map[graph.NodeID][]graph.NodeID, n)

		for _, u := range nodes {
			for _, w := range outNeighbors[u] {
				inNeighbors[w] = append(inNeighbors[w], u)
			}
		}
	}

	type job struct {
		v graph.NodeID
	}

	jobs := make(chan job, workers*2)
	var wg sync.WaitGroup
	var mu sync.Mutex

	b := func(u, v graph.NodeID) int {
		sum := 0

		if g.HasEdge(u, v) {
			sum++
		}
		if g.HasEdge(v, u) {
			sum++
		}

		return sum
	}

	worker := func() {
		defer wg.Done()

		if isDirected {
			for j := range jobs {
				v := j.v

				kIn := len(inNeighbors[v])
				kOut := len(outNeighbors[v])
				kTot := kIn + kOut

				if kTot < 2 {
					mu.Lock()
					cc[v] = 0
					mu.Unlock()
					continue
				}

				outSet := make(map[graph.NodeID]struct{}, kOut)
				for _, w := range outNeighbors[v] {
					outSet[w] = struct{}{}
				}

				m := 0
				for _, u := range inNeighbors[v] {
					if _, ok := outSet[u]; ok {
						m++
					}
				}

				den := float64(kTot*(kTot-1) - 2*m)
				if den <= 0 {
					mu.Lock()
					cc[v] = 0
					mu.Unlock()
					continue
				}

				totSet := make(map[graph.NodeID]struct{}, kTot)
				for _, u := range outNeighbors[v] {
					totSet[u] = struct{}{}
				}
				for _, u := range inNeighbors[v] {
					totSet[u] = struct{}{}
				}

				tot := make([]graph.NodeID, 0, len(totSet))
				for u := range totSet {
					if u != v {
						tot = append(tot, u)
					}
				}

				T := 0
				for i := 0; i < len(tot); i++ {
					jj := tot[i]
					bvj := b(v, jj)
					if bvj == 0 {
						continue
					}

					for k := 0; k < len(tot); k++ {
						if i == k {
							continue
						}

						kk := tot[k]
						bkv := b(kk, v)
						if bkv == 0 {
							continue
						}

						bjk := b(jj, kk)
						if bjk == 0 {
							continue
						}

						T += bvj * bjk * bkv
					}
				}

				t := float64(T) / 2.0
				cv := t / den

				mu.Lock()
				cc[v] = cv
				mu.Unlock()
			}

			return
		}

		for j := range jobs {
			v := j.v
			neis := outNeighbors[v]
			k := len(neis)

			if k < 2 {
				mu.Lock()
				cc[v] = 0
				mu.Unlock()
				continue
			}

			links := 0
			for i := 0; i < k; i++ {
				vi := neis[i]

				for j2 := i + 1; j2 < k; j2++ {
					vj := neis[j2]

					if g.HasEdge(vi, vj) || g.HasEdge(vj, vi) {
						links++
					}
				}
			}

			cv := float64(links) / float64(k*(k-1)/2)

			mu.Lock()
			cc[v] = cv
			mu.Unlock()
		}
	}

	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go worker()
	}

	for _, v := range nodes {
		jobs <- job{v: v}
	}
	close(jobs)

	wg.Wait()

	sum := 0.0
	for _, v := range nodes {
		sum += cc[v]
	}

	gcc := sum / float64(n)

	return gcc, cc, nil
}
