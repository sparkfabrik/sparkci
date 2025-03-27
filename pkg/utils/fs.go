package utils

import (
	"fmt"
	"os"
)

func WriteTempFile(name string, data string) (*os.File, error) {
	if name == "" {
		name = "tempfile"
	}

	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tmpFile.WriteString(data); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tmpFile, nil
}
