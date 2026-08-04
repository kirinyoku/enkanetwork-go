# Enka.Network API - Zenless Zone Zero

This document describes the structure and formulas for the Zenless Zone Zero API response provided by Enka.Network. 

## Table of Content

- [Data Structure Overview](#data-structure-overview)
- [Profile & Player Info](#profile--player-info)
- [Characters (AvatarData)](#characters-avatardata)
- [W-Engines (Weapon)](#w-engines-weapon)
- [Drive Discs (Equipment)](#drive-discs-equipment)
- [Definitions](#definitions)
- [Formulas](#formulas)
- [Icons and Localizations](#icons-and-localizations)

---

## Data Structure Overview

> [!NOTE]
> API responses always contain a minimal amount of data (mainly base values and IDs). To get localized names, descriptions, or icons, you must map these IDs to parsed JSON files in [API-docs/store/zzz](https://github.com/EnkaNetwork/API-docs/tree/master/store/zzz). If you need more granular game data, refer to the [ZenlessData](https://git.mero.moe/dimbreath/ZenlessData) repository.

The JSON response is structured as a hierarchical tree. This overview will help you understand how the data is nested.

```text
Profile
├── PlayerInfo
│   ├── SocialDetail
│   │   ├── ProfileDetail
│   │   │   └── TitleInfo
│   │   └── MedalList
│   └── ShowcaseDetail
│       └── AvatarList
│           ├── SkillLevelList
│           ├── Weapon
│           └── EquippedList
│               └── Equipment
│                   ├── MainPropertyList
│                   └── RandomPropertyList
├── ttl
├── uid
├── region
└── owner
```

---

## Profile & Player Info

### Profile

| Name | Type | Description |
| :--- | :--- | :---------- |
| uid | `string` | Player UID |
| ttl | `int` | Seconds until next cache update |
| region | `string` | Player's server region |
| [owner](../enka/api.md#owner) | `*Owner` | Enka profile owner details (if applicable) |
| [PlayerInfo](#playerinfo) | `PlayerInfo` | Container for Profile and Showcase info |

### PlayerInfo

| Name | Type | Description |
| :--- | :--- | :--------- | 
| [SocialDetail](#socialdetail) | `*SocialDetail` | Public profile information, including the player's signature, level, selected namecard, and displayed badges. |
| [ShowcaseDetail](#showcasedetail) | `*ShowcaseDetail` | Detailed information about the agents the player has chosen to display on their in-game profile showcase. Includes their stats, skills, and equipment. |

### SocialDetail

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Desc | `string` | The custom signature or bio text written by the player on their profile. |
| [ProfileDetail](#profiledetail) | `*ProfileDetail` | Core player profile details (Level, UID, Nickname, etc.). |
| [MedalList](#medallist) | `[]Medal` | List of in-game badges the player has chosen to display on their profile. | 

### ProfileDetail

| Name | Type | Description |
| :--- | :--- | :--------- |
| Uid | `int64` | The player's unique Zenless Zone Zero account identifier. |
| Nickname | `string` | The display name chosen by the player to represent their account. |
| ProfileId | `int` | The unique ID of the player's selected avatar image or profile picture. |
| Level | `int` | The player's current account level, known as the Inter-Knot Level. | 
| Title | `int` | The primary Title ID (check [store/zzz/titles.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/titles.json)). |
| CallingCardId | `int` | The ID of the namecard background currently equipped to the profile (check [store/zzz/namecards.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/namecards.json)). |
| AvatarId | `int` | The ID identifying the main sibling protagonist (Wise or Belle) selected during the start of the game. |
| PlatformType | `int` | The numeric ID representing the hardware platform (PC, PS5, or Mobile) currently associated with the account. |
| [TitleInfo](#titleinfo) | `*TitleInfo` | Detailed metadata regarding the player's currently equipped title and its associated dynamic properties. |

#### TitleInfo

| Name | Type | Description |
| :--- | :--- | :--------- |
| Title | `int` | Base Title ID. |
| FullTitle | `int` | Full Title ID, which accounts for specific tiers, rarities, or visual upgrades of the base title (check [store/zzz/titles.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/titles.json)). |
| Args | `[]any` | Dynamic values injected into the title's text (e.g., a specific completion time or score). |
| KMOHDEAKEFG | `[]any` | Unknown game field, currently unmapped. |

### MedalList (Badges)

| Name | Type | Description |
| :--- | :--- | :--------- |
| MedalType | `int` | Badge Category/Type (e.g., Shiyu Defense, Hollow Zero) (check [store/zzz/medals.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/medals.json)). |
| Value | `int` | Progress number or tier associated with the badge. |
| MedalIcon | `int` | ID of the visual icon used for the badge (check [store/zzz/medals.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/medals.json)). | 
| MedalScore | `int` | Additional score metric for the badge (e.g., Endless Tower points or specific challenge score). |

---

## Characters (AvatarData)

### ShowcaseDetail

| Name | Type | Description |
| :--- | :--- | :--------- | 
| [AvatarList](#avatarlist) | `[]AvatarData` | List of detailed agent configurations (stats, weapons, discs) currently visible on the player's profile showcase. |

### AvatarList (AvatarData)

> [!TIP]
> The API returns **base** skill levels (how much the player invested manually). To get the true in-game levels, developers must manually add bonuses based on the agent's Mindscape (`TalentLevel`). If the agent has `TalentLevel >= 3`, manually increase the affected skill levels by 2. If `TalentLevel >= 5`, increase the remaining affected skill levels by 2.

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Id | `int` | Unique Agent ID (check [store/zzz/avatars.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/avatars.json)). |
| Exp | `int` | Current experience points of the agent towards their next level. |
| Level | `int` | The agent's current overall level. |
| PromotionLevel | `int` | The agent's ascension phase, which dictates their maximum level cap. |
| TalentLevel | `int` | Number of unlocked Mindscape Cinema nodes (duplicates obtained) for this agent. |
| SkinId | `int` | ID of the equipped outfit/skin for the agent. |
| CoreSkillEnhancement | `int` | The level of the agent's Core Passive skill enhancement (nodes A through F represented as an integer). |
| TalentToggleList | `[]bool` | Visual toggles for Mindscape Cinema effects. |
| WeaponEffectState | `int` | Indicates if the equipped W-Engine's special visual effect is active `[0: None, 1: OFF, 2: ON]`. |
| IsHidden | `bool` | Whether the agent's detailed stats are hidden from public view in the showcase. |
| IsFavorite | `bool` | Whether the player has marked this agent as a favorite. |
| ClaimedRewardList | `[]int` | IDs representing the ascension rewards the player has claimed for this agent. |
| ObtainmentTimestamp | `int64` | Unix timestamp of when the player first obtained this agent. |
| WeaponUid | `int` | Unique instance ID of the W-Engine equipped on this agent. |
| [Weapon](#w-engines-weapon) | `*Weapon` | Detailed data about the equipped W-Engine. | 
| SkillLevelList | `[]SkillLevel` | Array detailing the investment levels for the agent's combat skills. |
| [EquippedList](#drive-discs-equipment) | `[]EquippedItem` | Array of Drive Discs currently equipped on the agent. |
| IsUpgradeUnlocked | `bool` | Indicates whether the agent's Potential Vision mechanic is unlocked. |
| UpgradeId | `int` | The specific Potential Vision upgrade node/buff ID currently active on the agent. |

### SkillLevelList

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Level | `int` | Skill Level |
| Index | `int` | Skill Index (Check [Skills](#skills) for definitions) |

---

## W-Engines (Weapon)

> [!NOTE]
> **`Id` vs `Uid`**: `Id` is the static identifier of the item type (e.g., `14102` for Steel Cushion). The `Uid` field is unique to a *specific instance* of that item on the player's account. If two agents in a showcase share a W-Engine with the exact same `Uid`, it means the player swapped a single item between them. This allows developers to deduplicate the player's inventory.

For base stats, refer to [store/zzz/weapons.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/refs/heads/master/store/zzz/weapons.json). See [Formulas](#formulas) to calculate actual values.

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Uid | `int` | Unique instance ID of this specific W-Engine on the player's account. |
| Id | `int` | Static item type ID of the W-Engine. |
| Exp | `int` | Current experience points invested into the W-Engine. |
| Level | `int` | The W-Engine's current level. |
| BreakLevel | `int` | The refinement/modification level of the W-Engine (obtained by feeding duplicate copies). |
| UpgradeLevel | `int` | The ascension phase of the W-Engine, determining its level cap. |
| IsAvailable | `bool` | Whether the W-Engine is currently usable. |
| IsLocked | `bool` | Whether the player has locked the W-Engine in their inventory to prevent accidental deletion. |

---

## Drive Discs (Equipment)

### EquippedList

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Slot | `int` | Slot index (1-6) |
| [Equipment](#equipment) | `*Equipment` | Drive Disc data |

### Equipment

> [!NOTE]
> **`Id` vs `Uid`**: Just like W-Engines, `Id` is the static item type, while `Uid` identifies the specific instance of that disc on the account. Use `Uid` to deduplicate discs if a player swaps them between agents.
> 
> Rarity of drive discs can be found in [store/zzz/equipment.json](https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/equipments.json).

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Uid | `int` | Unique instance ID of this specific Drive Disc on the player's account. |
| Id | `int` | Static item type ID of the Drive Disc. |
| Exp | `int` | Current experience points invested into the Drive Disc. |
| Level | `int` | The Drive Disc's current upgrade level `[0-15]`. |
| BreakLevel | `int` | The total number of times random substats were upgraded as the Drive Disc was leveled up. |
| IsLocked | `bool` | Whether the player has locked the Drive Disc to prevent accidental deletion. |
| IsAvailable | `bool` | Whether the Drive Disc is currently usable. |
| IsTrash | `bool` | Whether the player has marked this Drive Disc as junk/trash in their inventory. |
| [MainPropertyList](#property) | `[]Property` | The primary, guaranteed stat of the Drive Disc. |
| [RandomPropertyList](#property) | `[]Property` | The list of randomly rolled substats on the Drive Disc. |


### Property

See [Formulas](#formulas) to calculate actual stat values based on these properties.

| Name | Type | Description |
| :--- | :--- | :---------- | 
| PropertyId | `int` | Property ID (check [Definitions](#property-id)) |
| PropertyValue | `int` | Property Base Value |
| PropertyLevel | `int` | Amount of rolls, only matters if substat |

---

## Definitions

### Property Id

Refer to the table below and [store/zzz/property.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/refs/heads/master/store/zzz/property.json) for static data mapping.

| ID | Description | ID | Description |
| :--- | :--- | :--- | :--- |
| `11101` | HP `[Base]` | `30501` | Energy Regen `[Base]` |
| `11102` | HP% | `30502` | Energy Regen% |
| `11103` | HP `[Flat]` | `30503` | Energy Regen `[Flat]` |
| `12101` | ATK `[Base]` | `31201` | Anomaly Proficiency `[Base]` |
| `12102` | ATK% | `31203` | Anomaly Proficiency `[Flat]` |
| `12103` | ATK `[Flat]` | `31401` | Anomaly Mastery `[Base]` |
| `12201` | Impact `[Base]` | `31402` | Anomaly Mastery% |
| `12202` | Impact% | `31403` | Anomaly Mastery `[Flat]` |
| `13101` | Def `[Base]` | `31501` | Physical DMG Bonus `[Base]` |
| `13102` | Def% | `31503` | Physical DMG Bonus `[Flat]` |
| `13103` | Def `[Flat]` | `31601` | Fire DMG Bonus `[Base]` |
| `20101` | Crit Rate `[Base]` | `31603` | Fire DMG Bonus `[Flat]` |
| `20103` | Crit Rate `[Flat]` | `31701` | Ice DMG Bonus `[Base]` |
| `21101` | Crit DMG `[Base]` | `31703` | Ice DMG Bonus `[Flat]` |
| `21103` | Crit DMG `[Flat]` | `31801` | Electric DMG Bonus `[Base]` |
| `23101` | Pen Ratio `[Base]` | `31803` | Electric DMG Bonus `[Flat]` |
| `23103` | Pen Ratio `[Flat]` | `31901` | Ether DMG Bonus `[Base]` |
| `23201` | PEN  `[Base]` | `31903` | Ether DMG Bonus `[Flat]` |
| `23203` | PEN `[Flat]` | | |

### Badge Type 

| Type | Description |
| :--- | :---------- |
| `1` | Shiyu Defense |
| `2` | Simulated Battle Tower |
| `3` | Deadly Assault |
| `4` | Simulated Battle Tower - Last Stand | 

### Skills

| Index | Description|
| :--- | :--------- |
| `0` | Basic Attack |
| `1` | Special Attack |
| `2` | Dash  |
| `3` | Ultimate |
| `5` | Core Skill |
| `6` | Assist |

---

## Formulas

> [!IMPORTANT]
> While working with character, weapon, and disc stats, you must use these formulas to calculate their final values.

### Agent Stats  

To calculate the base stats of an Agent, you need to use [store/zzz/avatars.json](https://github.com/EnkaNetwork/API-docs/blob/master/store/zzz/avatars.json).  

> [!TIP]
> It is recommended to `math.Floor` results before summing them up with stats from other sources.

```text
BaseTotalValue = BaseProps[PropertyId] + GrowthValue + PromotionValue + CoreEnhancementValue

GrowthValue          = (GrowthProps[PropertyId] * (Avatar.Level - 1)) / 10000
PromotionValue       = PromotionProps[Avatar.PromotionLevel - 1][PropertyId]
CoreEnhancementValue = CoreEnhancementProps[Avatar.CoreSkillEnhancement][PropertyId]
```

### Game-Accurate Formulas

#### W-Engine  

To calculate exact W-Engine stats, you need to use:  
- [WeaponLevelTemplateTb.json](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/FileCfg/WeaponLevelTemplateTb.json)  
- [WeaponStarTemplateTb.json](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/FileCfg/WeaponStarTemplateTb.json)  

**Main Stat:**  
```text
Result = MainStat.PropertyValue * (1 + WeaponLevel.FIELD_XXX / 10000 + WeaponStar.FIELD_YYY / 10000)

// Example (Steel Cushion Level 60, BreakLevel 5): 
// 684 = 46 * (1 + 94090 / 10000 + 44610 / 10000)
```

**Secondary Stat:**  
```text
Result = SecondaryStat.PropertyValue * (1 + WeaponStar.FIELD_ZZZ / 10000)

// Example (BreakLevel 5): 
// 2400 = 960 * (1 + 15000 / 10000)
```

#### Drive Disc  

To calculate exact Drive Disc stats, you need to use [EquipmentLevelTemplateTb.json](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/FileCfg/EquipmentLevelTemplateTb.json). This file determines the Drive Disc value based on its level and rarity.  

**Main Stat:**  
```text
Result = MainStat.PropertyValue * (1 + EquipmentLevel.Field_XXX / 10000)

// Example (Level 14, Rarity 4 [S-rank]): 
// 2090 = 550 * (1 + 28000 / 10000)
```

## Icons and Localizations

All the icon names and localized text mappings are included in the parsed data at [API-docs/store/zzz](https://github.com/EnkaNetwork/API-docs/tree/master/store/zzz).

- **Localizations:** For the names used in Enka.Network, refer to [store/zzz/locs.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/refs/heads/master/store/zzz/locs.json).
- **Extensive Text Data:** For any additional info about names, descriptions, etc., check the [TextMap Data](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/TextMap) in the ZenlessData repository (only includes languages supported by the game).
- **Icons & Graphics:** For additional image info, check the [ZenlessData](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/) repository.