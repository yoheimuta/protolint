package config

// MessageNamesExcludePrepositionsOption represents the option for the MESSAGE_NAMES_EXCLUDE_PREPOSITIONS rule.
type MessageNamesExcludePrepositionsOption struct {
	Prepositions []string `yaml:"prepositions"`
}
