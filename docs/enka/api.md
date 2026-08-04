# Enka.Network API - Profile Structures

This document describes the JSON structures returned by the **Enka.Network Profiles API**.

## Data Structure Overview

```text
Owner
├── id
├── is_staff
├── hash
├── username
└── profile (*PatreonProfile)
    ├── bio
    ├── level
    ├── avatar
    └── image_url

Hoyo
├── user (*Owner)
├── uid
├── uid_public
├── public
├── live_public
├── verified
├── player_info
├── hash
├── region
├── order
├── avatar_order
├── hoyo_type
└── live_data_hash

Build
├── id
├── name
├── avatar_id
├── owner
├── avatar_data (*AvatarData)
├── order
├── live
├── settings (*Settings)
│   ├── adaptiveColor
│   ├── artSource
│   ├── caption
│   ├── honkardWidth
│   └── transform
├── public
├── image
├── hoyo_type
└── hoyo
```

---

## Enka Profile

### Owner

| Name | Type | Description |
| :--- | :--- | :---------- |
| id | `int` | Enka.Network User ID. |
| is_staff | `*bool` | Indicates whether the user is an Enka.Network staff member. |
| hash | `string` | The user's unique Enka profile hash (used for routing/URLs). |
| username | `string` | The user's chosen Enka.Network username. |
| profile | `*PatreonProfile` | Additional profile data, typically present if the user is a Patreon supporter. |

### PatreonProfile (EnkaProfile)

| Name | Type | Description |
| :--- | :--- | :---------- |
| bio | `string` | The user's custom biography written on their Enka.Network profile. |
| level | `int` | The user's Patreon membership tier/level. |
| avatar | `string` | The file name or identifier for the user's custom avatar on Enka. |
| image_url | `string` | The full URL to the user's profile picture (usually sourced from Patreon). |

---

## Linked Accounts (Hoyos)

Users can link multiple game accounts to their Enka profile. These are called "Hoyos".

### Hoyo

| Name | Type | Description |
| :--- | :--- | :---------- |
| user | `*Owner` | The Enka user who owns this linked game account. |
| uid | `*int` | The game account UID (if provided/public). |
| uid_public | `*bool` | Whether the user has chosen to make this account's UID public. |
| public | `*bool` | Whether the Hoyo account is public. |
| live_public | `*bool` | Whether the live build showcase for this account is public. |
| verified | `*bool` | Whether the user has verified ownership of this game account. |
| player_info | `*PlayerInfo` | Game-specific player information (similar to showcase PlayerInfo). |
| hash | `string` | A unique string identifier for this specific linked account. |
| region | `string` | The server region of the game account. |
| order | `string` | The sorting order of this account on the user's profile. |
| avatar_order | `map[string]int` | Custom sorting order for the characters/avatars inside this account. |
| hoyo_type | `int` | API game type discriminator (e.g., `0` = Genshin, `1` = HSR, `2` = ZZZ). |
| live_data_hash | `int` | Hash of the current live data for cache invalidation. |

---

## Saved Builds

Users can save snapshots of their characters on Enka.Network.

### Build

> [!NOTE]
> The exact structure of `avatar_data` within a `Build` will depend on the game (`HoyoType`). For example, a ZZZ build will contain ZZZ-specific `zzz.AvatarData`, while a Genshin build will contain Genshin-specific `genshin.AvatarInfo`.

| Name | Type | Description |
| :--- | :--- | :---------- |
| id | `int` | The unique internal database ID of the saved build on Enka.Network. |
| name | `string` | The custom name given to the build by the user. |
| avatar_id | `string` | The string identifier of the character/agent this build belongs to. |
| owner | `string` | Owner identifier associated with the build. |
| avatar_data | `*AvatarData` | The full, game-specific character data snapshot (e.g., stats, weapons, artifacts at the time of saving). |
| order | `string` | The sorting order of the saved build on the user's Enka profile. |
| live | `*bool` | If true, indicates this is not a saved build but rather a live snapshot retrieved directly from the game's showcase when the user clicked "refresh". |
| [settings](#settings) | `*Settings` | The visual configuration and settings used for rendering the build card. |
| public | `*bool` | Whether the user has marked this specific build as public or private. |
| image | `*string` | The URL of the custom art/image assigned to the build card. |
| hoyo_type | `int` | An API discriminator indicating the game/account type (e.g., `0` for Genshin, `1` for HSR). |
| hoyo | `string` | A unique identifier (hash) linking the build to a specific HoYoverse account. |

### Settings

| Name | Type | Description |
| :--- | :--- | :---------- |
| adaptiveColor | `*bool` | Whether the build card should automatically adapt its color palette based on the character/image. |
| artSource | `*string` | A link or text crediting the original source/artist of the custom image used in the build card. |
| caption | `*string` | A custom caption or note written on the build card. |
| honkardWidth | `*float64` | Visual layout configuration: the specific width of the card rendering. |
| transform | `*string` | CSS-like transformation properties applied to position/scale the custom image within the build card. |
