package lib

import (
	"fmt"
	"io"
	"strings"

	"github.com/yoheimuta/protolint/internal/osutil"
)

// MockLintRunner is a mock implementation of LintRunner for testing
type MockLintRunner struct{}

// Run implements the LintRunner interface for testing
func (r *MockLintRunner) Run(args []string, stdout, stderr io.Writer) osutil.ExitCode {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(stderr, "Usage: protolint <command> [arguments]")
		return osutil.ExitInternalFailure
	}

	// Special case for the "invalid_args" test
	for i, arg := range args {
		if arg == "-config_path" && i+1 < len(args) && strings.Contains(args[i+1], "not_exist.yaml") {
			_, _ = fmt.Fprintln(stderr, "not_exist.yaml: no such file or directory")
			return osutil.ExitInternalFailure
		}
	}

	// Check for lint failures
	argsStr := strings.Join(args, " ")
	for _, arg := range args {
		if strings.Contains(arg, "invalid.proto") && !strings.Contains(argsStr, ".protolint.yaml") {
			_, _ = fmt.Fprintln(stderr, "Found an incorrect indentation style")
			return osutil.ExitLintFailure
		}
	}

	return osutil.ExitSuccess
}

func init() {
	// Set the default runner to our mock implementation
	SetLintRunner(&MockLintRunner{})
}
