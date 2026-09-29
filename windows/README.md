# Windows — user awam

1. Klik kanan `bios-reader.exe` → **Run as Administrator**
2. Password langsung tampil. Jendela tidak langsung tutup (tunggu Enter).
3. Syarat: boot UEFI, Windows 7+ / 10 / 11.

Jika SmartScreen / antivirus blokir `.exe` (wajar untuk exe baru tanpa sertifikat):
- Alternatif: klik kanan `bios-reader.bat` → Run as Administrator (jalan via PowerShell, source terbuka di `bios_pw_read_windows.ps1`).

Catatan: Dell/HP/Lenovo baru yang simpan hash di EEPROM akan terbaca `tidak ditemukan` — itu batas hardware, bukan bug.
