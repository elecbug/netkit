package analyzer

import (
	"container/heap"
	"math"

	"github.com/elecbug/netkit/v2/graph"
)

// Modularity computes Newman-Girvan modularity Q.
//
// If a partition is provided in config, Q is computed for that partition.
// If no partition is provided, a greedy Clauset-Newman-Moore-style modularity
// maximization is used to generate a partition first.
//
// All computations are performed on an undirected projection of the graph.
// Self-loops are ignored.
//
// The function returns 0 for empty graphs, graphs without edges, or invalid
// numerical results.
func (a *Analyzer) Modularity() (float64, error) {
	if a == nil || a.baseGraph == nil {
		return 0, nil
	}

	if a.cfg == nil || a.cfg.Modularity == nil {
		return 0, nil
	}

	g := a.baseGraph

	part := a.cfg.Modularity.Partition
	if part == nil {
		part = greedyModularityCommunitiesNX(g)
	}

	if len(part) == 0 {
		return 0, nil
	}

	return modularityQNX(g, part), nil
}

// greedyModularityCommunitiesNX implements Clauset-Newman-Moore greedy modularity maximization.
func greedyModularityCommunitiesNX(g *graph.Graph) map[graph.NodeID]int {
	ids := g.Nodes()
	n := len(ids)

	if n == 0 {
		return nil
	}

	idxOf := make(map[graph.NodeID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	edges, deg, m := undirectedEdgesAndDegrees(g, ids, idxOf)
	if m == 0 {
		part := make(map[graph.NodeID]int, n)
		for i, u := range ids {
			part[u] = i
		}
		return part
	}

	inv2m := 1.0 / (2.0 * float64(m))

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

	e := make(map[int]map[int]float64, n)
	a := make(map[int]float64, n)

	for c := range active {
		a[c] = 0
	}

	for i := 0; i < n; i++ {
		a[commOf[i]] += float64(deg[i]) * inv2m
	}

	for _, pr := range edges {
		i, j := pr[0], pr[1]
		ci, cj := commOf[i], commOf[j]

		if ci == cj {
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

	pq := &candidatePQ{}
	heap.Init(pq)

	for c, nbrs := range e {
		for d, w := range nbrs {
			if c >= d {
				continue
			}

			delta := 2.0*w - 2.0*(a[c]*a[d])
			heap.Push(pq, mergeCandidate{
				c:      c,
				d:      d,
				deltaQ: delta,
				verC:   version[c],
				verD:   version[d],
			})
		}
	}

	for pq.Len() > 0 {
		top := heap.Pop(pq).(mergeCandidate)
		c, d := top.c, top.d

		if !active[c] || !active[d] || top.verC != version[c] || top.verD != version[d] {
			continue
		}

		curW := 0.0
		if e[c] != nil {
			curW = e[c][d]
		}

		curDelta := 2.0*curW - 2.0*(a[c]*a[d])
		if curDelta <= 1e-12 {
			break
		}

		members[c] = append(members[c], members[d]...)
		delete(members, d)

		for _, idx := range members[c] {
			commOf[idx] = c
		}

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

		active[d] = false
		version[c]++
		version[d]++

		for x := range e[c] {
			if !active[x] || x == c {
				continue
			}

			delta := 2.0*e[c][x] - 2.0*(a[c]*a[x])
			cc, dd := minInt(c, x), maxInt(c, x)

			heap.Push(pq, mergeCandidate{
				c:      cc,
				d:      dd,
				deltaQ: delta,
				verC:   version[cc],
				verD:   version[dd],
			})
		}
	}

	labelOf := make(map[int]int, len(members))
	next := 0

	for commID := range members {
		labelOf[commID] = next
		next++
	}

	part := make(map[graph.NodeID]int, n)
	for oldComm, label := range labelOf {
		for _, idx := range members[oldComm] {
			part[ids[idx]] = label
		}
	}

	return part
}

// modularityQNX computes the modularity Q of a given partition of the graph.
func modularityQNX(g *graph.Graph, partition map[graph.NodeID]int) float64 {
	ids := g.Nodes()
	n := len(ids)

	if n == 0 {
		return 0
	}

	idxOf := make(map[graph.NodeID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	edges, deg, m := undirectedEdgesAndDegrees(g, ids, idxOf)
	if m == 0 {
		return 0
	}

	internalEdges := make(map[int]float64)
	degreeSum := make(map[int]float64)

	for i, u := range ids {
		c := partition[u]
		degreeSum[c] += float64(deg[i])
	}

	for _, e := range edges {
		i, j := e[0], e[1]
		u := ids[i]
		v := ids[j]

		cu := partition[u]
		cv := partition[v]

		if cu == cv {
			internalEdges[cu]++
		}
	}

	mFloat := float64(m)
	twoM := 2.0 * mFloat

	q := 0.0

	for c, kSum := range degreeSum {
		lc := internalEdges[c]

		q += (lc / mFloat) - ((kSum / twoM) * (kSum / twoM))
	}

	if math.IsNaN(q) || math.IsInf(q, 0) {
		return 0
	}

	return q
}

// undirectedEdgesAndDegrees returns the list of undirected edges (as pairs of node indices), the degree
// of each node, and the total number of edges in the undirected projection of the graph.
func undirectedEdgesAndDegrees(
	g *graph.Graph,
	ids []graph.NodeID,
	idxOf map[graph.NodeID]int,
) ([][2]int, []int, int) {
	n := len(ids)

	seen := make(map[int]map[int]bool, n)
	edges := make([][2]int, 0)
	deg := make([]int, n)

	neighbors := func(u graph.NodeID) []graph.NodeID {
		uNode, err := g.Node(u)
		if err != nil {
			return nil
		}

		return uNode.Neighbors()
	}

	for i := 0; i < n; i++ {
		u := ids[i]

		for _, v := range neighbors(u) {
			j, ok := idxOf[v]
			if !ok {
				continue
			}

			if i == j {
				continue
			}

			x, y := i, j
			if x > y {
				x, y = y, x
			}

			if seen[x] == nil {
				seen[x] = make(map[int]bool)
			}

			if seen[x][y] {
				continue
			}

			seen[x][y] = true
			edges = append(edges, [2]int{x, y})

			deg[x]++
			deg[y]++
		}
	}

	m := len(edges)
	return edges, deg, m
}

// mergeCandidate represents a potential merge of two communities in the CNM algorithm, along with the change
// in modularity (deltaQ) that would result from the merge, and the version numbers of the communities to ensure validity of the candidate.
type mergeCandidate struct {
	c, d   int
	deltaQ float64
	verC   int
	verD   int
}

// candidatePQ implements a priority queue for merge candidates based on their deltaQ values.
type candidatePQ []mergeCandidate

// Len returns the number of elements in the priority queue.
func (pq candidatePQ) Len() int {
	return len(pq)
}

// Less compares two elements in the priority queue based on their deltaQ values.
func (pq candidatePQ) Less(i, j int) bool {
	return pq[i].deltaQ > pq[j].deltaQ
}

// Swap swaps two elements in the priority queue.
func (pq candidatePQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds an element to the priority queue.
func (pq *candidatePQ) Push(x interface{}) {
	*pq = append(*pq, x.(mergeCandidate))
}

// Pop removes and returns the last element of the priority queue.
func (pq *candidatePQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// minInt returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// maxInt returns the maximum of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
