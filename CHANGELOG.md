# Changelog

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
