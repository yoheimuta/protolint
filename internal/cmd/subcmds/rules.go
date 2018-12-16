package subcmds

import (
	"github.com/yoheimuta/protolinter/internal/addon/rules"
	"github.com/yoheimuta/protolinter/internal/linter/rule"
)

// NewAllRules creates new rules.
func NewAllRules() []rule.Rule {
	return []rule.Rule{
		rules.NewEnumNamesUpperCamelCaseRule(),
		rules.NewRPCNamesUpperCamelCaseRule(),
	}
}
