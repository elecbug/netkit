package g_algorithm

import (
	"runtime"
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// ClusteringCoefficientAll computes local clustering coefficients for all nodes.
// - If g.IsBidirectional()==false (directed): Fagiolo (2007) directed clustering (matches NetworkX).
// - If g.IsBidirectional()==true (undirected): standard undirected clustering.
// Returns map[node.ID]float64 with a value for every node in g.
func ClusteringCoefficient(g *graph.Graph, cfg *config.Config) map[node.ID]float64 {
	res := make(map[node.ID]float64)
	if g == nil {
		return res
	}

	nodes := g.Nodes()
	n := len(nodes)
	if n == 0 {
		return res
	}

	workers := 0
	if cfg != nil && cfg.Workers > 0 {
		workers = cfg.Workers
	}
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	// Build helper structures
	// outNeighbors[v] = slice of out-neighbors of v (exclude self)
	// inNeighbors[v]  = slice of in-neighbors of v (exclude self) - only needed for directed
	outNeighbors := make(map[node.ID][]node.ID, n)
	for _, v := range nodes {
		ns := g.Neighbors(v)
		buf := make([]node.ID, 0, len(ns))
		for _, w := range ns {
			if w != v {
				buf = append(buf, w)
			}
		}
		outNeighbors[v] = buf
	}

	isDirected := !g.IsBidirectional()
	var inNeighbors map[node.ID][]node.ID
	if isDirected {
		inNeighbors = make(map[node.ID][]node.ID, n)
		for _, u := range nodes {
			for _, w := range outNeighbors[u] {
				// u -> w, so u is in-neighbor of w
				inNeighbors[w] = append(inNeighbors[w], u)
			}
		}
	}

	type job struct{ v node.ID }
	jobs := make(chan job, workers*2)
	var wg sync.WaitGroup
	var mu sync.Mutex // protects res map

	// Edge multiplicity for Fagiolo: b(u,v) = a_uv + a_vu ∈ {0,1,2}
	b := func(u, v node.ID) int {
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
			// Fagiolo directed clustering (matches NetworkX's directed_clustering)
			for j := range jobs {
				v := j.v
				// k_in, k_out
				kIn := len(inNeighbors[v])
				kOut := len(outNeighbors[v])
				kTot := kIn + kOut

				// m_v: number of mutual (reciprocal) dyads with v
				// Compute intersection of in(v) and out(v)
				if kTot < 2 {
					mu.Lock()
					res[v] = 0.0
					mu.Unlock()
					continue
				}
				outSet := make(map[node.ID]struct{}, kOut)
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
					res[v] = 0.0
					mu.Unlock()
					continue
				}

				// T_v = sum_{j != k} b(v,j) * b(j,k) * b(k,v)
				// with j,k in tot = in(v) ∪ out(v)
				totSet := make(map[node.ID]struct{}, kTot) // upper bound
				for _, u := range outNeighbors[v] {
					totSet[u] = struct{}{}
				}
				for _, u := range inNeighbors[v] {
					totSet[u] = struct{}{}
				}
				// Make a slice to iterate
				tot := make([]node.ID, 0, len(totSet))
				for u := range totSet {
					if u != v { // guard (shouldn't be in set anyway)
						tot = append(tot, u)
					}
				}
				// Accumulate ordered pairs (j,k), j!=k
				var T int
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
				// t_v = T / 2
				t := float64(T) / 2.0
				cv := t / den

				mu.Lock()
				res[v] = cv
				mu.Unlock()
			}
		} else {
			// Undirected clustering (standard)
			for j := range jobs {
				v := j.v
				neis := outNeighbors[v] // in undirected storage, this is symmetric
				k := len(neis)
				if k < 2 {
					mu.Lock()
					res[v] = 0.0
					mu.Unlock()
					continue
				}
				// Count edges among neighbors (each neighbor pair that is connected)
				links := 0
				for i := 0; i < k; i++ {
					vi := neis[i]
					for j2 := i + 1; j2 < k; j2++ {
						vj := neis[j2]
						// in bidirectional storage, one direction check is enough
						if g.HasEdge(vi, vj) || g.HasEdge(vj, vi) {
							links++
						}
					}
				}
				cv := float64(links) / float64(k*(k-1)/2)
				mu.Lock()
				res[v] = cv
				mu.Unlock()
			}
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

	return res
}
