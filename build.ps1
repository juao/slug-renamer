# build.ps1 - Script de build completo para Windows
# Executar no PowerShell: .\build.ps1
#
# Pre-requisitos:
#   1. Go 1.22+        -> https://go.dev/dl/
#   2. MinGW (windres) -> choco install mingw
#   3. NSIS 3.x        -> choco install nsis

$ErrorActionPreference = "Stop"
$AppName   = "slug-renamer"
$OutputDir = ".\dist"

Write-Host ""
Write-Host "Slug Renamer - Build Script" -ForegroundColor Cyan
Write-Host "---------------------------------" -ForegroundColor DarkGray

New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null

Write-Host "[1/4] Baixando dependencias..." -ForegroundColor Yellow
go mod tidy

Write-Host "[2/4] Compilando manifesto e icone (windres)..." -ForegroundColor Yellow
Push-Location .\cmd
try {
    windres slug-renamer.rc -O coff -o slug-renamer.syso
    if ($LASTEXITCODE -ne 0) { throw "windres falhou. Verifique se o MinGW esta no PATH." }
    Write-Host "  OK: cmd\slug-renamer.syso" -ForegroundColor Green
} finally {
    Pop-Location
}

Write-Host "[3/4] Compilando slug-renamer.exe..." -ForegroundColor Yellow
$env:GOOS   = "windows"
$env:GOARCH = "amd64"
go build -o "$OutputDir\$AppName.exe" -ldflags "-s -w -H windowsgui" .\cmd\...
Write-Host "  OK: $OutputDir\$AppName.exe" -ForegroundColor Green

Write-Host "[4/4] Gerando instalador com NSIS..." -ForegroundColor Yellow
$env:PATH += ";C:\Program Files (x86)\NSIS"
Copy-Item "$OutputDir\$AppName.exe" ".\installer\$AppName.exe" -Force
makensis .\installer\installer.nsi
Move-Item ".\installer\$AppName-setup.exe" "$OutputDir\$AppName-setup.exe" -Force
Remove-Item ".\installer\$AppName.exe" -ErrorAction SilentlyContinue
Write-Host "  OK: $OutputDir\$AppName-setup.exe" -ForegroundColor Green

Write-Host ""
Write-Host "Build concluido! Arquivos em: $OutputDir" -ForegroundColor Green
