package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

const (
	defaultSuffix = "UNSPECIFIED"
)

// EnumFieldNamesZeroValueEndWithRule verifies that the zero value enum should have the suffix (e.g. "UNSPECIFIED", "INVALID").
// See https://developers.google.com/protocol-buffers/docs/style#enums.
type EnumFieldNamesZeroValueEndWithRule struct {
	suffix string
}

// NewEnumFieldNamesZeroValueEndWithRule creates a new EnumFieldNamesZeroValueEndWithRule.
func NewEnumFieldNamesZeroValueEndWithRule(
	suffix string,
) EnumFieldNamesZeroValueEndWithRule {
	if len(suffix) == 0 {
		suffix = defaultSuffix
	}
	return EnumFieldNamesZeroValueEndWithRule{
		suffix: suffix,
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

// Apply applies the rule to the proto.
func (r EnumFieldNamesZeroValueEndWithRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &enumFieldNamesZeroValueEndWithVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(),
		suffix:         r.suffix,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type enumFieldNamesZeroValueEndWithVisitor struct {
	*visitor.BaseAddVisitor
	suffix string
}

// VisitEnumField checks the enum field.
func (v *enumFieldNamesZeroValueEndWithVisitor) VisitEnumField(field *parser.EnumField) bool {
	if field.Number == "0" && !strings.HasSuffix(field.Ident, v.suffix) {
		v.AddFailuref(field.Meta.Pos, "EnumField name %q with zero value should have the suffix %q", field.Ident, v.suffix)
	}
	return false
}
