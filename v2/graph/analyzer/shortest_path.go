package analyzer

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/elecbug/netkit/v2/graph"
)

// ComputeAllShortestPaths computes and caches all shortest paths for the current graph.
func (a *Analyzer) ComputeAllShortestPaths() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	currentHash := a.baseGraph.Hash()

	if a.graphHash == currentHash && len(a.allShortestPaths) > 0 {
		return nil
	}

	paths, err := allShortestPaths(a.baseGraph)
	if err != nil {
		return err
	}

	a.allShortestPaths = paths
	a.graphHash = currentHash

	return nil
}

// ShortestPaths returns cached shortest paths between start and end.
// If the cache is stale, it recomputes all shortest paths first.
func (a *Analyzer) ShortestPaths(start, end graph.NodeID) ([]graph.Path, error) {
	if err := a.ComputeAllShortestPaths(); err != nil {
		return nil, err
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	byStart, ok := a.allShortestPaths[start]
	if !ok {
		return nil, fmt.Errorf("start node %s not found", start)
	}

	paths, ok := byStart[end]
	if !ok {
		return nil, fmt.Errorf("no path found from %s to %s", start, end)
	}

	result := make([]graph.Path, len(paths))
	copy(result, paths)

	return result, nil
}

// Path computes the path and total distance for a given sequence of nodes. It returns a Path containing
// the total distance and the sequence of nodes. If any edge in the path does not exist, it returns a
// Path marked as infinite (unreachable) and an error indicating which edge is missing.
func allShortestPaths(g *graph.Graph) (map[graph.NodeID]map[graph.NodeID][]graph.Path, error) {
	result := make(map[graph.NodeID]map[graph.NodeID][]graph.Path)

	if !g.Weighted {
		for _, start := range g.Nodes() {
			pathsFromStart, err := allShortestPathsFromStart(g, start)
			if err != nil {
				return nil, err
			}

			result[start] = pathsFromStart
		}

		return result, nil
	}

	for _, start := range g.Nodes() {
		pathsFromStart, err := allWeightedShortestPathsFromStart(g, start)
		if err != nil {
			return nil, err
		}

		result[start] = pathsFromStart
	}

	return result, nil
}

// allShortestPathsFromStart computes all shortest paths from the given start node to all reachable nodes in the graph.
func allShortestPathsFromStart(g *graph.Graph, start graph.NodeID) (map[graph.NodeID][]graph.Path, error) {
	dist := make(map[graph.NodeID]int)
	preds := make(map[graph.NodeID]map[graph.NodeID]bool)

	queue := []graph.NodeID{start}
	dist[start] = 0

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		nodeV, err := g.Node(v)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %w", v, err)
		}

		for _, w := range nodeV.Neighbors() {
			nextDist := dist[v] + 1

			oldDist, seen := dist[w]
			if !seen {
				dist[w] = nextDist
				queue = append(queue, w)

				if preds[w] == nil {
					preds[w] = make(map[graph.NodeID]bool)
				}
				preds[w][v] = true

				continue
			}

			if oldDist == nextDist {
				if preds[w] == nil {
					preds[w] = make(map[graph.NodeID]bool)
				}
				preds[w][v] = true
			}
		}
	}

	result := make(map[graph.NodeID][]graph.Path)

	for _, end := range g.Nodes() {
		if _, ok := dist[end]; !ok {
			continue
		}

		paths, err := buildPathsFromPreds(g, start, end, preds)
		if err != nil {
			return nil, err
		}

		result[end] = paths
	}

	return result, nil
}

// buildPathsFromPreds constructs all shortest paths from start to end using the predecessors map.
func buildPathsFromPreds(
	g *graph.Graph,
	start, end graph.NodeID,
	preds map[graph.NodeID]map[graph.NodeID]bool,
) ([]graph.Path, error) {
	if start == end {
		path, err := g.Path(start)
		if err != nil {
			return nil, err
		}
		return []graph.Path{*path}, nil
	}

	var rawPaths [][]graph.NodeID
	cur := []graph.NodeID{end}

	var dfs func(u graph.NodeID)
	dfs = func(u graph.NodeID) {
		if u == start {
			seq := make([]graph.NodeID, len(cur))
			for i := range cur {
				seq[i] = cur[len(cur)-1-i]
			}
			rawPaths = append(rawPaths, seq)
			return
		}

		for p := range preds[u] {
			cur = append(cur, p)
			dfs(p)
			cur = cur[:len(cur)-1]
		}
	}

	dfs(end)

	paths := make([]graph.Path, 0, len(rawPaths))

	for _, seq := range rawPaths {
		path, err := g.Path(seq...)
		if err != nil {
			return nil, fmt.Errorf("failed to create path for sequence %v: %w", seq, err)
		}

		paths = append(paths, *path)
	}

	return paths, nil
}

// allWeightedShortestPathsFromStart computes all shortest paths from the given start
// node to all reachable nodes in a weighted graph using Dijkstra's algorithm.
func allWeightedShortestPathsFromStart(g *graph.Graph, start graph.NodeID) (map[graph.NodeID][]graph.Path, error) {
	const eps = 1e-9

	dist := make(map[graph.NodeID]float64)
	preds := make(map[graph.NodeID]map[graph.NodeID]bool)

	for _, id := range g.Nodes() {
		dist[id] = math.Inf(1)
	}

	dist[start] = 0

	pq := &dijkstraPQ{}
	heap.Init(pq)
	heap.Push(pq, dijkstraItem{
		id:   start,
		dist: 0,
	})

	for pq.Len() > 0 {
		item := heap.Pop(pq).(dijkstraItem)
		v := item.id

		if item.dist > dist[v]+eps {
			continue
		}

		nodeV, err := g.Node(v)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %w", v, err)
		}

		for _, w := range nodeV.Neighbors() {
			weight, err := g.EdgeWeight(v, w)
			if err != nil {
				return nil, fmt.Errorf("failed to get edge weight %s -> %s: %w", v, w, err)
			}

			if weight < 0 {
				return nil, fmt.Errorf("negative edge weight is not supported by Dijkstra: %s -> %s = %f", v, w, weight)
			}

			nextDist := dist[v] + float64(weight)

			if nextDist < dist[w]-eps {
				dist[w] = nextDist

				preds[w] = map[graph.NodeID]bool{
					v: true,
				}

				heap.Push(pq, dijkstraItem{
					id:   w,
					dist: nextDist,
				})

				continue
			}

			if math.Abs(nextDist-dist[w]) <= eps {
				if preds[w] == nil {
					preds[w] = make(map[graph.NodeID]bool)
				}

				preds[w][v] = true
			}
		}
	}

	result := make(map[graph.NodeID][]graph.Path)

	for _, end := range g.Nodes() {
		if math.IsInf(dist[end], 1) {
			continue
		}

		paths, err := buildPathsFromPreds(g, start, end, preds)
		if err != nil {
			return nil, err
		}

		result[end] = paths
	}

	return result, nil
}
