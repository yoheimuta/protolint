package config_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/protolint/internal/linter/config"
	"gopkg.in/yaml.v2"
)

func TestQuoteConsistentOption_UnmarshalYAML(t *testing.T) {
	for _, test := range []struct {
		name                      string
		inputConfig               []byte
		wantQuoteConsistentOption config.QuoteConsistentOption
		wantExistErr              bool
	}{
		{
			name: "not found supported quote",
			inputConfig: []byte(`
quote: backtick
`),
			wantExistErr: true,
		},
		{
			name: "empty config",
			inputConfig: []byte(`
`),
		},
		{
			name: "quote: double",
			inputConfig: []byte(`
quote: double
`),
			wantQuoteConsistentOption: config.QuoteConsistentOption{
				Quote: config.DoubleQuote,
			},
		},
		{
			name: "quote: single",
			inputConfig: []byte(`
quote: single
`),
			wantQuoteConsistentOption: config.QuoteConsistentOption{
				Quote: config.SingleQuote,
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var got config.QuoteConsistentOption

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

			if !reflect.DeepEqual(got, test.wantQuoteConsistentOption) {
				t.Errorf("got %v, but want %v", got, test.wantQuoteConsistentOption)
			}
		})
	}
}
