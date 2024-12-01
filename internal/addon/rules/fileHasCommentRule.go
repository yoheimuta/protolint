package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// FileHasCommentRule verifies that a file starts with a doc comment.
type FileHasCommentRule struct {
	RuleWithSeverity
}

// NewFileHasCommentRule creates a new FileHasCommentRule.
func NewFileHasCommentRule(severity rule.Severity) FileHasCommentRule {
	return FileHasCommentRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
	}
}

// ID returns the ID of this rule.
func (r FileHasCommentRule) ID() string {
	return "FILE_HAS_COMMENT"
}

// Purpose returns the purpose of this rule.
func (r FileHasCommentRule) Purpose() string {
	return "Verifies that a file starts with a doc comment."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r FileHasCommentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r FileHasCommentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &fileHasCommentVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type fileHasCommentVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitSyntax checks the syntax.
func (v *fileHasCommentVisitor) VisitSyntax(s *parser.Syntax) bool {
	if !hasComment(s.Comments) {
		v.AddFailuref(s.Meta.Pos, `File should start with a doc comment`)
	}
	return false
}

// VisitEdition checks the syntax.
func (v *fileHasCommentVisitor) VisitEdition(s *parser.Edition) bool {
	if !hasComment(s.Comments) {
		v.AddFailuref(s.Meta.Pos, `File should start with a doc comment`)
	}
	return false
}
