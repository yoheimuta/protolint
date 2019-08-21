package shared

import (
	"context"

	"github.com/yoheimuta/protolint/internal/addon/plugin/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"ruleSet": &RuleSetGRPCPlugin{},
}

// RuleSetGRPCPlugin is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type RuleSetGRPCPlugin struct {
	plugin.Plugin
	Impl RuleSet
}

// GRPCServer registers this plugin.
func (p *RuleSetGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterRuleSetServiceServer(s, &GRPCServer{server: p.Impl})
	return nil
}

// GRPCClient returns the interface implementation for the plugin.
func (p *RuleSetGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewRuleSetServiceClient(c)}, nil
}
