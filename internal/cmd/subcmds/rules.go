package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/config"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

func defaultRules(
	option config.RulesOption,
	fixMode bool,
) []rule.Rule {
	maxLineLength := option.MaxLineLength
	indent := option.Indent

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
		rules.NewIndentRule(
			indent.Style,
			indent.Newline,
			fixMode,
		),
	}
}

// DefaultRuleIDs are the default rule ids.
func DefaultRuleIDs() []string {
	emptyOption := config.RulesOption{}

	var ids []string
	for _, rule := range defaultRules(emptyOption, false) {
		ids = append(ids, rule.ID())
	}
	return ids
}

// NewAllRules creates new rules.
func NewAllRules(
	option config.RulesOption,
	fixMode bool,
) []rule.Rule {
	return defaultRules(option, fixMode)
}
