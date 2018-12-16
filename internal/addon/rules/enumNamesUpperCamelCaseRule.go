package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolinter/internal/linter/report"
	"github.com/yoheimuta/protolinter/internal/strs"
)

// EnumNamesUpperCamelCaseRule checks that enum type names are CamelCase (with an initial capital).
type EnumNamesUpperCamelCaseRule struct{}

// NewEnumNamesUpperCamelCaseRule creates a new EnumNamesUpperCamelCaseRule.
func NewEnumNamesUpperCamelCaseRule() EnumNamesUpperCamelCaseRule {
	return EnumNamesUpperCamelCaseRule{}
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
		v.addFailuref(enum.Meta.Pos, "Enum name %q must be UpperCamelCase.", enum.EnumName)
	}
	return false
}
