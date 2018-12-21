package visitor

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolinter/internal/linter/report"
)

type extendedVisitor interface {
	parser.Visitor

	// OnStart is called when visiting is started.
	OnStart(*parser.Proto) error
	// Finally is called when visiting is done.
	Finally() error
	// Failures returns the accumulated failures.
	Failures() []report.Failure
}

// RunVisitor dispatches the call to the visitor.
func RunVisitor(
	visitor extendedVisitor,
	proto *parser.Proto,
) ([]report.Failure, error) {
	if err := visitor.OnStart(proto); err != nil {
		return nil, err
	}
	proto.Accept(visitor)
	if err := visitor.Finally(); err != nil {
		return nil, err
	}
	return visitor.Failures(), nil
}
