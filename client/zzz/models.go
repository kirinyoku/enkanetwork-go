package zzz

import (
	"encoding/json"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/models"
)

// --------------------------------------------IMPORTANT--------------------------------------------
// For detailed information on properties, refer to the EnkaNetwork API — Zenless
// Zone Zero documentation (https://github.com/kirinyoku/enkanetwork-go/blob/master/docs/zzz/api.md).
// -------------------------------------------------------------------------------------------------
//
// Profile
// ├── PlayerInfo
// │   ├── SocialDetail
// │   │   ├── ProfileDetail
// │   │   │   └── TitleInfo
// │   │   └── MedalList
// │   └── ShowcaseDetail
// │       └── AvatarList
// │           ├── SkillLevelList
// │           ├── Weapon
// │           └── EquippedList
// │               └── Equipment
// │                   ├── MainPropertyList
// │                   └── RandomPropertyList
// ├── ttl
// ├── uid
// ├── region
// └── owner

// Profile represents the root structure of the response containing player information
// for Zenless Zone Zero.
type Profile struct {
	// PlayerInfo contains basic information about the game account from the player's showcase.
	PlayerInfo PlayerInfo `json:"PlayerInfo"`
	// TTL indicates the seconds remaining until the next request to the game. Until
	// the TTL expires, the endpoint returns cached data — but such requests still
	// count toward the rate limit.
	TTL int `json:"ttl"`
	// Owner is the Enka profile associated with the provided UID. The response includes
	// an Owner if:
	//   1. The user has an account on the site;
	//   2. The user has added their UID to their profile;
	//   3. The user has verified that the UID belongs to them;
	//   4. The user has set their profile visibility to "public".
	Owner *models.Owner `json:"owner,omitempty"`
	// UID is the player's UID in Zenless Zone Zero.
	UID string `json:"uid,omitempty"`
	// Region is the player's server region (e.g., "Asia", "Europe", "America").
	Region string `json:"region,omitempty"`
	// Raw contains the original API response.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown top-level API fields preserved during unmarshaling.
	Extra map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON preserves unknown top-level fields for API drift tolerance.
func (p *Profile) UnmarshalJSON(data []byte) error {
	type Alias Profile
	var profile Alias
	if err := json.Unmarshal(data, &profile); err != nil {
		return err
	}

	raw, extra, err := core.PreserveUnknownJSON(data, "PlayerInfo", "ttl", "owner", "uid", "region")
	if err != nil {
		return err
	}

	profile.Raw = raw
	profile.Extra = extra
	*p = Profile(profile)
	return nil
}

// MarshalJSON writes known fields plus unknown fields preserved in Extra.
func (p Profile) MarshalJSON() ([]byte, error) {
	type Alias Profile
	return core.MergeKnownExtraAndRawJSON(Alias(p), p.Extra, p.Raw)
}

// CacheTTL implements the core.Cacheable interface.
func (p Profile) CacheTTL() time.Duration {
	return time.Duration(p.TTL) * time.Second
}

// Build contains information about a specific agent's build in Zenless Zone Zero.
type Build struct {
	// ID is the ID of the build.
	ID int `json:"id"`
	// Name is the name of the build.
	Name string `json:"name"`
	// AvatarID is the ID of the agent.
	AvatarID string `json:"avatar_id"`
	// AvatarData is the agent data (*zzz.AvatarData).
	AvatarData *AvatarData `json:"avatar_data"`
	// Order is the order of the saved build on the Enka.
	Order int `json:"order"`
	// Live indicates it is not a saved build but one retrieved from the game's showcase when the "refresh" button is clicked. During
	// an update, all old live builds are deleted, and new ones are created. Updates
	// are user-initiated, so this data may not be up to date.
	Live *bool `json:"live,omitempty"`
	// Settings contains build-specific configuration data.
	Settings *Settings `json:"settings"`
	// Public indicates whether the build is public.
	Public *bool `json:"public,omitempty"`
	// Image is the URL of the build image.
	Image *string `json:"image"`
	// HoyoType is the game type discriminator (2 = Zenless Zone Zero).
	HoyoType int `json:"hoyo_type"`
	// Hoyo is the unique hoyo identifier (hoyo_hash).
	Hoyo string `json:"hoyo"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (b *Build) UnmarshalJSON(data []byte) error {
	type Alias Build
	var build Alias
	if err := json.Unmarshal(data, &build); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	build.Raw = raw
	build.Extra = extra
	*b = Build(build)
	return nil
}

func (b Build) MarshalJSON() ([]byte, error) {
	type Alias Build
	return core.MergeKnownExtraAndRawJSON(Alias(b), b.Extra, b.Raw)
}

// PlayerInfo contains basic information about a player's game account.
type PlayerInfo struct {
	// SocialDetail contains public profile information, including the player's signature,
	// level, selected namecard, and displayed badges.
	SocialDetail *SocialDetail `json:"SocialDetail"`
	// ShowcaseDetail contains detailed information about the agents the player has chosen to
	// display on their in-game profile showcase. Includes their stats,
	// skills, and equipment.
	ShowcaseDetail *ShowcaseDetail `json:"ShowcaseDetail"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (p *PlayerInfo) UnmarshalJSON(data []byte) error {
	type Alias PlayerInfo
	var info Alias
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	info.Raw = raw
	info.Extra = extra
	*p = PlayerInfo(info)
	return nil
}

func (p PlayerInfo) MarshalJSON() ([]byte, error) {
	type Alias PlayerInfo
	return core.MergeKnownExtraAndRawJSON(Alias(p), p.Extra, p.Raw)
}

// ShowcaseDetail contains a list of agents in the player’s showcase.
type ShowcaseDetail struct {
	// AvatarList is a list of detailed agent configurations (stats, weapons, discs) currently visible on the player's profile showcase.
	AvatarList []AvatarData `json:"AvatarList"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (s *ShowcaseDetail) UnmarshalJSON(data []byte) error {
	type Alias ShowcaseDetail
	var showcase Alias
	if err := json.Unmarshal(data, &showcase); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	showcase.Raw = raw
	showcase.Extra = extra
	*s = ShowcaseDetail(showcase)
	return nil
}

func (s ShowcaseDetail) MarshalJSON() ([]byte, error) {
	type Alias ShowcaseDetail
	return core.MergeKnownExtraAndRawJSON(Alias(s), s.Extra, s.Raw)
}

// AvatarData contains detailed information about an agent.
type AvatarData struct {
	// ID is the unique Agent ID (see https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/avatars.json).
	ID int `json:"Id"`
	// Exp is the current experience points of the agent towards their next level.
	Exp int `json:"Exp"`
	// Level is the agent's current overall level.
	Level int `json:"Level"`
	// PromotionLevel is the agent's ascension phase, which dictates their maximum level cap.
	PromotionLevel int `json:"PromotionLevel"`
	// TalentLevel is the number of unlocked Mindscape Cinema nodes (duplicates obtained) for this agent.
	TalentLevel int `json:"TalentLevel"`
	// SkinID is the ID of the equipped outfit/skin for the agent.
	SkinID int `json:"SkinId"`
	// CoreSkillEnhancement is the level of the agent's Core Passive skill enhancement (nodes A through F represented as an integer).
	CoreSkillEnhancement int `json:"CoreSkillEnhancement"`
	// TalentToggleList represents visual toggles for Mindscape Cinema effects.
	TalentToggleList []bool `json:"TalentToggleList"`
	// WeaponEffectState indicates if the equipped W-Engine's special visual effect is active [0: None, 1: OFF, 2: ON].
	WeaponEffectState int `json:"WeaponEffectState"`
	// IsHidden indicates whether the agent's detailed stats are hidden from public view in the showcase.
	IsHidden *bool `json:"IsHidden,omitempty"`
	// ClaimedRewardList contains IDs representing the ascension rewards the player has claimed for this agent.
	ClaimedRewardList []int `json:"ClaimedRewardList"`
	// ObtainmentTimestamp is the Unix timestamp of when the player first obtained this agent.
	ObtainmentTimestamp int64 `json:"ObtainmentTimestamp"`
	// Weapon contains detailed data about the equipped W-Engine.
	Weapon *Weapon `json:"Weapon"`
	// SkillLevelList is an array detailing the investment levels for the agent's combat skills.
	SkillLevelList []SkillLevel `json:"SkillLevelList"`
	// EquippedList is an array of Drive Discs currently equipped on the agent.
	EquippedList []EquippedItem `json:"EquippedList"`
	// IsFavorite indicates whether the player has marked this agent as a favorite.
	IsFavorite *bool `json:"IsFavorite,omitempty"`
	// WeaponUID is the unique instance ID of the W-Engine equipped on this agent.
	WeaponUID int `json:"WeaponUid"`
	// IsUpgradeUnlocked indicates whether the agent's Potential Vision mechanic is unlocked.
	IsUpgradeUnlocked *bool `json:"IsUpgradeUnlocked,omitempty"`
	// UpgradeID is the specific Potential Vision upgrade node/buff ID currently active on the agent.
	UpgradeID int `json:"UpgradeId"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (a *AvatarData) UnmarshalJSON(data []byte) error {
	type Alias AvatarData
	var avatar Alias
	if err := json.Unmarshal(data, &avatar); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	avatar.Raw = raw
	avatar.Extra = extra
	*a = AvatarData(avatar)
	return nil
}

func (a AvatarData) MarshalJSON() ([]byte, error) {
	type Alias AvatarData
	return core.MergeKnownExtraAndRawJSON(Alias(a), a.Extra, a.Raw)
}

// SkillLevel contains information about an agent’s skill level.
type SkillLevel struct {
	// Level is the Skill Level.
	Level int `json:"Level"`
	// Index is the Skill Index (see "Skills" section in docs/zzz/api.md).
	Index int `json:"Index"`
}

// EquippedItem contains information about an equipped Drive Disc.
type EquippedItem struct {
	// Slot is the slot index (1-6).
	Slot int `json:"Slot"`
	// Equipment contains Drive Disc data.
	Equipment *Equipment `json:"Equipment"`
}

// Equipment contains information about a Drive Disc.
type Equipment struct {
	// UID is the unique instance ID of this specific Drive Disc on the player's account.
	UID int `json:"Uid"`
	// ID is the static item type ID of the Drive Disc.
	ID int `json:"Id"`
	// Exp is the current experience points invested into the Drive Disc.
	Exp int `json:"Exp"`
	// Level is the Drive Disc's current upgrade level [0-15].
	Level int `json:"Level"`
	// BreakLevel is the total number of times random substats were upgraded as the Drive Disc was leveled up.
	BreakLevel int `json:"BreakLevel"`
	// IsLocked indicates whether the player has locked the Drive Disc to prevent accidental deletion.
	IsLocked *bool `json:"IsLocked,omitempty"`
	// IsAvailable indicates whether the Drive Disc is currently usable.
	IsAvailable *bool `json:"IsAvailable,omitempty"`
	// IsTrash indicates whether the player has marked this Drive Disc as "trash" in their inventory.
	IsTrash *bool `json:"IsTrash,omitempty"`
	// MainPropertyList is the primary, guaranteed stat of the Drive Disc.
	MainPropertyList []Property `json:"MainPropertyList"`
	// RandomPropertyList is the list of randomly rolled substats on the Drive Disc.
	RandomPropertyList []Property `json:"RandomPropertyList"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (e *Equipment) UnmarshalJSON(data []byte) error {
	type Alias Equipment
	var equipment Alias
	if err := json.Unmarshal(data, &equipment); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	equipment.Raw = raw
	equipment.Extra = extra
	*e = Equipment(equipment)
	return nil
}

func (e Equipment) MarshalJSON() ([]byte, error) {
	type Alias Equipment
	return core.MergeKnownExtraAndRawJSON(Alias(e), e.Extra, e.Raw)
}

// Property contains information about a Drive Disc’s stat.
type Property struct {
	// PropertyID is the Property ID (see "Property Id" section in docs/zzz/api.md).
	PropertyID int `json:"PropertyId"`
	// PropertyValue is the Property Base Value.
	PropertyValue int `json:"PropertyValue"`
	// PropertyLevel is the amount of rolls, only matters if substat.
	PropertyLevel int `json:"PropertyLevel"`
}

// Weapon contains information about a W-Engine.
type Weapon struct {
	// UID is the unique instance ID of this specific W-Engine on the player's account.
	UID int `json:"Uid"`
	// ID is the static item type ID of the W-Engine.
	ID int `json:"Id"`
	// Exp is the current experience points invested into the W-Engine.
	Exp int `json:"Exp"`
	// Level is the W-Engine's current level.
	Level int `json:"Level"`
	// BreakLevel is the refinement/modification level of the W-Engine (obtained by feeding duplicate copies).
	BreakLevel int `json:"BreakLevel"`
	// UpgradeLevel is the ascension phase of the W-Engine, determining its level cap.
	UpgradeLevel int `json:"UpgradeLevel"`
	// IsAvailable indicates whether the W-Engine is currently usable.
	IsAvailable *bool `json:"IsAvailable,omitempty"`
	// IsLocked indicates whether the player has locked the W-Engine in their inventory to prevent accidental deletion.
	IsLocked *bool `json:"IsLocked,omitempty"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (w *Weapon) UnmarshalJSON(data []byte) error {
	type Alias Weapon
	var weapon Alias
	if err := json.Unmarshal(data, &weapon); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	weapon.Raw = raw
	weapon.Extra = extra
	*w = Weapon(weapon)
	return nil
}

func (w Weapon) MarshalJSON() ([]byte, error) {
	type Alias Weapon
	return core.MergeKnownExtraAndRawJSON(Alias(w), w.Extra, w.Raw)
}

// SocialDetail contains social profile information.
type SocialDetail struct {
	// MedalList is a list of in-game badges the player has chosen to display on their profile.
	MedalList []Medal `json:"MedalList"`
	// ProfileDetail contains core player profile details (Level, UID, Nickname, etc.).
	ProfileDetail *ProfileDetail `json:"ProfileDetail"`
	// Desc is the custom signature or bio text written by the player on their profile.
	Desc string `json:"Desc"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (s *SocialDetail) UnmarshalJSON(data []byte) error {
	type Alias SocialDetail
	var social Alias
	if err := json.Unmarshal(data, &social); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	social.Raw = raw
	social.Extra = extra
	*s = SocialDetail(social)
	return nil
}

func (s SocialDetail) MarshalJSON() ([]byte, error) {
	type Alias SocialDetail
	return core.MergeKnownExtraAndRawJSON(Alias(s), s.Extra, s.Raw)
}

// Medal contains information about a badge.
type Medal struct {
	// Value is the progress number or tier associated with the badge.
	Value int `json:"Value"`
	// MedalIcon is the ID of the visual icon used for the badge (see https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/medals.json).
	MedalIcon int `json:"MedalIcon"`
	// MedalType is the Badge Category/Type (e.g., Shiyu Defense, Hollow Zero) (see https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/medals.json).
	MedalType int `json:"MedalType"`
	// MedalScore is an additional score metric for the badge (e.g., Endless Tower points or specific challenge score).
	MedalScore int `json:"MedalScore"`
}

// ProfileDetail contains detailed player profile information.
type ProfileDetail struct {
	// UID is the player's unique Zenless Zone Zero account identifier.
	UID int64 `json:"Uid"`
	// Nickname is the display name chosen by the player to represent their account.
	Nickname string `json:"Nickname"`
	// ProfileID is the unique ID of the player's selected avatar image.
	ProfileID int `json:"ProfileId"`
	// Level is the player's current account level, known as the Inter-Knot Level.
	Level int `json:"Level"`
	// Title is the primary Title ID (see https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/titles.json).
	Title int `json:"Title"`
	// CallingCardID is the ID of the namecard background currently equipped to the profile (see https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/namecards.json).
	CallingCardID int `json:"CallingCardId"`
	// AvatarID is the ID identifying the main sibling protagonist (Wise or Belle) selected during the start of the game.
	AvatarID int `json:"AvatarId"`
	// TitleInfo contains detailed metadata regarding the player's currently equipped title and its associated dynamic properties.
	TitleInfo *TitleInfo `json:"TitleInfo"`
	// PlatformType is the numeric ID representing the hardware platform (PC, PS5, or Mobile) currently associated with the account.
	PlatformType int `json:"PlatformType"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (p *ProfileDetail) UnmarshalJSON(data []byte) error {
	type Alias ProfileDetail
	var profile Alias
	if err := json.Unmarshal(data, &profile); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	profile.Raw = raw
	profile.Extra = extra
	*p = ProfileDetail(profile)
	return nil
}

func (p ProfileDetail) MarshalJSON() ([]byte, error) {
	type Alias ProfileDetail
	return core.MergeKnownExtraAndRawJSON(Alias(p), p.Extra, p.Raw)
}

// TitleInfo contains title-related information.
type TitleInfo struct {
	// Title is the Base Title ID.
	Title int `json:"Title"`
	// FullTitle is the Full Title ID (see https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/titles.json).
	FullTitle int `json:"FullTitle"`
	// Args are dynamic values injected into the title's text (e.g., a specific completion time or score).
	Args []any `json:"Args"`
	// KMOHDEAKEFG is an unknown game field, currently unmapped.
	KMOHDEAKEFG []any `json:"KMOHDEAKEFG,omitempty"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (t *TitleInfo) UnmarshalJSON(data []byte) error {
	type Alias TitleInfo
	var title Alias
	if err := json.Unmarshal(data, &title); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	title.Raw = raw
	title.Extra = extra
	*t = TitleInfo(title)
	return nil
}

func (t TitleInfo) MarshalJSON() ([]byte, error) {
	type Alias TitleInfo
	return core.MergeKnownExtraAndRawJSON(Alias(t), t.Extra, t.Raw)
}

// Settings represents build-specific configuration options.
type Settings struct {
	// AdaptiveColor indicates whether adaptive color is enabled.
	AdaptiveColor *bool `json:"adaptiveColor,omitempty"`
	// ArtSource is the source of the image.
	ArtSource *string `json:"artSource,omitempty"`
	// Caption is the caption of the build.
	Caption *string `json:"caption,omitempty"`
	// HonkardWidth is the width of the image.
	HonkardWidth *float64 `json:"honkardWidth,omitempty"`
	// Transform is the transformation applied to the image.
	Transform *string `json:"transform,omitempty"`
	// Raw contains the original API object.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	type Alias Settings
	var settings Alias
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	settings.Raw = raw
	settings.Extra = extra
	*s = Settings(settings)
	return nil
}

func (s Settings) MarshalJSON() ([]byte, error) {
	type Alias Settings
	return core.MergeKnownExtraAndRawJSON(Alias(s), s.Extra, s.Raw)
}
