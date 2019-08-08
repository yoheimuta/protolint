package reporters

import (
	"fmt"
	"io"

	"github.com/yoheimuta/protolint/internal/linter/report"
)

// PlainReporter prints failures as it is.
type PlainReporter struct{}

// Report writes failures to w.
func (r PlainReporter) Report(w io.Writer, fs []report.Failure) error {
	for _, failure := range fs {
		_, err := fmt.Fprintln(w, failure)
		if err != nil {
			return err
		}
	}
	return nil
}
