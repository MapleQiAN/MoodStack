package app

import (
	"encoding/base64"
	"fmt"
)

// BiometricInfo contains information about biometric support
type BiometricInfo struct {
	Supported bool   `json:"supported"`
	Message   string `json:"message"`
}

// EnableBiometricForUser enables biometric authentication for a user
func EnableBiometricForUser(userID uint, password string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// Decode salt and verify password
	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return fmt.Errorf("解码salt失败: %v", err)
	}

	if !VerifyPassword(password, user.PasswordHash, salt) {
		return fmt.Errorf("密码验证失败")
	}

	// Store the current encryption key (derived from password) as the biometric key
	encryptionKey := DeriveKey(password, salt)
	bioKeyStr := base64.StdEncoding.EncodeToString(encryptionKey)

	if err := gormDB.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"biometric_enabled": true,
		"biometric_key":     bioKeyStr,
	}).Error; err != nil {
		return fmt.Errorf("更新用户生物识别设置失败: %v", err)
	}

	return nil
}

// DisableBiometricForUser disables biometric authentication for a user
func DisableBiometricForUser(userID uint) error {
	if err := gormDB.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"biometric_enabled": false,
		"biometric_key":     "",
	}).Error; err != nil {
		return fmt.Errorf("禁用生物识别失败: %v", err)
	}
	return nil
}

// GetUserByID retrieves user by ID
func GetUserByID(userID uint) (*User, error) {
	var user User
	if err := gormDB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
