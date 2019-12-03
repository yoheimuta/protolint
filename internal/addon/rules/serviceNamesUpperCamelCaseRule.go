package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// ServiceNamesUpperCamelCaseRule verifies that all service names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#services.
type ServiceNamesUpperCamelCaseRule struct{}

// NewServiceNamesUpperCamelCaseRule creates a new ServiceNamesUpperCamelCaseRule.
func NewServiceNamesUpperCamelCaseRule() ServiceNamesUpperCamelCaseRule {
	return ServiceNamesUpperCamelCaseRule{}
}

// ID returns the ID of this rule.
func (r ServiceNamesUpperCamelCaseRule) ID() string {
	return "SERVICE_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r ServiceNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all service names are CamelCase (with an initial capital)."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r ServiceNamesUpperCamelCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r ServiceNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &serviceNamesUpperCamelCaseVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type serviceNamesUpperCamelCaseVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitService checks the service.
func (v *serviceNamesUpperCamelCaseVisitor) VisitService(service *parser.Service) bool {
	if !strs.IsUpperCamelCase(service.ServiceName) {
		v.AddFailuref(service.Meta.Pos, "Service name %q must be UpperCamelCase", service.ServiceName)
	}
	return false
}
