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

	rs := map[string]report.Reporter{
		"plain": reporters.PlainReporter{},
	}
	if r, ok := rs[value]; ok {
		f.raw = value
		f.reporter = r
		return nil
	}
	return fmt.Errorf(`available reporters are "plain"`)
}
