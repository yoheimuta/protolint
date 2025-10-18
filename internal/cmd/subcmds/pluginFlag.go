package subcmds

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"

	"github.com/hashicorp/go-plugin"
)

// PluginFlag manages a flag for plugins.
type PluginFlag struct {
	raws []string
}

// String implements flag.Value.
func (f *PluginFlag) String() string {
	return fmt.Sprint(strings.Join(f.raws, ","))
}

// Set implements flag.Value.
func (f *PluginFlag) Set(value string) error {
	f.raws = append(f.raws, value)
	return nil
}

// BuildPlugins builds all plugins.
func (f *PluginFlag) BuildPlugins(verbose bool) ([]shared.RuleSet, error) {
	var plugins []shared.RuleSet

	for _, value := range f.raws {
		level := hclog.Warn
		if verbose {
			level = hclog.Trace
		}
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: shared.Handshake,
			Plugins:         shared.PluginMap,
			Cmd:             exec.Command(value),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolGRPC,
			},
			Logger: hclog.New(&hclog.LoggerOptions{
				Output: hclog.DefaultOutput,
				Level:  level,
				Name:   "plugin",
			}),
			// To cleanup. See. https://github.com/yoheimuta/protolint/issues/237
			Managed: true,
		})

		rpcClient, err := client.Client()
		if err != nil {
			return nil, fmt.Errorf("failed client.Client(), err=%s", err)
		}

		ruleSet, err := rpcClient.Dispense("ruleSet")
		if err != nil {
			return nil, fmt.Errorf("failed Dispense, err=%s", err)
		}
		plugins = append(plugins, ruleSet.(shared.RuleSet))
	}
	return plugins, nil
}
