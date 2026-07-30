package zzz

import (
	"context"
	"fmt"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
)

// Client extends core.Client to provide ZZZ-specific functionality for player
// profile requests. It serves as the primary tool for interacting with the EnkaNetwork
// API in this package.
//
// The Client embeds an immutable core client configured at construction time.
// Provide HTTP, cache, User-Agent, retry, and base URL settings through Options
// when calling New. Once created, use the Client to call methods like GetProfile
// to fetch player data.
type Client struct {
	*core.Client // Embedded read-only shared client configuration
}

// Options configures a Zenless Zone Zero API client.
type Options = core.Options

// RetryOptions configures retry behavior for a Zenless Zone Zero API client.
type RetryOptions = core.RetryOptions

// New creates a new Zenless Zone Zero API client.
func New(options Options) *Client {
	c := core.NewClient(options)

	return &Client{
		Client: c,
	}
}

// GetProfile fetches the full player profile for the given UID using EnkaNetwork API.
// The profile includes detailed information about the player, such as their nickname,
// level, agents, equipment, etc., as defined in the Profile struct.
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
//   - uid: The player's UID, which must be a 9 or 10-digit string (e.g., "1301806568").
//
// Returns:
//   - *Profile: A pointer to the Profile struct if the request is successful.
//   - error: An error if the request fails.
//
// Possible errors include:
//   - ErrInvalidUIDFormat: If the UID is not a 9 or 10-digit number.
//   - ErrPlayerNotFound: If the player does not exist.
//   - ErrRateLimited: If the rate limit is exceeded after retries.
//   - ErrServerMaintenance: If the API is under maintenance.
//   - ErrServerError: For general server errors.
//   - ErrServiceUnavailable: If the API is completely unavailable.
//
// Example:
//
//	ctx := context.Background()
//	profile, err := client.GetProfile(ctx, "1301806568")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Player Nickname:", profile.PlayerInfo.SocialDetail.ProfileDetail.Nickname)
//	fmt.Println("World Level:", profile.PlayerInfo.SocialDetail.ProfileDetail.Level)
func (c *Client) GetProfile(ctx context.Context, uid string) (*Profile, error) {
	if !isValidUID(uid) {
		return nil, ErrInvalidUIDFormat
	}

	key := fmt.Sprintf("zzz_%s", uid)

	if c.Cache() != nil {
		if cached, ok := c.Cache().Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/zzz/uid/%s", c.BaseURL(), uid)
	profile, err := core.FetchWithRetry[Profile](ctx, c.Fetcher(), url)
	if err == nil && c.Cache() != nil {
		c.Cache().Set(key, profile, time.Duration(profile.TTL)*time.Second)
	}

	return profile, err
}

// isValidUID checks if the provided UID is a valid 9 or 10-digit number.
// ZZZ UID can only be 9 or 10 digits (e.g., "1301806568").
// This function is used internally to validate UIDs before making requests.
//
// Parameters:
//   - uid: The UID string to validate.
//
// Returns:
//   - true if the UID is a 9 or 10-digit number, false otherwise.
func isValidUID(uid string) bool {
	if len(uid) != 9 && len(uid) != 10 {
		return false
	}
	for _, r := range uid {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
