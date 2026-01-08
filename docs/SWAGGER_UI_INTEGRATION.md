# Swagger UI Integration Guide

## 概述

WEBPRISM 現在整合了完整的 Swagger UI，讓您可以透過瀏覽器互動式地探索和測試所有 API 端點。

## 功能亮點

### ✅ 已完成的功能

1. **Swagger UI 介面**
   - 使用 Swagger UI 5.10.5 (最新穩定版)
   - 透過 CDN 加載，無需本地檔案
   - 完全響應式設計

2. **OpenAPI 規格合併**
   - 自動合併所有服務的 swagger.json 檔案
   - 單一統一的 API 文件
   - 包含所有端點、參數和資料模型

3. **服務器整合**
   - 無縫整合到 HTTP server
   - 自動路由配置
   - 啟動時顯示 Swagger UI URL

4. **支援的服務**
   - SpecService - API 規格管理
   - AuthService - 認證配置
   - ProxyService - 請求代理
   - HealthService - 健康檢查

## 快速開始

### 1. 啟動服務器

```bash
# 啟動資料庫
make db-up

# 建置並啟動服務器
make build
make run-server
```

### 2. 訪問 Swagger UI

在瀏覽器中開啟：

```text
http://localhost:8080/swagger-ui/
```

### 3. 訪問 OpenAPI 規格

直接下載 JSON 格式的 OpenAPI 規格：

```bash
curl http://localhost:8080/swagger.json > webprism-api.json
```

## 使用範例

### 測試版本端點

1. 在 Swagger UI 中找到 **HealthService**
2. 展開 `GET /api/v1/version`
3. 點擊 "Try it out"
4. 點擊 "Execute"
5. 查看回應：

```json
{
  "version": "799cf17-dirty",
  "gitCommit": "799cf1715f75638031e732ef387c2146e658304f",
  "gitBranch": "Add-CLAUDE.md",
  "buildTime": "2026-01-08T09:17:31Z",
  "goVersion": "go1.24.11",
  "platform": "darwin/arm64"
}
```

### 上傳 API 規格

1. 找到 **SpecService**
2. 展開 `POST /api/v1/specs`
3. 點擊 "Try it out"
4. 填入請求 Body：

```json
{
  "name": "My API",
  "specContent": "{\"openapi\":\"3.0.0\",\"info\":{\"title\":\"My API\",\"version\":\"1.0.0\"}}",
  "format": "SPEC_FORMAT_JSON"
}
```

5. 點擊 "Execute"
6. 查看回應中的 spec ID

### 設定認證

1. 找到 **AuthService**
2. 展開 `POST /api/v1/specs/{specId}/auth`
3. 點擊 "Try it out"
4. 填入 `specId` (使用上一步獲得的 ID)
5. 填入請求 Body：

```json
{
  "authType": "AUTH_TYPE_API_KEY",
  "credentials": {
    "api_key": "your-api-key-here",
    "header_name": "X-API-Key"
  }
}
```

6. 點擊 "Execute"

## 技術細節

### 架構

```text
HTTP Server
├── /api/v1/*           → API 端點 (gRPC-Gateway)
├── /swagger-ui/        → Swagger UI 介面 (HTML)
└── /swagger.json       → 合併的 OpenAPI 規格 (JSON)
```

### OpenAPI 規格合併

系統會自動合併以下檔案：

- `gen/openapi/v1/spec.swagger.json`
- `gen/openapi/v1/auth.swagger.json`
- `gen/openapi/v1/proxy.swagger.json`
- `gen/openapi/v1/health.swagger.json`

合併包括：

- **Paths** - 所有端點定義
- **Definitions** - 所有資料模型
- **Tags** - 服務分組
- **Metadata** - API 資訊和版本

### 實作檔案

1. **swagger.go** - Swagger 整合邏輯
   - `mergeSwaggerSpecs()` - 合併所有規格檔案
   - `serveSwaggerUI()` - 提供 Swagger UI HTML
   - `serveSwaggerJSON()` - 提供合併後的 JSON
   - `registerSwaggerHandlers()` - 註冊路由

2. **server.go** - HTTP 服務器整合
   - 註冊 Swagger 路由
   - 更新啟動訊息顯示 URL

3. **test-swagger.sh** - 測試腳本
   - 驗證 Swagger UI 可訪問性
   - 驗證 OpenAPI 規格正確性

## 測試

### 自動化測試

執行測試腳本驗證 Swagger 整合：

```bash
./scripts/test-swagger.sh
```

預期輸出：

```text
🧪 Testing Swagger UI Integration
==================================

📌 Starting WEBPRISM server in background...
⏳ Waiting for server to start...

🔍 Testing Swagger endpoints...

1. Testing /swagger.json
   ✅ /swagger.json is accessible

2. Testing /swagger-ui/
   ✅ /swagger-ui/ is accessible

3. Testing /swagger-ui (without trailing slash)
   ✅ /swagger-ui is accessible

✨ Test completed!
```

### 手動測試

#### 測試 1: 驗證 Swagger UI 載入

```bash
curl -s http://localhost:8080/swagger-ui/ | grep -q "WEBPRISM" && echo "✅ OK" || echo "❌ Failed"
```

#### 測試 2: 驗證 OpenAPI 規格

```bash
curl -s http://localhost:8080/swagger.json | jq '.info.title'
# 預期輸出: "WEBPRISM API"
```

#### 測試 3: 驗證所有服務都包含在內

```bash
curl -s http://localhost:8080/swagger.json | jq '.tags[].name'
# 預期輸出:
# "SpecService"
# "AuthService"
# "ProxyService"
# "HealthService"
```

## 疑難排解

### Swagger UI 無法載入

**症狀**: 訪問 `/swagger-ui/` 顯示 404

**解決方案**:

1. 檢查服務器是否正在運行：

```bash
curl http://localhost:8080/swagger.json
```

2. 檢查服務器日誌是否有錯誤：

```bash
make run-server
# 查看日誌輸出
```

### OpenAPI 規格不完整

**症狀**: Swagger UI 中缺少某些端點

**解決方案**:

1. 重新生成 protobuf 和 OpenAPI 規格：

```bash
make proto
```

2. 驗證所有 swagger.json 檔案存在：

```bash
ls -la gen/openapi/v1/
```

預期看到：

```text
-rw-r--r--  auth.swagger.json
-rw-r--r--  health.swagger.json
-rw-r--r--  proxy.swagger.json
-rw-r--r--  spec.swagger.json
```

### CORS 錯誤

**症狀**: 從其他來源訪問時出現 CORS 錯誤

**說明**: CORS 已經預設啟用，允許所有來源。如果仍有問題，檢查瀏覽器控制台的具體錯誤訊息。

## 進階使用

### 匯入到其他工具

#### Postman

1. 下載 OpenAPI 規格：

```bash
curl http://localhost:8080/swagger.json -o webprism-api.json
```

2. 在 Postman 中：
   - File → Import
   - 選擇 `webprism-api.json`
   - 點擊 Import

#### Insomnia

1. 下載 OpenAPI 規格
2. 在 Insomnia 中：
   - Application → Preferences → Data → Import Data
   - 選擇檔案
   - 匯入

#### VS Code REST Client

直接使用 OpenAPI 規格 URL：

```http
# @name GetVersion
GET http://localhost:8080/api/v1/version
```

### 自訂 Swagger UI

如果需要自訂 Swagger UI 的外觀，可以修改 `internal/adapters/http/swagger.go` 中的 HTML 模板。

例如，添加自訂 CSS：

```go
html := `<!DOCTYPE html>
<html lang="en">
<head>
    <style>
        .topbar { background-color: #1e3a8a !important; }
    </style>
    ...
</head>
...
```

## 未來改進

### 計劃中的功能

- [ ] 離線 Swagger UI（嵌入式檔案）
- [ ] API 認證支援（在 Swagger UI 中）
- [ ] 請求/回應範例
- [ ] 多版本 API 文件支援
- [ ] ReDoc 替代介面
- [ ] API 變更日誌整合

### 貢獻

如果您想改進 Swagger UI 整合，歡迎提交 PR：

1. Fork 專案
2. 創建功能分支
3. 提交變更
4. 發送 Pull Request

## 參考資料

- [Swagger UI 官方文檔](https://swagger.io/tools/swagger-ui/)
- [OpenAPI 規格](https://spec.openapis.org/oas/v3.1.0)
- [gRPC-Gateway OpenAPI](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_openapi_output/)

---

**最後更新**: 2026-01-08
**版本**: v1.0.0
**作者**: Jacky Chang
