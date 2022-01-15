package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/fixer"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// EnumNamesUpperCamelCaseRule verifies that all enum names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumNamesUpperCamelCaseRule struct {
	fixMode bool
}

// NewEnumNamesUpperCamelCaseRule creates a new EnumNamesUpperCamelCaseRule.
func NewEnumNamesUpperCamelCaseRule(
	fixMode bool,
) EnumNamesUpperCamelCaseRule {
	return EnumNamesUpperCamelCaseRule{
		fixMode: fixMode,
	}
}

// ID returns the ID of this rule.
func (r EnumNamesUpperCamelCaseRule) ID() string {
	return "ENUM_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r EnumNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all enum names are CamelCase (with an initial capital)."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumNamesUpperCamelCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r EnumNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &enumNamesUpperCamelCaseVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type enumNamesUpperCamelCaseVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitEnum checks the enum.
func (v *enumNamesUpperCamelCaseVisitor) VisitEnum(enum *parser.Enum) bool {
	name := enum.EnumName
	if !strs.IsUpperCamelCase(name) {
		expected := strs.ToUpperCamelCase(name)
		v.AddFailuref(enum.Meta.Pos, "Enum name %q must be UpperCamelCase like %q", name, expected)

		err := v.Fixer.SearchAndReplace(enum.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
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
