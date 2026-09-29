#!/bin/bash
# BIOS Password Reader - Linux (tanpa Python, tinggal klik/enter)
# Untuk user awam: klik kanan > Run, atau: bash bios-reader.sh
# Hanya untuk laptop milik sendiri / seizin pemilik.

echo "=== BIOS Password Reader (Linux) ==="
echo "Butuh: boot UEFI. Tidak butuh install apa-apa."
echo ""

EFIDIR="/sys/firmware/efi/efivars"
if [ ! -d "$EFIDIR" ]; then
  echo "[-] $EFIDIR tidak ada. Pastikan boot UEFI, bukan Legacy BIOS."
  read -p "Tekan Enter untuk keluar..."
  exit 1
fi

found=0
for f in "$EFIDIR"/SystemSupervisorPw-* "$EFIDIR"/SystemUserPw-*; do
  [ -e "$f" ] || continue
  found=1
  base=$(basename "$f")
  echo "[+] $base"
  hex=$(xxd -p "$f" 2>/dev/null | tr -d '\n')
  echo "    raw: $hex"
  # payload setelah 4 byte attr (8 hex char)
  payload=${hex:8}
  len_hex=${payload:0:2}
  len=$((16#$len_hex))
  pw_hex=${payload:2:$((len*2))}
  # hex -> ascii
  pw=$(echo "$pw_hex" | xxd -r -p 2>/dev/null)
  if [ -z "$pw" ]; then
    echo "    [-] kosong / tidak terbaca"
  else
    echo "    [+] kemungkinan password: '$pw'"
    rest_hex=${payload:$((2+len*2))}
    if [ -n "$rest_hex" ]; then
      rest=$(echo "$rest_hex" | xxd -r -p 2>/dev/null | tr -cd '[:print:]')
      if [ -n "$rest" ]; then
        echo "    [?] byte sisa: '$rest' (jika '$pw' gagal, coba '$pw$rest')"
      fi
    fi
  fi
  echo ""
done

if [ "$found" -eq 0 ]; then
  echo "[-] Variabel SystemSupervisorPw/SystemUserPw tidak ditemukan."
  echo "    Artinya: belum diset (normal), bukan InsydeH2O plaintext,"
  echo "    atau password tersimpan di EEPROM (Dell/HP/Lenovo baru)."
fi

echo ""
echo "[*] Masuk BIOS (F2/Del), masukkan password di atas."
echo "[*] Hapus via: Security > Set Supervisor Password > old diisi, new dikosongkan."
read -p "Tekan Enter untuk keluar..."
