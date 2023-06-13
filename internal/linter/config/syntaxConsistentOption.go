package config

// SyntaxConsistentOption represents the option for the SYNTAX_CONSISTENT rule.
type SyntaxConsistentOption struct {
	CustomizableSeverityOption
	Version string `yaml:"version"`
}
