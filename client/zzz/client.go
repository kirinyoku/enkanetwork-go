// Package zzz provides a simple client for interacting with the EnkaNetwork API
// for fetching Zenless Zone Zero player game profiles. It is designed to be user-friendly.
//
// The API provides data about players, such as their nickname, level, agents,
// equipment, etc. This client simplifies access to this data with methods like
// GetProfile and GetPlayerInfo.
//
// To use this package:
//  1. Create a Client instance using NewClient, optionally providing a custom HTTP
//     client, cache, and User-Agent string.
//  2. Call methods like GetProfile to fetch player data.
//  3. Handle errors returned by the methods, which provide clear information about
//     issues such as invalid UID, player not found, or rate limit exceeded.
//  4. Use a context to control request timeouts or cancellation as needed.
//
// Important rules for using the EnkaNetwork API:
//   - Avoid mass requests or iterating through UIDs, as this may overload the API
//     and result in rate limiting (HTTP 429).
//   - Set a User-Agent header to identify the application, aiding the API provider
//     in troubleshooting issues.
//   - Cache responses locally using the TTL value returned by the API to minimize
//     unnecessary requests, as cached responses still count toward rate limits.
//   - If a rate limit (429) is encountered, the client retries up to three times,
//     but code should be optimized to reduce requests.
//
// For more details, see the EnkaNetwork API documentation:
//   - https://api.enka.network/#/api
//   - https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md
package zzz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/common"
)

// Client extends common.Client to provide ZZZ-specific functionality for player
// profile requests. It serves as the primary tool for interacting with the EnkaNetwork
// API in this package.
//
// The Client struct embeds common.Client, inheriting shared features, including:
//   - An HTTP client for sending API requests.
//   - An optional cache to store responses and reduce API calls.
//   - A User-Agent string to identify the application in requests.
//
// Create a Client using the NewClient function, which allows customization of these
// settings. Once created, use the Client to call methods like GetProfile to fetch
// player data.
type Client struct {
	*common.Client // Embeds common.Client for shared HTTP and caching functionality
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
func NewClient(httpClient *http.Client, cache common.Cache, userAgent string) *Client {
	return &Client{
		Client: common.NewClient(httpClient, cache, userAgent),
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
//   - *Profile: A pointer to the player's profile if the request is successful.
//   - error: An error if the request fails. Possible errors include:
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
//	fmt.Println("Player Nickname:", profile.PlayerInfo.Nickname)
//	fmt.Println("World Level:", profile.PlayerInfo.WorldLevel)
func (c *Client) GetProfile(ctx context.Context, uid string) (*Profile, error) {
	if !isValidUID(uid) {
		return nil, common.ErrInvalidUIDFormat
	}

	if c.Cache != nil {
		key := "zzz_" + uid
		if cached, ok := c.Cache.Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/zzz/uid/%s/", common.BaseURL, uid)
	profile, err := c.fetchProfileWithRetry(ctx, url)
	if err == nil && c.Cache != nil {
		key := "zzz_" + uid
		c.Cache.Set(key, profile, time.Duration(profile.TTL)*time.Second)
	}
	return profile, err
}

// GetPlayerInfo fetches limited player profile information for the given UID.
// Unlike GetProfile, this method uses the "?info" query parameter to retrieve only basic information
// about the player (without detailed information about the agents in the showcase),
// which can be faster and use fewer API resources.
//
// The behavior is similar to GetProfile: it checks the cache first, makes an HTTP
// request if needed, retries on 429 errors, and caches the response using the ttl
// value from the API.
//
// Parameters:
//   - ctx: A context.Context to control the request's timeout or cancellation.
//   - uid: The player's UID, which must be a 9 or 10-digit string.
//
// Returns:
//   - *Profile: A pointer to the player's limited profile if successful.
//   - error: An error if the request fails (same possible errors as GetProfile).
//
// Example:
//
//	ctx := context.Background()
//	profile, err := client.GetPlayerInfo(ctx, "1301806568")
//	if err != nil {
//	    fmt.Println("Error:", err)
//	    return
//	}
//	fmt.Println("Player Nickname:", profile.PlayerInfo.Nickname)
func (c *Client) GetPlayerInfo(ctx context.Context, uid string) (*Profile, error) {
	if !isValidUID(uid) {
		return nil, common.ErrInvalidUIDFormat
	}

	if c.Cache != nil {
		key := "zzz_" + uid + "_info"
		if cached, ok := c.Cache.Get(key); ok {
			if profile, ok := cached.(*Profile); ok {
				return profile, nil
			}
		}
	}

	url := fmt.Sprintf("%s/zzz/uid/%s/?info", common.BaseURL, uid)
	profile, err := c.fetchProfileWithRetry(ctx, url)
	if err == nil && c.Cache != nil {
		key := "zzz_" + uid + "_info"
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
func (c *Client) fetchProfileWithRetry(ctx context.Context, url string) (*Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	var retries int
	const maxRetries = 3

	for {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			var profile Profile
			if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
				return nil, fmt.Errorf("failed to decode response: %w", err)
			}
			return &profile, nil

		case http.StatusTooManyRequests:
			if retries >= maxRetries {
				return nil, common.ErrRateLimited
			}
			retries++

			retryAfter := resp.Header.Get("Retry-After")
			var waitTime time.Duration
			if retryAfter != "" {
				waitTime, _ = time.ParseDuration(retryAfter + "s")
			} else {
				waitTime = 5 * time.Second
			}

			time.Sleep(waitTime)
			continue

		case http.StatusBadRequest:
			return nil, common.ErrInvalidUIDFormat
		case http.StatusNotFound:
			return nil, common.ErrPlayerNotFound
		case http.StatusFailedDependency:
			return nil, common.ErrServerMaintenance
		case http.StatusInternalServerError:
			return nil, common.ErrServerError
		case http.StatusServiceUnavailable:
			return nil, common.ErrServiceUnavailable
		default:
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}
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
