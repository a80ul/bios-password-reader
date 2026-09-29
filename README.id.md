# BIOS Password Reader — Versi Indonesia

> Lupa password BIOS? Santai, banyak yang ngalamin. 😅
> Tool kecil ini bacain lagi password Supervisor/User **kalau** BIOS kamu nyimpennya masih plaintext di NVRAM UEFI (bug klasik InsydeH2O, CVE-2021-43613).
>
> 👉 **Khusus laptop milik sendiri / seizin pemilik ya.** Bukan buat bobol laptop orang.
>
> English version: [README.md](README.md).

Udah kebukti jalan di: **Axioo MyBook Pro K5 (InsydeH2O)** — `SystemSupervisorPw` kebaca `123`. Iya, sesimpel itu.

## Laptop aku bisa nggak?

| Kondisi kamu | Hasilnya |
|---|---|
| InsydeH2O plaintext (Axioo, Acer/HP/Lenovo lama, RedmiBook, dll) | ✅ Password langsung muncul |
| Dell / HP / Lenovo baru, AMI / Phoenix (hash di EEPROM, bukan NVRAM) | ❌ Muncul `tidak ditemukan` — wajar, harus flash hardware |
| Boot Legacy BIOS | ❌ Harus UEFI |
| Password memang belum diset | ❌ Variabel nggak ada = normal, aman |

Nggak ada brute-force, nggak ada bobol paksa. Cuma scan `*SupervisorPw*`, `*UserPw*`, `*Password*` terus tampilin isinya.

## Folder mana buat aku?

```
bios-password-reader/
├── linux/                  ← USER LINUX mulai sini
│   ├── bios-reader-linux   ← tinggal jalanin, nggak perlu Python
│   └── bios-reader.sh      ← versi cadangan, nggak perlu install
├── windows/                ← USER WINDOWS mulai sini
│   ├── bios-reader.exe     ← double-click, Run as Administrator
│   ├── bios-reader.bat     ← cadangan kalau .exe diblokir
│   └── bios_pw_read_windows.ps1
├── src/                    ← SOURCE TERBUKA (buat oprekers)
│   ├── main.go, reader_linux.go, reader_windows.go, go.mod
│   ├── bios_pw_read_linux.py
│   └── bios_pw_read_windows.ps1
├── README.md               ← versi English
├── README.id.md            ← kamu lagi baca ini
└── LICENSE
```

Source sengaja disimpen **sebelum di-compile** biar siapa aja bisa audit, benerin, atau kirim PR.

## Cara pakai — Linux (nggak perlu Python, beneran)

```bash
cd linux
chmod +x bios-reader-linux bios-reader.sh
./bios-reader-linux
# atau
bash bios-reader.sh
```

Nanti muncul kayak gini:
```
[+] SystemSupervisorPw (7f9102df-...)
    raw: 070000000331323367
    [+] kemungkinan password: '123'
```

Terus reboot → spam `F2` → ketik `123` → masuk `Security > Set Supervisor Password` → isi old, new dikosongin buat hapus.

## Cara pakai — Windows (tinggal exe)

1. Klik kanan `windows/bios-reader.exe` → **Run as Administrator**
2. Password langsung nongol di jendela hitam. Nggak langsung ketutup, nunggu Enter dulu.
3. Syarat: boot UEFI, Windows 7+ / 10 / 11. `.exe`-nya nggak butuh PowerShell.

Kalau SmartScreen ngeblokir exe (wajar, exe baru belum sertifikat):
```
klik kanan windows/bios-reader.bat → Run as Administrator
```

## Build dari source (buat dev)

Butuh Go 1.22+:

```bash
cd src
go vet ./...
GOOS=linux GOARCH=amd64 go build -o ../linux/bios-reader-linux .
GOOS=windows GOARCH=amd64 go build -o ../windows/bios-reader.exe .
```

Script Python / PowerShell aslinya tetap ada di `src/` buat referensi.

## Keamanan & tanggung jawab

- Read-only: cuma baca NVRAM, nggak nulis/flash. Risiko brick kecil, tapi tetap **tanpa garansi**.
- Tolong jangan dipake di laptop curian, inventaris kantor tanpa izin IT, atau buat bypass keamanan orang lain.
- Pemilik repo **nggak bertanggung jawab** atas modifikasi komunitas maupun penyalahgunaan buat kejahatan.

## Lisensi / Copyright

© 2026 — Bebas dipake, dimodifikasi, disebar ulang, **dengan satu syarat: JANGAN buat kejahatan / akses tanpa izin**.
Lihat `LICENSE` lengkapnya. Modifikasi kamu = tanggung jawab kamu.

---
PR welcome: tambah GUID vendor baru, parser hash baru, atau benerin terjemahan. 🙌
