# EnkaNetwork Go API Wrapper

A lightweight, type-safe Go wrapper for the [EnkaNetwork API](https://api.enka.network/#/api), providing seamless access to player data for HoYoverse games, including:

- **Genshin Impact**
- **Honkai: Star Rail**
- **Zenless Zone Zero**

[![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/enkanetwork-go.svg)](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Multi-game Support**: Unified interface for Genshin Impact, Honkai: Star Rail, and Zenless Zone Zero.
- **Type Safety**: Strongly typed structs for reliable data handling.
- **Context Integration**: Full support for `context.Context` to manage request cancellation and timeouts.
- **Caching**: Built-in caching to reduce API calls and improve performance.
- **Robust Error Handling**: Detailed error types for common API scenarios.
- **Concurrency Safety**: Thread-safe for use in concurrent Go applications.

## Prerequisites

Before using this library, ensure the following:

- **Go Version**: Go 1.18 or later.
- **API Familiarity**: Review the official EnkaNetwork API documentation:
  - [EnkaNetwork API Documentation](https://api.enka.network/#/api)
  - [Genshin Impact API Details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md)
  - [Zenless Zone Zero API Details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md)
  - [EnkaNetwork Status](https://status.enka.network/) (check for API availability)

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Clients](#clients)
  - [Genshin Impact](#genshin-impact)
  - [Honkai: Star Rail](#honkai-star-rail)
  - [Zenless Zone Zero](#zenless-zone-zero)
  - [EnkaNetwork](#enka-network)
- [Advanced Features](#advanced-features)
  - [Caching](#caching)
  - [Error Handling](#error-handling)
  - [Context Support](#context-support)
- [API Usage Guidelines](#api-usage-guidelines)
- [Contributing](#contributing)
- [License](#license)
- [Changelog](#changelog)

## Installation

Install the library using the following command:

```bash
go get github.com/kirinyoku/enkanetwork-go@latest
```

## Quick Start

The following example demonstrates how to fetch a Genshin Impact player profile using the library:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

func main() {
	// Initialize a Genshin Impact client with a custom User-Agent
	client := genshin.NewClient(nil, nil, "my-app/1.0")

	// Create a context for the request
	ctx := context.Background()

	// Fetch player profile (replace "123456789" with a valid 9-digit UID)
	profile, err := client.GetProfile(ctx, "123456789")
	if err != nil {
		log.Fatalf("Failed to fetch profile: %v", err)
	}

	// Display player details
	fmt.Printf("Nickname: %s\n", profile.PlayerInfo.Nickname)
	fmt.Printf("Adventure Rank: %d\n", profile.PlayerInfo.Level)
	fmt.Printf("Signature: %s\n", profile.PlayerInfo.Signature)
}
```

### Notes
- **User-Agent**: Replace `"my-app/1.0"` with a unique identifier for your application (e.g., `app-name/version`).
- **UID**: Use a valid 9-digit UID for the respective game. Invalid UIDs will return `ErrInvalidUIDFormat`.
- **Context**: Always pass a `context.Context` to control timeouts and cancellations.

## Clients

The library provides dedicated clients for each supported game and the EnkaNetwork platform. Each client offers a consistent interface with methods like `GetProfile` and `GetPlayerInfo`.

### Genshin Impact

The Genshin Impact client retrieves player data, including profiles and character details.

```go
import "github.com/kirinyoku/enkanetwork-go/client/genshin"

// Initialize client
client := genshin.NewClient(nil, nil, "my-app/1.0")

// Fetch full player profile
profile, err := client.GetProfile(ctx, "123456789")
if err != nil {
	log.Fatalf("Failed to fetch profile: %v", err)
}

// Fetch basic player info
playerInfo, err := client.GetPlayerInfo(ctx, "123456789")
if err != nil {
	log.Fatalf("Failed to fetch player info: %v", err)
}

fmt.Printf("Nickname: %s\n", playerInfo.Nickname)
```

### Honkai: Star Rail

The Honkai: Star Rail client retrieves player and trailblazer data.

```go
import "github.com/kirinyoku/enkanetwork-go/client/hsr"

// Initialize client
client := hsr.NewClient(nil, nil, "my-app/1.0")

// Fetch full player profile
profile, err := client.GetProfile(ctx, "123456789")
if err != nil {
	log.Fatalf("Failed to fetch profile: %v", err)
}

// Fetch basic player info
playerInfo, err := client.GetPlayerInfo(ctx, "123456789")
if err != nil {
	log.Fatalf("Failed to fetch player info: %v", err)
}

fmt.Printf("Nickname: %s\n", playerInfo.Nickname)
```

### Zenless Zone Zero

The Zenless Zone Zero client retrieves player and agent data.

```go
import "github.com/kirinyoku/enkanetwork-go/client/zzz"

// Initialize client
client := zzz.NewClient(nil, nil, "my-app/1.0")

// Fetch full player profile
profile, err := client.GetProfile(ctx, "123456789")
if err != nil {
	log.Fatalf("Failed to fetch profile: %v", err)
}

// Fetch basic player info
playerInfo, err := client.GetPlayerInfo(ctx, "123456789")
if err != nil {
	log.Fatalf("Failed to fetch player info: %v", err)
}

fmt.Printf("Nickname: %s\n", playerInfo.Nickname)
```

### EnkaNetwork

The EnkaNetwork client provides access to user-specific features, such as Enka profiles and linked HoYoverse game data.

```go
import "github.com/kirinyoku/enkanetwork-go/client/enka"

// Initialize client
client := enka.NewClient(nil, nil, "my-app/1.0")

// Fetch EnkaNetwork user profile
userProfile, err := client.GetUserProfile(ctx, "enka_username")
if err != nil {
	log.Fatalf("Failed to fetch user profile: %v", err)
}

// Fetch HoYoverse profile linked to Enka account
hoyoProfile, err := client.GetUserProfileHoyo(ctx, "enka_username", "hoyo_hash")
if err != nil {
	log.Fatalf("Failed to fetch HoYoverse profile: %v", err)
}

fmt.Printf("Username: %s\n", userProfile.Username)
```

### Client Configuration

All clients accept optional `http.Client` and `Cache` parameters for customization:

```go
import "net/http"

// Custom HTTP client with timeout
httpClient := &http.Client{
	Timeout: 10 * time.Second,
}

// Initialize client with custom HTTP client and cache
client := genshin.NewClient(httpClient, cache, "my-app/1.0")
```

## Advanced Features

### Caching

To reduce API requests and improve performance, the library supports custom caching via the `Cache` interface.

```go
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
}
```

**Example: In-memory cache**

```go
package main

import (
	"log"
	"time"
	"context"
	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

// Simple in-memory cache
type simpleCache struct {
	data map[string]cacheEntry
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

func (c *simpleCache) Get(key string) (interface{}, bool) {
	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.value, true
}

func (c *simpleCache) Set(key string, value interface{}, expiration time.Duration) {
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(expiration),
	}
}

func main() {
	// Initialize cache
	cache := &simpleCache{data: make(map[string]cacheEntry)}

	// Initialize client with cache
	client := genshin.NewClient(nil, cache, "my-app/1.0")

	// Fetch profile with caching
	profile, err := client.GetProfile(context.Background(), "123456789")
	if err != nil {
		log.Fatalf("Failed to fetch profile: %v", err)
	}
}
```

**Best Practices**:
- Use the `ttl` field from API responses to set cache expiration.
- For production, consider persistent storage like Redis or a database.
- Ensure cache keys are unique to avoid conflicts across clients.

### Error Handling

The library provides specific error types in the `common` package to handle API scenarios:

- `ErrInvalidUIDFormat`: UID is not a valid 9-digit number.
- `ErrPlayerNotFound`: Player with the given UID does not exist.
- `ErrUserNotFound`: EnkaNetwork user not found.
- `ErrInvalidUsername`: Invalid EnkaNetwork username.
- `ErrRateLimited`: API rate limit exceeded (after retries).
- `ErrServerMaintenance`: API is under maintenance.
- `ErrServerError`: General server error.
- `ErrServiceUnavailable`: API is unavailable.

**Example: Handling errors**

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kirinyoku/enkanetwork-go/client/genshin"
	"github.com/kirinyoku/enkanetwork-go/common"
)

func main() {
	client := genshin.NewClient(nil, nil, "my-app/1.0")
	profile, err := client.GetProfile(context.Background(), "invalid")
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidUIDFormat):
			fmt.Println("Error: Invalid UID format. Use a 9-digit UID.")
		case errors.Is(err, common.ErrPlayerNotFound):
			fmt.Println("Error: Player not found.")
		case errors.Is(err, common.ErrRateLimited):
			fmt.Println("Error: Rate limit exceeded. Please try again later.")
		case errors.Is(err, common.ErrServerMaintenance):
			fmt.Println("Error: API is under maintenance. Check https://status.enka.network/.")
		default:
			fmt.Printf("Unexpected error: %v\n", err)
		}
		return
	}
	fmt.Printf("Player: %s\n", profile.PlayerInfo.Nickname)
}
```

**Tips**:
- Use `errors.Is` for type-safe error checking.
- Implement exponential backoff for `ErrRateLimited` in production.

### Context Support

All API methods accept a `context.Context` for managing timeouts and cancellations.

**Example: Timeout**

```go
import (
	"context"
	"log"
	"time"
	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

func main() {
	client := genshin.NewClient(nil, nil, "my-app/1.0")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile, err := client.GetProfile(ctx, "123456789")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	log.Printf("Player: %s", profile.PlayerInfo.Nickname)
}
```

**Example: Cancellation**

```go
import (
	"context"
	"log"
	"time"
	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

func main() {
	client := genshin.NewClient(nil, nil, "my-app/1.0")
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		profile, err := client.GetProfile(ctx, "123456789")
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		log.Printf("Player: %s", profile.PlayerInfo.Nickname)
	}()

	// Cancel the request after 2 seconds
	time.Sleep(2 * time.Second)
	cancel()
}
```

**Best Practices**:
- Always use a timeout in production to prevent hanging requests.
- Use cancellation for long-running operations that may need to be interrupted.

## API Usage Guidelines

To ensure compliance with the EnkaNetwork API and avoid issues:

1. **Rate Limits**:
   - The API enforces strict rate limits. The library retries on `429 Too Many Requests` errors, but caching is critical to minimize requests.
   - Check [EnkaNetwork Status](https://status.enka.network/) for current rate limit details.

2. **User-Agent**:
   - Provide a unique User-Agent (e.g., `"my-app/1.0"`) to identify your application.
   - Avoid generic or missing User-Agents, as they may lead to request blocking.

3. **Caching**:
   - Respect the `ttl` value in API responses for cache expiration.
   - Cached responses count toward rate limits, so implement efficient caching.
   - Use unique cache keys per client and endpoint to avoid conflicts.

4. **Error Handling**:
   - Handle all errors listed in the [Error Handling](#error-handling) section.
   - Implement retry logic with exponential backoff for `ErrRateLimited`.

5. **Monitoring**:
   - Regularly check [EnkaNetwork Status](https://status.enka.network/) for maintenance or downtime.
   - Log errors and monitor API responses for unexpected behavior.

6. **Best Practices**:
   - Use context timeouts for all API calls in production.
   - Test with valid UIDs during development to avoid `ErrPlayerNotFound`.
   - Keep your library version up-to-date to benefit from bug fixes and improvements.

## Contributing

This is my first library, and due to my lack of experience, it is far from perfect, so I would welcome your contributions! Here's how you can help

1. **Bug Reports**
   - Open an issue for any bugs you find
   - Include steps to reproduce the issue
   - Provide error messages and relevant code

2. **Feature Requests**
   - Open an issue describing the feature
   - Explain why it would be useful
   - Provide examples if possible

3. **Code Contributions**
   - Fork the repository
   - Create a feature branch
   - Submit a pull request
   - Include tests for new features

## License

This library is licensed under the [MIT License](LICENSE). See the [LICENSE](LICENSE) file for details.

## Changelog

For a detailed history of changes, including new features and bug fixes, see the [CHANGELOG](CHANGELOG.md).

## References

- [EnkaNetwork API Documentation](https://api.enka.network/#/api)
- [Genshin Impact API](https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md)
- [Zenless Zone Zero API](https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md)
- [EnkaNetwork Status](https://status.enka.network/)
- [GitHub Repository](https://github.com/kirinyoku/enkanetwork-go)