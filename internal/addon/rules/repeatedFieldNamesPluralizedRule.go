package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
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
}

// NewRepeatedFieldNamesPluralizedRule creates a new RepeatedFieldNamesPluralizedRule.
func NewRepeatedFieldNamesPluralizedRule(
	pluralRules map[string]string,
	singularRules map[string]string,
	uncountableRules []string,
	irregularRules map[string]string,
) RepeatedFieldNamesPluralizedRule {
	return RepeatedFieldNamesPluralizedRule{
		pluralRules:      pluralRules,
		singularRules:    singularRules,
		uncountableRules: uncountableRules,
		irregularRules:   irregularRules,
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

	v := &repeatedFieldNamesPluralizedCaseVisitor{
		BaseAddVisitor:  visitor.NewBaseAddVisitor(r.ID()),
		pluralizeClient: c,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type repeatedFieldNamesPluralizedCaseVisitor struct {
	*visitor.BaseAddVisitor
	pluralizeClient *strs.PluralizeClient
}

// VisitField checks the field.
func (v *repeatedFieldNamesPluralizedCaseVisitor) VisitField(field *parser.Field) bool {
	got := field.FieldName
	want := v.pluralizeClient.ToPlural(got)
	if field.IsRepeated && got != want {
		v.AddFailuref(field.Meta.Pos, "Repeated field name %q must be pluralized name %q", got, want)
	}
	return false
}

// VisitGroupField checks the group field.
func (v *repeatedFieldNamesPluralizedCaseVisitor) VisitGroupField(field *parser.GroupField) bool {
	got := field.GroupName
	want := v.pluralizeClient.ToPlural(got)
	if field.IsRepeated && got != want {
		v.AddFailuref(field.Meta.Pos, "Repeated group name %q must be pluralized name %q", got, want)
	}
	return true
}
