package lint

import (
	"flag"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds"
	"github.com/yoheimuta/protolint/linter/autodisable"

	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"

	"github.com/yoheimuta/protolint/internal/linter/report/reporters"

	"github.com/yoheimuta/protolint/internal/linter/report"
)

// Flags represents a set of lint flag parameters.
type Flags struct {
	*flag.FlagSet

	FilePaths                 []string
	ConfigPath                string
	ConfigDirPath             string
	FixMode                   bool
	Reporter                  report.Reporter
	AutoDisableType           autodisable.PlacementType
	OutputFilePath            string
	Verbose                   bool
	NoErrorOnUnmatchedPattern bool
	Plugins                   []shared.RuleSet
}

// NewFlags creates a new Flags.
func NewFlags(
	args []string,
) (Flags, error) {
	f := Flags{
		FlagSet:         flag.NewFlagSet("lint", flag.ExitOnError),
		Reporter:        reporters.PlainReporter{},
		AutoDisableType: autodisable.Noop,
	}
	var rf reporterFlag
	var af autoDisableFlag
	var pf subcmds.PluginFlag

	f.StringVar(
		&f.ConfigPath,
		"config_path",
		"",
		"path/to/protolint.yaml. Note that if both are set, config_dir_path is ignored.",
	)
	f.StringVar(
		&f.ConfigDirPath,
		"config_dir_path",
		"",
		"path/to/the_directory_including_protolint.yaml",
	)
	f.BoolVar(
		&f.FixMode,
		"fix",
		false,
		"mode that the command line automatically fix some of the problems",
	)
	f.Var(
		&rf,
		"reporter",
		`formatter to output results in the specific format. Available reporters are "plain"(default), "junit", "json", and "unix".`,
	)
	f.Var(
		&af,
		"auto_disable",
		`mode that the command line automatically disable some of the problems. Available auto_disable are "next" and "this".`,
	)
	f.StringVar(
		&f.OutputFilePath,
		"output_file",
		"",
		"path/to/output.txt",
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
	f.BoolVar(
		&f.NoErrorOnUnmatchedPattern,
		"no-error-on-unmatched-pattern",
		false,
		"exits with 0 when no file is matched",
	)

	_ = f.Parse(args)
	if rf.reporter != nil {
		f.Reporter = rf.reporter
	}
	if af.autoDisableType != 0 {
		f.AutoDisableType = af.autoDisableType
	}

	plugins, err := pf.BuildPlugins(f.Verbose)
	if err != nil {
		return Flags{}, err
	}
	f.Plugins = plugins

	f.FilePaths = f.Args()
	return f, nil
}
