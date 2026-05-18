package analyzer_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/elecbug/netkit/v2/graph"
	"github.com/elecbug/netkit/v2/graph/analyzer"
	"github.com/elecbug/netkit/v2/graph/standard"
)

// TestShortestPaths tests the functionality of the Analyzer's shortest path computations, including
// cache management and performance with different parallel core counts.
func TestShortestPaths(t *testing.T) {
	fmt.Println("Test Shortest Paths")
	testComputeShortestPath(t)
}

// testComputeShortestPath sets up a simple graph and tests the ComputeAllShortestPaths and ShortestPaths
// methods of the Analyzer to verify that it correctly computes and caches shortest paths, and that
// it updates the cache when the graph changes. It checks for correct path results and proper error handling when paths are removed.
func testComputeShortestPath(t *testing.T) {
	fmt.Println("- Test Compute Shortest Paths")
	g := graph.New(true, true)
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	g.AddNode("D")
	g.AddEdge("A", "B", graph.NewWeight(1))
	g.AddEdge("B", "C", graph.NewWeight(1))
	g.AddEdge("A", "C", graph.NewWeight(2))
	g.AddEdge("C", "D", graph.NewWeight(1))

	a := analyzer.NewAnalyzer(g, 1, analyzer.DefaultConfig())

	paths, err := a.ShortestPaths("A", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}

	path0, err := g.Path("A", "B", "C", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path1, err := g.Path("A", "C", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if (!equalPathSlices(paths[0], *path0) && !equalPathSlices(paths[1], *path1)) &&
		(!equalPathSlices(paths[0], *path1) && !equalPathSlices(paths[1], *path0)) {
		t.Errorf("expected path %v and %v, got %v and %v", *path0, *path1, paths[0], paths[1])
	}

	g.RemoveEdge("A", "C")

	paths, err = a.ShortestPaths("A", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}

	if !equalPathSlices(paths[0], *path0) {
		t.Errorf("expected path %v, got %v", *path0, paths[0])
	}

	g.RemoveEdge("B", "C")

	paths, err = a.ShortestPaths("A", "D")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if len(paths) != 0 {
		t.Fatalf("expected 0 paths, got %d", len(paths))
	}
}

// TestPerformance creates a larger random graph and tests the performance of the ShortestPaths method with different
// parallel core counts. It measures the time taken to compute shortest paths and to retrieve cached results, ensuring
// that the method works correctly and efficiently under various conditions.
func TestPerformance(t *testing.T) {
	fmt.Println("Test Performance")

	g, err := standard.ErdosRenyiGraph(
		42,
		false,
		func(from, to *graph.Node) *graph.Weight { return graph.NewWeight(rand.Float64() * 100) },
		1000,
		0.01,
	)
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	a := analyzer.NewAnalyzer(g, 1, analyzer.DefaultConfig())

	startTime := time.Now()
	paths, err := a.ShortestPaths("0", "999")
	// NOTE: paths may be empty if 0 and 999 are not connected (P(no path) ≈ 0.009%),
	// but we just want to test the performance of the method,
	// so we won't fail the test in that case.
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	} else if len(paths) > 0 {
		fmt.Printf("  - Found shortest paths from 0 to 999: %v\n", paths[0])
	}
	duration := time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths: %v\n", duration)

	a = analyzer.NewAnalyzer(g, 4, analyzer.DefaultConfig())

	startTime = time.Now()
	pathsCompared, err := a.ShortestPaths("0", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths with 4 cores: %v\n", duration)

	if paths[0].TotalDistance() != pathsCompared[0].TotalDistance() {
		t.Errorf("expected total distance %v, got %v", paths[0].TotalDistance(), pathsCompared[0].TotalDistance())
	}

	a = analyzer.NewAnalyzer(g, 16, analyzer.DefaultConfig())

	startTime = time.Now()
	pathsCompared, err = a.ShortestPaths("0", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths with 16 cores: %v\n", duration)

	if paths[0].TotalDistance() != pathsCompared[0].TotalDistance() {
		t.Errorf("expected total distance %v, got %v", paths[0].TotalDistance(), pathsCompared[0].TotalDistance())
	}

	a = analyzer.NewAnalyzer(g, 32, analyzer.DefaultConfig())

	startTime = time.Now()
	pathsCompared, err = a.ShortestPaths("0", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths with 32 cores: %v\n", duration)

	if paths[0].TotalDistance() != pathsCompared[0].TotalDistance() {
		t.Errorf("expected total distance %v, got %v", paths[0].TotalDistance(), pathsCompared[0].TotalDistance())
	}

	startTime = time.Now()
	_, err = a.ShortestPaths("0", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("  - Time taken to retrieve cached shortest paths: %v\n", duration)
}

// equalPathSlices compares two graph.Path values by node sequence and total distance.
func equalPathSlices(a, b graph.Path) bool {
	nodesA := a.Nodes()
	nodesB := b.Nodes()

	if len(nodesA) != len(nodesB) {
		return false
	}

	if a.TotalDistance() != b.TotalDistance() {
		return false
	}

	for i := range nodesA {
		if nodesA[i] != nodesB[i] {
			return false
		}
	}

	return true
}

func TestAnaylzer(t *testing.T) {
	fmt.Println("Test Analyzer")

	results := make(map[string]interface{})
	g, err := standard.ErdosRenyiGraph(time.Now().Nanosecond(), false, nil, 1000, 0.01)
	if err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}

	cfg := analyzer.DefaultConfig()
	a := analyzer.NewAnalyzer(g, 16, cfg)

	t.Run("ShortestPaths", func(t *testing.T) {
		results["shortest_paths"] = make(map[string]any)

		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if i == j {
					continue
				}

				paths, err := a.ShortestPaths(graph.NodeID(fmt.Sprintf("%d", i*100)), graph.NodeID(fmt.Sprintf("%d", j*100)))
				if err != nil {
					t.Fatalf("Failed to compute shortest paths: %v", err)
				}

				results["shortest_paths"].(map[string]any)[fmt.Sprintf("%d->%d", i*100, j*100)] = paths
			}
		}
	})
	t.Run("BetweennessCentrality", func(t *testing.T) {
		res, err := a.BetweennessCentrality()
		if err != nil {
			t.Fatalf("Failed to compute betweenness centrality: %v", err)
		}

		results["betweenness_centrality"] = res
	})
	t.Run("ClosenessCentrality", func(t *testing.T) {
		res, err := a.ClosenessCentrality()
		if err != nil {
			t.Fatalf("Failed to compute closeness centrality: %v", err)
		}

		results["closeness_centrality"] = res
	})
	t.Run("ClusteringCoefficient", func(t *testing.T) {
		gcc, ccs, err := a.ClusteringCoefficient()
		if err != nil {
			t.Fatalf("Failed to compute clustering coefficient: %v", err)
		}

		results["clustering_coefficient"] = map[string]any{
			"global": gcc,
			"local":  ccs,
			"average": func() float64 {
				sum := 0.0
				for _, v := range ccs {
					sum += v
				}
				return sum / float64(len(ccs))
			}(),
		}
	})
	t.Run("DegreeAssortativityCoefficient", func(t *testing.T) {
		res, err := a.DegreeAssortativityCoefficient()
		if err != nil {
			t.Fatalf("Failed to compute degree assortativity coefficient: %v", err)
		}

		results["degree_assortativity_coefficient"] = res
	})
	t.Run("DegreeCentrality", func(t *testing.T) {
		res, err := a.DegreeCentrality()
		if err != nil {
			t.Fatalf("Failed to compute degree centrality: %v", err)
		}

		results["degree_centrality"] = res
	})
	t.Run("Diameter", func(t *testing.T) {
		res, weight, err := a.Diameter()
		if err != nil {
			t.Fatalf("Failed to compute diameter: %v", err)
		}

		results["diameter"] = res
		results["diameter_weight"] = weight
	})
	t.Run("EdgeBetweennessCentrality", func(t *testing.T) {
		res, err := a.EdgeBetweennessCentrality()
		if err != nil {
			t.Fatalf("Failed to compute edge betweenness centrality: %v", err)
		}
		results["edge_betweenness_centrality"] = res
	})
	t.Run("EigenvectorCentrality", func(t *testing.T) {
		res, err := a.EigenvectorCentrality()
		if err != nil {
			t.Fatalf("Failed to compute eigenvector centrality: %v", err)
		}
		results["eigenvector_centrality"] = res
	})
	t.Run("Modularity", func(t *testing.T) {
		res, err := a.Modularity()
		if err != nil {
			t.Fatalf("Failed to compute modularity: %v", err)
		}
		results["modularity"] = res
	})
	t.Run("PageRank", func(t *testing.T) {
		res, err := a.PageRank()
		if err != nil {
			t.Fatalf("Failed to compute page rank: %v", err)
		}
		results["page_rank"] = res
	})

	jsonResults, err := json.MarshalIndent(results, "", "  ")

	if err != nil {
		t.Fatalf("Failed to marshal results: %v", err)
	}

	os.WriteFile("metrics.log", jsonResults, fs.ModePerm)

	gJson, err := g.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize graph: %v", err)
	}

	os.WriteFile("graph.log", []byte(gJson), fs.ModePerm)
}
