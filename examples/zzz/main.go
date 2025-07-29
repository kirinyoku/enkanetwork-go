// Example of using the Zenless Zone Zero client with a custom HTTP client and in-memory cache.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
)

// Simple in-memory cache
type Cache struct {
	data map[string]cacheEntry
	mu   sync.RWMutex
}

type cacheEntry struct {
	value     any
	expiresAt time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]cacheEntry),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.value, true
}

func (c *Cache) Set(key string, value any, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(expiration),
	}
}

func main() {
	// Create a context with a 15-second timeout to prevent hanging indefinitely
	// This ensures the program won't run forever if the API is unresponsive
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Configure a custom HTTP client with sensible defaults:
	// - 10 second timeout for individual requests
	// - Connection pooling with up to 10 idle connections
	// - 30 second idle connection timeout
	// - Compression enabled for better performance
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,               // Maximum number of idle connections to keep alive
			IdleConnTimeout:    30 * time.Second, // How long to keep idle connections
			DisableCompression: false,            // Enable gzip compression for responses
		},
	}

	// Initialize an in-memory cache to store API responses
	// This reduces the number of API calls and improves performance
	cache := NewCache()

	// Create a new Enka client with our custom HTTP client and cache
	// The User-Agent string helps Enka Network identify the source of API requests
	client := zzz.NewClient(httpClient, cache, "enkanetwork-go/1.0")

	// Define the UID of the player to fetch.
	const uid = "618285856"

	// Fetch the full Profile, which includes basic player info and characters showcase.
	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		switch err {
		case zzz.ErrInvalidUIDFormat:
			log.Fatalf("Invalid UID format %q: %v", uid, err)
		case zzz.ErrPlayerNotFound:
			log.Fatalf("Player not found for UID %q: %v", uid, err)
		case zzz.ErrRateLimited:
			log.Fatalf("Rate limit exceeded: %v", err)
		case zzz.ErrServerMaintenance:
			log.Fatalf("Server under maintenance: %v", err)
		default:
			log.Fatalf("Unexpected error fetching profile: %v", err)
		}
	}

	// Display basic player details.
	fmt.Printf("Player Nickname: %s\n", profile.PlayerInfo.SocialDetail.ProfileDetail.Nickname)
	fmt.Printf("World Level: %d\n", profile.PlayerInfo.SocialDetail.ProfileDetail.Level)

	// Display character showcase details.
	fmt.Printf("Characters in showcase (%d):\n", len(profile.PlayerInfo.ShowcaseDetail.AvatarList))
	for _, avatar := range profile.PlayerInfo.ShowcaseDetail.AvatarList {
		fmt.Printf("- ID: %d, Level: %d\n", avatar.ID, avatar.Level)
	}

	// -----------------------------------------------------------------------
	// WARNING: This is just a demonstration of cache access.
	// In production code, you should NOT access the cache directly like this.
	// The cache is meant for internal use by the client methods only.
	// This example is shown purely for educational purposes.
	// -----------------------------------------------------------------------
	cacheKey := fmt.Sprintf("zzz_%s", uid)
	data, ok := cache.Get(cacheKey)
	if !ok {
		log.Fatalf("Failed to get cached profile: %v", err)
	}

	if cachedProfile, ok := data.(*zzz.Profile); ok {
		fmt.Printf("Cached username: %s\n", cachedProfile.PlayerInfo.SocialDetail.ProfileDetail.Nickname)
	}
}
