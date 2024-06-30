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
		rules.NewFileHasCommentRule(
			option.FileHasComment.Severity,
		),
		rules.NewSyntaxConsistentRule(
			syntaxConsistent.Severity,
			syntaxConsistent.Version,
		),
		rules.NewFileNamesLowerSnakeCaseRule(
			fileNamesLowerSnakeCase.Severity,
			fileNamesLowerSnakeCase.Excludes,
			fixMode,
		),
		rules.NewQuoteConsistentRule(
			option.QuoteConsistentOption.Severity,
			option.QuoteConsistentOption.Quote,
			fixMode,
		),
		rules.NewOrderRule(
			option.Order.Severity,
			fixMode,
		),
		rules.NewIndentRule(
			indent.Severity,
			indent.Style,
			indent.NotInsertNewline,
			fixMode,
		),
		rules.NewMaxLineLengthRule(
			maxLineLength.Severity,
			maxLineLength.MaxChars,
			maxLineLength.TabChars,
		),
		rules.NewPackageNameLowerCaseRule(
			option.PackageNameLowerCase.Severity,
			fixMode,
		),
		rules.NewImportsSortedRule(
			option.ImportsSorted.Severity,
			fixMode,
		),
		rules.NewEnumFieldNamesPrefixRule(
			option.EnumFieldNamesPrefix.Severity,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumFieldNamesUpperSnakeCaseRule(
			option.EnumFieldNamesUpperSnakeCase.Severity,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumFieldNamesZeroValueEndWithRule(
			enumFieldNamesZeroValueEndWith.Severity,
			enumFieldNamesZeroValueEndWith.Suffix,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumFieldsHaveCommentRule(
			enumFieldsHaveComment.Severity,
			enumFieldsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewEnumNamesUpperCamelCaseRule(
			option.EnumFieldNamesUpperSnakeCase.Severity,
			fixMode,
			autoDisableType,
		),
		rules.NewEnumsHaveCommentRule(
			enumsHaveComment.Severity,
			enumsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewFieldNamesLowerSnakeCaseRule(
			option.FieldNamesLowerSnakeCase.Severity,
			fixMode,
			autoDisableType,
		),
		rules.NewFieldNamesExcludePrepositionsRule(
			fieldNamesExcludePrepositions.Severity,
			fieldNamesExcludePrepositions.Prepositions,
			fieldNamesExcludePrepositions.Excludes,
		),
		rules.NewFieldsHaveCommentRule(
			fieldsHaveComment.Severity,
			fieldsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewProto3FieldsAvoidRequiredRule(
			option.Proto3FieldsAvoidRequired.Severity,
			fixMode,
		),
		rules.NewProto3GroupsAvoidRule(
			option.Proto3GroupsAvoid.Severity,
			autoDisableType,
		),
		rules.NewRepeatedFieldNamesPluralizedRule(
			repeatedFieldNamesPluralized.Severity,
			repeatedFieldNamesPluralized.PluralRules,
			repeatedFieldNamesPluralized.SingularRules,
			repeatedFieldNamesPluralized.UncountableRules,
			repeatedFieldNamesPluralized.IrregularRules,
			fixMode,
			autoDisableType,
		),
		rules.NewMessageNamesUpperCamelCaseRule(
			option.MessageNamesUpperCamelCase.Severity,
			fixMode,
			autoDisableType,
		),
		rules.NewMessageNamesExcludePrepositionsRule(
			messageNamesExcludePrepositions.Severity,
			messageNamesExcludePrepositions.Prepositions,
			messageNamesExcludePrepositions.Excludes,
		),
		rules.NewMessagesHaveCommentRule(
			messagesHaveComment.Severity,
			messagesHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewRPCNamesUpperCamelCaseRule(
			option.RPCNamesUpperCamelCase.Severity,
			fixMode,
			autoDisableType,
		),
		rules.NewRPCNamesCaseRule(
			option.RPCNamesCaseOption.Severity,
			option.RPCNamesCaseOption.Convention,
		),
		rules.NewRPCsHaveCommentRule(
			rpcsHaveComment.Severity,
			rpcsHaveComment.ShouldFollowGolangStyle,
		),
		rules.NewServiceNamesUpperCamelCaseRule(
			option.ServiceNamesUpperCamelCase.Severity,
			fixMode,
			autoDisableType,
		),
		rules.NewServiceNamesEndWithRule(
			option.ServiceNamesEndWith.Severity,
			serviceNamesEndWith.Text,
		),
		rules.NewServicesHaveCommentRule(
			option.ServicesHaveComment.Severity,
			servicesHaveComment.ShouldFollowGolangStyle,
		),
	}
}
