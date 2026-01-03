#!/bin/bash
# WEBPRISM Docker-Only Setup Script
# 適用於已安裝 Docker 和 Go，只需用 Docker 運行 PostgreSQL 的使用者

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 WEBPRISM Docker 環境設置${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "此腳本適用於："
echo "  ✓ 已安裝 Go 和 protoc"
echo "  ✓ 已安裝 Docker"
echo "  ✓ 只需使用 Docker 運行 PostgreSQL"
echo ""

# Check prerequisites
MISSING=()

if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go 未安裝${NC}"
    MISSING+=("Go")
fi

if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker 未安裝${NC}"
    MISSING+=("Docker")
fi

if ! command -v protoc &> /dev/null; then
    echo -e "${RED}✗ protoc 未安裝${NC}"
    MISSING+=("protoc")
fi

if [ ${#MISSING[@]} -gt 0 ]; then
    echo ""
    echo -e "${RED}✗ 缺少必要工具: ${MISSING[*]}${NC}"
    echo ""
    echo "請先執行對應平台的安裝腳本："
    echo "  macOS:   ./scripts/setup-macos.sh"
    echo "  Ubuntu:  ./scripts/setup-ubuntu.sh"
    echo "  Windows: .\\scripts\\setup-windows.ps1"
    exit 1
fi

echo -e "${GREEN}✓ 前置需求檢查通過${NC}"
echo ""

# Check Docker daemon
if ! docker ps &> /dev/null 2>&1; then
    echo -e "${RED}✗ Docker daemon 未運行${NC}"
    echo ""
    echo "請啟動 Docker："
    echo "  macOS/Windows: 開啟 Docker Desktop"
    echo "  Linux: sudo systemctl start docker"
    exit 1
fi

echo -e "${GREEN}✓ Docker daemon 運行中${NC}"
echo ""

# Install Go tools
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}🔧 安裝 Go 開發工具${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

GOPATH=$(go env GOPATH)
echo "GOPATH: $GOPATH"
echo ""

echo "安裝 protoc 插件和 mockery..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install github.com/vektra/mockery/v2@latest

echo ""
echo -e "${GREEN}✓ Go 工具安裝完成${NC}"

# Check GOPATH in PATH
if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
    echo ""
    echo -e "${YELLOW}⚠️  警告: GOPATH/bin 不在您的 PATH 中${NC}"
    echo ""
    echo "請將以下內容加入您的 shell 配置檔:"
    echo ""
    echo -e "${BLUE}export PATH=\"\$PATH:$GOPATH/bin\"${NC}"
fi

# Setup third-party dependencies
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📦 設置第三方依賴${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$SCRIPT_DIR/setup-third-party.sh"

# Setup environment file
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}⚙️  配置環境變數${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env"
ENV_EXAMPLE="$PROJECT_ROOT/.env.example"

if [ ! -f "$ENV_FILE" ]; then
    if [ -f "$ENV_EXAMPLE" ]; then
        echo "建立 .env 檔案..."
        cp "$ENV_EXAMPLE" "$ENV_FILE"
        echo -e "${GREEN}✓ .env 檔案已建立${NC}"
        echo ""
        echo -e "${YELLOW}⚠️  請編輯 .env 並填入機密資訊:${NC}"
        echo "   - WEBPRISM_DATABASE_PASSWORD"
        echo "   - WEBPRISM_SECURITY_ENCRYPTION_KEY (32 bytes)"
        echo ""
        echo "   生成加密金鑰: openssl rand -base64 32 | head -c 32"
    else
        echo -e "${YELLOW}⚠️  找不到 .env.example${NC}"
    fi
else
    echo -e "${GREEN}✓ .env 檔案已存在${NC}"
fi

# Test Docker Compose
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}🐳 驗證 Docker Compose${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

cd "$PROJECT_ROOT"

if [ -f "docker-compose.yml" ]; then
    echo "檢查 docker-compose.yml..."
    docker compose config > /dev/null 2>&1
    echo -e "${GREEN}✓ Docker Compose 配置有效${NC}"
else
    echo -e "${YELLOW}⚠️  找不到 docker-compose.yml${NC}"
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✓ Docker 環境設置完成！${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "下一步:"
echo "  1. 編輯 .env 填入機密資訊"
echo "  2. 生成程式碼: make proto"
echo "  3. 下載依賴: go mod download"
echo "  4. 啟動 Docker 服務: make docker-up"
echo "  5. 執行測試: make test"
echo ""
echo "或直接執行: make setup && make docker-up"
echo ""
echo "💡 提示:"
echo "  - PostgreSQL 將在 Docker 中運行 (localhost:5432)"
echo "  - 不需要在本機安裝 PostgreSQL"
echo "  - 使用 'make docker-down' 停止服務"
echo "  - 使用 'make docker-logs' 查看日誌"
echo ""
