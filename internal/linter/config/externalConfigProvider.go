package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yoheimuta/protolint/internal/stringsutil"
)

type configLoader interface {
	LoadExternalConfig() (*ExternalConfig, error)
}

func loadFileContent(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("read %s, but the content is empty", file)
	}

	return data, nil
}

// GetExternalConfig provides the externalConfig.
func GetExternalConfig(
	filePath string,
	dirPath string,
) (*ExternalConfig, error) {
	reader, err := getExternalConfigLoader(filePath, dirPath)
	if err != nil {
		if len(filePath) == 0 && len(dirPath) == 0 {
			return nil, nil
		}
		return nil, err
	}

	return reader.LoadExternalConfig()
}

func getLoaderFromExtension(filePath string) (configLoader, error) {
	if strings.HasSuffix(filePath, externalConfigFileExtension) || strings.HasSuffix(filePath, externalConfigFileExtension2) {
		return yamlConfigLoader{filePath: filePath}, nil
	}
	if strings.HasSuffix(filePath, packageJsonFileNameForJsExtension) {
		return jsonConfigLoader{filePath: filePath}, nil
	}
	if strings.HasSuffix(filePath, pyProjectTomlFileNameForPyExtension) {
		return tomlConfigLoader{filePath: filePath}, nil
	}

	return nil, fmt.Errorf("%s is not a valid support file extension", filePath)
}

func getExternalConfigLoader(
	filePath string,
	dirPath string,
) (configLoader, error) {
	if 0 < len(filePath) {
		return getLoaderFromExtension(filePath)
	}

	dirPaths := []string{dirPath}
	if len(dirPath) == 0 {
		absPath, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		absPath = filepath.Dir(absPath)
		for !stringsutil.ContainsStringInSlice(absPath, dirPaths) {
			dirPaths = append(dirPaths, absPath)
			absPath = filepath.Dir(absPath)
		}
	}

	var checkedPaths []string
	// use protolint native files for default
	for _, dir := range dirPaths {
		for _, name := range []string{
			externalConfigFileName,
			externalConfigFileName2,
		} {
			for _, ext := range []string{
				externalConfigFileExtension,
				externalConfigFileExtension2,
			} {
				filePath := filepath.Join(dir, name+ext)
				checkedPaths = append(checkedPaths, filePath)
				if _, err := os.Stat(filePath); err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return nil, err
				}
				return yamlConfigLoader{filePath: filePath}, nil
			}
		}
	}

	// after checking for protolint yaml files, go for package.json of npm
	for _, dir := range dirPaths {
		filePath := filepath.Join(dir, packageJsonFileNameForJs)
		checkedPaths = append(checkedPaths, filePath)
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		return jsonConfigLoader{filePath: filePath}, nil
	}

	// after checking for protolint yaml and npm.json files, go for pyproject.toml of python
	for _, dir := range dirPaths {
		filePath := filepath.Join(dir, pyProjectTomlFileNameForPy)
		checkedPaths = append(checkedPaths, filePath)
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		return tomlConfigLoader{filePath: filePath}, nil
	}

	return nil, fmt.Errorf("not found config file by searching `%s`", strings.Join(checkedPaths, ","))
}
