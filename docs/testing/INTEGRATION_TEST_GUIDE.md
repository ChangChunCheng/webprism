# WEBPRISM 整合測試指南

本指南說明如何執行 WEBPRISM 的整合測試，包含與真實外部 API 的完整測試流程。

---

## 🎯 測試概述

WEBPRISM 提供完整的整合測試腳本，測試以下功能：
- ✅ 上傳 OpenAPI/Swagger 規格
- ✅ 配置 API 認證（Bearer Token、API Key）
- ✅ 執行代理請求（Path Parameters、Query Parameters）
- ✅ 驗證回應資料
- ✅ 查詢健康狀態
- ✅ 測試資料清理

---

## 📋 測試腳本說明

### test-petstore-full.sh

完整的 Petstore API 整合測試，包含 10 個測試步驟：

**位置**: `scripts/test-petstore-full.sh`

**測試步驟**:
1. **健康檢查** - 驗證服務正常運行
2. **上傳 OpenAPI 3.0 規格** - 測試 OpenAPI 3.0.4 支援
3. **驗證 Spec 上傳成功** - 確認規格正確解析
4. **配置認證 (API Key)** - 設定 Petstore API Key
5. **執行代理請求 #1** - `addPet` (POST with Body)
6. **執行代理請求 #2** - `findPetsByStatus` (GET with Query Parameters)
7. **執行代理請求 #3** - `getPetById` (GET with Path Parameters)
8. **上傳 Swagger 2.0 規格** - 測試 Swagger 2.0 支援
9. **配置認證 (Bearer Token)** - 設定 Bearer Token
10. **執行代理請求 #4** - 測試 Swagger 2.0 proxy

**日誌輸出**: `/tmp/webprism-test-petstore-full-YYYYMMDD_HHMMSS.log`

---

## 🚀 執行測試

### 前置條件

1. **啟動 WEBPRISM 服務**:
```bash
# 啟動資料庫
make db-up

# 啟動服務（在另一個終端）
make run-server
```

2. **確認服務運行**:
```bash
curl http://localhost:8080/api/v1/health
```

### 執行完整測試

```bash
# 執行 Petstore 完整測試
./scripts/test-petstore-full.sh
```

### 查看測試日誌

```bash
# 日誌會自動儲存在 /tmp/ 目錄
ls -lh /tmp/webprism-test-petstore-full-*.log

# 查看最新日誌
cat $(ls -t /tmp/webprism-test-petstore-full-*.log | head -1)
```

---

## 📊 測試輸出範例

### 成功的測試輸出

```
🧪 WEBPRISM Petstore Integration Test
======================================

📋 Test Configuration
  Base URL: http://localhost:8080/api/v1
  Log File: /tmp/webprism-test-petstore-full-20260108_220842.log

────────────────────────────────────────────────────────────

✅ Step 1: Health Check - PASSED
📊 Response:
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2026-01-08T22:08:42Z"
}

────────────────────────────────────────────────────────────

✅ Step 2: Upload OpenAPI 3.0 Spec - PASSED
📋 Spec ID: 01JGWQ9...
📝 Spec Details:
{
  "id": "01JGWQ9...",
  "name": "Swagger Petstore - OpenAPI 3.0",
  "version": "1.0.20",
  "baseUrl": "https://petstore3.swagger.io/api/v3"
}

────────────────────────────────────────────────────────────

✅ Step 3: Get Uploaded Spec - PASSED

────────────────────────────────────────────────────────────

✅ Step 4: Configure API Key Auth - PASSED
🔐 Auth Config ID: 01JGWQ9...

────────────────────────────────────────────────────────────

✅ Step 5: Execute Proxy Request #1 (addPet) - PASSED
📊 Response Status: 200
📦 Response Body:
{
  "id": 9223372036854753000,
  "name": "doggie",
  "category": {
    "id": 1,
    "name": "Dogs"
  },
  "photoUrls": ["string"],
  "tags": [
    {
      "id": 0,
      "name": "string"
    }
  ],
  "status": "available"
}

────────────────────────────────────────────────────────────

✅ Step 6: Execute Proxy Request #2 (findPetsByStatus) - PASSED
📊 Response Status: 200
📝 Note: Response contains array of pets (truncated in log)

────────────────────────────────────────────────────────────

✅ Step 7: Execute Proxy Request #3 (getPetById) - PASSED
📊 Response Status: 200
📦 Response Body:
{
  "id": 1,
  "category": {
    "id": 1,
    "name": "cat"
  },
  "name": "dog",
  "photoUrls": [],
  "tags": [],
  "status": "sold"
}

────────────────────────────────────────────────────────────

✅ Step 8: Upload Swagger 2.0 Spec - PASSED
📋 Spec ID: 01JGWQA...

────────────────────────────────────────────────────────────

✅ Step 9: Configure Bearer Token Auth - PASSED
🔐 Auth Config ID: 01JGWQA...

────────────────────────────────────────────────────────────

✅ Step 10: Execute Proxy Request #4 (getPetById via Swagger 2.0) - PASSED
📊 Response Status: 200

════════════════════════════════════════════════════════════

✅ All tests passed! 🎉

📊 Test Summary:
   Total Steps: 10
   Passed: 10
   Failed: 0
   Duration: ~15 seconds

📚 View detailed log at: /tmp/webprism-test-petstore-full-20260108_220842.log
```

---

## 🧪 測試涵蓋範圍

### OpenAPI 版本支援

| 版本 | 狀態 | 測試檔案 |
|------|------|----------|
| OpenAPI 3.0.4 | ✅ 支援 | `testdata/petstore-openapi3.json` |
| Swagger 2.0 | ✅ 支援 | `testdata/petstore-full.json` |
| OpenAPI 3.1.x | ✅ 支援 | (自動轉換) |

### 認證類型支援

| 認證類型 | 狀態 | 測試步驟 |
|----------|------|----------|
| API Key (Header) | ✅ 支援 | Step 4 |
| Bearer Token | ✅ 支援 | Step 9 |
| Basic Auth | ✅ 支援 | (已單元測試) |
| OAuth2 | 🚧 計畫中 | - |

### 參數類型支援

| 參數類型 | 狀態 | 測試步驟 |
|----------|------|----------|
| Path Parameters | ✅ 支援 | Step 7 |
| Query Parameters | ✅ 支援 | Step 6 |
| Header Parameters | ✅ 支援 | Step 4 (API Key) |
| Request Body | ✅ 支援 | Step 5 |

### HTTP 方法支援

| 方法 | 狀態 | 測試步驟 |
|------|------|----------|
| GET | ✅ 支援 | Steps 6, 7, 10 |
| POST | ✅ 支援 | Step 5 |
| PUT | ✅ 支援 | (已單元測試) |
| DELETE | ✅ 支援 | (已單元測試) |
| PATCH | ✅ 支援 | (已單元測試) |

---

## 🔧 自訂測試

### 使用自己的 API

1. **準備 OpenAPI 規格檔案**:
```bash
# 將您的 OpenAPI 規格放入 testdata/
cp your-api-spec.json testdata/
```

2. **修改測試腳本**:
```bash
# 複製測試腳本範本
cp scripts/test-petstore-full.sh scripts/test-your-api.sh

# 編輯腳本，修改以下變數：
# - SPEC_FILE: 您的規格檔案路徑
# - BASE_URL: 您的 API 基礎 URL
# - AUTH_TYPE: 認證類型
# - 測試請求內容
```

3. **執行測試**:
```bash
chmod +x scripts/test-your-api.sh
./scripts/test-your-api.sh
```

---

## 🐛 故障排除

### 問題 1: 連線失敗

**症狀**:
```
❌ Step 1: Health Check - FAILED
curl: (7) Failed to connect to localhost port 8080
```

**解決方案**:
```bash
# 確認服務是否運行
ps aux | grep webprism-server

# 如果沒有運行，啟動服務
make run-server

# 確認端口沒有被占用
lsof -i :8080
```

### 問題 2: 規格上傳失敗

**症狀**:
```
❌ Step 2: Upload Spec - FAILED
Error: invalid OpenAPI format
```

**解決方案**:
```bash
# 驗證 OpenAPI 規格格式
# 使用線上工具: https://editor.swagger.io/

# 檢查規格檔案是否為有效 JSON
jq . testdata/your-spec.json

# 確認規格包含必要欄位
jq '.openapi, .swagger, .info, .paths' testdata/your-spec.json
```

### 問題 3: 代理請求失敗

**症狀**:
```
❌ Step 7: Execute Proxy Request - FAILED
Response Status: 404
```

**可能原因**:
1. **Operation ID 不存在**: 檢查 OpenAPI 規格中的 `operationId`
2. **參數錯誤**: 確認參數名稱和類型正確
3. **Base URL 錯誤**: 確認 OpenAPI 規格中的 server URL 正確

**解決方案**:
```bash
# 查看完整日誌
cat /tmp/webprism-test-petstore-full-*.log

# 檢查服務日誌
# (如果服務在前台運行，直接查看終端輸出)

# 驗證 Operation ID 存在
jq '.paths[][] | select(.operationId != null) | .operationId' testdata/your-spec.json
```

### 問題 4: 認證失敗

**症狀**:
```
Response Status: 401
Error: Unauthorized
```

**解決方案**:
```bash
# 確認認證配置正確
curl -X GET "http://localhost:8080/api/v1/auth/$SPEC_ID"

# 檢查 API Key/Token 是否有效
# (嘗試直接調用外部 API 驗證)

# 更新認證配置
curl -X PUT "http://localhost:8080/api/v1/auth/$SPEC_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "authType": "bearer",
    "credentials": {
      "token": "your-new-token"
    }
  }'
```

---

## 📊 效能測試

### 基準測試結果

在標準開發環境（MacBook Pro M1, 16GB RAM）上的測試結果：

| 測試項目 | 平均時間 | 說明 |
|----------|----------|------|
| 健康檢查 | ~50ms | HTTP GET 請求 |
| 上傳規格 | ~200ms | 包含規格解析和資料庫儲存 |
| 配置認證 | ~100ms | 加密儲存認證資料 |
| 代理請求 | ~500-1500ms | 取決於外部 API 回應時間 |
| 完整測試 | ~15秒 | 10 個步驟的總時間 |

---

## 🎯 最佳實踐

### 1. 使用環境變數管理機密資訊

```bash
# 不要在腳本中硬編碼 API Token
export PETSTORE_API_KEY="your-api-key"

# 在腳本中使用環境變數
AUTH_REQUEST=$(cat <<EOF
{
  "specId": "$SPEC_ID",
  "authType": "api_key",
  "credentials": {
    "key": "$PETSTORE_API_KEY",
    "keyLocation": "header",
    "keyName": "api_key"
  }
}
EOF
)
```

### 2. 使用 jq 解析 JSON 回應

```bash
# 擷取特定欄位
SPEC_ID=$(echo "$RESPONSE" | jq -r '.id')

# 檢查欄位是否存在
if echo "$RESPONSE" | jq -e '.error' > /dev/null 2>&1; then
    echo "Error found in response"
fi
```

### 3. 記錄完整日誌

```bash
# 使用 tee 同時輸出到終端和日誌檔
LOG_FILE="/tmp/test-$(date +%Y%m%d_%H%M%S).log"
exec 1> >(tee -a "$LOG_FILE")
exec 2>&1

echo "Test started at $(date)"
```

### 4. 清理測試資料

```bash
# 測試結束後清理
curl -X DELETE "http://localhost:8080/api/v1/specs/$SPEC_ID"

# 或使用 trap 確保清理執行
trap "cleanup_test_data" EXIT
```

---

## 📚 相關文件

- [單元測試指南](./QUICK_START.md) - 執行單元測試
- [測試概述](./README.md) - 測試系統總覽
- [API 文件](../../api/proto/v1/) - gRPC API 定義
- [使用指南](../../docs/user-guide/02-QUICKSTART.md) - 基本使用教學

---

## 🔄 持續改進

### 計畫中的測試增強

- [ ] 自動化 CI/CD 整合測試
- [ ] 效能壓力測試（1000+ RPS）
- [ ] 並發測試
- [ ] 錯誤注入測試
- [ ] 更多真實 API 測試案例

### 貢獻測試案例

如果您有公開的 API 並希望加入測試案例：

1. Fork 本專案
2. 添加您的 OpenAPI 規格到 `testdata/`
3. 創建測試腳本到 `scripts/test-your-api.sh`
4. 提交 Pull Request

---

**最後更新**: 2026-01-08
**維護者**: Jacky Chang
