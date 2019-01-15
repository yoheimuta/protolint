package file

import (
	"os"

	protoparser "github.com/yoheimuta/go-protoparser"
	"github.com/yoheimuta/go-protoparser/parser"
)

// ProtoFile is a Protocol Buffer file.
type ProtoFile struct {
	// The path to the .proto file.
	// Must be absolute.
	// Must be cleaned.
	path string
	// The path to display in output.
	// This will be relative to the working directory, or the absolute path
	// if the file was outside the working directory.
	displayPath string
}

// NewProtoFile creates a new proto file.
func NewProtoFile(
	path string,
	displayPath string,
) ProtoFile {
	return ProtoFile{
		path:        path,
		displayPath: displayPath,
	}
}

// Parse parses a Protocol Buffer file.
func (f ProtoFile) Parse(
	debug bool,
) (_ *parser.Proto, err error) {
	reader, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := reader.Close()
		if err != nil {
			return
		}
		if closeErr != nil {
			err = closeErr
		}
	}()

	proto, err := protoparser.Parse(
		reader,
		protoparser.WithFilename(f.displayPath),
		protoparser.WithBodyIncludingComments(true),
		protoparser.WithDebug(debug),
	)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// Path returns the path to the .proto file.
func (f ProtoFile) Path() string {
	return f.path
}

// DisplayPath returns the path to display in output.
func (f ProtoFile) DisplayPath() string {
	return f.displayPath
}
