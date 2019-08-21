package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// FieldsHaveCommentRule verifies that all fields have a comment.
type FieldsHaveCommentRule struct {
	// Golang style comments should begin with the name of the thing being described.
	// See https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
	shouldFollowGolangStyle bool
}

// NewFieldsHaveCommentRule creates a new FieldsHaveCommentRule.
func NewFieldsHaveCommentRule(
	shouldFollowGolangStyle bool,
) FieldsHaveCommentRule {
	return FieldsHaveCommentRule{
		shouldFollowGolangStyle: shouldFollowGolangStyle,
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
		BaseAddVisitor:          visitor.NewBaseAddVisitor(r.ID()),
		shouldFollowGolangStyle: r.shouldFollowGolangStyle,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type fieldsHaveCommentVisitor struct {
	*visitor.BaseAddVisitor
	shouldFollowGolangStyle bool
}

// VisitField checks the field.
func (v *fieldsHaveCommentVisitor) VisitField(field *parser.Field) bool {
	n := field.FieldName
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(field.Comments, n) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComment(field.Comments) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment`, n)
	}
	return false
}

// VisitMapField checks the map field.
func (v *fieldsHaveCommentVisitor) VisitMapField(field *parser.MapField) bool {
	n := field.MapName
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(field.Comments, n) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComment(field.Comments) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment`, n)
	}
	return false
}

// VisitOneofField checks the oneof field.
func (v *fieldsHaveCommentVisitor) VisitOneofField(field *parser.OneofField) bool {
	n := field.FieldName
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(field.Comments, n) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComment(field.Comments) {
		v.AddFailuref(field.Meta.Pos, `Field %q should have a comment`, n)
	}
	return false
}
