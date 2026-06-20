$ErrorActionPreference = "Stop"

Write-Host "=========================================="
Write-Host " Building Aether Daemon Alpha 0.01 "
Write-Host "=========================================="

$RootDir = Get-Location
$SrcDir = ".\src\daemon-go"
$OutDir = "$RootDir\bin"

if (!(Test-Path -Path $OutDir)) {
    New-Item -ItemType Directory -Path $OutDir | Out-Null
}

Set-Location $SrcDir

Write-Host "Building for Windows (amd64)..."
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "$OutDir\aether-daemon-win-amd64.exe" ".\cmd\aether-daemon"

Write-Host "Building for macOS (M1/arm64)..."
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -ldflags="-s -w" -o "$OutDir\aether-daemon-mac-arm64" ".\cmd\aether-daemon"

Set-Location $RootDir

Write-Host "=========================================="
Write-Host " Build Complete! Binaries are in .\bin\"
Write-Host "=========================================="
