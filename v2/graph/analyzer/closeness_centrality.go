package analyzer

import (
	"github.com/elecbug/netkit/v2/graph"
)

// ClosenessCentrality computes closeness centrality for every node.
//
// The distance from u to reachable nodes is computed from cached shortest paths.
// If graph paths carry weights, TotalDistance() is used as the path distance.
//
// Semantics:
//   - Undirected graph:
//     standard closeness centrality
//   - Directed graph + Reverse=false:
//     in-closeness, using distances from other nodes to u
//   - Directed graph + Reverse=true:
//     out-closeness, using distances from u to other nodes
//
// If WfImproved is enabled, the result applies the NetworkX Wasserman-Faust
// correction:
//
//	C(u) = ((r - 1) / sumDist) * ((r - 1) / (n - 1))
//
// where r is the number of reachable nodes including u.
func (a *Analyzer) ClosenessCentrality() (map[graph.NodeID]float64, error) {
	out := make(map[graph.NodeID]float64)

	if a == nil || a.baseGraph == nil {
		return out, nil
	}

	if err := a.computeAllShortestPaths(); err != nil {
		return nil, err
	}

	g := a.baseGraph

	wfImproved := true
	reverse := false
	if a.cfg != nil && a.cfg.Closeness != nil {
		wfImproved = a.cfg.Closeness.WfImproved
		reverse = a.cfg.Closeness.Reverse
	}

	ids := g.Nodes()
	N := len(ids)

	if N <= 1 {
		for _, u := range ids {
			out[u] = 0
		}
		return out, nil
	}

	// Type: map[start]map[end]int or compatible numeric distance type
	all := a.allShortestPaths

	applyNX := func(nReach int, sumDist float64, isUndirected bool) float64 {
		if sumDist <= 0.0 || N <= 1 || nReach <= 1 {
			return 0.0
		}

		base := float64(nReach-1) / sumDist

		if wfImproved {
			return base * (float64(nReach-1) / float64(N-1))
		}

		if isUndirected {
			return base * (float64(nReach-1) / float64(N-1))
		}

		return base
	}

	isUndirected := !g.IsDirected()

	for _, u := range ids {
		sumDist := 0.0
		r := 1 // reachable nodes including u

		if isUndirected || reverse {
			// OUT-closeness on G.
			row := all[u]
			if row != nil {
				for _, v := range ids {
					if v == u {
						continue
					}

					ps, ok := row[v]
					if !ok || len(ps) == 0 {
						continue
					}

					dist := ps[0].TotalDistance()
					if dist <= 0 {
						continue
					}

					r++
					sumDist += float64(dist)
				}
			}
		} else {
			// IN-closeness on G.
			for _, v := range ids {
				if v == u {
					continue
				}

				row := all[v]
				if row == nil {
					continue
				}

				ps, ok := row[u]
				if !ok || len(ps) == 0 {
					continue
				}

				dist := ps[0].TotalDistance()
				if dist <= 0 {
					continue
				}

				r++
				sumDist += float64(dist)
			}
		}

		out[u] = applyNX(r, sumDist, isUndirected)
	}

	return out, nil
}
