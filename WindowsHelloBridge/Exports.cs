using System;
using System.Runtime.InteropServices;
using Windows.Security.Credentials.UI;
using Windows.System;
using System.Threading.Tasks;

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
}

