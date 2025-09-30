package standard_graph

import (
	"fmt"
	"math/rand"
)

const randCode = 42

var seed = int64(randCode)

// SetSeed sets the seed for random operations in the graph.
func SetSeed(value int64) {
	seed = value
}

// SetSeedRandom sets the seed to a random value for non-deterministic behavior.
func SetSeedRandom() {
	seed = randCode
}

func genRand() *rand.Rand {
	localSeed := seed

	if seed == randCode {
		localSeed = rand.Int63()
	}

	return rand.New(rand.NewSource(localSeed))
}

func toString(id int) string {
	return fmt.Sprintf("%d", id)
}
