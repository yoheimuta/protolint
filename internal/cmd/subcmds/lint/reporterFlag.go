package lint

import (
	"fmt"
	"strings"

	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/linter/report/reporters"
)

type reporterFlag struct {
	raw      string
	reporter report.Reporter
}

type reporterStreamFlag struct {
	reporterFlag
	targetFile string
}

type reporterStreamFlags []reporterStreamFlag

func (f *reporterStreamFlag) String() string {
	return fmt.Sprint(f.raw)
}

func (f *reporterStreamFlag) Set(value string) error {
	if f.reporter != nil {
		return fmt.Errorf("reporter is already set")
	}

	valueSplit := strings.SplitN(value, ":", 2)
	if len(valueSplit) != 2 {
		return fmt.Errorf("cannot find output file in %s", value)
	}

	reporterName := valueSplit[0]
	outputFile := valueSplit[1]

	r, err := GetReporter(reporterName)
	if err != nil {
		return err
	}

	f.raw = value
	f.reporter = r
	f.targetFile = outputFile

	return nil
}

func (fs *reporterStreamFlags) String() string {
	var items []string
	for _, flag := range *fs {
		items = append(items, flag.String())
	}

	return strings.Join(items, " ")
}

func (fs *reporterStreamFlags) Set(value string) error {
	var r reporterStreamFlag
	err := r.Set(value)
	if err != nil {
		return err
	}

	*fs = append(*fs, r)
	return nil
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
		"sarif": reporters.SarifReporter{},
	}
	if r, ok := rs[value]; ok {
		return r, nil
	}
	return nil, fmt.Errorf(`available reporters are "plain", "junit", "json", "sarif", and "unix"`)
}
