package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/plugin"
	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/config"
	internalrule "github.com/yoheimuta/protolint/internal/linter/rule"
)

// NewAllRules creates new all rules.
func NewAllRules(
	option config.RulesOption,
	fixMode bool,
	verbose bool,
	plugins []shared.RuleSet,
) (internalrule.Rules, error) {
	rs := newAllInternalRules(option, fixMode)

	es, err := plugin.GetExternalRules(plugins, fixMode, verbose)
	if err != nil {
		return nil, err
	}
	rs = append(rs, es...)
	return rs, nil
}

func newAllInternalRules(
	option config.RulesOption,
	fixMode bool,
) internalrule.Rules {
	syntaxConsistent := option.SyntaxConsistent
	fileNamesLowerSnakeCase := option.FileNamesLowerSnakeCase
	indent := option.Indent
	maxLineLength := option.MaxLineLength
	enumFieldNamesZeroValueEndWith := option.EnumFieldNamesZeroValueEndWith
	serviceNamesEndWith := option.ServiceNamesEndWith
	fieldNamesExcludePrepositions := option.FieldNamesExcludePrepositions
	messageNamesExcludePrepositions := option.MessageNamesExcludePrepositions
	messagesHaveComment := option.MessagesHaveComment
	servicesHaveComment := option.ServicesHaveComment
	rpcsHaveComment := option.RPCsHaveComment
	fieldsHaveComment := option.FieldsHaveComment
	enumsHaveComment := option.EnumsHaveComment
	enumFieldsHaveComment := option.EnumFieldsHaveComment
	repeatedFieldNamesPluralized := option.RepeatedFieldNamesPluralized

	return internalrule.Rules{
		rules.NewFileHasCommentRule(),
		rules.NewSyntaxConsistentRule(
			syntaxConsistent.Version,
		),
		rules.NewFileNamesLowerSnakeCaseRule(
			fileNamesLowerSnakeCase.Excludes,
		),
		rules.NewQuoteConsistentRule(option.QuoteConsistentOption.Quote, fixMode),
		rules.NewOrderRule(),
		rules.NewIndentRule(
			indent.Style,
			indent.Newline,
			indent.NotInsertNewline,
			fixMode,
		),
		rules.NewMaxLineLengthRule(
			maxLineLength.MaxChars,
			maxLineLength.TabChars,
		),

		rules.NewPackageNameLowerCaseRule(),

		rules.NewImportsSortedRule(
			fixMode,
		),

		rules.NewEnumFieldNamesPrefixRule(),
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

		rules.NewFieldNamesLowerSnakeCaseRule(),
		rules.NewFieldNamesExcludePrepositionsRule(
			fieldNamesExcludePrepositions.Prepositions,
			fieldNamesExcludePrepositions.Excludes,
		),
		rules.NewFieldsHaveCommentRule(
			fieldsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewProto3FieldsAvoidRequiredRule(),
		rules.NewProto3GroupsAvoidRule(),
		rules.NewRepeatedFieldNamesPluralizedRule(
			repeatedFieldNamesPluralized.PluralRules,
			repeatedFieldNamesPluralized.SingularRules,
			repeatedFieldNamesPluralized.UncountableRules,
			repeatedFieldNamesPluralized.IrregularRules,
		),

		rules.NewMessageNamesUpperCamelCaseRule(),
		rules.NewMessageNamesExcludePrepositionsRule(
			messageNamesExcludePrepositions.Prepositions,
			messageNamesExcludePrepositions.Excludes,
		),
		rules.NewMessagesHaveCommentRule(
			messagesHaveComment.ShouldFollowGolangStyle,
		),

		rules.NewRPCNamesUpperCamelCaseRule(),
		rules.NewRPCNamesCaseRule(option.RPCNamesCaseOption.Convention),
		rules.NewRPCsHaveCommentRule(
			rpcsHaveComment.ShouldFollowGolangStyle,
		),

		rules.NewServiceNamesUpperCamelCaseRule(),
		rules.NewServiceNamesEndWithRule(
			serviceNamesEndWith.Text,
		),
		rules.NewServicesHaveCommentRule(
			servicesHaveComment.ShouldFollowGolangStyle,
		),
	}
}
