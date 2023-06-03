package visitor

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/fixer"
)

// BaseFixableVisitor represents a base visitor which can fix failures.
type BaseFixableVisitor struct {
	*BaseAddVisitor

	Fixer     fixer.Fixer
	finallyFn func() error
}

// NewBaseFixableVisitor creates a BaseFixableVisitor.
func NewBaseFixableVisitor(
	ruleID string,
	fixMode bool,
	proto *parser.Proto,
	severity string,
) (*BaseFixableVisitor, error) {
	f, err := fixer.NewFixing(fixMode, proto)
	if err != nil {
		return nil, err
	}
	return &BaseFixableVisitor{
		BaseAddVisitor: NewBaseAddVisitor(ruleID, severity),
		Fixer:          f,
		finallyFn:      f.Finally,
	}, nil
}

// Finally fixes the proto file by overwriting it.
func (v *BaseFixableVisitor) Finally() error {
	err := v.finallyFn()
	if err != nil {
		return err
	}
	return v.BaseAddVisitor.Finally()
}
