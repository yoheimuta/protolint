package config

import "github.com/yoheimuta/protolint/linter/rule"

// CustomizableSeverityOption represents an option where the
// severity of a rule can be configured via yaml.
type CustomizableSeverityOption struct {
	Severity rule.Severity `yaml:"severity" toml:"severity"`
}
