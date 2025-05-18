package rules

import (
	"sort"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// ImportsSortedRule enforces sorted imports.
type ImportsSortedRule struct {
	RuleWithSeverity
	fixMode bool
}

// NewImportsSortedRule creates a new ImportsSortedRule.
func NewImportsSortedRule(
	severity rule.Severity,
	fixMode bool,
) ImportsSortedRule {
	return ImportsSortedRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
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
	base, err := visitor.NewBaseFixableVisitor(r.ID(), true, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &importsSortedVisitor{
		BaseFixableVisitor: base,
		fixMode:            r.fixMode,
		sorter:             new(importSorter),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type importsSortedVisitor struct {
	*visitor.BaseFixableVisitor
	fixMode bool
	sorter  *importSorter
}

func (v importsSortedVisitor) VisitImport(i *parser.Import) (next bool) {
	v.sorter.add(i)
	return false
}

func (v importsSortedVisitor) Finally(proto *parser.Proto) error {
	notSorted := v.sorter.notSortedImports()

	v.Fixer.ReplaceAll(func(lines []string) []string {
		var fixedLines []string

		for i, line := range lines {
			if invalid, ok := notSorted[i+1]; ok {
				v.AddFailuref(
					invalid.Meta.Pos,
					`Imports are not sorted.`,
				)
				line = lines[invalid.sortedLine-1]
			}
			fixedLines = append(fixedLines, line)
		}
		return fixedLines
	})
	if !v.fixMode {
		return nil
	}
	return v.BaseFixableVisitor.Finally(proto)
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
