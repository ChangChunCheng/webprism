# 🚀 WEBPRISM 快速開始指南

本指南將帶領你從零開始設置並運行 WEBPRISM。

---

## 📋 目錄

- [前置需求](#前置需求)
- [安裝步驟](#安裝步驟)
- [啟動模式](#啟動模式)
- [驗證安裝](#驗證安裝)
- [第一次使用](#第一次使用)
- [常見問題](#常見問題)

---

## 前置需求

請先安裝以下系統工具（**一次性設置**）:

| 工具 | 用途 | 安裝方式 |
|------|------|----------|
| **Go 1.21+** | 程式語言環境 | [下載安裝](https://go.dev/dl/) 或 `brew install go` |
| **protoc** | Protobuf 編譯器 | `brew install protobuf` |
| **Docker** | 容器化部署 | [Docker Desktop](https://www.docker.com/products/docker-desktop/) |

詳細安裝指南請參考 [DEPENDENCIES.md](../DEPENDENCIES.md)

---

## 安裝步驟

### 步驟 1: Clone 專案

```bash
git clone https://github.com/ChangChunCheng/webprism.git
cd webprism
```

### 步驟 2: 初始化專案

```bash
make setup
```

這個指令會：
- 安裝必要的開發工具（protoc-gen-go, mockery 等）
- 生成 Protobuf 程式碼
- 生成測試 Mock
- 下載 Go 依賴套件

> **💡 提示**: Clone 專案後，IDE 可能會顯示 import 錯誤，這是正常的。執行 `make setup` 會生成 `gen/` 和 `internal/mocks/` 目錄，錯誤就會消失。

### 步驟 3: 配置環境變數

```bash
cp .env.example .env
```

編輯 `.env` 檔案，填入必要資訊：

```bash
# 資料庫設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=webprism
DB_PASSWORD=your_password_here
DB_NAME=webprism
DB_SSLMODE=disable

# 伺服器設定
SERVER_PORT=8080
GRPC_PORT=50051

# 安全設定
ENCRYPTION_KEY=your-32-byte-encryption-key-here
```

---

## 啟動模式

WEBPRISM 提供三種啟動模式，選擇最適合你的方式：

### 🔥 方式 1: 混合模式（推薦）

**適合**: 大部分開發者 - PostgreSQL 用 Docker，WEBPRISM 本機運行

```bash
# 1. 啟動 PostgreSQL（Docker）
make db-up

# 2. 運行 WEBPRISM（本機）
make run-server
```

**優點**:
- ✅ 快速重啟
- ✅ 方便除錯
- ✅ 資料庫隔離

---

### 💻 方式 2: 完全本機

**適合**: 已自行配置 PostgreSQL 的開發者

```bash
# 1. 確保 PostgreSQL 已運行
# 2. 更新 .env 中的資料庫連線資訊

# 3. 運行 WEBPRISM
make run-local
```

**優點**:
- ✅ 不需要 Docker
- ✅ 資源佔用少

**注意**: 此專案不提供 PostgreSQL 安裝或管理指令，請自行確保 PostgreSQL 已正確設定並運行。

---

### 🐳 方式 3: 完全 Docker

**適合**: 生產環境部署、CI/CD、隔離測試

```bash
# 啟動所有服務（Docker）
make docker-up-full
```

**優點**:
- ✅ 完全隔離
- ✅ 接近生產環境
- ✅ 易於部署

---

## 驗證安裝

### 檢查服務狀態

```bash
# 檢查 HTTP 服務
curl http://localhost:8080/api/v1/specs

# 預期輸出：JSON 格式的 API 規格列表（初始為空陣列）
{"specs":[]}
```

### 查看日誌

```bash
# Docker 模式
make docker-logs

# 本機模式
# 查看終端機輸出
```

---

## 第一次使用

### 1. 上傳 API 規格

準備一個 OpenAPI 規格檔案（例如 `weather-api.json`），然後使用 CLI 上傳：

```bash
./bin/webprism spec upload weather-api.json --name "Weather API"
```

### 2. 設定 API 認證

```bash
# 設定 Bearer Token 認證
./bin/webprism auth set <spec-id> --type bearer --credentials token=your_api_token

# 設定 API Key 認證
./bin/webprism auth set <spec-id> --type apikey --credentials key=X-API-Key,value=your_api_key
```

### 3. 調用 API

```bash
# 透過 CLI 調用
./bin/webprism proxy call <spec-id> getWeather --param city=Taipei

# 透過 HTTP API 調用
curl -X POST http://localhost:8080/api/v1/proxy/call \
  -H "Content-Type: application/json" \
  -d '{
    "spec_id": "<spec-id>",
    "operation_id": "getWeather",
    "parameters": {"city": "Taipei"}
  }'
```

### 4. 查看已上傳的規格

```bash
# 使用 CLI
./bin/webprism spec list

# 使用 HTTP API
curl http://localhost:8080/api/v1/specs
```

---

## 常見問題

### Q1: 執行 `make setup` 時出現權限錯誤

**解決方式**:

```bash
# macOS/Linux
chmod +x scripts/*.sh
make setup
```

### Q2: Docker 無法啟動資料庫

**解決方式**:

```bash
# 檢查 Docker 是否運行
docker ps

# 如果沒有運行，啟動 Docker Desktop
# 然後重試
make db-up
```

### Q3: 連接資料庫失敗

**檢查清單**:

1. 確認 PostgreSQL 已啟動
2. 檢查 `.env` 中的資料庫設定
3. 確認資料庫密碼正確
4. 檢查防火牆設定

```bash
# 測試資料庫連線
psql -h localhost -p 5432 -U webprism -d webprism
```

### Q4: Port 已被佔用

**解決方式**:

```bash
# 查看佔用 Port 的程式
lsof -i :8080
lsof -i :50051
lsof -i :5432

# 停止佔用的程式或修改 .env 中的 Port 設定
```

### Q5: IDE 顯示 import 錯誤

這是正常的！執行 `make setup` 後會生成必要的檔案。

```bash
make setup
```

---

## 使用場景範例

### 場景 1: LLM 工具調用

讓 AI 模型（Claude、GPT）透過 MCP 協定直接調用任何有 OpenAPI 文件的服務。

```bash
# 上傳 API 規格
./bin/webprism spec upload weather-api.json --name "Weather API"

# AI 可透過 MCP 直接調用
# "幫我查詢台北今天的天氣"
```

### 場景 2: API 整合簡化

開發團隊只需上傳第三方 API 文件，即可立即獲得標準化的調用介面。

```bash
# 上傳 Stripe API 規格
./bin/webprism spec upload stripe-api.json

# 設定認證
./bin/webprism auth set <spec-id> --type bearer --credentials token=sk_test_...

# 直接調用
./bin/webprism proxy call <spec-id> createPayment --param amount=1000
```

### 場景 3: 多環境管理

統一管理開發、測試、生產環境的 API 憑證。

```bash
# 開發環境
./bin/webprism auth set <spec-id> --type bearer --credentials token=dev_token

# 生產環境
./bin/webprism auth set <spec-id> --type bearer --credentials token=prod_token
```

---

## 下一步

- 📖 **深入學習**: 閱讀 [使用指南](user-guide/02-QUICKSTART.md)
- 👨‍💻 **參與開發**: 閱讀 [開發指南](development/CONTRIBUTING.md)
- 🧪 **執行測試**: 閱讀 [測試指南](testing/QUICK_START.md)
- 📚 **瀏覽文件**: 查看 [文件中心](README.md)

---

**遇到問題？** 查看 [完整文件](README.md) 或提出 [Issue](https://github.com/ChangChunCheng/webprism/issues)
