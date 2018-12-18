package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolinter/internal/linter/report"
	"github.com/yoheimuta/protolinter/internal/strs"
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

// Apply applies the rule to the proto.
func (r ServiceNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	visitor := &serviceNamesUpperCamelCaseVisitor{
		baseAddVisitor: newBaseAddVisitor(),
	}
	return runVisitor(visitor, proto)
}

type serviceNamesUpperCamelCaseVisitor struct {
	*baseAddVisitor
}

// VisitService checks the service.
func (v *serviceNamesUpperCamelCaseVisitor) VisitService(service *parser.Service) bool {
	if !strs.IsUpperCamelCase(service.ServiceName) {
		v.addFailuref(service.Meta.Pos, "Service name %q must be UpperCamelCase", service.ServiceName)
	}
	return false
}
