package shared

import (
	"context"

	"github.com/yoheimuta/protolint/internal/addon/plugin/proto"
)

// GRPCServer is the implementation of RuleSet.
type GRPCServer struct {
	proto.UnimplementedRuleSetServiceServer
	server RuleSet
}

// ListRules returns all supported rules metadata.
func (s *GRPCServer) ListRules(_ context.Context, req *proto.ListRulesRequest) (*proto.ListRulesResponse, error) {
	return s.server.ListRules(req)
}

// Apply applies the rule to the proto.
func (s *GRPCServer) Apply(_ context.Context, req *proto.ApplyRequest) (*proto.ApplyResponse, error) {
	return s.server.Apply(req)
}
