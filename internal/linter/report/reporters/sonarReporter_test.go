package reporters_test

import (
	"bytes"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/linter/report/reporters"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

func TestSonarReporter_Report(t *testing.T) {
	tests := []struct {
		name          string
		inputFailures []report.Failure
		wantOutput    string
	}{
		{
			name: "Prints failures in Sonar format",
			inputFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
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
					string(rule.SeverityError),
					`EnumField name "SECOND.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
			},
			wantOutput: `[
  {
    "engineId": "protolint",
    "ruleId": "ENUM_NAMES_UPPER_CAMEL_CASE",
    "primaryLocation": {
      "message": "EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES",
      "filePath": "example.proto",
      "textRange": {
        "startLine": 5,
        "startColumn": 10
      }
    },
    "severity": "MAJOR",
    "issueType": "CODE_SMELL"
  },
  {
    "engineId": "protolint",
    "ruleId": "ENUM_NAMES_UPPER_CAMEL_CASE",
    "primaryLocation": {
      "message": "EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES",
      "filePath": "example.proto",
      "textRange": {
        "startLine": 10,
        "startColumn": 20
      }
    },
    "severity": "MAJOR",
    "issueType": "CODE_SMELL"
  }
]`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := reporters.SonarReporter{}.Report(buf, test.inputFailures)
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
