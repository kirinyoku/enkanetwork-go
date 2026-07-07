package genshin

import (
	"encoding/json"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/models"
)

// ------------------------------- IMPORTANT --------------------------------------
// For detailed information on properties, refer to the EnkaNetwork API — Genshin
// Impact documentation (https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md).
// -------------------------------------------------------------------------------

// Profile represents the root structure of the response containing player information
// and character data. It serves as the main container for all data returned by the
// EnkaNetwork API for Genshin Impact.
type Profile struct {
	// PlayerInfo contains basic information about the game account from the player's showcase
	PlayerInfo models.PlayerInfo `json:"playerInfo"`
	// AvatarInfoList contains detailed information for each character in the showcase.
	// If missing, the showcase is either hidden by the player or contains no characters.
	// The GetPlayerInfo method always returns AvatarInfoList as an empty slice.
	AvatarInfoList []AvatarInfo `json:"avatarInfoList,omitempty"`
	// Owner is the Enka profile associated with the provided UID.
	// The response includes an Owner if:
	//   1. The user has an account on the site;
	//   2. The user has added their UID to their profile;
	//   3. The user has verified that the UID belongs to them;
	//   4. The user has set their profile visibility to "public"
	Owner *models.Owner `json:"owner,omitempty"`
	// TTL indicates the seconds remaining until the next request to the game.
	// Until the TTL expires, the endpoint returns cached data — but such requests still
	// count toward the rate limit. Cache data locally and use the TTL to avoid
	// requesting the UID again until it expires
	TTL int `json:"ttl"`
	// UID is the player's UID in Genshin Impact
	UID string `json:"uid,omitempty"`
	// Region is the server region of the player (e.g., "NA", "EU", "Asia", "TW, HK, MO").
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

	raw, extra, err := core.PreserveUnknownJSON(data, "playerInfo", "avatarInfoList", "owner", "ttl", "uid", "region")
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

// AvatarInfo contains detailed information for characters in the showcase.
type AvatarInfo struct {
	AvatarID                int                        `json:"avatarId,omitempty"`                // Character ID
	PropMap                 map[string]Prop            `json:"propMap,omitempty"`                 // Maps prop_type to Prop, listing character properties (e.g., `{"4001": {"type": 4001, "ival": "90", "val": "90"}}`, where 4001 is character level)
	TalentIDList            []int                      `json:"talentIdList,omitempty"`            // List of constellation IDs (empty if constellation level is 0)
	FightPropMap            map[string]float64         `json:"fightPropMap,omitempty"`            // Maps IDs to values for character's combat properties (see https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md#fightprop for IDs)
	SkillDepotID            int                        `json:"skillDepotId,omitempty"`            // Character skill set ID
	InherentProudSkillList  []int                      `json:"inherentProudSkillList,omitempty"`  // List of unlocked skill IDs
	SkillLevelMap           map[string]int             `json:"skillLevelMap,omitempty"`           // Maps skill_id to level for character skills
	ProudSkillExtraLevelMap map[string]int             `json:"proudSkillExtraLevelMap,omitempty"` // Maps proud_skill_id to level for skills enhanced by constellations (e.g., 3rd and 5th)
	EquipList               []Equip                    `json:"equipList,omitempty"`               // List of equipment: weapon and artifacts
	FetterInfo              *FetterInfo                `json:"fetterInfo,omitempty"`              // Character friendship level information
	Raw                     json.RawMessage            `json:"-"`                                 // Raw contains the original API object
	Extra                   map[string]json.RawMessage `json:"-"`                                 // Extra contains unknown API fields
}

func (a *AvatarInfo) UnmarshalJSON(data []byte) error {
	type Alias AvatarInfo
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
	*a = AvatarInfo(avatar)
	return nil
}

func (a AvatarInfo) MarshalJSON() ([]byte, error) {
	type Alias AvatarInfo
	return core.MergeKnownExtraAndRawJSON(Alias(a), a.Extra, a.Raw)
}

// Build contains information about a specific character build in Genshin Impact.
type Build struct {
	ID         int         `json:"id,omitempty"`          // Unique identifier for the build
	Name       string      `json:"name,omitempty"`        // Name of the build
	AvatarID   string      `json:"avatar_id,omitempty"`   // ID of the character
	AvatarData *AvatarInfo `json:"avatar_data,omitempty"` // Character data (*genshin.AvatarInfo)
	Order      int         `json:"order,omitempty"`       // Order of the saved build on Enka
	// If a build has a live: true field, it indicates it is not a saved build but one
	// retrieved from the game’s showcase when the "refresh" button is clicked. During
	// an update, all old live builds are deleted, and new ones are created. Updates
	// are user-initiated, so this data may not be up to date
	Live     *bool                      `json:"live,omitempty"`
	Settings *Settings                  `json:"settings,omitempty"`  // Settings contains build-specific configuration data
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

// Equip contains detailed information about a character's equipment (weapon and artifacts).
type Equip struct {
	ItemID    int                        `json:"itemId,omitempty"`    // Equipment ID
	Reliquary *Reliquary                 `json:"reliquary,omitempty"` // Artifact base information
	Weapon    *Weapon                    `json:"weapon,omitempty"`    // Weapon base information
	Flat      any                        `json:"flat,omitempty"`      // Detailed information about the equipment
	Raw       json.RawMessage            `json:"-"`                   // Raw contains the original API object
	Extra     map[string]json.RawMessage `json:"-"`                   // Extra contains unknown API fields
}

func (e *Equip) UnmarshalJSON(data []byte) error {
	type Alias Equip
	var equip Alias
	if err := json.Unmarshal(data, &equip); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	equip.Raw = raw
	equip.Extra = extra
	*e = Equip(equip)
	return nil
}

func (e Equip) MarshalJSON() ([]byte, error) {
	type Alias Equip
	return core.MergeKnownExtraAndRawJSON(Alias(e), e.Extra, e.Raw)
}

// FetterInfo contains information about a character's friendship level.
type FetterInfo struct {
	ExpLevel int `json:"expLevel"` // Character's friendship level in the game
}

// FlatReliquary contains detailed information about an artifact.
type FlatReliquary struct {
	NameTextMapHash    models.StringNumber        `json:"nameTextMapHash,omitempty"`    // Hash for equipment name (see localizations: https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md#localizations)
	RankLevel          int                        `json:"rankLevel,omitempty"`          // Rarity level of the equipment
	ItemType           string                     `json:"itemType,omitempty"`           // Equipment type: weapon or artifact
	Icon               string                     `json:"icon,omitempty"`               // Equipment icon name (see https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md#icons-and-images)
	EquipType          string                     `json:"equipType,omitempty"`          // Artifact type
	SetID              int                        `json:"setId,omitempty"`              // Artifact set ID
	SetNameTextMapHash models.StringNumber        `json:"setNameTextMapHash,omitempty"` // Hash for artifact set name (see localizations: https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md#localizations)
	ReliquaryMainstat  *ReliquaryMainstat         `json:"reliquaryMainstat,omitempty"`  // Artifact main stat
	ReliquarySubstats  []ReliquarySubstat         `json:"reliquarySubStats,omitempty"`  // List of artifact substats
	Raw                json.RawMessage            `json:"-"`                            // Raw contains the original API object
	Extra              map[string]json.RawMessage `json:"-"`                            // Extra contains unknown API fields
}

func (f *FlatReliquary) UnmarshalJSON(data []byte) error {
	type Alias FlatReliquary
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
	*f = FlatReliquary(flat)
	return nil
}

func (f FlatReliquary) MarshalJSON() ([]byte, error) {
	type Alias FlatReliquary
	return core.MergeKnownExtraAndRawJSON(Alias(f), f.Extra, f.Raw)
}

// WeaponStat contains information about a weapon’s stat.
type WeaponStat struct {
	AppendPropID string  `json:"appendPropId,omitempty"` // Equipment append property name (see definitions: https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md#appendprop)
	StatValue    float64 `json:"statValue,omitempty"`    // Value of the property
}

// FlatWeapon contains detailed information about a weapon.
type FlatWeapon struct {
	NameTextMapHash models.StringNumber        `json:"nameTextMapHash,omitempty"` // Hash for equipment name (see localizations: https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md#localizations)
	ItemType        string                     `json:"itemType,omitempty"`        // Equipment type: weapon or artifact
	Icon            string                     `json:"icon,omitempty"`            // Equipment icon name (see https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md#icons-and-images)
	WeaponStats     []WeaponStat               `json:"weaponStats,omitempty"`     // List of weapon stats: base ATK and substat
	Raw             json.RawMessage            `json:"-"`                         // Raw contains the original API object
	Extra           map[string]json.RawMessage `json:"-"`                         // Extra contains unknown API fields
}

func (f *FlatWeapon) UnmarshalJSON(data []byte) error {
	type Alias FlatWeapon
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
	*f = FlatWeapon(flat)
	return nil
}

func (f FlatWeapon) MarshalJSON() ([]byte, error) {
	type Alias FlatWeapon
	return core.MergeKnownExtraAndRawJSON(Alias(f), f.Extra, f.Raw)
}

// Weapon contains information about a weapon’s level and refinement.
type Weapon struct {
	Level        int                        `json:"level,omitempty"`        // Weapon level
	PromoteLevel int                        `json:"promoteLevel,omitempty"` // Weapon ascension level
	AffixMap     map[string]int             `json:"affixMap,omitempty"`     // Maps to weapon refinement levels (0–4)
	Raw          json.RawMessage            `json:"-"`                      // Raw contains the original API object
	Extra        map[string]json.RawMessage `json:"-"`                      // Extra contains unknown API fields
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

// ReliquarySubstat contains information about an artifact’s substat.
type ReliquarySubstat struct {
	AppendPropID string  `json:"appendPropId,omitempty"` // Artifact append property name
	StatValue    float64 `json:"statValue,omitempty"`    // Value of the substat
}

// ReliquaryMainstat contains information about an artifact’s main stat.
type ReliquaryMainstat struct {
	MainPropID string  `json:"mainPropId,omitempty"` // Artifact main property name
	StatValue  float64 `json:"statValue,omitempty"`  // Value of the main stat
}

// Reliquary contains base information about an artifact.
type Reliquary struct {
	Level            int                        `json:"level,omitempty"`            // Artifact level [1-21]
	MainPropID       int                        `json:"mainPropId,omitempty"`       // Artifact main stat ID
	AppendPropIDList []int                      `json:"appendPropIdList,omitempty"` // List of artifact substat IDs
	Raw              json.RawMessage            `json:"-"`                          // Raw contains the original API object
	Extra            map[string]json.RawMessage `json:"-"`                          // Extra contains unknown API fields
}

func (r *Reliquary) UnmarshalJSON(data []byte) error {
	type Alias Reliquary
	var reliquary Alias
	if err := json.Unmarshal(data, &reliquary); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	reliquary.Raw = raw
	reliquary.Extra = extra
	*r = Reliquary(reliquary)
	return nil
}

func (r Reliquary) MarshalJSON() ([]byte, error) {
	type Alias Reliquary
	return core.MergeKnownExtraAndRawJSON(Alias(r), r.Extra, r.Raw)
}

// Prop contains information about a character property.
type Prop struct {
	Type int    `json:"type,omitempty"` // ID of the property type
	Ival string `json:"ival,omitempty"` // Ignore it
	Val  string `json:"val,omitempty"`  // Value of the property
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
