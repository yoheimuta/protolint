package lint

import (
	"io"
	"testing"
)

func TestCmdLint_Run_FailFastOnStdinWithModification(t *testing.T) {
	fixMode, _ := NewFlags([]string{"-fix", "-stdin_filename=file.proto", "-"})
	autoDisable, _ := NewFlags([]string{"-auto_disable=next", "-stdin_filename=file.proto", "-"})
	correct, _ := NewFlags([]string{"-stdin_filename=file.proto", "-"})

	tests := []struct {
		name        string
		flags       Flags
		expectError bool
	}{
		{
			name:        "stdin with fix mode should fail",
			flags:       fixMode,
			expectError: true,
		},
		{
			name:        "stdin with auto-disable should fail",
			flags:       autoDisable,
			expectError: true,
		},
		{
			name:        "stdin without modification should pass",
			flags:       correct,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := NewCmdLint(tt.flags, io.Discard, io.Discard)

			_, err := c.run()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("got unexpected error")
				}
			}
		})
	}
}
