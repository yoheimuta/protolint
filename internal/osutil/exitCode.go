package osutil

// ExitCode is a code for os.Exit()
type ExitCode int

// ExitCode constants.
const (
	ExitSuccess ExitCode = iota
	ExitFailure
)
