# WEBPRISM 快速開始指南

**僅需 5 分鐘！** 從零到啟動 WEBPRISM。

---

## ⚠️ 第一次使用？先安裝工具！

如果是第一次設定，需要先安裝必要工具：

```bash
# 1. 檢查工具是否已安裝
chmod +x scripts/check-tools.sh
./scripts/check-tools.sh

# 2. 如果缺少工具，安裝 protoc
brew install protobuf

# 3. 安裝 Go 插件
make install-tools

# 4. 再次檢查
./scripts/check-tools.sh
```

詳細安裝說明請參考：[INSTALL_TOOLS.md](INSTALL_TOOLS.md)

---

## 🚀 快速開始（使用 Docker）

如果您只想快速測試，使用內建的 PostgreSQL：

```bash
# 1. 初始化專案（生成程式碼、下載依賴）
make setup

# 2. 啟動所有服務（PostgreSQL + WEBPRISM）
make docker-up

# 3. 測試
curl http://localhost:8080/api/v1/specs
```

> **💡 提示**: 執行 `make setup` 前,IDE 可能顯示 import 錯誤,這是正常的。`make setup` 會生成 `gen/` 和 `internal/mocks/` 目錄。

完成！ 🎉

---

## 🔧 進階設定（使用本機 PostgreSQL）

### 前置步驟（僅需執行一次）

```bash
# 1. 設定環境變數
cp .env.example .env
# 編輯 .env 檔案，填入您的 PostgreSQL 資訊
nano .env

# 2. 載入環境變數
export $(cat .env | grep -v '^#' | xargs)

# 3. 安裝工具
brew install postgresql@16 protobuf golangci-lint golang-migrate

# 4. 啟動 PostgreSQL
brew services start postgresql@16

# 5. 創建資料庫
createdb webprism

# 6. 完整設定專案（安裝 Go 工具、生成程式碼、下載依賴）
make setup

# 7. 執行 Migration
make migrate-up

# 8. 編譯
make build
```

### 每次開發時執行

```bash
# 載入環境變數
export $(cat .env | grep -v '^#' | xargs)

# 啟動服務
make run-server
```

---

## 📝 第一次 API 測試

### 使用 CLI

```bash
# 1. 檢查版本
./bin/webprism version

# 2. 上傳 OpenAPI 規格（需要準備一個 JSON 檔案）
./bin/webprism spec upload petstore.json \
  --name "Petstore API" \
  --version "1.0.0"

# 3. 列出規格
./bin/webprism spec list

# 4. 設定認證
./bin/webprism auth set <spec-id> \
  --type bearer \
  --credentials token=your-token

# 5. 執行代理請求
./bin/webprism proxy call <spec-id> getPetById \
  --param petId=123
```

### 使用 HTTP API

```bash
# 1. 列出規格
curl http://localhost:8080/api/v1/specs

# 2. 上傳規格
curl -X POST http://localhost:8080/api/v1/specs \
  -H "Content-Type: application/json" \
  -d @petstore.json

# 3. 健康檢查
curl -X POST http://localhost:8080/api/v1/specs/<spec-id>/health-check
```

---

## 🆘 遇到問題？

### Protobuf 生成失敗

```bash
# 重新安裝工具
make install-tools

# 清理後重新生成
rm -rf gen/go/v1/*
make proto
```

### PostgreSQL 連線失敗

```bash
# 檢查服務
brew services list | grep postgresql

# 重啟服務
brew services restart postgresql@16
```

### 編譯錯誤

```bash
# 檢查未使用的變數
golangci-lint run ./...

# 格式化程式碼
make fmt
```

---

## 📚 更多資訊

- 完整設定: [MACOS_SETUP.md](MACOS_SETUP.md)
- 問題修正: [ISSUES_FIXED.md](ISSUES_FIXED.md)
- 清理指南: [CLEANUP_AND_FIX.md](CLEANUP_AND_FIX.md)
- 完整文件: [V1_COMPLETE_SUMMARY.md](V1_COMPLETE_SUMMARY.md)

---

**開始開發！** 🚀
