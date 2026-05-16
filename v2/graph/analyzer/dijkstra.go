package analyzer

import "github.com/elecbug/netkit/v2/graph"

// Dijkstra's algorithm implementation for weighted graphs
type dijkstraItem struct {
	id   graph.NodeID
	dist float64
}

// dijkstraPQ implements a priority queue for Dijkstra's algorithm based on the distance from the start node.
type dijkstraPQ []dijkstraItem

// Len returns the number of items in the priority queue.
func (pq dijkstraPQ) Len() int {
	return len(pq)
}

// Less compares two items in the priority queue based on their distance, returning true
// if the item at index i has a smaller distance than the item at index j.
func (pq dijkstraPQ) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

// Swap exchanges the items at indices i and j in the priority queue.
func (pq dijkstraPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds a new item to the priority queue. It appends the item to the end of the slice.
func (pq *dijkstraPQ) Push(x any) {
	*pq = append(*pq, x.(dijkstraItem))
}

// Pop removes and returns the item with the smallest distance from the priority queue.
// container/heap moves the smallest item to the end before calling this method.
func (pq *dijkstraPQ) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
