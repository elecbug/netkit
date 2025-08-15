package graph

import "math"

// Distance represents the weight of an edge in a graph.
type Distance int

// INF_DISTANCE is a sentinel representing an unreachable or infinite distance.
const INF_DISTANCE = Distance(math.MaxInt)
