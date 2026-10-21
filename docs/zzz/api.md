# Enka.Network API - Zenless Zone Zero

This document describes the structure and formulas for the Zenless Zone Zero API response provided by Enka.Network. 

## Table of Content

- [Data Structure Overview](#data-structure-overview)
- [Profile & Player Info](#profile--player-info)
- [Characters (AvatarData)](#characters-avatardata)
- [W-Engines (Weapon)](#w-engines-weapon)
- [Drive Discs (Equipment)](#drive-discs-equipment)
- [Definitions](#definitions)
  - [Property ID](#property-id)
  - [Badge Type](#badge-type)
  - [Skills](#skills)
- [Formulas](#formulas)
  - [Agent Base Stats](#agent-base-stats)
  - [Total Combat Stats Pipeline](#total-combat-stats-pipeline)
  - [Rounding Rules](#rounding-rules)
  - [W-Engine Stats](#w-engine-stats)
  - [Drive Disc Stats](#drive-disc-stats)
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
| [MedalList](#medallist-badges) | `[]Medal` | List of in-game badges the player has chosen to display on their profile. | 

### ProfileDetail

| Name | Type | Description |
| :--- | :--- | :--------- |
| Uid | `int64` | The player's unique Zenless Zone Zero account identifier. |
| Nickname | `string` | The display name chosen by the player to represent their account. |
| ProfileId | `int` | The unique ID of the player's selected avatar image or profile picture (check [store/zzz/pfps.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/pfps.json)). |
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
| MedalType | `int` | Badge Category/Type (e.g., Shiyu Defense, Hollow Zero) (check [store/zzz/medals.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/medals.json) and [Badge Type](#badge-type)). |
| Value | `int` | Progress number or tier associated with the badge. |
| MedalIcon | `int` | ID of the visual icon used for the badge (check [store/zzz/medals.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/medals.json)). | 
| MedalScore | `*int` | Additional score metric for the badge (e.g., Endless Ascension points or specific challenge score). |

---

## Characters (AvatarData)

### ShowcaseDetail

| Name | Type | Description |
| :--- | :--- | :--------- | 
| [AvatarList](#avatarlist-avatardata) | `[]AvatarData` | List of detailed agent configurations (stats, weapons, discs) currently visible on the player's profile showcase. |

### AvatarList (AvatarData)

> [!TIP]
> The API returns **base** skill levels (how much the player invested manually). To get the true in-game levels, developers must manually add bonuses based on the agent's Mindscape Cinema (`TalentLevel`).
> - `TalentLevel >= 3` (Cinema 3 / M3): Adds +2 levels to designated active skills (e.g., Basic Attack, Dodge, Assist, Special, Chain).
> - `TalentLevel >= 5` (Cinema 5 / M5): Adds +2 levels to the remaining active combat skills.

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Id | `int` | Unique Agent ID (check [store/zzz/avatars.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/avatars.json)). |
| Exp | `int` | Current experience points of the agent towards their next level. |
| Level | `int` | The agent's current overall level `[1-60]`. |
| PromotionLevel | `int` | The agent's ascension phase `[0-5]`, which dictates their maximum level cap (Level 10/20/30/40/50/60). |
| TalentLevel | `int` | Number of unlocked Mindscape Cinema nodes `[0-6]` (duplicates obtained) for this agent. |
| SkinId | `int` | ID of the equipped outfit/skin for the agent. If `0`, the default skin is used. |
| CoreSkillEnhancement | `int` | The level of the agent's Core Passive skill enhancement `[0-6]`, representing Core Nodes A through F. |
| TalentToggleList | `[]bool` | Visual toggles for Mindscape Cinema effects. |
| WeaponEffectState | `int` | Indicates if the equipped W-Engine's special visual effect is active `[0: None, 1: OFF, 2: ON]`. |
| IsHidden | `bool` | Whether the agent's detailed stats are hidden from public view in the showcase. |
| IsFavorite | `bool` | Whether the player has marked this agent as a favorite. |
| ClaimedRewardList | `[]int` | IDs representing the ascension rewards the player has claimed for this agent. |
| ObtainmentTimestamp | `int64` | Unix timestamp of when the player first obtained this agent. |
| WeaponUid | `int` | Unique instance ID of the W-Engine equipped on this agent. |
| [Weapon](#w-engines-weapon) | `*Weapon` | Detailed data about the equipped W-Engine. | 
| [SkillLevelList](#skilllevellist) | `[]SkillLevel` | Array detailing the investment levels for the agent's combat skills. |
| [EquippedList](#drive-discs-equipment) | `[]EquippedItem` | Array of Drive Discs currently equipped on the agent. |
| IsUpgradeUnlocked | `*bool` | Indicates whether the agent's Potential Vision mechanic is unlocked (version 1.4+ mechanic). |
| UpgradeId | `int` | The specific Potential Vision upgrade node/buff ID currently active on the agent. |

### SkillLevelList

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Level | `int` | Base Skill Level (before Mindscape Cinema bonuses). |
| Index | `int` | Skill Slot Index (see [Skills](#skills) for index definitions). |

---

## W-Engines (Weapon)

> [!NOTE]
> **`Id` vs `Uid`**: `Id` is the static identifier of the item type (e.g., `14102` for Steel Cushion). The `Uid` field is unique to a *specific instance* of that item on the player's account. If two agents in a showcase share a W-Engine with the exact same `Uid`, it means the player swapped a single item between them. This allows developers to deduplicate the player's inventory.

For base stats, refer to [store/zzz/weapons.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/weapons.json). See [Formulas](#formulas) to calculate actual values.

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Uid | `int` | Unique instance ID of this specific W-Engine on the player's account. |
| Id | `int` | Static item type ID of the W-Engine. |
| Exp | `int` | Current experience points invested into the W-Engine. |
| Level | `int` | The W-Engine's current level `[1-60]`. |
| BreakLevel | `int` | The ascension phase of the W-Engine `[0-5]`, determining its level cap (10/20/30/40/50/60). |
| UpgradeLevel | `int` | The refinement / modification level `[1-5]` (M1–M5), obtained by feeding duplicate copies. Scales the passive effect and secondary stat. |
| IsAvailable | `bool` | Whether the W-Engine is currently usable. |
| IsLocked | `bool` | Whether the player has locked the W-Engine in their inventory to prevent accidental deletion. |

---

## Drive Discs (Equipment)

### EquippedList

| Name | Type | Description |
| :--- | :--- | :--------- | 
| Slot | `int` | Partition slot index `[1-6]`. |
| [Equipment](#equipment) | `*Equipment` | Drive Disc data. |

### Equipment

> [!NOTE]
> **`Id` vs `Uid`**: Just like W-Engines, `Id` is the static item type, while `Uid` identifies the specific instance of that disc on the account. Use `Uid` to deduplicate discs if a player swaps them between agents.
> 
> Drive disc set details and rarities can be found in [store/zzz/equipments.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/equipments.json).

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

### Drive Disc Partition Slots

Each partition slot (1 through 6) has designated main stat possibilities:

| Slot | Fixed / Variable | Available Main Stats |
| :--- | :--- | :--- |
| **Partition 1** | Fixed | HP `[Flat]` |
| **Partition 2** | Fixed | ATK `[Flat]` |
| **Partition 3** | Fixed | DEF `[Flat]` |
| **Partition 4** | Variable | HP%, ATK%, DEF%, CRIT Rate, CRIT DMG, Anomaly Proficiency |
| **Partition 5** | Variable | HP%, ATK%, DEF%, PEN Ratio, Attribute DMG Bonus (Physical, Fire, Ice, Electric, Ether, Wind) |
| **Partition 6** | Variable | HP%, ATK%, DEF%, Impact%, Anomaly Mastery, Energy Regen% |

### Property

See [Formulas](#formulas) to calculate actual stat values based on these properties.

| Name | Type | Description |
| :--- | :--- | :---------- | 
| PropertyId | `int` | Property ID (check [Definitions](#property-id)). |
| PropertyValue | `int` | Property base value (for percent stats, divide by `10000.0`). |
| PropertyLevel | `int` | Upgrade roll count (1 for base roll, plus additional rolls for substat upgrades). |

---

## Definitions

### Property ID

Refer to the table below and [store/zzz/property.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/property.json) for static data mapping.

Property IDs follow a structured numbering scheme:
- `xx01`: Base stat / initial value
- `xx02`: Percentage modifier (Ratio)
- `xx03`: Flat stat modifier
- `xx04`: Additional percentage modifier
- `xx05`: Additional flat modifier / Final attribute bonus

| ID | Description | ID | Description |
| :--- | :--- | :--- | :--- |
| `11101` | HP `[Base]` | `30501` | Energy Regen `[Base]` |
| `11102` | HP% | `30502` | Energy Regen% |
| `11103` | HP `[Flat]` | `30503` | Energy Regen `[Flat]` |
| `11104` | HP% `[Bonus]` | `30601` | Healing Ratio `[Base]` |
| `11105` | HP `[Flat Bonus]` | `30603` | Healing Ratio `[Flat]` |
| `11901` | Max RP `[Base]` (Rupture) | `30701`-`31105` | All DMG Bonus (`AddedDamageRatio`) |
| `11903` | Max RP `[Flat]` | `31201` | Anomaly Proficiency `[Base]` |
| `12101` | ATK `[Base]` | `31202` | Anomaly Proficiency% |
| `12102` | ATK% | `31203` | Anomaly Proficiency `[Flat]` |
| `12103` | ATK `[Flat]` | `31301` | Anomaly Accumulation Ratio% |
| `12104` | ATK% `[Bonus]` | `31401` | Anomaly Mastery `[Base]` |
| `12105` | ATK `[Flat Bonus]` | `31402` | Anomaly Mastery% |
| `12201` | Impact `[Base]` | `31403` | Anomaly Mastery `[Flat]` |
| `12202` | Impact% | `31501` | Physical DMG Bonus `[Base]` |
| `12203` | Impact `[Flat]` | `31505` | Physical DMG Bonus `[Total]` |
| `12301` | Sheer Force `[Base]` (Rupture) | `31601` | Fire DMG Bonus `[Base]` |
| `12303` | Sheer Force `[Flat]` | `31605` | Fire DMG Bonus `[Total]` |
| `13101` | DEF `[Base]` | `31701` | Ice DMG Bonus `[Base]` |
| `13102` | DEF% | `31705` | Ice DMG Bonus `[Total]` |
| `13103` | DEF `[Flat]` | `31801` | Electric DMG Bonus `[Base]` |
| `20101` | Crit Rate `[Base]` | `31805` | Electric DMG Bonus `[Total]` |
| `20103` | Crit Rate `[Flat]` | `31901` | Ether DMG Bonus `[Base]` |
| `21101` | Crit DMG `[Base]` | `31905` | Ether DMG Bonus `[Total]` |
| `21103` | Crit DMG `[Flat]` | `32001` | Adrenaline Auto-Accumulation `[Base]` |
| `23101` | PEN Ratio `[Base]` | `32002` | Adrenaline Auto-Accumulation% |
| `23103` | PEN Ratio `[Flat]` | `32003` | Adrenaline Auto-Accumulation `[Flat]` |
| `23201` | Flat PEN `[Base]` | `32201` | Skip DEF DMG Ratio `[Base]` |
| `23203` | Flat PEN `[Flat]` | `32301` | Wind DMG Bonus `[Base]` |
| | | `32305` | Wind DMG Bonus `[Total]` |

### Badge Type 

Badge categories mapped from in-game medal definitions (see [store/zzz/medals.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/medals.json)):

| Type | Name |
| :--- | :--- |
| `1` | Shiyu Defense: Waypoint / Adaption |
| `2` | Endless Ascension |
| `3` | Deadly Assault: Line Breaker |
| `4` | Endless Last Stand |
| `5` | Blessed by Glory |
| `6` | Shiyu Defense: Lingering Sun |
| `7` | Deadly Assault: Disintegration |
| `8` | Withering Garden / Hollow Zero |
| `9` | Hollow Zero Exploration |
| `10` | Crowned by Waves |
| `11` | Annihilation |
| `12` | Simulated Annihilation |
| `13` | Labyrinth Fortune |
| `14` | Deadly Assault: Purge |
| `15` | Deadly Assault: Annihilator |

### Skills

The `SkillLevelList` array maps to specific combat abilities using these slot indices:

| Index | Skill Category | In-Game Abilities |
| :--- | :--- | :--- |
| `0` | Basic Attack | Basic Attack combos |
| `1` | Dodge | Dodge, Dash Attack, Dodge Counter |
| `2` | Assist | Quick Assist, Defensive Assist, Evasive Assist, Assist Follow-Up |
| `3` | Special Attack | Special Attack, EX Special Attack |
| `4` / `6` | Chain Attack & Ultimate | Chain Attack, Ultimate (Primary slot is `6`, Alternate slot is `4`) |

> [!NOTE]
> **Core Passive / Skill Enhancements** are not part of `SkillLevelList`. The Core Passive level is stored directly in `AvatarData.CoreSkillEnhancement` `[0-6]` (representing Core nodes A through F).

---

## Formulas

> [!IMPORTANT]
> When calculating final combat stats for display, apply stat growth, W-Engine scaling, Drive Disc properties, and set bonuses in the specified order.

### Agent Base Stats

To calculate the innate base stats of an Agent, use [store/zzz/avatars.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/avatars.json):

```text
AvatarBaseStat = BaseProps[PropertyId] + GrowthValue + PromotionValue + CoreEnhancementValue

GrowthValue          = (GrowthProps[PropertyId] * (Avatar.Level - 1)) / 10000
PromotionValue       = PromotionProps[Avatar.PromotionLevel - 1][PropertyId]
CoreEnhancementValue = CoreEnhancementProps[Avatar.CoreSkillEnhancement][PropertyId]
```

### Total Combat Stats Pipeline

1. **Innate Base Stats**: Calculate the agent's base stats using the formula above.
2. **Base ATK Addition**: W-Engine Base ATK is added directly to the Agent's Base ATK before any ATK% multipliers:
   $$\text{BaseATK}_{\text{total}} = \lfloor \text{AvatarBaseATK} \rfloor + \lfloor \text{WEngineBaseATK} \rfloor$$
3. **Stat Aggregation**: Accumulate all percentage multipliers and flat bonuses from W-Engine secondary stats, Drive Disc main stats, Drive Disc substats, and 2-piece set bonuses.
4. **Final Combat Stats**:
   $$\text{TotalStat} = \lfloor \text{BaseStat} \times (1 + \sum \text{PercentBonuses}) + \sum \text{FlatBonuses} \rfloor$$
5. **Rupture Mechanics**: Agents with the Rupture specialty convert HP and ATK into Sheer Force:
   $$\text{SheerForce} = \lfloor \text{TotalHP} \times 0.1 \rfloor + \lfloor \text{TotalATK} \times 0.3 \rfloor + \sum \text{SheerForceBonuses}$$

### Rounding Rules

- **Floored Stats (`math.Floor`)**:
  - `HP`, `ATK`, `DEF`, `Impact`, `Anomaly Proficiency`, `Anomaly Mastery`, `Flat PEN`, `Sheer Force`.
- **Exact Decimal Stats (NOT floored)**:
  - `CRIT Rate` (e.g., `5.0%` = `0.05`), `CRIT DMG` (e.g., `50.0%` = `0.50`), `PEN Ratio` (e.g., `24.0%` = `0.24`), `Attribute DMG Bonus` (e.g., `30.0%` = `0.30`), and `Energy Regen` (energy recovered per second, e.g., `1.20`).

---

### W-Engine Stats

To calculate exact W-Engine stats, refer to data from [ZenlessData](https://git.mero.moe/dimbreath/ZenlessData):
- [WeaponLevelTemplateTb.json](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/FileCfg/WeaponLevelTemplateTb.json) (Level growth multiplier `AHMDJCIHNKG` and substat divisor `IDBKOAPHGLC`)
- [WeaponStarTemplateTb.json](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/FileCfg/WeaponStarTemplateTb.json) (Ascension/Break star modifiers `NMFHJKEFLOG` and `FCLIIPBDDKP`)

#### Main Stat (Base ATK)

```text
Result = MainStat.PropertyValue * (1 + WeaponLevel.AHMDJCIHNKG / 10000 + WeaponStar.NMFHJKEFLOG / 10000)

// Example: Steel Cushion (S-Rank Level 60, BreakLevel 5 / Ascension 5):
// 684 = 46 * (1 + 94090 / 10000 + 44610 / 10000)
//     = 46 * (1 + 9.4090 + 4.4610) = 46 * 14.87 = 684
```

#### Secondary Stat

W-Engine secondary stats scale based on the level denominator `IDBKOAPHGLC` (`sub_stat_denominator`) and star phase:

$$\text{LevelMultiplier} = \frac{10000}{\text{WeaponLevel.IDBKOAPHGLC}}$$

```text
Result = SecondaryStat.PropertyValue * (10000 / WeaponLevel.IDBKOAPHGLC) * (1 + WeaponStar.FCLIIPBDDKP / 10000)

// Example: Steel Cushion (Level 60, BreakLevel 5 / Ascension 5, Base CRIT Rate 960 = 9.6%):
// Level 60 IDBKOAPHGLC = 4000 (Multiplier = 10000 / 4000 = 2.5)
// BreakLevel 5 FCLIIPBDDKP = 0
// 2400 (24.0%) = 960 * 2.5 * 1.0 = 2400
```

---

### Drive Disc Stats

To calculate exact Drive Disc stats, use [EquipmentLevelTemplateTb.json](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/FileCfg/EquipmentLevelTemplateTb.json) (Multiplier `JEKGLLBALFE` based on disc rarity and upgrade level `0-15`).

#### Main Stat

```text
Result = MainStat.PropertyValue * (1 + EquipmentLevel.JEKGLLBALFE / 10000)

// Example: S-Rank Disc (Rarity 4) at Level 15 (JEKGLLBALFE = 30000 = +300%):
// Level 15 Flat ATK: 316 = 79 * (1 + 30000 / 10000) = 79 * 4.0 = 316
// Level 15 CRIT Rate: 2400 (24.0%) = 600 * (1 + 30000 / 10000) = 600 * 4.0 = 2400
```

---

## Icons and Localizations

All icon names and localized text mappings are included in the parsed data at [API-docs/store/zzz](https://github.com/EnkaNetwork/API-docs/tree/master/store/zzz).

- **Localizations:** For names and descriptions used in Enka.Network, refer to [store/zzz/locs.json](https://raw.githubusercontent.com/EnkaNetwork/API-docs/master/store/zzz/locs.json).
- **TextMap Data:** For complete multi-language text dumps, check the [TextMap Data](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/TextMap) in the ZenlessData repository.
- **Icons & Graphics:** Assets can be fetched via Enka CDN (`https://enka.network/ui/zzz/<icon_name>.png`) or directly from [ZenlessData](https://git.mero.moe/dimbreath/ZenlessData/src/branch/master/).