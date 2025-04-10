package graph_algorithm_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/elecbug/go-dspkg/graph"
	ga "github.com/elecbug/go-dspkg/graph/graph_algorithm"
	"github.com/elecbug/go-dspkg/graph/graph_type"
)

func TestAlgorithm(t *testing.T) {
	// g := graph.NewGraph(graph.UNDIRECTED_UNWEIGHTED, 0)
	// g.Update()
	// g.IsUpdated()

	cap := 1000
	g := graph.NewGraph(graph_type.UNDIRECTED_UNWEIGHTED, cap)

	for i := 0; i < cap; i++ {
		g.AddNode()
	}

	for i := 0; i < cap*6; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
		from := graph.NodeID(r.Intn(g.NodeCount()))

		r = rand.New(rand.NewSource(time.Now().UnixNano() + int64(i*i)))
		to := graph.NodeID(r.Intn(g.NodeCount()))

		g.AddEdge(from, to)
	}

	t.Log("ready")
	// t.Logf("\n%s\n", g.String())

	{
		u := ga.NewUnit(g)
		s := time.Now()

		t.Logf("\nAverageShortestPathLength: %v\n", spew.Sdump(u.AverageShortestPathLength()))
		t.Logf("\nDiameter: %v\n", spew.Sdump(u.Diameter()))
		t.Logf("\nPercentile 0.25: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.25)))
		t.Logf("\nPercentile 0.5: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.5)))
		t.Logf("\nPercentile 0.75: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.75)))
		t.Logf("\nPercentile 0.95: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.95)))
		// t.Logf("\nShortestPath: %v\n", spew.Sdump(u.ShortestPath(0, graph.NodeID(cap-1))))
		// t.Logf("\nBetweennessCentrality: %v\n", spew.Sdump(u.BetweennessCentrality()))
		// t.Logf("\nDegreeCentrality: %v\n", spew.Sdump(u.DegreeCentrality()))
		// t.Logf("\nClusteringCoefficient: %v\n", spew.Sdump(u.ClusteringCoefficient()))
		// t.Logf("\nEigenvectorCentrality: %v\n", spew.Sdump(u.EigenvectorCentrality(1000, 1e-6)))
		// t.Logf("\nRichClubCoefficient: %v\n", spew.Sdump(u.RichClubCoefficient(2)))
		// t.Logf("\nGlobalEfficiency: %v\n", spew.Sdump(u.GlobalEfficiency()))
		// t.Logf("\nLocalEfficiency: %v\n", spew.Sdump(u.LocalEfficiency()))

		duration := time.Since(s)
		t.Logf("execution time: %s", duration)
	}
	{
		u := ga.NewParallelUnit(g, 5)
		s := time.Now()

		t.Logf("\nAverageShortestPathLength: %v\n", spew.Sdump(u.AverageShortestPathLength()))
		t.Logf("\nDiameter: %v\n", spew.Sdump(u.Diameter()))
		t.Logf("\nPercentile 0.25: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.25)))
		t.Logf("\nPercentile 0.5: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.5)))
		t.Logf("\nPercentile 0.75: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.75)))
		t.Logf("\nPercentile 0.95: %v\n", spew.Sdump(u.PercentileShortestPathLength(0.95)))
		// t.Logf("\nShortestPath: %v\n", spew.Sdump(pu.ShortestPath(0, graph.NodeID(cap-1))))
		// t.Logf("\nBetweennessCentrality: %v\n", spew.Sdump(u.BetweennessCentrality()))
		// t.Logf("\nDegreeCentrality: %v\n", spew.Sdump(u.DegreeCentrality()))
		// t.Logf("\nClusteringCoefficient: %v\n", spew.Sdump(u.ClusteringCoefficient()))
		// t.Logf("\nEigenvectorCentrality: %v\n", spew.Sdump(u.EigenvectorCentrality(1000, 1e-6)))
		// t.Logf("\nRichClubCoefficient: %v\n", spew.Sdump(u.RichClubCoefficient(2)))
		// t.Logf("\nGlobalEfficiency: %v\n", spew.Sdump(u.GlobalEfficiency()))
		// t.Logf("\nLocalEfficiency: %v\n", spew.Sdump(u.LocalEfficiency()))

		duration := time.Since(s)
		t.Logf("execution time: %s", duration)
	}
}
