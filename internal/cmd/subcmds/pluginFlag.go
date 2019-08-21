package subcmds

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"

	"github.com/hashicorp/go-plugin"
)

// PluginFlag manages a flag for plugins.
type PluginFlag struct {
	raws    []string
	plugins []shared.RuleSet
}

// String implements flag.Value.
func (f *PluginFlag) String() string {
	return fmt.Sprint(strings.Join(f.raws, ","))
}

// Set implements flag.Value.
func (f *PluginFlag) Set(value string) error {
	f.raws = append(f.raws, value)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", value),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})

	rpcClient, err := client.Client()
	if err != nil {
		return fmt.Errorf("failed client.Client(), err=%s", err)
	}

	ruleSet, err := rpcClient.Dispense("ruleSet")
	if err != nil {
		return fmt.Errorf("failed Dispense, err=%s", err)
	}
	f.plugins = append(f.plugins, ruleSet.(shared.RuleSet))
	return nil
}

// Plugins returns all plugins.
func (f *PluginFlag) Plugins() []shared.RuleSet {
	return f.plugins
}
