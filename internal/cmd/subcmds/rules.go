package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/plugin"
	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/config"
	internalrule "github.com/yoheimuta/protolint/internal/linter/rule"
	"github.com/yoheimuta/protolint/linter/autodisable"
)

// NewAllRules creates new all rules.
func NewAllRules(
	option config.RulesOption,
	fixMode bool,
	autoDisableType autodisable.PlacementType,
	verbose bool,
	plugins []shared.RuleSet,
) (internalrule.Rules, error) {
	rs := newAllInternalRules(option, fixMode, autoDisableType)

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
	autoDisableType autodisable.PlacementType,
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
			fixMode,
		),
		rules.NewQuoteConsistentRule(option.QuoteConsistentOption.Quote, fixMode),
		rules.NewOrderRule(fixMode),
		rules.NewIndentRule(
			indent.Style,
			indent.NotInsertNewline,
			fixMode,
		),
		rules.NewMaxLineLengthRule(
			maxLineLength.MaxChars,
			maxLineLength.TabChars,
		),

		rules.NewPackageNameLowerCaseRule(fixMode),

		rules.NewImportsSortedRule(
			fixMode,
		),

		rules.NewEnumFieldNamesPrefixRule(fixMode, autoDisableType),
		rules.NewEnumFieldNamesUpperSnakeCaseRule(fixMode, autoDisableType),
		rules.NewEnumFieldNamesZeroValueEndWithRule(
			enumFieldNamesZeroValueEndWith.Suffix,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumFieldsHaveCommentRule(
			enumFieldsHaveComment.ShouldFollowGolangStyle,
		),

		rules.NewEnumNamesUpperCamelCaseRule(fixMode, autoDisableType),
		rules.NewEnumsHaveCommentRule(
			enumsHaveComment.ShouldFollowGolangStyle,
		),

		rules.NewFieldNamesLowerSnakeCaseRule(fixMode, autoDisableType),
		rules.NewFieldNamesExcludePrepositionsRule(
			fieldNamesExcludePrepositions.Prepositions,
			fieldNamesExcludePrepositions.Excludes,
		),
		rules.NewFieldsHaveCommentRule(
			fieldsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewProto3FieldsAvoidRequiredRule(fixMode),
		rules.NewProto3GroupsAvoidRule(),
		rules.NewRepeatedFieldNamesPluralizedRule(
			repeatedFieldNamesPluralized.PluralRules,
			repeatedFieldNamesPluralized.SingularRules,
			repeatedFieldNamesPluralized.UncountableRules,
			repeatedFieldNamesPluralized.IrregularRules,
			fixMode,
		),

		rules.NewMessageNamesUpperCamelCaseRule(fixMode),
		rules.NewMessageNamesExcludePrepositionsRule(
			messageNamesExcludePrepositions.Prepositions,
			messageNamesExcludePrepositions.Excludes,
		),
		rules.NewMessagesHaveCommentRule(
			messagesHaveComment.ShouldFollowGolangStyle,
		),

		rules.NewRPCNamesUpperCamelCaseRule(fixMode),
		rules.NewRPCNamesCaseRule(option.RPCNamesCaseOption.Convention),
		rules.NewRPCsHaveCommentRule(
			rpcsHaveComment.ShouldFollowGolangStyle,
		),

		rules.NewServiceNamesUpperCamelCaseRule(fixMode),
		rules.NewServiceNamesEndWithRule(
			serviceNamesEndWith.Text,
		),
		rules.NewServicesHaveCommentRule(
			servicesHaveComment.ShouldFollowGolangStyle,
		),
	}
}
