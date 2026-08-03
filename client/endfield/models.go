package endfield

import (
	"encoding/json"
	"time"

	"github.com/kirinyoku/enkanetwork-go/internal/core"
	"github.com/kirinyoku/enkanetwork-go/models"
)

// Profile represents the root structure of the response containing player information
// for Arknights Endfield.
type Profile struct {
	// PlayerInfo contains detailed information about the game account and characters
	PlayerInfo *PlayerInfo `json:"playerInfo,omitempty"`
	// TTL indicates the seconds remaining until the next request to the game. Until
	// the TTL expires, the endpoint returns cached data — but such requests still
	// count toward the rate limit
	TTL int `json:"ttl,omitempty"`
	// Owner is the Enka profile associated with the provided UID. The response includes
	// an Owner if:
	//   1. The user has an account on the site;
	//   2. The user has added their UID to their profile;
	//   3. The user has verified that the UID belongs to them;
	//   4. The user has set their profile visibility to "public"
	Owner *models.Owner `json:"owner,omitempty"`
	// UID is the player's UID in Arknights Endfield.
	UID string `json:"uid,omitempty"`
	// Region is the player's server region.
	Region string `json:"region,omitempty"`
	// Raw contains the original API response.
	Raw json.RawMessage `json:"-"`
	// Extra contains unknown top-level API fields preserved during unmarshaling.
	Extra map[string]json.RawMessage `json:"-"`
}

func (p *Profile) UnmarshalJSON(data []byte) error {
	type Alias Profile
	var profile Alias
	if err := json.Unmarshal(data, &profile); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	profile.Raw = raw
	profile.Extra = extra
	*p = Profile(profile)
	return nil
}

func (p Profile) MarshalJSON() ([]byte, error) {
	type Alias Profile
	return core.MergeKnownExtraAndRawJSON(Alias(p), p.Extra, p.Raw)
}

func (p Profile) CacheTTL() time.Duration {
	return time.Duration(p.TTL) * time.Second
}

// PlayerInfo contains information about the game account and the characters on display.
type PlayerInfo struct {
	BusinessCard *BusinessCard              `json:"businessCard,omitempty"`
	CharData     []CharData                 `json:"charData"`
	Raw          json.RawMessage            `json:"-"`
	Extra        map[string]json.RawMessage `json:"-"`
}

func (p *PlayerInfo) UnmarshalJSON(data []byte) error {
	type Alias PlayerInfo
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	alias.Raw = raw
	alias.Extra = extra
	*p = PlayerInfo(alias)
	return nil
}

func (p PlayerInfo) MarshalJSON() ([]byte, error) {
	type Alias PlayerInfo
	return core.MergeKnownExtraAndRawJSON(Alias(p), p.Extra, p.Raw)
}

// BusinessCard contains the player's social and profile settings.
type BusinessCard struct {
	Name                          string                         `json:"name"`
	ShortID                       string                         `json:"shortId"`
	Gender                        int                            `json:"gender"`
	BusinessCardTopicID           int                            `json:"businessCardTopicId"`
	Signature                     string                         `json:"signature"`
	UserAvatarID                  int                            `json:"userAvatarId"`
	UserAvatarFrameID             int                            `json:"userAvatarFrameId"`
	AdventureLevel                int                            `json:"adventureLevel"`
	WorldLevel                    int                            `json:"worldLevel"`
	MainMissionID                 string                         `json:"mainMissionId"`
	CreateTime                    int64                          `json:"createTime"`
	PlatformRoleID                string                         `json:"platformRoleId"`
	DomainDev                     *DomainDev                     `json:"domainDev,omitempty"`
	Achievement                   *Achievement                   `json:"achievement,omitempty"`
	Statistic                     *Statistic                     `json:"statistic,omitempty"`
	CharList                      []CharListEntry                `json:"charList"`
	ContingencyContractBestRecord *ContingencyContractBestRecord `json:"contingencyContractBestRecord,omitempty"`
	BusinessCardExpandFlag        bool                           `json:"businessCardExpandFlag"`
	Raw                           json.RawMessage                `json:"-"`
	Extra                         map[string]json.RawMessage     `json:"-"`
}

func (b *BusinessCard) UnmarshalJSON(data []byte) error {
	type Alias BusinessCard
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	alias.Raw = raw
	alias.Extra = extra
	*b = BusinessCard(alias)
	return nil
}

func (b BusinessCard) MarshalJSON() ([]byte, error) {
	type Alias BusinessCard
	return core.MergeKnownExtraAndRawJSON(Alias(b), b.Extra, b.Raw)
}

type DomainDev struct {
	Domains []Domain `json:"domains"`
}

type Domain struct {
	DomainID string `json:"domainId"`
	Level    int    `json:"level"`
}

type Achievement struct {
	Display  []AchievementDisplay `json:"display"`
	InfoList []AchievementInfo    `json:"infoList"`
}

type AchievementDisplay struct {
	Key   int `json:"key"`
	Value int `json:"value"`
}

type AchievementInfo struct {
	AchieveNumID int  `json:"achieveNumId"`
	Level        int  `json:"level"`
	IsPlated     bool `json:"isPlated"`
}

type Statistic struct {
	CharNum   int `json:"charNum"`
	WeaponNum int `json:"weaponNum"`
	DocNum    int `json:"docNum"`
}

type CharListEntry struct {
	TemplateID     string `json:"templateId"`
	Level          int    `json:"level"`
	PotentialLevel int    `json:"potentialLevel"`
}

type ContingencyContractBestRecord struct {
	ActivityID int `json:"activityId"`
	BestScore  int `json:"bestScore"`
}

// CharData contains detailed information about a specific character build.
type CharData struct {
	TemplateID      int                        `json:"templateId"`
	Level           int                        `json:"level"`
	Exp             int                        `json:"exp"`
	PotentialLevel  int                        `json:"potentialLevel"`
	Equip           []EquipEntry               `json:"equip"`
	Weapon          *Weapon                    `json:"weapon,omitempty"`
	SkillInfo       *SkillInfo                 `json:"skillInfo,omitempty"`
	EquipMedicineID int                        `json:"equipMedicineId"`
	Talent          *Talent                    `json:"talent,omitempty"`
	PotentialCg     []any                      `json:"potentialCg"`
	Raw             json.RawMessage            `json:"-"`
	Extra           map[string]json.RawMessage `json:"-"`
}

func (c *CharData) UnmarshalJSON(data []byte) error {
	type Alias CharData
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	alias.Raw = raw
	alias.Extra = extra
	*c = CharData(alias)
	return nil
}

func (c CharData) MarshalJSON() ([]byte, error) {
	type Alias CharData
	return core.MergeKnownExtraAndRawJSON(Alias(c), c.Extra, c.Raw)
}

type EquipEntry struct {
	Key   int    `json:"key"`
	Value *Equip `json:"value,omitempty"`
}

type Equip struct {
	TemplateID int   `json:"templateid"`
	Enhance    []any `json:"enhance"`
}

type Weapon struct {
	TemplateID     int                        `json:"templateId"`
	Exp            int                        `json:"exp"`
	WeaponLv       int                        `json:"weaponLv"`
	RefineLv       int                        `json:"refineLv"`
	BreakthroughLv int                        `json:"breakthroughLv"`
	AttachedGem    *AttachedGem               `json:"attachedGem,omitempty"`
	Raw            json.RawMessage            `json:"-"`
	Extra          map[string]json.RawMessage `json:"-"`
}

func (w *Weapon) UnmarshalJSON(data []byte) error {
	type Alias Weapon
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	alias.Raw = raw
	alias.Extra = extra
	*w = Weapon(alias)
	return nil
}

func (w Weapon) MarshalJSON() ([]byte, error) {
	type Alias Weapon
	return core.MergeKnownExtraAndRawJSON(Alias(w), w.Extra, w.Raw)
}

type AttachedGem struct {
	TemplateID int       `json:"templateId"`
	TotalCost  int       `json:"totalCost"`
	Terms      []GemTerm `json:"terms"`
	DomainID   int       `json:"domainId"`
}

type GemTerm struct {
	TermNumID int `json:"termNumId"`
	Cost      int `json:"cost"`
}

type SkillInfo struct {
	LevelInfo             []SkillLevelInfo           `json:"levelInfo"`
	NormalSkill           string                     `json:"normalSkill"`
	UltimateSkill         string                     `json:"ultimateSkill"`
	ComboSkill            string                     `json:"comboSkill"`
	DispNormalAttackSkill string                     `json:"dispNormalAttackSkill"`
	Raw                   json.RawMessage            `json:"-"`
	Extra                 map[string]json.RawMessage `json:"-"`
}

func (s *SkillInfo) UnmarshalJSON(data []byte) error {
	type Alias SkillInfo
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	alias.Raw = raw
	alias.Extra = extra
	*s = SkillInfo(alias)
	return nil
}

func (s SkillInfo) MarshalJSON() ([]byte, error) {
	type Alias SkillInfo
	return core.MergeKnownExtraAndRawJSON(Alias(s), s.Extra, s.Raw)
}

type SkillLevelInfo struct {
	SkillID            string `json:"skillId"`
	SkillLevel         int    `json:"skillLevel"`
	SkillMaxLevel      int    `json:"skillMaxLevel"`
	SkillEnhancedLevel int    `json:"skillEnhancedLevel"`
}

type Talent struct {
	LatestBreakNode         string                     `json:"latestBreakNode"`
	AttrNodes               []string                   `json:"attrNodes"`
	LatestPassiveSkillNodes []string                   `json:"latestPassiveSkillNodes"`
	LatestFactorySkillNodes []string                   `json:"latestFactorySkillNodes"`
	Raw                     json.RawMessage            `json:"-"`
	Extra                   map[string]json.RawMessage `json:"-"`
}

func (t *Talent) UnmarshalJSON(data []byte) error {
	type Alias Talent
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	raw, extra, err := core.PreserveUnknownJSONForStruct(data, Alias{})
	if err != nil {
		return err
	}
	alias.Raw = raw
	alias.Extra = extra
	*t = Talent(alias)
	return nil
}

func (t Talent) MarshalJSON() ([]byte, error) {
	type Alias Talent
	return core.MergeKnownExtraAndRawJSON(Alias(t), t.Extra, t.Raw)
}
