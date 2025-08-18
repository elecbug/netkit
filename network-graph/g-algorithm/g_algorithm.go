// Package g_algorithm provides graph algorithms for network analysis.
package g_algorithm

import (
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/path"
)

var cachedAllShortestPaths = make(map[string]path.GraphPaths)
var cachedAllShortestPathLengths = make(map[string]path.PathLength)
var cacheMu sync.RWMutex
