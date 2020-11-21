package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProtoSet represents a set of .proto files.
type ProtoSet struct {
	protoFiles []ProtoFile
}

// NewProtoSet creates a new ProtoSet.
func NewProtoSet(
	targetPaths []string,
) (ProtoSet, error) {
	fs, err := collectAllProtoFilesFromArgs(targetPaths)
	if err != nil {
		return ProtoSet{}, err
	}
	if len(fs) == 0 {
		return ProtoSet{}, fmt.Errorf("not found protocol buffer files in %v", targetPaths)
	}

	return ProtoSet{
		protoFiles: fs,
	}, nil
}

// ProtoFiles returns proto files.
func (s ProtoSet) ProtoFiles() []ProtoFile {
	return s.protoFiles
}

func collectAllProtoFilesFromArgs(
	targetPaths []string,
) ([]ProtoFile, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	absCwd, err := absClean(cwd)
	if err != nil {
		return nil, err
	}
	// Eval a possible symlink for the cwd to calculate the correct relative paths in the next step.
	if newPath, err := filepath.EvalSymlinks(absCwd); err == nil {
		absCwd = newPath
	}

	var fs []ProtoFile
	for _, path := range targetPaths {
		absTarget, err := absClean(path)
		if err != nil {
			return nil, err
		}

		f, err := collectAllProtoFiles(absCwd, absTarget)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f...)
	}
	return fs, nil
}

func collectAllProtoFiles(
	absWorkDirPath string,
	absPath string,
) ([]ProtoFile, error) {
	var fs []ProtoFile

	err := filepath.Walk(
		absPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) != ".proto" {
				return nil
			}

			displayPath, err := filepath.Rel(absWorkDirPath, path)
			if err != nil {
				displayPath = path
			}
			displayPath = filepath.Clean(displayPath)
			fs = append(fs, NewProtoFile(path, displayPath))
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

// absClean returns the cleaned absolute path of the given path.
func absClean(path string) (string, error) {
	if path == "" {
		return path, nil
	}
	if !filepath.IsAbs(path) {
		return filepath.Abs(path)
	}
	return filepath.Clean(path), nil
}
