package rules

import "github.com/yoheimuta/protolint/linter/rule"

type RuleWithSeverity struct {
	severity rule.Severity
}

func NewRuleWithSeverity(
	severity rule.Severity,
) RuleWithSeverity {
	return RuleWithSeverity{severity: severity}
}

func (r RuleWithSeverity) Severity() rule.Severity {
	return r.severity
}
