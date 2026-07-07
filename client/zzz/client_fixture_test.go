package zzz

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
	data := []byte(`{"PlayerInfo":{"SocialDetail":{"MedalList":null,"ProfileDetail":null,"Desc":"fixture"},"ShowcaseDetail":null},"ttl":300,"uid":"1301806568","futureField":{"enabled":true}}`)

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
	data := []byte(`{"PlayerInfo":{"SocialDetail":{"MedalList":[],"ProfileDetail":{"Uid":1,"Nickname":"Proxy","ProfileId":1,"Level":60,"Title":1,"CallingCardId":1,"AvatarId":1,"TitleInfo":{"Title":1,"FullTitle":1,"Args":[],"KMOHDEAKEFG":[],"futureTitle":true},"PlatformType":1,"futureProfile":true},"Desc":"fixture","futureSocial":true},"ShowcaseDetail":{"AvatarList":[{"Id":1021,"Exp":0,"Level":60,"PromotionLevel":6,"TalentLevel":0,"SkinId":0,"CoreSkillEnhancement":0,"TalentToggleList":[],"WeaponEffectState":0,"ClaimedRewardList":[],"ObtainmentTimestamp":0,"Weapon":{"Uid":1,"Id":1,"Exp":0,"Level":60,"BreakLevel":5,"UpgradeLevel":1,"futureWeapon":true},"SkillLevelList":[],"EquippedList":[{"Slot":1,"Equipment":{"Uid":1,"Id":1,"Exp":0,"Level":15,"BreakLevel":5,"MainPropertyList":[],"RandomPropertyList":[],"futureEquipment":true}}],"WeaponUid":1,"UpgradeId":0,"futureAvatar":true}],"futureShowcase":true},"futurePlayer":true},"ttl":300}`)

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatalf("failed to unmarshal profile: %v", err)
	}
	if _, ok := profile.PlayerInfo.Extra["futurePlayer"]; !ok {
		t.Fatal("expected playerInfo futurePlayer in Extra")
	}
	if _, ok := profile.PlayerInfo.SocialDetail.Extra["futureSocial"]; !ok {
		t.Fatal("expected socialDetail futureSocial in Extra")
	}
	if _, ok := profile.PlayerInfo.SocialDetail.ProfileDetail.Extra["futureProfile"]; !ok {
		t.Fatal("expected profileDetail futureProfile in Extra")
	}
	if _, ok := profile.PlayerInfo.SocialDetail.ProfileDetail.TitleInfo.Extra["futureTitle"]; !ok {
		t.Fatal("expected titleInfo futureTitle in Extra")
	}
	if _, ok := profile.PlayerInfo.ShowcaseDetail.Extra["futureShowcase"]; !ok {
		t.Fatal("expected showcaseDetail futureShowcase in Extra")
	}
	avatar := profile.PlayerInfo.ShowcaseDetail.AvatarList[0]
	if _, ok := avatar.Extra["futureAvatar"]; !ok {
		t.Fatal("expected avatar futureAvatar in Extra")
	}
	if _, ok := avatar.Weapon.Extra["futureWeapon"]; !ok {
		t.Fatal("expected weapon futureWeapon in Extra")
	}
	if _, ok := avatar.EquippedList[0].Equipment.Extra["futureEquipment"]; !ok {
		t.Fatal("expected equipment futureEquipment in Extra")
	}

	got, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("failed to marshal profile: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestOptionalBoolPresenceIsPreserved(t *testing.T) {
	data := []byte(`{"Uid":7392,"Id":33041,"Exp":0,"Level":15,"BreakLevel":5,"IsAvailable":false,"IsLocked":false,"IsTrash":false,"MainPropertyList":[],"RandomPropertyList":[]}`)

	var equipment Equipment
	if err := json.Unmarshal(data, &equipment); err != nil {
		t.Fatalf("failed to unmarshal equipment: %v", err)
	}

	got, err := json.Marshal(equipment)
	if err != nil {
		t.Fatalf("failed to marshal equipment: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func TestOptionalBoolAbsenceIsPreserved(t *testing.T) {
	data := []byte(`{"Uid":7392,"Id":33041,"Exp":0,"Level":15,"BreakLevel":5,"MainPropertyList":[],"RandomPropertyList":[]}`)

	var equipment Equipment
	if err := json.Unmarshal(data, &equipment); err != nil {
		t.Fatalf("failed to unmarshal equipment: %v", err)
	}

	got, err := json.Marshal(equipment)
	if err != nil {
		t.Fatalf("failed to marshal equipment: %v", err)
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
