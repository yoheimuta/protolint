package rules

import (
	"bytes"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
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
type OrderRule struct {
	fixMode bool
}

// NewOrderRule creates a new OrderRule.
func NewOrderRule(
	fixMode bool,
) OrderRule {
	return OrderRule{
		fixMode: fixMode,
	}
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
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &orderVisitor{
		BaseFixableVisitor: base,
		state:              initialOrderState,
		machine:            newOrderStateTransition(),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type orderVisitor struct {
	*visitor.BaseFixableVisitor
	state   orderState
	machine orderStateTransition

	formatter formatter
}

func (v *orderVisitor) Finally() error {
	if 0 < len(v.Failures()) {
		shouldFixed := true
		v.Fixer.ReplaceContent(func(content []byte) []byte {
			newContent := v.formatter.format(content)
			if bytes.Equal(content, newContent) {
				shouldFixed = false
			}
			return newContent
		})

		// TODO: BaseFixableVisitor.Finally should run the base Finally first, and then the fixing later.
		if shouldFixed {
			return v.BaseFixableVisitor.Finally()
		}
	}
	return nil
}

func (v *orderVisitor) VisitSyntax(s *parser.Syntax) bool {
	next := v.machine.transit(v.state, syntaxVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(s.Meta.Pos, "Syntax should be located at the top. Check if the file is ordered in the correct manner.")
	}
	v.state = syntaxOrderState
	v.formatter.syntax = s
	return false
}

func (v *orderVisitor) VisitPackage(p *parser.Package) bool {
	next := v.machine.transit(v.state, packageVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(p.Meta.Pos, "The order of Package is invalid. Check if the file is ordered in the correct manner.")
	}
	v.state = packageOrderState
	v.formatter.pkg = p
	return false
}

func (v *orderVisitor) VisitImport(i *parser.Import) bool {
	next := v.machine.transit(v.state, importsVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(i.Meta.Pos, "The order of Import is invalid. Check if the file is ordered in the correct manner.")
	}
	v.state = importsOrderState
	v.formatter.addImports(i)
	return false
}

func (v *orderVisitor) VisitOption(o *parser.Option) bool {
	next := v.machine.transit(v.state, fileOptionsVisitEvent)
	if next == invalidOrderState {
		v.AddFailuref(o.Meta.Pos, "The order of Option is invalid. Check if the file is ordered in the correct manner.")
	}
	v.state = fileOptionsOrderState
	v.formatter.addOptions(o)
	return false
}

func (v *orderVisitor) VisitMessage(m *parser.Message) bool {
	v.state = everythingElseOrderState
	v.formatter.addMisc(m)
	return false
}

func (v *orderVisitor) VisitEnum(e *parser.Enum) bool {
	v.state = everythingElseOrderState
	v.formatter.addMisc(e)
	return false
}

func (v *orderVisitor) VisitService(s *parser.Service) bool {
	v.state = everythingElseOrderState
	v.formatter.addMisc(s)
	return false
}

func (v *orderVisitor) VisitExtend(e *parser.Extend) bool {
	v.state = everythingElseOrderState
	v.formatter.addMisc(e)
	return false
}

// State Checker
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

// Formatter
type indexedVisitee struct {
	index   int
	visitee parser.Visitee
}

// NOTE: This check is not used at the moment.
// If no one requests to put the same element in a row as much as possible,
// we should delete this wrap struct, indexedVisitee.
func (i indexedVisitee) isContiguous(a indexedVisitee) bool {
	return i.index-a.index == 1
}

type formatter struct {
	syntax  *parser.Syntax
	pkg     *parser.Package
	imports []indexedVisitee
	options []indexedVisitee
	misc    []indexedVisitee
}

func (f formatter) index() int {
	idx := 0
	if f.syntax != nil {
		idx = 1
	}
	if f.pkg != nil {
		idx++
	}
	return idx + len(f.imports) + len(f.options) + len(f.misc)
}

func (f *formatter) addImports(t *parser.Import) {
	f.imports = append(f.imports, indexedVisitee{f.index(), t})
}

func (f *formatter) addOptions(t *parser.Option) {
	f.options = append(f.options, indexedVisitee{f.index(), t})
}

func (f *formatter) addMisc(t parser.Visitee) {
	f.misc = append(f.misc, indexedVisitee{f.index(), t})
}

type line struct {
	startPos meta.Position
	endPos   meta.Position
}

func newLine(meta meta.Meta, comments []*parser.Comment, inline *parser.Comment) line {
	var l line
	l.startPos = meta.Pos
	if 0 < len(comments) {
		l.startPos = comments[0].Meta.Pos
	}
	l.endPos = meta.LastPos
	if inline != nil {
		l.endPos = inline.Meta.LastPos
	}
	return l
}

func newVisiteeLine(elm parser.Visitee) line {
	switch e := elm.(type) {
	case *parser.Syntax:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	case *parser.Package:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	case *parser.Import:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	case *parser.Option:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	case *parser.Message:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	case *parser.Extend:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	case *parser.Enum:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	case *parser.Service:
		return newLine(e.Meta, e.Comments, e.InlineComment)
	}
	return line{}
}

func (l line) hasEmptyLine(prev line) bool {
	return 1 < l.startPos.Line-prev.endPos.Line
}

type writer struct {
	content    []byte
	newContent []byte
}

func (w *writer) write(l line) {
	w.newContent = append(w.newContent, w.content[l.startPos.Offset:l.endPos.Offset+1]...)
}

func (w *writer) writeN(l line) {
	w.write(l)
	w.newContent = append(w.newContent, "\n"...)
}

func (w *writer) writeNN(l line) {
	w.write(l)
	w.newContent = append(w.newContent, "\n\n"...)
}

func (w *writer) writeOnlyN() {
	w.newContent = append(w.newContent, "\n"...)
}

func (w *writer) removeLastRedundantN() {
	if bytes.Equal(w.newContent[len(w.newContent)-2:len(w.newContent)], []byte("\n\n")) {
		w.newContent = w.newContent[0 : len(w.newContent)-1]
	}
}

func (f formatter) format(content []byte) []byte {
	w := writer{content: content}

	sl := newVisiteeLine(f.syntax)
	w.writeNN(sl)

	if f.pkg != nil {
		pl := newVisiteeLine(f.pkg)
		w.writeNN(pl)
	}

	visitees := [][]indexedVisitee{f.imports, f.options, f.misc}
	for i, vs := range visitees {
		var ls []line
		for _, elm := range vs {
			ls = append(ls, newVisiteeLine(elm.visitee))
		}

		for i, l := range ls {
			// There are any empty lines between both ls
			if 0 < i && l.hasEmptyLine(ls[i-1]) {
				w.writeOnlyN()
			}
			w.writeN(l)
		}

		if 0 < len(ls) && i < len(visitees)-1 {
			w.writeOnlyN()
		}
	}

	w.removeLastRedundantN()

	return w.newContent
}
