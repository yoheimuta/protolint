package fixer_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yoheimuta/protolint/linter/fixer"
)

// TestBaseFixingFinallyKeepsUnmodifiedFileIntact guards against rewriting a file
// whose content the fixer did not change. See https://github.com/yoheimuta/protolint/issues/452
func TestBaseFixingFinallyKeepsUnmodifiedFileIntact(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "example.proto")
	content := []byte("syntax = \"proto3\";\n\nmessage Example {\n  string name = 1;\n}\n")
	if err := os.WriteFile(fileName, content, 0o600); err != nil {
		t.Fatal(err)
	}

	// Push the timestamp into the past so that a rewrite is observable
	// regardless of the filesystem timestamp granularity.
	past := time.Now().Add(-1 * time.Hour)
	if err := os.Chtimes(fileName, past, past); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(fileName)
	if err != nil {
		t.Fatal(err)
	}

	f, err := fixer.NewBaseFixing(fileName)
	if err != nil {
		t.Fatal(err)
	}
	// No edit is recorded: the rules found nothing to fix.
	if err = f.Finally(); err != nil {
		t.Fatal(err)
	}

	after, err := os.Stat(fileName)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("content changed: got %q, but want %q", got, content)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Errorf(
			"mtime of an unmodified file changed: before %v, after %v",
			before.ModTime(), after.ModTime(),
		)
	}
}

// TestBaseFixingFinallyWritesModifiedFile is the control: a real change must still be written.
func TestBaseFixingFinallyWritesModifiedFile(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "example.proto")
	content := []byte("syntax = \"proto3\";\n")
	if err := os.WriteFile(fileName, content, 0o600); err != nil {
		t.Fatal(err)
	}

	f, err := fixer.NewBaseFixing(fileName)
	if err != nil {
		t.Fatal(err)
	}
	f.ReplaceText(1, "proto3", "proto2")
	if err = f.Finally(); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if want := "syntax = \"proto2\";\n"; string(got) != want {
		t.Errorf("got %q, but want %q", got, want)
	}
}
