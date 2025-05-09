package reporters

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yoheimuta/protolint/linter/report"
)

// MCPReporter prints failures in MCP friendly JSON format.
type MCPReporter struct{}

// Report writes failures to w in MCP friendly format.
func (r MCPReporter) Report(w io.Writer, fs []report.Failure) error {
	// Group failures by file
	fileFailures := make(map[string][]map[string]interface{})

	for _, failure := range fs {
		filePath := failure.Pos().Filename

		failureInfo := map[string]interface{}{
			"rule_id":  failure.RuleID(),
			"message":  failure.Message(),
			"line":     failure.Pos().Line,
			"column":   failure.Pos().Column,
			"severity": failure.Severity(),
		}

		fileFailures[filePath] = append(fileFailures[filePath], failureInfo)
	}

	// Convert to array of results
	var fileResults []map[string]interface{} = []map[string]interface{}{}
	for filePath, failures := range fileFailures {
		fileResults = append(fileResults, map[string]interface{}{
			"file_path": filePath,
			"failures":  failures,
		})
	}

	result := map[string]interface{}{
		"results": fileResults,
	}

	bs, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(bs))
	return err
}
