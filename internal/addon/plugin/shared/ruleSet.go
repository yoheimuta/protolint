package shared

import "github.com/yoheimuta/protolint/internal/addon/plugin/proto"

// RuleSet is the interface that we're exposing as a plugin.
type RuleSet interface {
	ListRules(*proto.ListRulesRequest) (*proto.ListRulesResponse, error)
	Apply(*proto.ApplyRequest) (*proto.ApplyResponse, error)
}
