package network_graph_test

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"reflect"
	"testing"

	algo "github.com/elecbug/go-dspkg/network-graph/g-algorithm"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

func TestSimple(t *testing.T) {
	// Create a sample graph
	g := graph.New(true)

	n1 := node.ID("1")
	n2 := node.ID("2")
	n3 := node.ID("3")
	n4 := node.ID("4")

	g.AddNode(n1)
	g.AddNode(n2)
	g.AddNode(n3)
	g.AddNode(n4)

	g.AddEdge(n1, n2)
	g.AddEdge(n2, n3)
	g.AddEdge(n4, n3)

	tests := []struct {
		start   node.ID
		end     node.ID
		want    path.Path
		wantErr bool
	}{
		{start: n1, end: n3, want: *path.NewPath(n1, n2, n3), wantErr: false},
		{start: n1, end: n1, want: *path.NewPath(n1), wantErr: false},
		{start: n2, end: n1, want: *path.NewPath(n2, n1), wantErr: false},
		{start: n3, end: n4, want: *path.NewPath(n3, n4), wantErr: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("From %s to %s", tt.start, tt.end), func(t *testing.T) {
			got := algo.ShortestPath(g, tt.start, tt.end)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShortestPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraph(t *testing.T) {
	g := graph.New(true)

	for i := 0; i < 100; i++ {
		g.AddNode(node.ID(fmt.Sprintf("%d", i)))
	}

	for i := 0; i < 4000; i++ {
		g.AddEdge(node.ID(fmt.Sprintf("%d", rand.Intn(100))), node.ID(fmt.Sprintf("%d", rand.Intn(100))))
	}

	str, err := graph.Save(g)

	if err != nil {
		t.Fatalf("Failed to save graph: %v", err)
	}

	os.WriteFile("log", []byte(str), fs.ModePerm)

	g2, err := graph.Load(str)

	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	if !reflect.DeepEqual(g, g2) {
		t.Errorf("Loaded graph is not equal to original graph")
	}

	t.Run("AllShortestPaths-First", func(t *testing.T) {
		sps := algo.AllShortestPaths(g, &algo.Config{Workers: 16})
		t.Logf("Path 0 to 99: %+v\n", sps["0"]["99"])
	})

	t.Run("AllShortestPaths-Cached", func(t *testing.T) {
		algo.AllShortestPaths(g, &algo.Config{Workers: 16})
	})

	t.Run("Betweenness Centrality", func(t *testing.T) {
		t.Logf("Average Betweenness Centrality: %+v\n", algo.ToGlobal(g, algo.BetweennessCentrality, &algo.Config{Workers: 16}))
	})

	t.Run("Clustering Coefficient", func(t *testing.T) {
		t.Logf("Average Clustering Coefficient: %+v\n", algo.ToGlobal(g, algo.ClusteringCoefficient, &algo.Config{Workers: 16}))
	})

	t.Run("Degree Centrality", func(t *testing.T) {
		t.Logf("Average Degree Centrality: %+v\n", algo.ToGlobal(g, algo.DegreeCentrality, &algo.Config{Workers: 16}))
	})

	t.Run("Diameter", func(t *testing.T) {
		t.Logf("Diameter: %+v\n", algo.Diameter(g, &algo.Config{Workers: 16}))
	})

	t.Run("Eigenvector Centrality", func(t *testing.T) {
		t.Logf("Average Eigenvector Centrality: %+v\n", algo.ToGlobal(g, algo.EigenvectorCentrality, &algo.Config{Workers: 16}))
	})

	t.Run("Rich Club Coefficient", func(t *testing.T) {
		t.Logf("Average Rich Club Coefficient: %+v\n", algo.ToGlobal(g, algo.RichClubCoefficient, &algo.Config{Workers: 16}))
	})
}
