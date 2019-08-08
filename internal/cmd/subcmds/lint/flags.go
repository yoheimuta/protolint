package lint

import (
	"flag"

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
		"formatter to output results in the specific format. Available reporters are 'plain'(default).",
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
	return f
}
