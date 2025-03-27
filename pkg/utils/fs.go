package utils

import (
	"fmt"
	"os"
)

func WriteTempFile(name string, data string) (*os.File, error) {
	if name == "" {
		name = "tempfile"
	}
	if data == "" {
		return nil, fmt.Errorf("empty string provided")
	}
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(data); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tmpFile, nil
}
