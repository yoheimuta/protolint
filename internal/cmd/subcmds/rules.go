package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

// NewAllRules creates new rules.
func NewAllRules() []rule.Rule {
	return []rule.Rule{
		rules.NewEnumFieldNamesUpperSnakeCaseRule(),
		rules.NewEnumNamesUpperCamelCaseRule(),
		rules.NewFieldNamesLowerSnakeCaseRule(),
		rules.NewMessageNamesUpperCamelCaseRule(),
		rules.NewRPCNamesUpperCamelCaseRule(),
		rules.NewServiceNamesUpperCamelCaseRule(),
	}
}
