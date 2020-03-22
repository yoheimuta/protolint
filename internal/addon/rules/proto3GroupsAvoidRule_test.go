package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestProto3GroupsAvoidRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto without groups",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{},
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
							&parser.GroupField{
								GroupName: "song_Name",
							},
							&parser.GroupField{
								GroupName: "song.name",
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
							&parser.GroupField{
								GroupName: "song_Name",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.GroupField{
								GroupName: "song.name",
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
					"PROTO3_GROUPS_AVOID",
					`Group "song_Name" should be avoided for proto3`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"PROTO3_GROUPS_AVOID",
					`Group "song.name" should be avoided for proto3`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewProto3GroupsAvoidRule()

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
