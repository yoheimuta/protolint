package subcmds

import (
	"github.com/yoheimuta/protolint/internal/addon/plugin"
	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/config"
	internalrule "github.com/yoheimuta/protolint/internal/linter/rule"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/rule"
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
		rules.NewFileHasCommentRule(
			rule.Severity_Error,
		),
		rules.NewSyntaxConsistentRule(
			rule.Severity_Error,
			syntaxConsistent.Version,
		),
		rules.NewFileNamesLowerSnakeCaseRule(
			rule.Severity_Error,
			fileNamesLowerSnakeCase.Excludes,
			fixMode,
		),
		rules.NewQuoteConsistentRule(
			rule.Severity_Error,
			option.QuoteConsistentOption.Quote,
			fixMode,
		),
		rules.NewOrderRule(
			rule.Severity_Error,
			fixMode,
		),
		rules.NewIndentRule(
			rule.Severity_Error,
			indent.Style,
			indent.NotInsertNewline,
			fixMode,
		),
		rules.NewMaxLineLengthRule(
			rule.Severity_Error,
			maxLineLength.MaxChars,
			maxLineLength.TabChars,
		),
		rules.NewPackageNameLowerCaseRule(
			rule.Severity_Error,
			fixMode,
		),
		rules.NewImportsSortedRule(
			rule.Severity_Error,
			fixMode,
		),
		rules.NewEnumFieldNamesPrefixRule(
			rule.Severity_Error,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumFieldNamesUpperSnakeCaseRule(
			rule.Severity_Error,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumFieldNamesZeroValueEndWithRule(
			rule.Severity_Error,
			enumFieldNamesZeroValueEndWith.Suffix,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumFieldsHaveCommentRule(
			rule.Severity_Error,
			enumFieldsHaveComment.ShouldFollowGolangStyle,
		),

		rules.NewEnumNamesUpperCamelCaseRule(
			rule.Severity_Error,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumsHaveCommentRule(
			rule.Severity_Error,
			enumsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewFieldNamesLowerSnakeCaseRule(
			rule.Severity_Error,
			fixMode,
			autoDisableType,
		),
		rules.NewFieldNamesExcludePrepositionsRule(
			rule.Severity_Error,
			fieldNamesExcludePrepositions.Prepositions,
			fieldNamesExcludePrepositions.Excludes,
		),
		rules.NewFieldsHaveCommentRule(
			rule.Severity_Error,
			fieldsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewProto3FieldsAvoidRequiredRule(
			rule.Severity_Error,
			fixMode,
		),
		rules.NewProto3GroupsAvoidRule(
			rule.Severity_Error,
			autoDisableType,
		),
		rules.NewRepeatedFieldNamesPluralizedRule(
			rule.Severity_Error,
			repeatedFieldNamesPluralized.PluralRules,
			repeatedFieldNamesPluralized.SingularRules,
			repeatedFieldNamesPluralized.UncountableRules,
			repeatedFieldNamesPluralized.IrregularRules,
			fixMode,
			autoDisableType,
		),
		rules.NewMessageNamesUpperCamelCaseRule(
			rule.Severity_Error,
			fixMode,
			autoDisableType,
		),
		rules.NewMessageNamesExcludePrepositionsRule(
			rule.Severity_Error,
			messageNamesExcludePrepositions.Prepositions,
			messageNamesExcludePrepositions.Excludes,
		),
		rules.NewMessagesHaveCommentRule(
			rule.Severity_Error,
			messagesHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewRPCNamesUpperCamelCaseRule(
			rule.Severity_Error,
			fixMode,
			autoDisableType,
		),
		rules.NewRPCNamesCaseRule(
			rule.Severity_Error,
			option.RPCNamesCaseOption.Convention,
		),
		rules.NewRPCsHaveCommentRule(
			rule.Severity_Error,
			rpcsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewServiceNamesUpperCamelCaseRule(
			rule.Severity_Error,
			fixMode,
			autoDisableType,
		),
		rules.NewServiceNamesEndWithRule(
			rule.Severity_Error,
			serviceNamesEndWith.Text,
		),
		rules.NewServicesHaveCommentRule(
			rule.Severity_Error,
			servicesHaveComment.ShouldFollowGolangStyle,
		),
	}
}
