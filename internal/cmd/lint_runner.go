package cmd

import (
	"io"

	"github.com/yoheimuta/protolint/internal/libinternal"
	"github.com/yoheimuta/protolint/internal/osutil"
)

// CmdLintRunner implements the LintRunner interface for cmd package
type CmdLintRunner struct{}

// NewCmdLintRunner creates a new CmdLintRunner
func NewCmdLintRunner() *CmdLintRunner {
	return &CmdLintRunner{}
}

// Run executes the lint command
func (r *CmdLintRunner) Run(args []string, stdout, stderr io.Writer) osutil.ExitCode {
	return Do(args, stdout, stderr)
}

// Initialize registers the cmd lint runner with the internal library
func Initialize() {
	libinternal.SetLintRunner(NewCmdLintRunner())
}
