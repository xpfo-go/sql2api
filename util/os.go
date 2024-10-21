package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFileIfNotExist(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	dirPath := filepath.Dir(filePath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("create dir failed: %w", err)
	}

	if _, err := os.Create(filePath); err != nil {
		return fmt.Errorf("create file failed: %w", err)
	}

	return nil
}
