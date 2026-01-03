# WEBPRISM 工具安裝指南

本指南說明如何安裝 WEBPRISM 開發所需的所有工具。

> **注意**: 本文件詳細說明手動安裝步驟。快速安裝請參考 [DEPENDENCIES.md](../../DEPENDENCIES.md)

---

## 📦 工具分類

### 系統級工具（需手動安裝）

以下工具需要安裝到系統環境:

1. **Go 1.21+** - 程式語言環境
2. **protoc** - Protocol Buffers 編譯器（C++ 實現的二進制工具）
3. **Docker & Docker Compose** - 容器化部署（推薦）
4. **PostgreSQL 16+** - 資料庫（可選，可用 Docker 替代）

### Go 工具（自動安裝到 GOPATH/bin）

以下工具通過 `make install-tools` 自動安裝，**不影響系統環境**:

1. **protoc-gen-go** - Go 語言的 Protobuf 插件
2. **protoc-gen-go-grpc** - gRPC 的 Go 插件
3. **protoc-gen-grpc-gateway** - grpc-gateway 插件
4. **protoc-gen-openapiv2** - OpenAPI 文件生成插件
5. **mockery** - Mock 生成工具

---

## 🚀 macOS 快速安裝

### 方法 1: 使用 Homebrew（推薦）

```bash
# 1. 安裝 Homebrew（如果還沒有）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. 安裝 protoc
brew install protobuf

# 3. 驗證安裝
protoc --version
# 應該顯示: libprotoc 3.x.x 或更高版本

# 4. 安裝其他工具
brew install golangci-lint
brew install golang-migrate
brew install postgresql@16  # 可選

# 5. 啟動 PostgreSQL（如果安裝了）
brew services start postgresql@16
```

### 方法 2: 手動安裝 protoc

如果 Homebrew 有問題，可以手動下載：

```bash
# 下載最新版本
PROTOC_VERSION="25.1"
cd /tmp
curl -LO "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-universal_binary.zip"

# 解壓縮
unzip "protoc-${PROTOC_VERSION}-osx-universal_binary.zip" -d protoc

# 移動到系統路徑
sudo mv protoc/bin/protoc /usr/local/bin/
sudo mv protoc/include/* /usr/local/include/

# 清理
rm -rf protoc "protoc-${PROTOC_VERSION}-osx-universal_binary.zip"

# 驗證
protoc --version
```

---

## 🔧 安裝 Go 工具

安裝完系統工具後，使用 Makefile 安裝 Go 工具到 GOPATH/bin:

```bash
# 進入專案目錄
cd webprism

# 安裝所有 Go 工具（推薦）
make install-tools

# 這會將工具安裝到: $(go env GOPATH)/bin
# 通常是: ~/go/bin/
```

安裝的工具包括:
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-grpc-gateway
- protoc-gen-openapiv2
- mockery

**重要**: 確保 `$(go env GOPATH)/bin` 在您的 PATH 中:

```bash
# 添加到 ~/.zshrc 或 ~/.bashrc
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

---

## ✅ 驗證安裝

執行以下命令驗證所有工具都已正確安裝：

```bash
# 檢查 protoc
protoc --version
# 預期: libprotoc 3.x.x 或更高

# 檢查 Go 插件
which protoc-gen-go
which protoc-gen-go-grpc
which protoc-gen-grpc-gateway
which protoc-gen-openapiv2

# 檢查其他工具
golangci-lint --version
migrate -version

# 檢查 PostgreSQL（如果安裝了）
psql --version
```

所有命令都應該成功執行並顯示版本號。

---

## 🎯 完整設定流程

安裝完工具後，按以下順序執行：

```bash
# 1. 進入專案目錄
cd /Users/ChangChunCheng/projects/webprism

# 2. 安裝 Go 插件
make install-tools

# 3. 生成 Protobuf 程式碼
make proto

# 4. 驗證生成結果
ls -la gen/go/v1/
# 應該看到 .pb.go 和 .pb.gw.go 檔案

# 6. 下載 Go 依賴
go mod tidy

# 7. 編譯專案
make build

# 8. 驗證編譯結果
ls -la bin/
# 應該看到: webprism-server, webprism-mcp, webprism
```

---

## 🐛 常見問題

### 問題 1: `protoc: command not found`

**原因**: protoc 未安裝或不在 PATH 中

**解決方案**:
```bash
# 安裝 protoc
brew install protobuf

# 如果已安裝但找不到，檢查 PATH
echo $PATH

# 如果使用 zsh，確保 Homebrew 在 PATH 中
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zshrc
source ~/.zshrc
```

### 問題 2: `protoc-gen-go: program not found or is not executable`

**原因**: Go 插件未安裝或 GOPATH/bin 不在 PATH 中

**解決方案**:
```bash
# 確保 GOPATH/bin 在 PATH 中
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# 重新安裝插件
make install-tools
```

### 問題 3: `permission denied` 執行插件

**原因**: 插件沒有執行權限

**解決方案**:
```bash
# 為所有 Go 插件添加執行權限
chmod +x $(go env GOPATH)/bin/protoc-gen-*
```

### 問題 4: Homebrew 安裝緩慢

**解決方案**: 使用國內鏡像（如果在中國）
```bash
# 替換 Homebrew 源
export HOMEBREW_BREW_GIT_REMOTE="https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git"
export HOMEBREW_CORE_GIT_REMOTE="https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git"
```

### 問題 5: M1/M2 Mac 的 PATH 問題

**原因**: Apple Silicon Mac 的 Homebrew 安裝在 `/opt/homebrew`

**解決方案**:
```bash
# 在 ~/.zshrc 中添加
eval "$(/opt/homebrew/bin/brew shellenv)"

# 重新載入
source ~/.zshrc
```

---

## 📋 工具版本參考

以下是測試過的工具版本組合：

```
protoc: 25.1 或更高
Go: 1.21 或更高
protoc-gen-go: v1.31.0 或更高
protoc-gen-go-grpc: v1.3.0 或更高
protoc-gen-grpc-gateway: v2.19.0 或更高
golangci-lint: v1.55.0 或更高
golang-migrate: v4.17.0 或更高
PostgreSQL: 14, 15, 或 16
```

---

## 🔍 檢查腳本

創建一個檢查腳本來驗證所有工具：

```bash
# 創建檢查腳本
cat > scripts/check-tools.sh << 'EOF'
#!/bin/bash

echo "🔍 檢查 WEBPRISM 開發工具..."
echo

# 顏色定義
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

check_command() {
    if command -v $1 &> /dev/null; then
        echo -e "${GREEN}✓${NC} $1: $(command -v $1)"
        $1 --version 2>&1 | head -1
    else
        echo -e "${RED}✗${NC} $1: 未安裝"
        return 1
    fi
    echo
}

# 檢查所有工具
check_command protoc
check_command protoc-gen-go
check_command protoc-gen-go-grpc
check_command protoc-gen-grpc-gateway
check_command protoc-gen-openapiv2
check_command golangci-lint
check_command migrate
check_command psql

echo "檢查完成！"
EOF

chmod +x scripts/check-tools.sh

# 執行檢查
./scripts/check-tools.sh
```

---

## 🎓 更多資源

- **Protocol Buffers**: https://protobuf.dev/
- **gRPC**: https://grpc.io/
- **grpc-gateway**: https://grpc-ecosystem.github.io/grpc-gateway/
- **golangci-lint**: https://golangci-lint.run/
- **golang-migrate**: https://github.com/golang-migrate/migrate

---

**下一步**: 安裝完所有工具後，請參考 [QUICKSTART.md](QUICKSTART.md) 繼續設定。
