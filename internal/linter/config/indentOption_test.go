package config_test

import (
	"reflect"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/yoheimuta/protolint/internal/linter/config"
)

func TestIndentOption_UnmarshalYAML(t *testing.T) {
	for _, test := range []struct {
		name             string
		inputConfig      []byte
		wantIndentOption config.IndentOption
		wantExistErr     bool
	}{
		{
			name: "not found supported style",
			inputConfig: []byte(`
style: 4-space
`),
			wantExistErr: true,
		},
		{
			name: "style: tab",
			inputConfig: []byte(`
style: tab
`),
			wantIndentOption: config.IndentOption{
				Style: "\t",
			},
		},
		{
			name: "style: 4",
			inputConfig: []byte(`
style: 4
`),
			wantIndentOption: config.IndentOption{
				Style: strings.Repeat(" ", 4),
			},
		},
		{
			name: "style: 2",
			inputConfig: []byte(`
style: 2
`),
			wantIndentOption: config.IndentOption{
				Style: strings.Repeat(" ", 2),
			},
		},
		{
			name: "not found supported newline",
			inputConfig: []byte(`
style: tab
newline: linefeed
`),
			wantExistErr: true,
		},
		{
			name: "newline: \n",
			inputConfig: []byte(`
newline: "\n"
`),
			wantIndentOption: config.IndentOption{
				Newline: "\n",
			},
		},
		{
			name: "newline: \r",
			inputConfig: []byte(`
newline: "\r"
`),
			wantIndentOption: config.IndentOption{
				Newline: "\r",
			},
		},
		{
			name: "newline: \r\n",
			inputConfig: []byte(`
newline: "\r\n"
`),
			wantIndentOption: config.IndentOption{
				Newline: "\r\n",
			},
		},
		{
			name: "support not_insert_newline",
			inputConfig: []byte(`
not_insert_newline: true
`),
			wantIndentOption: config.IndentOption{
				NotInsertNewline: true,
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var got config.IndentOption

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

			if !reflect.DeepEqual(got, test.wantIndentOption) {
				t.Errorf("got %v, but want %v", got, test.wantIndentOption)
			}
		})
	}
}
