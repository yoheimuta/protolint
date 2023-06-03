package config

// MessageNamesExcludePrepositionsOption represents the option for the MESSAGE_NAMES_EXCLUDE_PREPOSITIONS rule.
type MessageNamesExcludePrepositionsOption struct {
	CustomizableSeverityOption
	Prepositions []string `yaml:"prepositions"`
	Excludes     []string `yaml:"excludes"`
}
