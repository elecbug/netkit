// package algorithm

package algorithm

import (
	"container/heap"
	"math"

	"github.com/elecbug/netkit/network-graph/algorithm/config"
	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
)

// Modularity computes Newman-Girvan modularity Q.
// If cfg.Modularity.Partition == nil, it runs greedy CNM to obtain a partition first.
// All computations are performed on an undirected projection (to match NetworkX).
func Modularity(g *graph.Graph, cfg *config.Config) float64 {
	if g == nil || cfg == nil || cfg.Modularity == nil {
		return 0.0
	}
	if cfg.Modularity.Partition == nil {
		cfg.Modularity.Partition = GreedyModularityCommunitiesNX(g)
	}
	part := cfg.Modularity.Partition
	if len(part) == 0 {
		return 0.0
	}
	return modularityQNX(g, part)
}

// GreedyModularityCommunitiesNX implements the Clauset–Newman–Moore greedy
// modularity maximization, aligned with NetworkX conventions:
// - work on an undirected projection
// - m = number_of_undirected_edges
// - inv2m = 1/(2m); each undirected edge contributes inv2m twice (symmetrically)
// Returns a partition map: nodeID -> compact community label.
func GreedyModularityCommunitiesNX(g *graph.Graph) map[node.ID]int {
	ids := g.Nodes()
	n := len(ids)
	if n == 0 {
		return nil
	}

	// Build stable indices
	idxOf := make(map[node.ID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	// Build undirected edge set and degrees for the projection
	edges, deg, m := undirectedEdgesAndDegrees(g, ids, idxOf)
	if m == 0 {
		part := make(map[node.ID]int, n)
		for i, u := range ids {
			part[u] = i
		}
		return part
	}
	inv2m := 1.0 / (2.0 * float64(m))

	// Initialize communities: each node is its own community
	commOf := make([]int, n)
	members := make(map[int][]int, n)
	active := make(map[int]bool, n)
	version := make(map[int]int, n)
	for i := 0; i < n; i++ {
		commOf[i] = i
		members[i] = []int{i}
		active[i] = true
		version[i] = 0
	}

	// e[c][d]: fraction of edges between communities c and d (over 2m), symmetric
	// a[c]   : sum_j e[c][j]
	e := make(map[int]map[int]float64, n)
	a := make(map[int]float64, n)
	for c := range active {
		a[c] = 0
	}

	// Initialize a from node degrees: a_c += k_i / (2m) for i in c (singleton now)
	for i := 0; i < n; i++ {
		a[commOf[i]] += float64(deg[i]) * inv2m
	}

	// Initialize e from undirected edges (each contributes inv2m to both e[c][d] and e[d][c])
	for _, pr := range edges {
		i, j := pr[0], pr[1] // i<j ensured
		ci, cj := commOf[i], commOf[j]
		if ci == cj {
			// Initially impossible; keep for completeness
			continue
		}
		w := inv2m
		if e[ci] == nil {
			e[ci] = make(map[int]float64)
		}
		if e[cj] == nil {
			e[cj] = make(map[int]float64)
		}
		e[ci][cj] += w
		e[cj][ci] += w
	}

	// Priority queue of merges by ΔQ = 2*(e_cd - a_c * a_d), largest first
	pq := &candidatePQ{}
	heap.Init(pq)

	// Seed
	for c, nbrs := range e {
		for d, w := range nbrs {
			if c >= d {
				continue
			}
			delta := 2.0*(w) - 2.0*(a[c]*a[d])
			heap.Push(pq, mergeCandidate{c: c, d: d, deltaQ: delta, verC: version[c], verD: version[d]})
		}
	}

	// Greedy loop
	for pq.Len() > 0 {
		top := heap.Pop(pq).(mergeCandidate)
		c, d := top.c, top.d
		if !active[c] || !active[d] || top.verC != version[c] || top.verD != version[d] {
			continue
		}
		// Recompute current ΔQ
		curW := 0.0
		if e[c] != nil {
			curW = e[c][d]
		}
		curDelta := 2.0*(curW) - 2.0*(a[c]*a[d])
		if curDelta <= 1e-12 {
			// No further positive gain
			break
		}

		// Merge d into c
		// 1) membership
		members[c] = append(members[c], members[d]...)
		delete(members, d)
		for _, idx := range members[c] {
			commOf[idx] = c
		}

		// 2) update e and a
		neighborsD := e[d]
		if e[c] == nil {
			e[c] = make(map[int]float64)
		}
		for x, w := range neighborsD {
			if x == c {
				continue
			}
			e[c][x] += w
			if e[x] == nil {
				e[x] = make(map[int]float64)
			}
			e[x][c] += w
			delete(e[x], d)
		}
		delete(e[c], d)
		delete(e, d)

		a[c] += a[d]
		delete(a, d)

		// 3) deactivate d, bump versions
		active[d] = false
		version[c]++
		version[d]++

		// 4) push new candidates (c,x)
		for x := range e[c] {
			if !active[x] || x == c {
				continue
			}
			delta := 2.0*(e[c][x]) - 2.0*(a[c]*a[x])
			cc, dd := minInt(c, x), maxInt(c, x)
			heap.Push(pq, mergeCandidate{c: cc, d: dd, deltaQ: delta, verC: version[cc], verD: version[dd]})
		}
	}

	// Compact community labels
	labelOf := make(map[int]int, len(members))
	next := 0
	for commID := range members {
		labelOf[commID] = next
		next++
	}

	part := make(map[node.ID]int, n)
	for oldComm, label := range labelOf {
		for _, idx := range members[oldComm] {
			part[ids[idx]] = label
		}
	}
	return part
}

// modularityQNX evaluates Q with the NetworkX-compatible undirected definition:
//
//	m = number_of_undirected_edges
//	inv2m = 1/(2m)
//	Q = (1/2m) * sum_{(i,j) in undirected edges, c(i)=c(j)} [ 1 - (k_i k_j)/(2m) ]
func modularityQNX(g *graph.Graph, partition map[node.ID]int) float64 {
	ids := g.Nodes()
	n := len(ids)
	if n == 0 {
		return 0.0
	}
	idxOf := make(map[node.ID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	edges, deg, m := undirectedEdgesAndDegrees(g, ids, idxOf)
	if m == 0 {
		return 0.0
	}
	inv2m := 1.0 / (2.0 * float64(m))

	var sum float64
	for _, pr := range edges {
		i, j := pr[0], pr[1] // i<j
		ui, uj := ids[i], ids[j]
		if partition[ui] != partition[uj] {
			continue
		}
		ki := float64(deg[i])
		kj := float64(deg[j])
		sum += 1.0 - (ki*kj)*inv2m
	}
	Q := sum * inv2m
	if math.IsNaN(Q) || math.IsInf(Q, 0) {
		return 0.0
	}
	return Q
}

// undirectedEdgesAndDegrees builds the undirected projection (unique i<j pairs),
// returns: list of edges as pairs of indices, degree per node in the projection, and m (=|E|).
// An undirected edge contributes exactly once as (i<j).
func undirectedEdgesAndDegrees(g *graph.Graph, ids []node.ID, idxOf map[node.ID]int) ([][2]int, []int, int) {
	n := len(ids)
	seen := make(map[int]map[int]bool, n)
	edges := make([][2]int, 0)
	deg := make([]int, n)

	for i := 0; i < n; i++ {
		u := ids[i]
		for _, v := range g.Neighbors(u) {
			j := idxOf[v]
			if i == j {
				continue
			}
			a, b := i, j
			if a > b {
				a, b = b, a
			}
			if seen[a] == nil {
				seen[a] = make(map[int]bool)
			}
			if seen[a][b] {
				continue
			}
			seen[a][b] = true
			edges = append(edges, [2]int{a, b})
			// increase degree in projection
			deg[a]++
			deg[b]++
		}
	}
	m := len(edges)
	return edges, deg, m
}

//
// ---------- Helpers & PQ ----------
//

type mergeCandidate struct {
	c, d   int     // community IDs (ensure c<d when pushing to PQ)
	deltaQ float64 // modularity gain if merged
	verC   int     // version stamps for staleness check
	verD   int
}

type candidatePQ []mergeCandidate

func (pq candidatePQ) Len() int            { return len(pq) }
func (pq candidatePQ) Less(i, j int) bool  { return pq[i].deltaQ > pq[j].deltaQ } // max-heap
func (pq candidatePQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *candidatePQ) Push(x interface{}) { *pq = append(*pq, x.(mergeCandidate)) }
func (pq *candidatePQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
