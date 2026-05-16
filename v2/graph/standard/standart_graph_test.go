package standard_test

import (
	"fmt"
	"math"
	"sync"
	"testing"

	"github.com/elecbug/netkit/v2/graph/standard"
)

// TestBarabasiAlbertGraph tests the Barabási-Albert graph generation function.
func TestBarabasiAlbertGraph(t *testing.T) {
	fmt.Println("Test Barabási-Albert Graph")

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
	fmt.Println("Test Erdős-Rényi Graph")

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

// TestRandomGeometricGraph tests the random geometric graph generation function.
func TestRandomGeometricGraph(t *testing.T) {
	fmt.Println("Test Random Geometric Graph")

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

			g, err := standard.RandomGeometricGraph(i, false, standard.Unweighted(), n, r)
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

// TestRandomRegularGraph tests the random regular graph generation function.
func TestRandomRegularGraph(t *testing.T) {
	fmt.Println("Test Random Regular Graph")

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

			g, err := standard.RandomRegularGraph(i, false, standard.Unweighted(), n, int(k))
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

// TestWattsStrogatzGraph tests the Watts-Strogatz graph generation function.
func TestWattsStrogatzGraph(t *testing.T) {
	fmt.Println("Test Watts-Strogatz Graph")

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

			g, err := standard.WattsStrogatzGraph(i, false, standard.Unweighted(), n, int(k), beta)
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

// TestGenerateFromConfig tests the StandardGraph function with various configurations.
func TestGenerateFromConfig(t *testing.T) {
	fmt.Println("Test StandardGraph with Config")

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
		g, err := standard.StandardGraph(i, false, standard.Unweighted(), config)
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
		if _, err := standard.StandardGraph(i, false, standard.Unweighted(), config); err == nil {
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
