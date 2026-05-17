package standard_test

import (
	"fmt"
	"math"
	"sync"
	"testing"

	"github.com/elecbug/netkit/v2/graph"
	"github.com/elecbug/netkit/v2/graph/standard"
)

func TestStandardGraph(t *testing.T) {
	fmt.Println("Test Standard Graph Generation")
	testGridGraph(t)
	testTriangleHexGraph(t)
	testBarabasiAlbertGraph(t)
	testErdosRenyiGraph(t)
	testRandomGeometricGraph(t)
	testRandomRegularGraph(t)
	testWattsStrogatzGraph(t)

	fmt.Println("Test Standard Graph Generation from Config")
	testGenerateFromConfig(t)
}

// testBarabasiAlbertGraph tests the Barabási-Albert graph generation function.
func testBarabasiAlbertGraph(t *testing.T) {
	fmt.Println("- Test Barabási-Albert Graph")

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

			g, err := standard.BarabasiAlbertGraph(i, false, standard.Unweighted, n, m)
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

// testErdosRenyiGraph tests the Erdős-Rényi graph generation function.
func testErdosRenyiGraph(t *testing.T) {
	fmt.Println("- Test Erdős-Rényi Graph")

	trial := 100
	n := 1000
	p := 0.2

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

			g, err := standard.ErdosRenyiGraph(i, false, standard.Unweighted, n, p)
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
	if !checkRange(maxDegree, expectedMaxDegree, 0.1, 0.1) {
		t.Errorf(
			"expected max degree around %f-%f, got %f",
			expectedMaxDegree*(1-0.1),
			expectedMaxDegree*(1+0.1),
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

// testRandomGeometricGraph tests the random geometric graph generation function.
func testRandomGeometricGraph(t *testing.T) {
	fmt.Println("- Test Random Geometric Graph")

	trial := 100
	n := 1000
	targetDegree := 20.0
	r, err := standard.RForRandomGeometricGraph(targetDegree, n)
	if err != nil {
		t.Fatalf("failed to calculate radius for random geometric graph: %v", err)
	}

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

			g, err := standard.RandomGeometricGraph(i, false, standard.Unweighted, n, r)
			if err != nil {
				t.Errorf("failed to generate random geometric graph: %v", err)
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
	expectedAvgDegree := targetDegree
	if !checkRange(avgDegree, expectedAvgDegree, 0.1, 0.1) {
		t.Errorf(
			"expected average degree around %f-%f, got %f",
			expectedAvgDegree*(1-0.1),
			expectedAvgDegree*(1+0.1),
			avgDegree,
		)
	}

	maxDegree := float64(totalMaxDegree) / float64(trial)
	expectedMaxDegree := float64(poissonUpperExtreme(n, targetDegree))
	if !checkRange(maxDegree, expectedMaxDegree, 0.1, 0.1) {
		t.Errorf(
			"expected max degree around %f-%f, got %f",
			expectedMaxDegree*(1-0.1),
			expectedMaxDegree*(1+0.1),
			maxDegree,
		)
	}

	minDegree := float64(totalMinDegree) / float64(trial)
	expectedMinDegree := float64(poissonLowerExtreme(n, targetDegree/2))
	if !checkRange(minDegree, expectedMinDegree, 1.0, 1.0) {
		t.Errorf(
			"expected min degree around %f-%f, got %f",
			expectedMinDegree*(1-1.0),
			expectedMinDegree*(1+1.0),
			minDegree,
		)
	}
}

// testRandomRegularGraph tests the random regular graph generation function.
func testRandomRegularGraph(t *testing.T) {
	fmt.Println("- Test Random Regular Graph")

	trial := 100
	n := 1000
	k := 20.0

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

			g, err := standard.RandomRegularGraph(i, false, standard.Unweighted, n, int(k))
			if err != nil {
				t.Errorf("failed to generate random regular graph: %v", err)
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
	expectedAvgDegree := k
	if !checkRange(avgDegree, expectedAvgDegree, 0.0, 0.0) {
		t.Errorf(
			"expected average degree around %f-%f, got %f",
			expectedAvgDegree*(1-0.0),
			expectedAvgDegree*(1+0.0),
			avgDegree,
		)
	}

	maxDegree := float64(totalMaxDegree) / float64(trial)
	expectedMaxDegree := k
	if !checkRange(maxDegree, expectedMaxDegree, 0.0, 0.0) {
		t.Errorf(
			"expected max degree around %f-%f, got %f",
			expectedMaxDegree*(1-0.0),
			expectedMaxDegree*(1+0.0),
			maxDegree,
		)
	}

	minDegree := float64(totalMinDegree) / float64(trial)
	expectedMinDegree := k
	if !checkRange(minDegree, expectedMinDegree, 0.0, 0.0) {
		t.Errorf(
			"expected min degree around %f-%f, got %f",
			expectedMinDegree*(1-0.0),
			expectedMinDegree*(1+0.0),
			minDegree,
		)
	}
}

// testWattsStrogatzGraph tests the Watts-Strogatz graph generation function.
func testWattsStrogatzGraph(t *testing.T) {
	fmt.Println("- Test Watts-Strogatz Graph")

	trial := 100
	n := 1000
	k := 20.0
	beta := 0.1

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

			g, err := standard.WattsStrogatzGraph(i, false, standard.Unweighted, n, int(k), beta)
			if err != nil {
				t.Errorf("failed to generate Watts-Strogatz graph: %v", err)
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
	expectedAvgDegree := k*(1-beta) + k*beta
	if !checkRange(avgDegree, expectedAvgDegree, 0.1, 0.1) {
		t.Errorf(
			"expected average degree around %f-%f, got %f",
			expectedAvgDegree*(1-0.1),
			expectedAvgDegree*(1+0.1),
			avgDegree,
		)
	}

	q := k / 2
	sigma := math.Sqrt(q*beta*(1-beta) + beta*q)
	extreme := sigma * math.Sqrt(2*math.Log(float64(n)))

	maxDegree := float64(totalMaxDegree) / float64(trial)
	expectedMaxDegree := k + extreme
	if !checkRange(maxDegree, expectedMaxDegree, 0.1, 0.1) {
		t.Errorf(
			"expected max degree around %f-%f, got %f",
			expectedMaxDegree*(1-0.1),
			expectedMaxDegree*(1+0.1),
			maxDegree,
		)
	}

	minDegree := float64(totalMinDegree) / float64(trial)
	expectedMinDegree := k - extreme
	if !checkRange(minDegree, expectedMinDegree, 0.1, 0.1) {
		t.Errorf(
			"expected min degree around %f-%f, got %f",
			expectedMinDegree*(1-0.1),
			expectedMinDegree*(1+0.1),
			minDegree,
		)
	}
}

// testGridGraph tests the grid graph generation function.
func testGridGraph(t *testing.T) {
	fmt.Println("- Test Grid Graph")

	rows := 10
	cols := 10
	torus := false

	g, err := standard.GridGraph(0, false, standard.Unweighted, rows, cols, torus)
	if err != nil {
		t.Fatalf("failed to generate grid graph: %v", err)
	}

	expectedNodes := rows * cols
	if len(g.Nodes()) != expectedNodes {
		t.Errorf("expected %d nodes, got %d", expectedNodes, len(g.Nodes()))
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			id := graph.NodeID(fmt.Sprintf("%d", i*cols+j))
			node, err := g.Node(id)
			if err != nil {
				t.Errorf("failed to get node: %v", err)
				continue
			}

			x, okX := node.Tag("x")
			y, okY := node.Tag("y")
			if !okX || !okY {
				t.Errorf("node %s is missing tags", id)
				continue
			}

			if x != fmt.Sprintf("%d", i) || y != fmt.Sprintf("%d", j) {
				t.Errorf("node %s has incorrect tags: x=%s, y=%s", id, x, y)
			}

			expectedDegree := 4
			if i == 0 || i == rows-1 {
				expectedDegree--
			}
			if j == 0 || j == cols-1 {
				expectedDegree--
			}

			if node.Degree() != expectedDegree {
				t.Errorf("node %s has degree %d, expected %d", id, node.Degree(), expectedDegree)
			}
		}
	}
}

// testTriangleHexGraph tests the triangle hex graph generation function.
func testTriangleHexGraph(t *testing.T) {
	fmt.Println("- Test Triangle Hex Graph")

	edge := 3

	g, err := standard.TriangleHexGraph(0, false, standard.Unweighted, edge)
	if err != nil {
		t.Fatalf("failed to generate triangle hex graph: %v", err)
	}

	expectedNodes := 3*edge*(edge+1)/2 + 1
	if len(g.Nodes()) != expectedNodes {
		t.Errorf("expected %d nodes, got %d", expectedNodes, len(g.Nodes()))
	}

	for i := 0; i < expectedNodes; i++ {
		nodeID := graph.NodeID(fmt.Sprintf("%d", i))

		node, err := g.Node(nodeID)
		if err != nil {
			t.Errorf("failed to get node: %v", err)
			continue
		}

		q, okQ := node.Tag("q")
		r, okR := node.Tag("r")
		if !okQ || !okR {
			t.Errorf("node %s is missing tags", nodeID)
			continue
		}

		qInt := 0
		rInt := 0
		fmt.Sscanf(q, "%d", &qInt)
		fmt.Sscanf(r, "%d", &rInt)

		expectedDegree := degreeOfTriangleHex(qInt, rInt, edge)

		if node.Degree() != expectedDegree {
			t.Errorf("node %s (q=%d, r=%d) has degree %d, expected %d", nodeID, qInt, rInt, node.Degree(), expectedDegree)
		}
	}
}

// testGenerateFromConfig tests the StandardGraph function with various configurations.
func testGenerateFromConfig(t *testing.T) {
	fmt.Println("- Test StandardGraph with Config")

	configs := []standard.GraphConfig{
		{
			Type: standard.ErdosRenyi,
			Params: map[string]interface{}{
				"n": 1000,
				"p": 0.2,
			},
		},
		{
			Type: standard.BarabasiAlbert,
			Params: map[string]interface{}{
				"n": 1000,
				"m": 3,
			},
		},
		{
			Type: standard.WattsStrogatz,
			Params: map[string]interface{}{
				"n":    1000,
				"k":    20,
				"beta": 0.1,
			},
		},
	}

	for i, config := range configs {
		g, err := standard.StandardGraph(i, false, standard.Unweighted, config)
		if err != nil {
			t.Errorf("failed to generate graph for config %d: %v", i, err)
			continue
		}

		if len(g.Nodes()) != 1000 {
			t.Errorf("expected 1000 nodes for config %d, got %d", i, len(g.Nodes()))
		}
	}

	invalidConfigs := []standard.GraphConfig{
		{
			Type: standard.ErdosRenyi,
			Params: map[string]interface{}{
				"n": -1,
				"p": 0.2,
			},
		},
		{
			Type: standard.BarabasiAlbert,
			Params: map[string]interface{}{
				"n": 10,
				"m": 20,
			},
		},
		{
			Type: standard.WattsStrogatz,
			Params: map[string]interface{}{
				"n":    1000,
				"k":    20,
				"beta": -0.1,
			},
		},
	}

	for i, config := range invalidConfigs {
		if _, err := standard.StandardGraph(i, false, standard.Unweighted, config); err == nil {
			t.Errorf("expected error for invalid config %d, but got none", i)
		}
	}
}

// checkRange checks if value is within the range of target*(1-lower) and target*(1+upper).
func checkRange(value, target float64, lower, upper float64) bool {
	return value >= target*(1-lower) && value <= target*(1+upper)
}

// poissonCDF calculates the cumulative distribution function of the Poisson distribution for k and lambda.
func poissonCDF(k int, lambda float64) float64 {
	if k < 0 {
		return 0
	}

	term := math.Exp(-lambda)
	sum := term

	for i := 1; i <= k; i++ {
		term *= lambda / float64(i)
		sum += term
	}

	return sum
}

// poissonUpperExtreme calculates the upper extreme value for a Poisson distribution with mean lambda and n samples.
func poissonUpperExtreme(n int, lambda float64) int {
	target := 1.0 - 1.0/float64(n)

	for k := 0; k < 100000; k++ {
		if poissonCDF(k, lambda) >= target {
			return k
		}
	}

	return -1
}

// poissonLowerExtreme calculates the lower extreme value for a Poisson distribution with mean lambda and n samples.
func poissonLowerExtreme(n int, lambda float64) int {
	target := 1.0 / float64(n)

	for k := 0; k < 100000; k++ {
		if poissonCDF(k, lambda) >= target {
			return k
		}
	}

	return -1
}

// degreeOfTriangleHex calculates the degree of a node in the triangle hex graph based on its q and r coordinates and the edge length.
func degreeOfTriangleHex(q, r, n int) int {
	dirs := [][2]int{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
		{1, -1},
		{-1, 1},
	}

	degree := 0

	for _, d := range dirs {
		nq := q + d[0]
		nr := r + d[1]

		if exists(nq, nr, n) {
			degree++
		}
	}

	return degree
}

// exists checks if the coordinates (q, r) are valid for a node in the triangle hex graph with edge length n.
func exists(q, r, n int) bool {
	limit := n - 1

	return q >= -limit && q <= limit &&
		r >= -limit && r <= limit &&
		q+r >= -limit && q+r <= limit
}
