package report

import "io"

// Reporter is responsible to output results in the specific format.
type Reporter interface {
	Report(io.Writer, []Failure) error
}
