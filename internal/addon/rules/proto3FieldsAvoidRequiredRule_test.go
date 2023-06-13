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

func TestProto3FieldsAvoidRequiredRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto without fields",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{},
				},
			},
		},
		{
			name: "no failures for proto with not required field names",
			inputProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
				},
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "song_name",
							},
							&parser.Field{
								IsRepeated: true,
								FieldName:  "singer",
							},
							&parser.Field{
								IsOptional: true,
								FieldName:  "singer",
							},
							&parser.MapField{
								MapName: "song_name2",
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "song_name3",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto with required field names for proto2",
			inputProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto2",
				},
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								IsRequired: true,
								FieldName:  "song_name",
							},
							&parser.Field{
								IsRequired: true,
								FieldName:  "singer",
							},
						},
					},
				},
			},
		},
		{
			name: "failures for proto with required field names for proto3",
			inputProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
				},
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								IsRequired: true,
								FieldName:  "song_Name",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.Field{
								IsRequired: true,
								FieldName:  "song.name",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     10,
										Column:   20,
									},
								},
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"PROTO3_FIELDS_AVOID_REQUIRED",
					`Field "song_Name" should avoid required for proto3`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"PROTO3_FIELDS_AVOID_REQUIRED",
					`Field "song.name" should avoid required for proto3`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewProto3FieldsAvoidRequiredRule(rule.SeverityError, false)

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

func TestProto3FieldsAvoidRequiredRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "avoid_required.proto",
			wantFilename:  "avoid_required.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "avoid_required.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewProto3FieldsAvoidRequiredRule(rule.SeverityError, true)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
