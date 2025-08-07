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
