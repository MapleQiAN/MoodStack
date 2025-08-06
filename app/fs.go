package app

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo represents information about an uploaded file
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
	Path string `json:"path"`
}

// UploadFile saves an uploaded file and returns its information
func UploadFile(content []byte, filename string) (*FileInfo, error) {
	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create uploads directory: %v", err)
	}

	// Generate unique filename using hash
	hash := sha256.Sum256(content)
	ext := filepath.Ext(filename)
	baseFilename := strings.TrimSuffix(filename, ext)
	uniqueFilename := fmt.Sprintf("%s_%x%s", baseFilename, hash[:8], ext)

	filePath := filepath.Join(uploadsDir, uniqueFilename)

	// Write file
	if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %v", err)
	}

	fileInfo := &FileInfo{
		Name: filename,
		Size: int64(len(content)),
		Type: ext,
		Path: filePath,
	}

	return fileInfo, nil
}

// ReadFile reads a file and returns its content
func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

// GetSupportedFileTypes returns list of supported file types
func GetSupportedFileTypes() []string {
	return []string{".txt", ".md", ".docx", ".pdf"}
}

// IsFileTypeSupported checks if file type is supported
func IsFileTypeSupported(fileType string) bool {
	supportedTypes := GetSupportedFileTypes()
	for _, t := range supportedTypes {
		if strings.ToLower(fileType) == strings.ToLower(t) {
			return true
		}
	}
	return false
}
