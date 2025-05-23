// Package cache provides caching mechanisms for API responses
// to improve performance and reduce external API calls.
package cache

import (
	"math"
	"sort"
	"sync"
	"time"
)

// Item represents a cached item with expiration
type Item struct {
	Value      interface{}
	Expiration int64
}

// Expired checks if the item has expired
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// TTLCache is a thread-safe cache with time-based expiration
type TTLCache struct {
	items           map[string]Item
	mu              sync.RWMutex
	defaultTTL      time.Duration
	cleanupInterval time.Duration
	maxItems        int
	stopCleanup     chan bool
}

// NewTTLCache creates a new cache with the specified TTL and cleanup interval
// maxItems specifies the maximum number of items before oldest are evicted
func NewTTLCache(defaultTTL, cleanupInterval time.Duration, maxItems int) *TTLCache {
	cache := &TTLCache{
		items:           make(map[string]Item),
		defaultTTL:      defaultTTL,
		cleanupInterval: cleanupInterval,
		maxItems:        maxItems,
		stopCleanup:     make(chan bool),
	}

	// Start the cleanup process
	cache.startCleanupTimer()

	return cache
}

// Set adds an item to the cache with the default TTL
func (c *TTLCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL adds an item to the cache with a specific TTL
func (c *TTLCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	var expiration int64

	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}

	// If we're over capacity, remove oldest items
	if c.maxItems > 0 && len(c.items) > c.maxItems {
		c.evictOldest()
	}
}

// Get retrieves an item from the cache
// Returns the item and a bool indicating if the item was found
func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	// Check if the item has expired
	if item.Expired() {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache
func (c *TTLCache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Count returns the number of items in the cache
func (c *TTLCache) Count() int {
	c.mu.RLock()
	count := len(c.items)
	c.mu.RUnlock()
	return count
}

// Clear removes all items from the cache
func (c *TTLCache) Clear() {
	c.mu.Lock()
	c.items = make(map[string]Item)
	c.mu.Unlock()
}

// evictOldest removes the oldest items when cache exceeds maxItems
// This function assumes the lock is already held
func (c *TTLCache) evictOldest() {
	// Create a slice of keys and their expiration times
	type keyExpiration struct {
		key        string
		expiration int64
	}

	// Calculate how many items to remove
	itemsToRemove := len(c.items) - c.maxItems
	if itemsToRemove <= 0 {
		return
	}

	// Collect all key expirations
	keyExpirations := make([]keyExpiration, 0, len(c.items))
	for k, v := range c.items {
		// Use MaxInt64 for items without expiration to treat them as lowest eviction priority
		exp := v.Expiration
		if exp == 0 {
			exp = math.MaxInt64
		}
		keyExpirations = append(keyExpirations, keyExpiration{k, exp})
	}

	// Sort by expiration time (oldest first)
	sort.Slice(keyExpirations, func(i, j int) bool {
		return keyExpirations[i].expiration < keyExpirations[j].expiration
	})

	// Delete the oldest items
	for i := 0; i < itemsToRemove; i++ {
		delete(c.items, keyExpirations[i].key)
	}
}

// startCleanupTimer starts the cleanup timer
func (c *TTLCache) startCleanupTimer() {
	if c.cleanupInterval <= 0 {
		return
	}

	ticker := time.NewTicker(c.cleanupInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.deleteExpired()
			case <-c.stopCleanup:
				ticker.Stop()
				return
			}
		}
	}()
}

// deleteExpired deletes all expired items
func (c *TTLCache) deleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && v.Expiration < now {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Stop stops the cleanup timer
func (c *TTLCache) Stop() {
	close(c.stopCleanup)
}

// Global cache instance
var (
	globalCache     *TTLCache
	globalCacheOnce sync.Once
)

// GetGlobalCache returns the global cache instance
func GetGlobalCache() *TTLCache {
	globalCacheOnce.Do(func() {
		// 5 minute TTL, cleanup every minute, max 1000 items
		globalCache = NewTTLCache(5*time.Minute, time.Minute, 1000)
	})
	return globalCache
}
