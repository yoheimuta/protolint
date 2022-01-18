package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/lexer/scanner"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/fixer"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// RepeatedFieldNamesPluralizedRule verifies that repeated field names are pluralized names.
// See https://developers.google.com/protocol-buffers/docs/style#repeated-fields.
type RepeatedFieldNamesPluralizedRule struct {
	pluralRules      map[string]string
	singularRules    map[string]string
	uncountableRules []string
	irregularRules   map[string]string
	fixMode          bool
}

// NewRepeatedFieldNamesPluralizedRule creates a new RepeatedFieldNamesPluralizedRule.
func NewRepeatedFieldNamesPluralizedRule(
	pluralRules map[string]string,
	singularRules map[string]string,
	uncountableRules []string,
	irregularRules map[string]string,
	fixMode bool,
) RepeatedFieldNamesPluralizedRule {
	return RepeatedFieldNamesPluralizedRule{
		pluralRules:      pluralRules,
		singularRules:    singularRules,
		uncountableRules: uncountableRules,
		irregularRules:   irregularRules,
		fixMode:          fixMode,
	}
}

// ID returns the ID of this rule.
func (r RepeatedFieldNamesPluralizedRule) ID() string {
	return "REPEATED_FIELD_NAMES_PLURALIZED"
}

// Purpose returns the purpose of this rule.
func (r RepeatedFieldNamesPluralizedRule) Purpose() string {
	return "Verifies that repeated field names are pluralized names."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RepeatedFieldNamesPluralizedRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r RepeatedFieldNamesPluralizedRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	c := strs.NewPluralizeClient()
	for k, v := range r.pluralRules {
		c.AddPluralRule(k, v)
	}
	for k, v := range r.singularRules {
		c.AddSingularRule(k, v)
	}
	for _, w := range r.uncountableRules {
		c.AddUncountableRule(w)
	}
	for k, v := range r.irregularRules {
		c.AddIrregularRule(k, v)
	}

	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &repeatedFieldNamesPluralizedVisitor{
		BaseFixableVisitor: base,
		pluralizeClient:    c,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type repeatedFieldNamesPluralizedVisitor struct {
	*visitor.BaseFixableVisitor
	pluralizeClient *strs.PluralizeClient
}

// VisitField checks the field.
func (v *repeatedFieldNamesPluralizedVisitor) VisitField(field *parser.Field) bool {
	got := field.FieldName
	want := v.pluralizeClient.ToPlural(got)
	if field.IsRepeated && got != want {
		v.AddFailuref(field.Meta.Pos, "Repeated field name %q must be pluralized name %q", got, want)

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
				NewText: []byte(want),
			}
		})
		if err != nil {
			panic(err)
		}
	}
	return false
}

// VisitGroupField checks the group field.
func (v *repeatedFieldNamesPluralizedVisitor) VisitGroupField(field *parser.GroupField) bool {
	got := field.GroupName
	want := v.pluralizeClient.ToPlural(got)
	if field.IsRepeated && got != want {
		v.AddFailuref(field.Meta.Pos, "Repeated group name %q must be pluralized name %q", got, want)

		err := v.Fixer.SearchAndReplace(field.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
			switch lex.Token {
			case scanner.TREPEATED, scanner.TREQUIRED, scanner.TOPTIONAL:
			default:
				lex.UnNext()
			}
			lex.NextKeyword()
			lex.Next()
			return fixer.TextEdit{
				Pos:     lex.Pos.Offset,
				End:     lex.Pos.Offset + len(lex.Text) - 1,
				NewText: []byte(want),
			}
		})
		if err != nil {
			panic(err)
		}
	}
	return true
}
