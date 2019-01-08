package lint

import (
	"github.com/yoheimuta/protolint/internal/cmd/subcmds"
	"github.com/yoheimuta/protolint/internal/linter/config"
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

// CmdLintConfig is a config for lint command.
type CmdLintConfig struct {
	external config.ExternalConfig
}

// NewCmdLintConfig creates a new CmdLintConfig.
func NewCmdLintConfig(
	externalConfig config.ExternalConfig,
) CmdLintConfig {
	return CmdLintConfig{
		external: externalConfig,
	}
}

// GenRules generates rules which are applied to the filename path.
func (c CmdLintConfig) GenRules(
	f file.ProtoFile,
) ([]rule.HasApply, error) {
	var hasApplies []rule.HasApply
	for _, r := range subcmds.NewAllRules() {
		if c.external.SkipRule(r.ID(), f.DisplayPath(), subcmds.DefaultRuleIDs()) {
			continue
		}
		hasApplies = append(hasApplies, r)
	}

	return hasApplies, nil
}
