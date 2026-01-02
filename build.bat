@echo off
echo Building Zig Kernel...
pushd internal\zig
zig build-lib main.zig -dynamic -O ReleaseFast
if %errorlevel% neq 0 (
    echo Zig build failed!
    popd
    exit /b %errorlevel%
)
move /Y main.dll ..\..\ssn_kernel.dll >nul
del main.lib main.pdb
popd

echo Deploying Kernel to cmd examples...
copy /Y ssn_kernel.dll cmd\rps-context\ssn_kernel.dll >nul
copy /Y ssn_kernel.dll cmd\rps-noise\ssn_kernel.dll >nul
copy /Y ssn_kernel.dll cmd\rps-10k\ssn_kernel.dll >nul

echo Build and Deployment Complete.
