package mcp

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestLintFilesTool_GetInfo(t *testing.T) {
	tool := NewLintFilesTool()
	info := tool.GetInfo()

	// Check basic properties
	if info.Name != "lint-files" {
		t.Errorf("Expected tool name 'lint-files', got '%s'", info.Name)
	}

	if info.Description == "" {
		t.Error("Expected non-empty description")
	}

	// Check schema
	schema, ok := info.InputSchema.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected schema to be map[string]interface{}")
	}

	// Check schema type
	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	// Check properties
	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected properties to be map[string]interface{}")
	}

	// Check required files property
	filesProperty, ok := properties["files"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected files property to be map[string]interface{}")
	}

	if filesProperty["type"] != "array" {
		t.Errorf("Expected files property type 'array', got '%v'", filesProperty["type"])
	}

	// Check required fields
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatalf("Expected required to be []string")
	}

	if len(required) == 0 || required[0] != "files" {
		t.Errorf("Expected 'files' to be required, got %v", required)
	}
}

func TestLintFilesTool_Execute_InvalidArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			name:    "invalid JSON",
			args:    `{"files": `,
			wantErr: true,
		},
		{
			name:    "empty files array",
			args:    `{"files": []}`,
			wantErr: true,
		},
		{
			name:    "missing files field",
			args:    `{"config_path": "/path/to/config.yaml"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := NewLintFilesTool()
			_, err := tool.Execute(json.RawMessage(tt.args))

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test that the arguments are properly parsed
func TestLintFilesArgs_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    LintFilesArgs
		wantErr bool
	}{
		{
			name: "basic args",
			args: `{"files": ["/path/to/file.proto"]}`,
			want: LintFilesArgs{
				Files: []string{"/path/to/file.proto"},
			},
			wantErr: false,
		},
		{
			name: "with config path",
			args: `{"files": ["/path/to/file.proto"], "config_path": "/path/to/config.yaml"}`,
			want: LintFilesArgs{
				Files:      []string{"/path/to/file.proto"},
				ConfigPath: "/path/to/config.yaml",
			},
			wantErr: false,
		},
		{
			name: "with fix",
			args: `{"files": ["/path/to/file.proto"], "fix": true}`,
			want: LintFilesArgs{
				Files: []string{"/path/to/file.proto"},
				Fix:   true,
			},
			wantErr: false,
		},
		{
			name: "with all options",
			args: `{"files": ["/path/to/file1.proto", "/path/to/file2.proto"], "config_path": "/path/to/config.yaml", "fix": true}`,
			want: LintFilesArgs{
				Files:      []string{"/path/to/file1.proto", "/path/to/file2.proto"},
				ConfigPath: "/path/to/config.yaml",
				Fix:        true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got LintFilesArgs
			err := json.Unmarshal([]byte(tt.args), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
