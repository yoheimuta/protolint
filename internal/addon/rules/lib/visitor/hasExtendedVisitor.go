package visitor

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

type hasExtendedVisitor interface {
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
	visitor hasExtendedVisitor,
	proto *parser.Proto,
	ruleID string,
) ([]report.Failure, error) {
	v := newExtendedDisableRuleVisitor(
		visitor,
		ruleID,
	)

	if err := v.OnStart(proto); err != nil {
		return nil, err
	}
	proto.Accept(v)
	if err := visitor.Finally(); err != nil {
		return nil, err
	}
	return v.Failures(), nil
}
