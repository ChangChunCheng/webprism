#!/bin/bash
# WEBPRISM 工具檢查腳本 - 支援跨平台 (Windows/macOS/Ubuntu)

# 顏色定義
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 取得腳本目錄
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 載入平台檢測
source "$SCRIPT_DIR/detect-platform.sh" 2>/dev/null || true

echo "🔍 檢查 WEBPRISM 開發工具..."
echo ""

# 顯示平台資訊
OS=$(detect_os 2>/dev/null || echo "unknown")
echo -e "${BLUE}Platform:${NC} $OS ($(uname -m))"
echo ""

ALL_OK=true

check_command() {
    local cmd=$1
    local name=${2:-$1}

    if command -v "$cmd" &> /dev/null; then
        echo -e "${GREEN}✓${NC} $name"
        if [ "$3" = "version" ]; then
            $cmd --version 2>&1 | head -1 | sed 's/^/  /'
        fi
    else
        echo -e "${RED}✗${NC} $name ${YELLOW}(未安裝)${NC}"
        ALL_OK=false
    fi
}

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "必要工具"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

check_command protoc "Protocol Buffers 編譯器" version
check_command go "Go 語言" version

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Go 插件"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

check_command protoc-gen-go "protoc-gen-go"
check_command protoc-gen-go-grpc "protoc-gen-go-grpc"
check_command protoc-gen-grpc-gateway "protoc-gen-grpc-gateway"
check_command protoc-gen-openapiv2 "protoc-gen-openapiv2"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "開發工具"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

check_command golangci-lint "golangci-lint" version
check_command migrate "golang-migrate"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "資料庫（可選）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

check_command psql "PostgreSQL" version

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ "$ALL_OK" = true ]; then
    echo -e "${GREEN}✓ 所有必要工具都已安裝！${NC}"
    echo ""
    echo "下一步："
    echo "  1. 執行 'make proto' 生成 Protobuf 程式碼"
    echo "  2. 執行 'go mod tidy' 下載依賴"
    echo "  3. 執行 'make build' 編譯專案"
else
    echo -e "${RED}✗ 有些工具尚未安裝${NC}"
    echo ""
    echo "請參考以下安裝指南："
    echo ""

    case "$OS" in
        macos)
            echo -e "${BLUE}📱 macOS 安裝指令:${NC}"
            echo "  # 使用 Homebrew 安裝"
            echo "  brew install go protobuf postgresql@16 golangci-lint golang-migrate"
            echo ""
            echo "  # 安裝 Go 插件"
            echo "  make install-tools"
            echo ""
            echo "  # 或執行一鍵安裝腳本"
            echo "  ./scripts/setup-macos.sh"
            ;;
        ubuntu)
            echo -e "${BLUE}🐧 Ubuntu/Debian 安裝指令:${NC}"
            echo "  # 執行 Ubuntu 安裝腳本"
            echo "  ./scripts/setup-ubuntu.sh"
            echo ""
            echo "  # 或手動安裝"
            echo "  sudo apt-get update"
            echo "  sudo apt-get install -y protobuf-compiler postgresql-16"
            echo "  make install-tools"
            ;;
        windows)
            echo -e "${BLUE}🪟 Windows 安裝指令:${NC}"
            echo "  # 建議使用 WSL2 並按 Ubuntu 方式安裝"
            echo "  # 或在 PowerShell 中執行："
            echo "  .\\scripts\\setup-windows.ps1"
            echo ""
            echo "  # 手動安裝連結："
            echo "  # Go: https://go.dev/dl/"
            echo "  # protoc: https://github.com/protocolbuffers/protobuf/releases"
            echo "  # PostgreSQL: https://www.postgresql.org/download/windows/"
            ;;
        *)
            echo -e "${YELLOW}⚠️  未識別的作業系統${NC}"
            echo "  請參考 DEPENDENCIES.md 手動安裝"
            ;;
    esac

    echo ""
    echo "詳細說明請參考: DEPENDENCIES.md"
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
