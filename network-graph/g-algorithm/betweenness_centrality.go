package g_algorithm

import (
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// BetweennessCentrality computes the betweenness centrality for a node in the graph.
func BetweennessCentrality(g *graph.Graph, cfg *config.Config) map[node.ID]float64 {
	nodes := g.Nodes()
	n := len(nodes)

	if n < 3 {
		result := make(map[node.ID]float64)

		for _, v := range nodes {
			result[v] = 0.0
		}

		return result
	}

	workers := 1

	if cfg != nil && cfg.Workers > 1 {
		workers = cfg.Workers

		if workers > n {
			workers = n
		}
	}

	bcGlobal := make(map[node.ID]float64, n)

	var mu sync.Mutex

	tasks := make(chan node.ID, n)

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()

			localBC := make(map[node.ID]float64, n)

			for s := range tasks {
				stack := make([]node.ID, 0, n)
				preds := make(map[node.ID][]node.ID, n)
				sigma := make(map[node.ID]float64, n)
				dist := make(map[node.ID]int, n)

				for _, v := range nodes {
					dist[v] = -1
				}

				sigma[s] = 1
				dist[s] = 0

				// BFS
				q := []node.ID{s}
				for len(q) > 0 {
					v := q[0]
					q = q[1:]
					stack = append(stack, v)

					for _, w := range g.Neighbors(v) {
						if dist[w] < 0 {
							dist[w] = dist[v] + 1
							q = append(q, w)
						}

						if dist[w] == dist[v]+1 {
							sigma[w] += sigma[v]
							preds[w] = append(preds[w], v)
						}
					}
				}

				delta := make(map[node.ID]float64, n)
				for len(stack) > 0 {
					w := stack[len(stack)-1]
					stack = stack[:len(stack)-1]

					for _, v := range preds[w] {
						if sigma[w] > 0 {
							delta[v] += (sigma[v] / sigma[w]) * (1.0 + delta[w])
						}
					}

					if w != s {
						localBC[w] += delta[w]
					}
				}
			}

			if len(localBC) > 0 {
				mu.Lock()

				for v, val := range localBC {
					bcGlobal[v] += val
				}

				mu.Unlock()
			}
		}()
	}

	for _, s := range nodes {
		tasks <- s
	}

	close(tasks)
	wg.Wait()

	for v := range bcGlobal {
		bcGlobal[v] *= 0.5
	}

	norm := 2.0 / float64((n-1)*(n-2))

	for k := range bcGlobal {
		bcGlobal[k] *= norm
	}

	return bcGlobal
}
