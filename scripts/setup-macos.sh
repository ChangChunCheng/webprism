#!/bin/bash
# WEBPRISM macOS Environment Setup Script

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🍎 WEBPRISM macOS 環境設置${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo -e "${YELLOW}⚠️  Homebrew 未安裝${NC}"
    echo ""
    echo "Homebrew 是 macOS 上最流行的套件管理工具"
    echo ""
    read -p "是否要安裝 Homebrew? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        echo -e "${BLUE}📦 安裝 Homebrew...${NC}"
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

        # Add Homebrew to PATH for Apple Silicon Macs
        if [[ $(uname -m) == 'arm64' ]]; then
            echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
            eval "$(/opt/homebrew/bin/brew shellenv)"
        fi

        echo -e "${GREEN}✓ Homebrew 安裝完成${NC}"
    else
        echo -e "${RED}✗ 需要 Homebrew 才能繼續${NC}"
        echo "手動安裝指令: /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
        exit 1
    fi
fi

echo -e "${GREEN}✓ Homebrew 已安裝${NC}"
echo ""

# Update Homebrew
echo -e "${BLUE}🔄 更新 Homebrew...${NC}"
brew update

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📦 檢查必要工具${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

TOOLS_TO_INSTALL=()

# Check Go
if command -v go &> /dev/null; then
    echo -e "${GREEN}✓ Go 已安裝${NC} ($(go version))"
else
    echo -e "${YELLOW}⚪ Go 未安裝${NC}"
    TOOLS_TO_INSTALL+=("go")
fi

# Check protoc
if command -v protoc &> /dev/null; then
    echo -e "${GREEN}✓ protoc 已安裝${NC} ($(protoc --version))"
else
    echo -e "${YELLOW}⚪ protoc 未安裝${NC}"
    TOOLS_TO_INSTALL+=("protobuf")
fi

# Check Docker
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓ Docker 已安裝${NC} ($(docker --version))"
else
    echo -e "${YELLOW}⚪ Docker 未安裝${NC}"
    echo -e "${BLUE}   提示: Docker Desktop 需要手動下載安裝${NC}"
    echo -e "${BLUE}   下載: https://www.docker.com/products/docker-desktop/${NC}"
fi

# Check PostgreSQL (optional)
if command -v psql &> /dev/null; then
    echo -e "${GREEN}✓ PostgreSQL 已安裝${NC} ($(psql --version))"
else
    echo -e "${YELLOW}⚪ PostgreSQL 未安裝 (可選)${NC}"
    read -p "   是否要安裝 PostgreSQL? (建議使用 Docker) (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        TOOLS_TO_INSTALL+=("postgresql@16")
    fi
fi

# Check optional tools
echo ""
echo "可選開發工具:"

if command -v golangci-lint &> /dev/null; then
    echo -e "${GREEN}✓ golangci-lint 已安裝${NC}"
else
    echo -e "${YELLOW}⚪ golangci-lint 未安裝 (可選)${NC}"
    read -p "   是否要安裝? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        TOOLS_TO_INSTALL+=("golangci-lint")
    fi
fi

if command -v migrate &> /dev/null; then
    echo -e "${GREEN}✓ golang-migrate 已安裝${NC}"
else
    echo -e "${YELLOW}⚪ golang-migrate 未安裝 (可選)${NC}"
    read -p "   是否要安裝? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        TOOLS_TO_INSTALL+=("golang-migrate")
    fi
fi

# Install tools if needed
if [ ${#TOOLS_TO_INSTALL[@]} -gt 0 ]; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}📦 安裝工具${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "即將安裝: ${TOOLS_TO_INSTALL[*]}"
    echo ""

    for tool in "${TOOLS_TO_INSTALL[@]}"; do
        echo -e "${BLUE}安裝 $tool...${NC}"
        brew install "$tool"
    done

    echo ""
    echo -e "${GREEN}✓ 工具安裝完成${NC}"
fi

# Install Go tools
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}🔧 安裝 Go 開發工具${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if command -v go &> /dev/null; then
    echo "安裝 protoc 插件和 mockery 到 GOPATH..."
    echo ""

    # Get GOPATH
    GOPATH=$(go env GOPATH)
    echo "GOPATH: $GOPATH"
    echo ""

    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
    go install github.com/vektra/mockery/v2@latest

    echo ""
    echo -e "${GREEN}✓ Go 工具安裝完成${NC}"

    # Check if GOPATH/bin is in PATH
    if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
        echo ""
        echo -e "${YELLOW}⚠️  警告: GOPATH/bin 不在您的 PATH 中${NC}"
        echo ""
        echo "請將以下內容加入您的 shell 配置檔 (~/.zshrc 或 ~/.bashrc):"
        echo ""
        echo -e "${BLUE}export PATH=\"\$PATH:$GOPATH/bin\"${NC}"
        echo ""
        echo "然後執行: source ~/.zshrc  (或 source ~/.bashrc)"
    fi
else
    echo -e "${RED}✗ Go 未安裝，無法安裝 Go 工具${NC}"
    echo "請先安裝 Go 後再執行此腳本"
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ macOS 環境設置完成！${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "下一步:"
echo "  1. 驗證安裝: ./scripts/check-tools.sh"
echo "  2. 設置專案: make setup"
echo "  3. 啟動服務: make docker-up"
echo ""
echo "💡 提示:"
echo "  - 如未安裝 Docker Desktop，請從官網下載："
echo "    https://www.docker.com/products/docker-desktop/"
echo "  - 使用 Docker 管理 PostgreSQL 更方便 (無需本機安裝)"
echo "  - 詳細文件: DEPENDENCIES.md"
echo ""
