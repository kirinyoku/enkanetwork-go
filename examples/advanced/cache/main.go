// Advanced example of custom HTTP, retry, and cache configuration.
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

func (c *memoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.value, true
}

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
		Cache:     newMemoryCache(),
		UserAgent: "enkanetwork-go-example/1.0",
		Retry: &genshin.RetryOptions{
			MaxAttempts: 2,
			Delay:       2 * time.Second,
		},
	})

	profile, err := client.GetProfile(ctx, "618285856")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Nickname: %s\n", profile.PlayerInfo.Nickname)
	fmt.Printf("Showcase Characters: %d\n", len(profile.AvatarInfoList))
}
