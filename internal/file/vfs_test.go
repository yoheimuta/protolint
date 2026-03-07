package file_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/yoheimuta/protolint/internal/file"
)

func TestVFS_Concurrency(t *testing.T) {
	var wg sync.WaitGroup
	workers := 100

	for i := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			path := fmt.Sprintf("virtual-%d.proto", id)
			content := fmt.Appendf(nil, "content-%d", id)

			// Wtite to the map with mutex
			file.SetVirtualFile(path, content)

			// Read from the map
			got, err := file.ReadFile(path)
			if err != nil {
				t.Errorf("worker %d: failed to read: %v", id, err)
				return
			}

			if !bytes.Equal(got, content) {
				t.Errorf("worker %d: data corruption! expected %s, got %s", id, content, got)
			}
		}(i)
	}
	wg.Wait()
}

func TestOpen_VirtualAndPhysical(t *testing.T) {
	// 1. Test Virtual Access
	virtualPath := "virtual_test.proto"
	virtualContent := []byte("syntax = 'proto3';")
	file.SetVirtualFile(virtualPath, virtualContent)
	defer file.SetVirtualFile(virtualPath, nil)

	r, err := file.Open(virtualPath)
	if err != nil {
		t.Fatalf("Open(virtual) failed: %v", err)
	}

	got, err := io.ReadAll(r)
	_ = r.Close() // Essential to test if it closes without panic
	if err != nil {
		t.Fatalf("ReadAll(virtual) failed: %v", err)
	}
	if !bytes.Equal(got, virtualContent) {
		t.Errorf("expected %s, got %s", virtualContent, got)
	}

	// 2. Test Physical Fallback (using a real file)
	tmpFile, _ := os.CreateTemp("", "physical.proto")
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	physicalContent := []byte("physical content")
	if err := os.WriteFile(tmpFile.Name(), physicalContent, 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	r2, err := file.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Open(physical) failed: %v", err)
	}
	got2, _ := io.ReadAll(r2)
	_ = r2.Close()
	if !bytes.Equal(got2, physicalContent) {
		t.Errorf("expected %s, got %s", physicalContent, got2)
	}
}
