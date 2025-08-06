package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

// App struct
type App struct {
	ctx        context.Context
	diariesDir string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		diariesDir: "diaries",
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Create diaries directory if it doesn't exist
	if err := os.MkdirAll(a.diariesDir, 0755); err != nil {
		fmt.Printf("Failed to create diaries directory: %v\n", err)
	}
}

// GetDiariesList returns all diary entries
func (a *App) GetDiariesList() ([]Diary, error) {
	var diaries []Diary

	// Read all files in diaries directory
	files, err := ioutil.ReadDir(a.diariesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(a.diariesDir, file.Name())
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
func (a *App) GetDiaryByID(id string) (*Diary, error) {
	filePath := filepath.Join(a.diariesDir, id+".json")

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
func (a *App) SaveDiary(diary *Diary) error {
	diary.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(diary, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(a.diariesDir, diary.ID+".json")
	return ioutil.WriteFile(filePath, data, 0644)
}

// UploadDiary uploads a diary file and converts it to markdown
func (a *App) UploadDiary(filename string, content []byte) (*Diary, error) {
	// Check file type
	ext := strings.ToLower(filepath.Ext(filename))
	supportedTypes := []string{".txt", ".md", ".docx", ".pdf"}
	supported := false
	for _, t := range supportedTypes {
		if ext == t {
			supported = true
			break
		}
	}

	if !supported {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// Convert to markdown
	markdownContent, err := a.convertToMarkdown(content, ext)
	if err != nil {
		return nil, fmt.Errorf("转换文件失败: %v", err)
	}

	// Generate unique ID
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("生成ID失败: %v", err)
	}

	// Extract title
	title := a.extractTitle(markdownContent, filename)

	// Create diary entry
	diary := &Diary{
		ID:        id,
		Title:     title,
		Content:   markdownContent,
		FileName:  filename,
		FileType:  ext,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      []string{},
	}

	// Save diary
	if err := a.SaveDiary(diary); err != nil {
		return nil, fmt.Errorf("保存日记失败: %v", err)
	}

	return diary, nil
}

// convertToMarkdown converts various file formats to Markdown
func (a *App) convertToMarkdown(content []byte, fileType string) (string, error) {
	switch strings.ToLower(fileType) {
	case ".txt":
		return a.convertTxtToMarkdown(content), nil
	case ".md":
		return string(content), nil
	case ".docx":
		return a.convertDocxToMarkdown(content), nil
	case ".pdf":
		return a.convertPdfToMarkdown(content), nil
	default:
		return "", fmt.Errorf("不支持的文件类型: %s", fileType)
	}
}

// convertTxtToMarkdown converts plain text to markdown
func (a *App) convertTxtToMarkdown(content []byte) string {
	text := string(content)
	lines := strings.Split(text, "\n")
	var markdownLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			markdownLines = append(markdownLines, "")
			continue
		}

		// If line looks like a title (short and no period at end)
		if len(line) < 100 && !strings.HasSuffix(line, ".") && !strings.HasSuffix(line, "!") && !strings.HasSuffix(line, "?") {
			if len(markdownLines) == 0 || (len(markdownLines) > 0 && markdownLines[len(markdownLines)-1] == "") {
				markdownLines = append(markdownLines, "## "+line)
				continue
			}
		}

		markdownLines = append(markdownLines, line)
	}

	return strings.Join(markdownLines, "\n")
}

// convertDocxToMarkdown converts DOCX to markdown (simplified version)
func (a *App) convertDocxToMarkdown(content []byte) string {
	return "# 从DOCX文件导入的内容\n\n*注意：完整的DOCX解析功能将在后续版本中实现*\n\n```\n原始文件内容需要专门的DOCX解析库来处理\n```"
}

// convertPdfToMarkdown converts PDF to markdown (simplified version)
func (a *App) convertPdfToMarkdown(content []byte) string {
	return "# 从PDF文件导入的内容\n\n*注意：完整的PDF解析功能将在后续版本中实现*\n\n```\n原始文件内容需要专门的PDF解析库来处理\n```"
}

// extractTitle extracts a title from content
func (a *App) extractTitle(content string, filename string) string {
	lines := strings.Split(content, "\n")

	// Look for the first non-empty line as title
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Remove markdown headers if present
			line = strings.TrimPrefix(line, "# ")
			line = strings.TrimPrefix(line, "## ")
			line = strings.TrimPrefix(line, "### ")

			// Limit title length
			if len(line) > 50 {
				line = line[:47] + "..."
			}

			return line
		}
	}

	// If no content found, use filename
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

// generateID generates a random ID
func generateID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
