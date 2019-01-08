package lint

import (
	"github.com/yoheimuta/protolint/internal/cmd/subcmds"
	"github.com/yoheimuta/protolint/internal/linter/config"
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

// CmdLintConfig is a config for lint command.
type CmdLintConfig struct {
	c config.ExternalConfig
}

// NewCmdLintConfig creates a new CmdLintConfig.
func NewCmdLintConfig(
	externalConfig config.ExternalConfig,
) CmdLintConfig {
	return CmdLintConfig{
		c: externalConfig,
	}
}

// GenRules generates rules which are applied to the filename path.
func (c CmdLintConfig) GenRules(
	f file.ProtoFile,
) ([]rule.HasApply, error) {
	var hasApplies []rule.HasApply
	for _, r := range subcmds.NewAllRules() {
		hasApplies = append(hasApplies, r)
	}
	return hasApplies, nil
}
