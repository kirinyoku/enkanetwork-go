package models

// AvatarDetail contains detailed character information for Honkai: Star Rail
type AvatarDetail struct {
	Pos           int         `json:"pos,omitempty"`           // Character position
	Rank          int         `json:"rank,omitempty"`          // Character rank
	Level         int         `json:"level,omitempty"`         // Character level
	Assist        int         `json:"_assist,omitempty"`       // Assist status
	AvatarID      int         `json:"avatarId,omitempty"`      // Character ID
	Equipment     *Equipment  `json:"equipment,omitempty"`     // Equipped equipment
	Promotion     int         `json:"promotion,omitempty"`     // Promotion level
	RelicList     []Relic     `json:"relicList,omitempty"`     // List of equipped relics
	SkillTreeList []SkillTree `json:"skillTreeList,omitempty"` // List of skill tree nodes
}

// SkillTree represents a skill tree node for a character in Honkai: Star Rail
type SkillTree struct {
	Level   int `json:"level,omitempty"`   // Node level
	PointID int `json:"pointId,omitempty"` // Node point ID
}
