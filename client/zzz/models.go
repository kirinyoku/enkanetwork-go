package zzz

import "github.com/kirinyoku/enkanetwork-go/models"

// ------------------------------- IMPORTANT --------------------------------------
// For detailed information on properties, refer to the EnkaNetwork API — Zenless
// Zone Zero documentation (https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md).
// -------------------------------------------------------------------------------

// Profile represents the root structure of the response containing player information
// for Zenless Zone Zero.
type Profile struct {
	// PlayerInfo contains basic information about the game account from the player's showcase
	PlayerInfo PlayerInfo `json:"PlayerInfo"`
	// TTL indicates the seconds remaining until the next request to the game. Until
	// the TTL expires, the endpoint returns cached data — but such requests still
	// count toward the rate limit
	TTL int `json:"ttl"` // Seconds until next update
	// Owner is the Enka profile associated with the provided UID. The response includes
	// an Owner if:
	//   1. The user has an account on the site;
	//   2. The user has added their UID to their profile;
	//   3. The user has verified that the UID belongs to them;
	//   4. The user has set their profile visibility to "public"
	Owner *models.Owner `json:"owner,omitempty"`
	// UID is the player's UID in Zenless Zone Zero.
	UID string `json:"uid,omitempty"`
}

// Build contains information about a specific character build in Zenless Zone Zero.
type Build struct {
	ID         int         `json:"id"`          // ID of the build
	Name       string      `json:"name"`        // Name of the build
	AvatarID   string      `json:"avatar_id"`   // ID of the agent
	AvatarData *AvatarData `json:"avatar_data"` // Agent data (*zzz.AvatarData)
	Order      int         `json:"order"`       // Order of the saved build on the Enka
	// If a build has a live: true field, it indicates it is not a saved build but one
	// retrieved from the game’s showcase when the "refresh" button is clicked. During
	// an update, all old live builds are deleted, and new ones are created. Updates
	// are user-initiated, so this data may not be up to date.
	Live     bool      `json:"live"`
	Settings *Settings `json:"settings"`  // Settings contains build-specific configuration data
	Public   bool      `json:"public"`    // Whether the build is public
	Image    *string   `json:"image"`     // URL of the build image
	HoyoType int       `json:"hoyo_type"` // ID of the Hoyo game (0 for Genshin, 1 for HSR, 2 for ZZZ)
	Hoyo     string    `json:"hoyo"`      // Unique hoyo identifier (hoyo_hash)
}

// PlayerInfo contains basic information about a player's game account.
type PlayerInfo struct {
	SocialDetail   *SocialDetail   `json:"SocialDetail"`   // Social profile details
	ShowcaseDetail *ShowcaseDetail `json:"ShowcaseDetail"` // Agent showcase details
}

// ShowcaseDetail contains a list of agents in the player’s showcase.
type ShowcaseDetail struct {
	AvatarList []AvatarData `json:"AvatarList"` // List of agents
}

// AvatarData contains detailed information about an agent.
type AvatarData struct {
	ID                   int            `json:"Id"`                   // Agent ID
	Exp                  int            `json:"Exp"`                  // Agent experience
	Level                int            `json:"Level"`                // Agent level
	PromotionLevel       int            `json:"PromotionLevel"`       // Agent promotion level
	TalentLevel          int            `json:"TalentLevel"`          // Agent mindscape level
	SkinID               int            `json:"SkinId"`               // Agent skin ID
	CoreSkillEnhancement int            `json:"CoreSkillEnhancement"` // Core skill unlocked enhancements (A, B, C, D, E, F)
	TalentToggleList     []bool         `json:"TalentToggleList"`     // Mindscape Cinema visual toggles
	WeaponEffectState    int            `json:"WeaponEffectState"`    // W-Engine signature special effect state (0: None, 1: OFF, 2: ON)
	ClaimedRewardList    []int          `json:"ClaimedRewardList"`    // Agent promotion rewards
	ObtainmentTimestamp  int64          `json:"ObtainmentTimestamp"`  // Agent obtainment timestamp
	Weapon               *Weapon        `json:"Weapon"`               // Equipped W-Engine
	SkillLevelList       []SkillLevel   `json:"SkillLevelList"`       // List of agent skill levels (see definitions: https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md#skills for indexes)
	EquippedList         []EquippedItem `json:"EquippedList"`         // List of equipped Drive Discs
	IsFavorite           bool           `json:"IsFavorite"`           // Whether the agent is marked as favorite
	WeaponUID            int            `json:"WeaponUid"`            // W-Engine UID
	// I couldn't find any information about these fields in the Enka API documentation.
	// If you have any information, please let me know.
	IsUpgradeUnlocked bool `json:"IsUpgradeUnlocked"` // ??
	UpgradeID         int  `json:"UpgradeId"`         // ??
}

// SkillLevel contains information about an agent’s skill level.
type SkillLevel struct {
	Level int `json:"Level"` // Skill level
	Index int `json:"Index"` // Skill index
}

// EquippedItem contains information about an equipped Drive Disc.
type EquippedItem struct {
	Slot      int        `json:"Slot"`      // Slot index
	Equipment *Equipment `json:"Equipment"` // Drive Disc data
}

// Equipment contains information about a Drive Disc.
type Equipment struct {
	UID                int        `json:"Uid"`                // Drive Disc UID
	ID                 int        `json:"Id"`                 // Drive Disc ID
	Exp                int        `json:"Exp"`                // Drive Disc experience
	Level              int        `json:"Level"`              // Drive Disc level [0-15]
	BreakLevel         int        `json:"BreakLevel"`         // Number of random stat procs
	IsLocked           bool       `json:"IsLocked"`           // Whether the Drive Disc is locked
	IsAvailable        bool       `json:"IsAvailable"`        // Whether the Drive Disc is available
	IsTrash            bool       `json:"IsTrash"`            // Whether the Drive Disc is marked as trash
	MainPropertyList   []Property `json:"MainPropertyList"`   // Drive Disc main stat (see Stat definitions: https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md#Stat for additional information)
	RandomPropertyList []Property `json:"RandomPropertyList"` // Drive Disc substats (see Stat definitions: https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md#Stat for additional information)
}

// Property contains information about a Drive Disc’s stat.
type Property struct {
	PropertyID    int `json:"PropertyId"`    // Property ID (see definitions: https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md#property-id for IDs)
	PropertyValue int `json:"PropertyValue"` // Property base value
	PropertyLevel int `json:"PropertyLevel"` // Number of rolls (applies only to substats)
}

// Weapon contains information about a W-Engine.
type Weapon struct {
	UID          int  `json:"Uid"`          // W-Engine UID
	ID           int  `json:"Id"`           // W-Engine ID
	Exp          int  `json:"Exp"`          // W-Engine experience
	Level        int  `json:"Level"`        // W-Engine level
	BreakLevel   int  `json:"BreakLevel"`   // W-Engine modification level
	UpgradeLevel int  `json:"UpgradeLevel"` // W-Engine phase level
	IsAvailable  bool `json:"IsAvailable"`  // Whether the W-Engine is available
	IsLocked     bool `json:"IsLocked"`     // Whether the W-Engine is locked
}

// SocialDetail contains social profile information.
type SocialDetail struct {
	MedalList     []Medal        `json:"MedalList"`     // List of badges
	ProfileDetail *ProfileDetail `json:"ProfileDetail"` // Profile details
	Desc          string         `json:"Desc"`          // Profile signature
}

// Medal contains information about a badge.
type Medal struct {
	Value      int `json:"Value"`      // Progress number
	MedalIcon  int `json:"MedalIcon"`  // Icon ID
	MedalType  int `json:"MedalType"`  // Badge type (see https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md#badge-type)
	MedalScore int `json:"MedalScore"` // Badge score
}

// ProfileDetail contains detailed player profile information.
type ProfileDetail struct {
	UID           int64      `json:"Uid"`           // Player UID
	Nickname      string     `json:"Nickname"`      // Player Nickname
	ProfileID     int        `json:"ProfileId"`     // Profile Picture ID
	Level         int        `json:"Level"`         // Inter-Knot Level
	Title         int        `json:"Title"`         // Title ID
	CallingCardID int        `json:"CallingCardId"` // Namecard ID
	AvatarID      int        `json:"AvatarId"`      // Main Character ID (Wise or Belle)
	TitleInfo     *TitleInfo `json:"TitleInfo"`     // Title information
	PlatformType  int        `json:"PlatformType"`  // Platform Type
}

// TitleInfo contains title-related information.
type TitleInfo struct {
	// I couldn't find any information about these fields in the Enka API documentation.
	// If you have any information, please let me know.
	Title     int   `json:"Title"`     // ??
	FullTitle int   `json:"FullTitle"` // ??
	Args      []any `json:"Args"`      // ??
}

// Settings represents build-specific configuration options.
type Settings struct {
	AdaptiveColor *bool    `json:"adaptiveColor,omitempty"` // Whether adaptive color is enabled
	ArtSource     *string  `json:"artSource,omitempty"`     // Source of the image
	Caption       *string  `json:"caption,omitempty"`       // Caption of the build
	HonkardWidth  *float64 `json:"honkardWidth,omitempty"`  // Width of the image
	Transform     *string  `json:"transform,omitempty"`     // Transformation applied to the image
}
