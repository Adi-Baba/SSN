@echo off
echo [BUILD] Building SSN Kernel (Zig)...
cd internal\zig
zig build-lib main.zig -dynamic -O ReleaseFast -femit-bin=../../cmd/rps-context/ssn_kernel.dll
if %errorlevel% neq 0 (
    echo [ERROR] Zig build failed!
    exit /b %errorlevel%
)

cd ..\..
echo [BUILD] Success. Kernel DLL deployed to:
echo   - cmd/rps-context/
copy cmd\rps-context\ssn_kernel.dll cmd\rps-10k\ssn_kernel.dll >nul
echo   - cmd/rps-10k/
copy cmd\rps-context\ssn_kernel.dll cmd\rps-noise\ssn_kernel.dll >nul
echo   - cmd/rps-noise/
