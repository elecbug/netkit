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

	cap := 200
	g := graph.NewGraph(graph_type.UNDIRECTED_UNWEIGHTED, cap)

	for i := 0; i < cap; i++ {
		g.AddNode()
	}

	for i := 0; i < g.NodeCount()*g.NodeCount()/10; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
		from := graph.NodeID(r.Intn(g.NodeCount()))

		r = rand.New(rand.NewSource(time.Now().UnixNano() + int64(i*i)))
		to := graph.NodeID(r.Intn(g.NodeCount()))

		g.AddEdge(from, to)
	}

	t.Logf("\n%s\n", g.String())

	{
		u := ga.NewUnit(g)
		s := time.Now()

		t.Logf("\nShortestPath: %v\n", spew.Sdump(u.ShortestPath(0, graph.NodeID(cap-1))))
		t.Logf("\nAverageShortestPathLength: %v\n", spew.Sdump(u.AverageShortestPathLength()))
		t.Logf("\nBetweennessCentrality: %v\n", spew.Sdump(u.BetweennessCentrality()))
		t.Logf("\nClusteringCoefficient: %v\n", spew.Sdump(u.ClusteringCoefficient()))
		t.Logf("\nDegreeCentrality: %v\n", spew.Sdump(u.DegreeCentrality()))
		t.Logf("\nDiameter: %v\n", spew.Sdump(u.Diameter()))
		t.Logf("\nEigenvectorCentrality: %v\n", spew.Sdump(u.EigenvectorCentrality(1000, 1e-6)))
		t.Logf("\nGlobalEfficiency: %v\n", spew.Sdump(u.GlobalEfficiency()))
		t.Logf("\nLocalEfficiency: %v\n", spew.Sdump(u.LocalEfficiency()))
		t.Logf("\nPercentileShortestPathLength: %v\n", spew.Sdump(u.PercentileShortestPathLength(30)))
		t.Logf("\nRichClubCoefficient: %v\n", spew.Sdump(u.RichClubCoefficient(2)))

		duration := time.Since(s)
		t.Logf("execution time: %s", duration)
	}
	{
		pu := ga.NewParallelUnit(g, 20)
		s := time.Now()

		t.Logf("\nShortestPath: %v\n", spew.Sdump(pu.ShortestPath(0, graph.NodeID(cap-1))))
		t.Logf("\nAverageShortestPathLength: %v\n", spew.Sdump(pu.AverageShortestPathLength()))
		t.Logf("\nBetweennessCentrality: %v\n", spew.Sdump(pu.BetweennessCentrality()))
		t.Logf("\nClusteringCoefficient: %v\n", spew.Sdump(pu.ClusteringCoefficient()))
		t.Logf("\nDegreeCentrality: %v\n", spew.Sdump(pu.DegreeCentrality()))
		t.Logf("\nDiameter: %v\n", spew.Sdump(pu.Diameter()))
		t.Logf("\nEigenvectorCentrality: %v\n", spew.Sdump(pu.EigenvectorCentrality(1000, 1e-6)))
		t.Logf("\nGlobalEfficiency: %v\n", spew.Sdump(pu.GlobalEfficiency()))
		t.Logf("\nLocalEfficiency: %v\n", spew.Sdump(pu.LocalEfficiency()))
		t.Logf("\nPercentileShortestPathLength: %v\n", spew.Sdump(pu.PercentileShortestPathLength(30)))
		t.Logf("\nRichClubCoefficient: %v\n", spew.Sdump(pu.RichClubCoefficient(2)))

		duration := time.Since(s)
		t.Logf("execution time: %s", duration)
	}
}
