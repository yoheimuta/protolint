package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/fixer"
	"github.com/yoheimuta/protolint/linter/rule"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// FieldNamesLowerSnakeCaseRule verifies that all field names are underscore_separated_names.
// See https://developers.google.com/protocol-buffers/docs/style#message-and-field-names.
type FieldNamesLowerSnakeCaseRule struct {
	RuleWithSeverity
	fixMode         bool
	autoDisableType autodisable.PlacementType
}

// NewFieldNamesLowerSnakeCaseRule creates a new FieldNamesLowerSnakeCaseRule.
func NewFieldNamesLowerSnakeCaseRule(
	severity rule.Severity,
	fixMode bool,
	autoDisableType autodisable.PlacementType,
) FieldNamesLowerSnakeCaseRule {
	if autoDisableType != autodisable.Noop {
		fixMode = false
	}
	return FieldNamesLowerSnakeCaseRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
		autoDisableType:  autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r FieldNamesLowerSnakeCaseRule) ID() string {
	return "FIELD_NAMES_LOWER_SNAKE_CASE"
}

// Purpose returns the purpose of this rule.
func (r FieldNamesLowerSnakeCaseRule) Purpose() string {
	return "Verifies that all field names are underscore_separated_names."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r FieldNamesLowerSnakeCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r FieldNamesLowerSnakeCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &fieldNamesLowerSnakeCaseVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type fieldNamesLowerSnakeCaseVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitField checks the field.
func (v *fieldNamesLowerSnakeCaseVisitor) VisitField(field *parser.Field) bool {
	name := field.FieldName
	if !strs.IsLowerSnakeCase(name) {
		expected := strs.ToLowerSnakeCase(name)
		v.AddFailuref(field.Meta.Pos, "Field name %q must be underscore_separated_names like %q", name, expected)

		err := v.Fixer.SearchAndReplace(field.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
			switch lex.Token {
			case scanner.TREPEATED, scanner.TREQUIRED, scanner.TOPTIONAL:
			default:
				lex.UnNext()
			}
			parseType(lex)
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

// VisitMapField checks the map field.
func (v *fieldNamesLowerSnakeCaseVisitor) VisitMapField(field *parser.MapField) bool {
	name := field.MapName
	if !strs.IsLowerSnakeCase(name) {
		expected := strs.ToLowerSnakeCase(name)
		v.AddFailuref(field.Meta.Pos, "Field name %q must be underscore_separated_names like %q", name, expected)

		err := v.Fixer.SearchAndReplace(field.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
			lex.Next()
			lex.Next()
			lex.Next()
			parseType(lex)
			lex.Next()
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

// VisitOneofField checks the oneof field.
func (v *fieldNamesLowerSnakeCaseVisitor) VisitOneofField(field *parser.OneofField) bool {
	name := field.FieldName
	if !strs.IsLowerSnakeCase(name) {
		expected := strs.ToLowerSnakeCase(name)
		v.AddFailuref(field.Meta.Pos, "Field name %q must be underscore_separated_names like %q", name, expected)

		err := v.Fixer.SearchAndReplace(field.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			parseType(lex)
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

// Below codes are copied from go-protoparser.
var typeConstants = map[string]struct{}{
	"double":   {},
	"float":    {},
	"int32":    {},
	"int64":    {},
	"uint32":   {},
	"uint64":   {},
	"sint32":   {},
	"sint64":   {},
	"fixed32":  {},
	"fixed64":  {},
	"sfixed32": {},
	"sfixed64": {},
	"bool":     {},
	"string":   {},
	"bytes":    {},
}

func parseType(lex *lexer.Lexer) {
	lex.Next()
	if _, ok := typeConstants[lex.Text]; ok {
		return
	}
	lex.UnNext()

	_, _, err := lex.ReadMessageType()
	if err != nil {
		panic(err)
	}
}
