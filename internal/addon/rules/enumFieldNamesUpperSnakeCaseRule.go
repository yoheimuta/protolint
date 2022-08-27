package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/fixer"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// EnumFieldNamesUpperSnakeCaseRule verifies that all enum field names are CAPITALS_WITH_UNDERSCORES.
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumFieldNamesUpperSnakeCaseRule struct {
	fixMode         bool
	autoDisableType autodisable.PlacementType
}

// NewEnumFieldNamesUpperSnakeCaseRule creates a new EnumFieldNamesUpperSnakeCaseRule.
func NewEnumFieldNamesUpperSnakeCaseRule(
	fixMode bool,
	autoDisableType autodisable.PlacementType,
) EnumFieldNamesUpperSnakeCaseRule {
	if autoDisableType != autodisable.Noop {
		fixMode = false
	}
	return EnumFieldNamesUpperSnakeCaseRule{
		fixMode:         fixMode,
		autoDisableType: autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r EnumFieldNamesUpperSnakeCaseRule) ID() string {
	return "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE"
}

// Purpose returns the purpose of this rule.
func (r EnumFieldNamesUpperSnakeCaseRule) Purpose() string {
	return "Verifies that all enum field names are CAPITALS_WITH_UNDERSCORES."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumFieldNamesUpperSnakeCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r EnumFieldNamesUpperSnakeCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &enumFieldNamesUpperSnakeCaseVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type enumFieldNamesUpperSnakeCaseVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitEnumField checks the enum field.
func (v *enumFieldNamesUpperSnakeCaseVisitor) VisitEnumField(field *parser.EnumField) bool {
	name := field.Ident
	if !strs.IsUpperSnakeCase(name) {
		expected := strs.ToUpperSnakeCase(name)
		v.AddFailuref(field.Meta.Pos, "EnumField name %q must be CAPITALS_WITH_UNDERSCORES like %q", name, expected)

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
