package file_test

import (
	"os"
	"testing"

	"github.com/yoheimuta/protolint/internal/linter/file"
)

func TestProtoFile_ParseCachingAndReset(t *testing.T) {
	content := []byte("syntax = 'proto3';")
	tmpFile, _ := os.CreateTemp("", "test.proto")
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	if err := os.WriteFile(tmpFile.Name(), content, 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	f := file.NewProtoFile(tmpFile.Name(), "test.proto")

	// Parse and fill the cache in
	p1, err := f.Parse(false)
	if err != nil {
		t.Fatal(err)
	}

	// Change the file on disk
	newContent := []byte("syntax = 'proto3'; message New {}")
	if err := os.WriteFile(tmpFile.Name(), newContent, 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	// Parsing with no reset should return p1
	p2, _ := f.Parse(false)
	if p1 != p2 {
		t.Errorf("expected cached protocol buffer definition, but got a new one")
	}

	// Reset cache and data (imitate a modifying mode)
	f.ResetCache()
	f.ResetData()

	// Parsing after reset should return a new data
	p3, _ := f.Parse(false)
	if p1 == p3 {
		t.Errorf("expected a new protocol buffer definition after reset, but got the old one")
	}
}
