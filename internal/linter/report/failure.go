package report

import (
	"fmt"

	"github.com/yoheimuta/go-protoparser/parser/meta"
)

// Failure represents a lint error information.
type Failure struct {
	pos     meta.Position
	message string
}

// Failuref creates a new Failure and the formatting works like fmt.Sprintf.
func Failuref(
	pos meta.Position,
	format string,
	a ...interface{},
) Failure {
	return Failure{
		pos:     pos,
		message: fmt.Sprintf(format, a...),
	}
}

// String stringifies Failure.
func (f Failure) String() string {
	return fmt.Sprintf("[%s], %s", f.pos, f.message)
}
