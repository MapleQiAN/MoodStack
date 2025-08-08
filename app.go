package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"MoodStack/app"
)

// App struct
type App struct {
	ctx            context.Context
	currentUser    *app.User
	currentSession *app.Session
	encryptionKey  []byte
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	if err := app.InitDatabase(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
	}

	// Create diaries directory if it doesn't exist
	if err := app.EnsureDiariesDir(); err != nil {
		fmt.Printf("Failed to create diaries directory: %v\n", err)
	}

	// Run database migrations
	if err := app.RunMigrations(); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
	}

	// Cleanup expired sessions on startup
	if err := app.CleanupExpiredSessions(); err != nil {
		fmt.Printf("Failed to cleanup expired sessions: %v\n", err)
	}
}

// GetEmotionAnalysisHistory returns emotion analysis history for current user
func (a *App) GetEmotionAnalysisHistory() ([]app.EmotionAnalysis, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	return app.GetUserEmotionTrends(a.currentUser.ID, 0)
}

// GetDiariesList returns all diary entries
func (a *App) GetDiariesList() ([]app.Diary, error) {
	if a.currentUser == nil {
		return app.GetDiariesList() // Fallback to file-based storage
	}
	return app.GetEncryptedDiariesList(a.currentUser.ID, a.encryptionKey)
}

// GetDiaryByID returns a specific diary by ID
func (a *App) GetDiaryByID(id string) (*app.Diary, error) {
	if a.currentUser == nil {
		return app.GetDiaryByID(id) // Fallback to file-based storage
	}
	return app.GetEncryptedDiaryByID(id, a.currentUser.ID, a.encryptionKey)
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
	if a.currentUser != nil {
		// Save to encrypted database
		if err := app.SaveEncryptedDiary(diary, a.currentUser.ID, a.encryptionKey); err != nil {
			return nil, fmt.Errorf("保存加密日记失败: %v", err)
		}
	} else {
		// Fallback to file-based storage
		if err := app.SaveDiary(diary); err != nil {
			return nil, fmt.Errorf("保存日记失败: %v", err)
		}
	}

	return diary, nil
}

// CreateDiaryWithEncryption creates a new diary entry with specified encryption options
func (a *App) CreateDiaryWithEncryption(title, content string, encryptionOptions app.DiaryEncryptionOptions) (*app.Diary, error) {
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

	// Save diary with encryption options
	if a.currentUser != nil {
		// Save to encrypted database with options
		if err := app.SaveEncryptedDiaryWithOptions(diary, a.currentUser.ID, a.encryptionKey, &encryptionOptions); err != nil {
			return nil, fmt.Errorf("保存加密日记失败: %v", err)
		}
	} else {
		// Fallback to file-based storage (ignoring encryption options)
		if err := app.SaveDiary(diary); err != nil {
			return nil, fmt.Errorf("保存日记失败: %v", err)
		}
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
	if a.currentUser != nil {
		// Save to encrypted database
		if err := app.SaveEncryptedDiary(diary, a.currentUser.ID, a.encryptionKey); err != nil {
			return nil, fmt.Errorf("保存加密日记失败: %v", err)
		}
	} else {
		// Fallback to file-based storage
		if err := app.SaveDiary(diary); err != nil {
			return nil, fmt.Errorf("保存日记失败: %v", err)
		}
	}

	return diary, nil
}

// UpdateDiary updates an existing diary entry
func (a *App) UpdateDiary(diary app.Diary) error {
	// Check if diary exists
	if a.currentUser != nil {
		_, err := app.GetEncryptedDiaryByID(diary.ID, a.currentUser.ID, a.encryptionKey)
		if err != nil {
			return fmt.Errorf("日记不存在: %v", err)
		}
		// Save to encrypted database
		if err := app.SaveEncryptedDiary(&diary, a.currentUser.ID, a.encryptionKey); err != nil {
			return fmt.Errorf("更新加密日记失败: %v", err)
		}
	} else {
		_, err := app.GetDiaryByID(diary.ID)
		if err != nil {
			return fmt.Errorf("日记不存在: %v", err)
		}
		// Fallback to file-based storage
		if err := app.SaveDiary(&diary); err != nil {
			return fmt.Errorf("更新日记失败: %v", err)
		}
	}

	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SearchDiaries searches for diaries by a query string
func (a *App) SearchDiaries(query string) ([]app.Diary, error) {
	if a.currentUser == nil {
		return app.SearchDiaries(query)
	}
	return app.SearchEncryptedDiaries(query, a.currentUser.ID, a.encryptionKey)
}

// SearchDiariesWithContext searches for diaries and returns detailed search results with context
func (a *App) SearchDiariesWithContext(query string) ([]app.SearchResult, error) {
	if a.currentUser == nil {
		return app.SearchDiariesWithContext(query)
	}
	return app.SearchEncryptedDiariesWithContext(query, a.currentUser.ID, a.encryptionKey)
}

// Authentication and security methods

// AuthStatusResult represents the authentication status
type AuthStatusResult struct {
	IsAuthenticated bool `json:"isAuthenticated"`
	HasUsers        bool `json:"hasUsers"`
	RequireSetup    bool `json:"requireSetup"`
}

// CheckAuthStatus checks if user is authenticated
func (a *App) CheckAuthStatus() (*AuthStatusResult, error) {
	hasUsers, err := app.HasAnyUsers()
	if err != nil {
		return nil, err
	}

	return &AuthStatusResult{
		IsAuthenticated: a.currentUser != nil,
		HasUsers:        hasUsers,
		RequireSetup:    !hasUsers,
	}, nil
}

// CreateFirstUser creates the first user in the system
func (a *App) CreateFirstUser(username, password string) (*app.AuthResult, error) {
	// Check if any users exist
	hasUsers, err := app.HasAnyUsers()
	if err != nil {
		return nil, err
	}

	if hasUsers {
		return &app.AuthResult{
			Success: false,
			Message: "系统已有用户，请使用登录功能",
		}, nil
	}

	// Create user
	_, err = app.CreateUser(username, password)
	if err != nil {
		return &app.AuthResult{
			Success: false,
			Message: fmt.Sprintf("创建用户失败: %v", err),
		}, nil
	}

	// Authenticate immediately
	return a.AuthenticateUser(username, password)
}

// AuthenticateUser authenticates with username and password
func (a *App) AuthenticateUser(username, password string) (*app.AuthResult, error) {
	result, err := app.AuthenticateWithPassword(username, password)
	if err != nil {
		return nil, err
	}

	if result.Success {
		a.currentUser = result.User
		a.currentSession = result.Session

		// Derive encryption key from password
		salt, err := base64.StdEncoding.DecodeString(result.User.Salt)
		if err != nil {
			return &app.AuthResult{
				Success: false,
				Message: "解析用户数据失败",
			}, nil
		}
		a.encryptionKey = app.DeriveKey(password, salt)
	}

	return result, nil
}

// AuthenticateWithBiometric authenticates using biometric
func (a *App) AuthenticateWithBiometric(username string) (*app.AuthResult, error) {
	result, err := app.AuthenticateWithBiometric(username)
	if err != nil {
		return nil, err
	}

	if result.Success {
		a.currentUser = result.User
		a.currentSession = result.Session

		// For biometric auth, we need to derive key from stored biometric key
		// This uses the same key derivation as password-based auth for unified encryption
		if result.User.BiometricKey != "" {
			bioKey, err := base64.StdEncoding.DecodeString(result.User.BiometricKey)
			if err == nil {
				a.encryptionKey = bioKey
			}
		} else {
			// If no biometric key is stored, use the user's salt to generate a key
			// This maintains compatibility with existing encrypted diaries
			salt, err := base64.StdEncoding.DecodeString(result.User.Salt)
			if err == nil {
				// Use a default key derivation for biometric mode
				a.encryptionKey = app.DeriveKey("biometric_default", salt)
			}
		}
	}

	return result, nil
}

// Logout logs out the current user
func (a *App) Logout() error {
	if a.currentSession != nil {
		if err := app.DeleteSession(a.currentSession.ID); err != nil {
			return err
		}
	}

	a.currentUser = nil
	a.currentSession = nil
	a.encryptionKey = nil

	return nil
}

// GetCurrentUser returns the current authenticated user
func (a *App) GetCurrentUser() *app.User {
	return a.currentUser
}

// CheckBiometricSupport checks if biometric authentication is supported
func (a *App) CheckBiometricSupport() (*app.BiometricInfo, error) {
	return app.BiometricSupport()
}

// EnableBiometric enables biometric authentication for current user
func (a *App) EnableBiometric(password string) error {
	if a.currentUser == nil {
		return fmt.Errorf("用户未登录")
	}

	return app.EnableBiometricForUser(a.currentUser.ID, password)
}

// DisableBiometric disables biometric authentication for current user
func (a *App) DisableBiometric() error {
	if a.currentUser == nil {
		return fmt.Errorf("用户未登录")
	}

	return app.DisableBiometricForUser(a.currentUser.ID)
}

// (removed) TestBiometricAPI 已移除，改用原生DLL方案

// GetFirstUser returns the first user for single-user mode
func (a *App) GetFirstUser() (*app.User, error) {
	return app.GetFirstUser()
}

// Migration methods

// CheckMigrationStatus checks if data migration is needed
func (a *App) CheckMigrationStatus() (*app.MigrationStatus, error) {
	return app.GetMigrationStatus()
}

// MigrateData migrates existing file-based diaries to encrypted database
func (a *App) MigrateData() error {
	if a.currentUser == nil {
		return fmt.Errorf("用户未登录")
	}

	if a.encryptionKey == nil {
		return fmt.Errorf("加密密钥未初始化")
	}

	return app.MigrateFileBasedDiariesToDatabase(int(a.currentUser.ID), a.encryptionKey)
}

// Flexible encryption methods

// GetDiaryWithPassword retrieves a diary using an individual password
func (a *App) GetDiaryWithPassword(diaryID, password string) (*app.Diary, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	return app.GetEncryptedDiaryWithPassword(diaryID, a.currentUser.ID, password)
}

// GetDiaryEncryptionInfo returns encryption information for a diary
func (a *App) GetDiaryEncryptionInfo(diaryID string) (*app.DiaryEncryptionInfo, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	return app.GetDiaryEncryptionInfo(diaryID, a.currentUser.ID)
}

// UpdateDiaryWithEncryption updates an existing diary with new encryption options
func (a *App) UpdateDiaryWithEncryption(diary app.Diary, encryptionOptions app.DiaryEncryptionOptions) error {
	if a.currentUser == nil {
		return fmt.Errorf("用户未登录")
	}

	// Update the diary's update time
	diary.UpdatedAt = time.Now()

	// Save with new encryption options
	return app.SaveEncryptedDiaryWithOptions(&diary, a.currentUser.ID, a.encryptionKey, &encryptionOptions)
}

// Emotion Analysis methods

// AnalyzeDiaryEmotion analyzes emotion for a specific diary
func (a *App) AnalyzeDiaryEmotion(diaryID string, useAI bool, ollamaURL string) (*app.EmotionAnalysisResult, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	// Get the diary content
	diary, err := app.GetEncryptedDiaryByID(diaryID, a.currentUser.ID, a.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("获取日记失败: %v", err)
	}

	// Analyze emotion
	return app.AnalyzeDiaryEmotion(diaryID, a.currentUser.ID, diary.Content, useAI, ollamaURL)
}

// GetDiaryEmotionAnalysis gets existing emotion analysis for a diary
func (a *App) GetDiaryEmotionAnalysis(diaryID string) (*app.EmotionAnalysis, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	return app.GetEmotionAnalysis(diaryID, a.currentUser.ID)
}

// GetUserEmotionTrends gets emotion trends for the current user
func (a *App) GetUserEmotionTrends(days int) ([]app.EmotionAnalysis, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	return app.GetUserEmotionTrends(a.currentUser.ID, days)
}

// GetUserEmotionStatistics gets aggregated emotion statistics for the current user
func (a *App) GetUserEmotionStatistics(days int) (map[string]interface{}, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	return app.GetEmotionStatistics(a.currentUser.ID, days)
}

// AnalyzeAllDiariesEmotion analyzes emotion for all user's diaries
func (a *App) AnalyzeAllDiariesEmotion(useAI bool, ollamaURL string) (map[string]interface{}, error) {
	if a.currentUser == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	// Get all diaries
	diaries, err := app.GetEncryptedDiariesList(a.currentUser.ID, a.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("获取日记列表失败: %v", err)
	}

	successCount := 0
	failureCount := 0
	var lastError error

	// Analyze each diary
	for _, diary := range diaries {
		// Skip if already analyzed (unless we're forcing re-analysis)
		existing, err := app.GetEmotionAnalysis(diary.ID, a.currentUser.ID)
		if err == nil && existing != nil {
			successCount++
			continue
		}

		// Analyze emotion
		_, err = app.AnalyzeDiaryEmotion(diary.ID, a.currentUser.ID, diary.Content, useAI, ollamaURL)
		if err != nil {
			failureCount++
			lastError = err
		} else {
			successCount++
		}
	}

	result := map[string]interface{}{
		"total":     len(diaries),
		"success":   successCount,
		"failures":  failureCount,
		"completed": successCount == len(diaries),
	}

	if lastError != nil {
		result["lastError"] = lastError.Error()
	}

	return result, nil
}
