package hsr

import (
	"encoding/json"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/models"
)

// Profile represents the root structure of the response containing player information
// and character data. It serves as the main container for all data returned by the
// EnkaNetwork API for Honkai: Star Rail.
type Profile struct {
	// DetailInfo contains detailed information about the player's account and characters
	DetailInfo *DetailInfo `json:"detailInfo,omitempty"`
	// TTL indicates the seconds remaining until the next request to the game. Until
	// the TTL expires, the endpoint returns cached data — but such requests still
	// count toward the rate limit
	TTL int `json:"ttl,omitempty"`
	// Owner is the Enka profile associated with the provided UID. The response includes
	// an Owner if:
	//   1. The user has an account on the site;
	//   2. The user has added their UID to their profile;
	//   3. The user has verified that the UID belongs to them;
	//   4. The user has set their profile visibility to "public"
	Owner *models.Owner `json:"owner,omitempty"`
	// UID is the unique identifier for the player's account
	UID string `json:"uid,omitempty"`
	// Region indicates the server region of the player (e.g., "ASIA", "USA", "EUROPE")
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

	raw, extra, err := core.PreserveUnknownJSON(data, "detailInfo", "ttl", "owner", "uid", "region")
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

// Build contains information about a specific character build in Honkai: Star Rail.
type Build struct {
	ID         int           `json:"id,omitempty"`          // Unique identifier for the build
	Name       string        `json:"name,omitempty"`        // Name of the build
	AvatarID   string        `json:"avatar_id,omitempty"`   // ID of the character
	AvatarData *AvatarDetail `json:"avatar_data,omitempty"` // Character data (*hsr.AvatarDetail)
	Order      int           `json:"order,omitempty"`       // Order of the saved build on the Enka
	// If a build has a live: true field, it indicates it is not a saved build but one
	// retrieved from the game’s showcase when the "refresh" button is clicked. During
	// an update, all old live builds are deleted, and new ones are created. Updates
	// are user-initiated, so this data may not be up to date
	Live     *bool                      `json:"live,omitempty"`
	Settings Settings                   `json:"settings,omitempty"`  // Settings contains build-specific configuration data
	Public   *bool                      `json:"public,omitempty"`    // Whether the build is public
	Image    *string                    `json:"image,omitempty"`     // URL of the build image
	HoyoType int                        `json:"hoyo_type,omitempty"` // API game/account type discriminator
	Hoyo     string                     `json:"hoyo,omitempty"`      // Unique hoyo identifier (hoyo_hash)
	Raw      json.RawMessage            `json:"-"`                   // Raw contains the original API object
	Extra    map[string]json.RawMessage `json:"-"`                   // Extra contains unknown API fields
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

// DetailInfo contains detailed information about the player's account and characters.
type DetailInfo struct {
	WorldLevel         int                 `json:"worldLevel,omitempty"`         // Player's current world level
	PrivacySettingInfo *PrivacySettingInfo `json:"privacySettingInfo,omitempty"` // Player's privacy settings
	HeadIcon           int                 `json:"headIcon,omitempty"`           // ID of the player's profile icon
	AvatarDetailList   []AvatarDetail      `json:"avatarDetailList,omitempty"`   // List of detailed character information
	Platform           string              `json:"platform,omitempty"`           // Platform where the account is registered
	RecordInfo         *RecordInfo         `json:"recordInfo,omitempty"`         // Player's achievement and collection records
	UID                int                 `json:"uid,omitempty"`                // Player's unique identifier
	Level              int                 `json:"level,omitempty"`              // Player's account level
	Nickname           string              `json:"nickname,omitempty"`           // Player's chosen nickname
	IsDisplayAvatar    *bool               `json:"isDisplayAvatar,omitempty"`    // Whether the player's avatar is displayed
	FriendCount        int                 `json:"friendCount,omitempty"`        // Number of friends the player has
	PersonalCardID     int                 `json:"personalCardId,omitempty"`     // ID of the player's personal card

	// There is no information about this field in the EnkaNetwork HSR API doc.
	// In practice, I always get its value as []map.
	HeadFrameInfo any                        `json:"headFrameInfo,omitempty"`
	Raw           json.RawMessage            `json:"-"` // Raw contains the original API object
	Extra         map[string]json.RawMessage `json:"-"` // Extra contains unknown API fields
}

func (d *DetailInfo) UnmarshalJSON(data []byte) error {
	type Alias DetailInfo
	var detail Alias
	if err := json.Unmarshal(data, &detail); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	detail.Raw = raw
	detail.Extra = extra
	*d = DetailInfo(detail)
	return nil
}

func (d DetailInfo) MarshalJSON() ([]byte, error) {
	type Alias DetailInfo
	return core.MergeKnownExtraAndRawJSON(Alias(d), d.Extra, d.Raw)
}

// PrivacySettingInfo contains the player’s privacy settings for various game features.
type PrivacySettingInfo struct {
	DisplayCollection   *bool `json:"displayCollection,omitempty"`   // Whether collections are visible
	DisplayRecord       *bool `json:"displayRecord,omitempty"`       // Whether records are visible
	DisplayRecordTeam   *bool `json:"displayRecordTeam,omitempty"`   // Whether team records are visible
	DisplayOnlineStatus *bool `json:"displayOnlineStatus,omitempty"` // Whether online status is visible
	DisplayDiary        *bool `json:"displayDiary,omitempty"`        // Whether diary is visible
}

// AvatarDetail contains detailed information about a character in the player’s account.
type AvatarDetail struct {
	RelicList     []Relic                    `json:"relicList,omitempty"`     // List of relics equipped on the character
	Rank          int                        `json:"rank,omitempty"`          // Character's rank
	Level         int                        `json:"level,omitempty"`         // Character's level
	Promotion     int                        `json:"promotion,omitempty"`     // Character's promotion level
	SkillTreeList []SkillTree                `json:"skillTreeList,omitempty"` // List of skill tree points
	Equipment     *Equipment                 `json:"equipment,omitempty"`     // Character's equipment
	AvatarID      int                        `json:"avatarId,omitempty"`      // Character's unique identifier
	EnhancedID    int                        `json:"enhancedId,omitempty"`    // Enhanced avatar identifier
	Assist        *bool                      `json:"_assist,omitempty"`       // Whether the character is in assist position
	Pos           int                        `json:"pos,omitempty"`           // Character's position in the team
	Raw           json.RawMessage            `json:"-"`                       // Raw contains the original API object
	Extra         map[string]json.RawMessage `json:"-"`                       // Extra contains unknown API fields
}

func (a *AvatarDetail) UnmarshalJSON(data []byte) error {
	type Alias AvatarDetail
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
	*a = AvatarDetail(avatar)
	return nil
}

func (a AvatarDetail) MarshalJSON() ([]byte, error) {
	type Alias AvatarDetail
	return core.MergeKnownExtraAndRawJSON(Alias(a), a.Extra, a.Raw)
}

// Relic represents a relic equipped on a character.
type Relic struct {
	MainAffixID  int                        `json:"mainAffixId,omitempty"`  // ID of the main affix
	SubAffixList []SubAffix                 `json:"subAffixList,omitempty"` // List of sub-affixes
	TID          int                        `json:"tid,omitempty"`          // Template ID of the relic
	Type         int                        `json:"type,omitempty"`         // Type of the relic
	Level        int                        `json:"level,omitempty"`        // Level of the relic
	Flat         *Flat                      `json:"_flat,omitempty"`        // Flat data containing properties and set information
	Raw          json.RawMessage            `json:"-"`                      // Raw contains the original API object
	Extra        map[string]json.RawMessage `json:"-"`                      // Extra contains unknown API fields
}

func (r *Relic) UnmarshalJSON(data []byte) error {
	type Alias Relic
	var relic Alias
	if err := json.Unmarshal(data, &relic); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	relic.Raw = raw
	relic.Extra = extra
	*r = Relic(relic)
	return nil
}

func (r Relic) MarshalJSON() ([]byte, error) {
	type Alias Relic
	return core.MergeKnownExtraAndRawJSON(Alias(r), r.Extra, r.Raw)
}

// SubAffix represents a sub-affix on a relic.
type SubAffix struct {
	AffixID int `json:"affixId,omitempty"` // ID of the sub-affix
	Cnt     int `json:"cnt,omitempty"`     // Count of the sub-affix
	Step    int `json:"step,omitempty"`    // Step level of the sub-affix
}

// Flat contains flat data for relics and equipment.
type Flat struct {
	Props   []models.Prop              `json:"props,omitempty"`   // List of properties
	SetName models.StringNumber        `json:"setName,omitempty"` // Name of the set
	SetID   int                        `json:"setID,omitempty"`   // ID of the set
	Raw     json.RawMessage            `json:"-"`                 // Raw contains the original API object
	Extra   map[string]json.RawMessage `json:"-"`                 // Extra contains unknown API fields
}

func (f *Flat) UnmarshalJSON(data []byte) error {
	type Alias Flat
	var flat Alias
	if err := json.Unmarshal(data, &flat); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	flat.Raw = raw
	flat.Extra = extra
	*f = Flat(flat)
	return nil
}

func (f Flat) MarshalJSON() ([]byte, error) {
	type Alias Flat
	return core.MergeKnownExtraAndRawJSON(Alias(f), f.Extra, f.Raw)
}

// SkillTree represents a skill tree point on a character.
type SkillTree struct {
	PointID int `json:"pointId,omitempty"` // ID of the skill tree point
	Level   int `json:"level,omitempty"`   // Level of the skill tree point
}

// EquipmentFlat contains flat data for equipment.
type EquipmentFlat struct {
	Name  models.StringNumber        `json:"name,omitempty"`  // Name of the equipment
	Props []models.Prop              `json:"props,omitempty"` // List of properties
	Raw   json.RawMessage            `json:"-"`               // Raw contains the original API object
	Extra map[string]json.RawMessage `json:"-"`               // Extra contains unknown API fields
}

func (f *EquipmentFlat) UnmarshalJSON(data []byte) error {
	type Alias EquipmentFlat
	var flat Alias
	if err := json.Unmarshal(data, &flat); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	flat.Raw = raw
	flat.Extra = extra
	*f = EquipmentFlat(flat)
	return nil
}

func (f EquipmentFlat) MarshalJSON() ([]byte, error) {
	type Alias EquipmentFlat
	return core.MergeKnownExtraAndRawJSON(Alias(f), f.Extra, f.Raw)
}

// Equipment represents a character's equipment (weapon).
type Equipment struct {
	Rank      int                        `json:"rank,omitempty"`      // Rank of the equipment
	TID       int                        `json:"tid,omitempty"`       // Template ID of the equipment
	Promotion int                        `json:"promotion,omitempty"` // Promotion level of the equipment
	Level     int                        `json:"level,omitempty"`     // Level of the equipment
	Flat      *EquipmentFlat             `json:"_flat,omitempty"`     // Flat data containing properties
	Raw       json.RawMessage            `json:"-"`                   // Raw contains the original API object
	Extra     map[string]json.RawMessage `json:"-"`                   // Extra contains unknown API fields
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

// RecordInfo contains various achievement and collection records for the player.
type RecordInfo struct {
	AchievementCount       int                        `json:"achievementCount,omitempty"`       // Number of achievements completed
	BookCount              int                        `json:"bookCount,omitempty"`              // Number of books collected
	AvatarCount            int                        `json:"avatarCount,omitempty"`            // Number of characters unlocked
	EquipmentCount         int                        `json:"equipmentCount,omitempty"`         // Number of equipment pieces
	MusicCount             int                        `json:"musicCount,omitempty"`             // Number of music tracks
	RelicCount             int                        `json:"relicCount,omitempty"`             // Number of relics collected
	ChallengeInfo          *any                       `json:"challengeInfo,omitempty"`          // Challenge-related information
	MaxRogueChallengeScore int                        `json:"maxRogueChallengeScore,omitempty"` // Highest score in rogue challenges
	Raw                    json.RawMessage            `json:"-"`                                // Raw contains the original API object
	Extra                  map[string]json.RawMessage `json:"-"`                                // Extra contains unknown API fields
}

func (r *RecordInfo) UnmarshalJSON(data []byte) error {
	type Alias RecordInfo
	var record Alias
	if err := json.Unmarshal(data, &record); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	record.Raw = raw
	record.Extra = extra
	*r = RecordInfo(record)
	return nil
}

func (r RecordInfo) MarshalJSON() ([]byte, error) {
	type Alias RecordInfo
	return core.MergeKnownExtraAndRawJSON(Alias(r), r.Extra, r.Raw)
}

// Settings represents build-specific configuration options.
type Settings struct {
	AdaptiveColor *bool                      `json:"adaptiveColor,omitempty"` // Whether adaptive color is enabled
	ArtSource     *string                    `json:"artSource,omitempty"`     // Source of the image
	Caption       *string                    `json:"caption,omitempty"`       // Caption of the build
	HonkardWidth  *float64                   `json:"honkardWidth,omitempty"`  // Width of the image
	Transform     *string                    `json:"transform,omitempty"`     // Transformation applied to the image
	Raw           json.RawMessage            `json:"-"`                       // Raw contains the original API object
	Extra         map[string]json.RawMessage `json:"-"`                       // Extra contains unknown API fields
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
