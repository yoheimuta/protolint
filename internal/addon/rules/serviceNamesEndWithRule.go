package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

// ServiceNamesEndWithRule verifies that all service names end with the specified value.
type ServiceNamesEndWithRule struct {
	text string
}

// NewServiceNamesEndWithRule creates a new ServiceNamesEndWithRule.
func NewServiceNamesEndWithRule(text string) ServiceNamesEndWithRule {
	return ServiceNamesEndWithRule{
		text: text,
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

// Apply applies the rule to the proto.
func (r ServiceNamesEndWithRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &serviceNamesEndWithVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(),
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
