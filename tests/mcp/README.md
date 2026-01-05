# MCP 測試環境

此目錄包含 WEBPRISM MCP 介面的測試工具和腳本。

## 📁 目錄結構

```
tests/mcp/
├── README.md                 # 本文件
├── test-mcp.js              # 自動化測試腳本
├── test-mcp-manual.sh       # 手動測試指南
├── package.json             # Node.js 依賴配置
├── package-lock.json        # 鎖定的依賴版本（ignored）
└── node_modules/            # Node.js 依賴（ignored）
```

## 🛠️ 環境設置

### 前置需求

- Node.js 18+ 已安裝
- PostgreSQL 資料庫已啟動（本地或 Docker）
- WEBPRISM MCP server 已編譯（`bin/webprism-mcp`）
- 環境變數已配置（`.env` 檔案）

### 初次設置

```bash
# 1. 進入測試目錄
cd tests/mcp

# 2. 安裝 Node.js 依賴
npm install

# 3. 返回專案根目錄
cd ../..
```

## 🧪 執行測試

### 方法 1: Web 互動測試（推薦 - 有真實感！）

啟動 Web UI 測試伺服器：

```bash
# 從專案根目錄執行
node tests/mcp/web-server.js
```

然後開啟瀏覽器訪問 `http://localhost:3000`

**特點**:
- 🌐 視覺化介面
- 🎯 互動式操作
- 📊 即時結果顯示
- 💡 內建測試範例

詳細說明請參考：[WEB_TEST_GUIDE.md](WEB_TEST_GUIDE.md)

### 方法 2: 自動化測試（快速驗證）

從**專案根目錄**執行：

```bash
node tests/mcp/test-mcp.js
```

**測試內容**:
- ✅ MCP server 啟動驗證
- ✅ 環境變數載入驗證
- ✅ JSON-RPC 協定測試（initialize）
- ✅ 工具列表驗證（6 個工具）
- ✅ 工具名稱和描述驗證

**預期輸出**:
```
=================================
WEBPRISM MCP 基礎連線測試
=================================

📡 正在啟動 MCP server...
✓ 已載入環境變數
  - WEBPRISM_DATABASE_PASSWORD: ***
  - WEBPRISM_SECURITY_ENCRYPTION_KEY: ***

...

✅ 所有 6 個工具都正確註冊！

=================================
測試摘要
=================================
✅ 通過: 7
❌ 失敗: 0
📊 總計: 7

🎉 所有測試通過！MCP 介面可用！
=================================
```

### 方法 3: 使用 MCP Inspector（官方工具）

如果需要互動式測試，可使用官方的 MCP Inspector：

```bash
# 從專案根目錄執行
npx --prefix tests/mcp @modelcontextprotocol/inspector ./bin/webprism-mcp
```

這會開啟瀏覽器視窗，提供視覺化的測試介面。

## 📝 測試腳本說明

### test-mcp.js

**功能**:
- 自動載入專案根目錄的 `.env` 檔案
- 啟動 MCP server（stdio 模式）
- 發送 JSON-RPC 請求並驗證回應
- 檢查所有 6 個工具是否正確註冊

**特點**:
- 相對路徑處理（可從專案根目錄執行）
- 超時保護（30 秒）
- 詳細的錯誤訊息
- 彩色輸出（emoji 標記）

### test-mcp-manual.sh

顯示手動測試的指令和步驟。

## 🔍 已驗證的 MCP Tools

| # | 工具名稱 | 功能 |
|---|----------|------|
| 1 | `upload_api_spec` | 上傳 OpenAPI 規格 |
| 2 | `list_api_specs` | 列出所有規格 |
| 3 | `get_api_spec` | 取得特定規格詳情 |
| 4 | `set_auth` | 設定認證配置 |
| 5 | `call_api` | 執行 API 代理呼叫 |
| 6 | `health_check` | 檢查 API 健康狀態 |

## ⚠️ 注意事項

### 1. 執行位置

**必須從專案根目錄執行測試腳本**：

```bash
# ✅ 正確
node tests/mcp/test-mcp.js

# ❌ 錯誤
cd tests/mcp && node test-mcp.js
```

測試腳本會自動處理相對路徑，找到：
- `../../.env` - 環境變數配置
- `../../bin/webprism-mcp` - MCP server 執行檔

### 2. 環境變數

確保 `.env` 檔案包含必要的配置：

```bash
WEBPRISM_DATABASE_PASSWORD=webprism
WEBPRISM_SECURITY_ENCRYPTION_KEY=dev-key-12345678901234567890abcd  # 必須 32 bytes
```

### 3. 資料庫連線

測試前確保 PostgreSQL 已啟動：

```bash
# Docker 模式
make db-up

# 或檢查本地 PostgreSQL
psql -U webprism -d webprism -c "SELECT 1"
```

### 4. Git 忽略

以下檔案已加入 `.gitignore`，不會被提交：
- `tests/mcp/node_modules/`
- `tests/mcp/package-lock.json`

## 🚀 快速開始

完整的測試流程：

```bash
# 1. 確保在專案根目錄
cd /path/to/webprism

# 2. 編譯 MCP server（如果未編譯）
make build

# 3. 啟動資料庫（如果未啟動）
make db-up

# 4. 設置測試環境（首次）
cd tests/mcp && npm install && cd ../..

# 5. 執行測試
node tests/mcp/test-mcp.js
```

## 📚 相關文件

- **MCP 測試報告**: [../../docs/testing/MCP_TEST_REPORT.md](../../docs/testing/MCP_TEST_REPORT.md)
- **MCP Server 實作**: [../../internal/adapters/mcp/](../../internal/adapters/mcp/)
- **專案文件**: [../../docs/README.md](../../docs/README.md)

## 🐛 故障排除

### 問題 1: "encryption key must be exactly 32 bytes"

**原因**: 環境變數未正確載入

**解決方案**:
```bash
# 檢查 .env 檔案是否存在
cat .env | grep WEBPRISM_SECURITY_ENCRYPTION_KEY

# 確保金鑰長度為 32 bytes
echo -n "dev-key-12345678901234567890abcd" | wc -c  # 應該是 32
```

### 問題 2: "Cannot find module"

**原因**: 未安裝 Node.js 依賴

**解決方案**:
```bash
cd tests/mcp
npm install
cd ../..
```

### 問題 3: 測試超時

**原因**: 資料庫未啟動或連線失敗

**解決方案**:
```bash
# 檢查 PostgreSQL
docker compose ps
# 或
pg_isready -h localhost -p 5432
```

## 📊 測試結果

最新測試結果請參閱：[MCP_TEST_REPORT.md](../../docs/testing/MCP_TEST_REPORT.md)
