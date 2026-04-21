# install.ps1 - install the prebuilt cursor-sync.exe from bin\windows_<arch>\.
#
# Usage:
#   .\install.ps1
#
# Env vars:
#   $env:INSTALL_DIR   install location (default: $env:USERPROFILE\.cursor-sync\bin)

$ErrorActionPreference = "Stop"

$ScriptDir = $PSScriptRoot
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { Join-Path $env:USERPROFILE ".cursor-sync\bin" }

switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { $arch = "amd64" }
    "ARM64" { $arch = "arm64" }
    default { throw "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)" }
}

$src = Join-Path $ScriptDir "bin\windows_$arch\cursor-sync.exe"
if (-not (Test-Path $src)) {
    throw "No prebuilt binary for windows_$arch at $src.`nRun .\build.sh (on a Unix machine with Go) to build it, then re-run .\install.ps1."
}

New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
$dest = Join-Path $InstallDir "cursor-sync.exe"
Copy-Item -Path $src -Destination $dest -Force

Write-Host ""
Write-Host "Installed cursor-sync to $dest"

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not ($userPath -split ";" | Where-Object { $_ -eq $InstallDir })) {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to your user PATH (open a new shell to pick it up)."
}
Write-Host "Run: cursor-sync --help"
