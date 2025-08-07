package app

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Diary represents a diary entry
type Diary struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	FileName  string    `json:"fileName"`
	FileType  string    `json:"fileType"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Tags      []string  `json:"tags"`
}

// SearchResult represents a search result with context
type SearchResult struct {
	Diary           Diary    `json:"diary"`
	MatchedSnippets []string `json:"matchedSnippets"`
	MatchType       string   `json:"matchType"` // "title" or "content" or "both"
}

// FileInfo represents information about an uploaded file
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
	Path string `json:"path"`
}

var diariesDir = "diaries"

// EnsureDiariesDir creates the diaries directory if it doesn't exist
func EnsureDiariesDir() error {
	return os.MkdirAll(diariesDir, 0755)
}

// GetDiariesList returns all diary entries
func GetDiariesList() ([]Diary, error) {
	var diaries []Diary

	files, err := ioutil.ReadDir(diariesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(diariesDir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}

			var diary Diary
			if err := json.Unmarshal(data, &diary); err != nil {
				continue
			}

			diaries = append(diaries, diary)
		}
	}

	return diaries, nil
}

// GetDiaryByID returns a specific diary by ID
func GetDiaryByID(id string) (*Diary, error) {
	filePath := filepath.Join(diariesDir, id+".json")

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var diary Diary
	if err := json.Unmarshal(data, &diary); err != nil {
		return nil, err
	}

	return &diary, nil
}

// SaveDiary saves a diary entry to file
func SaveDiary(diary *Diary) error {
	diary.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(diary, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(diariesDir, diary.ID+".json")
	return ioutil.WriteFile(filePath, data, 0644)
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

// SearchDiaries searches for diaries by a query string in title and content
func SearchDiaries(query string) ([]Diary, error) {
	allDiaries, err := GetDiariesList()
	if err != nil {
		return nil, err
	}

	var matchedDiaries []Diary
	lowerCaseQuery := strings.ToLower(query)

	for _, diary := range allDiaries {
		if strings.Contains(strings.ToLower(diary.Title), lowerCaseQuery) || strings.Contains(strings.ToLower(diary.Content), lowerCaseQuery) {
			matchedDiaries = append(matchedDiaries, diary)
		}
	}

	return matchedDiaries, nil
}

// SearchDiariesWithContext searches for diaries and returns detailed search results with context
func SearchDiariesWithContext(query string) ([]SearchResult, error) {
	if query == "" {
		return []SearchResult{}, nil
	}

	allDiaries, err := GetDiariesList()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	lowerCaseQuery := strings.ToLower(query)

	for _, diary := range allDiaries {
		titleMatch := strings.Contains(strings.ToLower(diary.Title), lowerCaseQuery)
		contentMatch := strings.Contains(strings.ToLower(diary.Content), lowerCaseQuery)

		if titleMatch || contentMatch {
			result := SearchResult{
				Diary:           diary,
				MatchedSnippets: []string{},
				MatchType:       "",
			}

			// 确定匹配类型
			if titleMatch && contentMatch {
				result.MatchType = "both"
			} else if titleMatch {
				result.MatchType = "title"
			} else {
				result.MatchType = "content"
			}

			// 提取匹配的内容片段
			if contentMatch {
				snippets := extractSnippets(diary.Content, query, 100) // 100字符上下文
				result.MatchedSnippets = snippets
			}

			results = append(results, result)
		}
	}

	return results, nil
}

// extractSnippets extracts text snippets around the matched query
func extractSnippets(content, query string, contextLength int) []string {
	lowerContent := strings.ToLower(content)
	lowerQuery := strings.ToLower(query)

	var snippets []string
	start := 0

	for {
		index := strings.Index(lowerContent[start:], lowerQuery)
		if index == -1 {
			break
		}

		actualIndex := start + index

		// 计算片段的开始和结束位置
		snippetStart := actualIndex - contextLength
		if snippetStart < 0 {
			snippetStart = 0
		}

		snippetEnd := actualIndex + len(query) + contextLength
		if snippetEnd > len(content) {
			snippetEnd = len(content)
		}

		snippet := content[snippetStart:snippetEnd]

		// 添加省略号
		if snippetStart > 0 {
			snippet = "..." + snippet
		}
		if snippetEnd < len(content) {
			snippet = snippet + "..."
		}

		snippets = append(snippets, snippet)

		// 移动到下一个可能的匹配位置
		start = actualIndex + len(query)

		// 限制返回的片段数量
		if len(snippets) >= 3 {
			break
		}
	}

	return snippets
}
