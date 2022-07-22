package list

import (
	"flag"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds"

	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"

	"github.com/yoheimuta/protolint/internal/linter/report/reporters"

	"github.com/yoheimuta/protolint/internal/linter/report"
)

// Flags represents a set of lint flag parameters.
type Flags struct {
	*flag.FlagSet

	Reporter report.Reporter
	Plugins  []shared.RuleSet
}

// NewFlags creates a new Flags.
func NewFlags(
	args []string,
) (Flags, error) {
	f := Flags{
		FlagSet:  flag.NewFlagSet("lint", flag.ExitOnError),
		Reporter: reporters.PlainReporter{},
	}
	var pf subcmds.PluginFlag

	f.Var(
		&pf,
		"plugin",
		`plugins to provide custom lint rule set. Note that it's necessary to specify it as path format'`,
	)

	_ = f.Parse(args)

	plugins, err := pf.BuildPlugins(false)
	if err != nil {
		return Flags{}, err
	}
	f.Plugins = plugins
	return f, nil
}
