package endfield

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
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
	data := []byte(`{"futureRootField":true,"owner":{"username":"test"},"playerInfo":{"businessCard":{"adventureLevel":0,"businessCardExpandFlag":false,"businessCardTopicId":0,"charList":null,"createTime":0,"gender":0,"mainMissionId":"","name":"test","platformRoleId":"","shortId":"","signature":"","userAvatarFrameId":0,"userAvatarId":0,"worldLevel":0},"charData":[],"futureField":true},"region":"EU/NA","ttl":42,"uid":"6105392891"}`)

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("failed to unmarshal profile: %v", err)
	}
	if profile.Raw == nil {
		t.Fatal("expected Raw to be set")
	}
	if _, ok := profile.Extra["futureRootField"]; !ok {
		t.Fatal("expected futureRootField in Extra")
	}
	if _, ok := profile.PlayerInfo.Extra["futureField"]; !ok {
		t.Fatal("expected futureField in PlayerInfo Extra")
	}

	got, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("failed to marshal profile: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestValidationDoesNotRequireNetwork(t *testing.T) {
	client := New(Options{UserAgent: "test-agent"})

	if _, err := client.GetProfile(context.Background(), ""); err != ErrInvalidUIDFormat {
		t.Fatalf("expected ErrInvalidUIDFormat, got %v", err)
	}

	if _, err := client.GetProfile(context.Background(), "invalid-uid-123"); err != ErrInvalidUIDFormat {
		t.Fatalf("expected ErrInvalidUIDFormat, got %v", err)
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
