package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/config"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

// NewAllRules creates new all rules.
func NewAllRules(
	option config.RulesOption,
	fixMode bool,
) rule.Rules {
	fileNamesLowerSnakeCase := option.FileNamesLowerSnakeCase
	enumFieldNamesZeroValueEndWith := option.EnumFieldNamesZeroValueEndWith
	maxLineLength := option.MaxLineLength
	indent := option.Indent
	serviceNamesEndWith := option.ServiceNamesEndWith
	fieldNamesExcludePrepositions := option.FieldNamesExcludePrepositions
	messageNamesExcludePrepositions := option.MessageNamesExcludePrepositions

	return rule.Rules{
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
		rules.NewServiceNamesEndWithRule(
			serviceNamesEndWith.Text,
		),
		rules.NewFieldNamesExcludePrepositionsRule(
			fieldNamesExcludePrepositions.Prepositions,
		),
		rules.NewMessageNamesExcludePrepositionsRule(
			messageNamesExcludePrepositions.Prepositions,
		),
	}
}
