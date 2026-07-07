package core

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClientStoresImmutableConfiguration(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Second}
	cache := testCache{}

	client := NewClient(Options{
		HTTPClient: httpClient,
		Cache:      cache,
		UserAgent:  "test-agent",
		BaseURL:    "https://example.test/api/",
		Retry: &RetryOptions{
			MaxAttempts: 1,
			Delay:       time.Second,
		},
	})

	if client.HTTPClient() != httpClient {
		t.Fatal("HTTPClient() did not return configured HTTP client")
	}
	if client.Cache() != cache {
		t.Fatal("Cache() did not return configured cache")
	}
	if client.UserAgent() != "test-agent" {
		t.Fatalf("UserAgent() = %q, want test-agent", client.UserAgent())
	}
	if client.BaseURL() != "https://example.test/api" {
		t.Fatalf("BaseURL() = %q, want trimmed base URL", client.BaseURL())
	}
	if client.Retry().MaxAttempts != 1 || client.Retry().Delay != time.Second {
		t.Fatalf("Retry() = %+v, want configured retry options", client.Retry())
	}
}

func TestNormalizeRetryOptionsUsesDefaults(t *testing.T) {
	got := NormalizeRetryOptions(nil)
	if got.MaxAttempts != DefaultMaxAttempts {
		t.Fatalf("MaxAttempts = %d, want %d", got.MaxAttempts, DefaultMaxAttempts)
	}
	if got.Delay != DefaultRetryDelay {
		t.Fatalf("Delay = %s, want %s", got.Delay, DefaultRetryDelay)
	}
}

type testCache struct{}

func (testCache) Get(string) (any, bool) {
	return nil, false
}

func (testCache) Set(string, any, time.Duration) {}

func TestNormalizeRetryOptionsPreservesConfiguredValues(t *testing.T) {
	got := NormalizeRetryOptions(&RetryOptions{
		MaxAttempts: 1,
		Delay:       time.Second,
	})
	if got.MaxAttempts != 1 {
		t.Fatalf("MaxAttempts = %d, want 1", got.MaxAttempts)
	}
	if got.Delay != time.Second {
		t.Fatalf("Delay = %s, want %s", got.Delay, time.Second)
	}
}

func TestNormalizeRetryOptionsFillsUnsetFields(t *testing.T) {
	got := NormalizeRetryOptions(&RetryOptions{MaxAttempts: 2})
	if got.MaxAttempts != 2 {
		t.Fatalf("MaxAttempts = %d, want 2", got.MaxAttempts)
	}
	if got.Delay != DefaultRetryDelay {
		t.Fatalf("Delay = %s, want %s", got.Delay, DefaultRetryDelay)
	}
}
