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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch player profile
	profile, err := client.GetProfile(ctx, "618285856")
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

The library provides dedicated clients for each supported game and the EnkaNetwork platform.

### Genshin Impact

Genshin Impact client offers a consistent interface with two methods: `GetProfile` and `GetPlayerInfo`.

#### GetProfile

`GetProfile` fetches the complete player profile, including both basic player information and detailed character data (if the showcase is public). This method makes an additional request to obtain the `AvatarInfoList` containing detailed character information.

**Returns**:
- `*Profile` containing:
  - `PlayerInfo`: Basic player information (nickname, level, signature, etc.)
  - `AvatarInfoList`: Detailed information about characters in the showcase
  - `Owner`: Enka profile information (if public and linked)
  - `TTL`: Cache duration in seconds
  - `UID`: The player's UID

**Example**:
```go
import "github.com/kirinyoku/enkanetwork-go/client/genshin"

// Initialize client
client := genshin.NewClient(nil, nil, "my-app/1.0")

// Fetch full player profile with character details
profile, err := client.GetProfile(ctx, "618285856")
if err != nil {
    log.Fatalf("Failed to fetch profile: %v", err)
}

// Access player information
fmt.Printf("Nickname: %s\n", profile.PlayerInfo.Nickname)
fmt.Printf("Adventure Rank: %d\n", profile.PlayerInfo.Level)

// Check if character showcase is available
if len(profile.AvatarInfoList) > 0 {
	fmt.Printf("Characters in showcase: %d\n", len(profile.AvatarInfoList))
	for _, char := range profile.AvatarInfoList {
		fmt.Printf("Character ID: %d (Level %s)\n", char.AvatarID, char.PropMap["4001"].Ival)
	}
}
```

#### GetPlayerInfo

`GetPlayerInfo` fetches only the basic player information without the detailed character data. This method is faster and has fewer rate limits since it makes fewer requests to the API.

**Returns**:
- `*Profile` containing:
  - `PlayerInfo`: Basic player information (nickname, level, signature, etc.)
  - `AvatarInfoList`: **ALWAYS EMPTY SLICE**
  - `Owner`: Enka profile information (if public and linked)
  - `TTL`: Cache duration in seconds
  - `UID`: The player's UID

**Example**:
```go
// Fetch basic player info (faster, fewer rate limits)
profile, err := client.GetPlayerInfo(ctx, "618285856")
if err != nil {
    log.Fatalf("Failed to fetch player info: %v", err)
}

fmt.Printf("Player: %s (AR %d)\n", 
    profile.PlayerInfo.Nickname, 
    profile.PlayerInfo.Level)
```

#### Error Handling

The Genshin Impact client provides a set of error types for common API scenarios.

**Example**:
```go
import "github.com/kirinyoku/enkanetwork-go/client/genshin"

client := genshin.NewClient(nil, nil, "my-app/1.0")
profile, err := client.GetProfile(ctx, "invalid")
if err != nil {
	switch {
	case errors.Is(err, genshin.ErrInvalidUIDFormat):
		fmt.Println("Error: Invalid UID format. Use a 9-digit UID.")
	case errors.Is(err, genshin.ErrPlayerNotFound):
		fmt.Println("Error: Player not found.")
	case errors.Is(err, genshin.ErrRateLimited):
		fmt.Println("Error: Rate limit exceeded. Please try again later.")
	case errors.Is(err, genshin.ErrServerMaintenance):
		fmt.Println("Error: API is under maintenance. Check https://status.enka.network/.")
	default:
		fmt.Printf("Unexpected error: %v\n", err)
	}
	return
}
fmt.Printf("Player: %s\n", profile.PlayerInfo.Nickname)
```

### Honkai: Star Rail

The Honkai: Star Rail client retrieves comprehensive player and characters data through a single `GetProfile` method. Unlike the Genshin Impact client, HSR's API provides all data in a single request. 

#### GetProfile

`GetProfile` fetches the complete player profile, including detailed character information, equipment, and trailblazer data in a single request.

**Returns**:
- `*Profile` containing:
  - `DetailInfo`: Player and characters information
  - `TTL`: Cache duration in seconds
  - `Owner`: Enka profile information (if public and linked)
  - `UID`: The player's UID

**Example**:
```go
import "github.com/kirinyoku/enkanetwork-go/client/hsr"

// Initialize client
client := hsr.NewClient(nil, nil, "my-app/1.0")

// Create a context for the request
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Fetch player profile
profile, err := client.GetProfile(ctx, "800579959")
if err != nil {
    log.Fatalf("Failed to fetch profile: %v", err)
}

// Access player information
fmt.Printf("Nickname: %s\n", profile.DetailInfo.Nickname)
fmt.Printf("Trailblaze Level: %d\n", profile.DetailInfo.Level)

// Access character information
if len(profile.DetailInfo.AvatarDetailList) > 0 {
	for _, char := range profile.DetailInfo.AvatarDetailList {
		fmt.Printf("- %d (Level %d)\n", char.AvatarID, char.Level)
	}
}
```

#### Error Handling

The Honkai: Star Rail client provides a set of error types for common API scenarios.

**Example**:
```go
import "github.com/kirinyoku/enkanetwork-go/client/hsr"

client := hsr.NewClient(nil, nil, "my-app/1.0")
profile, err := client.GetProfile(ctx, "invalid")
if err != nil {
	switch {
	case errors.Is(err, hsr.ErrInvalidUIDFormat):
		fmt.Println("Error: Invalid UID format. Use a 9-digit UID.")
	case errors.Is(err, hsr.ErrPlayerNotFound):
		fmt.Println("Error: Player not found.")
	case errors.Is(err, hsr.ErrRateLimited):
		fmt.Println("Error: Rate limit exceeded. Please try again later.")
	case errors.Is(err, hsr.ErrServerMaintenance):
		fmt.Println("Error: API is under maintenance. Check https://status.enka.network/.")
	default:
		fmt.Printf("Unexpected error: %v\n", err)
	}
	return
}
fmt.Printf("Player: %s\n", profile.DetailInfo.Nickname)
```

### Zenless Zone Zero

#### GetProfile

`GetProfile` fetches the complete player profile, including detailed character information, equipment, and trailblazer data in a single request.

**Returns**:
- `*Profile` containing:
  - `PlayerInfo`: Basic information about the game account from the player's showcase
  - `TTL`: Cache duration in seconds
  - `Owner`: Enka profile information (if public and linked)
  - `UID`: The player's UID

**Example**:
```go
import "github.com/kirinyoku/enkanetwork-go/client/zzz"

// Initialize client
client := zzz.NewClient(nil, nil, "my-app/1.0")

// Create a context for the request
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Fetch player profile
profile, err := client.GetProfile(ctx, "1504687050")
if err != nil {
    log.Fatalf("Failed to fetch profile: %v", err)
}

// Access player information
fmt.Printf("Nickname: %s\n", profile.PlayerInfo.SocialDetail.ProfileDetail.Nickname)
fmt.Printf("Trailblaze Level: %d\n", profile.PlayerInfo.SocialDetail.ProfileDetail.Level)

// Access characters information
if len(profile.PlayerInfo.ShowcaseDetail.AvatarList) > 0 {
	for _, avatar := range profile.PlayerInfo.ShowcaseDetail.AvatarList {
		fmt.Printf("Character ID: %d (Level %d)\n", avatar.ID, avatar.Level)
	}
}
```

#### Error Handling

The Zenless Zone Zero client provides a set of error types for common API scenarios.

**Example**:
```go
import "github.com/kirinyoku/enkanetwork-go/client/zzz"

client := zzz.NewClient(nil, nil, "my-app/1.0")
profile, err := client.GetProfile(ctx, "invalid")
if err != nil {
	switch {
	case errors.Is(err, zzz.ErrInvalidUIDFormat):
		fmt.Println("Error: Invalid UID format. Use a 9-digit or 10-digit UID.")
	case errors.Is(err, zzz.ErrPlayerNotFound):
		fmt.Println("Error: Player not found.")
	case errors.Is(err, zzz.ErrRateLimited):
		fmt.Println("Error: Rate limit exceeded. Please try again later.")
	case errors.Is(err, zzz.ErrServerMaintenance):
		fmt.Println("Error: API is under maintenance. Check https://status.enka.network/.")
	default:
		fmt.Printf("Unexpected error: %v\n", err)
	}
	return
}
fmt.Printf("Player: %s\n", profile.DetailInfo.Nickname)
```

### EnkaNetwork

The EnkaNetwork client provides access to EnkaNetwork user profiles and their linked game accounts. This client allows you to fetch user information, their verified game accounts, and character builds.

#### Initialization

```go
import "github.com/kirinyoku/enkanetwork-go/client/enka"

// Initialize client with default settings
client := enka.NewClient(nil, nil, "my-app/1.0")

// Or with custom HTTP client and cache
customClient := &http.Client{Timeout: 20 * time.Second}
cache := // your cache implementation
client := enka.NewClient(customClient, cache, "my-app/1.0")
```

#### GetUserProfile

Fetches the Enka user profile for the given username.

**Returns**:
- `*Owner` containing:
	- ID: The Enka user ID
	- Hash: The Enka user hash
	- Username: The Enka username
	- Profile: Patreon profile data for Patreon members
- `error`: An error if the request fails

**Example**:
```go
profile, err := client.GetUserProfile(ctx, "Algoinde")
if err != nil {
    log.Fatalf("Failed to fetch user profile: %v", err)
}
fmt.Printf("Username: %s\n", profile.Username)
fmt.Printf("Bio: %s\n", profile.Profile.Bio)
```

#### GetUserProfileHoyos

Fetches a list of verified and public game accounts (hoyos) linked to an EnkaNetwork profile.

**Returns**:
- `*Hoyos`: A map of game account hashes to their metadata (map[string]Hoyo)
- `error`: An error if the request fails

**Example**:
```go
hoyos, err := client.GetUserProfileHoyos(ctx, "Algoinde")
if err != nil {
	log.Fatalf("Failed to fetch hoyos: %v", err)
}

for hash, hoyo := range *hoyos {
	fmt.Printf("hoyo_hash: %s, region: %s, uid: %d\n", hash, hoyo.Region, hoyo.UID)
}
```

#### GetUserProfileHoyo

Fetches detailed information about a specific game account (hoyo).

**Parameters**:
- `username`: EnkaNetwork username
- `hoyo_hash`: The hash of the game account (obtained from GetUserProfileHoyos)

**Returns**:
- `*Hoyo` containing:
	- UID: The UID of the game account
	- UIDPublic: Whether the UID is public
	- Public: Whether the Hoyo account is public
	- Verified: Whether the Hoyo account is verified
	- PlayerInfo: Player information for the account
	- Hash: The hash of the game account
	- Region: The region of the game account
	- AvatarOrder: The order of the characters in the game account
	- Order: The order of the Hoyo account
	- LivePublic: Whether the live build is public
	- HoyoType: The ID of the Hoyo game (0 for Genshin, 1 for HSR, 2 for ZZZ)
- `error`: An error if the request fails

**Example**:
```go
hoyo, err := client.GetUserProfileHoyo(ctx, "Algoinde", "4Wjv2e")
if err != nil {
	log.Fatalf("Failed to fetch hoyo: %v", err)
}

switch hoyo.HoyoType {
case 0:
	fmt.Println("Genshin Impact")
case 1:
	fmt.Println("Honkai: Star Rail")
case 2:
	fmt.Println("Zenless Zone Zero")
default:
	fmt.Println("Unknown game")
}
```

#### GetUserProfileHoyoBuilds

Fetches character builds for a specific game account.

**Parameters**:
- `username`: EnkaNetwork username
- `hoyo_hash`: The hash of the game account

**Returns**:
- `*Builds`: A map of character IDs to their builds (map[string][]Build)
- `error`: An error if the request fails

**Example**:
```go
builds, err := client.GetUserProfileHoyoBuilds(ctx, "Algoinde", "4Wjv2e")
if err != nil {
    log.Fatalf("Failed to fetch builds: %v", err)
}
for charID, charBuilds := range *builds {
    fmt.Printf("Character %s has %d builds\n", charID, len(charBuilds))
}
```

#### Error Handling

The EnkaNetwork client provides a set of error types for common API scenarios.

**Example**:
```go
import "github.com/kirinyoku/enkanetwork-go/client/enka"

client := enka.NewClient(nil, nil, "my-app/1.0")
profile, err := client.GetUserProfile(ctx, "invalid")
if err != nil {
	switch {
	case errors.Is(err, enka.ErrInvalidUsername):
		fmt.Println("Error: Invalid username.")
	case errors.Is(err, enka.ErrPlayerNotFound):
		fmt.Println("Error: Player not found.")
	case errors.Is(err, enka.ErrRateLimited):
		fmt.Println("Error: Rate limit exceeded. Please try again later.")
	case errors.Is(err, enka.ErrServerMaintenance):
		fmt.Println("Error: API is under maintenance. Check https://status.enka.network/.")
	default:
		fmt.Printf("Unexpected error: %v\n", err)
	}
	return
}
fmt.Printf("Player: %s\n", profile.Username)
```

## Advanced Features

### Caching

To reduce API requests and improve performance, the library supports custom caching via the `Cache` interface.

```go
type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any, expiration time.Duration)
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
	value     any
	expiresAt time.Time
}

func (c *simpleCache) Get(key string) (any, bool) {
	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.value, true
}

func (c *simpleCache) Set(key string, value any, expiration time.Duration) {
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

	profile, err := client.GetProfile(ctx, "618285856")
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
		profile, err := client.GetProfile(ctx, "618285856")
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