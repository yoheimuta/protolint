package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// OrderRule verifies that all files should be ordered in the following manner:
// 1. Syntax
// 2. Package
// 3. Imports (sorted)
// 4. File options
// 5. Everything else
// See https://developers.google.com/protocol-buffers/docs/style#file-structure.
type OrderRule struct{}

// NewOrderRule creates a new OrderRule.
func NewOrderRule() OrderRule {
	return OrderRule{}
}

// ID returns the ID of this rule.
func (r OrderRule) ID() string {
	return "ORDER"
}

// Purpose returns the purpose of this rule.
func (r OrderRule) Purpose() string {
	return "Verifies that all files should be ordered in the specific manner."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r OrderRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r OrderRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &orderVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		state:          initialOrderState,
		machine:        newOrderStateTransition(),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type orderVisitor struct {
	*visitor.BaseAddVisitor
	state   orderState
	machine orderStateTransition
}

func (v *orderVisitor) VisitSyntax(s *parser.Syntax) bool {
	next := v.machine.transit(v.state, syntaxVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(s.Meta.Pos, "Syntax should be located at the top. Check if the file is ordered in the correct manner.")
	}
	v.state = syntaxOrderState
	return false
}

func (v *orderVisitor) VisitPackage(p *parser.Package) bool {
	next := v.machine.transit(v.state, packageVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(p.Meta.Pos, "The order of Package is invalid. Check if the file is ordered in the correct manner.")
	}
	v.state = packageOrderState
	return false
}

func (v *orderVisitor) VisitImport(i *parser.Import) bool {
	next := v.machine.transit(v.state, importsVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(i.Meta.Pos, "The order of Import is invalid. Check if the file is ordered in the correct manner.")
	}
	v.state = importsOrderState
	return false
}

func (v *orderVisitor) VisitOption(o *parser.Option) bool {
	next := v.machine.transit(v.state, fileOptionsVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(o.Meta.Pos, "The order of Option is invalid. Check if the file is ordered in the correct manner.")
	}
	v.state = fileOptionsOrderState
	return false
}

func (v *orderVisitor) VisitMessage(m *parser.Message) bool {
	v.state = everythingElseOrderState
	return false
}

func (v *orderVisitor) VisitEnum(e *parser.Enum) bool {
	v.state = everythingElseOrderState
	return false
}

func (v *orderVisitor) VisitService(s *parser.Service) bool {
	v.state = everythingElseOrderState
	return false
}

func (v *orderVisitor) VisitExtend(e *parser.Extend) bool {
	v.state = everythingElseOrderState
	return false
}

type orderState int

const (
	invalidOrderState orderState = iota
	initialOrderState
	syntaxOrderState
	packageOrderState
	importsOrderState
	fileOptionsOrderState
	everythingElseOrderState
)

type orderEvent int

const (
	syntaxVisitEvent orderEvent = iota
	packageVisitEvent
	importsVisitEvent
	fileOptionsVisitEvent
)

type orderInput struct {
	state orderState
	event orderEvent
}

type orderStateTransition struct {
	f map[orderInput]orderState
}

func newOrderStateTransition() orderStateTransition {
	return orderStateTransition{
		f: map[orderInput]orderState{
			{
				state: initialOrderState,
				event: syntaxVisitEvent,
			}: syntaxOrderState,

			{
				state: initialOrderState,
				event: packageVisitEvent,
			}: packageOrderState,
			{
				state: syntaxOrderState,
				event: packageVisitEvent,
			}: packageOrderState,

			{
				state: initialOrderState,
				event: importsVisitEvent,
			}: importsOrderState,
			{
				state: syntaxOrderState,
				event: importsVisitEvent,
			}: importsOrderState,
			{
				state: packageOrderState,
				event: importsVisitEvent,
			}: importsOrderState,
			{
				state: importsOrderState,
				event: importsVisitEvent,
			}: importsOrderState,

			{
				state: initialOrderState,
				event: fileOptionsVisitEvent,
			}: fileOptionsOrderState,
			{
				state: syntaxOrderState,
				event: fileOptionsVisitEvent,
			}: fileOptionsOrderState,
			{
				state: packageOrderState,
				event: fileOptionsVisitEvent,
			}: fileOptionsOrderState,
			{
				state: importsOrderState,
				event: fileOptionsVisitEvent,
			}: fileOptionsOrderState,
			{
				state: fileOptionsOrderState,
				event: fileOptionsVisitEvent,
			}: fileOptionsOrderState,
		},
	}
}

func (t orderStateTransition) transit(
	state orderState,
	event orderEvent,
) orderState {
	out, ok := t.f[orderInput{state: state, event: event}]
	if !ok {
		return invalidOrderState
	}
	return out
}
