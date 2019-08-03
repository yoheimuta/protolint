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
	importsSorted := option.ImportsSorted
	maxLineLength := option.MaxLineLength
	indent := option.Indent
	serviceNamesEndWith := option.ServiceNamesEndWith
	fieldNamesExcludePrepositions := option.FieldNamesExcludePrepositions
	messageNamesExcludePrepositions := option.MessageNamesExcludePrepositions
	messagesHaveComment := option.MessagesHaveComment
	servicesHaveComment := option.ServicesHaveComment
	rpcsHaveComment := option.RPCsHaveComment
	fieldsHaveComment := option.FieldsHaveComment
	enumsHaveComment := option.EnumsHaveComment
	enumFieldsHaveComment := option.EnumFieldsHaveComment

	return rule.Rules{
		rules.NewOrderRule(),
		rules.NewEnumFieldNamesUpperSnakeCaseRule(),
		rules.NewEnumFieldNamesZeroValueEndWithRule(
			enumFieldNamesZeroValueEndWith.Suffix,
		),
		rules.NewEnumFieldsHaveCommentRule(
			enumFieldsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewEnumNamesUpperCamelCaseRule(),
		rules.NewEnumsHaveCommentRule(
			enumsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewFileNamesLowerSnakeCaseRule(
			fileNamesLowerSnakeCase.Excludes,
		),
		rules.NewFieldNamesLowerSnakeCaseRule(),
		rules.NewFieldsHaveCommentRule(
			fieldsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewImportsSortedRule(
			importsSorted.Newline,
			fixMode,
		),
		rules.NewMessageNamesUpperCamelCaseRule(),
		rules.NewPackageNameLowerCaseRule(),
		rules.NewRPCNamesUpperCamelCaseRule(),
		rules.NewRPCsHaveCommentRule(
			rpcsHaveComment.ShouldFollowGolangStyle,
		),
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
		rules.NewServicesHaveCommentRule(
			servicesHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewFieldNamesExcludePrepositionsRule(
			fieldNamesExcludePrepositions.Prepositions,
			fieldNamesExcludePrepositions.Excludes,
		),
		rules.NewMessageNamesExcludePrepositionsRule(
			messageNamesExcludePrepositions.Prepositions,
			messageNamesExcludePrepositions.Excludes,
		),
		rules.NewMessagesHaveCommentRule(
			messagesHaveComment.ShouldFollowGolangStyle,
		),
	}
}
