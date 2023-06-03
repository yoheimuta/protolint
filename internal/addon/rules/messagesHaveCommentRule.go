package rules

import (
	"fmt"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// MessagesHaveCommentRule verifies that all messages have a comment.
type MessagesHaveCommentRule struct {
	RuleWithSeverity
	// Golang style comments should begin with the name of the thing being described.
	// See https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
	shouldFollowGolangStyle bool
}

// NewMessagesHaveCommentRule creates a new MessagesHaveCommentRule.
func NewMessagesHaveCommentRule(
	severity rule.Severity,
	shouldFollowGolangStyle bool,
) MessagesHaveCommentRule {
	return MessagesHaveCommentRule{
		RuleWithSeverity:        RuleWithSeverity{severity: severity},
		shouldFollowGolangStyle: shouldFollowGolangStyle,
	}
}

// ID returns the ID of this rule.
func (r MessagesHaveCommentRule) ID() string {
	return "MESSAGES_HAVE_COMMENT"
}

// Purpose returns the purpose of this rule.
func (r MessagesHaveCommentRule) Purpose() string {
	return "Verifies that all messages have a comment."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r MessagesHaveCommentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r MessagesHaveCommentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &messagesHaveCommentVisitor{
		BaseAddVisitor:          visitor.NewBaseAddVisitor(r.ID()),
		shouldFollowGolangStyle: r.shouldFollowGolangStyle,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type messagesHaveCommentVisitor struct {
	*visitor.BaseAddVisitor
	shouldFollowGolangStyle bool
}

// VisitMessage checks the message.
func (v *messagesHaveCommentVisitor) VisitMessage(message *parser.Message) bool {
	n := message.MessageName
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(message.Comments, n) {
		v.AddFailuref(message.Meta.Pos, `Message %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComments(message.Comments, message.InlineComment, message.InlineCommentBehindLeftCurly) {
		v.AddFailuref(message.Meta.Pos, `Message %q should have a comment`, n)
	}
	return true
}

func hasGolangStyleComment(
	comments []*parser.Comment,
	describedName string,
) bool {
	return hasComment(comments) &&
		strings.HasPrefix(comments[0].Lines()[0], fmt.Sprintf(" %s", describedName))
}

func hasComment(comments []*parser.Comment) bool {
	return 0 < len(comments)
}

func hasComments(comments []*parser.Comment, inlines ...*parser.Comment) bool {
	if 0 < len(comments) {
		return true
	}
	for _, inline := range inlines {
		if inline != nil {
			return true
		}
	}
	return false
}
