package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/fixer"
	"github.com/yoheimuta/protolint/linter/rule"

	"github.com/yoheimuta/protolint/linter/strs"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// EnumFieldNamesPrefixRule verifies that enum field names are prefixed with its ENUM_NAME_UPPER_SNAKE_CASE.
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumFieldNamesPrefixRule struct {
	RuleWithSeverity
	fixMode         bool
	autoDisableType autodisable.PlacementType
}

// NewEnumFieldNamesPrefixRule creates a new EnumFieldNamesPrefixRule.
func NewEnumFieldNamesPrefixRule(
	severity rule.Severity,
	fixMode bool,
	autoDisableType autodisable.PlacementType,
) EnumFieldNamesPrefixRule {
	if autoDisableType != autodisable.Noop {
		fixMode = false
	}
	return EnumFieldNamesPrefixRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
		autoDisableType:  autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r EnumFieldNamesPrefixRule) ID() string {
	return "ENUM_FIELD_NAMES_PREFIX"
}

// Purpose returns the purpose of this rule.
func (r EnumFieldNamesPrefixRule) Purpose() string {
	return `Verifies that enum field names are prefixed with its ENUM_NAME_UPPER_SNAKE_CASE.`
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumFieldNamesPrefixRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r EnumFieldNamesPrefixRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &enumFieldNamesPrefixVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type enumFieldNamesPrefixVisitor struct {
	*visitor.BaseFixableVisitor
	enumName string
}

// VisitEnum checks the enum.
func (v *enumFieldNamesPrefixVisitor) VisitEnum(enum *parser.Enum) bool {
	v.enumName = enum.EnumName
	return true
}

// VisitEnumField checks the enum field.
func (v *enumFieldNamesPrefixVisitor) VisitEnumField(field *parser.EnumField) bool {
	expectedPrefix := strs.ToUpperSnakeCase(v.enumName)
	if !strings.HasPrefix(field.Ident, expectedPrefix) {
		v.AddFailuref(field.Meta.Pos, "EnumField name %q should have the prefix %q", field.Ident, expectedPrefix)

		expected := expectedPrefix + "_" + field.Ident
		err := v.Fixer.SearchAndReplace(field.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.Next()
			return fixer.TextEdit{
				Pos:     lex.Pos.Offset,
				End:     lex.Pos.Offset + len(lex.Text) - 1,
				NewText: []byte(expected),
			}
		})
		if err != nil {
			panic(err)
		}
	}
	return false
}
