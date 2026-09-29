package lib

import (
	"io"

	"github.com/yoheimuta/protolint/internal/cmd"
	"github.com/yoheimuta/protolint/internal/file"
	"github.com/yoheimuta/protolint/internal/libinternal"
)

var (
	// ErrLintFailure error is returned when there is a linting error
	ErrLintFailure = libinternal.ErrLintFailure
	// ErrInternalFailure error is returned when there is a parsing, internal, or runtime error.
	ErrInternalFailure = libinternal.ErrInternalFailure
)

// LintRunner is an interface for running lint commands
type LintRunner = libinternal.LintRunner

// SetLintRunner sets the runner used by the Lint function
func SetLintRunner(runner LintRunner) {
	libinternal.SetLintRunner(runner)
}

// GetLintRunner returns the current lint runner
func GetLintRunner() LintRunner {
	return libinternal.GetLintRunner()
}

// Lint is used to lint Protocol Buffer files with the protolint tool.
// It takes an array of strings (args) representing command line arguments,
// as well as two io.Writer instances (stdout and stderr) to which the output of the command should be written.
// It returns an error in the case of a linting error (ErrLintFailure)
// or a parsing, internal, or runtime error (ErrInternalFailure).
// Otherwise, it returns nil on success.
//
// Note: This function automatically initializes the default lint runner if none is set,
// so you don't need to call cmd.Initialize() before using it.
func Lint(args []string, stdout, stderr io.Writer) error {
	file.ResetVFS()
	defer file.ResetVFS()

	// Auto-initialize if needed
	if libinternal.GetLintRunner() == nil {
		cmd.Initialize()
	}

	// Use the internal implementation
	return libinternal.Lint(args, stdout, stderr)
}
