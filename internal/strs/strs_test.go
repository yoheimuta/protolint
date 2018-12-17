package strs_test

import (
	"testing"

	"github.com/yoheimuta/protolinter/internal/strs"
)

func TestIsUpperCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "the first letter is not an uppercase character",
			input: "hello",
		},
		{
			name:  "_ is included",
			input: "Hello_world",
		},
		{
			name:  ". is included",
			input: "Hello.world",
		},
		{
			name:  "the first letter is an uppercase character",
			input: "Hello",
			want:  true,
		},
		{
			name:  "the first letter is an uppercase character and rest is a camel case",
			input: "HelloWorld",
			want:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.IsUpperCamelCase(test.input)
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestIsUpperSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name: "empty is not uppercase",
		},
		{
			name:  "includes lowercase characters",
			input: "hello",
		},
		{
			name:  "includes a lowercase character",
			input: "hELLO",
		},
		{
			name:  "all uppercase",
			input: "HELLO",
			want:  true,
		},
		{
			name:  "all uppercase with underscore",
			input: "FIRST_VALUE",
			want:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.IsUpperSnakeCase(test.input)
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}
