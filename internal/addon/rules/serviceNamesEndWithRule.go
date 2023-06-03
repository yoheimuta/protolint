package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// ServiceNamesEndWithRule verifies that all service names end with the specified value.
type ServiceNamesEndWithRule struct {
	RuleWithSeverity
	text string
}

// NewServiceNamesEndWithRule creates a new ServiceNamesEndWithRule.
func NewServiceNamesEndWithRule(
	severity rule.Severity,
	text string,
) ServiceNamesEndWithRule {
	return ServiceNamesEndWithRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		text:             text,
	}
}

// ID returns the ID of this rule.
func (r ServiceNamesEndWithRule) ID() string {
	return "SERVICE_NAMES_END_WITH"
}

// Purpose returns the purpose of this rule.
func (r ServiceNamesEndWithRule) Purpose() string {
	return "Verifies that all service names end with the specified value."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r ServiceNamesEndWithRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r ServiceNamesEndWithRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &serviceNamesEndWithVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		text:           r.text,
	}

	return visitor.RunVisitor(v, proto, r.ID())
}

type serviceNamesEndWithVisitor struct {
	*visitor.BaseAddVisitor
	text string
}

// VisitService checks the service.
func (v *serviceNamesEndWithVisitor) VisitService(service *parser.Service) bool {
	if !strings.HasSuffix(service.ServiceName, v.text) {
		v.AddFailuref(service.Meta.Pos, "Service name %q must end with %s", service.ServiceName, v.text)
	}
	return false
}
