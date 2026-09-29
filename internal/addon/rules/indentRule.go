package rules

import (
	"sort"
	"strings"
	"unicode"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

const (
	// Use an indent of 2 spaces.
	// See https://developers.google.com/protocol-buffers/docs/style#standard-file-formatting
	defaultStyle = "  "
)

// IndentRule enforces a consistent indentation style.
type IndentRule struct {
	RuleWithSeverity
	style            string
	notInsertNewline bool
	fixMode          bool
}

// NewIndentRule creates a new IndentRule.
func NewIndentRule(
	severity rule.Severity,
	style string,
	notInsertNewline bool,
	fixMode bool,
) IndentRule {
	if len(style) == 0 {
		style = defaultStyle
	}

	return IndentRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		style:            style,
		notInsertNewline: notInsertNewline,
		fixMode:          fixMode,
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
	base, err := visitor.NewBaseFixableVisitor(r.ID(), true, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &indentVisitor{
		BaseFixableVisitor: base,
		style:              r.style,
		fixMode:            r.fixMode,
		notInsertNewline:   r.notInsertNewline,
		indentFixes:        make(map[int][]indentFix),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type indentFix struct {
	currentChars int
	replacement  string
	level        int
	pos          meta.Position
	isLast       bool
}

type indentVisitor struct {
	*visitor.BaseFixableVisitor
	style        string
	currentLevel int

	fixMode          bool
	notInsertNewline bool
	indentFixes      map[int][]indentFix
}

func (v indentVisitor) Finally(proto *parser.Proto) error {
	if v.fixMode {
		return v.fix(proto)
	}
	return nil
}

func (v indentVisitor) VisitEnum(e *parser.Enum) (next bool) {
	v.validateIndentLeading(e.Meta.Pos)
	defer func() { v.validateIndentLast(e.Meta.LastPos) }()
	for _, comment := range e.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range e.EnumBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitEnumField(f *parser.EnumField) (next bool) {
	v.validateIndentLeading(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitExtend(e *parser.Extend) (next bool) {
	v.validateIndentLeading(e.Meta.Pos)
	defer func() { v.validateIndentLast(e.Meta.LastPos) }()
	for _, comment := range e.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range e.ExtendBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitField(f *parser.Field) (next bool) {
	v.validateIndentLeading(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitGroupField(f *parser.GroupField) (next bool) {
	v.validateIndentLeading(f.Meta.Pos)
	defer func() { v.validateIndentLast(f.Meta.LastPos) }()
	for _, comment := range f.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range f.MessageBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitImport(i *parser.Import) (next bool) {
	v.validateIndentLeading(i.Meta.Pos)
	for _, comment := range i.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitMapField(m *parser.MapField) (next bool) {
	v.validateIndentLeading(m.Meta.Pos)
	for _, comment := range m.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitMessage(m *parser.Message) (next bool) {
	v.validateIndentLeading(m.Meta.Pos)
	defer func() { v.validateIndentLast(m.Meta.LastPos) }()
	for _, comment := range m.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range m.MessageBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitOneof(o *parser.Oneof) (next bool) {
	v.validateIndentLeading(o.Meta.Pos)
	defer func() { v.validateIndentLast(o.Meta.LastPos) }()
	for _, comment := range o.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, field := range o.OneofFields {
		field.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitOneofField(f *parser.OneofField) (next bool) {
	v.validateIndentLeading(f.Meta.Pos)
	for _, comment := range f.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitOption(o *parser.Option) (next bool) {
	v.validateIndentLeading(o.Meta.Pos)
	for _, comment := range o.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitPackage(p *parser.Package) (next bool) {
	v.validateIndentLeading(p.Meta.Pos)
	for _, comment := range p.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitReserved(r *parser.Reserved) (next bool) {
	v.validateIndentLeading(r.Meta.Pos)
	for _, comment := range r.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitRPC(r *parser.RPC) (next bool) {
	v.validateIndentLeading(r.Meta.Pos)
	defer func() {
		line := v.Fixer.Lines()[r.Meta.LastPos.Line-1]
		runes := []rune(line)
		for i := r.Meta.LastPos.Column - 2; 0 < i; i-- {
			r := runes[i]
			if r == '{' || r == ')' {
				// skip validating the indentation when the line ends with {}, {};, or );
				return
			}
			if r == '}' || unicode.IsSpace(r) {
				continue
			}
			break
		}
		v.validateIndentLast(r.Meta.LastPos)
	}()
	for _, comment := range r.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range r.Options {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitService(s *parser.Service) (next bool) {
	v.validateIndentLeading(s.Meta.Pos)
	defer func() { v.validateIndentLast(s.Meta.LastPos) }()
	for _, comment := range s.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}

	defer v.nest()()
	for _, body := range s.ServiceBody {
		body.Accept(v)
	}
	return false
}

func (v indentVisitor) VisitSyntax(s *parser.Syntax) (next bool) {
	v.validateIndentLeading(s.Meta.Pos)
	for _, comment := range s.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) VisitEdition(s *parser.Edition) (next bool) {
	v.validateIndentLeading(s.Meta.Pos)
	for _, comment := range s.Comments {
		v.validateIndentLeading(comment.Meta.Pos)
	}
	return false
}

func (v indentVisitor) validateIndentLeading(
	pos meta.Position,
) {
	v.validateIndent(pos, false)
}

func (v indentVisitor) validateIndentLast(
	pos meta.Position,
) {
	v.validateIndent(pos, true)
}

func (v indentVisitor) validateIndent(
	pos meta.Position,
	isLast bool,
) {
	line := v.Fixer.Lines()[pos.Line-1]
	runes := []rune(line)
	column := pos.Column
	if isLast {
		// deal with last position followed by ';'. See https://github.com/yoheimuta/protolint/issues/99
		for 1 < column && runes[column-1] == ';' {
			column--
		}
	}
	// An element is only responsible for the line's indentation when nothing but
	// whitespace precedes it. Otherwise it shares the line with an earlier element
	// or comment, and the spaces in between are not indentation. Counting every
	// space before the element as indentation reported a wrong current style and,
	// in fix mode, sliced the line at the wrong offset (see issue #409).
	isFirstOnLine := strings.TrimSpace(string(runes[:column-1])) == ""
	leading := line[:len(line)-len(strings.TrimLeftFunc(line, unicode.IsSpace))]
	// When newlines are not inserted, the whole line is re-indented as one unit,
	// so every element on it is judged by the line's indentation.
	judgedByLineIndent := isFirstOnLine || v.notInsertNewline

	indentation := strings.Repeat(v.style, v.currentLevel)
	v.indentFixes[pos.Line-1] = append(v.indentFixes[pos.Line-1], indentFix{
		currentChars: len(leading),
		replacement:  indentation,
		level:        v.currentLevel,
		pos:          pos,
		isLast:       isLast,
	})

	if judgedByLineIndent && leading == indentation {
		return
	}
	if 1 < len(v.indentFixes[pos.Line-1]) && v.notInsertNewline {
		return
	}
	// Only the element that starts the line owns its indentation, whichever of
	// the line's elements happens to be visited first. The others are told to
	// move to a line of their own (see issue #349).
	if judgedByLineIndent {
		v.AddFailuref(
			pos,
			`Found an incorrect indentation style "%s". "%s" is correct.`,
			leading,
			indentation,
		)
	} else {
		v.AddFailuref(
			pos,
			`Found a possible incorrect indentation style. Inserting a new line is recommended.`,
		)
	}
}

func (v *indentVisitor) nest() func() {
	v.currentLevel++
	return func() {
		v.currentLevel--
	}
}

func (v indentVisitor) fix(proto *parser.Proto) error {
	var shouldFixed bool

	v.Fixer.ReplaceAll(func(lines []string) []string {
		var fixedLines []string
		for i, line := range lines {
			lines := []string{line}
			if fixes, ok := v.indentFixes[i]; ok {
				// Fixes are recorded in visiting order, and a visitor validates an
				// element before the comments attached to it, so a comment that
				// precedes its element on the same line arrives after it. The split
				// below walks the line from right to left and needs the fixes in
				// column order (see issue #409).
				sort.SliceStable(fixes, func(a, b int) bool {
					return fixes[a].pos.Column < fixes[b].pos.Column
				})
				lines[0] = fixes[0].replacement + line[fixes[0].currentChars:]
				shouldFixed = true

				if 1 < len(fixes) && !v.notInsertNewline {
					// compose multiple lines in reverse order from right to left on one line.
					var rlines []string
					for j := len(fixes) - 1; 0 <= j; j-- {
						indentation := strings.Repeat(v.style, fixes[j].level)
						if fixes[j].isLast {
							// deal with last position followed by ';'. See https://github.com/yoheimuta/protolint/issues/99
							for line[fixes[j].pos.Column-1] == ';' {
								fixes[j].pos.Column--
							}
						}

						endColumn := len(line)
						if j < len(fixes)-1 {
							endColumn = fixes[j+1].pos.Column - 1
						}
						text := line[fixes[j].pos.Column-1 : endColumn]
						text = strings.TrimRightFunc(text, func(r rune) bool {
							// removing right spaces is a possible side effect that users do not expect,
							// but it's probably acceptable and usually recommended.
							return unicode.IsSpace(r)
						})

						rlines = append(rlines, indentation+text)
					}

					// sort the multiple lines in order
					lines = []string{}
					for j := len(rlines) - 1; 0 <= j; j-- {
						lines = append(lines, rlines[j])
					}
				}
			}
			fixedLines = append(fixedLines, lines...)
		}
		return fixedLines
	})

	if !shouldFixed {
		return nil
	}
	return v.BaseFixableVisitor.Finally(proto)
}
