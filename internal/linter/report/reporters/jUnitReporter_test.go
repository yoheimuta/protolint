package reporters_test

import (
	"bytes"
	"testing"

	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/linter/report/reporters"
)

func TestJUnitReporter_Report(t *testing.T) {
	tests := []struct {
		name          string
		inputFailures []report.Failure
		wantOutput    string
	}{
		{
			name: "Prints no failures in the JUnit XML format",
			wantOutput: `<?xml version="1.0" encoding="UTF-8"?>
  <testsuites>
      <testsuite tests="1" failures="0" time="0">
          <package>net.protolint</package>
          <testcase classname="net.protolint.ALL_RULES" name="All Rules" time="0"></testcase>
      </testsuite>
  </testsuites>
`,
		},
		{
			name: "Prints failures in the JUnit XML format",
			inputFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					`EnumField name "fIRST_VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					`EnumField name "SECOND.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
			},
			wantOutput: `<?xml version="1.0" encoding="UTF-8"?>
  <testsuites>
      <testsuite tests="2" failures="2" time="0">
          <package>net.protolint</package>
          <testcase classname="example" name="net.protolint.ENUM_NAMES_UPPER_CAMEL_CASE" time="0">
              <failure message="EnumField name &#34;fIRST_VALUE&#34; must be CAPITALS_WITH_UNDERSCORES" type="error"><![CDATA[line 5, col 10]]></failure>
          </testcase>
          <testcase classname="example" name="net.protolint.ENUM_NAMES_UPPER_CAMEL_CASE" time="0">
              <failure message="EnumField name &#34;SECOND.VALUE&#34; must be CAPITALS_WITH_UNDERSCORES" type="error"><![CDATA[line 10, col 20]]></failure>
          </testcase>
      </testsuite>
  </testsuites>
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := reporters.JUnitReporter{}.Report(buf, test.inputFailures)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if buf.String() != test.wantOutput {
				t.Errorf("got %s, but want %s", buf.String(), test.wantOutput)
			}
		})
	}
}
