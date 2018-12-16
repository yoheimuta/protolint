package file

import (
	"os"

	"log"

	"github.com/yoheimuta/go-protoparser"
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
func (f ProtoFile) Parse() (*parser.Proto, error) {
	reader, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = reader.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	proto, err := protoparser.Parse(reader, protoparser.WithFilename(f.displayPath))
	if err != nil {
		return nil, err
	}
	return proto, nil
}
