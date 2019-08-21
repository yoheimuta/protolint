package report

import (
	"io"

	"github.com/yoheimuta/protolint/linter/report"
)

// Reporter is responsible to output results in the specific format.
type Reporter interface {
	Report(io.Writer, []report.Failure) error
}
