#!/bin/bash
# WEBPRISM Ubuntu/Debian Environment Setup Script
# Supports Ubuntu 22.04+

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐧 WEBPRISM Ubuntu 環境設置${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check Ubuntu version
if [ -f /etc/os-release ]; then
    . /etc/os-release
    echo -e "系統: ${GREEN}$PRETTY_NAME${NC}"

    # Extract version number
    VERSION_NUM=$(echo $VERSION_ID | cut -d. -f1)
    if [ "$VERSION_NUM" -lt 22 ]; then
        echo -e "${YELLOW}⚠️  警告: 建議使用 Ubuntu 22.04 或更高版本${NC}"
        echo ""
    fi
else
    echo -e "${YELLOW}⚠️  無法檢測系統版本${NC}"
fi

echo ""

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  請不要使用 root 執行此腳本${NC}"
    echo "正確方式: ./scripts/setup-ubuntu.sh"
    exit 1
fi

# Update package list
echo -e "${BLUE}📦 更新套件列表...${NC}"
sudo apt-get update

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📦 檢查必要工具${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

PACKAGES_TO_INSTALL=()

# Check Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓ Go 已安裝${NC} ($GO_VERSION)"

    # Check Go version
    GO_VER_NUM=$(echo $GO_VERSION | sed 's/go//' | cut -d. -f2)
    if [ "$GO_VER_NUM" -lt 21 ]; then
        echo -e "${YELLOW}  ⚠️  建議升級到 Go 1.21 或更高版本${NC}"
    fi
else
    echo -e "${YELLOW}⚪ Go 未安裝${NC}"
    echo "   將安裝 Go 1.21+"
fi

# Check protoc
if command -v protoc &> /dev/null; then
    echo -e "${GREEN}✓ protoc 已安裝${NC} ($(protoc --version))"
else
    echo -e "${YELLOW}⚪ protoc 未安裝${NC}"
    PACKAGES_TO_INSTALL+=("protobuf-compiler")
fi

# Check git
if command -v git &> /dev/null; then
    echo -e "${GREEN}✓ git 已安裝${NC}"
else
    echo -e "${YELLOW}⚪ git 未安裝${NC}"
    PACKAGES_TO_INSTALL+=("git")
fi

# Check make
if command -v make &> /dev/null; then
    echo -e "${GREEN}✓ make 已安裝${NC}"
else
    echo -e "${YELLOW}⚪ make 未安裝${NC}"
    PACKAGES_TO_INSTALL+=("build-essential")
fi

# Check curl
if command -v curl &> /dev/null; then
    echo -e "${GREEN}✓ curl 已安裝${NC}"
else
    echo -e "${YELLOW}⚪ curl 未安裝${NC}"
    PACKAGES_TO_INSTALL+=("curl")
fi

# Check Docker
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✓ Docker 已安裝${NC} ($(docker --version))"

    # Check if user is in docker group
    if groups | grep -q docker; then
        echo -e "${GREEN}  ✓ 使用者在 docker 群組中${NC}"
    else
        echo -e "${YELLOW}  ⚠️  使用者不在 docker 群組中${NC}"
        echo "     執行以下命令加入 docker 群組："
        echo "     sudo usermod -aG docker $USER"
        echo "     然後登出重新登入"
    fi
else
    echo -e "${YELLOW}⚪ Docker 未安裝${NC}"
    echo "   將提供 Docker 安裝指引"
fi

# Check PostgreSQL (optional)
echo ""
echo "可選工具:"
if command -v psql &> /dev/null; then
    echo -e "${GREEN}✓ PostgreSQL 已安裝${NC} ($(psql --version))"
else
    echo -e "${YELLOW}⚪ PostgreSQL 未安裝 (可選)${NC}"
    read -p "   是否要安裝 PostgreSQL 16? (建議使用 Docker) (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        PACKAGES_TO_INSTALL+=("postgresql-16" "postgresql-client-16")
    fi
fi

# Install packages if needed
if [ ${#PACKAGES_TO_INSTALL[@]} -gt 0 ]; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}📦 安裝套件${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "即將安裝: ${PACKAGES_TO_INSTALL[*]}"
    echo ""

    # Add PostgreSQL repository if needed
    if [[ " ${PACKAGES_TO_INSTALL[@]} " =~ " postgresql-16 " ]]; then
        echo -e "${BLUE}加入 PostgreSQL 官方倉庫...${NC}"
        sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
        wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
        sudo apt-get update
    fi

    sudo apt-get install -y "${PACKAGES_TO_INSTALL[@]}"
    echo ""
    echo -e "${GREEN}✓ 套件安裝完成${NC}"
fi

# Install Go if not present
if ! command -v go &> /dev/null; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}📦 安裝 Go 1.21+${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    GO_VERSION="1.21.6"
    GO_ARCH="amd64"

    # Detect architecture
    if [ "$(uname -m)" = "aarch64" ]; then
        GO_ARCH="arm64"
    fi

    GO_TARBALL="go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    GO_URL="https://go.dev/dl/${GO_TARBALL}"

    echo "下載 Go ${GO_VERSION}..."
    wget -q --show-progress "$GO_URL"

    echo "安裝 Go..."
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "$GO_TARBALL"
    rm "$GO_TARBALL"

    # Add Go to PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
    fi

    export PATH=$PATH:/usr/local/go/bin

    echo -e "${GREEN}✓ Go 安裝完成${NC}"
    echo ""
    echo -e "${YELLOW}⚠️  請執行以下命令以更新 PATH:${NC}"
    echo "   source ~/.bashrc"
fi

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}🐳 安裝 Docker${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    read -p "是否要安裝 Docker? (推薦) (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "安裝 Docker..."

        # Install Docker using official script
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        rm get-docker.sh

        # Add user to docker group
        sudo usermod -aG docker $USER

        echo ""
        echo -e "${GREEN}✓ Docker 安裝完成${NC}"
        echo ""
        echo -e "${YELLOW}⚠️  請登出並重新登入以使 docker 群組生效${NC}"
        echo "   或執行: newgrp docker"
    else
        echo -e "${YELLOW}⚪ 跳過 Docker 安裝${NC}"
    fi
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
        echo -e "${YELLOW}⚠️  GOPATH/bin 不在您的 PATH 中${NC}"
        echo ""
        echo "已自動加入到 ~/.bashrc"
        echo "請執行: source ~/.bashrc"
    fi
else
    echo -e "${RED}✗ Go 未安裝，無法安裝 Go 工具${NC}"
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ Ubuntu 環境設置完成！${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "下一步:"
echo "  1. 重新載入環境變數: source ~/.bashrc"
echo "  2. 驗證安裝: ./scripts/check-tools.sh"
echo "  3. 設置專案: make setup"
echo "  4. 啟動服務: make docker-up"
echo ""
echo "💡 提示:"
echo "  - 如果安裝了 Docker，請登出重新登入以使群組權限生效"
echo "  - 使用 Docker 管理 PostgreSQL 更方便 (無需本機安裝)"
echo "  - 詳細文件: DEPENDENCIES.md"
echo ""
