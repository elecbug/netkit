package graph

import (
	"math"
)

// Distance represents the weight of an edge in a graph.
// It is defined as an unsigned integer for non-negative edge weights.
type Distance uint

// INF is a constant representing infinity.
// It is used to denote an unreachable state or maximum possible distance.
const INF = Distance(math.MaxUint)

// ToInt converts the Distance type to an int.
// This method is useful when an integer representation of the edge weight is required.
func (w Distance) ToInt() int {
	return int(w.uint())
}

// uint returns the Distance as a uint type.
// This allows direct access to the underlying unsigned integer value.
func (w Distance) uint() uint {
	return uint(w)
}
