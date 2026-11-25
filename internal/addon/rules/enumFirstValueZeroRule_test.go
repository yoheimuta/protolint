package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/setting_test"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

func TestEnumFirstValueZeroRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for empty proto",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{},
			},
		},
		{
			name: "no failures for empty enum",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "EmptyEnum",
					},
				},
			},
		},
		{
			name: "no failures when first enum value is zero",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "Status",
						EnumBody: []parser.Visitee{
							&parser.Option{OptionName: "allow_alias", Constant: "true"},
							&parser.Reserved{Ranges: []*parser.Range{{Begin: "2"}}},
							&parser.EmptyStatement{},
							&parser.EnumField{
								Ident:  "STATUS_UNSPECIFIED",
								Number: "0",
							},
							&parser.EnumField{
								Ident:  "STATUS_ACTIVE",
								Number: "1",
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for hex and signed zero first values",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "HexZero",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "HEX_ZERO_UNSPECIFIED",
								Number: "0x0",
							},
						},
					},
					&parser.Enum{
						EnumName: "SignedZero",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "SIGNED_ZERO_UNSPECIFIED",
								Number: "-0",
							},
						},
					},
				},
			},
		},
		{
			name: "failure when first enum value is positive non-zero, even if a later value is zero",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "Status",
						EnumBody: []parser.Visitee{
							&parser.Option{OptionName: "allow_alias", Constant: "true"},
							&parser.EnumField{
								Ident:  "STATUS_ACTIVE",
								Number: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   3,
									},
								},
							},
							&parser.EnumField{
								Ident:  "STATUS_UNSPECIFIED",
								Number: "0",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   130,
										Line:     6,
										Column:   3,
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
						Column:   3,
					},
					"ENUM_FIRST_VALUE_ZERO",
					string(rule.SeverityError),
					`First value "STATUS_ACTIVE" of enum "Status" should have numeric value 0, got 1`,
				),
			},
		},
		{
			name: "failure when first enum value is negative",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "NegativeEnum",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "NEGATIVE_VALUE",
								Number: "-1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   50,
										Line:     4,
										Column:   3,
									},
								},
							},
							&parser.EnumField{
								Ident:  "ZERO_VALUE",
								Number: "0",
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   50,
						Line:     4,
						Column:   3,
					},
					"ENUM_FIRST_VALUE_ZERO",
					string(rule.SeverityError),
					`First value "NEGATIVE_VALUE" of enum "NegativeEnum" should have numeric value 0, got -1`,
				),
			},
		},
		{
			name: "failures for nested enums inside messages alongside valid enums",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageName: "Outer",
						MessageBody: []parser.Visitee{
							&parser.Enum{
								EnumName: "ValidNested",
								EnumBody: []parser.Visitee{
									&parser.EnumField{
										Ident:  "VALID_NESTED_UNSPECIFIED",
										Number: "0",
									},
								},
							},
							&parser.Enum{
								EnumName: "InvalidNested",
								EnumBody: []parser.Visitee{
									&parser.EnumField{
										Ident:  "INVALID_NESTED_FIRST",
										Number: "2",
										Meta: meta.Meta{
											Pos: meta.Position{
												Filename: "example.proto",
												Offset:   200,
												Line:     12,
												Column:   5,
											},
										},
									},
									&parser.EnumField{
										Ident:  "INVALID_NESTED_SECOND",
										Number: "3",
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
						Offset:   200,
						Line:     12,
						Column:   5,
					},
					"ENUM_FIRST_VALUE_ZERO",
					string(rule.SeverityError),
					`First value "INVALID_NESTED_FIRST" of enum "InvalidNested" should have numeric value 0, got 2`,
				),
			},
		},
		{
			name: "disable:this on first field suppresses failure without flagging the second field",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "DisabledField",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FIRST",
								Number: "1",
								InlineComment: &parser.Comment{
									Raw: "// protolint:disable:this ENUM_FIRST_VALUE_ZERO",
								},
							},
							&parser.EnumField{
								Ident:  "SECOND",
								Number: "2",
							},
						},
					},
				},
			},
		},
		{
			name: "disable:next on enum suppresses failure for its first field",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "DisabledEnum",
						Comments: []*parser.Comment{
							{Raw: "// protolint:disable:next ENUM_FIRST_VALUE_ZERO"},
						},
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FIRST",
								Number: "1",
							},
							&parser.EnumField{
								Ident:  "SECOND",
								Number: "2",
							},
						},
					},
					&parser.Enum{
						EnumName: "CheckedAfterDisabled",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "AFTER_FIRST",
								Number: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   300,
										Line:     20,
										Column:   3,
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
						Offset:   300,
						Line:     20,
						Column:   3,
					},
					"ENUM_FIRST_VALUE_ZERO",
					string(rule.SeverityError),
					`First value "AFTER_FIRST" of enum "CheckedAfterDisabled" should have numeric value 0, got 1`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewEnumFirstValueZeroRule(rule.SeverityError, autodisable.Noop)

			got, err := r.Apply(test.inputProto)
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

func TestEnumFirstValueZeroRule_Apply_disable(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		inputPlacementType autodisable.PlacementType
		wantFilename       string
	}{
		{
			name:          "do nothing in case of no violations",
			inputFilename: "valid.proto",
			wantFilename:  "valid.proto",
		},
		{
			name:               "insert disable:next comments",
			inputFilename:      "invalid.proto",
			inputPlacementType: autodisable.Next,
			wantFilename:       "disable_next.proto",
		},
		{
			name:               "insert disable:this comments",
			inputFilename:      "invalid.proto",
			inputPlacementType: autodisable.ThisThenNext,
			wantFilename:       "disable_this.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewEnumFirstValueZeroRule(rule.SeverityError, test.inputPlacementType)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}

	for _, disabledFile := range []string{"disable_next.proto", "disable_this.proto"} {
		disabledFile := disabledFile
		t.Run("re-applying on "+disabledFile+" yields no failures", func(t *testing.T) {
			path := setting_test.TestDataPath("rules", "enumFirstValueZero", disabledFile)
			proto, err := file.NewProtoFile(path, path).Parse(false)
			if err != nil {
				t.Fatalf("parse %s: %v", disabledFile, err)
			}
			r := rules.NewEnumFirstValueZeroRule(rule.SeverityError, autodisable.Noop)
			got, err := r.Apply(proto)
			if err != nil {
				t.Fatalf("apply %s: %v", disabledFile, err)
			}
			if len(got) != 0 {
				t.Errorf("got failures %v, want none", got)
			}
		})
	}
}
