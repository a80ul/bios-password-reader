//go:build windows

package main

import (
	"encoding/hex"
	"syscall"
	"unsafe"
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procGetFwEnvVarW         = kernel32.NewProc("GetFirmwareEnvironmentVariableW")
	advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procOpenProcessToken     = advapi32.NewProc("OpenProcessToken")
	procLookupPrivilegeValue = advapi32.NewProc("LookupPrivilegeValueW")
	procAdjustTokenPriv      = advapi32.NewProc("AdjustTokenPrivileges")
)

const (
	tokenAdjustPrivileges = 0x20
	tokenQuery            = 0x8
	sePrivilegeEnabled    = 0x2
)

type tokenPrivs struct {
	Count uint32
	Luid  uint64
	Attr  uint32
}

func enablePrivilege(name string) {
	var tok syscall.Handle
	cur, _ := syscall.GetCurrentProcess()
	r1, _, _ := procOpenProcessToken.Call(uintptr(cur), uintptr(tokenAdjustPrivileges|tokenQuery), uintptr(unsafe.Pointer(&tok)))
	if r1 == 0 {
		return
	}
	defer syscall.CloseHandle(tok)
	var luid uint64
	n16, _ := syscall.UTF16PtrFromString(name)
	r1, _, _ = procLookupPrivilegeValue.Call(0, uintptr(unsafe.Pointer(n16)), uintptr(unsafe.Pointer(&luid)))
	if r1 == 0 {
		return
	}
	tp := tokenPrivs{Count: 1, Luid: luid, Attr: sePrivilegeEnabled}
	procAdjustTokenPriv.Call(uintptr(tok), 0, uintptr(unsafe.Pointer(&tp)), 0, 0, 0)
}

func getFwVar(name, guid string) ([]byte, bool) {
	n16, _ := syscall.UTF16PtrFromString(name)
	g16, _ := syscall.UTF16PtrFromString(guid)
	buf := make([]byte, 4096)
	ret, _, _ := procGetFwEnvVarW.Call(
		uintptr(unsafe.Pointer(n16)),
		uintptr(unsafe.Pointer(g16)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if ret == 0 {
		return nil, false
	}
	return buf[:ret], true
}

func readPasswords() []Result {
	enablePrivilege("SeSystemEnvironmentPrivilege")
	// Daftar variabel umum lintas vendor.
	// InsydeH2O (Axioo, Acer, HP lama, Lenovo lama, RedmiBook, dll): GUID 7f9102df-e999-4740-80a6-b2038512217b
	// Yang lain (Dell/HP baru) biasanya hash di EEPROM -> akan terbaca "tidak ditemukan", itu normal.
	targets := []struct{ Name, GUID string }{
		{"SystemSupervisorPw", "{7f9102df-e999-4740-80a6-b2038512217b}"},
		{"SystemUserPw", "{7f9102df-e999-4740-80a6-b2038512217b}"},
	}
	var out []Result
	for _, t := range targets {
		data, ok := getFwVar(t.Name, t.GUID)
		if !ok {
			continue
		}
		pw, note, found := parseInsyde(data)
		r := Result{
			Name:   t.Name,
			GUID:   t.GUID,
			RawHex: hex.EncodeToString(data),
			Pass:   pw,
			Note:   note,
			Found:  found,
		}
		empty := true
		for _, c := range data {
			if c != 0 {
				empty = false
				break
			}
		}
		if empty {
			r.Found = false
			r.Empty = true
			r.Note = "kosong (belum diset)"
		}
		out = append(out, r)
	}
	return out
}
