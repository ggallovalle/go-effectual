package luagen

import (
	"go/format"
	"os"
	"path/filepath"
)

func WriteSource(info *TypeInfo, source string) (string, error) {
	formatted, err := format.Source([]byte(source))
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(info.OutputFileName())
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return "", err
	}

	return outputPath, nil
}

func WriteToDir(info *TypeInfo, source, dir string) (string, error) {
	formatted, err := format.Source([]byte(source))
	if err != nil {
		return "", err
	}

	outputPath := filepath.Join(dir, info.OutputFileName())
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return "", err
	}

	return outputPath, nil
}
