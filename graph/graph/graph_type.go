package graph

// GraphType is an enumeration that defines the type of a graph.
// It specifies whether the graph is directed or undirected and whether it is weighted or unweighted.
type GraphType int

// Enumeration values for `GraphType`.
// These constants represent different types of graphs:
const (
	DIRECTED_UNWEIGHTED   GraphType = iota // A directed graph with unweighted edges.
	DIRECTED_WEIGHTED                      // A directed graph with weighted edges.
	UNDIRECTED_UNWEIGHTED                  // An undirected graph with unweighted edges.
	UNDIRECTED_WEIGHTED                    // An undirected graph with weighted edges.
)

// `String` converts a `GraphType` value to its string representation.
// This is useful for displaying the graph type in a human-readable format.
func (g GraphType) String() string {
	switch g {
	case DIRECTED_UNWEIGHTED:
		return "Directed Unweighted Graph" // Case for directed unweighted graph.
	case DIRECTED_WEIGHTED:
		return "Directed Weighted Graph" // Case for directed weighted graph.
	case UNDIRECTED_UNWEIGHTED:
		return "Undirected Unweighted Graph" // Case for undirected unweighted graph.
	case UNDIRECTED_WEIGHTED:
		return "Undirected Weighted Graph" // Case for undirected weighted graph.
	default:
		return "Unknown Graph Type" // Default case for unrecognized graph types.
	}
}
