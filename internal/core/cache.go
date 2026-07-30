package core

import "time"

// Cache defines an interface for caching API responses.
// Caching helps reduce the number of requests to the API, which is important because
// even cached responses from the API count toward rate limits. Users can implement
// this interface to provide their own caching mechanism, such as an in-memory cache
// or a database.
//
// Values are stored and returned as-is. The client may store pointers to response
// structs. Cache implementations that need mutation isolation should copy,
// serialize, or otherwise protect values on Set and Get.
type Cache interface {
	// Get retrieves a value from the cache by key.
	// Returns the cached value and true if found,
	// or nil and false if not found or expired.
	Get(key string) (any, bool)
	// Set stores a value in the cache with the given key and expiration time.
	// The expiration time determines how long the value remains valid.
	// The value is provided as-is and may be a pointer to a response struct.
	Set(key string, value any, expiration time.Duration)
}

// Cacheable defines an interface for models that can provide their own cache TTL.
type Cacheable interface {
	// CacheTTL returns the duration for which the model should be cached.
	CacheTTL() time.Duration
}
