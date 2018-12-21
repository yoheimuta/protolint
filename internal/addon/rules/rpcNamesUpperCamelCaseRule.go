package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolinter/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolinter/internal/linter/report"
	"github.com/yoheimuta/protolinter/internal/strs"
)

// RPCNamesUpperCamelCaseRule verifies that all rpc names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#services.
type RPCNamesUpperCamelCaseRule struct{}

// NewRPCNamesUpperCamelCaseRule creates a new RPCNamesUpperCamelCaseRule.
func NewRPCNamesUpperCamelCaseRule() RPCNamesUpperCamelCaseRule {
	return RPCNamesUpperCamelCaseRule{}
}

// ID returns the ID of this rule.
func (r RPCNamesUpperCamelCaseRule) ID() string {
	return "RPC_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r RPCNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all rpc names are CamelCase (with an initial capital)."
}

// Apply applies the rule to the proto.
func (r RPCNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &rpcNamesUpperCamelCaseVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(),
	}
	return visitor.RunVisitor(v, proto)
}

type rpcNamesUpperCamelCaseVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitRPC checks the rpc.
func (v *rpcNamesUpperCamelCaseVisitor) VisitRPC(rpc *parser.RPC) bool {
	if !strs.IsUpperCamelCase(rpc.RPCName) {
		v.AddFailuref(rpc.Meta.Pos, "RPC name %q must be UpperCamelCase", rpc.RPCName)
	}
	return false
}
