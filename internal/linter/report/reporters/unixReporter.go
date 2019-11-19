package reporters

import (
	"fmt"
	"io"

	"github.com/yoheimuta/protolint/linter/report"
)

// UnixReporter prints failures as it respects Unix output conventions
// those are frequently employed by preprocessors and compilers.
//
// The format is "FILENAME:LINE:COL: MESSAGE".
type UnixReporter struct{}

// Report writes failures to w.
func (r UnixReporter) Report(w io.Writer, fs []report.Failure) error {
	for _, failure := range fs {
		unix := fmt.Sprintf(
			"%s: %s",
			failure.Pos(),
			failure.Message(),
		)
		_, err := fmt.Fprintln(w, unix)
		if err != nil {
			return err
		}
	}
	return nil
}
