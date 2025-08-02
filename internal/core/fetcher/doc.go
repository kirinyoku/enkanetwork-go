// Package fetcher provides a generic HTTP client with retry logic for API requests.
//
// The fetcher package implements a flexible and reusable HTTP client that handles:
//   - Automatic retries with backoff for failed requests
//   - Rate limiting with Retry-After header support
//   - Common error handling patterns
//   - Request timeouts and cancellation via context
//
// The main type, Fetcher[T], is a generic client that can be used to make HTTP requests
// and unmarshal JSON responses into the specified type T.
package fetcher
