package rules

import (
	"strings"
	"unicode"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// UnderscoreRule verifies the underscore rules of the official style guide: an identifier must not begin or
// end with an underscore, and an underscore must be followed by a letter, so neither a digit nor a second
// underscore may follow it.
// See https://protobuf.dev/programming-guides/style/#message_and_field_names.
type UnderscoreRule struct {
	RuleWithSeverity
	autoDisableType autodisable.PlacementType
}

// NewUnderscoreRule creates a new UnderscoreRule.
func NewUnderscoreRule(
	severity rule.Severity,
	autoDisableType autodisable.PlacementType,
) UnderscoreRule {
	return UnderscoreRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		autoDisableType:  autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r UnderscoreRule) ID() string {
	return "UNDERSCORE"
}

// Purpose returns the purpose of this rule.
func (r UnderscoreRule) Purpose() string {
	return `Verifies that identifiers do not begin or end with an underscore and that every underscore is followed by a letter.`
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r UnderscoreRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r UnderscoreRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), false, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &underscoreVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type underscoreVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitField checks the field name.
func (v *underscoreVisitor) VisitField(field *parser.Field) bool {
	v.check("Field name", field.FieldName, field.Meta.Pos)
	return false
}

// VisitMapField checks the map field name.
func (v *underscoreVisitor) VisitMapField(field *parser.MapField) bool {
	v.check("Map field name", field.MapName, field.Meta.Pos)
	return false
}

// VisitOneofField checks the oneof field name.
func (v *underscoreVisitor) VisitOneofField(field *parser.OneofField) bool {
	v.check("Oneof field name", field.FieldName, field.Meta.Pos)
	return false
}

// VisitEnumField checks the enum field name.
func (v *underscoreVisitor) VisitEnumField(field *parser.EnumField) bool {
	v.check("Enum field name", field.Ident, field.Meta.Pos)
	return false
}

func (v *underscoreVisitor) check(kind, name string, pos meta.Position) {
	if problem := underscoreProblem(name); problem != "" {
		v.AddFailuref(pos, "%s %q %s", kind, name, problem)
	}
}

// underscoreProblem explains why name breaks the style guide's underscore rules, or returns "" when it does not.
func underscoreProblem(name string) string {
	if strings.HasPrefix(name, "_") {
		return "must not begin with an underscore"
	}
	if strings.HasSuffix(name, "_") {
		return "must not end with an underscore"
	}
	runes := []rune(name)
	for i, r := range runes {
		if r != '_' {
			continue
		}
		if i+1 < len(runes) && !unicode.IsLetter(runes[i+1]) {
			return "must have every underscore followed by a letter"
		}
	}
	return ""
}
