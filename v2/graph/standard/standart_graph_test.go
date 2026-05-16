package standard_test

import (
	"math"
	"sync"
	"testing"

	"github.com/elecbug/netkit/v2/graph/standard"
)

// TestBarabasiAlbertGraph tests the Barabási-Albert graph generation function.
func TestBarabasiAlbertGraph(t *testing.T) {
	trial := 100
	n := 1000
	m := 3

	concurrency := 50
	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	var mu sync.Mutex

	totalDegree := 0
	totalMaxDegree := 0
	totalMinDegree := 0

	for i := 0; i < trial; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			g, err := standard.BarabasiAlbertGraph(i, false, standard.Unweighted(), n, m)
			if err != nil {
				t.Errorf("failed to generate Barabási-Albert graph: %v", err)
				return
			}

			if len(g.Nodes()) != n {
				t.Errorf("expected %d nodes, got %d", n, len(g.Nodes()))
			}

			localTotalDegree := 0
			localMaxDegree := 0
			localMinDegree := 1<<31 - 1

			for _, id := range g.Nodes() {
				node, err := g.Node(id)
				if err != nil {
					t.Errorf("failed to get node: %v", err)
					return
				}

				d := node.Degree()

				if d < m {
					t.Errorf("node %s has degree %d, expected at least %d", id, d, m)
				}

				localTotalDegree += d

				if d > localMaxDegree {
					localMaxDegree = d
				}

				if d < localMinDegree {
					localMinDegree = d
				}
			}

			mu.Lock()
			totalDegree += localTotalDegree
			totalMaxDegree += localMaxDegree
			totalMinDegree += localMinDegree
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	avgDegree := float64(totalDegree) / float64(n*trial)
	expectedAvgDegree := float64(2 * m)
	if !checkRange(avgDegree, expectedAvgDegree, 0.1, 0.1) {
		t.Errorf(
			"expected average degree around %f-%f, got %f",
			expectedAvgDegree*(1-0.1),
			expectedAvgDegree*(1+0.1),
			avgDegree,
		)
	}

	maxDegree := float64(totalMaxDegree) / float64(trial)
	expectedMaxDegree := float64(m) * math.Sqrt(float64(n))
	if !checkRange(maxDegree, expectedMaxDegree, 0.1, 0.1) {
		t.Errorf(
			"expected max degree around %f-%f, got %f",
			expectedMaxDegree*(1-0.1),
			expectedMaxDegree*(1+0.1),
			maxDegree,
		)
	}

	minDegree := float64(totalMinDegree) / float64(trial)
	expectedMinDegree := float64(m)
	if !checkRange(minDegree, expectedMinDegree, 0.1, 0.1) {
		t.Errorf("expected min degree around %f, got %f", expectedMinDegree, minDegree)
	}
}

// TestErdosRenyiGraph tests the Erdős-Rényi graph generation function.
func TestErdosRenyiGraph(t *testing.T) {
	trial := 100
	n := 1000
	p := 0.03

	concurrency := 50
	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	var mu sync.Mutex

	totalDegree := 0
	totalMaxDegree := 0
	totalMinDegree := 0

	for i := 0; i < trial; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			g, err := standard.ErdosRenyiGraph(i, false, standard.Unweighted(), n, p)
			if err != nil {
				t.Errorf("failed to generate Erdős-Rényi graph: %v", err)
				return
			}

			if len(g.Nodes()) != n {
				t.Errorf("expected %d nodes, got %d", n, len(g.Nodes()))
			}

			localTotalDegree := 0
			localMaxDegree := 0
			localMinDegree := 1<<31 - 1

			for _, id := range g.Nodes() {
				node, err := g.Node(id)
				if err != nil {
					t.Errorf("failed to get node: %v", err)
					return
				}

				d := node.Degree()

				localTotalDegree += d

				if d > localMaxDegree {
					localMaxDegree = d
				}

				if d < localMinDegree {
					localMinDegree = d
				}
			}

			mu.Lock()
			totalDegree += localTotalDegree
			totalMaxDegree += localMaxDegree
			totalMinDegree += localMinDegree
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	avgDegree := float64(totalDegree) / float64(n*trial)
	expectedAvgDegree := float64(n-1) * p
	if !checkRange(avgDegree, expectedAvgDegree, 0.1, 0.1) {
		t.Errorf(
			"expected average degree around %f-%f, got %f",
			expectedAvgDegree*(1-0.1),
			expectedAvgDegree*(1+0.1),
			avgDegree,
		)
	}

	muu := float64(n-1) * p
	sigma := math.Sqrt(float64(n-1) * p * (1 - p))
	extreme := sigma * math.Sqrt(2*math.Log(float64(n)))

	maxDegree := float64(totalMaxDegree) / float64(trial)
	expectedMaxDegree := muu + extreme
	if !checkRange(maxDegree, expectedMaxDegree, 0.25, 0.25) {
		t.Errorf(
			"expected max degree around %f-%f, got %f",
			expectedMaxDegree*(1-0.25),
			expectedMaxDegree*(1+0.25),
			maxDegree,
		)
	}

	minDegree := float64(totalMinDegree) / float64(trial)
	expectedMinDegree := muu - extreme
	if !checkRange(minDegree, expectedMinDegree, 0.25, 0.25) {
		t.Errorf(
			"expected min degree around %f-%f, got %f",
			expectedMinDegree*(1-0.25),
			expectedMinDegree*(1+0.25),
			minDegree,
		)
	}
}

// TestRandomGeometricGraph tests the random geometric graph generation function.
func TestRandomGeometricGraph(t *testing.T) {
	n := 100
	r := 0.1
	g, err := standard.RandomGeometricGraph(42, false, standard.Unweighted(), n, r)
	if err != nil {
		t.Fatalf("failed to generate random geometric graph: %v", err)
	}

	if len(g.Nodes()) != n {
		t.Errorf("expected %d nodes, got %d", n, len(g.Nodes()))
	}

	totalDegree := 0
	for _, id := range g.Nodes() {
		node, err := g.Node(id)
		if err != nil {
			t.Fatalf("failed to get node: %v", err)
		}
		d := node.Degree()
		totalDegree += d
	}

	expectedEdges := int(float64(n*(n-1)/2) * (3.14159 * r * r))
	if totalDegree/2 < expectedEdges/2 || totalDegree/2 > expectedEdges*2 {
		t.Errorf("expected around %d edges, got %d", expectedEdges, totalDegree/2)
	}
}

// checkRange checks if value is within the range of target*(1-lower) and target*(1+upper).
func checkRange(value, target float64, lower, upper float64) bool {
	return value >= target*(1-lower) && value <= target*(1+upper)
}
