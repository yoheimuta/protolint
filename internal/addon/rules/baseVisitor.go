package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolinter/internal/linter/report"
)

type baseVisitor struct{}

func (baseVisitor) OnStart(*parser.Proto) error { return nil }
func (baseVisitor) Finally() error              { return nil }

func (baseVisitor) VisitComment(*parser.Comment)                           {}
func (baseVisitor) VisitEmptyStatement(*parser.EmptyStatement) (next bool) { return true }
func (baseVisitor) VisitEnum(*parser.Enum) (next bool)                     { return true }
func (baseVisitor) VisitEnumField(*parser.EnumField) (next bool)           { return true }
func (baseVisitor) VisitField(*parser.Field) (next bool)                   { return true }
func (baseVisitor) VisitImport(*parser.Import) (next bool)                 { return true }
func (baseVisitor) VisitMapField(*parser.MapField) (next bool)             { return true }
func (baseVisitor) VisitMessage(*parser.Message) (next bool)               { return true }
func (baseVisitor) VisitOneof(*parser.Oneof) (next bool)                   { return true }
func (baseVisitor) VisitOneofField(*parser.OneofField) (next bool)         { return true }
func (baseVisitor) VisitOption(*parser.Option) (next bool)                 { return true }
func (baseVisitor) VisitPackage(*parser.Package) (next bool)               { return true }
func (baseVisitor) VisitReserved(*parser.Reserved) (next bool)             { return true }
func (baseVisitor) VisitRPC(*parser.RPC) (next bool)                       { return true }
func (baseVisitor) VisitService(*parser.Service) (next bool)               { return true }
func (baseVisitor) VisitSyntax(*parser.Syntax) (next bool)                 { return true }

type baseAddVisitor struct {
	baseVisitor
	failures []report.Failure
}

func newBaseAddVisitor() *baseAddVisitor {
	return &baseAddVisitor{}
}

// Failures returns the accumulated failures.
func (v *baseAddVisitor) Failures() []report.Failure {
	return v.failures
}

func (v *baseAddVisitor) addFailuref(
	pos meta.Position,
	format string,
	a ...interface{},
) {
	v.failures = append(v.failures, report.Failuref(pos, format, a...))
}

type extendedVisitor interface {
	parser.Visitor

	// OnStart is called when visiting is started.
	OnStart(*parser.Proto) error
	// Finally is called when visiting is done.
	Finally() error
	// Failures returns the accumulated failures.
	Failures() []report.Failure
}

func runVisitor(
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
