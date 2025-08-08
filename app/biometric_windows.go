//go:build windows

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// WindowsHello DLL bridge exports
var (
	whDLL                     *syscall.DLL
	procWHCheckAvailability   *syscall.Proc
	procWHRequestVerification *syscall.Proc
)

// InitializeWindowsHello loads the native DLL and locates exports
func InitializeWindowsHello(dllPath string) error {
	// 仅当 DLL 与两个过程都已就绪时才视为已初始化
	if whDLL != nil && procWHCheckAvailability != nil && procWHRequestVerification != nil {
		return nil
	}

	var err error
	var loaded *syscall.DLL

	candidatePaths := []string{}
	// 1) explicit path provided by caller
	if dllPath != "" {
		candidatePaths = append(candidatePaths, dllPath)
	}
	// 2) app subdirectory (prefer this to avoid colliding with managed DLLs in root)
	candidatePaths = append(candidatePaths, filepath.Join("app", "WindowsHelloBridge.dll"))
	// 3) current working directory
	candidatePaths = append(candidatePaths, "WindowsHelloBridge.dll")
	// 4) executable directory
	if exePath, e := os.Executable(); e == nil {
		exeDir := filepath.Dir(exePath)
		candidatePaths = append(candidatePaths,
			filepath.Join(exeDir, "WindowsHelloBridge.dll"),
			filepath.Join(exeDir, "..", "WindowsHelloBridge.dll"),
		)
	}
	// 5) common build output locations
	candidatePaths = append(candidatePaths,
		filepath.Join("build", "bin", "WindowsHelloBridge.dll"),
		filepath.Join("bin", "WindowsHelloBridge.dll"),
	)

	var firstErr error
	for _, p := range candidatePaths {
		if p == "" {
			continue
		}
		if _, statErr := os.Stat(p); statErr != nil {
			continue
		}
		// 先加载到临时变量，避免半初始化的全局状态
		loaded, err = syscall.LoadDLL(p)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		// 尝试解析导出到临时变量
		var checkProc, verifyProc *syscall.Proc
		checkProc, err = loaded.FindProc("WH_CheckAvailability")
		if err != nil {
			// 释放并尝试下一个路径
			_ = loaded.Release()
			loaded = nil
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		verifyProc, err = loaded.FindProc("WH_RequestVerification")
		if err != nil {
			_ = loaded.Release()
			loaded = nil
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		// 一切就绪后再写入全局变量
		whDLL = loaded
		procWHCheckAvailability = checkProc
		procWHRequestVerification = verifyProc
		return nil
	}

	if firstErr != nil {
		return fmt.Errorf("加载WindowsHelloBridge失败: %v", firstErr)
	}
	return fmt.Errorf("加载WindowsHelloBridge失败: 未找到可用的DLL路径")
}

// BiometricSupport checks if Windows Hello is available via DLL
func BiometricSupport() (*BiometricInfo, error) {
	if whDLL == nil || procWHCheckAvailability == nil || procWHRequestVerification == nil {
		// 默认从当前目录加载
		if err := InitializeWindowsHello("WindowsHelloBridge.dll"); err != nil {
			return &BiometricInfo{Supported: false, Message: fmt.Sprintf("无法加载Windows Hello模块: %v", err)}, nil
		}
	}
	// 双重校验，防止罕见情况下仍为空
	if procWHCheckAvailability == nil {
		return &BiometricInfo{Supported: false, Message: "Windows Hello检查函数未初始化"}, nil
	}
	r1, _, _ := procWHCheckAvailability.Call()
	code := int(r1)
	switch code {
	case 0:
		return &BiometricInfo{Supported: true, Message: "Windows Hello可用"}, nil
	case 1:
		return &BiometricInfo{Supported: false, Message: "设备不存在"}, nil
	case 2:
		return &BiometricInfo{Supported: false, Message: "未为当前用户配置"}, nil
	case 3:
		return &BiometricInfo{Supported: false, Message: "策略禁用"}, nil
	case 4:
		return &BiometricInfo{Supported: false, Message: "设备忙"}, nil
	default:
		return &BiometricInfo{Supported: false, Message: "状态未知"}, nil
	}
}

// AuthenticateWithBiometric authenticates the user via Windows Hello
func AuthenticateWithBiometric(username string) (*AuthResult, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return &AuthResult{Success: false, Message: "用户不存在"}, nil
	}
	if !user.BiometricEnabled {
		return &AuthResult{Success: false, Message: "该用户未启用生物识别认证"}, nil
	}

	if whDLL == nil || procWHCheckAvailability == nil || procWHRequestVerification == nil {
		if err := InitializeWindowsHello("WindowsHelloBridge.dll"); err != nil {
			return &AuthResult{Success: false, Message: fmt.Sprintf("无法加载Windows Hello模块: %v", err)}, nil
		}
	}
	if procWHRequestVerification == nil {
		return &AuthResult{Success: false, Message: "Windows Hello验证函数未初始化"}, nil
	}

	prompt := "MoodStack需要验证您的身份"
	utf16, err := syscall.UTF16PtrFromString(prompt)
	if err != nil {
		return &AuthResult{Success: false, Message: "提示字符串编码失败"}, nil
	}
	r1, _, _ := procWHRequestVerification.Call(uintptr(unsafe.Pointer(utf16)))
	code := int(r1)
	switch code {
	case 0:
		session, err := CreateSession(user.ID)
		if err != nil {
			return nil, fmt.Errorf("创建会话失败: %v", err)
		}
		return &AuthResult{Success: true, User: user, Session: session, Message: "生物识别认证成功"}, nil
	case 1:
		return &AuthResult{Success: false, Message: "生物识别设备不存在"}, nil
	case 2:
		return &AuthResult{Success: false, Message: "当前用户未配置Windows Hello"}, nil
	case 3:
		return &AuthResult{Success: false, Message: "Windows Hello被策略禁用"}, nil
	case 4:
		return &AuthResult{Success: false, Message: "生物识别设备正忙"}, nil
	case 5:
		return &AuthResult{Success: false, Message: "认证尝试次数已用完"}, nil
	case 6:
		return &AuthResult{Success: false, Message: "用户取消了认证"}, nil
	default:
		return &AuthResult{Success: false, Message: fmt.Sprintf("未知错误: %d", code)}, nil
	}
}
