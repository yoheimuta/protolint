package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

func TestValidServiceNamesEndWithRule_Apply(t *testing.T) {
	validTestCase := struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		name: "no failures for proto with valid service names",
		inputProto: &parser.Proto{
			ProtoBody: []parser.Visitee{
				&parser.Service{
					ServiceName: "SomeServiceService",
				},
				&parser.Service{
					ServiceName: "AnotherService",
				},
			},
		},
	}

	t.Run(validTestCase.name, func(t *testing.T) {
		rule := rules.NewServiceNamesEndWithRule("")

		_, err := rule.Apply(validTestCase.inputProto)
		if err != nil {
			t.Errorf("got err %v, but want nil", err)
			return
		}
	})
}

func TestInvalidServiceNamesEndWithRule_Apply(t *testing.T) {
	invalidTestCase := struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		name: "failures for proto with invalid service names",
		inputProto: &parser.Proto{
			ProtoBody: []parser.Visitee{
				&parser.Service{
					ServiceName: "SomeThing",
				},
				&parser.Service{
					ServiceName: "AnotherThing",
				},
			},
		},
		wantFailures: []report.Failure{
			report.Failuref(meta.Position{}, `Service name "SomeThing" must end with Service`),
			report.Failuref(meta.Position{}, `Service name "AnotherThing" must end with Service`),
		},
	}

	t.Run(invalidTestCase.name, func(t *testing.T) {
		rule := rules.NewServiceNamesEndWithRule("Service")

		got, err := rule.Apply(invalidTestCase.inputProto)
		if err != nil {
			t.Errorf("got err %v, but want nil", err)
			return
		}
		if !reflect.DeepEqual(got, invalidTestCase.wantFailures) {
			t.Errorf("got %v, but want %v", got, invalidTestCase.wantFailures)
		}
	})
}
