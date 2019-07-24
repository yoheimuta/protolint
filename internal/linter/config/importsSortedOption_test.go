package config_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/protolint/internal/linter/config"
	"gopkg.in/yaml.v2"
)

func TestImportsSortedOption_UnmarshalYAML(t *testing.T) {
	for _, test := range []struct {
		name                    string
		inputConfig             []byte
		wantImportsSortedOption config.ImportsSortedOption
		wantExistErr            bool
	}{
		{
			name: "not found supported newline",
			inputConfig: []byte(`
newline: linefeed
`),
			wantExistErr: true,
		},
		{
			name: "newline: \n",
			inputConfig: []byte(`
newline: "\n"
`),
			wantImportsSortedOption: config.ImportsSortedOption{
				Newline: "\n",
			},
		},
		{
			name: "newline: \r",
			inputConfig: []byte(`
newline: "\r"
`),
			wantImportsSortedOption: config.ImportsSortedOption{
				Newline: "\r",
			},
		},
		{
			name: "newline: \r\n",
			inputConfig: []byte(`
newline: "\r\n"
`),
			wantImportsSortedOption: config.ImportsSortedOption{
				Newline: "\r\n",
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var got config.ImportsSortedOption

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

			if !reflect.DeepEqual(got, test.wantImportsSortedOption) {
				t.Errorf("got %v, but want %v", got, test.wantImportsSortedOption)
			}
		})
	}
}
