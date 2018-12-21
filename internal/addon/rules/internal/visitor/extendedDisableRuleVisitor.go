package visitor

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolinter/internal/addon/rules/internal/disablerule"
	"github.com/yoheimuta/protolinter/internal/linter/report"
)

type extendedDisableRuleVisitor struct {
	inner       hasExtendedVisitor
	interpreter *disablerule.Interpreter
}

func newExtendedDisableRuleVisitor(
	inner hasExtendedVisitor,
	ruleID string,
) extendedDisableRuleVisitor {
	interpreter := disablerule.NewInterpreter(ruleID)
	return extendedDisableRuleVisitor{
		inner:       inner,
		interpreter: interpreter,
	}
}

func (v extendedDisableRuleVisitor) OnStart(p *parser.Proto) error  { return v.inner.OnStart(p) }
func (v extendedDisableRuleVisitor) Finally() error                 { return v.inner.Finally() }
func (v extendedDisableRuleVisitor) Failures() []report.Failure     { return v.inner.Failures() }
func (v extendedDisableRuleVisitor) VisitComment(c *parser.Comment) { v.inner.VisitComment(c) }
func (v extendedDisableRuleVisitor) VisitEmptyStatement(e *parser.EmptyStatement) (next bool) {
	return v.inner.VisitEmptyStatement(e)
}

func (v extendedDisableRuleVisitor) VisitEnum(e *parser.Enum) (next bool) {
	if v.interpreter.Interpret(e.Comments) {
		return true
	}
	return v.inner.VisitEnum(e)
}
func (v extendedDisableRuleVisitor) VisitEnumField(e *parser.EnumField) (next bool) {
	if v.interpreter.Interpret(e.Comments) {
		return true
	}
	return v.inner.VisitEnumField(e)
}
func (v extendedDisableRuleVisitor) VisitField(f *parser.Field) (next bool) {
	if v.interpreter.Interpret(f.Comments) {
		return true
	}
	return v.inner.VisitField(f)
}
func (v extendedDisableRuleVisitor) VisitImport(i *parser.Import) (next bool) {
	if v.interpreter.Interpret(i.Comments) {
		return true
	}
	return v.inner.VisitImport(i)
}
func (v extendedDisableRuleVisitor) VisitMapField(m *parser.MapField) (next bool) {
	if v.interpreter.Interpret(m.Comments) {
		return true
	}
	return v.inner.VisitMapField(m)
}
func (v extendedDisableRuleVisitor) VisitMessage(m *parser.Message) (next bool) {
	if v.interpreter.Interpret(m.Comments) {
		return true
	}
	return v.inner.VisitMessage(m)
}
func (v extendedDisableRuleVisitor) VisitOneof(o *parser.Oneof) (next bool) {
	if v.interpreter.Interpret(o.Comments) {
		return true
	}
	return v.inner.VisitOneof(o)
}
func (v extendedDisableRuleVisitor) VisitOneofField(o *parser.OneofField) (next bool) {
	if v.interpreter.Interpret(o.Comments) {
		return true
	}
	return v.inner.VisitOneofField(o)
}
func (v extendedDisableRuleVisitor) VisitOption(o *parser.Option) (next bool) {
	if v.interpreter.Interpret(o.Comments) {
		return true
	}
	return v.inner.VisitOption(o)
}
func (v extendedDisableRuleVisitor) VisitPackage(p *parser.Package) (next bool) {
	if v.interpreter.Interpret(p.Comments) {
		return true
	}
	return v.inner.VisitPackage(p)
}
func (v extendedDisableRuleVisitor) VisitReserved(r *parser.Reserved) (next bool) {
	if v.interpreter.Interpret(r.Comments) {
		return true
	}
	return v.inner.VisitReserved(r)
}
func (v extendedDisableRuleVisitor) VisitRPC(r *parser.RPC) (next bool) {
	if v.interpreter.Interpret(r.Comments) {
		return true
	}
	return v.inner.VisitRPC(r)
}
func (v extendedDisableRuleVisitor) VisitService(s *parser.Service) (next bool) {
	if v.interpreter.Interpret(s.Comments) {
		return true
	}
	return v.inner.VisitService(s)
}
func (v extendedDisableRuleVisitor) VisitSyntax(s *parser.Syntax) (next bool) {
	if v.interpreter.Interpret(s.Comments) {
		return true
	}
	return v.inner.VisitSyntax(s)
}
