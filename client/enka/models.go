package enka

import (
	"encoding/json"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/endfield"
	"github.com/kirinyoku/enkanetwork-go/client/genshin"
	"github.com/kirinyoku/enkanetwork-go/client/hsr"
	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/models"
)

// ------------------------------------------IMPORTANT------------------------------------------
// For detailed information on properties, refer to the EnkaNetwork API — Enka.Network Profiles
// documentation (https://github.com/EnkaNetwork/API-docs/blob/master/docs/enka/api.md).
// ---------------------------------------------------------------------------------------------
//
// Owner
// ├── id
// ├── is_staff
// ├── hash
// ├── username
// └── profile (*PatreonProfile)
//     ├── bio
//     ├── level
//     ├── avatar
//     └── image_url
//
// Hoyo
// ├── user (*Owner)
// ├── uid
// ├── uid_public
// ├── public
// ├── live_public
// ├── verified
// ├── player_info
// ├── hash
// ├── region
// ├── order
// ├── avatar_order
// ├── hoyo_type
// └── live_data_hash
//
// Build
// ├── id
// ├── name
// ├── avatar_id
// ├── owner
// ├── avatar_data (*AvatarDataWrapper)
// ├── live
// ├── settings (*Settings)
// │   ├── adaptiveColor
// │   ├── artSource
// │   ├── caption
// │   ├── honkardWidth
// │   └── transform
// ├── public
// ├── image
// ├── hoyo
// ├── order
// └── hoyo_type

// AvatarBuildsMap is a map where the key is the avatarID and the value is a slice
// of builds for that character, returned in random order. Each build includes an
// "order" field that can be used to sort them for display.
type AvatarBuildsMap map[string][]Build

// CacheTTL implements the core.Cacheable interface.
func (a AvatarBuildsMap) CacheTTL() time.Duration {
	return 5 * time.Minute
}

// HoyoType identifies which game an account or build belongs to.
// The API field is named "hoyo_type" even for non-HoYoverse games.
type HoyoType int

const (
	HoyoTypeGenshin HoyoType = iota
	HoyoTypeHSR
	HoyoTypeZZZ
	HoyoTypeArknightsEndfield
)

// Build contains information about a specific character build.
type Build struct {
	// ID is the ID of the build.
	ID int `json:"id,omitempty"`
	// Name is the name of the build.
	Name string `json:"name,omitempty"`
	// AvatarID is the ID of the avatar (character/agent).
	AvatarID string `json:"avatar_id,omitempty"`
	// Owner is the owner identifier associated with the build.
	Owner string `json:"owner,omitempty"`
	// AvatarData contains character information for supported games. Unsupported
	// game payloads are preserved as raw JSON.
	AvatarData AvatarDataWrapper `json:"avatar_data"`
	// Live indicates that it is not a saved build but one retrieved from the game’s showcase when the "refresh" button is clicked. During an update, all old live builds are deleted, and new ones are created. Updates are user-initiated, so this data may not be up to date.
	Live *bool `json:"live,omitempty"`
	// Settings contains build-specific configuration data.
	Settings Settings `json:"settings"`
	// Public indicates whether the build is public.
	Public *bool `json:"public,omitempty"`
	// Image is the URL of the build image.
	Image *string `json:"image"`
	// Hoyo is the unique hoyo identifier (hoyo_hash).
	Hoyo string `json:"hoyo,omitempty"`
	// Order is the order of the saved build on the Enka.
	Order string `json:"order,omitempty"`
	// HoyoType is the game type discriminator (0: Genshin, 1: HSR, 2: ZZZ, 3: Arknights Endfield).
	HoyoType HoyoType `json:"hoyo_type"`
}

// UnmarshalJSON decodes avatar_data using the build's hoyo_type discriminator.
func (b *Build) UnmarshalJSON(data []byte) error {
	type buildJSON struct {
		ID         int             `json:"id,omitempty"`
		Name       string          `json:"name,omitempty"`
		AvatarID   string          `json:"avatar_id,omitempty"`
		Owner      string          `json:"owner,omitempty"`
		AvatarData json.RawMessage `json:"avatar_data"`
		Live       *bool           `json:"live,omitempty"`
		Settings   Settings        `json:"settings"`
		Public     *bool           `json:"public,omitempty"`
		Image      *string         `json:"image"`
		Hoyo       string          `json:"hoyo,omitempty"`
		Order      string          `json:"order,omitempty"`
		HoyoType   HoyoType        `json:"hoyo_type"`
	}

	var decoded buildJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	*b = Build{
		ID:       decoded.ID,
		Name:     decoded.Name,
		AvatarID: decoded.AvatarID,
		Owner:    decoded.Owner,
		AvatarData: AvatarDataWrapper{
			Raw: append(json.RawMessage(nil), decoded.AvatarData...),
		},
		Live:     decoded.Live,
		Settings: decoded.Settings,
		Public:   decoded.Public,
		Image:    decoded.Image,
		Hoyo:     decoded.Hoyo,
		Order:    decoded.Order,
		HoyoType: decoded.HoyoType,
	}

	if len(decoded.AvatarData) == 0 || string(decoded.AvatarData) == "null" {
		return nil
	}

	switch decoded.HoyoType {
	case HoyoTypeGenshin:
		var avatar genshin.AvatarInfo
		if err := json.Unmarshal(decoded.AvatarData, &avatar); err != nil {
			return err
		}
		b.AvatarData.Genshin = &avatar
	case HoyoTypeHSR:
		var avatar hsr.AvatarDetail
		if err := json.Unmarshal(decoded.AvatarData, &avatar); err != nil {
			return err
		}
		b.AvatarData.HSR = &avatar
	case HoyoTypeZZZ:
		var avatar zzz.AvatarData
		if err := json.Unmarshal(decoded.AvatarData, &avatar); err != nil {
			return err
		}
		b.AvatarData.ZZZ = &avatar
	case HoyoTypeArknightsEndfield:
		var avatar endfield.CharData
		if err := json.Unmarshal(decoded.AvatarData, &avatar); err != nil {
			return err
		}
		b.AvatarData.Endfield = &avatar
	}

	return nil
}

// AvatarDataWrapper is a container struct that holds character data from different game clients.
// It is designed to support multiple games while maintaining a unified interface.
type AvatarDataWrapper struct {
	// Genshin holds character data specific to Genshin Impact.
	Genshin *genshin.AvatarInfo `json:"genshin,omitempty"`
	// HSR holds character data specific to Honkai: Star Rail.
	HSR *hsr.AvatarDetail `json:"hsr,omitempty"`
	// ZZZ holds character data specific to Zenless Zone Zero.
	ZZZ *zzz.AvatarData `json:"zzz,omitempty"`
	// Endfield holds character data specific to Arknights Endfield.
	Endfield *endfield.CharData `json:"endfield,omitempty"`
	// Raw contains the original JSON data for custom unmarshaling or debugging purposes.
	Raw json.RawMessage `json:"-"`
}

// UnmarshalJSON preserves raw avatar data when no build-level hoyo_type is available.
func (a *AvatarDataWrapper) UnmarshalJSON(data []byte) error {
	a.Genshin = nil
	a.HSR = nil
	a.ZZZ = nil
	a.Endfield = nil
	a.Raw = append(a.Raw[:0], data...)
	return nil
}

// MarshalJSON implements the json.Marshaler interface to provide custom JSON marshaling
// for the AvatarDataWrapper. This method serializes the appropriate game-specific
// avatar data based on which field is populated.
//
// The method checks each game-specific field in order of priority (Genshin -> HSR -> ZZZ -> Endfield),
// then falls back to the preserved raw JSON. If no data is present, it returns null.
func (a AvatarDataWrapper) MarshalJSON() ([]byte, error) {
	if a.Genshin != nil {
		return json.Marshal(a.Genshin)
	}

	if a.HSR != nil {
		return json.Marshal(a.HSR)
	}

	if a.ZZZ != nil {
		return json.Marshal(a.ZZZ)
	}

	if a.Endfield != nil {
		return json.Marshal(a.Endfield)
	}

	if len(a.Raw) > 0 {
		return a.Raw, nil
	}

	return []byte("null"), nil
}

// Hoyos is a map of Hoyo accounts and their metadata. The endpoint returns only
// verified and public accounts (users can hide accounts; unverified accounts are
// hidden by default). Each key is a unique identifier for a hoyo, which can be used
// in subsequent requests to retrieve information about the characters or builds of
// that game account.
type Hoyos map[string]Hoyo

// CacheTTL implements the core.Cacheable interface.
func (h Hoyos) CacheTTL() time.Duration {
	return 5 * time.Minute
}

// Owner represents an EnkaNetwork user profile associated with a game account.
type Owner = models.Owner

// PatreonProfile contains Patreon-related information for an Enka user.
type PatreonProfile = models.PatreonProfile

// Hoyo contains information about a specific Hoyo account.
type Hoyo struct {
	// User is the user information.
	User *Owner `json:"user,omitempty"`
	// UID is the UID of the game account, when the API provides it.
	UID *int `json:"uid,omitempty"`
	// UIDPublic indicates whether the UID is public.
	UIDPublic *bool `json:"uid_public,omitempty"`
	// Public indicates whether the Hoyo account is public.
	Public *bool `json:"public,omitempty"`
	// LivePublic indicates whether the live build is public.
	LivePublic *bool `json:"live_public,omitempty"`
	// Verified indicates whether the Hoyo account is verified.
	Verified *bool `json:"verified,omitempty"`
	// PlayerInfo is the player information for the account.
	PlayerInfo *models.PlayerInfo `json:"player_info,omitempty"`
	// Hash is the hash of the game account.
	Hash string `json:"hash,omitempty"`
	// Region is the region of the game account.
	Region string `json:"region,omitempty"`
	// Order is the order of the Hoyo account.
	Order string `json:"order"`
	// AvatarOrder is the order of the characters in the game account.
	AvatarOrder map[string]int `json:"avatar_order"`
	// HoyoType is the API game/account type discriminator.
	HoyoType HoyoType `json:"hoyo_type"`
	// LiveDataHash is the hash of the live data for the account.
	LiveDataHash int `json:"live_data_hash,omitempty"`
	// Raw contains the original API response.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown API fields.
	Extra map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON preserves unknown top-level fields for API drift tolerance.
func (h *Hoyo) UnmarshalJSON(data []byte) error {
	type Alias Hoyo
	var hoyo Alias
	if err := json.Unmarshal(data, &hoyo); err != nil {
		return err
	}

	raw, extra, err := core.PreserveUnknownJSON(
		data,
		"user",
		"uid",
		"uid_public",
		"public",
		"live_public",
		"verified",
		"player_info",
		"hash",
		"region",
		"order",
		"avatar_order",
		"hoyo_type",
		"live_data_hash",
	)
	if err != nil {
		return err
	}

	hoyo.Raw = raw
	hoyo.Extra = extra
	*h = Hoyo(hoyo)
	return nil
}

// MarshalJSON writes known fields plus unknown fields preserved in Extra.
func (h Hoyo) MarshalJSON() ([]byte, error) {
	type Alias Hoyo
	return core.MergeKnownExtraAndRawJSON(Alias(h), h.Extra, h.Raw)
}

// CacheTTL implements the core.Cacheable interface.
func (h Hoyo) CacheTTL() time.Duration {
	return 5 * time.Minute
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
}
