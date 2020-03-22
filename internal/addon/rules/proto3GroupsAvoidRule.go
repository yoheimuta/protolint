package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// Proto3GroupsAvoidRule verifies that all groups should be avoided for proto3.
type Proto3GroupsAvoidRule struct {
}

// NewProto3GroupsAvoidRule creates a new Proto3GroupsAvoidRule.
func NewProto3GroupsAvoidRule() Proto3GroupsAvoidRule {
	return Proto3GroupsAvoidRule{}
}

// ID returns the ID of this rule.
func (r Proto3GroupsAvoidRule) ID() string {
	return "PROTO3_GROUPS_AVOID"
}

// Purpose returns the purpose of this rule.
func (r Proto3GroupsAvoidRule) Purpose() string {
	return "Verifies that all groups should be avoided for proto3."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r Proto3GroupsAvoidRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r Proto3GroupsAvoidRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &proto3GroupsAvoidVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type proto3GroupsAvoidVisitor struct {
	*visitor.BaseAddVisitor
	isProto3 bool
}

// VisitSyntax checks the syntax.
func (v *proto3GroupsAvoidVisitor) VisitSyntax(s *parser.Syntax) bool {
	v.isProto3 = s.ProtobufVersion == "proto3"
	return false
}

// VisitField checks the field.
func (v *proto3GroupsAvoidVisitor) VisitGroupField(field *parser.GroupField) bool {
	if v.isProto3 {
		v.AddFailuref(field.Meta.Pos, `Group %q should be avoided for proto3`, field.GroupName)
	}
	return false
}
