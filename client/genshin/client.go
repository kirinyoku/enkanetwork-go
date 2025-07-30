package genshin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
)

// Client extends core.Client to provide Genshin-specific functionality for player
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
}

// NewClient creates a new Genshin Impact API client for making requests.
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
//   - A pointer to a new Genshin-specific Client instance ready to make API requests.
//
// Example:
//
//	// Create a client with default settings
//	client := genshin.NewClient(nil, nil, "my-app/1.0")
//	// Create a client with a custom HTTP client
//	customClient := &http.Client{Timeout: 20 * time.Second}
//	client := genshin.NewClient(customClient, nil, "my-app/1.0")
func NewClient(httpClient *http.Client, cache core.Cache, userAgent string) *Client {
	return &Client{
		Client: core.NewClient(httpClient, cache, userAgent),
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

	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/uid/%s", core.BaseURL, uid)
	profile, err := c.fetchProfileWithRetry(ctx, url)
	if err == nil && c.Cache != nil {
		c.Cache.Set(key, profile, time.Duration(profile.TTL)*time.Second)
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

	if c.Cache != nil {
		if cached, ok := c.Cache.Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/uid/%s?info", core.BaseURL, uid)
	profile, err := c.fetchProfileWithRetry(ctx, url)
	if err == nil && c.Cache != nil {
		c.Cache.Set(key, profile, time.Duration(profile.TTL)*time.Second)
	}

	return profile, err
}

// fetchProfileWithRetry is an internal helper function that fetches a player profile
// from the given URL with retry logic for handling rate limits (HTTP 429).
// It is used by GetProfile and GetPlayerInfo to make HTTP requests and process responses.
//
// The function:
//  1. Creates an HTTP request with the provided context and User-Agent header.
//  2. Sends the request and checks the response status code.
//  3. If the status is 200 (OK), decodes the response into a Profile struct.
//  4. If the status is 429 (Too Many Requests), retries up to 3 times, waiting for
//     the duration specified in the Retry-After header or 5 seconds by default.
//  5. For other status codes (400, 404, 424, 500, 503), returns the appropriate error.
//  6. If all retries fail due to rate limiting, returns an ErrRateLimited error.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - url: The URL to fetch the profile from.
//
// Returns:
//   - *Profile: A pointer to the player's profile if successful.
//   - error: An error if the request fails or retries are exhausted.
//
// Error handling includes specific error types for common HTTP status codes:
//   - 400: Invalid UID format
//   - 404: Player not found
//   - 424: Server under maintenance
//   - 429: Rate limited (handled automatically with retries)
//   - 500: Internal server error
//   - 503: Service unavailable
func (c *Client) fetchProfileWithRetry(ctx context.Context, url string) (*Profile, error) {
	const maxRetries = 3
	var profile Profile

	for attempt := 0; attempt < maxRetries; attempt++ {
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

		if resp.StatusCode == http.StatusOK {
			err = json.NewDecoder(resp.Body).Decode(&profile)
			if err != nil {
				return nil, fmt.Errorf("failed to decode profile: %w", err)
			}
			return &profile, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			var delay time.Duration

			if retryAfter != "" {
				if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
					delay = seconds
				} else {
					delay = 5 * time.Second
				}
			} else {
				delay = 5 * time.Second
			}
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		} else {
			switch resp.StatusCode {
			case 400:
				return nil, ErrInvalidUIDFormat
			case 404:
				return nil, ErrPlayerNotFound
			case 424:
				return nil, ErrServerMaintenance
			case 500:
				return nil, ErrServerError
			case 503:
				return nil, ErrServiceUnavailable
			default:
				return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
			}
		}
	}

	return nil, fmt.Errorf("rate limited: %w", ErrRateLimited)
}
