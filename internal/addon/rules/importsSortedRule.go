package rules

import (
	"sort"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/osutil"
)

// ImportsSortedRule enforces sorted imports.
type ImportsSortedRule struct {
	newline string
	fixMode bool
}

// NewImportsSortedRule creates a new ImportsSortedRule.
func NewImportsSortedRule(
	newline string,
	fixMode bool,
) ImportsSortedRule {
	if len(newline) == 0 {
		newline = defaultNewline
	}
	return ImportsSortedRule{
		newline: newline,
		fixMode: fixMode,
	}
}

// ID returns the ID of this rule.
func (r ImportsSortedRule) ID() string {
	return "IMPORTS_SORTED"
}

// Purpose returns the purpose of this rule.
func (r ImportsSortedRule) Purpose() string {
	return "Enforces sorted imports."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r ImportsSortedRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r ImportsSortedRule) Apply(
	proto *parser.Proto,
) ([]report.Failure, error) {
	fileName := proto.Meta.Filename
	lines, err := osutil.ReadAllLines(fileName, r.newline)
	if err != nil {
		return nil, err
	}

	v := &importsSortedVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		protoLines:     lines,
		fixMode:        r.fixMode,
		newline:        r.newline,
		protoFileName:  fileName,
		sorter:         new(importSorter),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type importsSortedVisitor struct {
	*visitor.BaseAddVisitor
	protoLines []string

	fixMode       bool
	newline       string
	protoFileName string
	sorter        *importSorter
}

func (v importsSortedVisitor) VisitImport(i *parser.Import) (next bool) {
	v.sorter.add(i)
	return false
}

func (v importsSortedVisitor) Finally() error {
	notSorted := v.sorter.notSortedImports()

	var fixedLines []string
	for i, line := range v.protoLines {
		if invalid, ok := notSorted[i+1]; ok {
			v.AddFailuref(
				invalid.Meta.Pos,
				`Imports are not sorted.`,
			)
			line = v.protoLines[invalid.sortedLine-1]
		}
		fixedLines = append(fixedLines, line)
	}

	if 0 < len(notSorted) && v.fixMode {
		return osutil.WriteLinesToExistingFile(v.protoFileName, fixedLines, v.newline)
	}
	return nil
}

type notSortedImport struct {
	*parser.Import
	sortedLine int
}

type importGroup []*parser.Import

func (g importGroup) isContiguous(i *parser.Import) bool {
	last := g[len(g)-1]
	return i.Meta.Pos.Line-last.Meta.Pos.Line == 1
}

func (g importGroup) sorted() importGroup {
	var s importGroup
	s = append(s, g...)
	sort.Slice(s, func(i, j int) bool { return s[i].Location < s[j].Location })
	return s
}

func (g importGroup) notSortedImports() map[int]*notSortedImport {
	is := make(map[int]*notSortedImport)
	s := g.sorted()

	for idx, i := range g {
		sorted := s[idx]
		if i.Location != sorted.Location {
			is[i.Meta.Pos.Line] = &notSortedImport{
				Import:     i,
				sortedLine: sorted.Meta.Pos.Line,
			}
		}
	}
	return is
}

type importSorter struct {
	groups []*importGroup
}

func (s *importSorter) add(i *parser.Import) {
	for _, g := range s.groups {
		if g.isContiguous(i) {
			*g = append(*g, i)
			return
		}
	}
	s.groups = append(s.groups, &importGroup{i})
}

func (s *importSorter) notSortedImports() map[int]*notSortedImport {
	is := make(map[int]*notSortedImport)
	for _, g := range s.groups {
		ps := g.notSortedImports()
		for line, p := range ps {
			is[line] = p
		}
	}
	return is
}
