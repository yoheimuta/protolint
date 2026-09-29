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
	// virtualFiles holds the raw byte content of files in memory.
	// It is primarily utilized to store and manage data streams passed via stdin.
	virtualFiles = make(map[string][]byte)
	// mu protects concurrent access to the virtualFiles map.
	mu sync.RWMutex
)

// IsStdin checks whether the given path corresponds to an active in-memory virtual file session.
// This allows linters and specific text rules to selectively skip or bypass side effects
// that are invalid for memory streams, such as physical file renaming or filesystem stats.
func IsStdin(path string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := virtualFiles[path]
	return ok
}

// SetVirtualFile registers or updates the in-memory content for a given file path.
// This is used to pre-populate the cache during arguments initialization or
// to accumulate consecutive automated transformations (fixing/auto-disabling) in-memory.
func SetVirtualFile(path string, data []byte) {
	mu.Lock()
	defer mu.Unlock()
	virtualFiles[path] = data
}

// Open opens the named file for reading.
// If the targeted path is registered within the virtual file system, it provides
// an in-memory io.ReadCloser wrapper to allow multiple reads. Otherwise, it
// transparently falls back to opening a physical file from the disk via os.Open.
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
// If the path is registered within the virtual file system, it returns the content
// directly from memory, ensuring subsequent rules evaluate the transformed buffers.
// Otherwise, it falls back to reading from the physical storage using os.ReadFile.
func ReadFile(path string) ([]byte, error) {
	mu.RLock()
	data, ok := virtualFiles[path]
	mu.RUnlock()

	if ok {
		return data, nil
	}

	return os.ReadFile(path)
}

// ResetVFS clears all virtual file entries.
// This is required to isolate independent execution cycles when protolint is invoked
// via the public lib.Lint() API inside a single process.
func ResetVFS() {
	mu.Lock()
	defer mu.Unlock()
	virtualFiles = make(map[string][]byte)
}
