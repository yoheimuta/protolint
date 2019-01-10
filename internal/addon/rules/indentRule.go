package rules

import (
	"bufio"
	"os"
	"strings"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

const (
	// 4 spaces
	defaultStyle = "    "
)

// IndentRule enforces a consistent indentation style.
type IndentRule struct {
	style string
}

// NewIndentRule creates a new IndentRule.
func NewIndentRule(
	style string,
) IndentRule {
	if len(style) == 0 {
		style = defaultStyle
	}

	return IndentRule{
		style: style,
	}
}

// ID returns the ID of this rule.
func (r IndentRule) ID() string {
	return "INDENT"
}

// Purpose returns the purpose of this rule.
func (r IndentRule) Purpose() string {
	return "Enforces a consistent indentation style."
}

// Apply applies the rule to the proto.
func (r IndentRule) Apply(
	proto *parser.Proto,
) ([]report.Failure, error) {
	lines, err := protoLines(proto)
	if err != nil {
		return nil, err
	}

	v := &indentVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(),
		style:          r.style,
		protoLines:     lines,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

func protoLines(
	proto *parser.Proto,
) (lines []string, err error) {
	fileName := proto.Meta.Filename
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := reader.Close()
		if err != nil {
			return
		}
		if closeErr != nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

type indentVisitor struct {
	*visitor.BaseAddVisitor
	style        string
	protoLines   []string
	currentLevel int
}

func (v indentVisitor) VisitEnum(e *parser.Enum) (next bool) {
	v.validateIndent(e.Meta.Pos)
	if e.Meta.Pos.Line < e.Meta.LastPos.Line {
		v.validateIndent(e.Meta.LastPos)
	}
	for _, comment := range e.Comments {
		v.validateIndent(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range e.EnumBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitEnumField(f *parser.EnumField) (next bool) {
	v.validateIndent(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitField(f *parser.Field) (next bool) {
	v.validateIndent(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitImport(i *parser.Import) (next bool) {
	v.validateIndent(i.Meta.Pos)
	for _, comment := range i.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitMapField(m *parser.MapField) (next bool) {
	v.validateIndent(m.Meta.Pos)
	for _, comment := range m.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitMessage(m *parser.Message) (next bool) {
	v.validateIndent(m.Meta.Pos)
	if m.Meta.Pos.Line < m.Meta.LastPos.Line {
		v.validateIndent(m.Meta.LastPos)
	}
	for _, comment := range m.Comments {
		v.validateIndent(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range m.MessageBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitOneof(o *parser.Oneof) (next bool) {
	v.validateIndent(o.Meta.Pos)
	if o.Meta.Pos.Line < o.Meta.LastPos.Line {
		v.validateIndent(o.Meta.LastPos)
	}
	for _, comment := range o.Comments {
		v.validateIndent(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, field := range o.OneofFields {
		field.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitOneofField(f *parser.OneofField) (next bool) {
	v.validateIndent(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitOption(o *parser.Option) (next bool) {
	v.validateIndent(o.Meta.Pos)
	for _, comment := range o.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitPackage(p *parser.Package) (next bool) {
	v.validateIndent(p.Meta.Pos)
	for _, comment := range p.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitReserved(r *parser.Reserved) (next bool) {
	v.validateIndent(r.Meta.Pos)
	for _, comment := range r.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitRPC(r *parser.RPC) (next bool) {
	v.validateIndent(r.Meta.Pos)
	if r.Meta.Pos.Line < r.Meta.LastPos.Line {
		v.validateIndent(r.Meta.LastPos)
	}
	for _, comment := range r.Comments {
		v.validateIndent(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range r.Options {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitService(s *parser.Service) (next bool) {
	v.validateIndent(s.Meta.Pos)
	if s.Meta.Pos.Line < s.Meta.LastPos.Line {
		v.validateIndent(s.Meta.LastPos)
	}
	for _, comment := range s.Comments {
		v.validateIndent(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range s.ServiceBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitSyntax(s *parser.Syntax) (next bool) {
	v.validateIndent(s.Meta.Pos)
	for _, comment := range s.Comments {
		v.validateIndent(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) validateIndent(
	pos meta.Position,
) {
	line := v.protoLines[pos.Line-1]
	leading := string([]rune(line)[0 : pos.Column-1])

	indentation := strings.Repeat(v.style, v.currentLevel)
	if leading == indentation {
		return
	}
	v.AddFailuref(
		pos,
		`Found an incorrect indentation style "%s". "%s" is correct.`,
		leading,
		indentation,
	)
}

func (v *indentVisitor) nest() func() {
	v.currentLevel++
	return func() {
		v.currentLevel--
	}
}
