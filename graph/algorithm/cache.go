package algorithm

import (
	"sync"
	"time"

	"github.com/elecbug/netkit/graph"
)

var cachedAllShortestPaths = make(map[string]graph.Paths)
var cachedAllShortestPathLengths = make(map[string]graph.PathLength)
var cacheMu sync.RWMutex

// ClearCache clears the cached shortest paths and their lengths.
func ClearCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cachedAllShortestPaths = make(map[string]graph.Paths)
	cachedAllShortestPathLengths = make(map[string]graph.PathLength)
}

// ClearCacheForGraph clears the cache for a specific graph.
func ClearCacheForGraph(g *graph.Graph) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	gh := g.Hash()

	delete(cachedAllShortestPaths, gh)
	delete(cachedAllShortestPathLengths, gh)
}

// AutoClearCache starts a goroutine that clears the cache at regular intervals defined by tick.
func AutoClearCache(tick time.Duration) {
	go func() {
		for {
			time.Sleep(tick)

			ClearCache()
		}
	}()
}
