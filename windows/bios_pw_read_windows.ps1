#Requires -RunAsAdministrator
<#
.SYNOPSIS
  InsydeH2O BIOS password reader (Windows PowerShell 5.1+ / pwsh 7+).
  Membaca plaintext supervisor/user password dari UEFI NVRAM via GetFirmwareEnvironmentVariableW.

.SYARAT
  - Windows UEFI (bukan Legacy BIOS), boot UEFI
  - Run as Administrator
  - Hanya untuk laptop milik sendiri / izin pemilik

.CARA PAKAI
  powershell -ExecutionPolicy Bypass -File bios_pw_read_windows.ps1

REFERENSI
  API: kernel32!GetFirmwareEnvironmentVariableW, GUID Insyde {7f9102df-e999-4740-80a6-b2038512217b}
#>

$ErrorActionPreference = "Stop"

# 1. Enable SeSystemEnvironmentPrivilege (wajib untuk baca UEFI var di Windows)
Add-Type -Language CSharp -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public class UefiPriv {
    [DllImport("advapi32.dll", SetLastError=true)]
    static extern bool OpenProcessToken(IntPtr h, int acc, out IntPtr tok);
    [DllImport("advapi32.dll", SetLastError=true)]
    static extern bool LookupPrivilegeValue(string sys, string name, out long luid);
    [StructLayout(LayoutKind.Sequential)]
    struct TOKEN_PRIVS { public int Count; public long Luid; public int Attr; }
    [DllImport("advapi32.dll", SetLastError=true)]
    static extern bool AdjustTokenPrivileges(IntPtr tok, bool dis, ref TOKEN_PRIVS p, int len, IntPtr prev, IntPtr relen);
    const int TOKEN_ADJUST_PRIVILEGES = 0x20;
    const int TOKEN_QUERY = 0x8;
    const int SE_PRIVILEGE_ENABLED = 0x2;
    public static void Enable(string name) {
        IntPtr tok = IntPtr.Zero;
        if (!OpenProcessToken(System.Diagnostics.Process.GetCurrentProcess().Handle, TOKEN_ADJUST_PRIVILEGES | TOKEN_QUERY, out tok))
            throw new System.ComponentModel.Win32Exception(Marshal.GetLastWin32Error());
        long luid = 0;
        if (!LookupPrivilegeValue(null, name, out luid))
            throw new System.ComponentModel.Win32Exception(Marshal.GetLastWin32Error());
        TOKEN_PRIVS p = new TOKEN_PRIVS();
        p.Count = 1; p.Luid = luid; p.Attr = SE_PRIVILEGE_ENABLED;
        if (!AdjustTokenPrivileges(tok, false, ref p, 0, IntPtr.Zero, IntPtr.Zero))
            throw new System.ComponentModel.Win32Exception(Marshal.GetLastWin32Error());
    }
}
'@

Add-Type -Language CSharp -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public class UefiRead {
    [DllImport("kernel32.dll", SetLastError=true, CharSet=CharSet.Unicode)]
    public static extern uint GetFirmwareEnvironmentVariableW(string lpName, string lpGuid, IntPtr pBuffer, uint nSize);
}
'@

try { [UefiPriv]::Enable("SeSystemEnvironmentPrivilege") } catch { Write-Warning "Gagal enable privilege: $_" }

function Read-UefiVar {
    param([string]$Name, [string]$Guid)
    $size = 4096
    $ptr = [Runtime.InteropServices.Marshal]::AllocHGlobal($size)
    try {
        $ret = [UefiRead]::GetFirmwareEnvironmentVariableW($Name, $Guid, $ptr, [uint32]$size)
        if ($ret -eq 0) {
            $err = [Runtime.InteropServices.Marshal]::GetLastWin32Error()
            return @{ Name=$Name; Guid=$Guid; Size=0; Error=$err; Data=$null }
        }
        $buf = New-Object byte[] $ret
        [Runtime.InteropServices.Marshal]::Copy($ptr, $buf, 0, $ret)
        return @{ Name=$Name; Guid=$Guid; Size=$ret; Error=0; Data=$buf }
    } finally {
        [Runtime.InteropServices.Marshal]::FreeHGlobal($ptr)
    }
}

# GUID Insyde supervisor (dari Linux: 7f9102df-e999-4740-80a6-b2038512217b)
$targets = @(
    @{ Name="SystemSupervisorPw"; Guid="{7f9102df-e999-4740-80a6-b2038512217b}" },
    @{ Name="SystemUserPw";       Guid="{7f9102df-e999-4740-80a6-b2038512217b}" }
)

$foundAny = $false
foreach ($t in $targets) {
    $r = Read-UefiVar -Name $t.Name -Guid $t.Guid
    Write-Host "`n[+] $($r.Name) $($r.Guid)"
    if ($r.Size -eq 0) {
        Write-Host "    [-] tidak terbaca. Win32Error=$($r.Error) (1=Legacy BIOS, 203=var tidak ada, 1314=belum admin)"
        continue
    }
    $foundAny = $true
    $hex = ($r.Data | ForEach-Object { $_.ToString("x2") }) -join " "
    Write-Host "    raw hex ($($r.Size) byte): $hex"
    if ($r.Size -lt 2) { Write-Host "    [-] terlalu pendek"; continue }
    $len = $r.Data[0]
    if ($len -gt ($r.Size - 1)) { $len = $r.Size - 1 }
    $pwBytes = $r.Data[1..$len]
    $pw = [Text.Encoding]::ASCII.GetString($pwBytes)
    Write-Host "    [+] kemungkinan password: '$pw'"
    if ($r.Size -gt ($len + 1)) {
        $rest = $r.Data[($len+1)..($r.Size-1)]
        $restStr = [Text.Encoding]::ASCII.GetString($rest)
        if ($restStr -match '^[\x20-\x7e]+$') {
            Write-Host "    [?] byte sisa printable: '$restStr' (jika '$pw' gagal, coba '$pw$restStr')"
        }
    }
}

if (-not $foundAny) {
    Write-Host "`n[-] Variabel tidak ditemukan. Kemungkinan bukan InsydeH2O atau password di EEPROM terpisah."
}
