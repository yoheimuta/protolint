package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// EnumFieldNamesUpperSnakeCaseRule verifies that all enum field names are CAPITALS_WITH_UNDERSCORES.
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumFieldNamesUpperSnakeCaseRule struct{}

// NewEnumFieldNamesUpperSnakeCaseRule creates a new EnumFieldNamesUpperSnakeCaseRule.
func NewEnumFieldNamesUpperSnakeCaseRule() EnumFieldNamesUpperSnakeCaseRule {
	return EnumFieldNamesUpperSnakeCaseRule{}
}

// ID returns the ID of this rule.
func (r EnumFieldNamesUpperSnakeCaseRule) ID() string {
	return "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE"
}

// Purpose returns the purpose of this rule.
func (r EnumFieldNamesUpperSnakeCaseRule) Purpose() string {
	return "Verifies that all enum field names are CAPITALS_WITH_UNDERSCORES."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumFieldNamesUpperSnakeCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r EnumFieldNamesUpperSnakeCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &enumFieldNamesUpperSnakeCaseVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type enumFieldNamesUpperSnakeCaseVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitEnumField checks the enum field.
func (v *enumFieldNamesUpperSnakeCaseVisitor) VisitEnumField(field *parser.EnumField) bool {
	if !strs.IsUpperSnakeCase(field.Ident) {
		v.AddFailuref(field.Meta.Pos, "EnumField name %q must be CAPITALS_WITH_UNDERSCORES", field.Ident)
	}
	return false
}
