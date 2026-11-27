package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

func newRPCForRequestResponseTest(name, request, response string) *parser.RPC {
	return &parser.RPC{
		RPCName: name,
		RPCRequest: &parser.RPCRequest{
			MessageType: request,
			Meta: meta.Meta{
				Pos: meta.Position{Offset: 10, Line: 2, Column: 3},
			},
		},
		RPCResponse: &parser.RPCResponse{
			MessageType: response,
			Meta: meta.Meta{
				Pos: meta.Position{Offset: 30, Line: 2, Column: 23},
			},
		},
	}
}

func TestRPCRequestResponseNamesEndWithRule_Apply(t *testing.T) {
	tests := []struct {
		name                string
		inputProto          *parser.Proto
		inputRequestSuffix  string
		inputResponseSuffix string
		inputExcludes       []string
		wantFailures        []report.Failure
	}{
		{
			name: "no failures for proto without rpc",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
				},
			},
		},
		{
			name: "no failures for proto with valid request and response names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							newRPCForRequestResponseTest("GetFoo", "GetFooRequest", "GetFooResponse"),
							newRPCForRequestResponseTest("ListBars", "foo.v1.ListBarsRequest", ".foo.v1.ListBarsResponse"),
						},
					},
				},
			},
		},
		{
			name: "failures for proto with invalid request and response names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							newRPCForRequestResponseTest("GetFoo", "GetFooReq", "Foo"),
							newRPCForRequestResponseTest("ListBars", "foo.v1.ListBarsRequest", "foo.v1.Bars"),
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{Offset: 10, Line: 2, Column: 3},
					"RPC_REQUEST_RESPONSE_NAMES_END_WITH",
					string(rule.SeverityError),
					`RPC "GetFoo" request message "GetFooReq" must end with Request`,
				),
				report.Failuref(
					meta.Position{Offset: 30, Line: 2, Column: 23},
					"RPC_REQUEST_RESPONSE_NAMES_END_WITH",
					string(rule.SeverityError),
					`RPC "GetFoo" response message "Foo" must end with Response`,
				),
				report.Failuref(
					meta.Position{Offset: 30, Line: 2, Column: 23},
					"RPC_REQUEST_RESPONSE_NAMES_END_WITH",
					string(rule.SeverityError),
					`RPC "ListBars" response message "foo.v1.Bars" must end with Response`,
				),
			},
		},
		{
			name: "failures for google.protobuf.Empty when it is not excluded",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							newRPCForRequestResponseTest("DeleteFoo", "DeleteFooRequest", "google.protobuf.Empty"),
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{Offset: 30, Line: 2, Column: 23},
					"RPC_REQUEST_RESPONSE_NAMES_END_WITH",
					string(rule.SeverityError),
					`RPC "DeleteFoo" response message "google.protobuf.Empty" must end with Response`,
				),
			},
		},
		{
			name: "no failures for excluded message types",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							newRPCForRequestResponseTest("DeleteFoo", "DeleteFooRequest", "google.protobuf.Empty"),
							newRPCForRequestResponseTest("Ping", ".google.protobuf.Empty", "google.protobuf.Empty"),
						},
					},
				},
			},
			inputExcludes: []string{"google.protobuf.Empty"},
		},
		{
			name: "no failures for proto with the custom suffixes",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							newRPCForRequestResponseTest("GetFoo", "GetFooReq", "GetFooRes"),
						},
					},
				},
			},
			inputRequestSuffix:  "Req",
			inputResponseSuffix: "Res",
		},
		{
			name: "failures for proto with the custom suffixes",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							newRPCForRequestResponseTest("GetFoo", "GetFooRequest", "GetFooResponse"),
						},
					},
				},
			},
			inputRequestSuffix:  "Req",
			inputResponseSuffix: "Res",
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{Offset: 10, Line: 2, Column: 3},
					"RPC_REQUEST_RESPONSE_NAMES_END_WITH",
					string(rule.SeverityError),
					`RPC "GetFoo" request message "GetFooRequest" must end with Req`,
				),
				report.Failuref(
					meta.Position{Offset: 30, Line: 2, Column: 23},
					"RPC_REQUEST_RESPONSE_NAMES_END_WITH",
					string(rule.SeverityError),
					`RPC "GetFoo" response message "GetFooResponse" must end with Res`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewRPCRequestResponseNamesEndWithRule(
				rule.SeverityError,
				test.inputRequestSuffix,
				test.inputResponseSuffix,
				test.inputExcludes,
			)

			got, err := rule.Apply(test.inputProto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if !reflect.DeepEqual(got, test.wantFailures) {
				t.Errorf("got %v, but want %v", got, test.wantFailures)
			}
		})
	}
}
