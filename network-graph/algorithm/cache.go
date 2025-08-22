package algorithm

import (
	"sync"

	"github.com/elecbug/netkit/network-graph/path"
)

var cachedAllShortestPaths = make(map[string]path.GraphPaths)
var cachedAllShortestPathLengths = make(map[string]path.PathLength)
var cacheMu sync.RWMutex

// CacheClear clears the cached shortest paths and their lengths.
func CacheClear() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cachedAllShortestPaths = make(map[string]path.GraphPaths)
	cachedAllShortestPathLengths = make(map[string]path.PathLength)
}
