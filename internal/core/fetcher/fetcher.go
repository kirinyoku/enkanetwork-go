package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core/errors"
)

const (
	maxRetries        = 3               // maxRetries defines the maximum number of retry attempts for failed requests
	defaultRetryDelay = 5 * time.Second // defaultRetryDelay is the default delay between retry attempts
)

// Fetcher is a generic HTTP client that handles request retries and error handling.
// The type parameter T specifies the type to unmarshal the JSON response into.
type Fetcher[T any] struct {
	client    *http.Client
	userAgent string
}

// NewFetcher creates a new Fetcher instance with the specified HTTP client and user agent.
// The HTTP client should be configured with appropriate timeouts and transport settings.
// The user agent string will be included in all requests.
func NewFetcher[T any](client *http.Client, userAgent string) *Fetcher[T] {
	return &Fetcher[T]{
		client:    client,
		userAgent: userAgent,
	}
}

// FetchWithRetry executes an HTTP GET request to the specified URL with retry logic for transient errors.
// It handles:
// - Request timeouts and cancellation via the provided context.
// - Automatic retries for server errors (500, 503) and rate limiting (429).
// - Rate limiting by respecting the Retry-After header if present.
// - Specific error mapping for common HTTP status codes (400, 404, 424, 500, 503).
//
// Parameters:
//   - ctx: Context for controlling request timeout and cancellation.
//   - url: The URL to fetch the resource from.
//
// Returns:
//   - *T: A pointer to the unmarshaled response body of type T on success.
//   - error: An error if the request fails after all retries or encounters a non-retryable error.
//
// Possible errors:
//   - errors.ErrInvalidUIDFormat: For 400 Bad Request
//   - errors.ErrPlayerNotFound: For 404 Not Found
//   - errors.ErrServerMaintenance: For 424 Failed Dependency
//   - errors.ErrServerError: For 500 Internal Server Error (if received outside retries)
//   - errors.ErrServiceUnavailable: For 503 Service Unavailable (if received outside retries)
//   - errors.ErrRateLimited: When retries are exhausted due to transient errors (429, 500, 503)
//
// The function attempts up to maxRetries times for transient errors (429, 500, 503).
// If retries are exhausted, it returns errors.ErrRateLimited.
// For other error status codes, it returns immediately with the corresponding error.
func (f *Fetcher[T]) FetchWithRetry(ctx context.Context, url string) (*T, error) {
	for attempt := range maxRetries {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", f.userAgent)

		resp, err := f.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result T

			err = json.Unmarshal(body, &result)
			if err != nil {
				return nil, fmt.Errorf("failed to decode profile: %w", err)
			}

			return &result, nil
		}

		// Check for retryable status codes: 429 (Too Many Requests), 500 (Internal Server Error), 503 (Service Unavailable)
		if resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode == http.StatusInternalServerError ||
			resp.StatusCode == http.StatusServiceUnavailable {
			// If not the last attempt, calculate delay and retry
			if attempt < maxRetries-1 {
				delay := defaultRetryDelay
				// For 429 and 503, attempt to parse Retry-After header for custom delay
				if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
					retryAfter := resp.Header.Get("Retry-After")
					if retryAfter != "" {
						delay = parseRetryAfter(retryAfter)
					}
				}
				// Wait for the calculated delay or exit if context is canceled
				select {
				case <-time.After(delay):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
		} else {
			switch resp.StatusCode {
			case 400:
				return nil, errors.ErrInvalidUIDFormat
			case 404:
				return nil, errors.ErrPlayerNotFound
			case 424:
				return nil, errors.ErrServerMaintenance
			case 500:
				return nil, errors.ErrServerError
			case 503:
				return nil, errors.ErrServiceUnavailable
			default:
				return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
			}
		}
	}

	return nil, errors.ErrRateLimited
}

// parseRetryAfter parses the Retry-After header value into a time.Duration.
// It handles both:
//   - Integer values (seconds)
//   - HTTP date strings (RFC 1123 format)
//
// If parsing fails or the date is in the past, it returns the defaultRetryDelay.
func parseRetryAfter(retryAfter string) time.Duration {
	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		return time.Duration(seconds) * time.Second
	}

	if date, err := time.Parse(time.RFC1123, retryAfter); err == nil {
		delay := time.Until(date)
		if delay < 0 {
			return 0 // Retry immediately if the date is in the past
		}
		return delay
	}

	return defaultRetryDelay // Default if parsing fails
}
