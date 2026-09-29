//go:build linux

package main

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

func readPasswords() []Result {
	dir := "/sys/firmware/efi/efivars"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []Result
	for _, e := range entries {
		name := e.Name()
		lower := strings.ToLower(name)
		// dukung BIOS apa saja yang menyimpan plaintext di UEFI var:
		// Insyde: SystemSupervisorPw / SystemUserPw, generik: *pw* *password* *supervisor*
		isCandidate := strings.Contains(lower, "supervisorpw") ||
			strings.Contains(lower, "systemuserpw") ||
			strings.Contains(lower, "hddpw") ||
			(strings.Contains(lower, "password") && !strings.Contains(lower, "passwordconfig"))
		if !isCandidate {
			continue
		}
		full := filepath.Join(dir, name)
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		if len(data) < 5 {
			continue
		}
		// file efivarfs = [attr 4 byte][payload...]
		// pisahkan nama dan GUID: Nama-GUID
		var varName, guid string
		if i := strings.LastIndex(name, "-"); i > 0 {
			varName = name[:i]
			guid = name[i+1:]
		} else {
			varName = name
		}
		payload := data[4:]
		pw, note, ok := parseInsyde(payload)
		r := Result{
			Name:   varName,
			GUID:   guid,
			RawHex: hex.EncodeToString(data),
			Pass:   pw,
			Note:   note,
			Found:  ok,
		}
		// deteksi kosong: payload semua nol
		empty := true
		for _, c := range payload {
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
