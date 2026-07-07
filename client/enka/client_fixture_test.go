package enka

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestOwnerFixtureRoundTrip(t *testing.T) {
	data := readFixture(t, "testdata/profile.json")

	var owner Owner
	if err := json.Unmarshal(data, &owner); err != nil {
		t.Fatalf("failed to unmarshal fixture: %v", err)
	}

	got, err := json.Marshal(owner)
	if err != nil {
		t.Fatalf("failed to marshal owner: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestOwnerPreservesUnknownFields(t *testing.T) {
	data := []byte(`{"id":749,"username":"Algoinde","futureField":{"enabled":true}}`)

	var owner Owner
	if err := json.Unmarshal(data, &owner); err != nil {
		t.Fatalf("failed to unmarshal owner: %v", err)
	}
	if owner.Raw == nil {
		t.Fatal("expected Raw to be set")
	}
	if _, ok := owner.Extra["futureField"]; !ok {
		t.Fatal("expected futureField in Extra")
	}

	got, err := json.Marshal(owner)
	if err != nil {
		t.Fatalf("failed to marshal owner: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestHoyosFixtureRoundTrip(t *testing.T) {
	data := readFixture(t, "testdata/hoyos.json")

	var hoyos Hoyos
	if err := json.Unmarshal(data, &hoyos); err != nil {
		t.Fatalf("failed to unmarshal fixture: %v", err)
	}

	got, err := json.Marshal(hoyos)
	if err != nil {
		t.Fatalf("failed to marshal hoyos: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestHoyoPreservesUnknownFields(t *testing.T) {
	data := []byte(`{"avatar_order":null,"hash":"2hOhJd","hoyo_type":3,"live_public":false,"order":"0.0000000000","uid_public":false,"futureField":{"enabled":true}}`)

	var hoyo Hoyo
	if err := json.Unmarshal(data, &hoyo); err != nil {
		t.Fatalf("failed to unmarshal hoyo: %v", err)
	}
	if hoyo.UID != nil {
		t.Fatalf("expected missing uid to stay nil, got %d", *hoyo.UID)
	}
	if hoyo.HoyoType != HoyoTypeArknightsEndfield {
		t.Fatalf("hoyo type = %d, want %d", hoyo.HoyoType, HoyoTypeArknightsEndfield)
	}
	if hoyo.Raw == nil {
		t.Fatal("expected Raw to be set")
	}
	if _, ok := hoyo.Extra["futureField"]; !ok {
		t.Fatal("expected futureField in Extra")
	}

	got, err := json.Marshal(hoyo)
	if err != nil {
		t.Fatalf("failed to marshal hoyo: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestHoyoOptionalBoolAbsenceIsPreserved(t *testing.T) {
	data := []byte(`{"avatar_order":null,"hash":"2hOhJd","hoyo_type":3,"order":"0.0000000000"}`)

	var hoyo Hoyo
	if err := json.Unmarshal(data, &hoyo); err != nil {
		t.Fatalf("failed to unmarshal hoyo: %v", err)
	}

	got, err := json.Marshal(hoyo)
	if err != nil {
		t.Fatalf("failed to marshal hoyo: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestBuildAvatarDataUsesHoyoTypeForHSR(t *testing.T) {
	data := []byte(`{"id":1,"name":"HSR build","avatar_id":"1301","avatar_data":{"avatarId":1301,"level":80,"rank":6,"pos":1},"settings":{},"image":null,"order":"0.0000000000","hoyo_type":1}`)

	var build Build
	if err := json.Unmarshal(data, &build); err != nil {
		t.Fatalf("failed to unmarshal build: %v", err)
	}
	if build.HoyoType != HoyoTypeHSR {
		t.Fatalf("hoyo type = %d, want %d", build.HoyoType, HoyoTypeHSR)
	}
	if build.AvatarData.Genshin != nil {
		t.Fatal("expected Genshin avatar data to stay nil")
	}
	if build.AvatarData.HSR == nil {
		t.Fatal("expected HSR avatar data to be populated")
	}
	if build.AvatarData.ZZZ != nil {
		t.Fatal("expected ZZZ avatar data to stay nil")
	}
	if build.AvatarData.HSR.AvatarID != 1301 || build.AvatarData.HSR.Rank != 6 {
		t.Fatalf("unexpected HSR avatar data: %+v", build.AvatarData.HSR)
	}

	got, err := json.Marshal(build)
	if err != nil {
		t.Fatalf("failed to marshal build: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestBuildAvatarDataUsesHoyoTypeForZZZ(t *testing.T) {
	data := []byte(`{"id":2,"name":"ZZZ build","avatar_id":"1021","avatar_data":{"Id":1021,"Level":60,"TalentLevel":3,"WeaponUid":42,"SkillLevelList":[]},"settings":{},"image":null,"order":"0.0000000000","hoyo_type":2}`)

	var build Build
	if err := json.Unmarshal(data, &build); err != nil {
		t.Fatalf("failed to unmarshal build: %v", err)
	}
	if build.HoyoType != HoyoTypeZZZ {
		t.Fatalf("hoyo type = %d, want %d", build.HoyoType, HoyoTypeZZZ)
	}
	if build.AvatarData.Genshin != nil {
		t.Fatal("expected Genshin avatar data to stay nil")
	}
	if build.AvatarData.HSR != nil {
		t.Fatal("expected HSR avatar data to stay nil")
	}
	if build.AvatarData.ZZZ == nil {
		t.Fatal("expected ZZZ avatar data to be populated")
	}
	if build.AvatarData.ZZZ.ID != 1021 || build.AvatarData.ZZZ.TalentLevel != 3 {
		t.Fatalf("unexpected ZZZ avatar data: %+v", build.AvatarData.ZZZ)
	}

	got, err := json.Marshal(build)
	if err != nil {
		t.Fatalf("failed to marshal build: %v", err)
	}

	var gotFields map[string]json.RawMessage
	if err := json.Unmarshal(got, &gotFields); err != nil {
		t.Fatalf("failed to unmarshal marshaled build: %v", err)
	}
	var avatarData map[string]any
	if err := json.Unmarshal(gotFields["avatar_data"], &avatarData); err != nil {
		t.Fatalf("failed to unmarshal marshaled avatar_data: %v", err)
	}
	if _, ok := avatarData["zzz"]; ok {
		t.Fatalf("avatar_data contains nested zzz wrapper: %s", gotFields["avatar_data"])
	}
	if avatarData["Id"] != float64(1021) {
		t.Fatalf("avatar_data Id = %v, want 1021", avatarData["Id"])
	}
}

func TestBuildAvatarDataPreservesRawForUnknownHoyoType(t *testing.T) {
	data := []byte(`{"id":3,"name":"future build","avatar_id":"9001","avatar_data":{"future":true,"value":9001},"settings":{},"image":null,"order":"0.0000000000","hoyo_type":9}`)

	var build Build
	if err := json.Unmarshal(data, &build); err != nil {
		t.Fatalf("failed to unmarshal build: %v", err)
	}
	if build.AvatarData.Genshin != nil || build.AvatarData.HSR != nil || build.AvatarData.ZZZ != nil {
		t.Fatalf("expected only raw avatar data, got %+v", build.AvatarData)
	}
	if build.AvatarData.Raw == nil {
		t.Fatal("expected raw avatar data to be preserved")
	}

	got, err := json.Marshal(build)
	if err != nil {
		t.Fatalf("failed to marshal build: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestEnkaPathParametersAreEscaped(t *testing.T) {
	var requestedPath string
	httpClient := &http.Client{
		Transport: enkaRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			requestedPath = r.URL.RequestURI()
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"hash":"hash","order":"0","avatar_order":{},"hoyo_type":0}`)),
				Request:    r,
			}, nil
		}),
	}
	client := New(Options{
		HTTPClient: httpClient,
		BaseURL:    "https://example.test/api",
		UserAgent:  "test-agent",
	})

	_, err := client.GetUserProfileHoyo(context.Background(), "name/with?query", "hash/with#fragment")
	if err != nil {
		t.Fatalf("GetUserProfileHoyo() error = %v", err)
	}

	want := "/api/profile/name%2Fwith%3Fquery/hoyos/hash%2Fwith%23fragment"
	if requestedPath != want {
		t.Fatalf("requested path = %q, want %q", requestedPath, want)
	}
}

func TestGetUserProfileHoyosUsesCachedValue(t *testing.T) {
	requests := 0
	httpClient := &http.Client{
		Transport: enkaRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			requests++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"abc123":{"hash":"abc123","order":"0","avatar_order":{},"hoyo_type":0}}`)),
				Request:    r,
			}, nil
		}),
	}

	client := New(Options{
		HTTPClient: httpClient,
		Cache:      newTestCache(),
		BaseURL:    "https://example.test/api",
		UserAgent:  "test-agent",
	})

	first, err := client.GetUserProfileHoyos(context.Background(), "Algoinde")
	if err != nil {
		t.Fatalf("first GetUserProfileHoyos() error = %v", err)
	}
	second, err := client.GetUserProfileHoyos(context.Background(), "Algoinde")
	if err != nil {
		t.Fatalf("second GetUserProfileHoyos() error = %v", err)
	}

	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if first["abc123"].Hash != "abc123" || second["abc123"].Hash != "abc123" {
		t.Fatalf("unexpected cached hoyos: first=%+v second=%+v", first, second)
	}
}

func TestGetUserProfileHoyoBuildsUsesCachedValue(t *testing.T) {
	requests := 0
	httpClient := &http.Client{
		Transport: enkaRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			requests++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"10000002":[{"id":1,"name":"build","avatar_id":"10000002","settings":{},"image":null,"order":"0","hoyo_type":0}]}`)),
				Request:    r,
			}, nil
		}),
	}

	client := New(Options{
		HTTPClient: httpClient,
		Cache:      newTestCache(),
		BaseURL:    "https://example.test/api",
		UserAgent:  "test-agent",
	})

	first, err := client.GetUserProfileHoyoBuilds(context.Background(), "Algoinde", "abc123")
	if err != nil {
		t.Fatalf("first GetUserProfileHoyoBuilds() error = %v", err)
	}
	second, err := client.GetUserProfileHoyoBuilds(context.Background(), "Algoinde", "abc123")
	if err != nil {
		t.Fatalf("second GetUserProfileHoyoBuilds() error = %v", err)
	}

	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
	if first["10000002"][0].Name != "build" || second["10000002"][0].Name != "build" {
		t.Fatalf("unexpected cached builds: first=%+v second=%+v", first, second)
	}
}

func TestValidationDoesNotRequireNetwork(t *testing.T) {
	client := New(Options{UserAgent: "test-agent"})

	if _, err := client.GetUserProfile(context.Background(), ""); err != ErrInvalidUsername {
		t.Fatalf("GetUserProfile expected ErrInvalidUsername, got %v", err)
	}

	if _, err := client.GetUserProfileHoyo(context.Background(), "", "hash"); err != ErrInvalidUsername {
		t.Fatalf("GetUserProfileHoyo expected ErrInvalidUsername, got %v", err)
	}

	if _, err := client.GetUserProfileHoyo(context.Background(), "Algoinde", ""); err != ErrInvalidHoyoHash {
		t.Fatalf("GetUserProfileHoyo expected ErrInvalidHoyoHash, got %v", err)
	}

	if _, err := client.GetUserProfileHoyoBuilds(context.Background(), "", "hash"); err != ErrInvalidUsername {
		t.Fatalf("GetUserProfileHoyoBuilds expected ErrInvalidUsername, got %v", err)
	}

	if _, err := client.GetUserProfileHoyoBuilds(context.Background(), "Algoinde", ""); err != ErrInvalidHoyoHash {
		t.Fatalf("GetUserProfileHoyoBuilds expected ErrInvalidHoyoHash, got %v", err)
	}
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

type enkaRoundTripFunc func(*http.Request) (*http.Response, error)

func (f enkaRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type testCache struct {
	values map[string]any
}

func newTestCache() *testCache {
	return &testCache{values: make(map[string]any)}
}

func (c *testCache) Get(key string) (any, bool) {
	value, ok := c.values[key]
	return value, ok
}

func (c *testCache) Set(key string, value any, _ time.Duration) {
	c.values[key] = value
}
