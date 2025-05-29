//go:build integration
// +build integration

// export RUN_INTEGRATION_TESTS=true
// go test -v ./clients/enka -tags=integration

package enka

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kirinyoku/enkanetwork-go/internal/common"
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

	client := NewClient(nil, nil, "test-agent")
	_, err := client.GetUserProfile(context.Background(), "")
	if err != common.ErrInvalidUsername {
		t.Errorf("expected ErrInvalidUsername, got %v", err)
	}
}

// TestGetUserProfileNotFound ensures GetUserProfile returns ErrUserNotFound for a non-existent username.
func TestGetUserProfileNotFound(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	client := NewClient(nil, nil, "test-agent")
	_, err := client.GetUserProfile(context.Background(), "nonexistentuser12345")
	if err != common.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

// TestCompareJSONResponseGetUserProfile ensures that the JSON response from the API matches the JSON
// generated from the Go structure returned by the client GetUserProfile method.
func TestCompareJSONResponseGetUserProfile(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	username := "Algoinde"
	client := NewClient(nil, nil, "test-agent")

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

	apiJSON = common.RemoveTTLField(apiJSON)
	clientJSON = common.RemoveTTLField(clientJSON)

	// Compare JSON responses
	if !cmp.Equal(apiJSON, clientJSON) {
		t.Errorf("JSON responses do not match. API JSON: %s\nClient JSON: %s", apiJSON, clientJSON)
	}
}

// TestCompareJSONResponseGetUserProfileHoyos ensures that the JSON response from the API matches the JSON
// generated from the Go structure returned by the client GetUserProfileHoyos method.
func TestCompareJSONResponseGetUserProfileHoyos(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	ctx := context.Background()
	username := "Algoinde"
	client := NewClient(nil, nil, "test-agent")

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

	apiJSON = common.RemoveTTLField(apiJSON)
	clientJSON = common.RemoveTTLField(clientJSON)

	// Compare JSON responses
	if !cmp.Equal(apiJSON, clientJSON) {
		t.Errorf("JSON responses do not match. API JSON: %s\nClient JSON: %s", apiJSON, clientJSON)
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
	client := NewClient(nil, nil, "test-agent")

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

	apiJSON = common.RemoveTTLField(apiJSON)
	clientJSON = common.RemoveTTLField(clientJSON)

	if !cmp.Equal(apiJSON, clientJSON) {
		t.Errorf("JSON responses do not match. API JSON: %s\nClient JSON: %s", apiJSON, clientJSON)
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
	client := NewClient(nil, nil, "test-agent")

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

	apiJSONBytes = common.RemoveTTLField(apiJSONBytes)

	var apiData Builds
	err = json.Unmarshal(apiJSONBytes, &apiData)
	if err != nil {
		t.Fatalf("failed to unmarshal API JSON into struct: %v", err)
	}

	clientJSONBytes, err := json.Marshal(builds)
	if err != nil {
		t.Fatalf("failed to marshal client response to JSON: %v", err)
	}

	clientJSONBytes = common.RemoveTTLField(clientJSONBytes)

	var clientData Builds
	err = json.Unmarshal(clientJSONBytes, &clientData)
	if err != nil {
		t.Fatalf("failed to unmarshal client marshaled JSON into struct: %v", err)
	}

	if !cmp.Equal(apiData, clientData) {
		t.Errorf("JSON responses do not match. API JSON: %s\nClient JSON: %s", apiJSONBytes, clientJSONBytes)
	}
}
