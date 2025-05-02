package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func SaveDocumentFile(originalFilename string, data []byte) (string, error) {
	
	baseDir := "uploads"
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create storage dir: %w", err)
	}

	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(originalFilename)
	name := fmt.Sprintf("%d%s", timestamp, ext)
	fullPath := filepath.Join(baseDir, name)

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fullPath, nil
}
func ReadDocumentFile(filePath string) ([]byte, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }
    return data, nil
}