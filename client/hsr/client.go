package hsr

import (
	"context"
	"fmt"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
)

// Client extends core.Client to provide HSR-specific functionality for player
// profile requests. It serves as the primary tool for interacting with the EnkaNetwork
// API in this package.
//
// The Client embeds an immutable core client configured at construction time.
// Provide HTTP, cache, User-Agent, retry, and base URL settings through Options
// when calling New. Once created, use the Client to call GetProfile to fetch
// player data.
type Client struct {
	*core.Client // Embedded read-only shared client configuration
}

// Options configures an HSR API client.
type Options = core.Options

// RetryOptions configures retry behavior for an HSR API client.
type RetryOptions = core.RetryOptions

// New creates a new HSR API client.
func New(options Options) *Client {
	c := core.NewClient(options)

	return &Client{
		Client: c,
	}
}

// GetProfile fetches the full player profile for the given UID using EnkaNetwork API.
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
//   - uid: The player's UID, which must be a 9-digit string (e.g., "800579959").
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
//	profile, err := client.GetProfile(ctx, "800579959")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//
// fmt.Println("Player Nickname:", profile.DetailInfo.Nickname)
// fmt.Println("World Level:", profile.DetailInfo.WorldLevel)
func (c *Client) GetProfile(ctx context.Context, uid string) (*Profile, error) {
	if !core.IsValidUID(uid) {
		return nil, ErrInvalidUIDFormat
	}

	key := fmt.Sprintf("hsr_%s", uid)
	url := fmt.Sprintf("%s/hsr/uid/%s", c.BaseURL(), uid)

	return core.FetchAndCache[Profile](ctx, c.Fetcher(), url, key, c.Cache())
}
