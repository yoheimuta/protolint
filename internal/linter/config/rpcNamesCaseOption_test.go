package config_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/protolint/internal/linter/config"
	"gopkg.in/yaml.v2"
)

func TestRPCNamesCaseOption_UnmarshalYAML(t *testing.T) {
	for _, test := range []struct {
		name                   string
		inputConfig            []byte
		wantRPCNamesCaseOption config.RPCNamesCaseOption
		wantExistErr           bool
	}{
		{
			name: "not found supported convention",
			inputConfig: []byte(`
convention: upper_camel_case
`),
			wantExistErr: true,
		},
		{
			name: "empty config",
			inputConfig: []byte(`
`),
		},
		{
			name: "convention: lower_camel_case",
			inputConfig: []byte(`
convention: lower_camel_case
`),
			wantRPCNamesCaseOption: config.RPCNamesCaseOption{
				Convention: config.ConventionLowerCamel,
			},
		},
		{
			name: "convention: upper_snake_case",
			inputConfig: []byte(`
convention: upper_snake_case
`),
			wantRPCNamesCaseOption: config.RPCNamesCaseOption{
				Convention: config.ConventionUpperSnake,
			},
		},
		{
			name: "convention: lower_snake_case",
			inputConfig: []byte(`
convention: lower_snake_case
`),
			wantRPCNamesCaseOption: config.RPCNamesCaseOption{
				Convention: config.ConventionLowerSnake,
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var got config.RPCNamesCaseOption

			err := yaml.UnmarshalStrict(test.inputConfig, &got)
			if test.wantExistErr {
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			}
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantRPCNamesCaseOption) {
				t.Errorf("got %v, but want %v", got, test.wantRPCNamesCaseOption)
			}
		})
	}
}
