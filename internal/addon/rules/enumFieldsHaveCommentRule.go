package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// EnumFieldsHaveCommentRule verifies that all enumFields have a comment.
type EnumFieldsHaveCommentRule struct {
	RuleWithSeverity
	// Golang style comments should begin with the name of the thing being described.
	// See https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
	shouldFollowGolangStyle bool
}

// NewEnumFieldsHaveCommentRule creates a new EnumFieldsHaveCommentRule.
func NewEnumFieldsHaveCommentRule(
	severity rule.Severity,
	shouldFollowGolangStyle bool,
) EnumFieldsHaveCommentRule {
	return EnumFieldsHaveCommentRule{
		RuleWithSeverity:        RuleWithSeverity{severity: severity},
		shouldFollowGolangStyle: shouldFollowGolangStyle,
	}
}

// ID returns the ID of this rule.
func (r EnumFieldsHaveCommentRule) ID() string {
	return "ENUM_FIELDS_HAVE_COMMENT"
}

// Purpose returns the purpose of this rule.
func (r EnumFieldsHaveCommentRule) Purpose() string {
	return "Verifies that all enum fields have a comment."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r EnumFieldsHaveCommentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r EnumFieldsHaveCommentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &enumFieldsHaveCommentVisitor{
		BaseAddVisitor:          visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
		shouldFollowGolangStyle: r.shouldFollowGolangStyle,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type enumFieldsHaveCommentVisitor struct {
	*visitor.BaseAddVisitor
	shouldFollowGolangStyle bool
}

// VisitEnumField checks the enumField.
func (v *enumFieldsHaveCommentVisitor) VisitEnumField(enumField *parser.EnumField) bool {
	n := enumField.Ident
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(enumField.Comments, n) {
		v.AddFailuref(enumField.Meta.Pos, `EnumField %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComments(enumField.Comments, enumField.InlineComment) {
		v.AddFailuref(enumField.Meta.Pos, `EnumField %q should have a comment`, n)
	}
	return false
}
