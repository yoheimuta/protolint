package lint

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

	ConfigDirPath string
	FixMode       bool
	Reporter      report.Reporter
	Verbose       bool
	Plugins       []shared.RuleSet
}

// NewFlags creates a new Flags.
func NewFlags(
	args []string,
) Flags {
	f := Flags{
		FlagSet:  flag.NewFlagSet("lint", flag.ExitOnError),
		Reporter: reporters.PlainReporter{},
	}
	var rf reporterFlag
	var pf subcmds.PluginFlag

	f.StringVar(
		&f.ConfigDirPath,
		"config_dir_path",
		"",
		"path/to/protolint.yaml",
	)
	f.BoolVar(
		&f.FixMode,
		"fix",
		false,
		"mode that the command line can automatically fix some of the problems",
	)
	f.Var(
		&rf,
		"reporter",
		`formatter to output results in the specific format. Available reporters are "plain"(default) and "junit".`,
	)
	f.Var(
		&pf,
		"plugin",
		`plugins to provide custom lint rule set. Note that it's necessary to specify it as path format'`,
	)
	f.BoolVar(
		&f.Verbose,
		"v",
		false,
		"verbose output that includes parsing process details",
	)

	_ = f.Parse(args)
	if rf.reporter != nil {
		f.Reporter = rf.reporter
	}
	f.Plugins = pf.Plugins()
	return f
}
