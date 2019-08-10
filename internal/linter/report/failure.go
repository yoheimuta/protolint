package report

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yoheimuta/go-protoparser/parser/meta"
)

// Failure represents a lint error information.
type Failure struct {
	pos     meta.Position
	message string
	ruleID  string
}

// Failuref creates a new Failure and the formatting works like fmt.Sprintf.
func Failuref(
	pos meta.Position,
	ruleID string,
	format string,
	a ...interface{},
) Failure {
	return Failure{
		pos:     pos,
		message: fmt.Sprintf(format, a...),
		ruleID:  ruleID,
	}
}

// String stringifies Failure.
func (f Failure) String() string {
	return fmt.Sprintf("[%s] %s", f.pos, f.message)
}

// Message returns a raw message.
func (f Failure) Message() string {
	return f.message
}

// Pos returns a raw position.
func (f Failure) Pos() meta.Position {
	return f.pos
}

// RuleID returns a rule ID.
func (f Failure) RuleID() string {
	return f.ruleID
}

// FilenameWithoutExt returns a filename without the extension.
func (f Failure) FilenameWithoutExt() string {
	name := f.pos.Filename
	return strings.TrimSuffix(name, filepath.Ext(name))
}
