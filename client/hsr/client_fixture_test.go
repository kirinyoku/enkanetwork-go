package hsr

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
	data := []byte(`{"detailInfo":{"nickname":"Maple"},"ttl":300,"uid":"807752192","futureField":{"enabled":true}}`)

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
	data := []byte(`{"detailInfo":{"nickname":"Maple","worldLevel":0,"futureDetail":true,"avatarDetailList":[{"avatarId":1301,"rank":0,"futureAvatar":{"enabled":true},"relicList":[{"tid":61011,"type":0,"futureRelic":1,"_flat":{"setID":0,"futureFlat":"yes"}}],"equipment":{"tid":23001,"level":0,"futureEquipment":false,"_flat":{"name":"123","futureEquipmentFlat":true}}}],"recordInfo":{"achievementCount":0,"futureRecord":42}},"ttl":300}`)

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("failed to unmarshal profile: %v", err)
	}
	if _, ok := profile.DetailInfo.Extra["futureDetail"]; !ok {
		t.Fatal("expected detailInfo futureDetail in Extra")
	}
	avatar := profile.DetailInfo.AvatarDetailList[0]
	if _, ok := avatar.Extra["futureAvatar"]; !ok {
		t.Fatal("expected avatar futureAvatar in Extra")
	}
	if _, ok := avatar.RelicList[0].Extra["futureRelic"]; !ok {
		t.Fatal("expected relic futureRelic in Extra")
	}
	if _, ok := avatar.RelicList[0].Flat.Extra["futureFlat"]; !ok {
		t.Fatal("expected flat futureFlat in Extra")
	}
	if _, ok := avatar.Equipment.Extra["futureEquipment"]; !ok {
		t.Fatal("expected equipment futureEquipment in Extra")
	}
	if _, ok := avatar.Equipment.Flat.Extra["futureEquipmentFlat"]; !ok {
		t.Fatal("expected equipment flat futureEquipmentFlat in Extra")
	}
	if _, ok := profile.DetailInfo.RecordInfo.Extra["futureRecord"]; !ok {
		t.Fatal("expected recordInfo futureRecord in Extra")
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
