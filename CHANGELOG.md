# Changelog

## [0.1.0] - 2025-05-29
### Added
- The initial version of the `EnkaNetwork Go Client` library for working with the EnkaNetwork API.
- Support for three HoYoverse games:
  - Genshin Impact.
  - Honkai: Star Rail.
  - Zenless Zone Zero.
- Enka Network client implemented to access EnkaNetwork profiles and data.
- Main features:
  - Get full player profile (`GetProfile`) and basic information (`GetPlayerInfo`) for all supported games.
  - Integration with `context.Context` to support canceling requests and timeouts.
  - Built-in support for caching API responses via the `Cache` interface.
  - Specialized error types to handle API scenarios such as `ErrInvalidUIDFormat`, `ErrPlayerNotFound`, `ErrRateLimited`, etc.
  - Security for competitive use by multiple goroutines.
  - Strong typing of responses for all game data.
- Added documentation with examples of use in `README.md`.