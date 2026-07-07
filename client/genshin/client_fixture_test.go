package genshin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestProfileFixtureRoundTrip(t *testing.T) {
	data := readFixture(t, "testdata/profile.json")

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("failed to unmarshal fixture: %v", err)
	}

	got, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("failed to marshal profile: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestProfilePreservesUnknownFields(t *testing.T) {
	data := []byte(`{"playerInfo":{"nickname":"Algoinde"},"ttl":300,"uid":"618285856","futureField":{"enabled":true}}`)

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("failed to unmarshal profile: %v", err)
	}
	if profile.Raw == nil {
		t.Fatal("expected Raw to be set")
	}
	if _, ok := profile.Extra["futureField"]; !ok {
		t.Fatal("expected futureField in Extra")
	}

	got, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("failed to marshal profile: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestProfilePreservesNestedUnknownFields(t *testing.T) {
	data := []byte(`{"playerInfo":{"nickname":"Algoinde","nameCardId":0,"futurePlayer":true},"avatarInfoList":[{"avatarId":10000002,"skillDepotId":0,"futureAvatar":{"enabled":true},"equipList":[{"itemId":11501,"futureEquip":1,"weapon":{"level":0,"futureWeapon":"yes"},"reliquary":{"level":0,"mainPropId":0,"futureReliquary":false}}]}],"ttl":300}`)

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("failed to unmarshal profile: %v", err)
	}
	if _, ok := profile.PlayerInfo.Extra["futurePlayer"]; !ok {
		t.Fatal("expected playerInfo futurePlayer in Extra")
	}
	if _, ok := profile.AvatarInfoList[0].Extra["futureAvatar"]; !ok {
		t.Fatal("expected avatar futureAvatar in Extra")
	}
	if _, ok := profile.AvatarInfoList[0].EquipList[0].Extra["futureEquip"]; !ok {
		t.Fatal("expected equip futureEquip in Extra")
	}
	if _, ok := profile.AvatarInfoList[0].EquipList[0].Weapon.Extra["futureWeapon"]; !ok {
		t.Fatal("expected weapon futureWeapon in Extra")
	}
	if _, ok := profile.AvatarInfoList[0].EquipList[0].Reliquary.Extra["futureReliquary"]; !ok {
		t.Fatal("expected reliquary futureReliquary in Extra")
	}

	got, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("failed to marshal profile: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestValidationDoesNotRequireNetwork(t *testing.T) {
	client := New(Options{UserAgent: "test-agent"})

	if _, err := client.GetProfile(context.Background(), "123"); err != ErrInvalidUIDFormat {
		t.Fatalf("GetProfile expected ErrInvalidUIDFormat, got %v", err)
	}

	if _, err := client.GetPlayerInfo(context.Background(), "123"); err != ErrInvalidUIDFormat {
		t.Fatalf("GetPlayerInfo expected ErrInvalidUIDFormat, got %v", err)
	}
}

func TestOptionsBaseURLUsesCustomServer(t *testing.T) {
	requestedPath := ""
	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			requestedPath = r.URL.RequestURI()
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"playerInfo":{"nickname":"Traveler"},"ttl":60}`)),
				Request:    r,
			}, nil
		}),
	}

	client := New(Options{
		HTTPClient: httpClient,
		BaseURL:    "https://example.test/api",
		UserAgent:  "test-agent",
	})

	profile, err := client.GetProfile(context.Background(), "618285856")
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if profile.PlayerInfo.Nickname != "Traveler" {
		t.Fatalf("nickname = %q, want Traveler", profile.PlayerInfo.Nickname)
	}
	if requestedPath != "/api/uid/618285856" {
		t.Fatalf("requested path = %q, want /api/uid/618285856", requestedPath)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func readFixture(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", path, err)
	}
	return data
}

func assertJSONEqual(t *testing.T, wantJSON, gotJSON []byte) {
	t.Helper()

	var want, got any
	if err := json.Unmarshal(wantJSON, &want); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(gotJSON, &got); err != nil {
		t.Fatalf("failed to unmarshal actual JSON: %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("JSON mismatch\nwant: %s\ngot:  %s", wantJSON, gotJSON)
	}
}
