package config

import "github.com/yoheimuta/protolint/linter/rule"

type CustomizableSeverityOption struct {
	severity *rule.Severity `yaml:"severity"`
}

func (c CustomizableSeverityOption) Severity() rule.Severity {
	if c.severity == nil {
		return rule.Severity_Error
	}

	return *c.severity
}
