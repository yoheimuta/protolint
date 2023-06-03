package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/fixer"
	"github.com/yoheimuta/protolint/linter/rule"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// PackageNameLowerCaseRule verifies that the package name doesn't contain any uppercase letters.
// See https://developers.google.com/protocol-buffers/docs/style#packages.
type PackageNameLowerCaseRule struct {
	RuleWithSeverity
	fixMode bool
}

// NewPackageNameLowerCaseRule creates a new PackageNameLowerCaseRule.
func NewPackageNameLowerCaseRule(
	severity rule.Severity,
	fixMode bool,
) PackageNameLowerCaseRule {
	return PackageNameLowerCaseRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
	}
}

// ID returns the ID of this rule.
func (r PackageNameLowerCaseRule) ID() string {
	return "PACKAGE_NAME_LOWER_CASE"
}

// Purpose returns the purpose of this rule.
func (r PackageNameLowerCaseRule) Purpose() string {
	return "Verifies that the package name doesn't contain any uppercase letters."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r PackageNameLowerCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r PackageNameLowerCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &packageNameLowerCaseVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type packageNameLowerCaseVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitPackage checks the package.
func (v *packageNameLowerCaseVisitor) VisitPackage(p *parser.Package) bool {
	name := p.Name
	if !isPackageLowerCase(name) {
		expected := strings.ToLower(name)
		v.AddFailuref(p.Meta.Pos, "Package name %q must not contain any uppercase letter. Consider to change like %q.", name, expected)

		err := v.Fixer.SearchAndReplace(p.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
			ident, startPos, _ := lex.ReadFullIdent()
			return fixer.TextEdit{
				Pos:     startPos.Offset,
				End:     startPos.Offset + len(ident) - 1,
				NewText: []byte(expected),
			}
		})
		if err != nil {
			panic(err)
		}
	}
	return false
}

func isPackageLowerCase(packageName string) bool {
	return !strs.HasAnyUpperCase(packageName)
}
