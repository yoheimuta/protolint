package rules

import "github.com/yoheimuta/protolint/linter/rule"

// RuleWithSeverity represents a rule with a configurable severity.
type RuleWithSeverity struct {
	severity rule.Severity
}

// NewRuleWithSeverity takes a severity and adds it to a new instance
// of RuleWithSeverity
func NewRuleWithSeverity(
	severity rule.Severity,
) RuleWithSeverity {
	return RuleWithSeverity{severity: severity}
}

// Severity returns the configured severity.
func (r RuleWithSeverity) Severity() rule.Severity {
	return r.severity
}
