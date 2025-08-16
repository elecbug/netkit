package g_algorithm

import (
	"sync"

	"github.com/elecbug/go-dspkg/network-graph/path"
)

var cachedAllShortestPaths map[string]path.GraphPaths = make(map[string]path.GraphPaths)
var cacheMu sync.RWMutex
