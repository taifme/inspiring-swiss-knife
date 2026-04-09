# build.ps1 - Build Inspiring Swiss Knife for Windows
# Run from the project root:  .\build.ps1

param(
    [string]$Output = "inspiring-swiss-knife.exe",
    [switch]$Release
)

$ErrorActionPreference = "Stop"

Write-Host "⚔  Building Inspiring Swiss Knife..." -ForegroundColor Cyan

# Ensure Go is in PATH
$goPaths = @(
    "C:\Program Files\Go\bin",
    "$env:LOCALAPPDATA\Programs\Go\bin",
    "$env:GOPATH\bin"
)
foreach ($p in $goPaths) {
    if (Test-Path "$p\go.exe") {
        $env:PATH = "$p;$env:PATH"
        break
    }
}

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Error "Go not found. Install from https://go.dev/dl/ or run: winget install GoLang.Go"
    exit 1
}

Write-Host "Go version: $(go version)" -ForegroundColor Gray

# Tidy dependencies
Write-Host "  → Tidying modules..." -ForegroundColor Gray
go mod tidy

# Build flags
$ldflags = "-s -w"  # strip debug info for smaller binary
$buildArgs = @("build", "-ldflags=$ldflags", "-o", $Output, ".")

if ($Release) {
    Write-Host "  → Building RELEASE (optimized)..." -ForegroundColor Green
    $env:GOFLAGS = "-trimpath"
} else {
    Write-Host "  → Building DEBUG..." -ForegroundColor Yellow
}

& go @buildArgs

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed!"
    exit 1
}

$size = [math]::Round((Get-Item $Output).Length / 1MB, 2)
Write-Host "✓ Built: $Output ($size MB)" -ForegroundColor Green
Write-Host ""
Write-Host "Run it:  .\$Output" -ForegroundColor Cyan
Write-Host "Note:    Run as Administrator for tweaks to work correctly." -ForegroundColor Yellow
