package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/kirinyoku/enkanetwork-go/internal/core/errors"
)

// Fetcher is a generic HTTP client that handles request retries and error handling.
type Fetcher struct {
	client      *http.Client
	userAgent   string
	maxAttempts int
	retryDelay  time.Duration
}

// NewFetcher creates a new Fetcher instance with the specified HTTP client and user agent.
// The HTTP client should be configured with appropriate timeouts and transport settings.
// The user agent string will be included in all requests.
func NewFetcher(client *http.Client, userAgent string, retry RetryOptions) *Fetcher {
	if retry.MaxAttempts <= 0 {
		retry.MaxAttempts = DefaultMaxAttempts
	}
	if retry.Delay <= 0 {
		retry.Delay = DefaultRetryDelay
	}
	return &Fetcher{
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
func FetchWithRetry[T any](ctx context.Context, f *Fetcher, url string) (*T, error) {
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

		if resp.StatusCode == http.StatusOK {
			body, err := readAndCloseResponseBody(resp.Body)
			if err != nil {
				return nil, err
			}

			var result T
			if err := json.Unmarshal(body, &result); err != nil {
				return nil, fmt.Errorf("failed to decode profile: %w", err)
			}

			return &result, nil
		}
		if err := discardAndCloseResponseBody(resp.Body); err != nil {
			return nil, err
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
				if err := waitForRetry(ctx, delay); err != nil {
					return nil, err
				}
				continue
			}

			return nil, errorForStatus(resp.StatusCode)
		}

		return nil, errorForStatus(resp.StatusCode)
	}

	return nil, errors.ErrRateLimited
}

// FetchAndCache executes an HTTP GET request with retry logic and caches the result.
// It checks the cache first. If the item is found, it performs a safe type assertion
// and returns the typed result. If not found, it calls FetchWithRetry to get the data,
// and if successful and a cache is provided, it stores the result in the cache using
// the TTL provided by the model itself via the Cacheable interface.
func FetchAndCache[T Cacheable](ctx context.Context, f *Fetcher, url string, cacheKey string, cache Cache) (*T, error) {
	if cache != nil {
		if cached, ok := cache.Get(cacheKey); ok {
			if typed, ok := cached.(*T); ok {
				return typed, nil
			}
		}
	}

	result, err := FetchWithRetry[T](ctx, f, url)
	if err != nil {
		return nil, err
	}

	if cache != nil {
		cache.Set(cacheKey, result, (*result).CacheTTL())
	}

	return result, nil
}

func readAndCloseResponseBody(body io.ReadCloser) ([]byte, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		_ = body.Close()
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	return data, nil
}

func discardAndCloseResponseBody(body io.ReadCloser) error {
	_, err := io.Copy(io.Discard, body)
	if err != nil {
		_ = body.Close()
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := body.Close(); err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	return nil
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
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
//   - All HTTP date formats accepted by net/http, plus RFC 1123 dates using UTC
//
// If parsing fails, it returns the fallback delay. Dates in the past produce a
// zero delay so the next attempt can start immediately.
func parseRetryAfter(retryAfter string, fallback time.Duration) time.Duration {
	retryAfter = strings.TrimSpace(retryAfter)
	if seconds, err := strconv.ParseUint(retryAfter, 10, 64); err == nil {
		maxSeconds := uint64((time.Duration(1<<63 - 1)) / time.Second)
		if seconds > maxSeconds {
			return fallback
		}
		return time.Duration(seconds) * time.Second
	}

	date, err := http.ParseTime(retryAfter)
	if err != nil {
		// time.RFC1123 also accepts UTC, which older versions of this package
		// accepted even though HTTP dates conventionally use GMT.
		date, err = time.Parse(time.RFC1123, retryAfter)
	}
	if err == nil {
		delay := time.Until(date)
		if delay < 0 {
			return 0
		}
		return delay
	}

	return fallback
}
