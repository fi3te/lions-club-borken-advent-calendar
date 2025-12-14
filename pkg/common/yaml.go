package common

import (
	"os"

	"github.com/goccy/go-yaml"
)

func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func ReadYaml[T any](fileName string) (*T, error) {
	yml, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var value T
	if err := yaml.Unmarshal(yml, &value); err != nil {
		return nil, err
	}

	return &value, nil
}

func WriteYaml[T any](fileName string, value *T) error {
	bytes, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, bytes, os.ModePerm)
}
