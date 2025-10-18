package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds/lint"
	"github.com/yoheimuta/protolint/internal/cmd/subcmds/list"
	"github.com/yoheimuta/protolint/internal/osutil"
	"github.com/yoheimuta/protolint/mcp"
)

const (
	help = `
Protocol Buffer Linter Command.

Usage:
	protolint <command> [arguments]
	protolint --version
	protolint --mcp

The commands are:
	lint     lint protocol buffer files
	list     list all current lint rules being used
	version  print protolint version

The flags are:
	--version  print protolint version
	-v         print protolint version (when used as the only argument)
	--mcp      start as an MCP server
`
)

const (
	subCmdLint    = "lint"
	subCmdList    = "list"
	subCmdVersion = "version"
	mcpFlag       = "--mcp"
)

var (
	version  = "master"
	revision = "latest"
)

// Do runs the command logic.
func Do(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	// Check for --version and --mcp flags
	for _, arg := range args {
		if arg == "--version" || (arg == "-v" && len(args) == 1) {
			return doVersion(stdout)
		}
		if arg == mcpFlag {
			return doMCP(stdout, stderr)
		}
	}

	switch {
	case len(args) == 0:
		_, _ = fmt.Fprint(stderr, help)
		return osutil.ExitInternalFailure
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
		return doLint(args[1:], stdout, stderr)
	case subCmdList:
		return doList(args[1:], stdout, stderr)
	case subCmdVersion:
		return doVersion(stdout)
	default:
		return doLint(args, stdout, stderr)
	}
}

func doLint(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	if len(args) < 1 {
		_, _ = fmt.Fprintln(stderr, "protolint lint requires at least one argument. See Usage.")
		_, _ = fmt.Fprint(stderr, help)
		return osutil.ExitInternalFailure
	}

	flags, err := lint.NewFlags(args)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return osutil.ExitInternalFailure
	}
	if len(flags.Args()) < 1 {
		_, _ = fmt.Fprintln(stderr, "protolint lint requires at least one argument. See Usage.")
		_, _ = fmt.Fprint(stderr, help)
		return osutil.ExitInternalFailure
	}

	subCmd, err := lint.NewCmdLint(
		flags,
		stdout,
		stderr,
	)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		if flags.NoErrorOnUnmatchedPattern &&
			(strings.Contains(err.Error(), "not found protocol buffer files") ||
				strings.Contains(err.Error(), "system cannot find the file")) {
			return osutil.ExitSuccess
		}
		return osutil.ExitInternalFailure
	}
	return subCmd.Run()
}

func doList(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	flags, err := list.NewFlags(args)
	if err != nil {
		_, _ = fmt.Fprint(stderr, err)
		return osutil.ExitInternalFailure
	}
	subCmd := list.NewCmdList(
		flags,
		stdout,
		stderr,
	)
	return subCmd.Run()
}

func doVersion(
	stdout io.Writer,
) osutil.ExitCode {
	_, _ = fmt.Fprintln(stdout, "protolint version "+version+"("+revision+")")
	return osutil.ExitSuccess
}

func doMCP(
	stdout io.Writer,
	stderr io.Writer,
) osutil.ExitCode {
	server := mcp.NewServer(stdout, stderr)
	return server.Run()
}
