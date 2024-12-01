package visitor

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
)

type extendedAutoDisableVisitor struct {
	inner     HasExtendedVisitor
	automator autodisable.PlacementStrategy
}

func newExtendedAutoDisableVisitor(
	inner HasExtendedVisitor,
	ruleID string,
	protoFilename string,
	placementType autodisable.PlacementType,
) (*extendedAutoDisableVisitor, error) {
	automator, err := autodisable.NewPlacementStrategy(placementType, protoFilename, ruleID)
	if err != nil {
		return nil, err
	}

	return &extendedAutoDisableVisitor{
		inner:     inner,
		automator: automator,
	}, nil
}

func (v *extendedAutoDisableVisitor) OnStart(p *parser.Proto) error { return v.inner.OnStart(p) }
func (v *extendedAutoDisableVisitor) Finally() error {
	err := v.automator.Finalize()
	if err != nil {
		return err
	}
	return v.inner.Finally()
}
func (v *extendedAutoDisableVisitor) Failures() []report.Failure { return v.inner.Failures() }

func (v *extendedAutoDisableVisitor) VisitEdition(s *parser.Edition) (next bool) {
	return v.inner.VisitEdition(s)
}

func (v *extendedAutoDisableVisitor) VisitEmptyStatement(e *parser.EmptyStatement) (next bool) {
	return v.inner.VisitEmptyStatement(e)
}

func (v *extendedAutoDisableVisitor) VisitComment(c *parser.Comment) {
	v.inner.VisitComment(c)
}

func (v *extendedAutoDisableVisitor) VisitEnum(e *parser.Enum) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitEnum(e)
	}, func(offset int) {
		v.automator.Disable(offset, e.Comments, e.InlineCommentBehindLeftCurly)
	})
}

func (v *extendedAutoDisableVisitor) VisitEnumField(e *parser.EnumField) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitEnumField(e)
	}, func(offset int) {
		v.automator.Disable(offset, e.Comments, e.InlineComment)
	})
}

func (v *extendedAutoDisableVisitor) VisitExtend(m *parser.Extend) (next bool) {
	return v.inner.VisitExtend(m)
}

func (v *extendedAutoDisableVisitor) VisitExtensions(m *parser.Extensions) (next bool) {
	return v.inner.VisitExtensions(m)
}

func (v *extendedAutoDisableVisitor) VisitField(f *parser.Field) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitField(f)
	}, func(offset int) {
		v.automator.Disable(offset, f.Comments, f.InlineComment)
	})
}

func (v *extendedAutoDisableVisitor) VisitGroupField(m *parser.GroupField) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitGroupField(m)
	}, func(offset int) {
		v.automator.Disable(offset, m.Comments, m.InlineComment)
	})
}

func (v *extendedAutoDisableVisitor) VisitImport(i *parser.Import) (next bool) {
	return v.inner.VisitImport(i)
}

func (v *extendedAutoDisableVisitor) VisitMapField(m *parser.MapField) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitMapField(m)
	}, func(offset int) {
		v.automator.Disable(offset, m.Comments, m.InlineComment)
	})
}

func (v *extendedAutoDisableVisitor) VisitMessage(m *parser.Message) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitMessage(m)
	}, func(offset int) {
		v.automator.Disable(offset, m.Comments, m.InlineCommentBehindLeftCurly)
	})
}

func (v *extendedAutoDisableVisitor) VisitOneof(o *parser.Oneof) (next bool) {
	return v.inner.VisitOneof(o)
}

func (v *extendedAutoDisableVisitor) VisitOneofField(o *parser.OneofField) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitOneofField(o)
	}, func(offset int) {
		v.automator.Disable(offset, o.Comments, o.InlineComment)
	})
}

func (v *extendedAutoDisableVisitor) VisitOption(o *parser.Option) (next bool) {
	return v.inner.VisitOption(o)
}

func (v *extendedAutoDisableVisitor) VisitPackage(p *parser.Package) (next bool) {
	return v.inner.VisitPackage(p)
}

func (v *extendedAutoDisableVisitor) VisitReserved(r *parser.Reserved) (next bool) {
	return v.inner.VisitReserved(r)
}

func (v *extendedAutoDisableVisitor) VisitRPC(r *parser.RPC) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitRPC(r)
	}, func(offset int) {
		v.automator.Disable(offset, r.Comments, r.InlineComment)
	})
}

func (v *extendedAutoDisableVisitor) VisitService(s *parser.Service) (next bool) {
	return v.doIfFailure(func() bool {
		return v.inner.VisitService(s)
	}, func(offset int) {
		v.automator.Disable(offset, s.Comments, s.InlineCommentBehindLeftCurly)
	})
}

func (v *extendedAutoDisableVisitor) VisitSyntax(s *parser.Syntax) (next bool) {
	return v.inner.VisitSyntax(s)
}

func (v *extendedAutoDisableVisitor) doIfFailure(
	visit func() bool,
	disable func(int),
) bool {
	prev := v.inner.Failures()
	next := visit()
	curr := v.inner.Failures()
	if len(prev) == len(curr) {
		return next
	}
	disable(curr[len(curr)-1].Pos().Offset)
	return next
}
