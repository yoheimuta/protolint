package config

// SyntaxConsistentOption represents the option for the SYNTAX_CONSISTENT rule.
type SyntaxConsistentOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	Version                    string `yaml:"version" json:"version" toml:"version"`
}
