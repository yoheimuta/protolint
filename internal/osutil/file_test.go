package osutil_test

import (
	"testing"

	"github.com/yoheimuta/protolint/internal/osutil"
)

func TestDetectLineEnding(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         string
		wantExistErr bool
	}{
		{
			name:  "An empty string has no line ending",
			input: ``,
		},
		{
			name:  "a line string has no line ending",
			input: `first line`,
		},
		{
			name: "two lines have \n",
			input: `first line
second line`,
			want: "\n",
		},
		{
			name:  "two lines have \r",
			input: `first line` + "\r" + `second line`,
			want:  "\r",
		},
		{
			name:  "two lines have \r\n",
			input: `first line` + "\r\n" + `second line`,
			want:  "\r\n",
		},
		{
			name: "two lines have a mix of \n and \r, and \n is more",
			input: `first line
second line
third line` + "\r" + `forth line`,
			want: "\n",
		},
		{
			name: "two lines have a mix of \n and \r, and \r is more",
			input: `first line
second line` + "\r" + `third line` + "\r" + `forth line`,

			want: "\r",
		},
		{
			name: "two lines have a mix of \r\n and \r, and \r\n is more",
			input: `first line` + "\r\n" + `second line` + "\r\n" + `third line` +
				"\r" + `forth line`,
			want: "\r\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := osutil.DetectLineEnding(test.input)
			if test.wantExistErr && err == nil {
				t.Errorf("got err nil, but want err")
				return
			} else if err != nil {
				t.Errorf("got err, but want nil")
				return
			}
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}
