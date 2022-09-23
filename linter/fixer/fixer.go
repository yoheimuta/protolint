package fixer

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/osutil"
)

// TextEdit represents the replacement of the code between Pos and End with the new text.
type TextEdit struct {
	Pos     int
	End     int // Inclusive. If the target is abc, pos and end are 1 and 3, respectively.
	NewText []byte
}

// Fixer provides the ways to operate the proto content.
type Fixer interface {
	// NOTE: This method is insufficient to process unexpected multi-line contents.
	ReplaceText(line int, old, new string)
	ReplaceAll(proc func(lines []string) []string)

	SearchAndReplace(startPos meta.Position, lex func(lex *lexer.Lexer) TextEdit) error
	ReplaceContent(proc func(content []byte) []byte)

	Lines() []string
}

// Fixing adds the way to modify the proto file to Fixer.
type Fixing interface {
	Fixer
	Finally() error
}

// NewFixing creates a fixing, depending on fixMode.
func NewFixing(fixMode bool, proto *parser.Proto) (Fixing, error) {
	if fixMode {
		return NewBaseFixing(proto.Meta.Filename)
	}
	return NopFixing{}, nil
}

// BaseFixing implements Fixing.
type BaseFixing struct {
	content    []byte
	lineEnding string
	fileName   string
	textEdits  []TextEdit
}

// NewBaseFixing creates a BaseFixing.
func NewBaseFixing(protoFileName string) (*BaseFixing, error) {
	content, err := ioutil.ReadFile(protoFileName)
	if err != nil {
		return nil, err
	}

	// Regardless of the actual dominant line ending, the fixer will go with LF
	// because the parser recognizes only LF as a line ending.
	//
	// It will work for most cases like used LF, CRLF, and a mix of LF and CRLF.
	// See also https://github.com/yoheimuta/protolint/issues/280.
	lineEnding := "\n"

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

// ReplaceAll replaces the lines.
func (f *BaseFixing) ReplaceAll(proc func(lines []string) []string) {
	lines := strings.Split(string(f.content), f.lineEnding)
	lines = proc(lines)
	f.content = []byte(strings.Join(lines, f.lineEnding))
}

// SearchAndReplace locates text edits and replaces with them.
func (f *BaseFixing) SearchAndReplace(startPos meta.Position, lex func(lex *lexer.Lexer) TextEdit) error {
	r := bytes.NewReader(f.content)
	_, err := r.Seek(int64(startPos.Offset), 0)
	if err != nil {
		return err
	}

	l := lexer.NewLexer(r)
	t := lex(l)
	t.Pos += startPos.Offset
	t.End += startPos.Offset
	f.textEdits = append(f.textEdits, t)
	return nil
}

// ReplaceContent replaces entire content.
func (f *BaseFixing) ReplaceContent(proc func(content []byte) []byte) {
	f.content = proc(f.content)
}

// Lines returns the line format of f.content.
func (f *BaseFixing) Lines() []string {
	return strings.Split(string(f.content), f.lineEnding)
}

// Finally writes the fixed content to the file.
func (f *BaseFixing) Finally() error {
	diff := 0
	for _, t := range f.textEdits {
		t.Pos += diff
		t.End += diff
		f.content = append(f.content[:t.Pos], append(t.NewText, f.content[t.End+1:]...)...)
		diff += len(t.NewText) - (t.End - t.Pos + 1)
	}
	return osutil.WriteExistingFile(f.fileName, f.content)
}

// Replace records a textedit to replace the old with the next later.
func (f *BaseFixing) Replace(t TextEdit) {
	f.textEdits = append(f.textEdits, t)
}

// Content returns f.content.
func (f *BaseFixing) Content() []byte {
	return f.content
}

// LineEnding is a detected line ending.
func (f *BaseFixing) LineEnding() string {
	return f.lineEnding
}

// NopFixing does nothing.
type NopFixing struct{}

// ReplaceText noop
func (f NopFixing) ReplaceText(line int, old, new string) {}

// ReplaceAll noop
func (f NopFixing) ReplaceAll(proc func(lines []string) []string) {}

// SearchAndReplace noop
func (f NopFixing) SearchAndReplace(startPos meta.Position, lex func(lexer *lexer.Lexer) TextEdit) error {
	return nil
}

// ReplaceContent noop.
func (f NopFixing) ReplaceContent(proc func(content []byte) []byte) {}

// Lines noop.
func (f NopFixing) Lines() []string { return []string{} }

// Finally noop
func (f NopFixing) Finally() error { return nil }
