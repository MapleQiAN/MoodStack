package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"MoodStack/app"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Create diaries directory if it doesn't exist
	if err := app.EnsureDiariesDir(); err != nil {
		fmt.Printf("Failed to create diaries directory: %v\n", err)
	}
}

// GetDiariesList returns all diary entries
func (a *App) GetDiariesList() ([]app.Diary, error) {
	return app.GetDiariesList()
}

// GetDiaryByID returns a specific diary by ID
func (a *App) GetDiaryByID(id string) (*app.Diary, error) {
	return app.GetDiaryByID(id)
}

// UploadDiary uploads a diary file and converts it to markdown
func (a *App) UploadDiary(filename string, content []byte) (*app.Diary, error) {
	// Check file type
	ext := strings.ToLower(filepath.Ext(filename))
	if !app.IsFileTypeSupported(ext) {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// Convert to markdown
	markdownContent, err := app.ConvertToMarkdown(content, ext)
	if err != nil {
		return nil, fmt.Errorf("转换文件失败: %v", err)
	}

	// Generate unique ID
	id, err := app.GenerateID()
	if err != nil {
		return nil, fmt.Errorf("生成ID失败: %v", err)
	}

	// Extract title
	title := app.ExtractTitle(markdownContent, filename)

	// Create diary entry
	diary := &app.Diary{
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
	if err := app.SaveDiary(diary); err != nil {
		return nil, fmt.Errorf("保存日记失败: %v", err)
	}

	return diary, nil
}

// CreateDiary creates a new empty diary entry
func (a *App) CreateDiary(title, content string) (*app.Diary, error) {
	// Generate unique ID
	id, err := app.GenerateID()
	if err != nil {
		return nil, fmt.Errorf("生成ID失败: %v", err)
	}

	// Use provided title or generate default one
	diaryTitle := title
	if diaryTitle == "" {
		diaryTitle = "新日记"
	}

	// Create diary entry
	diary := &app.Diary{
		ID:        id,
		Title:     diaryTitle,
		Content:   content,
		FileName:  "",
		FileType:  ".md", // 标记为markdown格式
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tags:      []string{},
	}

	// Save diary
	if err := app.SaveDiary(diary); err != nil {
		return nil, fmt.Errorf("保存日记失败: %v", err)
	}

	return diary, nil
}

// UpdateDiary updates an existing diary entry
func (a *App) UpdateDiary(diary app.Diary) error {
	// Check if diary exists
	_, err := app.GetDiaryByID(diary.ID)
	if err != nil {
		return fmt.Errorf("日记不存在: %v", err)
	}

	// Save diary
	if err := app.SaveDiary(&diary); err != nil {
		return fmt.Errorf("更新日记失败: %v", err)
	}

	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
