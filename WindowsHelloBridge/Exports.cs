using System;
using System.Runtime.InteropServices;
using Windows.Security.Credentials.UI;
using Windows.System;
using System.Threading.Tasks;
using Windows.UI.Popups;

namespace WindowsHelloBridge;

public static class Exports
{
    // 返回码与 Go 端对齐：0=Verified, 1=DeviceNotPresent, 2=NotConfiguredForUser, 3=DisabledByPolicy, 4=DeviceBusy, 5=RetriesExhausted, 6=Canceled, <0=其他错误
    [UnmanagedCallersOnly(EntryPoint = "WH_CheckAvailability")]
    public static int WH_CheckAvailability()
    {
        try
        {
            // Windows Hello 需要在系统上可用且配置
            var availabilityTask = UserConsentVerifier.CheckAvailabilityAsync().AsTask();
            availabilityTask.Wait();
            var availability = availabilityTask.Result;

            return availability switch
            {
                UserConsentVerifierAvailability.Available => 0,
                UserConsentVerifierAvailability.DeviceNotPresent => 1,
                UserConsentVerifierAvailability.NotConfiguredForUser => 2,
                UserConsentVerifierAvailability.DisabledByPolicy => 3,
                UserConsentVerifierAvailability.DeviceBusy => 4,
                _ => -1,
            };
        }
        catch (Exception)
        {
            return -2;
        }
    }

    // message: UTF-16 指针；返回与 Go 端约定的结果码
    [UnmanagedCallersOnly(EntryPoint = "WH_RequestVerification")]
    public static unsafe int WH_RequestVerification(char* message)
    {
        try
        {
            string prompt = new string(message);
            var verifyTask = UserConsentVerifier.RequestVerificationAsync(prompt).AsTask();
            verifyTask.Wait();
            var result = verifyTask.Result;

            return result switch
            {
                UserConsentVerificationResult.Verified => 0,
                UserConsentVerificationResult.DeviceNotPresent => 1,
                UserConsentVerificationResult.NotConfiguredForUser => 2,
                UserConsentVerificationResult.DisabledByPolicy => 3,
                UserConsentVerificationResult.DeviceBusy => 4,
                UserConsentVerificationResult.RetriesExhausted => 5,
                UserConsentVerificationResult.Canceled => 6,
                _ => -1,
            };
        }
        catch (Exception)
        {
            return -2;
        }
    }

    // 支持父窗口句柄的验证函数：message: UTF-16 指针，parentHwnd: 父窗口句柄
    [UnmanagedCallersOnly(EntryPoint = "WH_RequestVerificationWithParent")]
    public static unsafe int WH_RequestVerificationWithParent(char* message, IntPtr parentHwnd)
    {
        try
        {
            string prompt = new string(message);
            
            // 如果提供了有效的父窗口句柄，设置窗口居中
            if (parentHwnd != IntPtr.Zero)
            {
                // 设置父窗口为前台窗口，确保Windows Hello对话框在其上方显示
                SetForegroundWindow(parentHwnd);
                
                // 获取父窗口位置和大小
                if (GetWindowRect(parentHwnd, out RECT parentRect))
                {
                    // 计算父窗口中心点
                    int parentCenterX = (parentRect.Left + parentRect.Right) / 2;
                    int parentCenterY = (parentRect.Top + parentRect.Bottom) / 2;
                    
                    // Windows Hello对话框通常会自动在活动窗口上方居中显示
                    // 通过设置父窗口为前台窗口来确保这一点
                }
            }
            
            var verifyTask = UserConsentVerifier.RequestVerificationAsync(prompt).AsTask();
            verifyTask.Wait();
            var result = verifyTask.Result;

            return result switch
            {
                UserConsentVerificationResult.Verified => 0,
                UserConsentVerificationResult.DeviceNotPresent => 1,
                UserConsentVerificationResult.NotConfiguredForUser => 2,
                UserConsentVerificationResult.DisabledByPolicy => 3,
                UserConsentVerificationResult.DeviceBusy => 4,
                UserConsentVerificationResult.RetriesExhausted => 5,
                UserConsentVerificationResult.Canceled => 6,
                _ => -1,
            };
        }
        catch (Exception)
        {
            return -2;
        }
    }

    // Windows API 导入
    [DllImport("user32.dll")]
    private static extern bool SetForegroundWindow(IntPtr hWnd);
    
    [DllImport("user32.dll")]
    private static extern bool GetWindowRect(IntPtr hWnd, out RECT lpRect);
    
    [StructLayout(LayoutKind.Sequential)]
    private struct RECT
    {
        public int Left;
        public int Top;
        public int Right;
        public int Bottom;
    }
}

