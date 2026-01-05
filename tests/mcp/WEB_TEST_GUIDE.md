# 🌐 WEBPRISM MCP Web 互動測試指南

## 🎯 為什麼需要 Web 測試？

與純粹的自動化測試不同，Web 互動測試讓你能：

- ✅ **親自操作** - 視覺化介面，直接點擊和輸入
- ✅ **即時反饋** - 立即看到 MCP 工具的執行結果
- ✅ **真實體驗** - 模擬 AI 如何使用這些工具
- ✅ **快速除錯** - 直接修改參數並重新測試

---

## 🚀 快速啟動

### 1️⃣ 啟動測試伺服器

```bash
# 從專案根目錄執行
node tests/mcp/web-server.js
```

### 2️⃣ 開啟瀏覽器

```
http://localhost:3000
```

### 3️⃣ 開始測試！

看到以下畫面就代表成功了：

```
╔════════════════════════════════════════════╗
║   WEBPRISM MCP 互動測試伺服器              ║
╚════════════════════════════════════════════╝

🌐 Web UI: http://localhost:3000

📋 功能:
  • 視覺化工具列表
  • 互動式工具測試
  • 即時結果顯示
  • JSON 格式化
```

---

## 🎨 Web UI 介面說明

### 頁面布局

```
┌─────────────────────────────────────────────────┐
│  🌐 WEBPRISM MCP 互動測試                       │
│  狀態: ✅ MCP Server 就緒                       │
├─────────────────────┬───────────────────────────┤
│  📋 可用工具        │  🎯 測試工具              │
│                     │                           │
│  • upload_api_spec  │  選擇的工具: list_api... │
│  • list_api_specs   │                           │
│  • get_api_spec     │  參數 (JSON 格式)         │
│  • set_auth         │  ┌─────────────────────┐ │
│  • call_api         │  │ {"limit": 10}       │ │
│  • health_check     │  └─────────────────────┘ │
│                     │                           │
│                     │  [執行工具]               │
├─────────────────────┴───────────────────────────┤
│  📊 執行結果                                    │
│  ┌──────────────────────────────────────────┐  │
│  │ {                                        │  │
│  │   "result": { ... }                      │  │
│  │ }                                        │  │
│  └──────────────────────────────────────────┘  │
├──────────────────────────────────────────────────┤
│  💡 快速測試範例                                │
│  [列出所有規格] [上傳測試規格] [健康檢查]      │
└──────────────────────────────────────────────────┘
```

### 使用步驟

1. **左側面板** - 點擊任一工具
2. **右側面板** - 編輯 JSON 參數（已預填範例）
3. **點擊「執行工具」** - 發送請求
4. **下方面板** - 查看格式化的結果

---

## 📝 測試範例

### 範例 1: 列出所有 API 規格

1. 點擊左側的 `list_api_specs`
2. 參數已自動填入：
   ```json
   {"limit": 10}
   ```
3. 點擊「執行工具」
4. 查看結果：
   ```json
   {
     "result": {
       "content": [
         {
           "type": "text",
           "text": "Found 0 API specs:\n[]"
         }
       ]
     }
   }
   ```

### 範例 2: 上傳 Petstore API 規格

1. 點擊左側的 `upload_api_spec`
2. 修改參數為實際的 OpenAPI 規格：
   ```json
   {
     "name": "Petstore API",
     "version": "1.0.0",
     "spec_data": {
       "openapi": "3.0.0",
       "info": {
         "title": "Petstore API",
         "version": "1.0.0"
       },
       "paths": {
         "/pets": {
           "get": {
             "operationId": "listPets",
             "summary": "List all pets",
             "responses": {
               "200": {
                 "description": "Success"
               }
             }
           }
         }
       }
     }
   }
   ```
3. 點擊「執行工具」
4. 記下回傳的 `spec_id`，後續測試會用到

### 範例 3: 取得規格詳情

1. 點擊左側的 `get_api_spec`
2. 使用上一步取得的 `spec_id`：
   ```json
   {
     "spec_id": "your-spec-id-here"
   }
   ```
3. 點擊「執行工具」
4. 查看完整的規格資訊和可用操作

### 範例 4: 設定認證

1. 點擊左側的 `set_auth`
2. 設定 API Key 認證：
   ```json
   {
     "spec_id": "your-spec-id-here",
     "auth_type": "api_key",
     "credentials": {
       "key": "test-api-key-12345",
       "header": "X-API-Key"
     }
   }
   ```
3. 點擊「執行工具」

### 範例 5: 執行 API 呼叫

1. 點擊左側的 `call_api`
2. 呼叫實際的 API：
   ```json
   {
     "spec_id": "your-spec-id-here",
     "operation_id": "listPets",
     "parameters": {
       "query": {
         "limit": "10"
       }
     }
   }
   ```
3. 點擊「執行工具」
4. 查看代理後的 API 回應

### 範例 6: 健康檢查

1. 點擊左側的 `health_check`
2. 檢查 API 健康狀態：
   ```json
   {
     "spec_id": "your-spec-id-here"
   }
   ```
3. 點擊「執行工具」
4. 查看健康狀態報告

---

## 🔍 完整測試流程

### 情境：測試 JSONPlaceholder API

這是一個公開的測試 API，非常適合用來測試 WEBPRISM。

#### Step 1: 上傳規格

```json
{
  "name": "JSONPlaceholder",
  "version": "1.0.0",
  "spec_data": {
    "openapi": "3.0.0",
    "info": {
      "title": "JSONPlaceholder API",
      "version": "1.0.0"
    },
    "servers": [
      {
        "url": "https://jsonplaceholder.typicode.com"
      }
    ],
    "paths": {
      "/posts": {
        "get": {
          "operationId": "getPosts",
          "summary": "Get all posts",
          "responses": {
            "200": {
              "description": "Success"
            }
          }
        }
      },
      "/posts/{id}": {
        "get": {
          "operationId": "getPost",
          "summary": "Get a post by ID",
          "parameters": [
            {
              "name": "id",
              "in": "path",
              "required": true,
              "schema": {
                "type": "integer"
              }
            }
          ],
          "responses": {
            "200": {
              "description": "Success"
            }
          }
        }
      }
    }
  }
}
```

#### Step 2: 查看上傳的規格

使用 `list_api_specs` 確認規格已上傳。

#### Step 3: 呼叫 API

```json
{
  "spec_id": "從 step 1 取得的 ID",
  "operation_id": "getPosts"
}
```

#### Step 4: 帶參數呼叫

```json
{
  "spec_id": "從 step 1 取得的 ID",
  "operation_id": "getPost",
  "parameters": {
    "path": {
      "id": "1"
    }
  }
}
```

---

## 🎯 測試重點

### 要驗證的功能

1. **工具列表** ✅
   - [ ] 6 個工具都有顯示
   - [ ] 描述正確
   - [ ] 點擊可選擇

2. **參數處理** ✅
   - [ ] JSON 格式驗證
   - [ ] 參數預填
   - [ ] 錯誤提示

3. **執行結果** ✅
   - [ ] 成功時顯示結果
   - [ ] 失敗時顯示錯誤
   - [ ] JSON 格式化正確

4. **完整流程** ✅
   - [ ] 上傳規格 → 列出規格 → 取得詳情
   - [ ] 設定認證 → 呼叫 API
   - [ ] 健康檢查

---

## 🛠️ Web Server 技術細節

### 架構

```
┌─────────────────────────────────────────┐
│         Web UI (瀏覽器)                 │
│  - HTML/CSS/JavaScript                  │
│  - 單頁應用                             │
└─────────────┬───────────────────────────┘
              │ HTTP API
              │ /api/status
              │ /api/tools/list
              │ /api/tools/call
┌─────────────┴───────────────────────────┐
│      Node.js HTTP Server                │
│  - web-server.js                        │
│  - Port 3000                            │
└─────────────┬───────────────────────────┘
              │ JSON-RPC (stdio)
┌─────────────┴───────────────────────────┐
│      MCP Server (Go)                    │
│  - bin/webprism-mcp                     │
│  - 6 tools                              │
└─────────────┬───────────────────────────┘
              │ SQL
┌─────────────┴───────────────────────────┐
│      PostgreSQL                         │
└─────────────────────────────────────────┘
```

### HTTP API 端點

| 端點 | 方法 | 說明 |
|------|------|------|
| `/` | GET | HTML 頁面 |
| `/api/status` | GET | MCP server 狀態 |
| `/api/tools/list` | GET | 列出所有工具 |
| `/api/tools/call` | POST | 呼叫指定工具 |

### 特點

- ✅ 單檔案實作（不需額外依賴）
- ✅ 自動啟動 MCP server
- ✅ 自動載入環境變數
- ✅ 響應式設計（支援手機）
- ✅ 即時狀態更新
- ✅ JSON 語法高亮

---

## 🐛 故障排除

### 問題 1: Web UI 顯示「MCP Server 啟動中...」

**原因**: MCP server 尚未完全啟動

**解決方案**: 等待 3-5 秒，頁面會自動更新

### 問題 2: 工具列表為空

**原因**: MCP server 未正確初始化

**解決方案**:
1. 檢查終端機是否有錯誤訊息
2. 確認資料庫已啟動
3. 確認 `.env` 配置正確

### 問題 3: 執行工具時出現錯誤

**原因**: JSON 格式錯誤或缺少必要參數

**解決方案**:
1. 檢查 JSON 格式是否正確
2. 參考範例修正參數
3. 查看錯誤訊息

### 問題 4: Port 3000 已被占用

**解決方案**:
```bash
# 找出占用的進程
lsof -i :3000

# 或修改 web-server.js 中的 PORT 變數
```

---

## 📚 相關文件

- **測試腳本**: [web-server.js](web-server.js)
- **自動化測試**: [test-mcp.js](test-mcp.js)
- **測試報告**: [../../docs/testing/MCP_TEST_REPORT.md](../../docs/testing/MCP_TEST_REPORT.md)

---

## 🎉 開始測試！

```bash
# 1. 啟動伺服器
node tests/mcp/web-server.js

# 2. 開啟瀏覽器
open http://localhost:3000

# 3. 開始互動測試！
```

祝測試順利！有任何問題歡迎隨時詢問。
