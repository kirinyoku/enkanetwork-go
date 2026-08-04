# EnkaNetwork Go

A lightweight Go wrapper for the [EnkaNetwork API](https://api.enka.network/#/api). It supports:

- **Genshin Impact**
- **Honkai: Star Rail**
- **Zenless Zone Zero**
- **Arknights Endfield**

[![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/enkanetwork-go.svg)](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go)
[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Live API Status

[![Genshin API](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-genshin.yml/badge.svg?branch=main)](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-genshin.yml) [![HSR API](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-hsr.yml/badge.svg?branch=main)](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-hsr.yml) [![ZZZ API](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-zzz.yml/badge.svg?branch=main)](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-zzz.yml) [![Endfield API](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-endfield.yml/badge.svg?branch=main)](https://github.com/kirinyoku/enkanetwork-go/actions/workflows/api-endfield.yml)

These badges reflect the real-time compatibility of this library with the EnkaNetwork API, verified daily via automated integration tests against live player profiles.

- **`passing`**: All data structures are fully up-to-date with the live API.
- **`failing`**: The API has changed (e.g. new fields). Thanks to *Drift Tolerance*, your application should continue working (unknown fields are caught in the `Extra` map), but a minor library update will be required to map the new fields into Go structs.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
  - [Caching](#caching)
  - [Retries](#retries)
- [Error Handling](#error-handling)
- [API Changes (Drift Tolerance)](#api-changes-drift-tolerance)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Fully Typed Models:** All known API responses are mapped to standard Go structs for safe and easy access.
- **Drift Tolerance:** New or unknown API fields are safely captured in an `Extra` map, so game updates won't break your code.
- **Accurate Zero Values:** Safely handles empty arrays, `0`, `false`, and `null` without losing data or panicking.
- **Smart Type Parsing:** Automatically handles fields that unpredictably switch between strings and numbers (e.g., `"123"` vs `123`).
- **Context Support:** Fully supports `context.Context` for request timeouts and cancellations.
- **Custom HTTP Clients:** Bring your own `http.Client` for proxies or custom transport settings.
- **Built-in Caching Support:** Easily plug in any caching layer via a simple interface to avoid rate limits.
- **Auto Retries:** Configurable automatic retries with backoff for handling temporary network issues.
- **Clear Error Handling:** Package-level errors make it easy to handle specific API failures.
- **Lightweight:** No external dependencies.

## Installation

> [!IMPORTANT]
> **Version 1.0.0 Update:** The repository was recently recreated, which caused version caching issues on `pkg.go.dev` with old legacy versions (e.g. `v0.5.5`). To reset the module proxy cache and provide a clean slate, the project has bumped directly to `v1.0.0`. All versions `<= v0.5.5` have been officially retracted.

```bash
go get github.com/kirinyoku/enkanetwork-go@latest
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kirinyoku/enkanetwork-go/client/genshin"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := genshin.New(genshin.Options{
		UserAgent: "my-app/1.0",
	})

	profile, err := client.GetProfile(ctx, "618285856")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Nickname:", profile.PlayerInfo.Nickname)
	fmt.Println("Adventure Rank:", profile.PlayerInfo.Level)
	fmt.Println("Showcase Characters:", len(profile.AvatarInfoList))
}
```

> [!NOTE]
> Clients for Honkai: Star Rail (`client/hsr`), Zenless Zone Zero (`client/zzz`), and Arknights Endfield (`client/endfield`) share the exact same API surface. Just swap the `genshin` import path with the game you need!

Full runnable examples are available in [`examples/`](examples).

## Configuration

Clients are configured at construction time through `Options`. After a client is
created, treat its configuration as read-only.

```go
client := genshin.New(genshin.Options{
	HTTPClient: &http.Client{
		Timeout: 10 * time.Second,
	},
	UserAgent: "my-app/1.0",
	Retry: &genshin.RetryOptions{
		MaxAttempts: 2,
		Delay:       2 * time.Second,
	},
})
```

If no HTTP client is provided, the library uses a default client with a
10-second timeout. If no retry options are provided, retryable responses are
attempted up to 3 times.

### Caching

Caching is optional. The library does not force a cache implementation; pass any
type that satisfies the cache interface:

```go
type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any, expiration time.Duration)
}
```

See [`examples/advanced/cache`](examples/advanced/cache/main.go) for a small
in-memory cache example.

### Retries

By default, clients make up to 3 attempts for retryable responses and respect
the `Retry-After` header when the API sends it. When `Retry-After` is not
available, the configured fallback delay is used.

Set `MaxAttempts` to `1` to disable retries:

```go
client := genshin.New(genshin.Options{
	UserAgent: "my-app/1.0",
	Retry: &genshin.RetryOptions{
		MaxAttempts: 1,
	},
})
```

See [`examples/advanced/retry`](examples/advanced/retry/main.go) for a focused
retry configuration example.

## Error Handling

Client packages expose errors that can be checked with `errors.Is`:

```go
profile, err := client.GetProfile(ctx, uid)
if err != nil {
	switch {
	case errors.Is(err, genshin.ErrInvalidUIDFormat):
		// UID is not valid for this endpoint.
	case errors.Is(err, genshin.ErrPlayerNotFound):
		// EnkaNetwork could not find this player.
	case errors.Is(err, genshin.ErrRateLimited):
		// Retry later or reduce request volume.
	default:
		// Handle other network, API, or decoding errors.
	}
	return
}
```

Common errors include invalid UID format, player not found, rate limiting,
server maintenance, server errors, and service unavailability.

## API Changes (Drift Tolerance)

Games update often, and the EnkaNetwork API changes with them. The API might add new fields or change data types.

Standard Go structs might fail to read the JSON if a type changes, or they might ignore new fields. To prevent errors and keep your app working, this library uses two strategies:

### 1. Catching New Fields (`Extra`)

When the API returns a new field that this library doesn't know about yet, it saves it in the `Extra map[string]json.RawMessage` field.

This means **you don't have to wait for a library update** to use new data. You can read it yourself right away:

```go
profile, err := client.GetProfile(ctx, "618285856")
if err != nil {
	log.Fatal(err)
}

// If the API adds a new field that is not yet in the Go structs:
if raw, ok := profile.Extra["someNewField"]; ok {
	// You can decode `raw` directly in your app
	_ = raw
}
```

### 2. Flexible Data Types

Sometimes the API returns a value as a string (`"123"`), but later changes it to a number (`123`). Standard Go JSON decoding will crash when this happens. To fix this, we use helper types like `models.StringNumber`. They can safely read both strings and numbers without breaking your app.

## Documentation

- [API reference on pkg.go.dev](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go)
- [Examples](examples)
- [Enka.Network Profiles API](docs/enka/api.md)
- [Zenless Zone Zero API](docs/zzz/api.md)
- [Changelog](CHANGELOG.md)

Useful EnkaNetwork links:

- [EnkaNetwork API documentation](https://api.enka.network/#/api)
- [Genshin Impact API details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md)
- [Zenless Zone Zero API details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md)
- [EnkaNetwork status](https://status.enka.network/)

## Contributing

Contributions are welcome! Please read the [Contributing Guide](CONTRIBUTING.md)
to get started. It covers everything from setting up your environment to
submitting a pull request — no prior open source experience required.

## License

Licensed under the [MIT License](LICENSE).
