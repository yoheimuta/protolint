package osutil

// ExitCode is a code for os.Exit()
type ExitCode int

// ExitCode constants.
const (
	ExitSuccess         ExitCode = iota
	ExitLintFailure              // Lint errors, exclusively.
	ExitInternalFailure          // All other errors: parsing, internal, runtime errors.
)
