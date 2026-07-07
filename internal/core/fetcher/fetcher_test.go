package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	coreerrors "github.com/kirinyoku/enkanetwork-go/internal/core/errors"
)

func TestFetchWithRetryReturnsRateLimitedAfter429Retries(t *testing.T) {
	attempts := 0
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			attempts++
			return retryableResponse(r, http.StatusTooManyRequests), nil
		}),
	}

	_, err := NewFetcher[struct{}](client, "test-agent", core.NormalizeRetryOptions(nil)).FetchWithRetry(context.Background(), "https://example.test")
	if err != coreerrors.ErrRateLimited {
		t.Fatalf("error = %v, want %v", err, coreerrors.ErrRateLimited)
	}
	if attempts != core.DefaultMaxAttempts {
		t.Fatalf("attempts = %d, want %d", attempts, core.DefaultMaxAttempts)
	}
}

func TestFetchWithRetryReturnsServiceUnavailableAfter503Retries(t *testing.T) {
	attempts := 0
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			attempts++
			return retryableResponse(r, http.StatusServiceUnavailable), nil
		}),
	}

	_, err := NewFetcher[struct{}](client, "test-agent", core.NormalizeRetryOptions(nil)).FetchWithRetry(context.Background(), "https://example.test")
	if err != coreerrors.ErrServiceUnavailable {
		t.Fatalf("error = %v, want %v", err, coreerrors.ErrServiceUnavailable)
	}
	if attempts != core.DefaultMaxAttempts {
		t.Fatalf("attempts = %d, want %d", attempts, core.DefaultMaxAttempts)
	}
}

func TestFetchWithRetryUsesConfiguredMaxAttempts(t *testing.T) {
	attempts := 0
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			attempts++
			return retryableResponse(r, http.StatusTooManyRequests), nil
		}),
	}

	retry := core.NormalizeRetryOptions(&core.RetryOptions{MaxAttempts: 2})
	_, err := NewFetcher[struct{}](client, "test-agent", retry).FetchWithRetry(context.Background(), "https://example.test")
	if err != coreerrors.ErrRateLimited {
		t.Fatalf("error = %v, want %v", err, coreerrors.ErrRateLimited)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestFetchWithRetryCanBeDisabled(t *testing.T) {
	attempts := 0
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			attempts++
			return retryableResponse(r, http.StatusServiceUnavailable), nil
		}),
	}

	retry := core.NormalizeRetryOptions(&core.RetryOptions{MaxAttempts: 1})
	_, err := NewFetcher[struct{}](client, "test-agent", retry).FetchWithRetry(context.Background(), "https://example.test")
	if err != coreerrors.ErrServiceUnavailable {
		t.Fatalf("error = %v, want %v", err, coreerrors.ErrServiceUnavailable)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestFetchWithRetryClosesBodyBeforeRetry(t *testing.T) {
	var firstBodyClosed int32
	attempts := 0
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				return retryableResponseWithBody(r, http.StatusServiceUnavailable, &trackingBody{
					Reader: strings.NewReader(`temporary`),
					closed: &firstBodyClosed,
				}), nil
			}
			if atomic.LoadInt32(&firstBodyClosed) == 0 {
				return nil, fmt.Errorf("first response body was not closed before retry")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"name":"ok"}`)),
				Request:    r,
			}, nil
		}),
	}

	got, err := NewFetcher[struct {
		Name string `json:"name"`
	}](client, "test-agent", core.NormalizeRetryOptions(nil)).FetchWithRetry(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("FetchWithRetry() error = %v", err)
	}
	if got.Name != "ok" {
		t.Fatalf("Name = %q, want ok", got.Name)
	}
}

func retryableResponse(r *http.Request, statusCode int) *http.Response {
	return retryableResponseWithBody(r, statusCode, io.NopCloser(strings.NewReader(`temporary`)))
}

func retryableResponseWithBody(r *http.Request, statusCode int, body io.ReadCloser) *http.Response {
	header := make(http.Header)
	header.Set("Retry-After", time.Now().Add(-time.Hour).UTC().Format(time.RFC1123))
	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       body,
		Request:    r,
	}
}

type trackingBody struct {
	*strings.Reader
	closed *int32
}

func (b *trackingBody) Close() error {
	atomic.StoreInt32(b.closed, 1)
	return nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
