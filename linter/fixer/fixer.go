package fixer

import (
	"io/ioutil"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/osutil"
)

// Fixer provides the ways to operate the proto content.
type Fixer interface {
	// NOTE: This method is insufficient to process unexpected multi-line contents.
	ReplaceText(line int, old, new string)
}

// Fixing adds the way to modify the proto file to Fixer.
type Fixing interface {
	Fixer
	Finally() error
}

// NewFixing creates a fixing, depending on fixMode.
func NewFixing(fixMode bool, proto *parser.Proto) (Fixing, error) {
	if fixMode {
		return newBaseFixing("", proto.Meta.Filename)
	}
	return NopFixing{}, nil
}

// BaseFixing implements Fixing.
type BaseFixing struct {
	content    []byte
	lineEnding string
	fileName   string
}

func newBaseFixing(optionalNewline, protoFileName string) (*BaseFixing, error) {
	content, err := ioutil.ReadFile(protoFileName)
	if err != nil {
		return nil, err
	}
	lineEnding := optionalNewline
	if len(lineEnding) == 0 {
		lineEnding, err = osutil.DetectLineEnding(string(content))
		if err != nil {
			return nil, err
		}
		if len(lineEnding) == 0 {
			lineEnding = "\n"
		}
	}
	return &BaseFixing{
		content:    content,
		lineEnding: lineEnding,
		fileName:   protoFileName,
	}, nil
}

// ReplaceText replaces the text at the line.
func (f *BaseFixing) ReplaceText(line int, old, new string) {
	lines := strings.Split(string(f.content), f.lineEnding)
	lines[line-1] = strings.Replace(lines[line-1], old, new, 1)
	f.content = []byte(strings.Join(lines, f.lineEnding))
}

// Finally writes the fixed content to the file.
func (f *BaseFixing) Finally() error {
	return osutil.WriteExistingFile(f.fileName, f.content)
}

// NopFixing does nothing.
type NopFixing struct{}

// ReplaceText noop
func (f NopFixing) ReplaceText(line int, old, new string) {}

// Finally noop
func (f NopFixing) Finally() error { return nil }
