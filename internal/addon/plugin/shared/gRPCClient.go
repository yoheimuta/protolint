package shared

import (
	"context"

	"github.com/yoheimuta/protolint/internal/addon/plugin/proto"
)

// GRPCClient is the implementation of RuleSet.
type GRPCClient struct {
	client proto.RuleSetServiceClient
}

// ListRules returns all supported rules metadata.
func (c *GRPCClient) ListRules(req *proto.ListRulesRequest) (*proto.ListRulesResponse, error) {
	return c.client.ListRules(context.Background(), req)
}

// Apply applies the rule to the proto.
func (c *GRPCClient) Apply(req *proto.ApplyRequest) (*proto.ApplyResponse, error) {
	return c.client.Apply(context.Background(), req)
}
