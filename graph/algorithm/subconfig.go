package algorithm

import "github.com/elecbug/netkit/graph/node"

// ClosenessCentralityConfig holds the configuration settings for the closeness centrality algorithm.
type ClosenessCentralityConfig struct {
	Reverse    bool
	WfImproved bool
}

// PageRankConfig holds the configuration settings for the PageRank algorithm.
type PageRankConfig struct {
	Alpha           float64              // damping, default 0.85
	MaxIter         int                  // default 100
	Tol             float64              // L1 error, default 1e-6
	Personalization *map[node.ID]float64 // p(u); if nil is uniform
	Dangling        *map[node.ID]float64 // d(u); if nil p(u)
	Reverse         bool
}

// BetweennessCentralityConfig holds the configuration settings for the edge betweenness centrality algorithm.
type BetweennessCentralityConfig struct {
	Normalized bool
}

// EdgeBetweennessCentralityConfig holds the configuration settings for the edge betweenness centrality algorithm.
type EdgeBetweennessCentralityConfig struct {
	Normalized bool
}

// EigenvectorCentralityConfig holds the configuration settings for the eigenvector centrality algorithm.
type EigenvectorCentralityConfig struct {
	MaxIter int
	Tol     float64
	Reverse bool
	NStart  *map[node.ID]float64 // initial vector; if nil, uniform distribution
}

// DegreeCentralityConfig holds the configuration settings for the degree centrality algorithm.
type DegreeCentralityConfig struct {
	Mode string
}

// AssortativityMode defines how degree pairs (j,k) are taken on each edge/arc.
// - Projected: ignore direction; use undirected degrees on both ends.
// - OutIn:     use out-degree(u) and in-degree(v) for each arc u->v
// - OutOut:    out-degree(u) and out-degree(v)
// - InIn:      in-degree(u) and in-degree(v)
// - InOut:     in-degree(u) and out-degree(v)
type AssortativityMode string

const (
	AssortativityProjected AssortativityMode = "projected"
	AssortativityOutIn     AssortativityMode = "out-in"
	AssortativityOutOut    AssortativityMode = "out-out"
	AssortativityInIn      AssortativityMode = "in-in"
	AssortativityInOut     AssortativityMode = "in-out"
)

// AssortativityCoefficientConfig holds the configuration settings for the assortativity coefficient algorithm.
type AssortativityCoefficientConfig struct {
	// Mode selects which degree pairing to use.
	// Defaults:
	//   - If graph is undirected: "projected"
	//   - If graph is directed:   "out-in"
	Mode AssortativityMode

	// IgnoreSelfLoops controls whether to ignore self loops (u==v).
	// Default: true
	IgnoreSelfLoops bool
}

// ModularityConfig holds the configuration settings for the modularity calculation.
type ModularityConfig struct {
	// Partition maps each node to its community ID.
	// If nil, algorithm will compute greedy modularity communities automatically.
	Partition map[node.ID]int
}
