package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/config"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

func defaultRules(
	option config.RuleOption,
) []rule.Rule {
	maxLineLength := option.MaxLineLength

	return []rule.Rule{
		rules.NewEnumFieldNamesUpperSnakeCaseRule(),
		rules.NewEnumNamesUpperCamelCaseRule(),
		rules.NewFieldNamesLowerSnakeCaseRule(),
		rules.NewMessageNamesUpperCamelCaseRule(),
		rules.NewRPCNamesUpperCamelCaseRule(),
		rules.NewServiceNamesUpperCamelCaseRule(),
		rules.NewMaxLineLengthRule(
			maxLineLength.MaxChars,
			maxLineLength.TabChars,
		),
	}
}

// DefaultRuleIDs are the default rule ids.
func DefaultRuleIDs() []string {
	emptyOption := config.RuleOption{}

	var ids []string
	for _, rule := range defaultRules(emptyOption) {
		ids = append(ids, rule.ID())
	}
	return ids
}

// NewAllRules creates new rules.
func NewAllRules(
	option config.RuleOption,
) []rule.Rule {
	return defaultRules(option)
}
