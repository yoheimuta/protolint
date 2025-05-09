package reporters

import (
	"bytes"
	"encoding/json"
	"testing"

	parser_meta "github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestMCPReporter_Report(t *testing.T) {
	tests := []struct {
		name     string
		failures []report.Failure
		want     map[string]interface{}
	}{
		{
			name:     "empty failures",
			failures: []report.Failure{},
			want: map[string]interface{}{
				"results": []interface{}{},
			},
		},
		{
			name: "single failure",
			failures: []report.Failure{
				report.Failuref(
					parser_meta.Position{
						Filename: "test.proto",
						Line:     10,
						Column:   15,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					"error",
					"Enum name must be UpperCamelCase",
				),
			},
			want: map[string]interface{}{
				"results": []interface{}{
					map[string]interface{}{
						"file_path": "test.proto",
						"failures": []interface{}{
							map[string]interface{}{
								"rule_id":  "ENUM_NAMES_UPPER_CAMEL_CASE",
								"message":  "Enum name must be UpperCamelCase",
								"line":     float64(10),
								"column":   float64(15),
								"severity": "error",
							},
						},
					},
				},
			},
		},
		{
			name: "multiple failures from the same file",
			failures: []report.Failure{
				report.Failuref(
					parser_meta.Position{
						Filename: "test.proto",
						Line:     10,
						Column:   15,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					"error",
					"Enum name must be UpperCamelCase",
				),
				report.Failuref(
					parser_meta.Position{
						Filename: "test.proto",
						Line:     20,
						Column:   5,
					},
					"INDENT",
					"warning",
					"Found an incorrect indentation style",
				),
			},
			want: map[string]interface{}{
				"results": []interface{}{
					map[string]interface{}{
						"file_path": "test.proto",
						"failures": []interface{}{
							map[string]interface{}{
								"rule_id":  "ENUM_NAMES_UPPER_CAMEL_CASE",
								"message":  "Enum name must be UpperCamelCase",
								"line":     float64(10),
								"column":   float64(15),
								"severity": "error",
							},
							map[string]interface{}{
								"rule_id":  "INDENT",
								"message":  "Found an incorrect indentation style",
								"line":     float64(20),
								"column":   float64(5),
								"severity": "warning",
							},
						},
					},
				},
			},
		},
		{
			name: "failures from multiple files",
			failures: []report.Failure{
				report.Failuref(
					parser_meta.Position{
						Filename: "test1.proto",
						Line:     10,
						Column:   15,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					"error",
					"Enum name must be UpperCamelCase",
				),
				report.Failuref(
					parser_meta.Position{
						Filename: "test2.proto",
						Line:     5,
						Column:   3,
					},
					"INDENT",
					"warning",
					"Found an incorrect indentation style",
				),
			},
			want: map[string]interface{}{
				"results": []interface{}{
					map[string]interface{}{
						"file_path": "test1.proto",
						"failures": []interface{}{
							map[string]interface{}{
								"rule_id":  "ENUM_NAMES_UPPER_CAMEL_CASE",
								"message":  "Enum name must be UpperCamelCase",
								"line":     float64(10),
								"column":   float64(15),
								"severity": "error",
							},
						},
					},
					map[string]interface{}{
						"file_path": "test2.proto",
						"failures": []interface{}{
							map[string]interface{}{
								"rule_id":  "INDENT",
								"message":  "Found an incorrect indentation style",
								"line":     float64(5),
								"column":   float64(3),
								"severity": "warning",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			r := MCPReporter{}
			err := r.Report(&buf, tt.failures)

			if err != nil {
				t.Errorf("MCPReporter.Report() error = %v", err)
				return
			}

			// Parse JSON output
			var got map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Errorf("Failed to parse JSON output: %v", err)
				return
			}

			// Compare actual vs expected
			assertEquivalentJSON(t, got, tt.want)
		})
	}
}

// assertEquivalentJSON compares two JSON-compatible objects for structural equality
func assertEquivalentJSON(t *testing.T, got, want map[string]interface{}) {
	t.Helper()

	// Compare the "results" array length
	gotResults, gotOk := got["results"].([]interface{})
	wantResults, wantOk := want["results"].([]interface{})

	if !gotOk || !wantOk {
		t.Errorf("Expected 'results' to be an array. Got %T, want %T", got["results"], want["results"])
		return
	}

	if len(gotResults) != len(wantResults) {
		t.Errorf("Results length mismatch. Got %d, want %d", len(gotResults), len(wantResults))
		return
	}

	// If there are no results, we're done
	if len(gotResults) == 0 {
		return
	}

	// For each file in wantResults, find matching file in gotResults
	for _, wantFile := range wantResults {
		wantFileMap := wantFile.(map[string]interface{})
		wantFilePath := wantFileMap["file_path"].(string)

		// Find matching file in gotResults
		var matchingGotFile map[string]interface{}
		for _, gotFile := range gotResults {
			gotFileMap := gotFile.(map[string]interface{})
			gotFilePath := gotFileMap["file_path"].(string)

			if gotFilePath == wantFilePath {
				matchingGotFile = gotFileMap
				break
			}
		}

		if matchingGotFile == nil {
			t.Errorf("Missing file %s in results", wantFilePath)
			continue
		}

		// Compare failures
		wantFailures := wantFileMap["failures"].([]interface{})
		gotFailures := matchingGotFile["failures"].([]interface{})

		if len(gotFailures) != len(wantFailures) {
			t.Errorf("Failures length mismatch for file %s. Got %d, want %d",
				wantFilePath, len(gotFailures), len(wantFailures))
			continue
		}

		// Compare each failure
		for i, wantFailure := range wantFailures {
			wantFailureMap := wantFailure.(map[string]interface{})
			gotFailureMap := gotFailures[i].(map[string]interface{})

			for key, wantValue := range wantFailureMap {
				gotValue, exists := gotFailureMap[key]
				if !exists {
					t.Errorf("Missing key %s in failure for file %s", key, wantFilePath)
					continue
				}

				if gotValue != wantValue {
					t.Errorf("Value mismatch for key %s in failure for file %s. Got %v, want %v",
						key, wantFilePath, gotValue, wantValue)
				}
			}
		}
	}
}
