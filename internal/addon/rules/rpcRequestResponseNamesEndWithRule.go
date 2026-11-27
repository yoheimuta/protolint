package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

const (
	defaultRPCRequestSuffix  = "Request"
	defaultRPCResponseSuffix = "Response"
)

// RPCRequestResponseNamesEndWithRule verifies that all rpc request and response message names end with the specified values.
type RPCRequestResponseNamesEndWithRule struct {
	RuleWithSeverity
	requestSuffix  string
	responseSuffix string
	excludes       []string
}

// NewRPCRequestResponseNamesEndWithRule creates a new RPCRequestResponseNamesEndWithRule.
func NewRPCRequestResponseNamesEndWithRule(
	severity rule.Severity,
	requestSuffix string,
	responseSuffix string,
	excludes []string,
) RPCRequestResponseNamesEndWithRule {
	if len(requestSuffix) == 0 {
		requestSuffix = defaultRPCRequestSuffix
	}
	if len(responseSuffix) == 0 {
		responseSuffix = defaultRPCResponseSuffix
	}
	return RPCRequestResponseNamesEndWithRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		requestSuffix:    requestSuffix,
		responseSuffix:   responseSuffix,
		excludes:         excludes,
	}
}

// ID returns the ID of this rule.
func (r RPCRequestResponseNamesEndWithRule) ID() string {
	return "RPC_REQUEST_RESPONSE_NAMES_END_WITH"
}

// Purpose returns the purpose of this rule.
func (r RPCRequestResponseNamesEndWithRule) Purpose() string {
	return "Verifies that all rpc request and response message names end with the specified suffixes."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RPCRequestResponseNamesEndWithRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r RPCRequestResponseNamesEndWithRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &rpcRequestResponseNamesEndWithVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
		requestSuffix:  r.requestSuffix,
		responseSuffix: r.responseSuffix,
		excludes:       r.excludes,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type rpcRequestResponseNamesEndWithVisitor struct {
	*visitor.BaseAddVisitor
	requestSuffix  string
	responseSuffix string
	excludes       []string
}

// isExcluded reports whether the message type is in excludes.
// A leading dot of a fully qualified name is ignored, so ".google.protobuf.Empty" matches "google.protobuf.Empty".
func (v *rpcRequestResponseNamesEndWithVisitor) isExcluded(messageType string) bool {
	name := strings.TrimPrefix(messageType, ".")
	for _, e := range v.excludes {
		if strings.TrimPrefix(e, ".") == name {
			return true
		}
	}
	return false
}

// VisitRPC checks the rpc.
func (v *rpcRequestResponseNamesEndWithVisitor) VisitRPC(rpc *parser.RPC) bool {
	if req := rpc.RPCRequest; req != nil {
		if !v.isExcluded(req.MessageType) && !strings.HasSuffix(req.MessageType, v.requestSuffix) {
			v.AddFailuref(req.Meta.Pos, "RPC %q request message %q must end with %s", rpc.RPCName, req.MessageType, v.requestSuffix)
		}
	}
	if res := rpc.RPCResponse; res != nil {
		if !v.isExcluded(res.MessageType) && !strings.HasSuffix(res.MessageType, v.responseSuffix) {
			v.AddFailuref(res.Meta.Pos, "RPC %q response message %q must end with %s", rpc.RPCName, res.MessageType, v.responseSuffix)
		}
	}
	return false
}
