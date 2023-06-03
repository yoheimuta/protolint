package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/internal/linter/config"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// RPCNamesCaseRule verifies that all rpc names conform to the specified convention.
type RPCNamesCaseRule struct {
	RuleWithSeverity
	convention config.ConventionType
}

// NewRPCNamesCaseRule creates a new RPCNamesCaseRule.
func NewRPCNamesCaseRule(
	severity rule.Severity,
	convention config.ConventionType,
) RPCNamesCaseRule {
	return RPCNamesCaseRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		convention:       convention,
	}
}

// ID returns the ID of this rule.
func (r RPCNamesCaseRule) ID() string {
	return "RPC_NAMES_CASE"
}

// Purpose returns the purpose of this rule.
func (r RPCNamesCaseRule) Purpose() string {
	return "Verifies that all rpc names conform to the specified convention."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RPCNamesCaseRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r RPCNamesCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &rpcNamesCaseVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		convention:     r.convention,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type rpcNamesCaseVisitor struct {
	*visitor.BaseAddVisitor
	convention config.ConventionType
}

// VisitRPC checks the rpc.
func (v *rpcNamesCaseVisitor) VisitRPC(rpc *parser.RPC) bool {
	if v.convention == config.ConventionLowerCamel && !strs.IsLowerCamelCase(rpc.RPCName) {
		v.AddFailuref(rpc.Meta.Pos, "RPC name %q must be LowerCamelCase", rpc.RPCName)
	} else if v.convention == config.ConventionUpperSnake && !strs.IsUpperSnakeCase(rpc.RPCName) {
		v.AddFailuref(rpc.Meta.Pos, "RPC name %q must be UpperSnakeCase", rpc.RPCName)
	} else if v.convention == config.ConventionLowerSnake && !strs.IsLowerSnakeCase(rpc.RPCName) {
		v.AddFailuref(rpc.Meta.Pos, "RPC name %q must be LowerSnakeCase", rpc.RPCName)
	}
	return false
}
