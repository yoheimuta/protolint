package report

import (
	"io"
	"os"

	"github.com/yoheimuta/protolint/linter/report"
)

const (
	WriteToConsole string = "-"
)

// Reporter is responsible to output results in the specific format.
type Reporter interface {
	Report(io.Writer, []report.Failure) error
}

type ReporterWithOutput struct {
	reporter   Reporter
	targetFile string
}

type ReportersWithOutput []ReporterWithOutput

func (ro ReporterWithOutput) ReportWithFallback(w io.Writer, failures []report.Failure) error {
	if ro.targetFile != WriteToConsole {
		var err error
		w, err = os.OpenFile(ro.targetFile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
	}

	return ro.reporter.Report(w, failures)
}

func (ros ReportersWithOutput) ReportWithFallback(w io.Writer, failures []report.Failure) error {
	for _, ro := range ros {
		err := ro.ReportWithFallback(w, failures)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewReporterWithOutput(r Reporter, targetFile string) *ReporterWithOutput {
	return &ReporterWithOutput{r, targetFile}
}
