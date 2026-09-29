# BIOS Password Reader

Baca password BIOS Supervisor/User yang tersimpan **plaintext** di UEFI NVRAM (bug InsydeH2O, CVE-2021-43613).
Untuk **laptop milik sendiri / seizin pemilik**. Tidak untuk bobol laptop orang lain.

Terbukti jalan di: **Axioo MyBook Pro K5 (InsydeH2O)** — variabel `SystemSupervisorPw` terbaca `123`.

## Dukung BIOS apa saja?

| Kondisi | Hasil |
|---|---|
| InsydeH2O plaintext (Axioo, Acer lama, HP lama, Lenovo lama, RedmiBook, dll) | ✅ Terbaca langsung |
| AMI / Phoenix / Dell / HP / Lenovo baru (hash di EEPROM, bukan NVRAM) | ❌ `tidak ditemukan` (normal, butuh bongkar/flash) |
| Boot Legacy BIOS | ❌ butuh UEFI |
| Password belum diset | ❌ variabel memang tidak ada (normal) |

Tool ini **best-effort**: scan semua variabel `*SupervisorPw*`, `*UserPw*`, `*Password*`. Tidak ada brute-force, tidak ada bypass paksa.

## Struktur folder (mana untuk siapa?)

```
bios-password-reader/
├── linux/                  ← UNTUK USER LINUX (awam)
│   ├── bios-reader-linux   ← tinggal klik / ./bios-reader-linux (tanpa Python)
│   └── bios-reader.sh      ← alternatif tanpa install: bash bios-reader.sh
├── windows/                ← UNTUK USER WINDOWS (awam)
│   ├── bios-reader.exe     ← tinggal double-click Run as Administrator
│   ├── bios-reader.bat     ← alternatif jika .exe diblokir
│   └── bios_pw_read_windows.ps1  ← source PowerShell
├── src/                    ← SOURCE TERBUKA (untuk developer / modifikasi)
│   ├── main.go, reader_linux.go, reader_windows.go, go.mod
│   ├── bios_pw_read_linux.py
│   └── bios_pw_read_windows.ps1
├── README.md
└── LICENSE
```

Source **sengaja disertakan sebelum compiler** agar komunitas bisa audit, perbaiki, dan pull-request.

## Cara pakai — Linux (user awam, tanpa Python)

```bash
cd linux
chmod +x bios-reader-linux bios-reader.sh
./bios-reader-linux
# atau
bash bios-reader.sh
```

Contoh output:
```
[+] SystemSupervisorPw (7f9102df-...)
    raw: 070000000331323367
    [+] kemungkinan password: '123'
```

Lalu reboot → `F2` → masukkan `123` → `Security > Set Supervisor Password` → old diisi, new dikosongkan untuk hapus.

## Cara pakai — Windows (user awam, tinggal exe)

1. Klik kanan `windows/bios-reader.exe` → **Run as Administrator**
2. Password langsung tampil di jendela console, tutJendela tidak langsung tertutup (tunggu Enter).
3. Syarat: boot UEFI, Windows 7+ / 10 / 11, PowerShell tidak dibutuhkan untuk `.exe`.

Alternatif (jika SmartScreen blokir exe):
```
klik kanan windows/bios-reader.bat → Run as Administrator
```

## Build dari source (developer)

Butuh Go 1.22+:

```bash
cd src
go vet ./...
GOOS=linux GOARCH=amd64 go build -o ../linux/bios-reader-linux .
GOOS=windows GOARCH=amd64 go build -o ../windows/bios-reader.exe .
```

Python / PowerShell asli tetap ada di `src/` sebagai referensi audit.

## Keamanan & tanggung jawab

- Tool hanya **membaca** NVRAM, tidak menulis/flash. Risiko brick minimal, tapi tetap **tanpa garansi**.
- Jangan pakai di laptop curian, inventaris kantor tanpa izin IT, atau untuk bypass keamanan pihak lain.
- Pemilik repo **tidak bertanggung jawab** atas modifikasi komunitas maupun penyalahgunaan untuk kejahatan.

## Lisensi / Copyright

© 2026 — Bebas dipakai, dimodifikasi, dan disebar ulang **dengan syarat: JANGAN untuk kejahatan / akses tanpa izin**.
Lihat `LICENSE` lengkap. Modifikasi Anda = tanggung jawab Anda.

---
Pull request welcome: tambah GUID vendor baru, parsing hash baru, atau terjemahan README.
