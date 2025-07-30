package models

// PlayerInfo contains basic information about the player's game account from their showcase.
type PlayerInfo struct {
	// -------------------------------------- Common Fields --------------------------------
	Nickname   string `json:"nickname,omitempty"`   // Player nickname
	Level      int    `json:"level,omitempty"`      // Player level
	Signature  string `json:"signature,omitempty"`  // Profile signature
	WorldLevel int    `json:"worldLevel,omitempty"` // Player world level
	// ------------------------------------ Genshin Impact ---------------------------------
	NameCardId           int              `json:"nameCardId,omitempty"`           // Profile namecard ID
	FinishAchievementNum int              `json:"finishAchievementNum,omitempty"` // Number of completed achievements
	ShowAvatarInfoList   []ShowAvatarInfo `json:"showAvatarInfoList,omitempty"`   // List of character information (IDs, levels, skins, constellations, elements).
	ShowNameCardIdList   []int            `json:"showNameCardIdList,omitempty"`   // List of namecard IDs
	ProfilePicture       *ProfilePicture  `json:"profilePicture,omitempty"`       // Player profile picture
	TheaterActIndex      int              `json:"theaterActIndex,omitempty"`      // Imaginarium Theater act
	TheaterModeIndex     int              `json:"theaterModeIndex,omitempty"`     // Imaginarium Theater difficulty mode
	TheaterStarIndex     int              `json:"theaterStarIndex,omitempty"`     // Imaginarium Theater stars earned
	IsShowAvatarTalent   bool             `json:"isShowAvatarTalent,omitempty"`   // Whether the constellation level is displayed
	FetterCount          int              `json:"fetterCount,omitempty"`          // Number of characters at maximum friendship level
	TowerStarIndex       int              `json:"towerStarIndex,omitempty"`       // Spiral Abyss stars earned
	TowerFloorIndex      int              `json:"towerFloorIndex,omitempty"`      // Spiral Abyss floor reached
	TowerLevelIndex      int              `json:"towerLevelIndex,omitempty"`      // Spiral Abyss chamber reached
	StygianIndex         int              `json:"stygianIndex,omitempty"`         // Stygian Onslaught difficulty mode
	StygianSeconds       int              `json:"stygianSeconds,omitempty"`       // Stygian Onslaught time in seconds
	// ------------------------------------ HONKAI: STAR RAIL ------------------------------------
	HeadIcon           int                 `json:"headIcon,omitempty"`           // Profile picture ID
	Birthday           int                 `json:"birthday,omitempty"`           // Player birthday
	Platform           string              `json:"platform,omitempty"`           // Platform (e.g. PC, Mobile)
	FriendCount        int                 `json:"friendCount,omitempty"`        // Number of friends
	IsDisplayAvatar    bool                `json:"isDisplayAvatar,omitempty"`    // Whether characters are displayed
	AvatarDetailList   []AvatarDetail      `json:"avatarDetailList,omitempty"`   // List of character details
	RecordInfo         *RecordInfo         `json:"recordInfo,omitempty"`         // Player record information
	PrivacySettingInfo *PrivacySettingInfo `json:"privacySettingInfo,omitempty"` // Player privacy settings
	// ------------------------------------ ZENLESS ZONE ZERO ------------------------------------
	Desc          string         `json:"Desc,omitempty"`          // Profile signature
	MedalList     []Medal        `json:"MedalList,omitempty"`     // List of badges
	ProfileDetail *ProfileDetail `json:"ProfileDetail,omitempty"` // Profile details
}

// Owner represents an EnkaNetwork user profile associated with a game account.
type Owner struct {
	ID       int             `json:"id,omitempty"`       // User ID
	Hash     string          `json:"hash,omitempty"`     // User hash
	Username string          `json:"username,omitempty"` // Enka username
	Profile  *PatreonProfile `json:"profile,omitempty"`  // Patreon profile data for Patreon members
}

// PatreonProfile contains Patreon-related information for an Enka user.
type PatreonProfile struct {
	Bio      string `json:"bio,omitempty"`       // User bio from Patreon
	Level    int    `json:"level,omitempty"`     // Patreon membership level
	Avatar   string `json:"avatar,omitempty"`    // Profile picture on Enka
	ImageURL string `json:"image_url,omitempty"` // Profile picture from Patreon
}

// Medal represents a badge in Zenless Zone Zero.
type Medal struct {
	Value     int `json:"Value,omitempty"`     // Progress number
	MedalIcon int `json:"MedalIcon,omitempty"` // Icon ID
	MedalType int `json:"MedalType,omitempty"` // Badge type (see https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md#badge-type)
}

// ProfileDetail contains detailed player profile information for Zenless Zone Zero.
type ProfileDetail struct {
	Uid           int       `json:"Uid,omitempty"`           // Player UID
	Level         int       `json:"Level,omitempty"`         // Inter-Knot Level
	Title         int       `json:"Title,omitempty"`         // Title ID
	AvatarId      int       `json:"AvatarId,omitempty"`      // Main Character ID (Wise or Belle)
	Nickname      string    `json:"Nickname,omitempty"`      // Player nickname
	ProfileId     int       `json:"ProfileId,omitempty"`     // Profile picture ID
	TitleInfo     TitleInfo `json:"TitleInfo,omitempty"`     // Title information
	PlatformType  int       `json:"PlatformType,omitempty"`  // Platform type (1: PC, 2: Mobile)
	CallingCardId int       `json:"CallingCardId,omitempty"` // Namecard ID
}

// TitleInfo contains title-related information for Zenless Zone Zero.
type TitleInfo struct {
	Title       int   `json:"Title,omitempty"`       // Title ID
	ECJPEHHALAO int   `json:"ECJPEHHALAO,omitempty"` // ????????
	HFKHLLBMPHM []any `json:"HFKHLLBMPHM,omitempty"` // ????????
}

// ShowAvatarInfo contains information about a character displayed in the player's showcase.
type ShowAvatarInfo struct {
	AvatarID    int `json:"avatarId,omitempty"`    // Character ID
	Level       int `json:"level,omitempty"`       // Character level
	EnergyType  int `json:"energyType,omitempty"`  // Character element ID
	CostumeId   int `json:"costumeId,omitempty"`   // ID of character's skin
	TalentLevel int `json:"talentLevel,omitempty"` // Character constellation level
}

// ProfilePicture represents a player’s profile picture.
type ProfilePicture struct {
	AvatarID int `json:"avatarId,omitempty"` // Character ID of profile picture
}

// RecordInfo contains player record statistics for Honkai: Star Rail.
type RecordInfo struct {
	BookCount              int            `json:"bookCount,omitempty"`              // Number of books collected
	MusicCount             int            `json:"musicCount,omitempty"`             // Number of music tracks collected
	RelicCount             int            `json:"relicCount,omitempty"`             // Number of relics collected
	AvatarCount            int            `json:"avatarCount,omitempty"`            // Number of characters owned
	ChallengeInfo          *ChallengeInfo `json:"challengeInfo,omitempty"`          // Challenge-related information
	EquipmentCount         int            `json:"equipmentCount,omitempty"`         // Number of equipment items
	AchievementCount       int            `json:"achievementCount,omitempty"`       // Number of achievements completed
	MaxRogueChallengeScore int            `json:"maxRogueChallengeScore,omitempty"` // Maximum rogue challenge score
}

// ChallengeInfo represents challenge-related data for Honkai: Star Rail (currently empty).
type ChallengeInfo struct {
}

// PrivacySettingInfo contains privacy settings for a Honkai: Star Rail player.
type PrivacySettingInfo struct {
	DisplayDiary        bool `json:"displayDiary,omitempty"`        // Whether the diary is displayed
	DisplayRecord       bool `json:"displayRecord,omitempty"`       // Whether records are displayed
	DisplayCollection   bool `json:"displayCollection,omitempty"`   // Whether collections are displayed
	DisplayRecordTeam   bool `json:"displayRecordTeam,omitempty"`   // Whether the record team is displayed
	DisplayOnlineStatus bool `json:"displayOnlineStatus,omitempty"` // Whether online status is displayed
}
