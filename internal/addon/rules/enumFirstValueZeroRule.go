package rules

import (
	"strconv"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// EnumFirstValueZeroRule verifies that the first value of an enum has a numeric value of zero.
// See https://protobuf.dev/programming-guides/dos-donts/#unspecified-enum.
type EnumFirstValueZeroRule struct {
	RuleWithSeverity
	autoDisableType autodisable.PlacementType
}

// NewEnumFirstValueZeroRule creates a new EnumFirstValueZeroRule.
func NewEnumFirstValueZeroRule(
	severity rule.Severity,
	autoDisableType autodisable.PlacementType,
) EnumFirstValueZeroRule {
	return EnumFirstValueZeroRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		autoDisableType:  autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r EnumFirstValueZeroRule) ID() string {
	return "ENUM_FIRST_VALUE_ZERO"
}

// Purpose returns the purpose of this rule.
func (r EnumFirstValueZeroRule) Purpose() string {
	return "Verifies that the first value of an enum has a numeric value of zero."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumFirstValueZeroRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r EnumFirstValueZeroRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &enumFirstValueZeroVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type enumFirstValueZeroVisitor struct {
	*visitor.BaseAddVisitor
	enumName   string
	firstField *parser.EnumField
}

// VisitEnum records the first enum field of the enum so VisitEnumField only checks that field.
func (v *enumFirstValueZeroVisitor) VisitEnum(enum *parser.Enum) bool {
	v.enumName = enum.EnumName
	v.firstField = nil
	for _, body := range enum.EnumBody {
		if field, ok := body.(*parser.EnumField); ok {
			v.firstField = field
			break
		}
	}
	return true
}

// VisitEnumField verifies that the first enum field has numeric value 0.
func (v *enumFirstValueZeroVisitor) VisitEnumField(field *parser.EnumField) bool {
	if field != v.firstField {
		return false
	}
	num, err := strconv.ParseInt(field.Number, 0, 64)
	if err != nil || num != 0 {
		v.AddFailuref(
			field.Meta.Pos,
			"First value %q of enum %q should have numeric value 0, got %s",
			field.Ident,
			v.enumName,
			field.Number,
		)
	}
	return false
}
