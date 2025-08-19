package network_graph_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"reflect"
	"testing"

	algo "github.com/elecbug/go-dspkg/network-graph/g-algorithm"
	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
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
		want    []path.Path
		wantErr bool
	}{
		{start: n1, end: n3, want: []path.Path{*path.New(n1, n2, n3)}, wantErr: false},
		{start: n1, end: n1, want: []path.Path{*path.New(n1)}, wantErr: false},
		{start: n2, end: n1, want: []path.Path{*path.New(n2, n1)}, wantErr: false},
		{start: n3, end: n4, want: []path.Path{*path.New(n3, n4)}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("From %s to %s", tt.start, tt.end), func(t *testing.T) {
			got := algo.ShortestPaths(g, tt.start, tt.end)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShortestPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathLengths(t *testing.T) {
	g := graph.New(false)

	nodeCount := 5000
	edgeCount := 10000

	for i := 0; i < nodeCount; i++ {
		g.AddNode(node.ID(fmt.Sprintf("%d", i)))
	}

	for i := 0; i < edgeCount; i++ {
		g.AddEdge(node.ID(fmt.Sprintf("%d", i)), node.ID(fmt.Sprintf("%d", rand.Intn(nodeCount))))
	}

	t.Run("CheckEqualShortestPaths", func(t *testing.T) {
		var got path.GraphPaths
		var want path.PathLength

		t.Run("WithPaths", func(t *testing.T) {
			got = algo.AllShortestPaths(g, &config.Config{Workers: 4})
		})
		gotLengths := got.OnlyLength()

		t.Run("WithoutPaths", func(t *testing.T) {
			want = algo.AllShortestPathLength(g, &config.Config{Workers: 4})
		})

		if !reflect.DeepEqual(gotLengths, want) {
			t.Errorf("AllShortestPathLength() = %v, want %v", gotLengths, want)
		}
	})

}

func TestBidirectionalGraph(t *testing.T) {
	g := graph.New(true)

	nodeCount := 1000
	edgeCount := 5000

	for i := 0; i < nodeCount; i++ {
		g.AddNode(node.ID(fmt.Sprintf("%d", i)))
	}

	for i := 0; i < edgeCount; i++ {
		g.AddEdge(node.ID(fmt.Sprintf("%d", rand.Intn(nodeCount))), node.ID(fmt.Sprintf("%d", rand.Intn(nodeCount))))
	}

	str, err := graph.Save(g)

	if err != nil {
		t.Fatalf("Failed to save graph: %v", err)
	}

	os.WriteFile("bidirectional.graph.log", []byte(str), fs.ModePerm)

	g2, err := graph.Load(str)

	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	if !reflect.DeepEqual(g, g2) {
		t.Errorf("Loaded graph is not equal to original graph")
	}

	graphMetrics(t, g, "bidirectional.")
}

func TestDirectionalGraph(t *testing.T) {
	g := graph.New(false)

	nodeCount := 1000
	edgeCount := 10000

	for i := 0; i < nodeCount; i++ {
		g.AddNode(node.ID(fmt.Sprintf("%d", i)))
	}

	for i := 0; i < edgeCount; i++ {
		g.AddEdge(node.ID(fmt.Sprintf("%d", rand.Intn(nodeCount))), node.ID(fmt.Sprintf("%d", rand.Intn(nodeCount))))
	}

	str, err := graph.Save(g)

	if err != nil {
		t.Fatalf("Failed to save graph: %v", err)
	}

	os.WriteFile("directional.graph.log", []byte(str), fs.ModePerm)

	g2, err := graph.Load(str)

	if err != nil {
		t.Fatalf("Failed to load graph: %v", err)
	}

	if !reflect.DeepEqual(g, g2) {
		t.Errorf("Loaded graph is not equal to original graph")
	}

	graphMetrics(t, g, "directional.")
}

func graphMetrics(t *testing.T, g *graph.Graph, text string) {
	results := make(map[string]interface{})
	cfg := config.Default()

	t.Run("ShortestPaths", func(t *testing.T) {
		// results["shortest_path_length"] = algo.AllShortestPathLength(g, cfg)
		results["shortest_paths"] = algo.AllShortestPaths(g, cfg).OnlyLength()
	})
	t.Run("BetweennessCentrality", func(t *testing.T) {
		results["betweenness_centrality"] = algo.BetweennessCentrality(g, cfg)
	})
	t.Run("ClosenessCentrality", func(t *testing.T) {
		results["closeness_centrality"] = algo.ClosenessCentrality(g, cfg)
	})
	t.Run("ClusteringCoefficient", func(t *testing.T) {
		results["clustering_coefficient"] = algo.ClusteringCoefficient(g, cfg)
	})
	t.Run("DegreeCentrality", func(t *testing.T) {
		results["degree_centrality"] = algo.DegreeCentrality(g, cfg)
	})
	t.Run("Diameter", func(t *testing.T) {
		results["diameter"] = algo.Diameter(g, cfg)
	})
	t.Run("EdgeBetweennessCentrality", func(t *testing.T) {
		results["edge_betweenness_centrality"] = algo.EdgeBetweennessCentrality(g, cfg)
	})
	t.Run("EigenvectorCentrality", func(t *testing.T) {
		results["eigenvector_centrality"] = algo.EigenvectorCentrality(g, cfg)
	})
	t.Run("PageRank", func(t *testing.T) {
		results["page_rank"] = algo.PageRank(g, cfg)
	})

	jsonResults, err := json.MarshalIndent(results, "", "  ")

	if err != nil {
		t.Fatalf("Failed to marshal results: %v", err)
	}

	os.WriteFile(text+"metrics.log", jsonResults, fs.ModePerm)
}
