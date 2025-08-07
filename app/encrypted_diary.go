package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// DiaryEncryptionOptions represents encryption options for a diary
type DiaryEncryptionOptions struct {
	Mode               string `json:"mode"`                         // 'unified', 'individual', 'biometric'
	IndividualPassword string `json:"individualPassword,omitempty"` // Only for individual mode
}

// SaveEncryptedDiary saves a diary entry with unified encryption
func SaveEncryptedDiary(diary *Diary, userID uint, encryptionKey []byte) error {
	return SaveEncryptedDiaryWithOptions(diary, userID, encryptionKey, &DiaryEncryptionOptions{
		Mode: "unified",
	})
}

// SaveEncryptedDiaryWithOptions saves a diary entry with specified encryption options
func SaveEncryptedDiaryWithOptions(diary *Diary, userID uint, masterKey []byte, options *DiaryEncryptionOptions) error {
	var encryptionKey []byte
	var encryptionSalt string

	// Determine encryption key based on mode
	switch options.Mode {
	case "unified":
		// Use the master key (derived from user password)
		encryptionKey = masterKey

	case "individual":
		// Generate a new salt for this diary's individual password
		if options.IndividualPassword == "" {
			return fmt.Errorf("individual password is required for individual encryption mode")
		}

		salt := make([]byte, saltSize)
		if _, err := rand.Read(salt); err != nil {
			return fmt.Errorf("failed to generate salt: %v", err)
		}

		encryptionKey = DeriveKey(options.IndividualPassword, salt)
		encryptionSalt = base64.StdEncoding.EncodeToString(salt)

	case "biometric":
		// Use the master key but mark as biometric mode
		encryptionKey = masterKey

	default:
		return fmt.Errorf("unsupported encryption mode: %s", options.Mode)
	}

	// Encrypt content
	encryptedContent, iv, err := EncryptData([]byte(diary.Content), encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt diary content: %v", err)
	}

	// Convert tags to JSON string
	tagsJSON, err := json.Marshal(diary.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %v", err)
	}

	ivStr := base64.StdEncoding.EncodeToString(iv)

	// Create encrypted diary entry
	encDiary := &EncryptedDiary{
		ID:               diary.ID,
		UserID:           userID,
		Title:            diary.Title,
		EncryptedContent: encryptedContent,
		IV:               ivStr,
		FileName:         diary.FileName,
		FileType:         diary.FileType,
		Tags:             string(tagsJSON),
		EncryptionMode:   options.Mode,
		EncryptionSalt:   encryptionSalt,
		CreatedAt:        diary.CreatedAt,
		UpdatedAt:        diary.UpdatedAt,
	}

	// Use GORM's Save method which handles both create and update
	if err := gormDB.Save(encDiary).Error; err != nil {
		return fmt.Errorf("failed to save encrypted diary: %v", err)
	}

	return nil
}

// GetEncryptedDiariesList returns all diary entries for a user (decrypted)
func GetEncryptedDiariesList(userID uint, encryptionKey []byte) ([]Diary, error) {
	var encDiaries []EncryptedDiary
	if err := gormDB.Where("user_id = ?", userID).Order("created_at DESC").Find(&encDiaries).Error; err != nil {
		return nil, fmt.Errorf("failed to query diaries: %v", err)
	}

	var diaries []Diary
	for _, encDiary := range encDiaries {
		// Parse tags JSON
		tagsJSON := encDiary.Tags

		// For 'unified' and 'biometric' modes, we can decrypt with the provided key
		// For 'individual' mode, we can't decrypt without the specific password
		if encDiary.EncryptionMode == "individual" {
			// Create a placeholder diary that shows it needs a password
			diary := Diary{
				ID:        encDiary.ID,
				Title:     encDiary.Title,
				Content:   "[此日记需要单独密码解锁]",
				FileName:  encDiary.FileName,
				FileType:  encDiary.FileType,
				Tags:      []string{},
				CreatedAt: encDiary.CreatedAt,
				UpdatedAt: encDiary.UpdatedAt,
			}
			diaries = append(diaries, diary)
			continue
		}

		// Decrypt content for unified and biometric modes
		content, err := decryptDiaryContent(encDiary.EncryptedContent, encDiary.IV, encryptionKey)
		if err != nil {
			continue
		}

		// Parse tags
		var tags []string
		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
				tags = []string{}
			}
		}

		diary := Diary{
			ID:        encDiary.ID,
			Title:     encDiary.Title,
			Content:   content,
			FileName:  encDiary.FileName,
			FileType:  encDiary.FileType,
			Tags:      tags,
			CreatedAt: encDiary.CreatedAt,
			UpdatedAt: encDiary.UpdatedAt,
		}

		diaries = append(diaries, diary)
	}

	return diaries, nil
}

// GetEncryptedDiaryWithPassword returns a diary decrypted with an individual password
func GetEncryptedDiaryWithPassword(diaryID string, userID uint, password string) (*Diary, error) {
	var encDiary EncryptedDiary
	if err := gormDB.Where("id = ? AND user_id = ?", diaryID, userID).First(&encDiary).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diary not found")
		}
		return nil, fmt.Errorf("failed to get diary: %v", err)
	}

	// Verify this is an individually encrypted diary
	if encDiary.EncryptionMode != "individual" {
		return nil, fmt.Errorf("this diary is not individually encrypted")
	}

	// Derive the encryption key from the password and stored salt
	salt, err := base64.StdEncoding.DecodeString(encDiary.EncryptionSalt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption salt: %v", err)
	}

	encryptionKey := DeriveKey(password, salt)

	// Decrypt content
	content, err := decryptDiaryContent(encDiary.EncryptedContent, encDiary.IV, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt diary content (incorrect password?): %v", err)
	}

	// Parse tags
	var tags []string
	if encDiary.Tags != "" {
		if err := json.Unmarshal([]byte(encDiary.Tags), &tags); err != nil {
			tags = []string{}
		}
	}

	diary := &Diary{
		ID:        encDiary.ID,
		Title:     encDiary.Title,
		Content:   content,
		FileName:  encDiary.FileName,
		FileType:  encDiary.FileType,
		Tags:      tags,
		CreatedAt: encDiary.CreatedAt,
		UpdatedAt: encDiary.UpdatedAt,
	}

	return diary, nil
}

// GetEncryptedDiaryByID returns a specific diary by ID (decrypted)
func GetEncryptedDiaryByID(diaryID string, userID uint, encryptionKey []byte) (*Diary, error) {
	var encDiary EncryptedDiary
	if err := gormDB.Where("id = ? AND user_id = ?", diaryID, userID).First(&encDiary).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diary not found")
		}
		return nil, fmt.Errorf("failed to get diary: %v", err)
	}

	// For individually encrypted diaries, we can't decrypt without the specific password
	if encDiary.EncryptionMode == "individual" {
		return &Diary{
			ID:        encDiary.ID,
			Title:     encDiary.Title,
			Content:   "[此日记需要单独密码解锁]",
			FileName:  encDiary.FileName,
			FileType:  encDiary.FileType,
			Tags:      []string{},
			CreatedAt: encDiary.CreatedAt,
			UpdatedAt: encDiary.UpdatedAt,
		}, nil
	}

	// Decrypt content for unified and biometric modes
	content, err := decryptDiaryContent(encDiary.EncryptedContent, encDiary.IV, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt diary content: %v", err)
	}

	// Parse tags
	var tags []string
	if encDiary.Tags != "" {
		if err := json.Unmarshal([]byte(encDiary.Tags), &tags); err != nil {
			tags = []string{}
		}
	}

	diary := &Diary{
		ID:        encDiary.ID,
		Title:     encDiary.Title,
		Content:   content,
		FileName:  encDiary.FileName,
		FileType:  encDiary.FileType,
		Tags:      tags,
		CreatedAt: encDiary.CreatedAt,
		UpdatedAt: encDiary.UpdatedAt,
	}

	return diary, nil
}

// GetDiaryEncryptionInfo returns encryption information for a diary
func GetDiaryEncryptionInfo(diaryID string, userID uint) (*DiaryEncryptionInfo, error) {
	var encDiary EncryptedDiary
	if err := gormDB.Select("encryption_mode, encryption_salt").Where("id = ? AND user_id = ?", diaryID, userID).First(&encDiary).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diary not found")
		}
		return nil, fmt.Errorf("failed to get diary encryption info: %v", err)
	}

	return &DiaryEncryptionInfo{
		Mode:    encDiary.EncryptionMode,
		HasSalt: encDiary.EncryptionSalt != "",
	}, nil
}

// DiaryEncryptionInfo contains encryption information for a diary
type DiaryEncryptionInfo struct {
	Mode    string `json:"mode"`    // 'unified', 'individual', 'biometric'
	HasSalt bool   `json:"hasSalt"` // Whether the diary has its own salt (individual mode)
}

// SearchEncryptedDiaries searches for diaries in encrypted database
func SearchEncryptedDiaries(query string, userID uint, encryptionKey []byte) ([]Diary, error) {
	// Get all diaries and search in decrypted content
	allDiaries, err := GetEncryptedDiariesList(userID, encryptionKey)
	if err != nil {
		return nil, err
	}

	var matchedDiaries []Diary
	lowerCaseQuery := strings.ToLower(query)

	for _, diary := range allDiaries {
		if strings.Contains(strings.ToLower(diary.Title), lowerCaseQuery) ||
			strings.Contains(strings.ToLower(diary.Content), lowerCaseQuery) {
			matchedDiaries = append(matchedDiaries, diary)
		}
	}

	return matchedDiaries, nil
}

// SearchEncryptedDiariesWithContext searches for diaries with context in encrypted database
func SearchEncryptedDiariesWithContext(query string, userID uint, encryptionKey []byte) ([]SearchResult, error) {
	if query == "" {
		return []SearchResult{}, nil
	}

	// Get all diaries and search in decrypted content
	allDiaries, err := GetEncryptedDiariesList(userID, encryptionKey)
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

			// Determine match type
			if titleMatch && contentMatch {
				result.MatchType = "both"
			} else if titleMatch {
				result.MatchType = "title"
			} else {
				result.MatchType = "content"
			}

			// Extract matched content snippets
			if contentMatch {
				snippets := extractSnippets(diary.Content, query, 100)
				result.MatchedSnippets = snippets
			}

			results = append(results, result)
		}
	}

	return results, nil
}

// DeleteEncryptedDiary deletes a diary from encrypted database
func DeleteEncryptedDiary(diaryID string, userID uint) error {
	result := gormDB.Where("id = ? AND user_id = ?", diaryID, userID).Delete(&EncryptedDiary{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete diary: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("diary not found or not owned by user")
	}

	return nil
}

// decryptDiaryContent decrypts diary content
func decryptDiaryContent(encryptedContent []byte, ivStr string, encryptionKey []byte) (string, error) {
	iv, err := base64.StdEncoding.DecodeString(ivStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode IV: %v", err)
	}

	decryptedData, err := DecryptData(encryptedContent, encryptionKey, iv)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt content: %v", err)
	}

	return string(decryptedData), nil
}
