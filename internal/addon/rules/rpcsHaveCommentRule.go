package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

// RPCsHaveCommentRule verifies that all rpcs have a comment.
type RPCsHaveCommentRule struct {
	// Golang style comments should begin with the name of the thing being described.
	// See https://github.com/golang/go/wiki/CodeReviewComments#comment-sentences
	shouldFollowGolangStyle bool
}

// NewRPCsHaveCommentRule creates a new RPCsHaveCommentRule.
func NewRPCsHaveCommentRule(
	shouldFollowGolangStyle bool,
) RPCsHaveCommentRule {
	return RPCsHaveCommentRule{
		shouldFollowGolangStyle: shouldFollowGolangStyle,
	}
}

// ID returns the ID of this rule.
func (r RPCsHaveCommentRule) ID() string {
	return "RPCS_HAVE_COMMENT"
}

// Purpose returns the purpose of this rule.
func (r RPCsHaveCommentRule) Purpose() string {
	return "Verifies that all rpcs have a comment."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RPCsHaveCommentRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r RPCsHaveCommentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &rpcsHaveCommentVisitor{
		BaseAddVisitor:          visitor.NewBaseAddVisitor(),
		shouldFollowGolangStyle: r.shouldFollowGolangStyle,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type rpcsHaveCommentVisitor struct {
	*visitor.BaseAddVisitor
	shouldFollowGolangStyle bool
}

// VisitRPC checks the rpc.
func (v *rpcsHaveCommentVisitor) VisitRPC(rpc *parser.RPC) bool {
	n := rpc.RPCName
	if v.shouldFollowGolangStyle && !hasGolangStyleComment(rpc.Comments, n) {
		v.AddFailuref(rpc.Meta.Pos, `RPC %q should have a comment of the form "// %s ..."`, n, n)
	} else if !hasComment(rpc.Comments) {
		v.AddFailuref(rpc.Meta.Pos, `RPC %q should have a comment`, n)
	}
	return false
}
