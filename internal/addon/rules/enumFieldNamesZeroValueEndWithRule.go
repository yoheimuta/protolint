package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/protolint/linter/fixer"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor"
)

const (
	defaultSuffix = "UNSPECIFIED"
)

// EnumFieldNamesZeroValueEndWithRule verifies that the zero value enum should have the suffix (e.g. "UNSPECIFIED", "INVALID").
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumFieldNamesZeroValueEndWithRule struct {
	suffix  string
	fixMode bool
}

// NewEnumFieldNamesZeroValueEndWithRule creates a new EnumFieldNamesZeroValueEndWithRule.
func NewEnumFieldNamesZeroValueEndWithRule(
	suffix string,
	fixMode bool,
) EnumFieldNamesZeroValueEndWithRule {
	if len(suffix) == 0 {
		suffix = defaultSuffix
	}
	return EnumFieldNamesZeroValueEndWithRule{
		suffix:  suffix,
		fixMode: fixMode,
	}
}

// ID returns the ID of this rule.
func (r EnumFieldNamesZeroValueEndWithRule) ID() string {
	return "ENUM_FIELD_NAMES_ZERO_VALUE_END_WITH"
}

// Purpose returns the purpose of this rule.
func (r EnumFieldNamesZeroValueEndWithRule) Purpose() string {
	return `Verifies that the zero value enum should have the suffix (e.g. "UNSPECIFIED", "INVALID").`
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumFieldNamesZeroValueEndWithRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r EnumFieldNamesZeroValueEndWithRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &enumFieldNamesZeroValueEndWithVisitor{
		BaseFixableVisitor: base,
		suffix:             r.suffix,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type enumFieldNamesZeroValueEndWithVisitor struct {
	*visitor.BaseFixableVisitor
	suffix string
}

// VisitEnumField checks the enum field.
func (v *enumFieldNamesZeroValueEndWithVisitor) VisitEnumField(field *parser.EnumField) bool {
	if field.Number == "0" && !strings.HasSuffix(field.Ident, v.suffix) {
		v.AddFailuref(field.Meta.Pos, "EnumField name %q with zero value should have the suffix %q", field.Ident, v.suffix)

		expected := field.Ident + "_" + v.suffix
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
