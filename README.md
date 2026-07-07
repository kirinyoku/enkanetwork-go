# EnkaNetwork Go

A lightweight Go wrapper for the [EnkaNetwork API](https://api.enka.network/#/api).
It fetches EnkaNetwork JSON responses and decodes them into typed Go structs for:

- **Genshin Impact**
- **Honkai: Star Rail**
- **Zenless Zone Zero**

[![Go Reference](https://pkg.go.dev/badge/github.com/kirinyoku/enkanetwork-go.svg)](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go)
[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Installation

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

Full runnable examples are available in [`examples/`](examples).

## Scope

This project is an API wrapper, not a game SDK.

The library exposes the data returned by EnkaNetwork in a Go-friendly form. It
does not translate game IDs into names, ship game metadata tables, calculate
builds, or replace EnkaNetwork API semantics with library-specific behavior.

For example, if EnkaNetwork returns a character ID, this wrapper returns that
character ID. Applications can choose how to map IDs to names, icons, or other
game data.

## Features

- **Typed-first models** for stable EnkaNetwork response fields.
- **API drift tolerance** through `Extra` fields on key response and nested
  models.
- **Presence preservation** for explicitly returned zero-like values such as
  `0`, `false`, empty strings, empty arrays, and `null`.
- **Flexible scalar helpers** for fields that may drift between JSON strings and
  numbers.
- **Context support** for cancellation and request timeouts.
- **Custom HTTP clients** through `Options.HTTPClient`.
- **Pluggable caching** through the `Cache` interface.
- **Configurable retries** through `Options.Retry`.
- **Package-level errors** for common API failures.

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

## API Drift Tolerance

EnkaNetwork responses can change after game patches. This library keeps typed
access for stable fields while preserving newly added or unknown fields through
`Extra map[string]json.RawMessage` on key API models.

```go
profile, err := client.GetProfile(ctx, "618285856")
if err != nil {
	log.Fatal(err)
}

if raw, ok := profile.Extra["newApiField"]; ok {
	// Decode raw if your application needs this field before it becomes typed.
	_ = raw
}
```

Some volatile scalar fields use helper types from `models`, such as
`models.StringNumber`, so the client can decode API values that may arrive as
either JSON strings or numbers.

## Documentation

- [API reference on pkg.go.dev](https://pkg.go.dev/github.com/kirinyoku/enkanetwork-go)
- [Examples](examples)
- [Changelog](CHANGELOG.md)

Useful EnkaNetwork links:

- [EnkaNetwork API documentation](https://api.enka.network/#/api)
- [Genshin Impact API details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/gi/api.md)
- [Zenless Zone Zero API details](https://github.com/EnkaNetwork/API-docs/blob/master/docs/zzz/api.md)
- [EnkaNetwork status](https://status.enka.network/)

## Versioning

This project follows semantic versioning after the public baseline release.

- **Patch releases** fix decoding, preservation, integration drift, or internal
  behavior without intentional public API changes.
- **Minor releases** add typed fields, helpers, endpoints, or other
  backward-compatible functionality.
- **Major releases** are reserved for deliberate breaking changes to public
  models, methods, or behavior.

Because EnkaNetwork is a live API, model updates may be needed after game
patches. The project prefers accurate, drift-tolerant model shapes over hiding
API changes behind misleading Go zero values.

## Contributing

Contributions are welcome. This project wraps a live API, so bug reports, drift
reports, fixture updates, model improvements, and documentation fixes are
especially helpful.

When opening an issue, include:

- The game/client you are using.
- The UID or endpoint shape involved, if it is safe to share.
- The error message or unexpected field behavior.
- The library version or commit.

When opening a pull request:

- Keep changes focused.
- Add or update tests for behavior changes.
- Preserve unknown API fields when adding or changing models.
- Run `go test ./...` before submitting when possible.

## Requirements

- Go 1.20 or newer.

## License

Licensed under the [MIT License](LICENSE).
