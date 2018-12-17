package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolinter/internal/linter/report"
	"github.com/yoheimuta/protolinter/internal/strs"
)

// EnumNamesUpperCamelCaseRule verifies that all enum names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumNamesUpperCamelCaseRule struct{}

// NewEnumNamesUpperCamelCaseRule creates a new EnumNamesUpperCamelCaseRule.
func NewEnumNamesUpperCamelCaseRule() EnumNamesUpperCamelCaseRule {
	return EnumNamesUpperCamelCaseRule{}
}

// ID returns the ID of this rule.
func (r EnumNamesUpperCamelCaseRule) ID() string {
	return "ENUM_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r EnumNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all enum names are CamelCase (with an initial capital)."
}

// Apply applies the rule to the proto.
func (r EnumNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	visitor := &enumNamesUpperCamelCaseVisitor{
		baseAddVisitor: newBaseAddVisitor(),
	}
	return runVisitor(visitor, proto)
}

type enumNamesUpperCamelCaseVisitor struct {
	*baseAddVisitor
}

// VisitEnum checks the enum.
func (v *enumNamesUpperCamelCaseVisitor) VisitEnum(enum *parser.Enum) bool {
	if !strs.IsUpperCamelCase(enum.EnumName) {
		v.addFailuref(enum.Meta.Pos, "Enum name %q must be UpperCamelCase", enum.EnumName)
	}
	return false
}
