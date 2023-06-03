package reporters_test

import (
	"bytes"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/linter/report/reporters"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

func TestSarifReporter_Report(t *testing.T) {
	tests := []struct {
		name          string
		inputFailures []report.Failure
		wantOutput    string
	}{
		{
			name: "Prints failures in JSON format",
			inputFailures: []report.Failure{
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.Severity_Error),
					`EnumField name "fIRST_VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.Severity_Warning),
					`EnumField name "SECOND.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
			},
			wantOutput: `{
  "runs": [
    {
      "artifacts": [
        {
          "location": {
            "uri": "example.proto"
          }
        }
      ],
      "results": [
        {
          "kind": "fail",
          "level": "error",
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "example.proto"
                },
                "region": {
                  "startColumn": 10,
                  "startLine": 5
                }
              }
            }
          ],
          "message": {
            "text": "EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES"
          },
          "ruleId": "ENUM_NAMES_UPPER_CAMEL_CASE"
        },
        {
          "kind": "fail",
          "level": "warning",
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "example.proto"
                },
                "region": {
                  "startColumn": 20,
                  "startLine": 10
                }
              }
            }
          ],
          "message": {
            "text": "EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES"
          },
          "ruleId": "ENUM_NAMES_UPPER_CAMEL_CASE"
        }
      ],
      "tool": {
        "driver": {
          "informationUri": "https://github.com/yoheimuta/protolint",
          "name": "protolint",
          "rules": [
            {
              "helpUri": "https://github.com/yoheimuta/protolint",
              "id": "ENUM_NAMES_UPPER_CAMEL_CASE"
            }
          ]
        }
      }
    }
  ],
  "version": "2.1.0"
}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := reporters.SarifReporter{}.Report(buf, test.inputFailures)
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
