package config

import "github.com/yoheimuta/protolint/linter/rule"

// CustomizableSeverityOption represents an option where the
// severity of a rule can be configured via yaml.
type CustomizableSeverityOption struct {
	// @ugly
	// In order to unmarshall from yaml the field in the struct must be public.
	// But we want to hide direct access to the field in order to
	// force the usage of function below which contains the fallback to the error severity.
	// This gives us a name clash so we added the internal suffix.
	// I am no go expert (yet) but i am sure that a more elegant solution exists
	SeverityInternal *rule.Severity `yaml:"severity" toml:"severity"`
}

// Severity returns the configured severity. If no severity
// is set, the default severity will be ERROR
func (c CustomizableSeverityOption) Severity() rule.Severity {
	if c.SeverityInternal == nil {
		return rule.SeverityError
	}

	return *c.SeverityInternal
}
