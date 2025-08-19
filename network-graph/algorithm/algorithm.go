// Package algorithm provides graph algorithms for network analysis.
package algorithm

import (
	"sync"

	"github.com/elecbug/netkit/network-graph/path"
)

var cachedAllShortestPaths = make(map[string]path.GraphPaths)
var cachedAllShortestPathLengths = make(map[string]path.PathLength)
var cacheMu sync.RWMutex
