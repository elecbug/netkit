package standard_graph

import (
	"fmt"
	"math/rand"
)

const randCode = 42

type StandardGraph struct {
	seed int64
}

func NewStandardGraph() *StandardGraph {
	return &StandardGraph{
		seed: randCode,
	}
}

// SetSeed sets the seed for random operations in the graph.
func (g *StandardGraph) SetSeed(value int64) {
	g.seed = value
}

// SetSeedRandom sets the seed to a random value for non-deterministic behavior.
func (g *StandardGraph) SetSeedRandom() {
	g.seed = randCode
}

func (g *StandardGraph) genRand() *rand.Rand {
	localSeed := g.seed

	if g.seed == randCode {
		localSeed = rand.Int63()
	}

	return rand.New(rand.NewSource(localSeed))
}

func toString(id int) string {
	return fmt.Sprintf("%d", id)
}
