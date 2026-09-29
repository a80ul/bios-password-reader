#!/usr/bin/env python3
"""
InsydeH2O BIOS password reader (Linux).
Membaca plaintext supervisor/user password dari efivarfs.
Butuh: boot UEFI, file /sys/firmware/efi/efivars ada.
Tidak butuh root untuk baca di kebanyakan distro, tapi butuh root di sebagian.

Cara pakai:
  python3 bios_pw_read_linux.py
  sudo python3 bios_pw_read_linux.py   (kalau permission denied)

Hanya untuk laptop milik sendiri / izin pemilik.
"""
import os
import glob
import sys

EFIVARS = "/sys/firmware/efi/efivars"
TARGETS = ["SystemSupervisorPw", "SystemUserPw", "SystemHddPw"]

def parse_insyde_payload(payload: bytes):
    # Format Insyde: [len 1 byte][password ASCII len bytes][extra/checksum...]
    # Contoh: 03 31 32 33 67 -> len=3, pw="123"
    if not payload:
        return None, "kosong"
    length = payload[0]
    raw = payload[1:]
    pw_bytes = raw[:length] if length <= len(raw) else raw
    try:
        pw = pw_bytes.decode("ascii")
    except Exception:
        pw = repr(pw_bytes)
    printable = all(32 <= c < 127 for c in pw_bytes)
    return pw, "" if printable else " (non-printable, mungkin bukan ASCII)"

def main():
    if not os.path.isdir(EFIVARS):
        print(f"[-] {EFIVARS} tidak ada. Pastikan boot UEFI, bukan Legacy BIOS.")
        sys.exit(1)

    found = False
    for target in TARGETS:
        for path in glob.glob(os.path.join(EFIVARS, target + "-*")):
            found = True
            print(f"[+] {os.path.basename(path)}")
            try:
                with open(path, "rb") as f:
                    data = f.read()
            except PermissionError:
                print("    [!] Permission denied, coba: sudo python3 bios_pw_read_linux.py")
                continue
            print(f"    raw hex: {data.hex()}")
            if len(data) < 5:
                print("    [-] terlalu pendek / tidak ada password (belum diset)")
                continue
            attr = data[:4]
            payload = data[4:]
            print(f"    attr: {attr.hex()}, payload: {payload.hex()}")
            pw, note = parse_insyde_payload(payload)
            if payload[1:1+payload[0]] == b"\x00" * payload[0] or all(c == 0 for c in payload):
                print("    [-] password kosong (belum diset)")
            else:
                print(f"    [+] kemungkinan password: '{pw}'{note}")
                # fallback: kalau ada byte sisa printable, tampilkan juga
                rest = payload[1+payload[0]:]
                if rest and all(32 <= c < 127 for c in rest[:8]):
                    print(f"    [?] byte sisa printable: {rest!r} (coba gabungan: '{pw + rest.decode('ascii', errors='ignore')}' jika '{pw}' gagal)")
            print()
    if not found:
        print("[-] Variabel SystemSupervisorPw/SystemUserPw tidak ditemukan.")
        print("    Kemungkinan: bukan InsydeH2O, atau password disimpan di EEPROM terpisah.")
        print("    Cek manual: ls /sys/firmware/efi/efivars/ | grep -i Pw")

if __name__ == "__main__":
    main()
