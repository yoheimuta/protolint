package lib

import (
	"errors"
	"io"

	"github.com/yoheimuta/protolint/internal/osutil"
)

var (
	// ErrLintFailure error is returned when there is a linting error
	ErrLintFailure = errors.New("lint error")
	// ErrInternalFailure error is returned when there is a parsing, internal, or runtime error.
	ErrInternalFailure = errors.New("parsing, internal or runtime errors")
)

// LintRunner is an interface for running lint commands
type LintRunner interface {
	Run(args []string, stdout, stderr io.Writer) osutil.ExitCode
}

var defaultRunner LintRunner

// SetLintRunner sets the runner used by the Lint function
func SetLintRunner(runner LintRunner) {
	defaultRunner = runner
}

// Lint is used to lint Protocol Buffer files with the protolint tool.
// It takes an array of strings (args) representing command line arguments,
// as well as two io.Writer instances (stdout and stderr) to which the output of the command should be written.
// It returns an error in the case of a linting error (ErrLintFailure)
// or a parsing, internal, or runtime error (ErrInternalFailure).
// Otherwise, it returns nil on success.
func Lint(args []string, stdout, stderr io.Writer) error {
	if defaultRunner == nil {
		return ErrInternalFailure
	}

	switch defaultRunner.Run(args, stdout, stderr) {
	case osutil.ExitSuccess:
		return nil

	case osutil.ExitLintFailure:
		return ErrLintFailure

	default:
		return ErrInternalFailure
	}
}
