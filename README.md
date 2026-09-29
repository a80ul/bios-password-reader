# BIOS Password Reader

> Forgot your BIOS password? Yeah, it happens. 😅
> This little tool reads back the Supervisor/User password **if** your BIOS left it as plaintext in UEFI NVRAM (the classic InsydeH2O bug, CVE-2021-43613).
>
> 👉 **For your own laptop / with the owner's permission only.** Not for breaking into other people's machines.
>
> Baca versi Indonesia di [README.id.md](README.id.md).

Tested on: **Axioo MyBook Pro K5 (InsydeH2O)** — `SystemSupervisorPw` came out as `123`. Yep, really that simple.

## Will it work on my BIOS?

| Your case | What happens |
|---|---|
| InsydeH2O plaintext (Axioo, old Acer/HP/Lenovo, RedmiBook, etc.) | ✅ Shows the password right away |
| New Dell / HP / Lenovo, AMI / Phoenix (hash in EEPROM, not NVRAM) | ❌ Says `not found` — that's normal, needs hardware flash |
| Legacy BIOS boot | ❌ Needs UEFI |
| No password set | ❌ No variable = nothing to show, all good |

No brute-force, no forced bypass. It just scans `*SupervisorPw*`, `*UserPw*`, `*Password*` and shows what's already there.

## Which folder is for me?

```
bios-password-reader/
├── linux/                  ← LINUX folks, start here
│   ├── bios-reader-linux   ← just run it, no Python needed
│   └── bios-reader.sh      ← alt version, no install needed
├── windows/                ← WINDOWS folks, start here
│   ├── bios-reader.exe     ← double-click, Run as Administrator
│   ├── bios-reader.bat     ← backup if .exe gets blocked
│   └── bios_pw_read_windows.ps1
├── src/                    ← OPEN SOURCE (devs, tinker here)
│   ├── main.go, reader_linux.go, reader_windows.go, go.mod
│   ├── bios_pw_read_linux.py
│   └── bios_pw_read_windows.ps1
├── README.md               ← you are here (English)
├── README.id.md            ← versi Indonesia
└── LICENSE
```

Source code is kept **before compiling** on purpose, so anyone can audit it, improve it, or send a PR.

## Quick start — Linux (no Python, promise)

```bash
cd linux
chmod +x bios-reader-linux bios-reader.sh
./bios-reader-linux
# or
bash bios-reader.sh
```

You'll see something like:
```
[+] SystemSupervisorPw (7f9102df-...)
    raw: 070000000331323367
    [+] possible password: '123'
```

Then reboot → press `F2` → type `123` → go to `Security > Set Supervisor Password` → fill old, leave new empty to clear it.

## Quick start — Windows (just the .exe)

1. Right-click `windows/bios-reader.exe` → **Run as Administrator**
2. Password pops up in the console. Window stays open until you press Enter.
3. Needs: UEFI boot, Windows 7+ / 10 / 11. No PowerShell needed for the `.exe`.

If SmartScreen blocks the exe (normal for a new unsigned exe):
```
right-click windows/bios-reader.bat → Run as Administrator
```

## Build from source (devs)

You need Go 1.22+:

```bash
cd src
go vet ./...
GOOS=linux GOARCH=amd64 go build -o ../linux/bios-reader-linux .
GOOS=windows GOARCH=amd64 go build -o ../windows/bios-reader.exe .
```

The original Python / PowerShell scripts are still in `src/` for reference.

## Safety & responsibility

- Read-only: it never writes/flashes anything. Brick risk is minimal, but still **no warranty**.
- Please don't use it on stolen laptops, office inventory without IT approval, or to bypass someone else's security.
- Repo owner is **not responsible** for community forks or misuse for crime.

## License / Copyright

© 2026 — Free to use, modify, and share, **on one condition: DON'T use it for crime / unauthorized access**.
See `LICENSE` for the full text. Your mods = your responsibility.

---
PRs welcome: new vendor GUIDs, new hash parsers, or better translations. 🙌
