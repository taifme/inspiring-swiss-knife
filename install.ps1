# ─────────────────────────────────────────────────────────────────────────────
# Inspiring Swiss Knife — One-Line Installer
#
# Usage (run as Administrator in PowerShell):
#   irm "https://isk.inspiringlivingsolutions.com/win" | iex
#   -- or --
#   irm "https://raw.githubusercontent.com/taifme/inspiring-swiss-knife/main/install.ps1" | iex
#
# What this script does:
#   1. Checks that PowerShell 5+ and internet access are available
#   2. Downloads the latest isk.exe from GitHub Releases
#   3. Verifies the SHA-256 checksum (when published alongside the release)
#   4. Launches isk.exe in the current terminal
#   5. Cleans up the temp file on exit
# ─────────────────────────────────────────────────────────────────────────────

#Requires -Version 5
$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$RepoOwner  = 'taifme'
$RepoName   = 'inspiring-swiss-knife'
$ApiUrl     = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"

function Write-Banner {
    $green  = [char]27 + '[92m'
    $cyan   = [char]27 + '[96m'
    $yellow = [char]27 + '[93m'
    $reset  = [char]27 + '[0m'
    Write-Host ""
    Write-Host "${cyan} ___  ________  ___  __${reset}"
    Write-Host "${cyan}|\  \|\   ____\|\  \|\  \${reset}"
    Write-Host "${cyan}\ \  \ \  \___|\ \  \/  /|_${reset}"
    Write-Host "${cyan} \ \  \ \_____  \ \   ___  \${reset}"
    Write-Host "${cyan}  \ \  \|____|\  \ \  \\ \  \${reset}"
    Write-Host "${cyan}   \ \__\____\_\  \ \__\\ \__\${reset}"
    Write-Host "${cyan}    \|__|\_________\|__| \|__|${reset}"
    Write-Host "${cyan}        \|_________|${reset}"
    Write-Host ""
    Write-Host "${cyan}  Inspiring Swiss Knife — Windows Onboarding & Optimization Tool${reset}"
    Write-Host "${yellow}  github.com/$RepoOwner/$RepoName${reset}"
    Write-Host ""
}

function Get-LatestRelease {
    Write-Host "  Fetching latest release info..." -ForegroundColor Cyan
    try {
        $headers = @{ 'User-Agent' = 'ISK-Installer/1.0' }
        $release = Invoke-RestMethod -Uri $ApiUrl -Headers $headers -UseBasicParsing
        return $release
    }
    catch {
        Write-Host "  ERROR: Could not reach GitHub API. Check your internet connection." -ForegroundColor Red
        Write-Host "  Manual download: https://github.com/$RepoOwner/$RepoName/releases" -ForegroundColor Yellow
        exit 1
    }
}

function Find-Asset($release) {
    # Try Windows AMD64 first
    $asset = $release.assets | Where-Object {
        $_.name -match '(?i)(windows|win).*(amd64|x86_64|x64).*\.exe$' -or
        $_.name -match '(?i)(amd64|x86_64|x64).*(windows|win).*\.exe$'
    } | Select-Object -First 1

    # Fallback: any .exe asset
    if (-not $asset) {
        $asset = $release.assets | Where-Object { $_.name -match '\.exe$' } | Select-Object -First 1
    }

    if (-not $asset) {
        Write-Host "  ERROR: No Windows executable found in release $($release.tag_name)." -ForegroundColor Red
        Write-Host "  Manual download: $($release.html_url)" -ForegroundColor Yellow
        exit 1
    }
    return $asset
}

function Get-Checksum($release, $assetName) {
    $sha = $release.assets | Where-Object { $_.name -eq "$assetName.sha256" } | Select-Object -First 1
    if (-not $sha) {
        $sha = $release.assets | Where-Object { $_.name -match 'sha256' } | Select-Object -First 1
    }
    return $sha
}

function Download-File($url, $dest) {
    Write-Host "  Downloading $(Split-Path $dest -Leaf)..." -ForegroundColor Cyan
    $headers = @{ 'User-Agent' = 'ISK-Installer/1.0' }
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $url -OutFile $dest -Headers $headers -UseBasicParsing
    $ProgressPreference = 'Continue'
}

function Verify-Checksum($filePath, $checksumUrl) {
    try {
        $headers = @{ 'User-Agent' = 'ISK-Installer/1.0' }
        $expected = (Invoke-RestMethod -Uri $checksumUrl -Headers $headers -UseBasicParsing).Trim().Split()[0].ToLower()
        $actual   = (Get-FileHash -Path $filePath -Algorithm SHA256).Hash.ToLower()
        if ($expected -ne $actual) {
            Write-Host "  ERROR: Checksum mismatch! Download may be corrupted." -ForegroundColor Red
            Write-Host "  Expected: $expected" -ForegroundColor Red
            Write-Host "  Got:      $actual" -ForegroundColor Red
            exit 1
        }
        Write-Host "  Checksum verified." -ForegroundColor Green
    }
    catch {
        Write-Host "  WARNING: Could not verify checksum — proceeding anyway." -ForegroundColor Yellow
    }
}

# ── Main ──────────────────────────────────────────────────────────────────────

Write-Banner

$release  = Get-LatestRelease
$asset    = Find-Asset $release
$tempDir  = $env:TEMP
$destPath = Join-Path $tempDir $asset.name

Write-Host "  Version:  $($release.tag_name)" -ForegroundColor Green
Write-Host "  Asset:    $($asset.name)" -ForegroundColor Green
Write-Host "  Size:     $([math]::Round($asset.size / 1MB, 1)) MB" -ForegroundColor Green
Write-Host ""

Download-File $asset.browser_download_url $destPath

# Verify checksum if available
$sha = Get-Checksum $release $asset.name
if ($sha) {
    Verify-Checksum $destPath $sha.browser_download_url
}

Write-Host ""
Write-Host "  Launching Inspiring Swiss Knife..." -ForegroundColor Green
Write-Host "  (Close the window or press [q] to quit)" -ForegroundColor DarkGray
Write-Host ""

# Run and wait
& $destPath

# Cleanup
Remove-Item $destPath -ErrorAction SilentlyContinue
