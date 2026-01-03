# WEBPRISM Windows Environment Setup Script
# PowerShell script for Windows users

Write-Host "🪟 WEBPRISM Windows 環境設置" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Host ""

# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "⚠️  警告: 建議以系統管理員身份執行此腳本" -ForegroundColor Yellow
    Write-Host ""
}

# Function to check if a command exists
function Test-CommandExists {
    param($command)
    $null = Get-Command $command -ErrorAction SilentlyContinue
    return $?
}

# Function to check and install using winget (Windows Package Manager)
function Install-WithWinget {
    param($packageId, $name)

    if (Test-CommandExists "winget") {
        Write-Host "📦 安裝 $name..." -ForegroundColor Yellow
        winget install --id=$packageId -e --silent
        return $true
    }
    return $false
}

Write-Host "💡 建議使用方式:" -ForegroundColor Green
Write-Host ""
Write-Host "方案 1: 使用 WSL2 (推薦)" -ForegroundColor Cyan
Write-Host "  - 安裝 WSL2: wsl --install"
Write-Host "  - 安裝 Ubuntu: wsl --install -d Ubuntu"
Write-Host "  - 在 WSL2 中執行: ./scripts/setup-ubuntu.sh"
Write-Host ""
Write-Host "方案 2: 使用 Docker Desktop (推薦)" -ForegroundColor Cyan
Write-Host "  - 只需安裝 Go 和 Docker Desktop"
Write-Host "  - 使用 Docker 運行 PostgreSQL"
Write-Host "  - 執行此腳本進行基本設置"
Write-Host ""
Write-Host "方案 3: 完整原生 Windows 安裝" -ForegroundColor Cyan
Write-Host "  - 安裝所有工具到 Windows"
Write-Host "  - 執行此腳本進行完整設置"
Write-Host ""

$choice = Read-Host "請選擇安裝方式 (1/2/3) 或按 Enter 繼續檢查當前環境"

if ($choice -eq "1") {
    Write-Host ""
    Write-Host "📝 WSL2 安裝步驟:" -ForegroundColor Green
    Write-Host "  1. 以系統管理員身份開啟 PowerShell"
    Write-Host "  2. 執行: wsl --install"
    Write-Host "  3. 重新啟動電腦"
    Write-Host "  4. 執行: wsl --install -d Ubuntu"
    Write-Host "  5. 設置 Ubuntu 使用者名稱和密碼"
    Write-Host "  6. 在 WSL2 中執行: cd /mnt/c/path/to/webprism && ./scripts/setup-ubuntu.sh"
    exit 0
}

Write-Host ""
Write-Host "🔍 檢查已安裝的工具..." -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Host ""

$allOk = $true

# Check Go
if (Test-CommandExists "go") {
    $goVersion = go version
    Write-Host "✓ Go: $goVersion" -ForegroundColor Green
} else {
    Write-Host "✗ Go 未安裝" -ForegroundColor Red
    $allOk = $false
}

# Check protoc
if (Test-CommandExists "protoc") {
    $protocVersion = protoc --version
    Write-Host "✓ protoc: $protocVersion" -ForegroundColor Green
} else {
    Write-Host "✗ protoc 未安裝" -ForegroundColor Red
    $allOk = $false
}

# Check Docker
if (Test-CommandExists "docker") {
    try {
        $dockerVersion = docker --version
        Write-Host "✓ Docker: $dockerVersion" -ForegroundColor Green
    } catch {
        Write-Host "✗ Docker 已安裝但未運行" -ForegroundColor Yellow
    }
} else {
    Write-Host "✗ Docker 未安裝" -ForegroundColor Yellow
}

# Check PostgreSQL (optional)
if (Test-CommandExists "psql") {
    $pgVersion = psql --version
    Write-Host "✓ PostgreSQL: $pgVersion" -ForegroundColor Green
} else {
    Write-Host "⚪ PostgreSQL 未安裝 (可選，可使用 Docker)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if (-not $allOk) {
    Write-Host ""
    Write-Host "📦 需要安裝以下工具:" -ForegroundColor Yellow
    Write-Host ""

    if (Test-CommandExists "winget") {
        Write-Host "使用 winget 自動安裝? (需要系統管理員權限)" -ForegroundColor Cyan
        $install = Read-Host "繼續? (y/N)"

        if ($install -eq "y" -or $install -eq "Y") {
            if (-not (Test-CommandExists "go")) {
                Install-WithWinget "GoLang.Go" "Go"
            }
            if (-not (Test-CommandExists "docker")) {
                Install-WithWinget "Docker.DockerDesktop" "Docker Desktop"
            }

            Write-Host ""
            Write-Host "⚠️  protoc 需要手動安裝:" -ForegroundColor Yellow
            Write-Host "   下載: https://github.com/protocolbuffers/protobuf/releases"
            Write-Host "   解壓並將 bin 目錄加入 PATH"
        }
    } else {
        Write-Host "請手動安裝以下工具:" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Go (1.21+):" -ForegroundColor White
        Write-Host "  下載: https://go.dev/dl/"
        Write-Host "  安裝後重新開啟 PowerShell"
        Write-Host ""
        Write-Host "Protocol Buffers:" -ForegroundColor White
        Write-Host "  下載: https://github.com/protocolbuffers/protobuf/releases"
        Write-Host "  解壓 protoc-xx.x-win64.zip"
        Write-Host "  將 bin 目錄加入系統 PATH"
        Write-Host ""
        Write-Host "Docker Desktop (推薦):" -ForegroundColor White
        Write-Host "  下載: https://www.docker.com/products/docker-desktop/"
        Write-Host ""
        Write-Host "PostgreSQL (可選):" -ForegroundColor White
        Write-Host "  下載: https://www.postgresql.org/download/windows/"
        Write-Host "  或使用 Docker: docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:16"
    }
} else {
    Write-Host ""
    Write-Host "✓ 所有必要工具都已安裝！" -ForegroundColor Green
    Write-Host ""
    Write-Host "下一步:" -ForegroundColor Cyan
    Write-Host "  1. 安裝 Go 工具: make install-tools"
    Write-Host "  2. 設置專案: make setup"
    Write-Host "  3. 啟動服務: make docker-up"
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
Write-Host ""
Write-Host "💡 提示:" -ForegroundColor Green
Write-Host "  - Windows 環境推薦使用 WSL2 或 Docker"
Write-Host "  - 詳細文件: DEPENDENCIES.md"
Write-Host "  - 問題回報: https://github.com/ChangChunCheng/webprism/issues"
