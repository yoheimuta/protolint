package reporters

import (
	"fmt"
	"io"
	"strings"

	"github.com/yoheimuta/protolint/linter/report"
)

// TscRport prints failures as string compatible to Type script compiler
//
// The format is "FILENAME(LINE,COL): SEVERITY RULE_ID: MESSAGE".
type TscReporter struct{}

func getTscSeverity(s string) string {
	if s == "note" {
		return "info"
	}

	return s
}

// Report writes failures to w.
func (r TscReporter) Report(w io.Writer, fs []report.Failure) error {
	for _, failure := range fs {
		tsc_output := fmt.Sprintf(
			"%s(%d,%d): %s %s: '%s'",
			failure.Pos().Filename,
			failure.Pos().Line,
			failure.Pos().Column,
			getTscSeverity(failure.Severity()),
			failure.RuleID(),
			strings.Trim(failure.Message(), `"`),
		)
		_, err := fmt.Fprintln(w, tsc_output)
		if err != nil {
			return err
		}
	}
	return nil
}
