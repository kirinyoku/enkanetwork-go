// Package core provides the foundational components for interacting with the EnkaNetwork API.
// It is designed to be used internally by the game-specific client packages (genshin, hsr, zzz).
//
// # Overview
//
// The core package includes:
//   - Base HTTP client with request handling and retry logic
//   - Caching interface for API response storage
//   - Common utilities for UID validation and response processing
//   - Shared constants
//
// # Usage
//
// This package is not intended to be used directly by end-users of the library.
// Instead, use one of the game-specific client packages:
//   - github.com/kirinyoku/enkanetwork-go/client/genshin
//   - github.com/kirinyoku/enkanetwork-go/client/hsr
//   - github.com/kirinyoku/enkanetwork-go/client/zzz
//
// # Caching
//
// The package provides a Cache interface that can be implemented to provide custom
// caching behavior.
//
// # Rate Limiting
//
// The client includes built-in support for handling rate limits with exponential backoff.
// By default, it will retry failed requests up to 3 times.
package core
