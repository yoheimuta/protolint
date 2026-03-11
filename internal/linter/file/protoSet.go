package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yoheimuta/protolint/internal/file"
)

// ProtoSet represents a set of .proto files.
type ProtoSet struct {
	protoFiles []*ProtoFile
}

// NewProtoSet creates a new ProtoSet.
func NewProtoSet(
	targetPaths []string,
	stdinFileName string,
) (ProtoSet, error) {
	stdinCount := 0
	for _, path := range targetPaths {
		if path == file.StdinPath {
			stdinCount++
		}
	}
	if stdinCount > 1 {
		return ProtoSet{}, fmt.Errorf("stdin (%s) can only be specified once", file.StdinPath)
	}

	fs, err := collectAllProtoFilesFromArgs(targetPaths, stdinFileName)
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
func (s *ProtoSet) ProtoFiles() []*ProtoFile {
	return s.protoFiles
}

func collectAllProtoFilesFromArgs(
	targetPaths []string,
	stdinFileName string,
) ([]*ProtoFile, error) {
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

	var fs []*ProtoFile
	for _, path := range targetPaths {
		if path == file.StdinPath {
			displayPath := file.StdinDisplayPath
			if stdinFileName != "" {
				displayPath = stdinFileName
			}
			fs = append(fs, NewProtoFile(file.StdinPath, stdinFilenameClean(absCwd, displayPath)))
			continue
		}

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
) ([]*ProtoFile, error) {
	var fs []*ProtoFile

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

// stdinFilenameClean returns a normalized path for content received via stdin.
// It attempts to make the targpath absolute and then calculates its relative path
// with respect to basepath (usually the current working directory).
//
// This normalization is crucial for the linter to correctly match the file
// against configuration-specific rules defined in .protolint.yaml, which typically
// use paths relative to the project root.
//
// If the path cannot be made relative (e.g., different drive letters on Windows),
// it returns the cleaned absolute path.
func stdinFilenameClean(basepath, targpath string) string {
	target, err := absClean(targpath)
	if err != nil {
		target = targpath
	}

	if relative, err := filepath.Rel(basepath, target); err == nil {
		return relative
	}

	return target
}
