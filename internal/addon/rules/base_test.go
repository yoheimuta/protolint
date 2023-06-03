package rules_test

import (
	"testing"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/rule"
)

func TestRulesWithSeverityHasSeverity(t *testing.T) {
	tests :=
		[]rule.Severity{
			rule.Severity_Note,
			rule.Severity_Warning,
			rule.Severity_Error,
		}
	for _, test := range tests {
		test := test
		t.Run(string(test), func(t *testing.T) {
			var rule_to_test rule.HasSeverity
			rule_to_test = rules.NewRuleWithSeverity(test)
			if rule_to_test.Severity() != test {
				t.Errorf("Rule should have %v severity, but got %v", test, rule_to_test.Severity())
			}
		})
	}
}
