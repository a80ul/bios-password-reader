@echo off
REM BIOS Password Reader - Windows launcher (alternatif jika .exe diblokir)
REM Klik kanan > Run as Administrator
REM Hanya untuk laptop milik sendiri / seizin pemilik.

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0bios_pw_read_windows.ps1"
pause
