// Package enka provides functionality for interacting with the EnkaNetwork API to
// fetch user profiles from the EnkaNetwork platform.
//
// The package ensures efficient requests through caching and provides clear, easy-to-handle
// errors for seamless integration into applications.
package enka

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/common"
)

// Client extends common.Client to provide Enka-specific functionality for user profile
// requests. It serves as the primary tool for interacting with the EnkaNetwork API in
// this package.
//
// The Client struct embeds common.Client, inheriting shared features, including:
// - An HTTP client for sending API requests.
// - An optional cache to store responses and reduce API calls.
// - A User-Agent string to identify the application in requests.
//
// Create a Client using the NewClient function, which allows customization of these
// settings. Once created, use the Client to call methods like GetUserProfile to fetch
// user data.
type Client struct {
	*common.Client // Embeds common.Client for shared HTTP and caching functionality
}

// NewClient creates a new Enka API client for making requests.
//
// This function allows customization of the client with a custom HTTP client, cache
// implementation, and User-Agent string. If not provided, default values are used:
// a standard HTTP client with a 10-second timeout, no cache, and a default User-Agent
// of "enkanetwork-go-client/1.0".
//
// Parameters:
//   - httpClient: An optional *http.Client for making HTTP requests. If nil, a default
//     client with a 10-second timeout is used.
//   - cache: An optional Cache implementation for storing responses. If nil, caching
//     is disabled.
//   - userAgent: A string to set as the User-Agent header in requests. If empty, the
//     default "enkanetwork-go-client/1.0" is used. A unique User-Agent, such as
//     "my-app/1.0", is recommended to identify the application.
//
// Returns:
//   - A pointer to a new Enka-specific Client instance ready to make API requests.
//
// Example:
//
//	// Create a client with default settings
//	client := enka.NewClient(nil, nil, "my-app/1.0")
//	// Create a client with a custom HTTP client
//	customClient := &http.Client{Timeout: 20 * time.Second}
//	client := enka.NewClient(customClient, nil, "my-app/1.0")
func NewClient(httpClient *http.Client, cache common.Cache, userAgent string) *Client {
	return &Client{
		Client: common.NewClient(httpClient, cache, userAgent),
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

	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			if owner, ok := cached.(*Owner); ok {
				return owner, nil
			}
		}
	}

	url := fmt.Sprintf("%s/profile/%s", common.BaseURL, username)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 404:
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}

	var owner Owner
	err = json.NewDecoder(resp.Body).Decode(&owner)
	if err != nil {
		return nil, err
	}

	if c.Cache != nil {
		c.Cache.Set(key, &owner, 5*time.Minute)
	}

	return &owner, nil
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

	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			if hoyos, ok := cached.(Hoyos); ok {
				return hoyos, nil
			}
		}
	}

	url := fmt.Sprintf("%s/profile/%s/hoyos/", common.BaseURL, username)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 404:
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}

	var hoyos Hoyos
	err = json.NewDecoder(resp.Body).Decode(&hoyos)
	if err != nil {
		return nil, err
	}

	if c.Cache != nil {
		c.Cache.Set(key, &hoyos, 5*time.Minute)
	}

	return hoyos, nil
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

	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			if hoyo, ok := cached.(*Hoyo); ok {
				return hoyo, nil
			}
		}
	}

	url := fmt.Sprintf("%s/profile/%s/hoyos/%s/?format=json", common.BaseURL, username, hoyo_hash)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 404:
			return nil, ErrHoyoAccountNotFound
		default:
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}

	var hoyo Hoyo
	err = json.NewDecoder(resp.Body).Decode(&hoyo)
	if err != nil {
		return nil, err
	}

	if c.Cache != nil {
		c.Cache.Set(key, &hoyo, 5*time.Minute)
	}

	return &hoyo, nil
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

	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			if builds, ok := cached.(AvatarBuildsMap); ok {
				return builds, nil
			}
		}
	}

	url := fmt.Sprintf("%s/profile/%s/hoyos/%s/builds/", common.BaseURL, username, hoyo_hash)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 404:
			return nil, ErrHoyoAccountBuildsNotFound
		default:
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
	}

	var builds AvatarBuildsMap
	err = json.NewDecoder(resp.Body).Decode(&builds)
	if err != nil {
		return nil, err
	}

	if c.Cache != nil {
		c.Cache.Set(key, builds, 5*time.Minute)
	}

	return builds, nil
}
