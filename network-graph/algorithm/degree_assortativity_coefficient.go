// package algorithm

package algorithm

import (
	"math"

	"github.com/elecbug/netkit/network-graph/algorithm/config"
	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
)

// DegreeAssortativityCoefficient computes Newman's degree assortativity coefficient (Pearson correlation)
// between degrees at the two ends of each edge/arc, with behavior controlled by cfg.Assortativity.
//
// Assumptions / Notes:
// - For undirected graphs, it counts each edge once using upper-triangle filtering by node index.
// - For directed graphs:
//   - Mode "projected": ignores direction and uses undirected degrees on both ends.
//   - Mode "out-in":    j = out(u), k = in(v) over arcs u->v
//   - Mode "out-out":   j = out(u), k = out(v)
//   - Mode "in-in":     j = in(u),  k = in(v)
//   - Mode "in-out":    j = in(u),  k = out(v)
//
// - Self-loops are ignored by default (configurable).
// - Replace neighbor getters with your Graph's actual API if different.
func DegreeAssortativityCoefficient(g *graph.Graph, cfg *config.Config) float64 {
	if g == nil {
		return 0.0
	}
	ids := g.Nodes()
	n := len(ids)
	if n == 0 {
		return 0.0
	}

	// Resolve config with sane defaults.
	assCfg := &config.AssortativityCoefficientConfig{
		Mode:            config.AssortativityProjected,
		IgnoreSelfLoops: true,
	}
	if cfg != nil && cfg.Assortativity != nil {
		if cfg.Assortativity.Mode != "" {
			assCfg.Mode = cfg.Assortativity.Mode
		}
		if cfg.Assortativity.IgnoreSelfLoops {
			assCfg.IgnoreSelfLoops = true
		}
	}
	isUndirected := g.IsBidirectional()

	// Default directed mode: "out-in" is common in literature.
	if !isUndirected && assCfg.Mode == config.AssortativityProjected {
		assCfg.Mode = config.AssortativityOutIn
	}

	// Build an index for upper-triangle filtering on undirected graphs.
	idxOf := make(map[node.ID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	// Degree caches
	// NOTE: Replace the neighbor getters with your graph API if needed.
	outDeg := make(map[node.ID]int, n)
	inDeg := make(map[node.ID]int, n)
	undeg := make(map[node.ID]int, n)

	if isUndirected {
		for _, u := range ids {
			nu := g.Neighbors(u) // TODO: replace if your API differs
			undeg[u] = len(nu)
		}
	} else {
		// If your Graph exposes OutNeighbors/InNeighbors, use them here.
		for _, u := range ids {
			// --- BEGIN: adapt these three lines to your graph API as needed ---
			nu := g.Neighbors(u) // out-neighbors if available; otherwise total neighbors
			outDeg[u] = len(nu)

			// If you have InNeighbors(u), use it; otherwise compute by a pass:
			// Here we do a simple O(m) pass below to fill inDeg properly.
			// --- END ---
		}
		// Build in-degree by scanning arcs once.
		for _, u := range ids {
			for _, v := range g.Neighbors(u) {
				inDeg[v]++
			}
		}
		// If you also need undirected degree for "projected" mode on directed graphs:
		for _, u := range ids {
			undeg[u] = outDeg[u] + inDeg[u] // counts reciprocal twice; OK for projection usage
		}
	}

	// Accumulators per Newman:
	// r = [E(jk) - E(0.5(j+k))^2] / [E(0.5(j^2 + k^2)) - E(0.5(j+k))^2]
	var m float64
	var sumJK, sumAvg, sumSq float64

	// Iterate edges/arcs
	for _, u := range ids {
		// Iterate neighbors (or out-neighbors on directed graphs)
		for _, v := range g.Neighbors(u) {
			if assCfg.IgnoreSelfLoops && u == v {
				continue
			}
			// Undirected: count each edge once
			if isUndirected {
				if idxOf[u] > idxOf[v] {
					continue
				}
			}

			var ju, jv float64
			if isUndirected {
				ju = float64(undeg[u])
				jv = float64(undeg[v])
			} else {
				switch assCfg.Mode {
				case config.AssortativityOutIn:
					ju = float64(outDeg[u])
					jv = float64(inDeg[v])
				case config.AssortativityOutOut:
					ju = float64(outDeg[u])
					jv = float64(outDeg[v])
				case config.AssortativityInIn:
					ju = float64(inDeg[u])
					jv = float64(inDeg[v])
				case config.AssortativityInOut:
					ju = float64(inDeg[u])
					jv = float64(outDeg[v])
				case config.AssortativityProjected:
					// Directed but projected (if user explicitly requests)
					ju = float64(undeg[u])
					jv = float64(undeg[v])
				default:
					// Fallback to out-in
					ju = float64(outDeg[u])
					jv = float64(inDeg[v])
				}
			}

			sumJK += ju * jv
			sumAvg += 0.5 * (ju + jv)
			sumSq += 0.5 * (ju*ju + jv*jv)
			m += 1.0
		}
	}

	if m == 0.0 {
		return 0.0
	}
	term1 := sumJK / m
	term2 := sumAvg / m
	term3 := sumSq / m

	num := term1 - term2*term2
	den := term3 - term2*term2
	if den == 0.0 || math.IsNaN(den) {
		return 0.0
	}
	return num / den
}
