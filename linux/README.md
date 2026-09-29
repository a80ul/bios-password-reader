# Linux — user awam

Tidak perlu Python. Pilih salah satu:

```bash
./bios-reader-linux
```

atau

```bash
bash bios-reader.sh
```

- Butuh boot UEFI (`/sys/firmware/efi/efivars` ada).
- Kalau `Permission denied`: `chmod +x bios-reader-linux` dulu.
- Kalau output `tidak ditemukan` = password belum diset ATAU bukan InsydeH2O plaintext (normal).
