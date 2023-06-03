package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// ServicesHaveCommentRule verifies that all services have a comment.
type ServicesHaveCommentRule struct {
	RuleWithSeverity
	// Golang style comments should begin with the name of the thing being described.
	// See https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
	shouldFollowGolangStyle bool
}

// NewServicesHaveCommentRule creates a new ServicesHaveCommentRule.
func NewServicesHaveCommentRule(
	severity rule.Severity,
	shouldFollowGolangStyle bool,
) ServicesHaveCommentRule {
	return ServicesHaveCommentRule{
		RuleWithSeverity:        RuleWithSeverity{severity: severity},
		shouldFollowGolangStyle: shouldFollowGolangStyle,
	}
}

// ID returns the ID of this rule.
func (r ServicesHaveCommentRule) ID() string {
	return "SERVICES_HAVE_COMMENT"
}

// Purpose returns the purpose of this rule.
func (r ServicesHaveCommentRule) Purpose() string {
	return "Verifies that all services have a comment."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r ServicesHaveCommentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r ServicesHaveCommentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &servicesHaveCommentVisitor{
		BaseAddVisitor:          visitor.NewBaseAddVisitor(r.ID()),
		shouldFollowGolangStyle: r.shouldFollowGolangStyle,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type servicesHaveCommentVisitor struct {
	*visitor.BaseAddVisitor
	shouldFollowGolangStyle bool
}

// VisitService checks the service.
func (v *servicesHaveCommentVisitor) VisitService(service *parser.Service) bool {
	n := service.ServiceName
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(service.Comments, n) {
		v.AddFailuref(service.Meta.Pos, `Service %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComments(service.Comments, service.InlineComment, service.InlineCommentBehindLeftCurly) {
		v.AddFailuref(service.Meta.Pos, `Service %q should have a comment`, n)
	}
	return false
}
