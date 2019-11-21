package lint

import (
	"fmt"

	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/linter/report/reporters"
)

type reporterFlag struct {
	raw      string
	reporter report.Reporter
}

func (f *reporterFlag) String() string {
	return fmt.Sprint(f.raw)
}

func (f *reporterFlag) Set(value string) error {
	if f.reporter != nil {
		return fmt.Errorf("reporter is already set")
	}

	r, err := GetReporter(value)
	if err != nil {
		return err
	}
	f.raw = value
	f.reporter = r
	return nil
}

// GetReporter returns a reporter from the specified key.
func GetReporter(value string) (report.Reporter, error) {
	rs := map[string]report.Reporter{
		"plain": reporters.PlainReporter{},
		"junit": reporters.JUnitReporter{},
		"unix":  reporters.UnixReporter{},
		"json":  reporters.JSONReporter{},
	}
	if r, ok := rs[value]; ok {
		return r, nil
	}
	return nil, fmt.Errorf(`available reporters are "plain", "junit", "json", and "unix"`)
}
