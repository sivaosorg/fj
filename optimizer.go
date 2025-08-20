package fj

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cached query result
type CacheEntry struct {
	Result    Context   `json:"result"`
	Timestamp time.Time `json:"timestamp"`
	HitCount  int64     `json:"hitCount"`
}

// QueryCache provides caching for query results using LRU eviction
type QueryCache struct {
	entries    map[string]*CacheEntry
	usage      []string // LRU tracking
	maxSize    int
	maxAge     time.Duration
	mu         sync.RWMutex
	stats      CacheStats
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	Hits       int64 `json:"hits"`
	Misses     int64 `json:"misses"`
	Evictions  int64 `json:"evictions"`
	Size       int   `json:"size"`
}

// NewQueryCache creates a new query cache with specified parameters
func NewQueryCache(maxSize int, maxAge time.Duration) *QueryCache {
	return &QueryCache{
		entries: make(map[string]*CacheEntry),
		usage:   make([]string, 0, maxSize),
		maxSize: maxSize,
		maxAge:  maxAge,
	}
}

// Get retrieves a cached result for the given query key
func (c *QueryCache) Get(key string) (Context, bool) {
	c.mu.RLock()
	entry, exists := c.entries[key]
	c.mu.RUnlock()
	
	if !exists {
		c.mu.Lock()
		c.stats.Misses++
		c.mu.Unlock()
		return Context{}, false
	}
	
	// Check if entry has expired
	if c.maxAge > 0 && time.Since(entry.Timestamp) > c.maxAge {
		c.mu.Lock()
		delete(c.entries, key)
		c.removeFromUsage(key)
		c.stats.Misses++
		c.stats.Evictions++
		c.mu.Unlock()
		return Context{}, false
	}
	
	// Update usage tracking and stats
	c.mu.Lock()
	entry.HitCount++
	c.moveToFront(key)
	c.stats.Hits++
	c.mu.Unlock()
	
	return entry.Result, true
}

// Put stores a query result in the cache
func (c *QueryCache) Put(key string, result Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if key already exists
	if _, exists := c.entries[key]; exists {
		c.entries[key].Result = result
		c.entries[key].Timestamp = time.Now()
		c.moveToFront(key)
		return
	}
	
	// If cache is at capacity, evict LRU entry
	if len(c.entries) >= c.maxSize {
		c.evictLRU()
	}
	
	// Add new entry
	entry := &CacheEntry{
		Result:    result,
		Timestamp: time.Now(),
		HitCount:  0,
	}
	
	c.entries[key] = entry
	c.usage = append([]string{key}, c.usage...)
	c.stats.Size = len(c.entries)
}

// Clear removes all entries from the cache
func (c *QueryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.entries = make(map[string]*CacheEntry)
	c.usage = c.usage[:0]
	c.stats = CacheStats{}
}

// Stats returns current cache statistics
func (c *QueryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	stats := c.stats
	stats.Size = len(c.entries)
	return stats
}

// moveToFront moves the key to the front of the usage list (most recently used)
func (c *QueryCache) moveToFront(key string) {
	// Remove from current position
	c.removeFromUsage(key)
	// Add to front
	c.usage = append([]string{key}, c.usage...)
}

// removeFromUsage removes a key from the usage tracking list
func (c *QueryCache) removeFromUsage(key string) {
	for i, k := range c.usage {
		if k == key {
			c.usage = append(c.usage[:i], c.usage[i+1:]...)
			break
		}
	}
}

// evictLRU removes the least recently used entry
func (c *QueryCache) evictLRU() {
	if len(c.usage) > 0 {
		lru := c.usage[len(c.usage)-1]
		delete(c.entries, lru)
		c.usage = c.usage[:len(c.usage)-1]
		c.stats.Evictions++
	}
}

// PathCache caches parsed path metadata to avoid re-parsing
type PathCache struct {
	entries map[string]*metadata
	mu      sync.RWMutex
	maxSize int
	stats   CacheStats
}

// NewPathCache creates a new path parsing cache
func NewPathCache(maxSize int) *PathCache {
	return &PathCache{
		entries: make(map[string]*metadata),
		maxSize: maxSize,
	}
}

// Get retrieves cached path metadata
func (c *PathCache) Get(path string) (*metadata, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if entry, exists := c.entries[path]; exists {
		c.stats.Hits++
		return entry, true
	}
	
	c.stats.Misses++
	return nil, false
}

// Put stores parsed path metadata in cache
func (c *PathCache) Put(path string, meta *metadata) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Simple eviction - clear cache when full
	if len(c.entries) >= c.maxSize {
		c.entries = make(map[string]*metadata)
		c.stats.Evictions++
	}
	
	// Create a copy of metadata to avoid race conditions
	metaCopy := *meta
	c.entries[path] = &metaCopy
	c.stats.Size = len(c.entries)
}

// Stats returns current cache statistics
func (c *PathCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	stats := c.stats
	stats.Size = len(c.entries)
	return stats
}

// Clear removes all entries from the path cache
func (c *PathCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.entries = make(map[string]*metadata)
	c.stats = CacheStats{}
}

// QueryOptimizer provides various optimizations for query execution
type QueryOptimizer struct {
	queryCache *QueryCache
	pathCache  *PathCache
	mu         sync.RWMutex
	enabled    bool
}

// NewQueryOptimizer creates a new query optimizer with default settings
func NewQueryOptimizer() *QueryOptimizer {
	return &QueryOptimizer{
		queryCache: NewQueryCache(1000, 5*time.Minute), // 1000 entries, 5 minute TTL
		pathCache:  NewPathCache(500),                   // 500 parsed paths
		enabled:    true,
	}
}

// IsEnabled returns whether optimization is enabled
func (o *QueryOptimizer) IsEnabled() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.enabled
}

// SetEnabled enables or disables query optimization
func (o *QueryOptimizer) SetEnabled(enabled bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.enabled = enabled
}

// OptimizeQuery applies various optimizations to query execution
func (o *QueryOptimizer) OptimizeQuery(ctx Context, query *Query) Context {
	if !o.IsEnabled() {
		// Bypass optimization - execute directly
		return o.executeWithoutCache(ctx, query)
	}
	
	// Generate cache key for this query
	cacheKey := o.generateCacheKey(ctx, query)
	
	// Try to get result from cache
	if result, found := o.queryCache.Get(cacheKey); found {
		return result
	}
	
	// Execute query with path caching optimization
	result := o.executeWithPathCache(ctx, query)
	
	// Store result in cache if it's valid
	if result.Exists() {
		o.queryCache.Put(cacheKey, result)
	}
	
	return result
}

// executeWithoutCache executes the query without any caching
func (o *QueryOptimizer) executeWithoutCache(ctx Context, query *Query) Context {
	processor := GetQueryProcessor()
	return processor.ExecuteQuery(ctx, query)
}

// executeWithPathCache executes the query using path caching optimization
func (o *QueryOptimizer) executeWithPathCache(ctx Context, query *Query) Context {
	// Check if path metadata is cached
	if !isEmptyString(query.Path) {
		if cachedMeta, found := o.pathCache.Get(query.Path); found {
			// Use cached metadata
			query.parsedPath = cachedMeta
		} else {
			// Parse and cache path metadata
			meta := analyzePath(query.Path)
			query.parsedPath = &meta
			o.pathCache.Put(query.Path, query.parsedPath)
		}
	}
	
	// Execute with cached path metadata
	processor := GetQueryProcessor()
	return processor.ExecuteQuery(ctx, query)
}

// generateCacheKey creates a unique cache key for the query and context
func (o *QueryOptimizer) generateCacheKey(ctx Context, query *Query) string {
	// Create a hash based on context and query parameters
	hasher := md5.New()
	
	// Add context data (simplified - could be more sophisticated)
	hasher.Write([]byte(ctx.unprocessed))
	
	// Add query parameters
	hasher.Write([]byte(query.Path))
	
	// Add filters
	for _, filter := range query.Filters {
		hasher.Write([]byte(fmt.Sprintf("%s%s%v", filter.Field, filter.Operator, filter.Value)))
	}
	
	// Add selections
	for _, selection := range query.Selections {
		hasher.Write([]byte(selection))
	}
	
	// Add other parameters
	hasher.Write([]byte(fmt.Sprintf("%s%d%d", query.OrderBy, query.Limit, query.Offset)))
	
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// GetCacheStats returns statistics for both query and path caches
func (o *QueryOptimizer) GetCacheStats() (queryStats, pathStats CacheStats) {
	return o.queryCache.Stats(), o.pathCache.Stats()
}

// ClearCaches clears both query and path caches
func (o *QueryOptimizer) ClearCaches() {
	o.queryCache.Clear()
	o.pathCache.Clear()
}

// Global optimizer instance
var (
	globalOptimizer *QueryOptimizer
	optimizerOnce   sync.Once
)

// GetQueryOptimizer returns the global query optimizer instance
func GetQueryOptimizer() *QueryOptimizer {
	optimizerOnce.Do(func() {
		globalOptimizer = NewQueryOptimizer()
	})
	return globalOptimizer
}

// OptimizedExecute executes a query using the global optimizer
func OptimizedExecute(ctx Context, query *Query) Context {
	optimizer := GetQueryOptimizer()
	return optimizer.OptimizeQuery(ctx, query)
}

// Performance monitoring and metrics
type PerformanceMetrics struct {
	QueryExecutionTime time.Duration `json:"queryExecutionTime"`
	CacheHitRatio      float64       `json:"cacheHitRatio"`
	PathCacheHitRatio  float64       `json:"pathCacheHitRatio"`
	MemoryUsage        int64         `json:"memoryUsage"`
}

// GetPerformanceMetrics returns current performance metrics
func (o *QueryOptimizer) GetPerformanceMetrics() PerformanceMetrics {
	queryStats := o.queryCache.Stats()
	pathStats := o.pathCache.Stats()
	
	var queryCacheHitRatio, pathCacheHitRatio float64
	
	if queryStats.Hits+queryStats.Misses > 0 {
		queryCacheHitRatio = float64(queryStats.Hits) / float64(queryStats.Hits+queryStats.Misses)
	}
	
	if pathStats.Hits+pathStats.Misses > 0 {
		pathCacheHitRatio = float64(pathStats.Hits) / float64(pathStats.Hits+pathStats.Misses)
	}
	
	return PerformanceMetrics{
		CacheHitRatio:     queryCacheHitRatio,
		PathCacheHitRatio: pathCacheHitRatio,
		MemoryUsage:       int64(queryStats.Size + pathStats.Size),
	}
}