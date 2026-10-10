package lint_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds/lint"
	"github.com/yoheimuta/protolint/internal/osutil"
)

func TestCmdLint_Run_FixMode(t *testing.T) {
	tests := []struct {
		name          string
		args          func(filePath string) []string
		protoContent  string
		wantExitCode  osutil.ExitCode
		wantRemaining string
		checkFile     func(t *testing.T, filePath string)
	}{
		{
			name: "exit 0 when -fix is enabled and all issues are fixable",
			args: func(filePath string) []string {
				return []string{"-fix", filePath}
			},
			protoContent: `syntax = "proto3";

message SearchRequest {
    string query = 1;
}
`,
			wantExitCode: osutil.ExitSuccess,
			checkFile: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if strings.Contains(string(content), "    string query = 1;") {
					t.Errorf("expected 4-space indent to be fixed, but still found it: %s", string(content))
				}
				if !strings.Contains(string(content), "  string query = 1;") {
					t.Errorf("expected 2-space indent after fix, got: %s", string(content))
				}
			},
		},
		{
			name: "exit 1 when -fix is enabled but unfixed issues remain",
			args: func(filePath string) []string {
				return []string{"-fix", filePath}
			},
			protoContent: `syntax = "proto3";

// This is an extremely long comment line that clearly exceeds the eighty characters limit for max line length.
message SearchRequest {
  string query = 1;
}
`,
			wantExitCode:  osutil.ExitLintFailure,
			wantRemaining: "The line length is",
		},
		{
			name: "exit 1 when -fix is enabled with both fixable and unfixable issues",
			args: func(filePath string) []string {
				return []string{"-fix", filePath}
			},
			protoContent: `syntax = "proto3";

// This is an extremely long comment line that clearly exceeds the eighty characters limit for max line length.
message SearchRequest {
    string query = 1;
}
`,
			wantExitCode:  osutil.ExitLintFailure,
			wantRemaining: "The line length is",
			checkFile: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if !strings.Contains(string(content), "  string query = 1;") {
					t.Errorf("expected fixable indent to be fixed, got: %s", string(content))
				}
			},
		},
		{
			name: "exit 1 when -fix is not enabled and fixable issues exist",
			args: func(filePath string) []string {
				return []string{filePath}
			},
			protoContent: `syntax = "proto3";

message SearchRequest {
    string query = 1;
}
`,
			wantExitCode:  osutil.ExitLintFailure,
			wantRemaining: "Found an incorrect indentation style",
			checkFile: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if !strings.Contains(string(content), "    string query = 1;") {
					t.Errorf("expected file to remain unchanged, got: %s", string(content))
				}
			},
		},
		{
			name: "exit 0 when -fix is enabled and file has no lint issues",
			args: func(filePath string) []string {
				return []string{"-fix", filePath}
			},
			protoContent: `syntax = "proto3";

message SearchRequest {
  string query = 1;
}
`,
			wantExitCode: osutil.ExitSuccess,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			filePath := filepath.Join(dir, "test.proto")
			if err := os.WriteFile(filePath, []byte(tc.protoContent), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			flags, err := lint.NewFlags(tc.args(filePath))
			if err != nil {
				t.Fatalf("failed to parse flags: %v", err)
			}

			var stdout, stderr bytes.Buffer
			cmdLint, err := lint.NewCmdLint(flags, &stdout, &stderr)
			if err != nil {
				t.Fatalf("failed to create CmdLint: %v", err)
			}

			exitCode := cmdLint.Run()
			if exitCode != tc.wantExitCode {
				t.Errorf("got exit code %v, want %v (stderr: %s)", exitCode, tc.wantExitCode, stderr.String())
			}

			if tc.wantRemaining != "" && !strings.Contains(stderr.String(), tc.wantRemaining) {
				t.Errorf("expected stderr to contain %q, but got %q", tc.wantRemaining, stderr.String())
			}

			if tc.checkFile != nil {
				tc.checkFile(t, filePath)
			}
		})
	}
}

func TestCmdLint_Run_FixMode_MultipleFiles(t *testing.T) {
	t.Run("exit 0 when all files contain only fixable issues", func(t *testing.T) {
		dir := t.TempDir()
		file1 := filepath.Join(dir, "file1.proto")
		file2 := filepath.Join(dir, "file2.proto")

		content1 := `syntax = "proto3";

message Req1 {
    string query = 1;
}
`
		content2 := `syntax = "proto3";

message Req2 {
    string filter = 1;
}
`
		if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
			t.Fatalf("failed to write file1: %v", err)
		}
		if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
			t.Fatalf("failed to write file2: %v", err)
		}

		flags, err := lint.NewFlags([]string{"-fix", file1, file2})
		if err != nil {
			t.Fatalf("failed to parse flags: %v", err)
		}

		var stdout, stderr bytes.Buffer
		cmdLint, err := lint.NewCmdLint(flags, &stdout, &stderr)
		if err != nil {
			t.Fatalf("failed to create CmdLint: %v", err)
		}

		exitCode := cmdLint.Run()
		if exitCode != osutil.ExitSuccess {
			t.Errorf("got exit code %v, want %v (stderr: %s)", exitCode, osutil.ExitSuccess, stderr.String())
		}
	})

	t.Run("exit 1 when one of multiple files has unfixed issues", func(t *testing.T) {
		dir := t.TempDir()
		file1 := filepath.Join(dir, "file1.proto")
		file2 := filepath.Join(dir, "file2.proto")

		content1 := `syntax = "proto3";

message Req1 {
    string query = 1;
}
`
		content2 := `syntax = "proto3";

// This is an extremely long comment line that clearly exceeds the eighty characters limit for max line length.
message Req2 {
  string filter = 1;
}
`
		if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
			t.Fatalf("failed to write file1: %v", err)
		}
		if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
			t.Fatalf("failed to write file2: %v", err)
		}

		flags, err := lint.NewFlags([]string{"-fix", file1, file2})
		if err != nil {
			t.Fatalf("failed to parse flags: %v", err)
		}

		var stdout, stderr bytes.Buffer
		cmdLint, err := lint.NewCmdLint(flags, &stdout, &stderr)
		if err != nil {
			t.Fatalf("failed to create CmdLint: %v", err)
		}

		exitCode := cmdLint.Run()
		if exitCode != osutil.ExitLintFailure {
			t.Errorf("got exit code %v, want %v (stderr: %s)", exitCode, osutil.ExitLintFailure, stderr.String())
		}
		if !strings.Contains(stderr.String(), "The line length is") {
			t.Errorf("expected stderr to contain 'The line length is', got %q", stderr.String())
		}
	})
}
