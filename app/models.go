package app

import (
	"encoding/json"
	"time"
)

// User represents a user in the system
type User struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Username         string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash     string    `gorm:"not null" json:"-"`
	Salt             string    `gorm:"not null" json:"-"`
	BiometricEnabled bool      `gorm:"default:false" json:"biometricEnabled"`
	BiometricKey     string    `json:"-"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Associations
	Sessions         []Session        `json:"-"`
	EncryptedDiaries []EncryptedDiary `json:"-"`
}

// Session represents an authentication session
type Session struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"userId"`
	SessionKey string    `gorm:"not null" json:"sessionKey"`
	ExpiresAt  time.Time `gorm:"not null;index" json:"expiresAt"`
	CreatedAt  time.Time `json:"createdAt"`

	// Associations
	User User `json:"-"`
}

// EncryptedDiary represents an encrypted diary entry in the database
type EncryptedDiary struct {
	ID               string    `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"userId"`
	Title            string    `gorm:"not null" json:"title"`
	EncryptedContent []byte    `gorm:"not null" json:"-"`
	IV               string    `gorm:"not null" json:"-"`
	FileName         string    `json:"fileName"`
	FileType         string    `json:"fileType"`
	Tags             string    `gorm:"type:text" json:"-"`
	EncryptionMode   string    `gorm:"not null;default:'unified'" json:"encryptionMode"`
	EncryptionSalt   string    `json:"-"`
	CreatedAt        time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Associations
	User User `json:"-"`
}

// GetTags returns the tags as a slice
func (ed *EncryptedDiary) GetTags() []string {
	if ed.Tags == "" {
		return []string{}
	}
	var tags []string
	if err := json.Unmarshal([]byte(ed.Tags), &tags); err != nil {
		return []string{}
	}
	return tags
}

// SetTags sets the tags from a slice
func (ed *EncryptedDiary) SetTags(tags []string) error {
	if tags == nil {
		tags = []string{}
	}
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	ed.Tags = string(data)
	return nil
}

// AppSetting represents application settings
type AppSetting struct {
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `gorm:"not null" json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName overrides the table name for EncryptedDiary
func (EncryptedDiary) TableName() string {
	return "encrypted_diaries"
}

// TableName overrides the table name for AppSetting
func (AppSetting) TableName() string {
	return "app_settings"
}

// EmotionAnalysis represents emotion analysis for a diary entry
type EmotionAnalysis struct {
	ID      string `gorm:"primaryKey" json:"id"`
	DiaryID string `gorm:"not null;index;uniqueIndex" json:"diaryId"`
	UserID  uint   `gorm:"not null;index" json:"userId"`

	// Emotion scores (0-1 scale)
	Joy      float64 `json:"joy"`
	Sadness  float64 `json:"sadness"`
	Anger    float64 `json:"anger"`
	Fear     float64 `json:"fear"`
	Love     float64 `json:"love"`
	Surprise float64 `json:"surprise"`
	Disgust  float64 `json:"disgust"`

	// Overall emotion
	DominantEmotion string  `json:"dominantEmotion"`
	Confidence      float64 `json:"confidence"`

	// Sentiment analysis
	SentimentScore float64 `json:"sentimentScore"` // -1 to 1 (negative to positive)
	SentimentLabel string  `json:"sentimentLabel"` // "positive", "negative", "neutral"

	// Keywords extracted from content
	Keywords string `gorm:"type:text" json:"-"` // JSON array of keywords

	// Analysis method used
	AnalysisMethod string `gorm:"not null" json:"analysisMethod"` // "programmatic", "ai"

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Associations
	Diary EncryptedDiary `gorm:"foreignKey:DiaryID;references:ID" json:"-"`
	User  User           `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

// GetKeywords returns the keywords as a slice
func (ea *EmotionAnalysis) GetKeywords() []string {
	if ea.Keywords == "" {
		return []string{}
	}
	var keywords []string
	if err := json.Unmarshal([]byte(ea.Keywords), &keywords); err != nil {
		return []string{}
	}
	return keywords
}

// SetKeywords sets the keywords from a slice
func (ea *EmotionAnalysis) SetKeywords(keywords []string) error {
	if keywords == nil {
		keywords = []string{}
	}
	data, err := json.Marshal(keywords)
	if err != nil {
		return err
	}
	ea.Keywords = string(data)
	return nil
}

// TableName overrides the table name for EmotionAnalysis
func (EmotionAnalysis) TableName() string {
	return "emotion_analyses"
}
