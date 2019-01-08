package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

var defaultRules = []rule.Rule{
	rules.NewEnumFieldNamesUpperSnakeCaseRule(),
	rules.NewEnumNamesUpperCamelCaseRule(),
	rules.NewFieldNamesLowerSnakeCaseRule(),
	rules.NewMessageNamesUpperCamelCaseRule(),
	rules.NewRPCNamesUpperCamelCaseRule(),
	rules.NewServiceNamesUpperCamelCaseRule(),
}

// DefaultRuleIDs are the default rule ids.
func DefaultRuleIDs() []string {
	var ids []string
	for _, rule := range defaultRules {
		ids = append(ids, rule.ID())
	}
	return ids
}

// NewAllRules creates new rules.
func NewAllRules() []rule.Rule {
	return defaultRules
}
