package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// ClosenessCentrality computes NetworkX-compatible closeness centrality.
//
// Semantics:
//   - Directed + Reverse=false  => OUT-closeness on G (matches nx.closeness_centrality(G))
//   - Directed + Reverse=true   => IN-closeness on G  (matches nx.closeness_centrality(G.reverse()))
//   - Undirected                => standard closeness.
//
// Distances are unweighted (#edges) and come from cached all-pairs shortest paths.
//
// Requirements:
// - AllShortestPaths(g, cfg) must respect directedness of g.
// - cfg.Closeness.WfImproved follows NetworkX default (true) unless overridden.
func ClosenessCentrality(g *graph.Graph, cfg *config.Config) map[node.ID]float64 {
	out := make(map[node.ID]float64)
	if g == nil {
		return out
	}

	// ---- read config (NetworkX defaults) ----
	wfImproved := true
	reverse := false
	if cfg != nil && cfg.Closeness != nil {
		// Field name follows earlier examples: WFImproved (not WfImproved)
		wfImproved = cfg.Closeness.WfImproved
		reverse = cfg.Closeness.Reverse
	}

	ids := g.Nodes()
	N := len(ids)
	if N <= 1 {
		for _, u := range ids {
			out[u] = 0
		}
		return out
	}

	// Use cached all-pairs shortest paths.
	// Type: map[start]map[end][]path.Path
	all := AllShortestPaths(g, cfg)

	// Exact NetworkX scaling helper.
	applyNX := func(nReach int, sumDist float64, isUndirected bool) float64 {
		// nReach = r = # of reachable nodes including the node itself.
		if sumDist <= 0.0 || N <= 1 || nReach <= 1 {
			return 0.0
		}
		base := float64(nReach-1) / sumDist
		if wfImproved {
			// cc = ((r-1)/sumDist) * ((r-1)/(N-1))
			return base * (float64(nReach-1) / float64(N-1))
		}
		// Legacy: scale by reachability only for undirected graphs.
		if isUndirected {
			return base * (float64(nReach-1) / float64(N-1))
		}
		return base
	}

	isUndirected := g.IsBidirectional()

	for _, u := range ids {
		var sumDist float64
		// r = #reachable including u. Start at 1 to count u (distance 0).
		r := 1

		if isUndirected || reverse {
			// OUT-closeness on G (undirected also uses this branch).
			row := all[u] // map[end][]path.Path
			if row != nil {
				for _, v := range ids {
					if v == u {
						continue
					}
					if ps, ok := row[v]; ok && len(ps) > 0 {
						// All shortest paths u->v have equal length; take any.
						d := len(ps[0].Nodes()) - 1
						if d > 0 {
							r++
							sumDist += float64(d)
						}
					}
				}
			}
		} else {
			// IN-closeness on G — i.e., OUT-closeness on G.reverse().
			for _, v := range ids {
				if v == u {
					continue
				}
				if row, ok := all[v]; ok {
					if ps, ok2 := row[u]; ok2 && len(ps) > 0 {
						d := len(ps[0].Nodes()) - 1
						if d > 0 {
							r++
							sumDist += float64(d)
						}
					}
				}
			}
		}

		out[u] = applyNX(r, sumDist, isUndirected)
	}

	return out
}
