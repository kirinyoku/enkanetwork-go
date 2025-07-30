package core

import (
	"net/http"
	"time"
)

// BaseURL is the root URL for the EnkaNetwork API, used as the starting point for all
// API requests. Each game (Genshin Impact, Honkai: Star Rail, Zenless Zone Zero) builds
// specific endpoints by adding paths to this URL.
const (
	BaseURL = "https://enka.network/api"
)

// Client represents an EnkaNetwork API client used to make requests to the API.
// It holds an HTTP client for sending requests, an optional cache for storing
// responses, and a User-Agent string to identify the client in API requests.
//
// Fields:
//   - HTTPClient: The HTTP client used for making requests. You can provide a custom
//     client with specific settings, like timeouts or proxies.
//   - Cache: An optional cache implementation to store API responses locally.
//   - UserAgent: A string sent in the User-Agent header of every request to identify
//     your application.
type Client struct {
	HTTPClient *http.Client // HTTP client for making requests
	Cache      Cache        // Optional cache for storing API responses
	UserAgent  string       // User-Agent string for HTTP requests
}

// NewClient creates and configures a new Client instance for making requests to the
// EnkaNetwork API. This function is used internally by game-specific client (e.g.,
// genshin.NewClient, hsr.NewClient) to set up the shared functionality needed for API
// communication. Users of the library don’t call this function directly; instead, they
// use the NewClient function provided by the game-specific package they’re working with,
// such as client/genshin.
//
// The function takes three parameters to customize the client:
//   - httpClient: An optional HTTP client for sending requests. If you provide nil, the
//     function creates a default HTTP client with a 10-second timeout, which means
//     requests will fail if the API doesn’t respond within 10 seconds. You can pass a
//     custom HTTP client with different settings, like a 30-second timeout or proxy
//     support, if needed.
//   - cache: An optional cache (implementing the Cache interface) for storing API
//     responses. If you provide nil, no caching will be used, and every request will go
//     directly to the API. Caching is recommended to reduce the number of requests and
//     stay within the API’s rate limits.
//   - userAgent: A string that identifies your application in API requests. If you provide
//     an empty string, the function sets a default User-Agent of
//     "enka-network-go-client/1.0". It’s a good idea to use a unique User-Agent, like
//     "my-game-app/1.0", to help the API team know who’s using their service.
//
// The function returns a pointer to a fully configured Client, ready to be used by
// game-specific client to make API requests.
func NewClient(httpClient *http.Client, cache Cache, userAgent string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	if userAgent == "" {
		userAgent = "enka-network-go-client/1.0"
	}
	return &Client{
		HTTPClient: httpClient,
		Cache:      cache,
		UserAgent:  userAgent,
	}
}
