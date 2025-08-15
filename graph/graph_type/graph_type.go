package graph_type

// GraphType describes whether a graph is directed/undirected and weighted/unweighted.
type GraphType int

// Enumeration values for GraphType.
const (
	DIRECTED_UNWEIGHTED   GraphType = iota // A directed graph with unweighted edges.
	DIRECTED_WEIGHTED                      // A directed graph with weighted edges.
	UNDIRECTED_UNWEIGHTED                  // An undirected graph with unweighted edges.
	UNDIRECTED_WEIGHTED                    // An undirected graph with weighted edges.
)

// String returns a human-readable representation of the graph type.
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
