package endfield

import (
	"context"
	"fmt"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
)

// Client extends core.Client to provide Arknights Endfield-specific functionality for player
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

// Options configures an Arknights Endfield API client.
type Options = core.Options

// RetryOptions configures retry behavior for an Arknights Endfield API client.
type RetryOptions = core.RetryOptions

// New creates a new Arknights Endfield API client.
func New(options Options) *Client {
	c := core.NewClient(options)

	return &Client{
		Client: c,
	}
}

// GetProfile fetches the player profile for the given UID using EnkaNetwork API.
// Since the exact data structure for Endfield is undocumented, the Profile struct
// primarily acts as a wrapper around the raw JSON, placing unknown fields into Extra.
//
// This method first checks if the profile is available in the cache (if a cache is
// provided). If not, it sends an HTTP GET request to the API. If the API returns a
// 429 (Too Many Requests) status, the client retries according to Options.Retry.
//
// If the request is successful, the profile is cached locally using the ttl value
// returned by the API.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - uid: The player's UID.
//
// Returns:
//   - *Profile: A pointer to the Profile struct if the request is successful.
//   - error: An error if the request fails.
//
// Possible errors include:
//   - ErrInvalidUIDFormat: If the UID is empty or contains non-digits.
//   - ErrPlayerNotFound: If the player does not exist.
//   - ErrRateLimited: If the rate limit is exceeded after retries.
//   - ErrServerMaintenance: If the API is under maintenance.
//   - ErrServerError: For general server errors.
//   - ErrServiceUnavailable: If the API is completely unavailable.
func (c *Client) GetProfile(ctx context.Context, uid string) (*Profile, error) {
	if !isValidUID(uid) {
		return nil, ErrInvalidUIDFormat
	}

	key := fmt.Sprintf("endfield_%s", uid)
	url := fmt.Sprintf("%s/ef/uid/%s", c.BaseURL(), uid)

	return core.FetchAndCache[Profile](ctx, c.Fetcher(), url, key, c.Cache())
}

// isValidUID checks if the provided UID is valid.
// For Arknights Endfield, we simply verify that it consists only of digits
// and is not empty, as the exact length is not strictly standardized across regions.
//
// Parameters:
//   - uid: The UID string to validate.
//
// Returns:
//   - true if the UID contains only digits and is not empty, false otherwise.
func isValidUID(uid string) bool {
	if len(uid) == 0 {
		return false
	}
	for _, r := range uid {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
