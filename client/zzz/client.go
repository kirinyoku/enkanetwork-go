package zzz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/internal/core/fetcher"
)

// Client extends core.Client to provide ZZZ-specific functionality for player
// profile requests. It serves as the primary tool for interacting with the EnkaNetwork
// API in this package.
//
// The Client struct embeds core.Client, inheriting shared features, including:
//   - An HTTP client for sending API requests.
//   - An optional cache to store responses and reduce API calls.
//   - A User-Agent string to identify the application in requests.
//
// Create a Client using the NewClient function, which allows customization of these
// settings. Once created, use the Client to call methods like GetProfile to fetch
// player data.
type Client struct {
	*core.Client // Embeds core.Client for shared HTTP and caching functionality
	fetcher      *fetcher.Fetcher[Profile]
}

// NewClient creates a new Zenless Zone Zero API client for making requests.
//
// This function allows you to customize the client by providing your own HTTP client,
// cache implementation, and User-Agent string. If you don't provide these, default
// values are used: a standard HTTP client with a 10-second timeout, no cache, and
// a default User-Agent of "enkanetwork-go-client/1.0".
//
// Parameters:
//   - httpClient: An optional *http.Client for making HTTP requests. If nil, a default
//     client with a 10-second timeout is used.
//   - cache: An optional Cache implementation for storing responses. If nil, caching
//     is disabled.
//   - userAgent: A string to set as the User-Agent header in requests. If empty, the
//     default "enkanetwork-go-client/1.0" is used. It's recommended to set a unique
//     User-Agent to identify your application, such as "my-app/1.0".
//
// Returns:
//   - A pointer to a new ZZZ-specific Client instance ready to make API requests.
//
// Example:
//
//	// Create a client with default settings
//	client := zzz.NewClient(nil, nil, "my-app/1.0")
//	// Create a client with a custom HTTP client
//	customClient := &http.Client{Timeout: 20 * time.Second}
//	client := zzz.NewClient(customClient, nil, "my-app/1.0")
func NewClient(httpClient *http.Client, cache core.Cache, userAgent string) *Client {
	c := core.NewClient(httpClient, cache, userAgent)

	return &Client{
		Client:  c,
		fetcher: fetcher.NewFetcher[Profile](c.HTTPClient, c.UserAgent),
	}
}

// GetProfile fetches the full player profile for the given UID using EnkaNetwork API.
// The profile includes detailed information about the player, such as their nickname,
// level, agents, equipment, etc., as defined in the Profile struct.
//
// This method first checks if the profile is available in the cache (if a cache is
// provided). If not, it sends an HTTP GET request to the API. If the API returns a
// 429 (Too Many Requests) status, the client will retry up to 3 times, waiting for
// the duration specified in the Retry-After header or 5 seconds by default.
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

	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/zzz/uid/%s", core.BaseURL, uid)
	profile, err := c.fetcher.FetchWithRetry(ctx, url)
	if err == nil && c.Cache != nil {
		c.Cache.Set(key, profile, time.Duration(profile.TTL)*time.Second)
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
