# Changelog

## [v1.0.0] - 2026-07-31

### Changed
- **Version Bump:** Skipped directly to `v1.0.0` to reset the Go module proxy cache (`pkg.go.dev`) after the repository was recreated.
- **Stable Release:** Formalized the current API as the stable `1.0.0` release.

### Deprecated
- **Retracted Legacy Versions:** All versions `<= v0.5.5` (published under the old repository) have been officially retracted in `go.mod` and should no longer be used.

## [v0.2.0] - 2026-07-30

### Changed
- Improved and simplified the internal architecture without breaking backward compatibility.
- Removed duplicated code by centralizing network requests (`Fetcher`) into the `core` package.
- Made caching cleaner and safer using Go Generics. Each game model now defines its own cache lifetime (`CacheTTL()`).
- Cleaned up the folder structure by moving errors directly into the `core` package for a flatter, simpler layout.

## [v0.1.2] - 2026-07-30

### Fixed
- Added `omitempty` tag to `KMOHDEAKEFG` field in ZZZ `TitleInfo` model to tolerate its removal from the API and fix integration test failures.

## [v0.1.1] - 2026-07-12

### Changed
- Improved JSON handling performance for flexible scalar types and unknown-field preservation.
- Reduced allocations when merging known, extra, and raw JSON fields.
- Reduced memory usage when draining non-successful HTTP responses.
- Improved retry timer handling to avoid lingering timers after cancellation.

### Fixed
- Preserved exact JSON number comparison in diff helpers for large integer values.
- Handled `Retry-After` parsing more strictly for invalid, negative, or overflowing values.

## [v0.1.0] - 2026-07-07

### Added

- Added initial EnkaNetwork Go API wrapper.
- Added clients for Genshin Impact, Honkai: Star Rail, Zenless Zone Zero, and EnkaNetwork account endpoints.
- Added typed-first models for stable EnkaNetwork response fields.
- Added unknown field preservation through `Extra` on key response and nested API models.
- Added raw JSON preservation for API drift tolerance.
- Added flexible scalar helpers for fields that may arrive as JSON strings or numbers.
- Added configurable `Options` for HTTP client, cache, User-Agent, retry behavior, and base URL.
- Added pluggable cache interface.
- Added configurable retry behavior with `Retry-After` support.
- Added package-level errors for common API failures.
- Added examples for basic client usage, custom cache, and retry configuration.
- Added maintainer documentation for API drift, project architecture, versioning, releases, changelog style, and commit style.

### Notes

- This is the first public pre-1.0 release.
- Public APIs may still change before `v1.0.0`, but changes should be documented in this changelog.
- The library is an API wrapper, not a game SDK.
