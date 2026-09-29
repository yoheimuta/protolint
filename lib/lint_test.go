package lib_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/yoheimuta/protolint/internal/setting_test"
	"github.com/yoheimuta/protolint/lib"
)

func TestLint(t *testing.T) {
	// Set the mock lint runner for testing
	originalRunner := lib.GetLintRunner() // Save the original runner to restore later
	lib.SetLintRunner(NewMockLintRunner())
	defer func() {
		// Restore the original runner after the test
		lib.SetLintRunner(originalRunner)
	}()

	tests := []struct {
		name            string
		inputArgs       []string
		wantStdoutRegex *regexp.Regexp
		wantStderrRegex *regexp.Regexp
		wantError       error
	}{
		{
			name:            "no args",
			wantStderrRegex: regexp.MustCompile(`[\S\s]*Usage:[\S\s]*protolint <command> \[arguments\][\S\s]*`),
			wantError:       lib.ErrInternalFailure,
		},
		{
			name: "invalid args",
			inputArgs: []string{
				"-config_path",
				setting_test.TestDataPath("lib", "not_exist.yaml"),
				setting_test.TestDataPath("lib", "valid.proto"),
			},
			wantStderrRegex: regexp.MustCompile(`[\S\s]*not_exist.yaml: no such file or directory`),
			wantError:       lib.ErrInternalFailure,
		},
		{
			name: "lint failures",
			inputArgs: []string{
				setting_test.TestDataPath("lib", "invalid.proto"),
			},
			wantStderrRegex: regexp.MustCompile(`[\S\s]*Found an incorrect indentation style[\S\s]*`),
			wantError:       lib.ErrLintFailure,
		},
		{
			name: "lint success",
			inputArgs: []string{
				setting_test.TestDataPath("lib", "valid.proto"),
			},
		},
		{
			name: "lint success by specifying a config file",
			inputArgs: []string{
				"-config_path",
				setting_test.TestDataPath("lib", ".protolint.yaml"),
				setting_test.TestDataPath("lib", "invalid.proto"),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			err := lib.Lint(test.inputArgs, &stdout, &stderr)
			if !errors.Is(err, test.wantError) {
				t.Errorf("got err %v, but want err %v", err, test.wantError)
			}

			if test.wantStdoutRegex != nil {
				if !test.wantStdoutRegex.MatchString(stdout.String()) {
					t.Errorf("got stdout %s, but want to match %v", stdout.String(), test.wantStdoutRegex)
				}
			} else if stdout.Len() > 0 {
				t.Errorf("got stdout %s, but want empty stdout", stdout.String())
			}

			if test.wantStderrRegex != nil {
				if !test.wantStderrRegex.MatchString(stderr.String()) {
					t.Errorf("got stderr %s, but want to match %v", stderr.String(), test.wantStderrRegex)
				}
			} else if stderr.Len() > 0 {
				t.Errorf("got stderr %s, but want empty stderr", stderr.String())
			}
		})
	}
}

func TestLint_StdinAndPhysicalFileCollision(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "protolint_stdin_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change cwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldCwd) }()

	relativeDir := "proto"
	displayPath := filepath.Join(relativeDir, "shared_name.proto")
	if err := os.MkdirAll(relativeDir, 0755); err != nil {
		t.Fatalf("failed to create relative dir: %v", err)
	}

	fsProtoContent := []byte(`syntax = "proto3";
message ValidMessage {
    string text = 1;
}`)
	if err := os.WriteFile(displayPath, fsProtoContent, 0644); err != nil {
		t.Fatalf("failed to write physical proto file: %v", err)
	}

	uniqueMarker := fmt.Sprintf("invalid_syntax_marker_%d", time.Now().Unix())
	stdinProtoContent := fmt.Appendf(nil, `syntax = "%s";`, uniqueMarker)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r1, w1, _ := os.Pipe()
	_, _ = w1.Write(stdinProtoContent)
	_ = w1.Close()
	os.Stdin = r1

	var stdout1, stderr1 bytes.Buffer
	args1 := []string{"lint", "-stdin_filename", displayPath, "-"}

	// Run #1 via stdin
	_ = lib.Lint(args1, &stdout1, &stderr1)

	if !strings.Contains(stderr1.String(), uniqueMarker) {
		t.Fatalf("First run (stdin) failed to report syntax error. Stderr: %s", stderr1.String())
	}

	var stdout2, stderr2 bytes.Buffer
	args2 := []string{"lint", displayPath}

	// Run #2 via the physical file with the same relative path
	err = lib.Lint(args2, &stdout2, &stderr2)

	if err != nil && err.Error() == lib.ErrInternalFailure.Error() {
		t.Fatalf("Collision detected! Second run read dirty stdin cache instead of physical file. Stderr: %s", stderr2.String())
	}

	if strings.Contains(stderr2.String(), uniqueMarker) {
		t.Fatal("Collision detected! The physical file run encountered a leaked stdin syntax error.")
	}
}
