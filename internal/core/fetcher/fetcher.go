package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/internal/core/errors"
)

// Fetcher is a generic HTTP client that handles request retries and error handling.
// The type parameter T specifies the type to unmarshal the JSON response into.
type Fetcher[T any] struct {
	client      *http.Client
	userAgent   string
	maxAttempts int
	retryDelay  time.Duration
}

// NewFetcher creates a new Fetcher instance with the specified HTTP client and user agent.
// The HTTP client should be configured with appropriate timeouts and transport settings.
// The user agent string will be included in all requests.
func NewFetcher[T any](client *http.Client, userAgent string, retry core.RetryOptions) *Fetcher[T] {
	if retry.MaxAttempts <= 0 {
		retry.MaxAttempts = core.DefaultMaxAttempts
	}
	if retry.Delay <= 0 {
		retry.Delay = core.DefaultRetryDelay
	}
	return &Fetcher[T]{
		client:      client,
		userAgent:   userAgent,
		maxAttempts: retry.MaxAttempts,
		retryDelay:  retry.Delay,
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
//   - errors.ErrServerError: When retries are exhausted due to 500 responses
//   - errors.ErrServiceUnavailable: When retries are exhausted due to 503 responses
//   - errors.ErrRateLimited: When retries are exhausted due to 429 responses
//
// The function attempts up to the configured maximum attempts for transient errors (429, 500, 503).
// If retries are exhausted, it returns the error matching the last retryable status.
// For other error status codes, it returns immediately with the corresponding error.
func (f *Fetcher[T]) FetchWithRetry(ctx context.Context, url string) (*T, error) {
	for attempt := 0; attempt < f.maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", f.userAgent)

		resp, err := f.client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		if err := resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("failed to close response body: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result T

			err = json.Unmarshal(body, &result)
			if err != nil {
				return nil, fmt.Errorf("failed to decode profile: %w", err)
			}

			return &result, nil
		}

		if isRetryableStatus(resp.StatusCode) {
			if attempt < f.maxAttempts-1 {
				delay := f.retryDelay
				if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
					retryAfter := resp.Header.Get("Retry-After")
					if retryAfter != "" {
						delay = parseRetryAfter(retryAfter, f.retryDelay)
					}
				}
				select {
				case <-time.After(delay):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}

			return nil, errorForStatus(resp.StatusCode)
		}

		return nil, errorForStatus(resp.StatusCode)
	}

	return nil, errors.ErrRateLimited
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusServiceUnavailable
}

func errorForStatus(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return errors.ErrInvalidUIDFormat
	case http.StatusNotFound:
		return errors.ErrPlayerNotFound
	case http.StatusFailedDependency:
		return errors.ErrServerMaintenance
	case http.StatusTooManyRequests:
		return errors.ErrRateLimited
	case http.StatusInternalServerError:
		return errors.ErrServerError
	case http.StatusServiceUnavailable:
		return errors.ErrServiceUnavailable
	default:
		return fmt.Errorf("unexpected status: %d", statusCode)
	}
}

// parseRetryAfter parses the Retry-After header value into a time.Duration.
// It handles both:
//   - Integer values (seconds)
//   - HTTP date strings (RFC 1123 format)
//
// If parsing fails or the date is in the past, it returns the fallback delay.
func parseRetryAfter(retryAfter string, fallback time.Duration) time.Duration {
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

	return fallback // Default if parsing fails
}
