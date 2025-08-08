//go:build !windows

package app

// Non-Windows stub implementations

func BiometricSupport() (*BiometricInfo, error) {
	return &BiometricInfo{Supported: false, Message: "Windows Hello仅在Windows系统上支持"}, nil
}

func AuthenticateWithBiometric(username string) (*AuthResult, error) {
	return &AuthResult{Success: false, Message: "Windows Hello仅在Windows系统上支持"}, nil
}
