package core

import (
	"net/http"
	"strings"
	"time"
)

// DefaultBaseURL is the root URL for the EnkaNetwork API.
const (
	DefaultBaseURL     = "https://enka.network/api"
	DefaultMaxAttempts = 3
	DefaultRetryDelay  = 5 * time.Second
)

// Options configures a shared EnkaNetwork API client.
type Options struct {
	HTTPClient *http.Client
	Cache      Cache
	UserAgent  string
	// Retry configures retry behavior for retryable HTTP responses.
	// Leave nil to use the default retry behavior.
	Retry *RetryOptions
	// BaseURL overrides the EnkaNetwork API base URL.
	// Leave empty for the official API. This is mainly intended for tests,
	// mocks, proxies, and advanced users. Do not set it from untrusted input.
	BaseURL string
}

// RetryOptions configures retry behavior for transient API failures.
type RetryOptions struct {
	// MaxAttempts is the total number of request attempts, including the first
	// request. Leave zero to use the default. Set to 1 to disable retries.
	MaxAttempts int
	// Delay is the fallback delay between attempts when Retry-After is not used.
	// Leave zero to use the default delay.
	Delay time.Duration
}

// Client represents an EnkaNetwork API client used to make requests to the API.
// It holds an HTTP client for sending requests, an optional cache for storing
// responses, and a User-Agent string to identify the client in API requests.
type Client struct {
	httpClient *http.Client
	cache      Cache
	userAgent  string
	baseURL    string
	retry      RetryOptions
	fetcher    *Fetcher
}

// NewClient creates and configures a shared Client instance.
func NewClient(options Options) *Client {
	httpClient := options.HTTPClient
	cache := options.Cache
	userAgent := options.UserAgent
	baseURL := strings.TrimRight(options.BaseURL, "/")

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	if userAgent == "" {
		userAgent = "enka-network-go-client/1.0"
	}

	retryOptions := NormalizeRetryOptions(options.Retry)

	return &Client{
		httpClient: httpClient,
		cache:      cache,
		userAgent:  userAgent,
		baseURL:    baseURL,
		retry:      retryOptions,
		fetcher:    NewFetcher(httpClient, userAgent, retryOptions),
	}
}

// HTTPClient returns the HTTP client configured at construction time.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// Cache returns the cache configured at construction time, if any.
func (c *Client) Cache() Cache {
	return c.cache
}

// UserAgent returns the User-Agent sent with API requests.
func (c *Client) UserAgent() string {
	return c.userAgent
}

// BaseURL returns the base API URL configured at construction time.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Retry returns the retry behavior configured at construction time.
func (c *Client) Retry() RetryOptions {
	return c.retry
}

// Fetcher returns the generic HTTP fetcher used by the client.
func (c *Client) Fetcher() *Fetcher {
	return c.fetcher
}

// NormalizeRetryOptions applies default retry settings where options are unset.
func NormalizeRetryOptions(options *RetryOptions) RetryOptions {
	retry := RetryOptions{
		MaxAttempts: DefaultMaxAttempts,
		Delay:       DefaultRetryDelay,
	}
	if options == nil {
		return retry
	}
	if options.MaxAttempts > 0 {
		retry.MaxAttempts = options.MaxAttempts
	}
	if options.Delay > 0 {
		retry.Delay = options.Delay
	}
	return retry
}
