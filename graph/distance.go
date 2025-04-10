package graph

// `Distance` represents the weight of an edge in a graph.
// It is defined as an unsigned integer for non-negative edge weights.
type Distance int

// `INF_DISTANCE` is a constant representing infinity.
// It is used to denote an unreachable state or maximum possible distance.
const INF_DISTANCE = Distance(-1)
