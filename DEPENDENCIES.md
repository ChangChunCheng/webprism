# WEBPRISM 依賴安裝指南

本文件說明 WEBPRISM 所需的**系統級依賴**。Go 模組依賴會自動通過 `go mod` 管理。

**支援平台**: Windows 10/11, macOS 12+, Ubuntu 22.04+

---

## 🚀 快速開始（選擇您的平台）

根據您的作業系統，執行對應的安裝腳本：

```bash
# macOS 使用者
./scripts/setup-macos.sh

# Ubuntu/Debian 使用者
./scripts/setup-ubuntu.sh

# Windows 使用者 (PowerShell)
.\scripts\setup-windows.ps1

# 已有 Go + Docker，只需設置環境
./scripts/setup-docker.sh
```

腳本會自動安裝所需工具並配置環境。完成後請參考 [README.md](README.md) 繼續設置專案。

---

## 📋 系統需求

### 必要工具

以下工具需要**使用者自行安裝**到系統中:

1. **Go 1.21+** - 程式語言環境
2. **Protocol Buffers (protoc)** - Protobuf 編譯器
3. **Docker & Docker Compose** - 容器化部署（推薦）

### 可選工具

1. **golangci-lint** - 程式碼檢查工具
2. **golang-migrate** - 資料庫遷移工具

### 資料庫需求

**PostgreSQL 16+** - 本專案需要 PostgreSQL 作為資料庫，您可以選擇：

- **推薦**: 使用 Docker 容器運行（`make db-up`）
- 使用自訂的 PostgreSQL 實例（需自行安裝和配置）

**注意**: 若選擇自訂 PostgreSQL，此專案不提供安裝或管理指令，請自行確保資料庫已正確設定並在 `.env` 中配置連線資訊。

### Go 工具依賴

以下工具會被安裝到 `GOPATH/bin` 中:

- `protoc-gen-go` - Go Protobuf 插件
- `protoc-gen-go-grpc` - Go gRPC 插件
- `protoc-gen-grpc-gateway` - gRPC Gateway 插件
- `protoc-gen-openapiv2` - OpenAPI 文件生成插件
- `mockery` - Mock 生成工具

這些工具透過 `make install-tools` 自動安裝到您的 GOPATH/bin 目錄。

**注意**: 這些工具會安裝在您的 Go 環境中,對所有 Go 專案可見。如需移除:
```bash
rm $(go env GOPATH)/bin/protoc-gen-*
rm $(go env GOPATH)/bin/mockery
```

### 第三方 Proto 依賴

WEBPRISM 需要 Google API 的 proto 定義檔案（用於 gRPC-Gateway）:

- **googleapis** - Google API 標準定義
  - 來源: https://github.com/googleapis/googleapis
  - 大小: 約 13MB（使用 sparse checkout）
  - 授權: Apache License 2.0

這些依賴**不會被提交到 Git**，會在專案設定時自動下載：

```bash
# 方式 1: 使用 Makefile（推薦）
make setup-third-party

# 方式 2: 直接執行腳本
./scripts/setup-third-party.sh
```

詳細說明請參考: [third_party/README.md](third_party/README.md)

---

## 🚀 平台安裝指南

### 方式 1: 使用自動化腳本（推薦）

選擇您的平台並執行對應的安裝腳本：

#### macOS

```bash
# 一鍵安裝所有工具
./scripts/setup-macos.sh

# 腳本會自動：
# - 檢查並安裝 Homebrew
# - 安裝 Go, protoc, Docker
# - 安裝 Go 開發工具到 GOPATH
# - 配置環境變數

# 注意：腳本不會安裝 PostgreSQL，推薦使用 Docker（make db-up）
```

#### Ubuntu 22.04+

```bash
# 一鍵安裝所有工具
./scripts/setup-ubuntu.sh

# 腳本會自動：
# - 安裝 Go 1.21+
# - 安裝 protoc, git, make
# - 安裝 Docker（推薦）
# - 安裝 Go 開發工具到 GOPATH
# - 配置環境變數

# 注意：腳本不會安裝 PostgreSQL，推薦使用 Docker（make db-up）
```

#### Windows 10/11

```powershell
# PowerShell 執行
.\scripts\setup-windows.ps1

# 腳本會：
# - 檢查已安裝工具
# - 提供安裝指引
# - 建議使用 WSL2 或 Docker Desktop
```

**Windows 最佳實踐**:

1. **推薦**: 使用 WSL2 + Ubuntu，然後執行 `./scripts/setup-ubuntu.sh`
2. 或使用 Docker Desktop + 手動安裝 Go 和 protoc

### 方式 2: 手動安裝

如果您偏好手動安裝，以下是各平台的詳細步驟：

#### macOS 手動安裝

```bash
# 1. 安裝 Homebrew（如果還沒有）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. 安裝必要工具
brew install go protobuf

# 3. 安裝 Docker Desktop（從官網下載）
# 下載：https://www.docker.com/products/docker-desktop/

# 4. 安裝 Go 工具
make install-tools

# 5. 驗證安裝
go version        # 應顯示 go1.21 或更高
protoc --version  # 應顯示 libprotoc 3.x.x 或更高
docker --version  # 驗證 Docker 安裝

# 注意：若要使用自訂 PostgreSQL，請自行安裝和配置
```

#### Ubuntu/Debian 手動安裝

```bash
# 1. 安裝基本工具
sudo apt-get update
sudo apt-get install -y curl wget git build-essential

# 2. 安裝 Go 1.21+
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# 3. 安裝 protoc
sudo apt-get install -y protobuf-compiler

# 4. 安裝 Docker（推薦）
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# 登出重新登入使群組生效

# 5. 安裝 Go 工具
make install-tools

# 6. 驗證安裝
go version
protoc --version
docker --version

# 注意：若要使用自訂 PostgreSQL，請自行安裝和配置
```

#### Windows 手動安裝

##### 選項 A: 使用 WSL2（推薦）

```powershell
# 1. 安裝 WSL2
wsl --install

# 2. 重新啟動電腦

# 3. 安裝 Ubuntu
wsl --install -d Ubuntu

# 4. 在 WSL2 中按照 Ubuntu 安裝步驟進行
```

##### 選項 B: 原生 Windows 安裝

1. **Go**: 從 <https://go.dev/dl/> 下載並安裝
2. **protoc**: 從 <https://github.com/protocolbuffers/protobuf/releases> 下載
   - 下載 `protoc-xx.x-win64.zip`
   - 解壓並將 `bin` 目錄加入系統 PATH
3. **Docker Desktop**: 從 <https://www.docker.com/products/docker-desktop/> 下載並安裝
4. **安裝 Go 工具**: 在專案目錄執行 `make install-tools`

### 方式 3: 只使用 Docker（最簡單）

如果您已安裝 Go、protoc 和 Docker，只需用 Docker 運行 PostgreSQL：

```bash
# 執行 Docker 環境設置
./scripts/setup-docker.sh

# 會自動：
# - 安裝 Go 工具到 GOPATH
# - 設置第三方依賴
# - 配置 .env 環境檔案
# - 驗證 Docker 環境
```

---

## 🔧 GOPATH 配置

Go 工具會被安裝到 `GOPATH/bin`，需要將其加入 PATH:

```bash
# 查看您的 GOPATH
go env GOPATH

# 將 GOPATH/bin 加入 PATH（添加到 ~/.zshrc 或 ~/.bashrc）
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# 驗證
echo $PATH | grep "$(go env GOPATH)/bin"
```

---

## 📦 Go 模組依賴

專案的 Go 依賴通過 `go.mod` 管理，會自動下載到本地模組快取:

```bash
# 下載所有依賴
go mod download

# 整理依賴（移除未使用的）
go mod tidy

# 驗證依賴
go mod verify
```

所有 Go 依賴會被下載到:
- **位置**: `$(go env GOMODCACHE)`（通常是 `~/go/pkg/mod/`）
- **不影響**: 全域 Go 環境
- **隔離性**: 每個專案獨立管理依賴版本

---

## ✅ 驗證安裝

完成安裝後，使用檢查腳本驗證所有工具：

```bash
# 所有平台（macOS/Ubuntu/Windows WSL2）
./scripts/check-tools.sh
```

腳本會自動檢測您的作業系統並顯示相應的安裝指引。

**預期輸出**:

```text
🔍 檢查 WEBPRISM 開發工具...

Platform: macos (x86_64)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
必要工具
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ Protocol Buffers 編譯器
  libprotoc 3.21.12

✓ Go 語言
  go version go1.21.6 darwin/amd64

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Go 插件
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ protoc-gen-go
✓ protoc-gen-go-grpc
✓ protoc-gen-grpc-gateway
✓ protoc-gen-openapiv2

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
資料庫（可選）
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ PostgreSQL
  psql (PostgreSQL) 16.1
```

---

## 🚀 使用方式說明

WEBPRISM 支援三種開發模式，根據您的環境選擇：

### 模式 1: 混合模式（推薦）

**說明**: PostgreSQL 用 Docker，WEBPRISM 本機運行

**適合**:

- 大部分開發者
- 不想在本機安裝 PostgreSQL
- 需要頻繁重啟服務除錯

**使用方式**:

```bash
# 啟動 PostgreSQL（Docker）
make db-up

# 運行 WEBPRISM（本機）
make run-server

# 停止 PostgreSQL
make db-down
```

### 模式 2: 完全本機

**說明**: 使用自訂的 PostgreSQL 實例

**適合**:

- 已自行配置 PostgreSQL
- 不想使用 Docker
- 資源受限的環境

**使用方式**:

```bash
# 確保您的 PostgreSQL 已運行並正確配置
# 在 .env 中設定連線資訊

# 運行 WEBPRISM
make run-local
```

**注意**: 此專案不提供 PostgreSQL 的安裝或管理指令，請自行確保資料庫已正確設定。

### 模式 3: 完全 Docker

**說明**: 所有服務都在 Docker 容器中

**適合**:

- 生產環境部署
- CI/CD 流程
- 需要完全隔離的測試環境

**使用方式**:

```bash
# 啟動所有服務（PostgreSQL + WEBPRISM）
make docker-up-full

# 查看日誌
make docker-logs

# 停止服務
make docker-down

# 清理資源（包含 volumes）
make docker-clean
```

### 快速對照表

| 模式 | PostgreSQL | WEBPRISM | 命令 | 適合情境 |
| ------ | ------ | ------ | ------ | ------ |
| 混合模式 | Docker | 本機 | `make db-up && make run-server` | 🔥 日常開發（推薦） |
| 完全本機 | 自訂 | 本機 | `make run-local` | 已自行配置 PostgreSQL |
| 完全 Docker | Docker | Docker | `make docker-up-full` | 生產/CI/CD |

---

## 📝 安裝腳本摘要

| 腳本 | 用途 | 適用平台 |
| ------ | ------ | ---------- |
| [setup-macos.sh](scripts/setup-macos.sh) | 自動安裝所有工具 | macOS 12+ |
| [setup-ubuntu.sh](scripts/setup-ubuntu.sh) | 自動安裝所有工具 | Ubuntu 22.04+ |
| [setup-windows.ps1](scripts/setup-windows.ps1) | 檢查工具並提供指引 | Windows 10/11 |
| [setup-docker.sh](scripts/setup-docker.sh) | 只設置 Docker 環境 | 所有平台 |
| [check-tools.sh](scripts/check-tools.sh) | 驗證工具安裝 | 所有平台 |
| [detect-platform.sh](scripts/detect-platform.sh) | 檢測作業系統 | 所有平台 |

---

## 📚 延伸閱讀

- **快速開始**: [README.md](README.md) - 三步啟動服務
- **詳細安裝說明**: [docs/user-guide/01-INSTALLATION.md](docs/user-guide/01-INSTALLATION.md)
- **快速開始指南**: [docs/user-guide/02-QUICKSTART.md](docs/user-guide/02-QUICKSTART.md)
- **macOS 特定設定**: [docs/user-guide/04-MACOS_SETUP.md](docs/user-guide/04-MACOS_SETUP.md)

---

## 🆘 常見問題

### Q: 我應該使用哪種安裝方式？

A: 依據您的情況選擇：

- **新手 / 快速開始**: 使用對應平台的自動化腳本（`setup-*.sh`）
- **已有 Go + Docker**: 使用 `./scripts/setup-docker.sh`
- **進階使用者 / 特殊需求**: 手動安裝並自行配置

### Q: Windows 使用者應該用 WSL2 還是原生 Windows？

A: 強烈推薦 WSL2

- ✅ **WSL2 優點**: 完整的 Linux 環境、腳本支援佳、開發體驗一致
- ⚠️ **原生 Windows**: 需手動下載安裝、路徑配置複雜、部分工具支援有限
- 💡 **最佳實踐**: WSL2 + Ubuntu + Docker Desktop

### Q: 必須安裝 PostgreSQL 嗎？

A: 專案需要 PostgreSQL，但不需要本機安裝

- **推薦**: 使用 `make db-up` 在 Docker 中運行 PostgreSQL
- **可選**: 使用自訂的 PostgreSQL 實例（需自行安裝和配置）
- **好處**: Docker 方式更簡單、隔離性好、易於清理
- **注意**: 此專案不提供 PostgreSQL 的安裝或管理指令

### Q: 為什麼需要安裝 protoc 到系統？

A: `protoc` 是 Protocol Buffers 的核心編譯器，是 C++ 實現的二進制工具，無法通過 Go 模組安裝。但 protoc 的 Go 插件（如 `protoc-gen-go`）會被安裝到 GOPATH/bin。

### Q: Go 工具會影響我的系統嗎？

A: Go 工具會安裝在 `GOPATH/bin` 中（通常是 `~/go/bin`），對您的所有 Go 專案可見，但不會影響系統全域環境（如 `/usr/local/bin`）。這是 Go 生態系統的標準做法。

如需移除特定工具:

```bash
rm $(go env GOPATH)/bin/protoc-gen-*
```

### Q: GOPATH 和 GOMODCACHE 的區別？

A:

- `GOPATH/bin`: 存放可執行工具（如 protoc-gen-go）
- `GOMODCACHE`: 存放下載的模組依賴（如 grpc-gateway 庫代碼）

兩者都是本地目錄，不影響全域環境。

### Q: 安裝腳本失敗了怎麼辦？

A: 依據錯誤類型處理：

1. **權限錯誤**: 某些步驟需要 `sudo`，請確保有管理員權限
2. **網路問題**: 檢查網路連線，某些下載可能需要 VPN
3. **版本不相容**: 確保系統版本符合要求（Ubuntu 22.04+, macOS 12+）
4. **手動安裝**: 如腳本持續失敗，可參考「手動安裝」章節逐步操作

### Q: 如何在多個 Go 版本間切換？

A: 使用 Go 版本管理工具

- **gvm** (Go Version Manager): <https://github.com/moovweb/gvm>
- **asdf**: <https://asdf-vm.com/> (支援多種語言)
- **手動**: 下載不同版本並修改 PATH

**注意**: WEBPRISM 需要 Go 1.21+

### Q: Docker Desktop 佔用太多資源怎麼辦？

A: 優化 Docker 資源配置

1. 開啟 Docker Desktop 設定
2. 調整 Resources 選項：
   - CPUs: 2-4 核心
   - Memory: 4-8 GB
   - Swap: 1-2 GB
3. 只運行必要的容器：`docker compose up -d db`（只啟動資料庫）

---

**安裝完成後**，請繼續閱讀 [README.md](README.md) 開始使用 WEBPRISM！
