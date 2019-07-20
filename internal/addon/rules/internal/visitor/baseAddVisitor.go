package visitor

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

// BaseAddVisitor represents a base visitor which can accumulate failures.
type BaseAddVisitor struct {
	BaseVisitor
	failures []report.Failure
}

// NewBaseAddVisitor creates a BaseAddVisitor.
func NewBaseAddVisitor() *BaseAddVisitor {
	return &BaseAddVisitor{}
}

// Failures returns the accumulated failures.
func (v *BaseAddVisitor) Failures() []report.Failure {
	return v.failures
}

// AddFailuref adds to the internal buffer and the formatting works like fmt.Sprintf.
func (v *BaseAddVisitor) AddFailuref(
	pos meta.Position,
	format string,
	a ...interface{},
) {
	v.failures = append(v.failures, report.Failuref(pos, format, a...))
}

// AddFailurefWithProtoMeta adds to the internal buffer and the formatting works like fmt.Sprintf.
func (v *BaseAddVisitor) AddFailurefWithProtoMeta(
	p *parser.ProtoMeta,
	format string,
	a ...interface{},
) {
	v.AddFailuref(
		meta.Position{
			Filename: p.Filename,
			Offset:   0,
			Line:     1,
			Column:   1,
		},
		format,
		a...,
	)
}
