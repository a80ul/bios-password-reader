package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
)

type Result struct {
	Name    string
	GUID    string
	RawHex  string
	Pass    string
	Note    string
	Found   bool
	Empty   bool
}

// parseInsyde mem-parsing format InsydeH2O:
// payload = [len 1 byte][password ASCII len bytes][sisa/checksum...]
func parseInsyde(payload []byte) (string, string, bool) {
	if len(payload) == 0 {
		return "", "kosong", false
	}
	if len(payload) == 1 {
		return "", "terlalu pendek", false
	}
	n := int(payload[0])
	raw := payload[1:]
	if n <= 0 || n > len(raw) {
		// fallback: tampilkan semua printable
		s := printableASCII(raw)
		if s == "" {
			return "", "tidak ada password / format tidak dikenal", false
		}
		return s, "format len tidak cocok, tampil mentah", true
	}
	pwBytes := raw[:n]
	// cek kosong (semua nol)
	empty := true
	for _, c := range pwBytes {
		if c != 0 {
			empty = false
			break
		}
	}
	if empty {
		return "", "kosong (belum diset)", false
	}
	pw := string(pwBytes)
	if !isPrintable(pwBytes) {
		return hex.EncodeToString(pwBytes), "non-ASCII (hash? bukan plaintext)", true
	}
	note := ""
	rest := raw[n:]
	if len(rest) > 0 && isPrintable(rest) && len(rest) <= 16 {
		// byte sisa kadang checksum printable, infokan sebagai fallback
		note = fmt.Sprintf("jika '%s' gagal, coba juga '%s%s'", pw, pw, string(rest))
	}
	return pw, note, true
}

func isPrintable(b []byte) bool {
	for _, c := range b {
		if c < 32 || c > 126 {
			return false
		}
	}
	return len(b) > 0
}

func printableASCII(b []byte) string {
	out := make([]byte, 0, len(b))
	for _, c := range b {
		if c >= 32 && c <= 126 {
			out = append(out, c)
		}
	}
	return string(out)
}

func main() {
	fmt.Println("=== BIOS Password Reader (InsydeH2O plaintext) ===")
	fmt.Printf("OS: %s |_butuh: boot UEFI_\n\n", runtime.GOOS)
	fmt.Println("Hanya untuk laptop milik sendiri / seizin pemilik.");
	results := readPasswords()
	if len(results) == 0 {
		fmt.Println("[-] Variabel password tidak ditemukan.")
		printNotFoundHelp()
		waitEnter()
		os.Exit(0)
	}
	anyPass := false
	for _, r := range results {
		fmt.Printf("[+] %s (%s)\n", r.Name, r.GUID)
		fmt.Printf("    raw: %s\n", r.RawHex)
		if !r.Found || r.Empty {
			fmt.Printf("    [-] %s\n\n", r.Note)
			continue
		}
		anyPass = true
		fmt.Printf("    [+] kemungkinan password: '%s'\n", r.Pass)
		if r.Note != "" {
			fmt.Printf("    [?] %s\n", r.Note)
		}
		fmt.Println()
	}
	if !anyPass {
		fmt.Println("[-] Tidak ada password plaintext yang terbaca (mungkin belum diset / bukan plaintext).")
		printNotFoundHelp()
	} else {
		fmt.Println("[*] Masuk BIOS (F2/Del), masukkan password di atas.")
		fmt.Println("[*] Untuk hapus: Security > Set Supervisor Password > isi old, new dikosongkan.")
	}
	waitEnter()
}

func printNotFoundHelp() {
	fmt.Println("Kemungkinan penyebab:")
	fmt.Println(" - Bukan InsydeH2O plaintext (Dell/HP/Lenovo baru simpan hash di EEPROM).")
	fmt.Println(" - Password memang belum diset (variabel tidak ada = normal).")
	fmt.Println(" - Boot Legacy BIOS (bukan UEFI) atau Secure Boot mengunci akses.")
	fmt.Println(" - Di Windows: belum Run as Administrator.")
}

func waitEnter() {
	// Agar user awam yang double-click sempat baca output
	fmt.Println("\nTekan Enter untuk keluar...")
	var b [1]byte
	os.Stdin.Read(b[:])
}
