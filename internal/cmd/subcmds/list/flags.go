package list

import (
	"flag"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds"

	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"
)

// Flags represents a set of lint flag parameters.
type Flags struct {
	*flag.FlagSet

	Plugins []shared.RuleSet
}

// NewFlags creates a new Flags.
func NewFlags(
	args []string,
) (Flags, error) {
	f := Flags{
		FlagSet: flag.NewFlagSet("list", flag.ExitOnError),
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
