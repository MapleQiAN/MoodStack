package app

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// BiometricSupport checks if biometric authentication is supported
func BiometricSupport() (*BiometricInfo, error) {
	if runtime.GOOS != "windows" {
		return &BiometricInfo{
			Supported: false,
			Message:   "生物识别认证仅在Windows系统上支持",
		}, nil
	}

	// Check if Windows Hello is available
	cmd := exec.Command("powershell", "-Command", `
		$winhello = Get-WindowsOptionalFeature -Online -FeatureName "WindowsHelloFace"
		if ($winhello.State -eq "Enabled") {
			Write-Output "enabled"
		} else {
			$devices = Get-PnpDevice | Where-Object {$_.FriendlyName -like "*biometric*" -or $_.FriendlyName -like "*fingerprint*" -or $_.FriendlyName -like "*face*"}
			if ($devices.Count -gt 0) {
				Write-Output "available"
			} else {
				Write-Output "unavailable"
			}
		}
	`)

	output, err := cmd.Output()
	if err != nil {
		return &BiometricInfo{
			Supported: false,
			Message:   "无法检测Windows Hello状态",
		}, nil
	}

	status := strings.TrimSpace(string(output))

	switch status {
	case "enabled":
		return &BiometricInfo{
			Supported: true,
			Message:   "Windows Hello已启用",
		}, nil
	case "available":
		return &BiometricInfo{
			Supported: true,
			Message:   "检测到生物识别设备，但Windows Hello未完全配置",
		}, nil
	default:
		return &BiometricInfo{
			Supported: false,
			Message:   "未检测到支持的生物识别设备",
		}, nil
	}
}

// BiometricInfo contains information about biometric support
type BiometricInfo struct {
	Supported bool   `json:"supported"`
	Message   string `json:"message"`
}

// AuthenticateWithBiometric performs biometric authentication using Windows Hello
func AuthenticateWithBiometric(username string) (*AuthResult, error) {
	// Check if biometric is supported
	bioInfo, err := BiometricSupport()
	if err != nil {
		return nil, err
	}

	if !bioInfo.Supported {
		return &AuthResult{
			Success: false,
			Message: bioInfo.Message,
		}, nil
	}

	// Get user to check if biometric is enabled
	user, err := GetUserByUsername(username)
	if err != nil {
		return &AuthResult{
			Success: false,
			Message: "用户不存在",
		}, nil
	}

	if !user.BiometricEnabled {
		return &AuthResult{
			Success: false,
			Message: "该用户未启用生物识别认证",
		}, nil
	}

	// Perform Windows Hello authentication
	success, err := performWindowsHelloAuth()
	if err != nil {
		return &AuthResult{
			Success: false,
			Message: fmt.Sprintf("生物识别认证失败: %v", err),
		}, nil
	}

	if !success {
		return &AuthResult{
			Success: false,
			Message: "生物识别认证失败",
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
		Message: "生物识别认证成功",
	}, nil
}

// performWindowsHelloAuth performs the actual Windows Hello authentication
func performWindowsHelloAuth() (bool, error) {
	// PowerShell script to prompt for Windows Hello authentication
	script := `
Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

try {
    # Try to use Windows Hello API
    $credential = [Windows.Security.Credentials.UI.UserConsentVerifier]::RequestVerificationAsync("MoodStack需要验证您的身份")
    $result = $credential.GetAwaiter().GetResult()
    
    if ($result -eq [Windows.Security.Credentials.UI.UserConsentVerificationResult]::Verified) {
        Write-Output "SUCCESS"
    } else {
        Write-Output "FAILED"
    }
} catch {
    # Fallback to credential prompt
    $cred = Get-Credential -Message "Windows Hello不可用，请输入Windows凭据进行验证" -UserName $env:USERNAME
    if ($cred) {
        Write-Output "SUCCESS"
    } else {
        Write-Output "FAILED"
    }
}
`

	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("PowerShell执行失败: %v", err)
	}

	result := strings.TrimSpace(string(output))
	return result == "SUCCESS", nil
}

// EnableBiometricForUser enables biometric authentication for a user
func EnableBiometricForUser(userID uint, password string) error {
	// Verify current password first
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
	// This ensures that biometric unlock can access the same encrypted data
	encryptionKey := DeriveKey(password, salt)
	bioKeyStr := base64.StdEncoding.EncodeToString(encryptionKey)

	// Update user record using GORM
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
