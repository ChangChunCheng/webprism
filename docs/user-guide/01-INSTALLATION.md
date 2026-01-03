# WEBPRISM 逐步安裝指南

**目標**: 從零開始，一步步完成 WEBPRISM 的完整設定。

**預計時間**: 15-20 分鐘

---

## 📋 前置檢查

```bash
# 確認 Go 已安裝
go version
# 應該顯示: go version go1.21.x 或更高

# 確認在正確的目錄
pwd
# 應該顯示: /Users/ChangChunCheng/projects/webprism
```

---

## 步驟 1: 安裝 Protocol Buffers 編譯器

```bash
# 使用 Homebrew 安裝
brew install protobuf

# 驗證安裝
protoc --version
# 應該顯示: libprotoc 3.x.x 或更高
```

**如果遇到問題**: 參考 [INSTALL_TOOLS.md](INSTALL_TOOLS.md)

---

## 步驟 2: 檢查所有工具

```bash
# 給檢查腳本執行權限
chmod +x scripts/check-tools.sh

# 執行檢查
./scripts/check-tools.sh
```

**預期輸出**: 至少 `protoc` 和 `go` 應該顯示綠色的 ✓

---

## 步驟 3: 安裝 Go 插件

```bash
# 使用 Makefile 安裝所有 Go 插件
make install-tools
```

這會安裝：
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-grpc-gateway
- protoc-gen-openapiv2
- golangci-lint
- gosec

**等待時間**: 約 2-3 分鐘

---

## 步驟 4: 再次檢查工具

```bash
# 確認所有工具都已安裝
./scripts/check-tools.sh
```

**預期結果**:
- ✓ Protocol Buffers 編譯器
- ✓ Go 語言
- ✓ protoc-gen-go
- ✓ protoc-gen-go-grpc
- ✓ protoc-gen-grpc-gateway
- ✓ protoc-gen-openapiv2
- ✓ golangci-lint

---

## 步驟 5: 清理專案結構

---

## 步驟 6: 生成 Protobuf 程式碼

```bash
# 生成程式碼
make proto
```

**預期輸出**:
```
Generating protobuf code...
Protobuf code generated successfully!
```

**驗證生成結果**:
```bash
ls -la gen/go/v1/
```

應該看到類似以下檔案：
- spec.pb.go
- spec_grpc.pb.go
- spec.pb.gw.go
- auth.pb.go
- auth_grpc.pb.go
- (更多檔案...)

---

## 步驟 7: 下載 Go 依賴

```bash
# 下載所有依賴
go mod tidy
```

**預期輸出**: 下載各種套件，不應該有錯誤

**驗證**:
```bash
go mod verify
# 應該顯示: all modules verified
```

---

## 步驟 8: 檢查程式碼（可選）

```bash
# 檢查是否有未使用的變數或其他問題
golangci-lint run ./...
```

**注意**: 可能會有一些警告，這是正常的。主要確認沒有嚴重錯誤。

---

## 步驟 9: 編譯專案

```bash
# 編譯所有二進位檔案
make build
```

**預期輸出**:
```
Building webprism-server...
Building webprism-mcp...
Building webprism-cli...
Build complete!
```

**驗證編譯結果**:
```bash
ls -la bin/
```

應該看到：
- webprism-server
- webprism-mcp
- webprism

---

## 步驟 10: 設定環境變數

```bash
# 複製環境變數範例
cp .env.example .env

# 編輯環境變數檔案
nano .env
```

**最低限度需要設定**:
```bash
WEBPRISM_DATABASE_HOST=localhost
WEBPRISM_DATABASE_USER=your_username
WEBPRISM_DATABASE_PASSWORD=your_password
WEBPRISM_DATABASE_DATABASE=webprism
WEBPRISM_SECURITY_ENCRYPTION_KEY=dev-key-12345678901234567890123
```

**生成安全的加密金鑰**:
```bash
# 生成 32 bytes 的金鑰
openssl rand -base64 32 | head -c 32
```

**載入環境變數**:
```bash
export $(cat .env | grep -v '^#' | xargs)
```

---

## 步驟 11: 設定 PostgreSQL

### 選項 A: 使用 Docker（推薦）

```bash
# 直接啟動 Docker Compose（包含 PostgreSQL）
make docker-up

# 跳到步驟 13
```

### 選項 B: 使用本機 PostgreSQL

```bash
# 安裝 PostgreSQL（如果還沒有）
brew install postgresql@16

# 啟動 PostgreSQL
brew services start postgresql@16

# 驗證運行
brew services list | grep postgresql
```

---

## 步驟 12: 創建資料庫（僅使用本機 PG）

```bash
# 創建資料庫
createdb webprism

# 或使用 psql
psql postgres -c "CREATE DATABASE webprism;"

# 驗證
psql -l | grep webprism
```

---

## 步驟 13: 執行資料庫遷移

```bash
# 執行 migration
make migrate-up
```

**預期輸出**:
```
Running database migrations up...
1/u init_schema (xx.xxxs)
```

**驗證表結構**:
```bash
psql webprism -c "\dt"
```

應該看到：
- api_specs
- auth_configs
- health_checks
- schema_migrations

---

## 步驟 14: 啟動服務

### 使用 Docker

```bash
# 已經在步驟 11 啟動了
# 查看日誌
make docker-logs
```

### 使用本機

```bash
# 啟動 HTTP + gRPC Server
make run-server
```

**預期輸出**:
```
Starting WEBPRISM server...
Starting gRPC server on port 9090
Starting HTTP server on port 8080
```

---

## 步驟 15: 測試服務

**開啟新的終端視窗**，執行以下測試：

```bash
# 1. 測試 HTTP API
curl http://localhost:8080/api/v1/specs

# 預期輸出: {"specs":null} 或 {"specs":[]}

# 2. 測試 CLI
./bin/webprism version

# 預期輸出: WEBPRISM v1.0.0

# 3. 測試 gRPC（需要 grpcurl）
brew install grpcurl
grpcurl -plaintext localhost:9090 list

# 預期輸出: 列出所有 gRPC 服務
```

---

## ✅ 完成檢查清單

確認以下所有項目都完成：

- [ ] protoc 已安裝並可執行
- [ ] Go 插件都已安裝（make install-tools）
- [ ] Protobuf 程式碼已生成（make proto）
- [ ] Go 依賴已下載（go mod tidy）
- [ ] 專案已成功編譯（make build）
- [ ] 環境變數已設定（.env）
- [ ] PostgreSQL 已啟動（Docker 或本機）
- [ ] 資料庫已創建（webprism）
- [ ] Migration 已執行（make migrate-up）
- [ ] 服務已啟動（make run-server 或 docker-up）
- [ ] HTTP API 測試成功
- [ ] CLI 測試成功

---

## 🎉 成功！

如果所有測試都通過，恭喜您已成功設定 WEBPRISM！

### 下一步

1. **閱讀功能文件**: [V1_COMPLETE_SUMMARY.md](V1_COMPLETE_SUMMARY.md)
2. **測試 API**: 參考文件中的「測試範例」章節
3. **上傳第一個規格**: 使用 CLI 或 HTTP API

---

## 🆘 遇到問題？

### 常見問題速查

| 問題 | 解決方案 |
|------|----------|
| protoc: command not found | `brew install protobuf` |
| protoc-gen-go 找不到 | 確認 `$GOPATH/bin` 在 PATH 中 |
| go mod tidy 錯誤 | 確認 go.mod 有 `replace` 指令 |
| Migration 失敗 | 檢查 PostgreSQL 連線和環境變數 |
| 編譯錯誤 | 確認已執行 `make proto` |
| 端口被占用 | 修改 config.yaml 的端口設定 |

### 詳細除錯指南

- [INSTALL_TOOLS.md](INSTALL_TOOLS.md) - 工具安裝問題
- [MACOS_SETUP.md](MACOS_SETUP.md) - macOS 設定問題
- [ISSUES_FIXED.md](ISSUES_FIXED.md) - 已知問題解決

### 獲取幫助

1. 執行 `./scripts/check-tools.sh` 檢查工具狀態
2. 檢查日誌輸出的錯誤訊息
3. 參考相關文件的「常見問題」章節

---

**祝您使用愉快！** 🚀
