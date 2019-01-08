package lint

import "flag"

// Flags represents a set of lint flag parameters.
type Flags struct {
	*flag.FlagSet

	ConfigDirPath string
}

// NewFlags creates a new Flags.
func NewFlags(
	args []string,
) Flags {
	f := Flags{
		FlagSet: flag.NewFlagSet("lint", flag.ExitOnError),
	}

	f.StringVar(&f.ConfigDirPath, "config_dir_path", "", "path/to/protolint.yaml")
	_ = f.Parse(args)
	return f
}
