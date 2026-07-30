package enka

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/internal/core/errors"
)

// Client extends core.Client to provide Enka-specific functionality for user profile
// requests. It serves as the primary tool for interacting with the EnkaNetwork API in
// this package.
//
// The Client embeds an immutable core client configured at construction time.
// Provide HTTP, cache, User-Agent, retry, and base URL settings through Options
// when calling New. Once created, use the Client to call methods like
// GetUserProfile to fetch user data.
type Client struct {
	*core.Client // Embedded read-only shared client configuration
}

// Options configures an Enka API client.
type Options = core.Options

// RetryOptions configures retry behavior for an Enka API client.
type RetryOptions = core.RetryOptions

// New creates a new Enka API client.
func New(options Options) *Client {
	c := core.NewClient(options)

	return &Client{
		Client: c,
	}
}

// GetUserProfile fetches the Enka user profile for the given username.
//
// Enka allows users to create a profile and link multiple game accounts to it.
// Users can verify ownership of a game account by including a confirmation code in
// their signature — ensuring the account is associated with their profile.
//
// A user profile contains information about an Enka account, such as the username,
// bio, and avatar, as defined in the Owner struct.
//
// Unlike GetProfile, this method does not use a TTL for caching because user profiles
// do not include a TTL value. Instead, successful responses are cached for a fixed
// duration of 5 minutes to reduce API requests.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - username: The username of the EnkaNetwork user (must not be empty).
//
// Returns:
//   - *Owner: A pointer to the user's profile if successful.
//   - error: An error if the request fails.
//
// Possible errors include:
//   - ErrInvalidUsername: If the username is empty.
//   - ErrUserNotFound: If the user does not exist.
//   - Other errors for network issues or unexpected HTTP status codes.
//
// Example:
//
//	ctx := context.Background()
//	owner, err := client.GetUserProfile(ctx, "Algoinde")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Username:", owner.Username)
//	fmt.Println("Bio:", owner.Profile.Bio)
func (c *Client) GetUserProfile(ctx context.Context, username string) (*Owner, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}

	key := fmt.Sprintf("user_%s", username)
	requestURL := fmt.Sprintf("%s/profile/%s", c.BaseURL(), url.PathEscape(username))

	owner, err := core.FetchAndCache[Owner](ctx, c.Fetcher(), requestURL, key, c.Cache())
	if err != nil {
		if err == errors.ErrPlayerNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return owner, nil
}

// GetUserProfileHoyos fetches a list of “hoyos” — verified and public game accounts
// (e.g., Genshin Impact accounts) and their metadata.
//
// The API returns only accounts that are verified and public (users can hide accounts;
// unverified accounts are hidden by default). Each key in the response is a unique
// identifier for a hoyo, which can be used for subsequent queries to retrieve
// information about the characters or builds of that game account.
//
// The behavior is similar to GetUserProfile: it checks the cache first, makes an HTTP
// request if needed, retries on 429 errors, and caches the response for a fixed
// duration of 5 minutes.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - username: The username of the EnkaNetwork user (must not be empty).
//
// Returns:
//   - Hoyos: Map where the key is the hoyo hash and the value is the Hoyo struct.
//   - error: An error if the request fails.
//
// Possible errors include:
//   - ErrInvalidUsername: If the username is empty.
//   - ErrUserNotFound: If the user does not exist.
//
// Example:
//
//	ctx := context.Background()
//	hoyos, err := client.GetUserProfileHoyos(ctx, "Algoinde")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Hoyos:", hoyos)
func (c *Client) GetUserProfileHoyos(ctx context.Context, username string) (Hoyos, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}

	key := fmt.Sprintf("user_%s_hoyos", username)
	requestURL := fmt.Sprintf("%s/profile/%s/hoyos", c.BaseURL(), url.PathEscape(username))

	hoyos, err := core.FetchAndCache[Hoyos](ctx, c.Fetcher(), requestURL, key, c.Cache())
	if err != nil {
		if err == errors.ErrPlayerNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return *hoyos, nil
}

// GetUserProfileHoyo fetches information about a specific Hoyo account.
//
// The behavior is similar to GetUserProfile: it checks the cache first, makes an HTTP
// request if needed, retries on 429 errors, and caches the response for a fixed
// duration of 5 minutes.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - username: The username of the EnkaNetwork user (must not be empty).
//   - hoyo_hash: The hash of the hoyo (must not be empty).
//
// Returns:
//   - *Hoyo: A pointer to the hoyo data if successful.
//   - error: An error if the request fails.
//
// Possible errors include:
//   - ErrInvalidUsername: If the username is empty.
//   - ErrInvalidHoyoHash: If the hoyo hash is empty.
//   - ErrHoyoAccountNotFound: If the hoyo account does not exist.
//
// Example:
//
//	ctx := context.Background()
//	hoyo, err := client.GetUserProfileHoyo(ctx, "Algoinde", "4Wjv2e")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Hoyo:", hoyo)
func (c *Client) GetUserProfileHoyo(ctx context.Context, username string, hoyo_hash string) (*Hoyo, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}

	if hoyo_hash == "" {
		return nil, ErrInvalidHoyoHash
	}

	key := fmt.Sprintf("user_%s_hoyos_%s", username, hoyo_hash)
	requestURL := fmt.Sprintf("%s/profile/%s/hoyos/%s", c.BaseURL(), url.PathEscape(username), url.PathEscape(hoyo_hash))

	hoyo, err := core.FetchAndCache[Hoyo](ctx, c.Fetcher(), requestURL, key, c.Cache())
	if err != nil {
		if err == errors.ErrPlayerNotFound {
			return nil, ErrHoyoAccountNotFound
		}
		return nil, err
	}

	return hoyo, nil
}

// GetUserProfileHoyoBuilds fetches character builds for a specific Hoyo account.
//
// The response is a map where the key is the character's avatarId, and the value is
// a slice of builds for that character, returned in random order. Each build includes
// an order field that can be used for sorting during display.
//
// If a build has a live: true field, it indicates a build retrieved from the showcase
// when the “update” button was clicked, rather than a saved build. When updating,
// all old live builds are deleted, and new ones are created. Updates are user-initiated,
// so this data may not be up to date.
//
// The behavior is similar to GetUserProfile: it checks the cache first, makes an HTTP
// request if needed, retries on 429 errors, and caches the response for a fixed
// duration of 5 minutes.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - username: The username of the EnkaNetwork user (must not be empty).
//   - hoyo_hash: The hash of the hoyo (must not be empty).
//
// Returns:
//   - AvatarBuildsMap: A map where the key is the avatarID and the value is a slice of builds for that character.
//   - error: An error if the request fails, such as ErrInvalidUsername or ErrHoyoAccountBuildsNotFound.
//
// Example:
//
//	ctx := context.Background()
//	avatarBuilds, err := client.GetUserProfileHoyoBuilds(ctx, "Algoinde", "4Wjv2e")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("avatarBuilds:", avatarBuilds)
func (c *Client) GetUserProfileHoyoBuilds(ctx context.Context, username string, hoyo_hash string) (AvatarBuildsMap, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}

	if hoyo_hash == "" {
		return nil, ErrInvalidHoyoHash
	}

	key := fmt.Sprintf("user_%s_hoyos_%s_builds", username, hoyo_hash)
	requestURL := fmt.Sprintf("%s/profile/%s/hoyos/%s/builds", c.BaseURL(), url.PathEscape(username), url.PathEscape(hoyo_hash))

	builds, err := core.FetchAndCache[AvatarBuildsMap](ctx, c.Fetcher(), requestURL, key, c.Cache())
	if err != nil {
		if err == errors.ErrPlayerNotFound {
			return nil, ErrHoyoAccountBuildsNotFound
		}
		return nil, err
	}

	return *builds, nil
}
