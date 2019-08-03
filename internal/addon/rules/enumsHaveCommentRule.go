package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

// EnumsHaveCommentRule verifies that all enums have a comment.
type EnumsHaveCommentRule struct {
	// Golang style comments should begin with the name of the thing being described.
	// See https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
	shouldFollowGolangStyle bool
}

// NewEnumsHaveCommentRule creates a new EnumsHaveCommentRule.
func NewEnumsHaveCommentRule(
	shouldFollowGolangStyle bool,
) EnumsHaveCommentRule {
	return EnumsHaveCommentRule{
		shouldFollowGolangStyle: shouldFollowGolangStyle,
	}
}

// ID returns the ID of this rule.
func (r EnumsHaveCommentRule) ID() string {
	return "ENUMS_HAVE_COMMENT"
}

// Purpose returns the purpose of this rule.
func (r EnumsHaveCommentRule) Purpose() string {
	return "Verifies that all enums have a comment."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumsHaveCommentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r EnumsHaveCommentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &enumsHaveCommentVisitor{
		BaseAddVisitor:          visitor.NewBaseAddVisitor(),
		shouldFollowGolangStyle: r.shouldFollowGolangStyle,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type enumsHaveCommentVisitor struct {
	*visitor.BaseAddVisitor
	shouldFollowGolangStyle bool
}

// VisitEnum checks the enum.
func (v *enumsHaveCommentVisitor) VisitEnum(enum *parser.Enum) bool {
	n := enum.EnumName
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(enum.Comments, n) {
		v.AddFailuref(enum.Meta.Pos, `Enum %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComment(enum.Comments) {
		v.AddFailuref(enum.Meta.Pos, `Enum %q should have a comment`, n)
	}
	return true
}
