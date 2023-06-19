package lint

import (
	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"
	"github.com/yoheimuta/protolint/internal/cmd/subcmds"
	"github.com/yoheimuta/protolint/internal/linter/config"
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/rule"
)

// CmdLintConfig is a config for lint command.
type CmdLintConfig struct {
	external        config.ExternalConfig
	fixMode         bool
	autoDisableType autodisable.PlacementType
	verbose         bool
	reporters       report.ReportersWithOutput
	plugins         []shared.RuleSet
}

// NewCmdLintConfig creates a new CmdLintConfig.
func NewCmdLintConfig(
	externalConfig config.ExternalConfig,
	flags Flags,
) CmdLintConfig {
	output := report.WriteToConsole
	if 0 < len(flags.OutputFilePath) {
		output = flags.OutputFilePath
	}

	var reporters report.ReportersWithOutput
	reporters = append(reporters, *report.NewReporterWithOutput(flags.Reporter, output))

	for _, additionalReporter := range flags.AdditionalReporters {
		r := *report.NewReporterWithOutput(additionalReporter.reporter, additionalReporter.targetFile)
		reporters = append(reporters, r)
	}

	return CmdLintConfig{
		external:        externalConfig,
		fixMode:         flags.FixMode,
		autoDisableType: flags.AutoDisableType,
		verbose:         flags.Verbose,
		reporters:       reporters,
		plugins:         flags.Plugins,
	}
}

// GenRules generates rules which are applied to the filename path.
func (c CmdLintConfig) GenRules(
	f file.ProtoFile,
) ([]rule.HasApply, error) {
	allRules, err := subcmds.NewAllRules(c.external.Lint.RulesOption, c.fixMode, c.autoDisableType, c.verbose, c.plugins)
	if err != nil {
		return nil, err
	}

	var defaultRuleIDs []string
	if c.external.Lint.Rules.AllDefault {
		defaultRuleIDs = allRules.IDs()
	} else {
		defaultRuleIDs = allRules.Default().IDs()
	}

	var hasApplies []rule.HasApply
	for _, r := range allRules {
		if c.external.ShouldSkipRule(r.ID(), f.DisplayPath(), defaultRuleIDs) {
			continue
		}
		hasApplies = append(hasApplies, r)
	}

	return hasApplies, nil
}
