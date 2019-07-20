package config

// FieldNamesExcludePrepositionsOption represents the option for the FIELD_NAMES_EXCLUDE_PREPOSITIONS rule.
type FieldNamesExcludePrepositionsOption struct {
	Prepositions []string `yaml:"prepositions"`
}
