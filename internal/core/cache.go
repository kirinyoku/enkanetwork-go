package core

import "time"

// Cache defines an interface for caching API responses.
// Caching helps reduce the number of requests to the API, which is important because
// even cached responses from the API count toward rate limits. Users can implement
// this interface to provide their own caching mechanism, such as an in-memory cache
// or a database.
type Cache interface {
	// Get retrieves a value from the cache by key.
	// Returns the cached value and true if found,
	// or nil and false if not found or expired.
	Get(key string) (any, bool)
	// Set stores a value in the cache with the given key and expiration time.
	// The expiration time determines how long the value remains valid.
	Set(key string, value any, expiration time.Duration)
}
