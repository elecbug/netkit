package analyzer_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/elecbug/netkit/v2/graph"
	"github.com/elecbug/netkit/v2/graph/analyzer"
	"github.com/elecbug/netkit/v2/graph/standard"
)

// TestShortestPaths tests the ShortestPaths method of the Analyzer to ensure it correctly finds the shortest path between two nodes in a graph.
func TestShortestPaths(t *testing.T) {
	fmt.Println("Test Shortest Paths")
	testComputeShortestPath(t)
	testPerformance(t)
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

	a := analyzer.NewAnalyzer(g, 1)

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

func testPerformance(t *testing.T) {
	fmt.Println("- Test Performance")

	g, err := standard.ErdosRenyiGraph(
		42,
		false,
		func(from, to graph.NodeID) *graph.Weight { return graph.NewWeight(rand.Float64() * 100) },
		1000,
		0.01,
	)
	if err != nil {
		t.Fatalf("failed to create graph: %v", err)
	}

	a := analyzer.NewAnalyzer(g, 1)

	startTime := time.Now()
	_, err = a.ShortestPaths("0", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration := time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths: %v\n", duration)

	a = analyzer.NewAnalyzer(g, 4)

	startTime = time.Now()
	_, err = a.ShortestPaths("0", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths with 4 cores: %v\n", duration)

	a = analyzer.NewAnalyzer(g, 16)

	startTime = time.Now()
	_, err = a.ShortestPaths("0", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths with 16 cores: %v\n", duration)

	a = analyzer.NewAnalyzer(g, 32)

	startTime = time.Now()
	_, err = a.ShortestPaths("0", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("  - Time taken to compute shortest paths with 32 cores: %v\n", duration)

	startTime = time.Now()
	_, err = a.ShortestPaths("0", "999")
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
