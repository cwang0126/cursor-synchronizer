# install.ps1 - download the latest cursor-sync binary for Windows.
#
# Usage:
#   irm https://raw.githubusercontent.com/cwang0126/cursor-synchronizer/main/install.ps1 | iex
#
# Env vars:
#   $env:REPO        owner/name (default: cwang0126/cursor-synchronizer)
#   $env:VERSION     release tag (default: latest)
#   $env:INSTALL_DIR install location (default: $env:USERPROFILE\.cursor-sync\bin)

$ErrorActionPreference = "Stop"

$Repo = if ($env:REPO) { $env:REPO } else { "cwang0126/cursor-synchronizer" }
$Version = if ($env:VERSION) { $env:VERSION } else { "latest" }
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { Join-Path $env:USERPROFILE ".cursor-sync\bin" }

switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { $goarch = "amd64" }
    "ARM64" { $goarch = "arm64" }
    default { throw "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)" }
}
$goos = "windows"

if ($Version -eq "latest") {
    $resp = Invoke-WebRequest "https://github.com/$Repo/releases/latest" -MaximumRedirection 0 -ErrorAction SilentlyContinue
    if ($resp.StatusCode -in 301, 302) {
        $Version = ($resp.Headers.Location -split "/tag/")[-1]
    } else {
        $Version = ((Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest").tag_name)
    }
    $UrlPrefix = "https://github.com/$Repo/releases/download/$Version"
} else {
    $UrlPrefix = "https://github.com/$Repo/releases/download/$Version"
}

$asset = "cursor-sync_${Version}_${goos}_${goarch}.zip"
$url = "$UrlPrefix/$asset"

Write-Host "Downloading $asset..."
$tmp = Join-Path $env:TEMP ("cursor-sync-" + [guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null
try {
    $zipPath = Join-Path $tmp $asset
    Invoke-WebRequest -Uri $url -OutFile $zipPath
    Expand-Archive -Path $zipPath -DestinationPath $tmp -Force

    $stage = "cursor-sync_${Version}_${goos}_${goarch}"
    $binSrc = Join-Path $tmp "$stage\cursor-sync.exe"

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    $dest = Join-Path $InstallDir "cursor-sync.exe"
    Copy-Item -Path $binSrc -Destination $dest -Force

    Write-Host ""
    Write-Host "Installed cursor-sync $Version to $dest"

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if (-not ($userPath -split ";" | Where-Object { $_ -eq $InstallDir })) {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        Write-Host "Added $InstallDir to your user PATH (open a new shell to pick it up)."
    }
    Write-Host "Run: cursor-sync --help"
}
finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
