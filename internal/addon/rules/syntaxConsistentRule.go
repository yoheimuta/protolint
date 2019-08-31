package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// SyntaxConsistentRule verifies that syntax is a specified version.
type SyntaxConsistentRule struct {
	version string
}

// NewSyntaxConsistentRule creates a new SyntaxConsistentRule.
func NewSyntaxConsistentRule(
	version string,
) SyntaxConsistentRule {
	if len(version) == 0 {
		version = "proto3"
	}
	return SyntaxConsistentRule{
		version: version,
	}
}

// ID returns the ID of this rule.
func (r SyntaxConsistentRule) ID() string {
	return "SYNTAX_CONSISTENT"
}

// Purpose returns the purpose of this rule.
func (r SyntaxConsistentRule) Purpose() string {
	return "Verifies that syntax is a specified version(default is proto3)."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r SyntaxConsistentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r SyntaxConsistentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &syntaxConsistentVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		version:        r.version,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type syntaxConsistentVisitor struct {
	*visitor.BaseAddVisitor
	version string
}

// VisitSyntax checks the syntax.
func (v *syntaxConsistentVisitor) VisitSyntax(s *parser.Syntax) bool {
	if s.ProtobufVersion != v.version {
		v.AddFailuref(s.Meta.Pos, "Syntax should be %q but was %q.", v.version, s.ProtobufVersion)
	}
	return false
}
