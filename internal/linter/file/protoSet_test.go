package file_test

import (
	"testing"

	"path/filepath"

	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/setting_test"
)

func TestNewProtoSet(t *testing.T) {
	tests := []struct {
		name             string
		inputTargetPaths []string
		wantProtoFiles   []file.ProtoFile
		wantExistErr     bool
	}{
		{
			name: "innerdir3 includes no files",
			inputTargetPaths: []string{
				setting_test.TestDataPath("testdir", "innerdir3"),
			},
			wantExistErr: true,
		},
		{
			name: "innerdir2 includes no proto files",
			inputTargetPaths: []string{
				setting_test.TestDataPath("testdir", "innerdir2"),
			},
			wantExistErr: true,
		},
		{
			name: "innerdir includes a proto file",
			inputTargetPaths: []string{
				setting_test.TestDataPath("testdir", "innerdir"),
			},
			wantProtoFiles: []file.ProtoFile{
				file.NewProtoFile(
					filepath.Join(setting_test.TestDataPath("testdir", "innerdir"), "/testinner.proto"),
					"../../../_testdata/testdir/innerdir/testinner.proto",
				),
			},
		},
		{
			name: "testdir includes proto files and inner dirs",
			inputTargetPaths: []string{
				setting_test.TestDataPath("testdir"),
			},
			wantProtoFiles: []file.ProtoFile{
				file.NewProtoFile(
					filepath.Join(setting_test.TestDataPath("testdir", "innerdir"), "/testinner.proto"),
					"../../../_testdata/testdir/innerdir/testinner.proto",
				),
				file.NewProtoFile(
					filepath.Join(setting_test.TestDataPath("testdir"), "/test.proto"),
					"../../../_testdata/testdir/test.proto",
				),
				file.NewProtoFile(
					filepath.Join(setting_test.TestDataPath("testdir"), "/test2.proto"),
					"../../../_testdata/testdir/test2.proto",
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := file.NewProtoSet(test.inputTargetPaths)
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

			for i, gotf := range got.ProtoFiles() {
				wantf := test.wantProtoFiles[i]
				if gotf.Path() != wantf.Path() {
					t.Errorf("got %v, but want %v", gotf.Path(), wantf.Path())
				}
				if gotf.DisplayPath() != wantf.DisplayPath() {
					t.Errorf("got %v, but want %v", gotf.DisplayPath(), wantf.DisplayPath())
				}
			}
		})
	}
}
