package config

import "github.com/yoheimuta/protolint/linter/rule"

// CustomizableSeverityOption represents an option where the
// severity of a rule can be configured via yaml.
type CustomizableSeverityOption struct {
	severity *rule.Severity `yaml:"severity" toml:"severity"`
}

// Severity returns the configured severity. If no severity
// is set, the default severity will be ERROR
func (c CustomizableSeverityOption) Severity() rule.Severity {
	if c.severity == nil {
		return rule.SeverityError
	}

	return *c.severity
}
