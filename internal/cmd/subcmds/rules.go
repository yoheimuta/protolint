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
	fileNamesLowerSnakeCase := option.FileNamesLowerSnakeCase
	enumFieldNamesZeroValueEndWith := option.EnumFieldNamesZeroValueEndWith
	maxLineLength := option.MaxLineLength
	indent := option.Indent

	return []rule.Rule{
		rules.NewEnumFieldNamesUpperSnakeCaseRule(),
		rules.NewEnumFieldNamesZeroValueEndWithRule(
			enumFieldNamesZeroValueEndWith.Suffix,
		),
		rules.NewEnumNamesUpperCamelCaseRule(),
		rules.NewFileNamesLowerSnakeCaseRule(
			fileNamesLowerSnakeCase.Excludes,
		),
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

func ruleIDs(rules []rule.Rule) []string {
	var ids []string
	for _, rule := range rules {
		ids = append(ids, rule.ID())
	}
	return ids
}

// DefaultRuleIDs are the default rule ids.
func DefaultRuleIDs() []string {
	emptyOption := config.RulesOption{}
	return ruleIDs(defaultRules(emptyOption, false))
}

// AllRuleIDs are the all rule ids.
func AllRuleIDs() []string {
	emptyOption := config.RulesOption{}
	return ruleIDs(NewAllRules(emptyOption, false))
}

// NewAllRules creates new rules.
func NewAllRules(
	option config.RulesOption,
	fixMode bool,
) []rule.Rule {
	serviceNamesEndWith := option.ServiceNamesEndWith
	fieldNamesExcludePrepositions := option.FieldNamesExcludePrepositions
	messageNamesExcludePrepositions := option.MessageNamesExcludePrepositions

	return append(
		defaultRules(option, fixMode),
		rules.NewServiceNamesEndWithRule(
			serviceNamesEndWith.Text,
		),
		rules.NewFieldNamesExcludePrepositionsRule(
			fieldNamesExcludePrepositions.Prepositions,
		),
		rules.NewMessageNamesExcludePrepositionsRule(
			messageNamesExcludePrepositions.Prepositions,
		),
	)
}
