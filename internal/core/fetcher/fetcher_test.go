package fetcher

import (
	"bytes"
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

func BenchmarkFetchWithRetryErrorBody(b *testing.B) {
	body := bytes.Repeat([]byte("x"), 1<<20)
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Header:     make(http.Header),
				Body:       io.NopCloser(readerOnly{Reader: bytes.NewReader(body)}),
				Request:    r,
			}, nil
		}),
	}
	fetcher := NewFetcher[struct{}](client, "test-agent", core.RetryOptions{MaxAttempts: 1})

	b.ReportAllocs()
	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fetcher.FetchWithRetry(context.Background(), "https://example.test"); err != coreerrors.ErrServiceUnavailable {
			b.Fatalf("error = %v, want %v", err, coreerrors.ErrServiceUnavailable)
		}
	}
}

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

func TestFetchWithRetryDrainsErrorBody(t *testing.T) {
	var bodyClosed int32
	body := &trackingBody{
		Reader: strings.NewReader("temporary response body"),
		closed: &bodyClosed,
	}
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return retryableResponseWithBody(r, http.StatusServiceUnavailable, body), nil
		}),
	}

	_, err := NewFetcher[struct{}](client, "test-agent", core.RetryOptions{MaxAttempts: 1}).FetchWithRetry(context.Background(), "https://example.test")
	if err != coreerrors.ErrServiceUnavailable {
		t.Fatalf("error = %v, want %v", err, coreerrors.ErrServiceUnavailable)
	}
	if body.Len() != 0 {
		t.Fatalf("response body has %d unread bytes", body.Len())
	}
	if atomic.LoadInt32(&bodyClosed) == 0 {
		t.Fatal("response body was not closed")
	}
}

func TestWaitForRetryReturnsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := waitForRetry(ctx, time.Hour); err != context.Canceled {
		t.Fatalf("waitForRetry() error = %v, want %v", err, context.Canceled)
	}
}

func TestWaitForRetryStopsDuringActiveWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelTimer := time.AfterFunc(10*time.Millisecond, cancel)
	defer cancelTimer.Stop()

	started := time.Now()
	if err := waitForRetry(ctx, time.Hour); err != context.Canceled {
		t.Fatalf("waitForRetry() error = %v, want %v", err, context.Canceled)
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("waitForRetry() took %v after cancellation", elapsed)
	}
}

func TestParseRetryAfter(t *testing.T) {
	fallback := 5 * time.Second
	tests := []struct {
		name  string
		value string
		want  time.Duration
	}{
		{name: "seconds", value: "15", want: 15 * time.Second},
		{name: "seconds with whitespace", value: " 15 ", want: 15 * time.Second},
		{name: "negative seconds", value: "-1", want: fallback},
		{name: "overflow", value: "9223372037", want: fallback},
		{name: "invalid", value: "later", want: fallback},
		{name: "past HTTP date", value: time.Now().Add(-time.Hour).UTC().Format(time.RFC1123), want: 0},
		{name: "RFC 850 date", value: "Sunday, 06-Nov-94 08:49:37 GMT", want: 0},
		{name: "ANSI C date", value: "Sun Nov  6 08:49:37 1994", want: 0},
		{name: "RFC 1123 UTC fallback", value: "Sun, 06 Nov 1994 08:49:37 UTC", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRetryAfter(tt.value, fallback); got != tt.want {
				t.Fatalf("parseRetryAfter(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}

	future := time.Now().Add(time.Hour).UTC().Format(time.RFC1123)
	if got := parseRetryAfter(future, fallback); got < 59*time.Minute || got > time.Hour {
		t.Fatalf("parseRetryAfter(future date) = %v, want about one hour", got)
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

type readerOnly struct {
	io.Reader
}
