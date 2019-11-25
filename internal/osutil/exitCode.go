package osutil

// ExitCode is a code for os.Exit()
type ExitCode int

// ExitCode constants.
const (
	ExitSuccess ExitCode = iota
	ExitInternalFailure
	ExitRuntimeFailure // runtime exceptions in go throw 2 by default
	ExitLintFailure
	ExitParseFailure
)
