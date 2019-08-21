package reporters

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/yoheimuta/go-protoparser/parser/meta"

	"github.com/yoheimuta/protolint/linter/report"
)

const (
	packageName = "net.protolint"
)

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many testcases.
type JUnitTestSuite struct {
	XMLName   xml.Name `xml:"testsuite"`
	Package   string   `xml:"package"`
	Tests     int      `xml:"tests,attr"`
	Failures  int      `xml:"failures,attr"`
	Time      string   `xml:"time,attr"`
	TestCases []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	ClassName string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",cdata"`
}

func constructTestCaseName(ruleID string) string {
	return packageName + "." + ruleID
}

func constructContents(pos meta.Position) string {
	return fmt.Sprintf("line %d, col %d", pos.Line, pos.Column)
}

// JUnitReporter prints failures in JUnit XML format.
type JUnitReporter struct{}

// Report writes failures to w.
func (r JUnitReporter) Report(w io.Writer, fs []report.Failure) error {
	suites := &JUnitTestSuites{}
	if 0 < len(fs) {
		var testcases []JUnitTestCase
		for _, f := range fs {
			testcase := JUnitTestCase{
				Name:      constructTestCaseName(f.RuleID()),
				ClassName: f.FilenameWithoutExt(),
				Time:      "0",
				Failure: &JUnitFailure{
					Message:  f.Message(),
					Type:     "error",
					Contents: constructContents(f.Pos()),
				},
			}
			testcases = append(testcases, testcase)
		}

		suite := JUnitTestSuite{
			Package:   packageName,
			Tests:     len(fs),
			Failures:  len(fs),
			Time:      "0",
			TestCases: testcases,
		}
		suites.Suites = append(suites.Suites, suite)
	} else {
		suites.Suites = []JUnitTestSuite{
			{
				Package: packageName,
				Tests:   1,
				Time:    "0",
				TestCases: []JUnitTestCase{
					{
						ClassName: constructTestCaseName("ALL_RULES"),
						Name:      "All Rules",
						Time:      "0",
					},
				},
			},
		}
	}

	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	err = enc.Encode(suites)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}
