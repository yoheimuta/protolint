package rules

import (
	"strings"

	"github.com/yoheimuta/protolint/linter/strs"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// EnumFieldNamesPrefixRule verifies that enum field names are prefixed with its ENUM_NAME_UPPER_SNAKE_CASE.
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumFieldNamesPrefixRule struct {
}

// NewEnumFieldNamesPrefixRule creates a new EnumFieldNamesPrefixRule.
func NewEnumFieldNamesPrefixRule() EnumFieldNamesPrefixRule {
	return EnumFieldNamesPrefixRule{}
}

// ID returns the ID of this rule.
func (r EnumFieldNamesPrefixRule) ID() string {
	return "ENUM_FIELD_NAMES_PREFIX"
}

// Purpose returns the purpose of this rule.
func (r EnumFieldNamesPrefixRule) Purpose() string {
	return `Verifies that enum field names are prefixed with its ENUM_NAME_UPPER_SNAKE_CASE.`
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumFieldNamesPrefixRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r EnumFieldNamesPrefixRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &enumFieldNamesPrefixVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type enumFieldNamesPrefixVisitor struct {
	*visitor.BaseAddVisitor
	enumName string
}

// VisitEnum checks the enum.
func (v *enumFieldNamesPrefixVisitor) VisitEnum(enum *parser.Enum) bool {
	v.enumName = enum.EnumName
	return true
}

// VisitEnumField checks the enum field.
func (v *enumFieldNamesPrefixVisitor) VisitEnumField(field *parser.EnumField) bool {
	expectedPrefix, err := strs.ToUpperSnakeCaseFromCamelCase(v.enumName)
	if err != nil {
		expectedPrefix = strings.ToUpper(v.enumName)
	}
	if !strings.HasPrefix(field.Ident, expectedPrefix) {
		v.AddFailuref(field.Meta.Pos, "EnumField name %q should have the prefix %q", field.Ident, expectedPrefix)
	}
	return false
}
