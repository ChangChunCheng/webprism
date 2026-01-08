# WEBPRISM Server Startup Guide

## 🎯 改進說明

WEBPRISM HTTP Server 已經升級，現在包含：
- ✅ 友好的啟動訊息顯示所有連線資訊
- ✅ **Swagger UI** - 互動式 API 文件介面
- ✅ **OpenAPI 規格** - 完整的 API 定義文件

## 📚 Swagger UI - 互動式 API 文件

啟動服務器後，您可以透過瀏覽器訪問完整的 API 文件：

### 快速訪問

- **Swagger UI (互動式介面)**: http://localhost:8080/swagger-ui/
- **OpenAPI Specification (JSON)**: http://localhost:8080/swagger.json

### Swagger UI 功能

✅ **互動式 API 瀏覽器** - 直接在瀏覽器中測試所有端點
✅ **完整的 API 文件** - 所有服務、端點和資料模型的詳細文件
✅ **請求/回應範例** - 查看請求格式和預期回應
✅ **立即試用** - 從 UI 執行真實的 API 呼叫
✅ **Schema 驗證** - 基於 OpenAPI 規格的自動完成和驗證

### 使用 Swagger UI

1. **開啟 Swagger UI** - 在瀏覽器訪問: http://localhost:8080/swagger-ui/

2. **探索 API**:
   - 瀏覽按服務分組的所有端點（SpecService、AuthService、ProxyService、HealthService）
   - 點擊任何端點查看詳細文件
   - 查看請求參數、Body Schema 和回應格式

3. **測試端點**:
   - 點擊任何端點上的 "Try it out"
   - 填入必要的參數
   - 點擊 "Execute" 執行真實的 API 呼叫
   - 查看回應狀態、Headers 和 Body

4. **下載 OpenAPI 規格**:
   - 使用 Swagger UI 中的 "Download" 按鈕，或
   - 直接訪問 http://localhost:8080/swagger.json

## 📋 啟動訊息預覽

當你啟動 WEBPRISM Server 時，會看到以下美觀的啟動訊息：

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                        🚀 WEBPRISM Server Started                             ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║  HTTP (REST API):                                                             ║
║    → http://localhost:8080                                                     ║
║                                                                               ║
║  📚 API Documentation:                                                        ║
║    → Swagger UI:  http://localhost:8080/swagger-ui/                           ║
║    → OpenAPI JSON: http://localhost:8080/swagger.json                          ║
║                                                                               ║
║  gRPC:                                                                        ║
║    → localhost:9090                                                            ║
║                                                                               ║
║  Available Endpoints:                                                         ║
║    • POST   /api/v1/specs                       - Upload API specification     ║
║    • GET    /api/v1/specs                       - List all specifications      ║
║    • GET    /api/v1/specs/{id}                  - Get specific specification   ║
║    • DELETE /api/v1/specs/{id}                  - Delete specification         ║
║    • POST   /api/v1/specs/{spec_id}/auth        - Set auth configuration       ║
║    • GET    /api/v1/specs/{spec_id}/auth        - Get auth configuration       ║
║    • POST   /api/v1/proxy                       - Execute proxy request        ║
║    • POST   /api/v1/specs/{spec_id}/health-check - Check API health             ║
║    • GET    /api/v1/specs/{spec_id}/health-check - Get health check history     ║
║    • GET    /api/v1/version                     - Get system version           ║
║                                                                               ║
║  Quick Start:                                                                 ║
║    curl http://localhost:8080/api/v1/version                                   ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

## 🚀 快速啟動

### 1. 使用 Docker Compose（推薦）

最簡單的方式是使用 Docker Compose 啟動完整環境：

```bash
# 啟動所有服務（PostgreSQL + WEBPRISM Server）
make docker-up-full

# 查看日誌
docker-compose -f deployments/docker-compose.full.yml logs -f webprism
```

### 2. 使用本地編譯

如果你想在本地運行（需要先啟動 PostgreSQL）：

```bash
# 1. 啟動 PostgreSQL
make docker-up

# 2. 設置環境變數（如果沒有 .env 文件）
cp .env.example .env

# 3. 啟動 Server
make run-server
```

### 3. 直接執行二進制文件

```bash
# 設置必要的環境變數
export WEBPRISM_DATABASE_PASSWORD=webprism
export WEBPRISM_SECURITY_ENCRYPTION_KEY=dev-key-12345678901234567890abcd

# 運行 server
./bin/webprism-server
```

## 🔗 快速測試連線

啟動成功後，你可以立即測試以下端點：

### 1. 健康檢查

```bash
curl http://localhost:8080/v1/health
```

預期回應：
```json
{
  "status": "healthy",
  "timestamp": "2026-01-07T15:30:00Z",
  "database": {
    "status": "healthy",
    "specs_count": 0
  }
}
```

### 2. 列出 API 規格

```bash
curl http://localhost:8080/v1/specs
```

### 3. 列出認證配置

```bash
curl http://localhost:8080/v1/auth/configs
```

## 📊 可用的端點

### API 規格管理

| 方法 | 端點 | 說明 |
|------|------|------|
| POST | `/v1/specs` | 上傳新的 OpenAPI 規格 |
| GET | `/v1/specs` | 列出所有規格 |
| GET | `/v1/specs/{id}` | 取得特定規格 |
| DELETE | `/v1/specs/{id}` | 刪除規格 |

#### 範例：上傳 OpenAPI 規格

```bash
curl -X POST http://localhost:8080/v1/specs \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My API",
    "version": "1.0.0",
    "openapi_spec": "{...OpenAPI JSON...}"
  }'
```

### 認證配置管理

| 方法 | 端點 | 說明 |
|------|------|------|
| POST | `/v1/auth/configs` | 建立認證配置 |
| GET | `/v1/auth/configs` | 列出所有配置 |
| GET | `/v1/auth/configs/{id}` | 取得特定配置 |
| PUT | `/v1/auth/configs/{id}` | 更新配置 |
| DELETE | `/v1/auth/configs/{id}` | 刪除配置 |

#### 範例：建立 API Key 認證

```bash
curl -X POST http://localhost:8080/v1/auth/configs \
  -H "Content-Type: application/json" \
  -d '{
    "spec_id": "your-spec-id",
    "auth_type": "api_key",
    "config": {
      "api_key": "your-api-key-here"
    }
  }'
```

#### 範例：建立 Basic Auth

```bash
curl -X POST http://localhost:8080/v1/auth/configs \
  -H "Content-Type: application/json" \
  -d '{
    "spec_id": "your-spec-id",
    "auth_type": "basic",
    "config": {
      "username": "your-username",
      "password": "your-password"
    }
  }'
```

### 代理請求執行

| 方法 | 端點 | 說明 |
|------|------|------|
| POST | `/v1/proxy/execute` | 執行 API 代理請求 |

#### 範例：執行代理請求

```bash
curl -X POST http://localhost:8080/v1/proxy/execute \
  -H "Content-Type: application/json" \
  -d '{
    "spec_id": "your-spec-id",
    "endpoint": "/users",
    "method": "GET",
    "headers": {
      "Accept": "application/json"
    }
  }'
```

### 健康檢查

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/v1/health` | 檢查服務健康狀態 |

## 🎨 自訂端口

如果你想更改預設端口，可以透過以下方式：

### 方法 1: 修改配置文件

編輯 `deployments/config.yaml`:

```yaml
server:
  http:
    port: 8081  # 改成你想要的端口
  grpc:
    port: 9091  # 改成你想要的端口
```

### 方法 2: 使用環境變數

```bash
export WEBPRISM_SERVER_HTTP_PORT=8081
export WEBPRISM_SERVER_GRPC_PORT=9091
./bin/webprism-server
```

## 📝 日誌輸出

啟動訊息會同時輸出到：

1. **Console (stdout)** - 美觀的 Banner 格式，方便人類閱讀
2. **Structured Logger** - JSON 格式的結構化日誌，方便機器解析

結構化日誌範例：
```json
{
  "level": "info",
  "msg": "HTTP server starting",
  "http_url": "http://localhost:8080",
  "grpc_addr": "localhost:9090",
  "health_check": "http://localhost:8080/v1/health",
  "time": "2026-01-07T15:30:00Z"
}
```

## 🔍 除錯技巧

### 查看詳細日誌

如果需要更詳細的日誌，可以設置日誌級別為 `debug`：

```bash
export WEBPRISM_LOGGING_LEVEL=debug
./bin/webprism-server
```

### 檢查端口是否被占用

```bash
# macOS/Linux
lsof -i :8080

# 或使用 netstat
netstat -an | grep 8080
```

### 測試 gRPC 連線

使用 `grpcurl` 測試 gRPC 端點：

```bash
# 列出可用的服務
grpcurl -plaintext localhost:9090 list

# 調用健康檢查
grpcurl -plaintext localhost:9090 webprism.v1.HealthService/Check
```

## 🎯 下一步

1. ✅ 閱讀 [快速開始指南](GETTING_STARTED.md)
2. ✅ 查看 [API 使用範例](user-guide/02-QUICKSTART.md)
3. ✅ 瀏覽 [完整 API 文件](../api/proto/v1/)
4. ✅ 加入開發：[貢獻指南](development/CONTRIBUTING.md)

## 💡 提示

- **始終檢查健康端點**：確保服務運行正常
- **使用結構化日誌**：便於追蹤問題
- **啟用 CORS**：已預設啟用，支援跨域請求
- **自動重載**：開發時可使用 `air` 等工具實現自動重載

## 🐛 常見問題

### Q: 端口已被占用怎麼辦？

A: 修改配置文件或使用環境變數更改端口。

### Q: 無法連接數據庫？

A: 確保 PostgreSQL 正在運行，並檢查環境變數 `WEBPRISM_DATABASE_PASSWORD` 是否設置正確。

### Q: 看不到啟動 Banner？

A: 確保你沒有將日誌輸出重定向到文件。Banner 會輸出到 stdout。

---

**相關文件**:
- [README.md](../README.md)
- [Getting Started](GETTING_STARTED.md)
- [API Documentation](../api/proto/v1/)
- [Configuration Guide](../deployments/config.yaml)
