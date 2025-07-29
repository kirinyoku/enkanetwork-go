// Example of using the Enka client with a custom HTTP client and in-memory cache.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/enka"
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
	client := enka.NewClient(httpClient, cache, "enkanetwork-go/1.0")

	// Define the Enka Network username to look up
	username := "Algoinde"

	// Fetch the user's profile information from EnkaNetwork
	// This includes basic information like bio, profile picture, etc.
	profile, err := client.GetUserProfile(ctx, username)
	if err != nil {
		switch err {
		case enka.ErrInvalidUsername:
			log.Fatalf("Invalid username format %q: %v", username, err)
		case enka.ErrUserNotFound:
			log.Fatalf("User not found for username %q: %v", username, err)
		default:
			log.Fatalf("Unexpected error fetching profile: %v", err)
		}
	}

	// Display the basic user information from the profile
	fmt.Printf("Profile:\n")
	fmt.Printf("Username: %s\n", profile.Username)
	fmt.Printf("Bio: %s\n", profile.Profile.Bio)

	// Fetch the user's hoyo accounts. A single EnkaNetwork user can have multiple hoyo accounts.
	// The endpoint will only return accounts that are both verified and public
	// (users can hide accounts; unverified accounts are hidden by default)
	hoyos, err := client.GetUserProfileHoyos(ctx, username)
	if err != nil {
		switch err {
		case enka.ErrInvalidUsername:
			log.Fatalf("Invalid username format %q: %v", username, err)
		case enka.ErrUserNotFound:
			log.Fatalf("User not found for username %q: %v", username, err)
		default:
			log.Fatalf("Unexpected error fetching hoyo accounts: %v", err)
		}
	}

	// Iterate through each hoyo account
	fmt.Printf("\nHoyo Accounts:\n")
	for hash, hoyo := range hoyos {
		// Display basic account information
		fmt.Printf("Account Hash: %s, Region: %s, UID: %d\n", hash, hoyo.Region, hoyo.UID)

		// Fetch detailed information for this specific hoyo account
		h, err := client.GetUserProfileHoyo(ctx, username, hash)
		if err != nil {
			switch err {
			case enka.ErrInvalidUsername:
				log.Fatalf("Invalid username format %q: %v", username, err)
			case enka.ErrInvalidHoyoHash:
				log.Fatalf("Invalid hoyo hash format %q: %v", hash, err)
			case enka.ErrHoyoAccountNotFound:
				log.Fatalf("Hoyo account not found for username %q and hash %q: %v", username, hash, err)
			default:
				log.Fatalf("Unexpected error fetching hoyo account details: %v", err)
			}
		}

		// Display player information from the hoyo account
		fmt.Printf("Nickname: %s\n", h.PlayerInfo.Nickname)
		fmt.Printf("World Level: %d\n", h.PlayerInfo.Level)

		// Fetch the character builds for this hoyo account
		// This includes all saved character builds across different games
		avatarBuilds, err := client.GetUserProfileHoyoBuilds(ctx, username, hash)
		if err != nil {
			switch err {
			case enka.ErrInvalidUsername:
				log.Fatalf("Invalid username format %q: %v", username, err)
			case enka.ErrInvalidHoyoHash:
				log.Fatalf("Invalid hoyo hash format %q: %v", hash, err)
			case enka.ErrHoyoAccountBuildsNotFound:
				log.Fatalf("Hoyo account builds not found for username %q and hash %q: %v", username, hash, err)
			default:
				log.Fatalf("Unexpected error fetching character builds: %v", err)
			}
		}

		// Display all character builds.
		// Each avatar can have multiple builds
		fmt.Printf("\nCharacter Builds:\n")
		for avatarID, builds := range avatarBuilds {
			fmt.Println("Builds for character ID:", avatarID)
			for _, build := range builds {
				switch build.HoyoType {
				case 0:
					fmt.Printf("Build ID %d (Genshin Impact)\n", build.ID)
				case 1:
					fmt.Printf("Build ID %d (Honkai: Star Rail)\n", build.ID)
				case 2:
					fmt.Printf("Build ID %d (Zenless Zone Zero)\n", build.ID)
				}
			}
			fmt.Println()
		}
	}

	// -----------------------------------------------------------------------
	// WARNING: This is just a demonstration of cache access.
	// In production code, you should NOT access the cache directly like this.
	// The cache is meant for internal use by the client methods only.
	// This example is shown purely for educational purposes.
	// -----------------------------------------------------------------------
	cacheKey := fmt.Sprintf("user_%s", username)
	data, ok := cache.Get(cacheKey)
	if !ok {
		log.Fatalf("Failed to get cached profile: %v", err)
	}

	if cachedProfile, ok := data.(*enka.Owner); ok {
		fmt.Printf("Cached username: %s\n", cachedProfile.Username)
	}
}
