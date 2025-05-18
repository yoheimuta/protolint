package visitor

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
)

// HasExtendedVisitor is a required interface given to RunVisitor.
type HasExtendedVisitor interface {
	parser.Visitor

	// Finally is called when visiting is done.
	Finally(*parser.Proto) error
	// Failures returns the accumulated failures.
	Failures() []report.Failure
}

// RunVisitor dispatches the call to the visitor.
func RunVisitor(
	visitor HasExtendedVisitor,
	proto *parser.Proto,
	ruleID string,
) ([]report.Failure, error) {
	return RunVisitorAutoDisable(visitor, proto, ruleID, autodisable.Noop)
}

// RunVisitorAutoDisable dispatches the call to the visitor.
func RunVisitorAutoDisable(
	visitor HasExtendedVisitor,
	proto *parser.Proto,
	ruleID string,
	autodisableType autodisable.PlacementType,
) ([]report.Failure, error) {
	// This check is just for existing test cases.
	protoFilename := ""
	if proto.Meta != nil {
		protoFilename = proto.Meta.Filename
	}
	autoDisabled, err := newExtendedAutoDisableVisitor(visitor, ruleID, protoFilename, autodisableType)
	if err != nil {
		return nil, err
	}
	disabled := newExtendedDisableRuleVisitor(
		autoDisabled,
		ruleID,
	)

	proto.Accept(disabled)
	if err := disabled.Finally(proto); err != nil {
		return nil, err
	}
	return disabled.Failures(), nil
}
