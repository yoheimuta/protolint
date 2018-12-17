package cmd

import (
	"fmt"
	"io"

	"github.com/yoheimuta/protolinter/internal/cmd/subcmds/lint"
	"github.com/yoheimuta/protolinter/internal/cmd/subcmds/list"
	"github.com/yoheimuta/protolinter/internal/osutil"
)

const (
	help = `
Protocol Buffer Linter Command.

Usage:
  pl .
  pl lint .
  pl lint example.proto example2.proto
`
)

const (
	subCmdLint = "lint"
	subCmdList = "list"
)

// Do runs the command logic.
func Do(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	switch {
	case len(args) == 0:
		_, _ = fmt.Fprint(stderr, help)
		return osutil.ExitFailure
	default:
		return doSub(
			args,
			stdout,
			stderr,
		)
	}
}

func doSub(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	switch args[0] {
	case subCmdLint:
		if len(args) < 2 {
			_, _ = fmt.Fprintln(stderr, "pl lint requires at least one argument. See Usage.")
			_, _ = fmt.Fprint(stderr, help)
			return osutil.ExitFailure
		}
		return doLint(args[1:], stdout, stderr)
	case subCmdList:
		return doList(stdout, stderr)
	default:
		return doLint(args, stdout, stderr)
	}
}

func doLint(
	restArgs []string,
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	subCmd, err := lint.NewCmdLint(
		restArgs,
		stdout,
		stderr,
	)
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return osutil.ExitFailure
	}
	return subCmd.Run()
}

func doList(
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	subCmd := list.NewCmdList(
		stdout,
		stderr,
	)
	return subCmd.Run()
}
