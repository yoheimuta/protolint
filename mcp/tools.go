// Package mcp implements the Model Context Protocol (MCP) server for protolint.
package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/yoheimuta/protolint/internal/libinternal"
)

// Tool defines the interface for MCP tools
type Tool interface {
	GetInfo() ToolInfo
	Execute(args json.RawMessage) (any, error)
}

// LintFilesTool is a tool for linting Proto files
type LintFilesTool struct{}

// NewLintFilesTool creates a new LintFilesTool
func NewLintFilesTool() *LintFilesTool {
	return &LintFilesTool{}
}

// GetInfo returns the tool information
func (t *LintFilesTool) GetInfo() ToolInfo {
	return ToolInfo{
		Name:        "lint-files",
		Description: "Lint and fix Protocol Buffer files using protolint",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"files": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
					},
					"description": "List of file paths to lint. The paths must be absolute.",
				},
				"config_path": map[string]any{
					"type":        "string",
					"description": "Path to protolint config file",
				},
				"fix": map[string]any{
					"type":        "boolean",
					"description": "Fix lint errors if possible. Default is false. It may return failures even if all errors are fixed, so you are strongly recommended to lint the target files again to see if they are fixed.",
				},
			},
			"required": []string{"files"},
		},
	}
}

// LintFilesArgs represents arguments for lint-files tool
type LintFilesArgs struct {
	Files      []string `json:"files"`
	ConfigPath string   `json:"config_path,omitempty"`
	Fix        bool     `json:"fix,omitempty"`
}

// Execute runs the lint-files tool
func (t *LintFilesTool) Execute(args json.RawMessage) (any, error) {
	var lintArgs LintFilesArgs
	if err := json.Unmarshal(args, &lintArgs); err != nil {
		return nil, fmt.Errorf("invalid arguments: %v", err)
	}

	if len(lintArgs.Files) == 0 {
		return nil, fmt.Errorf("no files specified")
	}

	// Construct command line arguments
	cmdArgs := []string{}

	// Use MCP reporter
	cmdArgs = append(cmdArgs, "--reporter", "mcp")

	if lintArgs.ConfigPath != "" {
		cmdArgs = append(cmdArgs, "--config_path", lintArgs.ConfigPath)
	}

	if lintArgs.Fix {
		cmdArgs = append(cmdArgs, "--fix")
	}

	// Add files at the end
	cmdArgs = append(cmdArgs, lintArgs.Files...)

	// Capture output to parse it
	var outputBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	// Run lint command
	err := libinternal.Lint(cmdArgs, &outputBuffer, &errorBuffer)

	// Determine exit code based on error
	exitCode := 0
	if err != nil {
		if err == libinternal.ErrLintFailure {
			exitCode = 1
		} else {
			exitCode = 2
			// Return error information if internal error occurred
			return map[string]any{
				"exit_code": exitCode,
				"error":     err.Error(),
				"stderr":    errorBuffer.String(),
			}, nil
		}
	}

	// Parse the JSON output from MCP reporter
	var result map[string]any
	if err := json.Unmarshal(errorBuffer.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse lint output: %v\n%s", err, errorBuffer.String())
	}

	// Add exit code to result
	result["exit_code"] = exitCode

	if lintArgs.Fix {
		// If fix is enabled, add a message indicating that you should lint again
		result["message"] = "Don't think these failures remain now. The fixer could fix all the failures already. Before you manually edit the file, you must lint the files again to see if they are fixed."
	}

	return result, nil
}
