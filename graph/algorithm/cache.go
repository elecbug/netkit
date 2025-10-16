package algorithm

import (
	"sync"
	"time"

	"github.com/elecbug/netkit/graph"
)

var cachedAllShortestPaths = make(map[string]graph.Paths)
var cachedAllShortestPathLengths = make(map[string]graph.PathLength)
var cacheMu sync.RWMutex

// CacheClear clears the cached shortest paths and their lengths.
func CacheClear() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cachedAllShortestPaths = make(map[string]graph.Paths)
	cachedAllShortestPathLengths = make(map[string]graph.PathLength)
}

// AutoCacheClear starts a goroutine that clears the cache at regular intervals defined by tick.
func AutoCacheClear(tick time.Duration) {
	go func() {
		for {
			time.Sleep(tick)

			CacheClear()
		}
	}()
}
