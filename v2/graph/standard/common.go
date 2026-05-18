package standard

import (
	"fmt"
	"math/rand"

	"github.com/elecbug/netkit/v2/graph"
)

// WeightedFunc defines a function type for generating edge weights based on node IDs.
type WeightedFunc func(from, to *graph.Node) *graph.Weight

// GraphType represents the type of graph to be generated.
type GraphType string

const (
	Grid            GraphType = "grid"
	TriangleHex     GraphType = "triangle_hex"
	ErdosRenyi      GraphType = "erdos_renyi"
	BarabasiAlbert  GraphType = "barabasi_albert"
	WattsStrogatz   GraphType = "watts_strogatz"
	RandomGeometric GraphType = "random_geometric"
	RandomRegular   GraphType = "random_regular"
	None            GraphType = "none"
)

// GraphConfig represents a configuration for graph generation, allowing for flexible parameters.
type GraphConfig struct {
	Type   GraphType
	Params map[string]interface{}
}

// generateRand creates a new rand.Rand instance based on the provided seed.
func generateRand(seed int) *rand.Rand {
	var randSource rand.Source
	if seed == 42 {
		randSource = rand.NewSource(rand.Int63())
	} else {
		randSource = rand.NewSource(int64(seed))
	}

	r := rand.New(randSource)

	return r
}

// StandardGraph generates a graph based on the provided configuration. It supports various graph types and parameters.
func StandardGraph(seed int, directed bool, weightFunc WeightedFunc, config GraphConfig) (*graph.Graph, error) {
	switch config.Type {
	case Grid:
		rows, ok := config.Params["rows"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'rows' for grid graph")
		}
		cols, ok := config.Params["cols"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'cols' for grid graph")
		}
		torus, ok := config.Params["torus"].(bool)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'torus' for grid graph")
		}
		return GridGraph(seed, directed, weightFunc, rows, cols, torus)
	case TriangleHex:
		edge, ok := config.Params["edge"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'edge' for triangle hex graph")
		}
		return TriangleHexGraph(seed, directed, weightFunc, edge)
	case ErdosRenyi:
		n, ok := config.Params["n"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'n' for Erdos-Renyi graph")
		}
		p, ok := config.Params["p"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'p' for Erdos-Renyi graph")
		}
		return ErdosRenyiGraph(seed, directed, weightFunc, n, p)
	case BarabasiAlbert:
		n, ok := config.Params["n"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'n' for Barabasi-Albert graph")
		}
		m, ok := config.Params["m"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'm' for Barabasi-Albert graph")
		}
		return BarabasiAlbertGraph(seed, directed, weightFunc, n, m)
	case WattsStrogatz:
		n, ok := config.Params["n"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'n' for Watts-Strogatz graph")
		}
		k, ok := config.Params["k"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'k' for Watts-Strogatz graph")
		}
		beta, ok := config.Params["beta"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'beta' for Watts-Strogatz graph")
		}
		return WattsStrogatzGraph(seed, directed, weightFunc, n, k, beta)
	case RandomGeometric:
		n, ok := config.Params["n"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'n' for Random Geometric graph")
		}

		_, okR := config.Params["r"].(float64)
		_, okK := config.Params["k"].(float64)

		if okR && !okK {
			radius, _ := config.Params["r"].(float64)
			return RandomGeometricGraph(seed, directed, weightFunc, n, radius)
		} else if !okR && okK {
			k, _ := config.Params["k"].(float64)
			radius, err := RForRandomGeometricGraph(k, n)
			if err != nil {
				return nil, err
			}

			return RandomGeometricGraph(seed, directed, weightFunc, n, radius)
		} else {
			return nil, fmt.Errorf("invalid parameters for random geometric graph: must provide either 'r' or 'k'")
		}
	case RandomRegular:
		n, ok := config.Params["n"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'n' for Random Regular graph")
		}
		k, ok := config.Params["k"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid parameter 'k' for Random Regular graph")
		}
		return RandomRegularGraph(seed, directed, weightFunc, n, k)
	case None:
		return graph.New(directed, true), nil
	default:
		return nil, fmt.Errorf("unsupported graph type: %s", config.Type)
	}
}
