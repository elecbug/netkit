package g_algorithm_test

import (
	"fmt"
	"reflect"
	"testing"

	algo "github.com/elecbug/go-dspkg/network-graph/g-algorithm"
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

func TestShortestPath(t *testing.T) {
	// Create a sample graph
	g := graph.New()

	n1 := node.ID("1")
	n2 := node.ID("2")
	n3 := node.ID("3")
	n4 := node.ID("4")

	g.AddNode(n1)
	g.AddNode(n2)
	g.AddNode(n3)
	g.AddNode(n4)

	g.AddEdge(n1, n2, false)
	g.AddEdge(n2, n3, false)
	g.AddEdge(n4, n3, true)

	tests := []struct {
		start   node.ID
		end     node.ID
		want    path.Path
		wantErr bool
	}{
		{start: n1, end: n3, want: *path.NewPath(n1, n2, n3), wantErr: false},
		{start: n1, end: n1, want: *path.NewPath(n1), wantErr: false},
		{start: n2, end: n1, want: *path.NewPath(), wantErr: false},
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
