package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

// ShortestPath computes a shortest path between start and end using BFS.
// It returns an empty path when no path exists.
func ShortestPath(graph *graph.Graph, start, end node.ID) path.Path {
	if start == end {
		return *path.NewPath(start)
	}

	queue := []node.ID{start}
	visited := make(map[node.ID]bool)
	parent := make(map[node.ID]node.ID)
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		neighbors := graph.GetEdges(current)

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)

				if neighbor == end {
					// Reconstruct path
					p := []node.ID{}

					for n := end; n != start; n = parent[n] {
						p = append([]node.ID{n}, p...)
					}

					p = append([]node.ID{start}, p...)

					return *path.NewPath(p...)
				}
			}
		}
	}

	return *path.NewPath() // No path found
}
