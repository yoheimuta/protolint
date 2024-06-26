package config

// FieldNamesExcludePrepositionsOption represents the option for the FIELD_NAMES_EXCLUDE_PREPOSITIONS rule.
type FieldNamesExcludePrepositionsOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	Prepositions               []string `yaml:"prepositions" json:"prepositions" toml:"prepositions"`
	Excludes                   []string `yaml:"excludes" json:"excludes" toml:"excludes"`
}
