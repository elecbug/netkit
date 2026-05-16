package standard

import (
	"math/rand"

	"github.com/elecbug/netkit/v2/graph"
)

// WeightedFunc defines a function type for generating edge weights based on node IDs.
type WeightedFunc func(from, to graph.NodeID) *graph.Weight

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

// Unweighted returns a WeightedFunc that generates unweighted edges (i.e., all weights are 1.0).
func Unweighted() WeightedFunc {
	return func(from, to graph.NodeID) *graph.Weight {
		return graph.NewWeight(1.0)
	}
}
