# MCP 測試環境快速設置指南

## 🚀 快速開始（3 步驟）

### 1️⃣ 安裝依賴

```bash
cd tests/mcp
npm install
cd ../..
```

### 2️⃣ 確保環境就緒

```bash
# 編譯 MCP server
make build

# 啟動資料庫
make db-up

# 確認 .env 已配置
cat .env | grep WEBPRISM
```

### 3️⃣ 執行測試

```bash
node tests/mcp/test-mcp.js
```

## ✅ 預期結果

```
🎉 所有測試通過！MCP 介面可用！
```

## 📚 完整文件

- **測試環境說明**: [README.md](README.md)
- **測試報告**: [../../docs/testing/MCP_TEST_REPORT.md](../../docs/testing/MCP_TEST_REPORT.md)

## ⚠️ 重要提醒

**必須從專案根目錄執行測試**：

```bash
# ✅ 正確
node tests/mcp/test-mcp.js

# ❌ 錯誤
cd tests/mcp && node test-mcp.js
```
