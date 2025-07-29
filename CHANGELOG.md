# Changelog

## [0.3.0] - 2025-07-29
### Changed
- Renamed `enka.Builds` type to `enka.AvatarBuildsMap` for better clarity and, accordingly, the `GetUserProfileHoyoBuilds` method of the `enka` package that returns it.
- The `GetUserProfileHoyos` method of the `enka` package has been updated and now returns the `enka.Hoyos` map directly, rather than a pointer to the map.
- The `GetUserProfile` method of the `enka` package has been updated and now returns the `*enka.Owner`, not `*common.Owner`.
- Improved cache key generation in game clients for better consistency
- Updated documentation

### Added
- Added `enka.Owner` and `enka.PatreonProfile` structs to support EnkaNetwork user profiles
- Added `examples/` directory with runnable demos for each client:
  - `examples/genshin/main.go` — Genshin Impact client
  - `examples/hsr/main.go`     — Honkai: Star Rail client
  - `examples/zzz/main.go`     — Zenless Zone Zero client
  - `examples/enka/main.go`    — EnkaNetwork user profiles

## [0.2.0] - 2025-07-28
### Changed
- Moved error definitions from `common` package to their respective game client packages
- Removed the `GetPlayerInfo` method from the `hsr` and `zzz` packages
- Updated documentation for all clients

### Added
- Added `StygianIndex` and `StygianSeconds` fields to `PlayerInfo` for Stygian Onslaught mode in Genshin Impact

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