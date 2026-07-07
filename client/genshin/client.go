package genshin

import (
	"context"
	"fmt"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/internal/core/fetcher"
)

// Client extends core.Client to provide Genshin-specific functionality for player
// profile requests. It serves as the primary tool for interacting with the EnkaNetwork
// API in this package.
//
// The Client embeds an immutable core client configured at construction time.
// Provide HTTP, cache, User-Agent, retry, and base URL settings through Options
// when calling New. Once created, use the Client to call methods like GetProfile
// to fetch player data.
type Client struct {
	*core.Client // Embedded read-only shared client configuration
	fetcher      *fetcher.Fetcher[Profile]
}

// Options configures a Genshin Impact API client.
type Options = core.Options

// RetryOptions configures retry behavior for a Genshin Impact API client.
type RetryOptions = core.RetryOptions

// New creates a new Genshin Impact API client.
func New(options Options) *Client {
	c := core.NewClient(options)

	return &Client{
		Client:  c,
		fetcher: fetcher.NewFetcher[Profile](c.HTTPClient(), c.UserAgent(), c.Retry()),
	}
}

// GetProfile fetches the full player profile for the given UID using EnkaNetwork API.
// The response will contain PlayerInfo and AvatarInfoList. PlayerInfo contains basic
// information about the game account. AvatarInfoList contains detailed information for
// each character in the showcase. If AvatarInfoList is missing, it means that the
// account's showcase is either hidden by the player or there are no characters there.
//
// This method first checks if the profile is available in the cache (if a cache is
// provided). If not, it sends an HTTP GET request to the API. If the API returns a
// 429 (Too Many Requests) status, the client retries according to Options.Retry.
// By default, it makes up to 3 attempts and waits for the duration specified in
// the Retry-After header or 5 seconds.
//
// If the request is successful, the profile is cached locally using the ttl value
// returned by the API, which indicates how long the data remains valid before the
// API queries the game again. Caching helps reduce the number of requests and
// respects the API's rate limits.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation. For
//     example, you can use context.WithTimeout to set a maximum duration for the request.
//   - uid: The player's UID, which must be a 9-digit string (e.g., "618285856").
//
// Returns:
//   - *Profile: A pointer to the Profile struct if the request is successful.
//   - error: An error if the request fails.
//
// Possible errors include:
//   - ErrInvalidUIDFormat: If the UID is not a 9-digit number.
//   - ErrPlayerNotFound: If the player does not exist.
//   - ErrRateLimited: If the rate limit is exceeded after retries.
//   - ErrServerMaintenance: If the API is under maintenance.
//   - ErrServerError: For general server errors.
//   - ErrServiceUnavailable: If the API is completely unavailable.
//
// Example:
//
//	ctx := context.Background()
//	profile, err := client.GetProfile(ctx, "618285856")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Player Nickname:", profile.PlayerInfo.Nickname)
//	fmt.Println("World Level:", profile.PlayerInfo.WorldLevel)
func (c *Client) GetProfile(ctx context.Context, uid string) (*Profile, error) {
	if !core.IsValidUID(uid) {
		return nil, ErrInvalidUIDFormat
	}

	key := fmt.Sprintf("genshin_%s", uid)

	if c.Cache() != nil {
		if cached, ok := c.Cache().Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/uid/%s", c.BaseURL(), uid)

	profile, err := c.fetcher.FetchWithRetry(ctx, url)
	if err == nil && c.Cache() != nil {
		c.Cache().Set(key, profile, time.Duration(profile.TTL)*time.Second)
	}

	return profile, err
}

// GetPlayerInfo fetches limited player profile information for the given UID.
// GetProfile always makes an additional request to obtain AvatarInfoList.
// If you only need PlayerInfo, use GetPlayerInfo — it works faster and has fewer rate limits.
//
// The behavior is similar to GetProfile: it checks the cache first, makes an HTTP
// request if needed, retries on 429 errors, and caches the response using the ttl
// value from the API.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - uid: The player's UID, which must be a 9-digit string.
//
// Returns:
//   - *Profile: A pointer to the Profile struct (AvatarInfoList is always empty slice) if the request is successful.
//   - error: An error if the request fails.
//
// Possible errors include:
//   - ErrInvalidUIDFormat: If the UID is not a 9-digit number.
//   - ErrPlayerNotFound: If the player does not exist.
//   - ErrRateLimited: If the rate limit is exceeded after retries.
//   - ErrServerMaintenance: If the API is under maintenance.
//   - ErrServerError: For general server errors.
//   - ErrServiceUnavailable: If the API is completely unavailable.
//
// Example:
//
//	ctx := context.Background()
//	profile, err := client.GetPlayerInfo(ctx, "618285856")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Player Nickname:", profile.PlayerInfo.Nickname)
func (c *Client) GetPlayerInfo(ctx context.Context, uid string) (*Profile, error) {
	if !core.IsValidUID(uid) {
		return nil, ErrInvalidUIDFormat
	}

	key := "genshin_" + uid + "_info"

	if c.Cache() != nil {
		if cached, ok := c.Cache().Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/uid/%s?info", c.BaseURL(), uid)

	profile, err := c.fetcher.FetchWithRetry(ctx, url)
	if err == nil && c.Cache() != nil {
		c.Cache().Set(key, profile, time.Duration(profile.TTL)*time.Second)
	}

	return profile, err
}
