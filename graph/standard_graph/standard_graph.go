// Package standard_graph provides a standard implementation of a graph data structure.
package standard_graph

// STANDARD_GRAPH_TYPE defines types of standard graphs.
type STANDARD_GRAPH_TYPE int

const (
	ERDOS_RENYI STANDARD_GRAPH_TYPE = iota
	RANDOM_REGULAR
	BARABASI_ALBERT
	WATTS_STROGATZ
	RANDOM_GEOMETRIC
	WAXMAN
)

// String returns the string representation of the STANDARD_GRAPH_TYPE.
func (s STANDARD_GRAPH_TYPE) String(onlyAlphabet bool) string {
	switch s {
	case ERDOS_RENYI:
		if onlyAlphabet {
			return "Erdos-Renyi"
		} else {
			return "Erdős-Rényi"
		}
	case RANDOM_REGULAR:
		return "Random Regular"
	case BARABASI_ALBERT:
		if onlyAlphabet {
			return "Barabasi-Albert"
		} else {
			return "Barabási-Albert"
		}
	case WATTS_STROGATZ:
		return "Watts-Strogatz"
	case RANDOM_GEOMETRIC:
		return "Random Geometric"
	case WAXMAN:
		return "Waxman"
	default:
		return "Unknown"
	}
}
