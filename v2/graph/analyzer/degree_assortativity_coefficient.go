package analyzer

import (
	"math"

	"github.com/elecbug/netkit/v2/graph"
)

// DegreeAssortativityCoefficient computes Newman's degree assortativity coefficient.
//
// The coefficient is the Pearson correlation between degree values at the two
// ends of each edge or arc.
//
// Semantics:
//   - Undirected graph:
//     each undirected edge is counted once
//   - Directed graph:
//     mode controls which degree pair is compared:
//     "out-in":  out(u) with in(v)
//     "out-out": out(u) with out(v)
//     "in-in":   in(u) with in(v)
//     "in-out":  in(u) with out(v)
//     "projected": direction is ignored and projected degrees are used
//
// Self-loops are ignored by default. If the variance term is zero, the function
// returns 0.
func (a *Analyzer) DegreeAssortativityCoefficient() (float64, error) {
	if a == nil || a.baseGraph == nil {
		return 0, nil
	}

	g := a.baseGraph
	ids := g.Nodes()
	n := len(ids)

	if n == 0 {
		return 0, nil
	}

	assCfg := &AssortativityCoefficientConfig{
		Mode:            AssortativityProjected,
		IgnoreSelfLoops: true,
	}

	if a.cfg != nil && a.cfg.Assortativity != nil {
		if a.cfg.Assortativity.Mode != "" {
			assCfg.Mode = a.cfg.Assortativity.Mode
		}

		if a.cfg.Assortativity.IgnoreSelfLoops {
			assCfg.IgnoreSelfLoops = true
		}
	}

	isUndirected := !g.IsDirected()

	if !isUndirected && assCfg.Mode == AssortativityProjected {
		assCfg.Mode = AssortativityOutIn
	}

	idxOf := make(map[graph.NodeID]int, n)
	for i, u := range ids {
		idxOf[u] = i
	}

	neighbors := func(u graph.NodeID) []graph.NodeID {
		uNode, err := g.Node(u)
		if err != nil {
			return nil
		}
		return uNode.Neighbors()
	}

	outDeg := make(map[graph.NodeID]int, n)
	inDeg := make(map[graph.NodeID]int, n)
	undeg := make(map[graph.NodeID]int, n)

	if isUndirected {
		for _, u := range ids {
			nu := neighbors(u)
			undeg[u] = len(nu)
		}
	} else {
		for _, u := range ids {
			nu := neighbors(u)
			outDeg[u] = len(nu)
		}

		for _, u := range ids {
			for _, v := range neighbors(u) {
				inDeg[v]++
			}
		}

		for _, u := range ids {
			undeg[u] = outDeg[u] + inDeg[u]
		}
	}

	var m float64
	var sumJK, sumAvg, sumSq float64

	for _, u := range ids {
		for _, v := range neighbors(u) {
			if assCfg.IgnoreSelfLoops && u == v {
				continue
			}

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
				case AssortativityOutIn:
					ju = float64(outDeg[u])
					jv = float64(inDeg[v])

				case AssortativityOutOut:
					ju = float64(outDeg[u])
					jv = float64(outDeg[v])

				case AssortativityInIn:
					ju = float64(inDeg[u])
					jv = float64(inDeg[v])

				case AssortativityInOut:
					ju = float64(inDeg[u])
					jv = float64(outDeg[v])

				case AssortativityProjected:
					ju = float64(undeg[u])
					jv = float64(undeg[v])

				default:
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

	if m == 0 {
		return 0, nil
	}

	term1 := sumJK / m
	term2 := sumAvg / m
	term3 := sumSq / m

	num := term1 - term2*term2
	den := term3 - term2*term2

	if den == 0 || math.IsNaN(den) {
		return 0, nil
	}

	return num / den, nil
}
