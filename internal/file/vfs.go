package file

import (
	"bytes"
	"io"
	"os"
	"sync"
)

const (
	// StdinPath is the magic symbol representing stdin.
	StdinPath = "-"
	// StdinDisplayPath is the default filename used for displaying stdin results.
	StdinDisplayPath = "stdin.proto"
)

var (
	// virtualFiles maps file paths to their raw byte content in memory.
	// This allows shadowing disk files or providing stdin content to rules.
	virtualFiles = make(map[string][]byte)
	// mu protects concurrent access to the virtualFiles map.
	mu sync.RWMutex
)

// SetVirtualFile registers raw byte content for a specific path in the virtual file system.
// This is typically used to cache stdin content so it can be read multiple times by different rules.
func SetVirtualFile(path string, data []byte) {
	mu.Lock()
	defer mu.Unlock()
	virtualFiles[path] = data
}

// Open opens a file for reading.
// It returns a reader from memory if the path is registered in the virtual file system,
// otherwise it falls back to opening a physical file from the disk.
func Open(path string) (io.ReadCloser, error) {
	mu.RLock()
	data, ok := virtualFiles[path]
	mu.RUnlock()

	if ok {
		return io.NopCloser(bytes.NewReader(data)), nil
	}

	return os.Open(path)
}

// ReadFile reads the entire named file into a byte slice.
// It returns data from memory if the path is registered in the virtual file system,
// otherwise it reads from the physical disk using os.ReadFile.
func ReadFile(path string) ([]byte, error) {
	mu.RLock()
	data, ok := virtualFiles[path]
	mu.RUnlock()

	if ok {
		return data, nil
	}

	return os.ReadFile(path)
}
