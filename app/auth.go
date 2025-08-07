package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AuthResult represents the result of authentication
type AuthResult struct {
	Success      bool     `json:"success"`
	User         *User    `json:"user,omitempty"`
	Session      *Session `json:"session,omitempty"`
	Message      string   `json:"message"`
	RequireSetup bool     `json:"requireSetup"`
}

// CreateUser creates a new user with password
func CreateUser(username, password string) (*User, error) {
	// Check if user already exists
	if exists, err := UserExists(username); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("user already exists")
	}

	// Generate salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}

	// Hash password
	passwordHash := HashPassword(password, salt)
	saltStr := base64.StdEncoding.EncodeToString(salt)

	// Create user using GORM
	user := &User{
		Username:         username,
		PasswordHash:     passwordHash,
		Salt:             saltStr,
		BiometricEnabled: false,
	}

	if err := gormDB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// AuthenticateWithPassword authenticates user with password
func AuthenticateWithPassword(username, password string) (*AuthResult, error) {
	// Get user from database
	user, err := GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &AuthResult{
				Success: false,
				Message: "用户不存在",
			}, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Decode salt
	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %v", err)
	}

	// Verify password
	if !VerifyPassword(password, user.PasswordHash, salt) {
		return &AuthResult{
			Success: false,
			Message: "密码错误",
		}, nil
	}

	// Create session
	session, err := CreateSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &AuthResult{
		Success: true,
		User:    user,
		Session: session,
		Message: "登录成功",
	}, nil
}

// GetUserByUsername retrieves user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := gormDB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return &user, nil
}

// UserExists checks if a user exists
func UserExists(username string) (bool, error) {
	var count int64
	if err := gormDB.Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasAnyUsers checks if there are any users in the system
func HasAnyUsers() (bool, error) {
	var count int64
	if err := gormDB.Model(&User{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateSession creates a new authentication session
func CreateSession(userID uint) (*Session, error) {
	// Generate session ID and key
	sessionID, err := GenerateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %v", err)
	}

	sessionKeyBytes := make([]byte, 32)
	if _, err := rand.Read(sessionKeyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate session key: %v", err)
	}
	sessionKey := hex.EncodeToString(sessionKeyBytes)

	// Session expires in 7 days
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Create session using GORM
	session := &Session{
		ID:         sessionID,
		UserID:     userID,
		SessionKey: sessionKey,
		ExpiresAt:  expiresAt,
	}

	if err := gormDB.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return session, nil
}

// ValidateSession validates a session and returns the user
func ValidateSession(sessionID, sessionKey string) (*User, error) {
	var session Session
	if err := gormDB.Where("id = ? AND session_key = ? AND expires_at > ?", sessionID, sessionKey, time.Now()).First(&session).Error; err != nil {
		return nil, err
	}

	var user User
	if err := gormDB.First(&user, session.UserID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteSession deletes a session (logout)
func DeleteSession(sessionID string) error {
	return gormDB.Where("id = ?", sessionID).Delete(&Session{}).Error
}

// CleanupExpiredSessions removes expired sessions
func CleanupExpiredSessions() error {
	return gormDB.Where("expires_at < ?", time.Now()).Delete(&Session{}).Error
}
