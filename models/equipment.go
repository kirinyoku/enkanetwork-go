package models

// Equipment contains information about a character’s equipment in Honkai: Star Rail.
type Equipment struct {
	TID       int   `json:"tid,omitempty"`       // Equipment type ID
	Rank      int   `json:"rank,omitempty"`      // Equipment rank
	Flat      *Flat `json:"_flat,omitempty"`     // Detailed equipment data
	Level     int   `json:"level,omitempty"`     // Equipment level
	Promotion int   `json:"promotion,omitempty"` // Promotion level
}

// Flat contains detailed metadata for equipment or relics in Honkai: Star Rail.
type Flat struct {
	Props   []Prop `json:"props,omitempty"`   // List of properties
	Name    string `json:"name,omitempty"`    // Equipment or relic name
	SetID   int    `json:"setID,omitempty"`   // Set ID
	SetName string `json:"setName,omitempty"` // Set name
}

// Prop represents a property of equipment or relics.
type Prop struct {
	Type  string  `json:"type,omitempty"`  // Property type
	Value float64 `json:"value,omitempty"` // Property value
}

// Relic contains information about a relic in Honkai: Star Rail.
type Relic struct {
	TID          int        `json:"tid,omitempty"`          // Relic type ID
	Type         int        `json:"type,omitempty"`         // Relic type
	Flat         *Flat      `json:"_flat,omitempty"`        // Detailed relic data
	Level        int        `json:"level,omitempty"`        // Relic level
	MainAffixID  int        `json:"mainAffixId,omitempty"`  // Main affix ID
	SubAffixList []SubAffix `json:"subAffixList,omitempty"` // List of sub-affixes
}

// SubAffix represents a sub-affix of a relic in Honkai: Star Rail
type SubAffix struct {
	Cnt     int `json:"cnt,omitempty"`     // Sub-affix count
	Step    int `json:"step,omitempty"`    // Sub-affix step
	AffixID int `json:"affixId,omitempty"` // Sub-affix ID
}
