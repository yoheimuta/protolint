package lib_test

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

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
