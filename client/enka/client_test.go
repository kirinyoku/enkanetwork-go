//go:build integration
// +build integration

// export RUN_INTEGRATION_TESTS=true
// go test -v ./client/enka -tags=integration

package enka

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
)

// TestMain sets up any global state for the integration tests.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestGetUserProfileInvalidUsername checks that GetUserProfile returns ErrInvalidUsername for an empty username.
func TestGetUserProfileInvalidUsername(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	client := New(Options{UserAgent: "test-agent"})
	_, err := client.GetUserProfile(context.Background(), "")
	if err != ErrInvalidUsername {
		t.Errorf("expected ErrInvalidUsername, got %v", err)
	}
}

// TestGetUserProfileNotFound ensures GetUserProfile returns ErrUserNotFound for a non-existent username.
func TestGetUserProfileNotFound(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	client := New(Options{UserAgent: "test-agent"})
	_, err := client.GetUserProfile(context.Background(), "nonexistentuser12345")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

// TestGetUserProfile ensures that the JSON response from the API matches the JSON
// generated from the Go structure returned by the client GetUserProfile method.
func TestGetUserProfile(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	username := "Algoinde"
	client := New(Options{UserAgent: "test-agent"})

	profile, err := client.GetUserProfile(ctx, username)
	if err != nil {
		t.Fatalf("failed to get profile from client: %v", err)
	}

	clientJSON, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("failed to marshal client response to JSON: %v", err)
	}

	url := fmt.Sprintf("https://enka.network/api/profile/%s/", username)
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

	if diff := core.JSONDiff(apiJSON, clientJSON); diff != "" {
		t.Errorf("JSON responses do not match: %s", diff)
	}
}

// TestGetUserProfileHoyos ensures that the JSON response from the API matches the JSON
// generated from the Go structure returned by the client GetUserProfileHoyos method.
func TestGetUserProfileHoyos(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	username := "Algoinde"
	client := New(Options{UserAgent: "test-agent"})

	hoyos, err := client.GetUserProfileHoyos(ctx, username)
	if err != nil {
		t.Fatalf("failed to get profile from client: %v", err)
	}

	clientJSON, err := json.Marshal(hoyos)
	if err != nil {
		t.Fatalf("failed to marshal client response to JSON: %v", err)
	}

	url := fmt.Sprintf("https://enka.network/api/profile/%s/hoyos/", username)
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

	if diff := core.JSONDiff(apiJSON, clientJSON); diff != "" {
		t.Errorf("JSON responses do not match: %s", diff)
	}
}

// TestGetUserProfileHoyo ensures that the JSON response from the API matches the JSON
// generated from the Go structure returned by the client GetUserProfileHoyo method.
func TestGetUserProfileHoyo(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	username := "Algoinde"
	client := New(Options{UserAgent: "test-agent"})

	hoyo, err := client.GetUserProfileHoyo(ctx, username, "4Wjv2e")
	if err != nil {
		t.Fatalf("failed to get profile from client: %v", err)
	}

	clientJSON, err := json.Marshal(hoyo)
	if err != nil {
		t.Fatalf("failed to marshal client response to JSON: %v", err)
	}

	url := fmt.Sprintf("https://enka.network/api/profile/%s/hoyos/%s/?format=json", username, "4Wjv2e")
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

	if diff := core.JSONDiff(apiJSON, clientJSON); diff != "" {
		t.Errorf("JSON responses do not match: %s", diff)
	}
}

// TestGetUserProfileHoyoBuilds ensures that the JSON response from the API matches the JSON
// generated from the Go structure returned by the client GetUserProfileHoyoBuilds method.
func TestGetUserProfileHoyoBuilds(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	username := "Algoinde"
	client := New(Options{UserAgent: "test-agent"})

	builds, err := client.GetUserProfileHoyoBuilds(ctx, username, "4Wjv2e")
	if err != nil {
		t.Fatalf("failed to get profile from client: %v", err)
	}

	url := fmt.Sprintf("https://enka.network/api/profile/%s/hoyos/%s/builds/", username, "4Wjv2e")
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

	apiJSONBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read API response: %v", err)
	}

	apiJSONBytes = core.RemoveTTLField(apiJSONBytes)

	clientJSONBytes, err := json.Marshal(builds)
	if err != nil {
		t.Fatalf("failed to marshal client response to JSON: %v", err)
	}

	clientJSONBytes = core.RemoveTTLField(clientJSONBytes)

	if diff := core.JSONDiff(apiJSONBytes, clientJSONBytes); diff != "" {
		t.Errorf("JSON responses do not match: %s", diff)
	}
}
