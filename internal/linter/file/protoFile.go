package file

import (
	"bytes"

	protoparser "github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/internal/file"
)

// ProtoFile is a Protocol Buffer file.
type ProtoFile struct {
	// The path to the .proto file.
	// For physical files, it must be absolute and cleaned.
	// For in-memory virtual files, it must be the canonical relative display path.
	path string
	// The path to display in output.
	// This will be relative to the working directory, or the absolute path
	// if the file was outside the working directory.
	displayPath string
	// cachedData holds the raw byte content of the file.
	// It is used to minimize physical filesystem access.
	cachedData []byte
	// cachedProto holds the parsed protocol buffer definition.
	// It is used to avoid redundant parsing across multiple rules.
	cachedProto *parser.Proto
}

// NewProtoFile creates a new proto file.
func NewProtoFile(
	path string,
	displayPath string,
) *ProtoFile {
	return &ProtoFile{
		path:        path,
		displayPath: displayPath,
	}
}

// Parse parses a Protocol Buffer file.
func (f *ProtoFile) Parse(
	debug bool,
) (*parser.Proto, error) {
	if f.cachedProto != nil {
		return f.cachedProto, nil
	}

	if f.cachedData == nil {
		data, err := file.ReadFile(f.path)
		if err != nil {
			return nil, err
		}
		f.cachedData = data
	}

	reader := bytes.NewReader(f.cachedData)

	proto, err := protoparser.Parse(
		reader,
		protoparser.WithFilename(f.displayPath),
		protoparser.WithBodyIncludingComments(true),
		protoparser.WithDebug(debug),
	)
	if err != nil {
		return nil, err
	}

	f.cachedProto = proto
	return f.cachedProto, nil
}

// Path returns the path to the .proto file.
func (f *ProtoFile) Path() string {
	return f.path
}

// DisplayPath returns the path to display in output.
func (f *ProtoFile) DisplayPath() string {
	return f.displayPath
}

// ResetData clears the cached raw byte content of the file.
// This invalidates the local buffer, forcing the linter to fetch the most
// up-to-date source state (either from the physical disk or from the
// mutated in-memory virtual filesystem layer) on the next execution cycle.
func (f *ProtoFile) ResetData() {
	f.cachedData = nil
}

// ResetCache clears the cached protocol buffer definition.
// This is used during modifying commands like 'lint -fix' to ensure the next
// rule works with a fresh parse if the previous rule modified the file.
func (f *ProtoFile) ResetCache() {
	f.cachedProto = nil
}
