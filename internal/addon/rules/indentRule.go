package rules

import (
	"strings"

	"github.com/yoheimuta/protolint/internal/osutil"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

const (
	// Use an indent of 2 spaces.
	// See https://developers.google.com/protocol-buffers/docs/style#standard-file-formatting
	defaultStyle = "  "

	defaultNewline = "\n"
)

// IndentRule enforces a consistent indentation style.
type IndentRule struct {
	style   string
	newline string
	fixMode bool
}

// NewIndentRule creates a new IndentRule.
func NewIndentRule(
	style string,
	newline string,
	fixMode bool,
) IndentRule {
	if len(style) == 0 {
		style = defaultStyle
	}
	if len(newline) == 0 {
		newline = defaultNewline
	}

	return IndentRule{
		style:   style,
		newline: newline,
		fixMode: fixMode,
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

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r IndentRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r IndentRule) Apply(
	proto *parser.Proto,
) ([]report.Failure, error) {
	fileName := proto.Meta.Filename
	lines, err := osutil.ReadAllLines(fileName, r.newline)
	if err != nil {
		return nil, err
	}

	v := &indentVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		style:          r.style,
		protoLines:     lines,
		fixMode:        r.fixMode,
		newline:        r.newline,
		protoFileName:  fileName,
		indentFixes:    make(map[int]indentFix),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type indentFix struct {
	currentChars int
	replacement  string
}

type indentVisitor struct {
	*visitor.BaseAddVisitor
	style        string
	protoLines   []string
	currentLevel int

	fixMode       bool
	newline       string
	protoFileName string
	indentFixes   map[int]indentFix
}

func (v indentVisitor) Finally() error {
	if v.fixMode {
		return v.fix()
	}
	return nil
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

	if v.fixMode {
		v.indentFixes[pos.Line-1] = indentFix{
			currentChars: len(leading),
			replacement:  indentation,
		}
	}
}

func (v *indentVisitor) nest() func() {
	v.currentLevel++
	return func() {
		v.currentLevel--
	}
}

func (v indentVisitor) fix() error {
	var shouldFixed bool

	var fixedLines []string
	for i, line := range v.protoLines {
		if fix, ok := v.indentFixes[i]; ok {
			line = fix.replacement + line[fix.currentChars:]
			shouldFixed = true
		}
		fixedLines = append(fixedLines, line)
	}

	if !shouldFixed {
		return nil
	}
	return osutil.WriteLinesToExistingFile(v.protoFileName, fixedLines, v.newline)
}
