package analyzer_test

import (
	"testing"

	"github.com/elecbug/netkit/v2/graph"
	"github.com/elecbug/netkit/v2/graph/analyzer"
)

// TestShortestPaths tests the ShortestPaths method of the Analyzer to ensure it correctly finds the shortest path between two nodes in a graph.
func TestShortestPaths(t *testing.T) {
	g := graph.New(true, true)
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	g.AddNode("D")
	g.AddEdge("A", "B", graph.NewWeight(1))
	g.AddEdge("B", "C", graph.NewWeight(1))
	g.AddEdge("A", "C", graph.NewWeight(2))
	g.AddEdge("C", "D", graph.NewWeight(1))

	a := analyzer.NewAnalyzer(g)

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

// equalPathSlices compares two graph.Path instances for equality, considering both the sequence of nodes and the total distance. It returns true if both paths are infinite or if they have the same nodes in the same order and the same total distance.
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
