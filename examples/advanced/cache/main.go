package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

// memoryCache is a simple thread-safe in-memory cache implementation.
//
// Recommendation: This implementation will grow indefinitely because it never deletes expired keys.
// In a production environment, you should:
//  1. Run a background goroutine (e.g., every minute) to periodically delete expired keys.
//  2. Alternatively, use a production-ready library like github.com/patrickmn/go-cache,
//     github.com/dgraph-io/ristretto, or a distributed cache like Redis.
type memoryCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	value     any
	expiresAt time.Time
}

func newMemoryCache() *memoryCache {
	return &memoryCache{
		entries: make(map[string]cacheEntry),
	}
}

// Get implements the Cache interface.
//
// Note: In an in-memory cache, returning `entry.value` returns a pointer to the original struct.
// You should treat the returned objects from the Enka API as read-only to avoid data races
// when multiple goroutines access the same cached profile.
func (c *memoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.value, true
}

// Set implements the Cache interface.
//
// The `expiration` duration is automatically calculated by the library based on the
// `TTL` field provided in the EnkaNetwork API response. You don't need to calculate it manually!
func (c *memoryCache) Set(key string, value any, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(expiration),
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := genshin.New(genshin.Options{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		// Pluggable Cache:
		// Once provided, the client will automatically check the cache before making an HTTP request.
		Cache:     newMemoryCache(),
		UserAgent: "enkanetwork-go-example/1.0",
		Retry: &genshin.RetryOptions{
			MaxAttempts: 2,
			Delay:       2 * time.Second,
		},
	})

	const uid = "618285856"

	// Request 1: Cache Miss
	// The first call will hit the API and store the result in the cache.
	fmt.Println("Fetching profile for the first time (Cache Miss)...")
	start := time.Now()
	profile1, err := client.GetProfile(ctx, uid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First request took: %v\n", time.Since(start))
	fmt.Printf("Nickname: %s\n\n", profile1.PlayerInfo.Nickname)

	// Request 2: Cache Hit
	// The second call for the same UID within the TTL will return instantly from memory.
	fmt.Println("Fetching profile again (Cache Hit)...")
	start = time.Now()
	profile2, err := client.GetProfile(ctx, uid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Second request took: %v\n", time.Since(start))
	fmt.Printf("Nickname: %s\n", profile2.PlayerInfo.Nickname)
}
