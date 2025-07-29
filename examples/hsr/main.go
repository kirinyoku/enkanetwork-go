// Example of using the HSR client with a custom HTTP client and in-memory cache.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/hsr"
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
	client := hsr.NewClient(httpClient, cache, "enkanetwork-go/1.0")

	// Define the UID of the player to fetch.
	const uid = "800579959"

	// Perform the API request to fetch the PlayerInfo by UID.
	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		switch err {
		case hsr.ErrInvalidUIDFormat:
			log.Fatalf("Invalid UID format %q: %v", uid, err)
		case hsr.ErrPlayerNotFound:
			log.Fatalf("Player not found for UID %q: %v", uid, err)
		case hsr.ErrRateLimited:
			log.Fatalf("Rate limit exceeded: %v", err)
		case hsr.ErrServerMaintenance:
			log.Fatalf("Server under maintenance: %v", err)
		default:
			log.Fatalf("Unexpected error fetching profile: %v", err)
		}
	}

	// Print some example information from the player's profile.
	log.Printf("Player Nickname: %s", profile.DetailInfo.Nickname)
	log.Printf("Level: %d", profile.DetailInfo.Level)
	log.Printf("World Level: %d", profile.DetailInfo.WorldLevel)

	// Check and display character showcase details if available.
	if len(profile.DetailInfo.AvatarDetailList) == 0 {
		fmt.Println("No character showcase is public or available.")
		return
	}

	// Display character showcase details.
	fmt.Printf("Characters in showcase (%d):\n", len(profile.DetailInfo.AvatarDetailList))
	for _, avatar := range profile.DetailInfo.AvatarDetailList {
		fmt.Printf("- %d (Level %d, Promotion %d)\n", avatar.AvatarID, avatar.Level, avatar.Promotion)
	}

	// -----------------------------------------------------------------------
	// WARNING: This is just a demonstration of cache access.
	// In production code, you should NOT access the cache directly like this.
	// The cache is meant for internal use by the client methods only.
	// This example is shown purely for educational purposes.
	// -----------------------------------------------------------------------
	cacheKey := fmt.Sprintf("hsr_%s", uid)
	data, ok := cache.Get(cacheKey)
	if !ok {
		log.Fatalf("Failed to get cached profile: %v", err)
	}

	if cachedProfile, ok := data.(*hsr.Profile); ok {
		fmt.Printf("Cached username: %s\n", cachedProfile.DetailInfo.Nickname)
	}
}
