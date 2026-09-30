# build.ps1 — builds the Go companion with the version from package.json
$ErrorActionPreference = "Stop"

# Path to the frontend's package.json — adjust to your repo layout
$pkgPath = ".\EFB\PackageSources\StreamDeckVR\package.json"

# Read version from package.json
$version = (Get-Content $pkgPath | ConvertFrom-Json).version

Write-Host "Building StreamDeckVR companion v$version..." -ForegroundColor Cyan

# Inject version via ldflags and build
go build -ldflags "-X main.Version=$version" -o streamdeckvr.exe .

if ($LASTEXITCODE -eq 0) {
    Write-Host "OK: streamdeckvr.exe v$version" -ForegroundColor Green

    # Verify it took (prints e.g. "main.Version=0.0.2" or blank if var not wired)
    go version -m streamdeckvr.exe | Select-String "build" | Out-Null
} else {
    Write-Host "Build FAILED" -ForegroundColor Red
    exit 1
}