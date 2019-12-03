package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/stringsutil"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// Default values are a conservative list picked out from all preposition candidates.
// See https://www.talkenglish.com/vocabulary/top-50-prepositions.aspx
var defaultPrepositions = []string{
	"of",
	"with",
	"at",
	"from",
	"into",
	"during",
	"including",
	"until",
	"against",
	"among",
	"throughout",
	"despite",
	"towards",
	"upon",
	"concerning",

	"to",
	"in",
	"for",
	"on",
	"by",
	"about",
}

// FieldNamesExcludePrepositionsRule verifies that all field names don't include prepositions (e.g. "for", "during", "at").
// It is assumed that the field names are underscore_separated_names.
// See https://cloud.google.com/apis/design/naming_convention#field_names.
type FieldNamesExcludePrepositionsRule struct {
	prepositions []string
	excludes     []string
}

// NewFieldNamesExcludePrepositionsRule creates a new FieldNamesExcludePrepositionsRule.
func NewFieldNamesExcludePrepositionsRule(
	prepositions []string,
	excludes []string,
) FieldNamesExcludePrepositionsRule {
	if len(prepositions) == 0 {
		prepositions = defaultPrepositions
	}
	return FieldNamesExcludePrepositionsRule{
		prepositions: prepositions,
		excludes:     excludes,
	}
}

// ID returns the ID of this rule.
func (r FieldNamesExcludePrepositionsRule) ID() string {
	return "FIELD_NAMES_EXCLUDE_PREPOSITIONS"
}

// Purpose returns the purpose of this rule.
func (r FieldNamesExcludePrepositionsRule) Purpose() string {
	return `Verifies that all field names don't include prepositions (e.g. "for", "during", "at").`
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r FieldNamesExcludePrepositionsRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r FieldNamesExcludePrepositionsRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &fieldNamesExcludePrepositionsVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		prepositions:   r.prepositions,
		excludes:       r.excludes,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type fieldNamesExcludePrepositionsVisitor struct {
	*visitor.BaseAddVisitor
	prepositions []string
	excludes     []string
}

// VisitField checks the field.
func (v *fieldNamesExcludePrepositionsVisitor) VisitField(field *parser.Field) bool {
	name := field.FieldName
	for _, e := range v.excludes {
		name = strings.Replace(name, e, "", -1)
	}

	parts := strs.SplitSnakeCaseWord(name)
	for _, p := range parts {
		if stringsutil.ContainsStringInSlice(p, v.prepositions) {
			v.AddFailuref(field.Meta.Pos, "Field name %q should not include a preposition %q", field.FieldName, p)
		}
	}
	return false
}

// VisitMapField checks the map field.
func (v *fieldNamesExcludePrepositionsVisitor) VisitMapField(field *parser.MapField) bool {
	name := field.MapName
	for _, e := range v.excludes {
		name = strings.Replace(name, e, "", -1)
	}

	parts := strs.SplitSnakeCaseWord(name)
	for _, p := range parts {
		if stringsutil.ContainsStringInSlice(p, v.prepositions) {
			v.AddFailuref(field.Meta.Pos, "Field name %q should not include a preposition %q", field.MapName, p)
		}
	}
	return false
}

// VisitOneofField checks the oneof field.
func (v *fieldNamesExcludePrepositionsVisitor) VisitOneofField(field *parser.OneofField) bool {
	name := field.FieldName
	for _, e := range v.excludes {
		name = strings.Replace(name, e, "", -1)
	}

	parts := strs.SplitSnakeCaseWord(name)
	for _, p := range parts {
		if stringsutil.ContainsStringInSlice(p, v.prepositions) {
			v.AddFailuref(field.Meta.Pos, "Field name %q should not include a preposition %q", field.FieldName, p)
		}
	}
	return false
}
