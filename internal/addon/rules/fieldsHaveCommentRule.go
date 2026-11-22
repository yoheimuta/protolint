package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// FieldsHaveCommentRule verifies that all fields have a comment.
type FieldsHaveCommentRule struct {
	RuleWithSeverity
	// Golang style comments should begin with the name of the thing being described.
	// See https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
	shouldFollowGolangStyle bool
	// Field names that are exempt from this rule.
	excludes []string
}

// NewFieldsHaveCommentRule creates a new FieldsHaveCommentRule.
func NewFieldsHaveCommentRule(
	severity rule.Severity,
	shouldFollowGolangStyle bool,
	excludes []string,
) FieldsHaveCommentRule {
	return FieldsHaveCommentRule{
		RuleWithSeverity:        RuleWithSeverity{severity: severity},
		shouldFollowGolangStyle: shouldFollowGolangStyle,
		excludes:                excludes,
	}
}

// ID returns the ID of this rule.
func (r FieldsHaveCommentRule) ID() string {
	return "FIELDS_HAVE_COMMENT"
}

// Purpose returns the purpose of this rule.
func (r FieldsHaveCommentRule) Purpose() string {
	return "Verifies that all fields have a comment."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r FieldsHaveCommentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r FieldsHaveCommentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &fieldsHaveCommentVisitor{
		BaseAddVisitor:          visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
		shouldFollowGolangStyle: r.shouldFollowGolangStyle,
		excludes:                r.excludes,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type fieldsHaveCommentVisitor struct {
	*visitor.BaseAddVisitor
	shouldFollowGolangStyle bool
	excludes                []string
}

func (v *fieldsHaveCommentVisitor) isExcluded(name string) bool {
	for _, e := range v.excludes {
		if e == name {
			return true
		}
	}
	return false
}

// VisitField checks the field.
func (v *fieldsHaveCommentVisitor) VisitField(field *parser.Field) bool {
	n := field.FieldName
	if v.isExcluded(n) {
		return false
	}
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(field.Comments, n) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComments(field.Comments, field.InlineComment) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment`, n)
	}
	return false
}

// VisitMapField checks the map field.
func (v *fieldsHaveCommentVisitor) VisitMapField(field *parser.MapField) bool {
	n := field.MapName
	if v.isExcluded(n) {
		return false
	}
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(field.Comments, n) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComments(field.Comments, field.InlineComment) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment`, n)
	}
	return false
}

// VisitOneofField checks the oneof field.
func (v *fieldsHaveCommentVisitor) VisitOneofField(field *parser.OneofField) bool {
	n := field.FieldName
	if v.isExcluded(n) {
		return false
	}
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(field.Comments, n) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComments(field.Comments, field.InlineComment) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment`, n)
	}
	return false
}
