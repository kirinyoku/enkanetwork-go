//go:build integration
// +build integration

// export RUN_INTEGRATION_TESTS=true
// go test -v ./client/zzz -tags=integration

package zzz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kirinyoku/enkanetwork-go/internal/core"
)

// TestMain sets up any global state for the integration tests.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestGetProfileInvalidUID checks that GetProfile returns ErrInvalidUIDFormat for an invalid UID.
func TestGetProfileInvalidUID(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	client := NewClient(nil, nil, "test-agent")
	_, err := client.GetProfile(context.Background(), "123")
	if err != ErrInvalidUIDFormat {
		t.Errorf("expected ErrInvalidUIDFormat, got %v", err)
	}
}

// TestGetProfileNotFound ensures GetProfile returns ErrPlayerNotFound for a non-existent UID.
func TestGetProfileNotFound(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	client := NewClient(nil, nil, "test-agent")
	_, err := client.GetProfile(context.Background(), "987654321")
	if err != ErrPlayerNotFound {
		t.Errorf("expected ErrPlayerNotFound, got %v", err)
	}
}

// TestGetProfile ensures that the JSON response from the API matches the JSON
// generated from the Go structure returned by the client GetProfile method.
func TestGetProfile(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	uid := "1301806568"
	client := NewClient(nil, nil, "test-agent")

	profile, err := client.GetProfile(ctx, uid)
	if err != nil {
		t.Fatalf("failed to get profile from client: %v", err)
	}

	clientJSON, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("failed to marshal client response to JSON: %v", err)
	}

	url := fmt.Sprintf("https://enka.network/api/zzz/uid/%s", uid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("User-Agent", "test-agent")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	apiJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read API response: %v", err)
	}

	apiJSON = core.RemoveTTLField(apiJSON)
	clientJSON = core.RemoveTTLField(clientJSON)

	if !cmp.Equal(apiJSON, clientJSON) {
		t.Errorf("JSON responses do not match. API JSON: %s\nClient JSON: %s", apiJSON, clientJSON)
	}
}
