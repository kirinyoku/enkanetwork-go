package enka

import (
	"encoding/json"

	"github.com/kirinyoku/enkanetwork-go/client/genshin"
	"github.com/kirinyoku/enkanetwork-go/client/hsr"
	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/enkanetwork-go/internal/common"
)

// AvatarBuildsMap is a map where the key is the avatarID and the value is a slice
// of builds for that character, returned in random order. Each build includes an
// "order" field that can be used to sort them for display.
type AvatarBuildsMap map[string][]Build

// Build contains information about a specific character build.
type Build struct {
	ID       int    `json:"id,omitempty"`        // ID of the build
	Name     string `json:"name,omitempty"`      // Name of the build
	AvatarID string `json:"avatar_id,omitempty"` // ID of the avatar (character/agent)
	// AvatarData contains character information, supporting multiple games through a
	// wrapper that holds data for Genshin Impact, Honkai: Star Rail, or Zenless Zone Zero
	AvatarData AvatarDataWrapper `json:"avatar_data"`
	// If a build has a live: true field, it indicates that it is not a saved build but
	// one retrieved from the game’s showcase when the "refresh" button is clicked.
	// During an update, all old live builds are deleted, and new ones are created.
	// Updates are user-initiated, so this data may not be up to date
	Live     bool     `json:"live,omitempty"`
	Settings Settings `json:"settings"`         // Settings contains build-specific configuration data
	Public   bool     `json:"public,omitempty"` // Whether the build is public
	Image    *string  `json:"image,omitempty"`  // URL of the build image
	Hoyo     string   `json:"hoyo,omitempty"`   // Unique hoyo identifier (hoyo_hash)
	Order    int      `json:"order,omitempty"`  // Order of the saved build on the Enka
	HoyoType int      `json:"hoyo_type"`        // ID of the Hoyo game (0 for Genshin, 1 for HSR, 2 for ZZZ)
}

// AvatarDataWrapper is a container struct that holds character data from different game clients.
// It is designed to support multiple games while maintaining a unified interface.
type AvatarDataWrapper struct {
	Genshin *genshin.AvatarInfo `json:"genshin,omitempty"` // Genshin holds character data specific to Genshin Impact
	HSR     *hsr.AvatarDetail   `json:"hsr,omitempty"`     // HSR holds character data specific to Honkai: Star Rail
	ZZZ     *zzz.AvatarData     `json:"zzz,omitempty"`     // ZZZ holds character data specific to Zenless Zone Zero
	Raw     json.RawMessage     `json:"-"`                 // Raw contains the original JSON data for custom unmarshaling or debugging purposes
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle custom JSON unmarshaling
// for AvatarDataWrapper. This method populates the appropriate game-specific avatar data
// field (Genshin, HSR, or ZZZ) based on the incoming JSON structure.
//
// The method attempts to unmarshal the input JSON into the Raw field first, followed by
// each game-specific field (Genshin, HSR, ZZZ). It returns an error if any unmarshaling
// attempt fails, leaving unmatched fields as nil. The Raw field preserves the original
// JSON data for custom processing or debugging.
//
// Parameters:
//   - data: The JSON-encoded byte slice containing the avatar data.
//
// Returns:
//   - error: An error if unmarshaling fails for the Raw field or any game-specific field.
func (a *AvatarDataWrapper) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &a.Raw); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &a.Genshin); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &a.HSR); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &a.ZZZ); err != nil {
		return err
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface to provide custom JSON marshaling
// for the AvatarDataWrapper. This method serializes the appropriate game-specific
// avatar data based on which field is populated.
//
// The method checks each game-specific field in order of priority (Genshin -> HSR -> ZZZ)
// and returns the JSON representation of the first non-nil field it encounters. If no
// game-specific data is present, it returns (nil, nil).
//
// Returns:
//   - []byte: The JSON-encoded data of the populated game-specific field
//   - error: Returns an error if marshaling fails, or nil if successful
func (a *AvatarDataWrapper) MarshalJSON() ([]byte, error) {
	if a.Genshin != nil {
		return json.Marshal(a.Genshin)
	}

	if a.HSR != nil {
		return json.Marshal(a.HSR)
	}

	if a.ZZZ != nil {
		return json.Marshal(a.ZZZ)
	}

	return nil, nil
}

// Hoyos is a map of Hoyo accounts and their metadata. The endpoint returns only
// verified and public accounts (users can hide accounts; unverified accounts are
// hidden by default). Each key is a unique identifier for a hoyo, which can be used
// in subsequent requests to retrieve information about the characters or builds of
// that game account.
type Hoyos map[string]Hoyo

// Hoyo contains information about a specific Hoyo account.
type Hoyo struct {
	UID         int                `json:"uid,omitempty"`          // UID of the game account
	UIDPublic   bool               `json:"uid_public,omitempty"`   // Whether the UID is public
	Public      bool               `json:"public,omitempty"`       // Whether the Hoyo account is public
	Verified    bool               `json:"verified,omitempty"`     // Whether the Hoyo account is verified
	PlayerInfo  *common.PlayerInfo `json:"player_info,omitempty"`  // Player information for the account
	Hash        string             `json:"hash,omitempty"`         // Hash of the game account
	Region      string             `json:"region,omitempty"`       // Region of the game account
	AvatarOrder map[string]int     `json:"avatar_order,omitempty"` // Order of the characters in the game account
	Order       int                `json:"order"`                  // Order of the Hoyo account
	LivePublic  bool               `json:"live_public"`            // Whether the live build is public
	HoyoType    int                `json:"hoyo_type"`              // ID of the Hoyo game (0 for Genshin, 1 for HSR, 2 for ZZZ)
}

// Settings represents build-specific configuration options.
type Settings struct {
	AdaptiveColor *bool    `json:"adaptiveColor,omitempty"` // Whether adaptive color is enabled
	ArtSource     *string  `json:"artSource,omitempty"`     // Source of the image
	Caption       *string  `json:"caption,omitempty"`       // Caption of the build
	HonkardWidth  *float64 `json:"honkardWidth,omitempty"`  // Width of the image
	Transform     *string  `json:"transform,omitempty"`     // Transformation applied to the image
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
