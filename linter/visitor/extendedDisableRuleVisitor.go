package visitor

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/linter/visitor/internal/disablerule"
)

// TODO: To work `enable comments` more precisely, this implementation has to be modified.
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

func (v extendedDisableRuleVisitor) OnStart(p *parser.Proto) error { return v.inner.OnStart(p) }
func (v extendedDisableRuleVisitor) Finally() error                { return v.inner.Finally() }
func (v extendedDisableRuleVisitor) Failures() []report.Failure    { return v.inner.Failures() }
func (v extendedDisableRuleVisitor) VisitEmptyStatement(e *parser.EmptyStatement) (next bool) {
	return v.inner.VisitEmptyStatement(e)
}

func (v extendedDisableRuleVisitor) VisitComment(c *parser.Comment) {
	if v.interpreter.Interpret([]*parser.Comment{c}) {
		return
	}
	v.inner.VisitComment(c)
}
func (v extendedDisableRuleVisitor) VisitEnum(e *parser.Enum) (next bool) {
	if v.interpreter.Interpret(e.Comments, e.InlineComment, e.InlineCommentBehindLeftCurly) {
		return true
	}
	return v.inner.VisitEnum(e)
}
func (v extendedDisableRuleVisitor) VisitEnumField(e *parser.EnumField) (next bool) {
	if v.interpreter.Interpret(e.Comments, e.InlineComment) {
		return true
	}
	return v.inner.VisitEnumField(e)
}
func (v extendedDisableRuleVisitor) VisitExtend(m *parser.Extend) (next bool) {
	if v.interpreter.Interpret(m.Comments, m.InlineComment, m.InlineCommentBehindLeftCurly) {
		return true
	}
	return v.inner.VisitExtend(m)
}
func (v extendedDisableRuleVisitor) VisitExtensions(m *parser.Extensions) (next bool) {
	if v.interpreter.Interpret(m.Comments, m.InlineComment) {
		return true
	}
	return v.inner.VisitExtensions(m)
}
func (v extendedDisableRuleVisitor) VisitField(f *parser.Field) (next bool) {
	if v.interpreter.Interpret(f.Comments, f.InlineComment) {
		return true
	}
	return v.inner.VisitField(f)
}
func (v extendedDisableRuleVisitor) VisitGroupField(m *parser.GroupField) (next bool) {
	if v.interpreter.Interpret(m.Comments, m.InlineComment, m.InlineCommentBehindLeftCurly) {
		return true
	}
	return v.inner.VisitGroupField(m)
}
func (v extendedDisableRuleVisitor) VisitImport(i *parser.Import) (next bool) {
	if v.interpreter.Interpret(i.Comments, i.InlineComment) {
		return true
	}
	return v.inner.VisitImport(i)
}
func (v extendedDisableRuleVisitor) VisitMapField(m *parser.MapField) (next bool) {
	if v.interpreter.Interpret(m.Comments, m.InlineComment) {
		return true
	}
	return v.inner.VisitMapField(m)
}
func (v extendedDisableRuleVisitor) VisitMessage(m *parser.Message) (next bool) {
	if v.interpreter.Interpret(m.Comments, m.InlineComment, m.InlineCommentBehindLeftCurly) {
		return true
	}
	return v.inner.VisitMessage(m)
}
func (v extendedDisableRuleVisitor) VisitOneof(o *parser.Oneof) (next bool) {
	if v.interpreter.Interpret(o.Comments, o.InlineComment, o.InlineCommentBehindLeftCurly) {
		return true
	}
	return v.inner.VisitOneof(o)
}
func (v extendedDisableRuleVisitor) VisitOneofField(o *parser.OneofField) (next bool) {
	if v.interpreter.Interpret(o.Comments, o.InlineComment) {
		return true
	}
	return v.inner.VisitOneofField(o)
}
func (v extendedDisableRuleVisitor) VisitOption(o *parser.Option) (next bool) {
	if v.interpreter.Interpret(o.Comments, o.InlineComment) {
		return true
	}
	return v.inner.VisitOption(o)
}
func (v extendedDisableRuleVisitor) VisitPackage(p *parser.Package) (next bool) {
	if v.interpreter.Interpret(p.Comments, p.InlineComment) {
		return true
	}
	return v.inner.VisitPackage(p)
}
func (v extendedDisableRuleVisitor) VisitReserved(r *parser.Reserved) (next bool) {
	if v.interpreter.Interpret(r.Comments, r.InlineComment) {
		return true
	}
	return v.inner.VisitReserved(r)
}
func (v extendedDisableRuleVisitor) VisitRPC(r *parser.RPC) (next bool) {
	if v.interpreter.Interpret(r.Comments, r.InlineComment) {
		return true
	}
	return v.inner.VisitRPC(r)
}
func (v extendedDisableRuleVisitor) VisitService(s *parser.Service) (next bool) {
	if v.interpreter.Interpret(s.Comments, s.InlineComment, s.InlineCommentBehindLeftCurly) {
		return true
	}
	return v.inner.VisitService(s)
}
func (v extendedDisableRuleVisitor) VisitSyntax(s *parser.Syntax) (next bool) {
	if v.interpreter.Interpret(s.Comments, s.InlineComment) {
		return true
	}
	return v.inner.VisitSyntax(s)
}
