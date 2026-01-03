# WebPrism 系統啟動與測試指南

**完整的系統啟動、驗證與測試流程**

本指南將引導您從啟動系統到完成第一個 API 代理測試的完整流程。

---

## 📋 前提條件

在開始之前，請確認您已完成：

- ✅ 閱讀並完成 [安裝指南](01-INSTALLATION.md)
- ✅ 所有必要工具已安裝（執行過 `./scripts/check-tools.sh`）
- ✅ 環境變數已設定（`.env` 檔案已配置）

---

## 🚀 系統啟動檢查清單

### 階段一：環境驗證 (5 分鐘)

#### 1.1 檢查工具狀態

```bash
# 執行工具檢查腳本
./scripts/check-tools.sh
```

**預期輸出**:
```
✓ Go version 1.21+ installed
✓ protoc installed
✓ protoc-gen-go installed
✓ protoc-gen-go-grpc installed
✓ golangci-lint installed
✓ mockery installed (optional)
✓ migrate installed
```

**如果有任何 ✗**: 請參考 [工具安裝指南](03-INSTALL_TOOLS.md)

#### 1.2 驗證 Go 環境

```bash
# 檢查 Go 版本
go version

# 檢查 GOPATH
go env GOPATH

# 檢查 Go module 狀態
go env GO111MODULE
```

**預期結果**:
- Go 版本 >= 1.21
- `GO111MODULE` = `on`

#### 1.3 檢查專案依賴

```bash
# 下載並整理依賴
go mod tidy

# 驗證依賴
go mod verify
```

**預期輸出**: `all modules verified`

---

### 階段二：編譯與建置 (5 分鐘)

#### 2.1 生成 Protobuf 程式碼

```bash
# 清理舊的生成檔案
rm -rf gen/go/v1/*

# 生成新的 Protobuf 程式碼
make proto
```

**預期輸出**:
```
Generating protobuf files...
Generated: gen/go/v1/spec_service.pb.go
Generated: gen/go/v1/spec_service_grpc.pb.go
...
```

**驗證生成結果**:
```bash
ls -la gen/go/v1/
```

應該看到 `.pb.go` 和 `_grpc.pb.go` 檔案。

#### 2.2 編譯專案

```bash
# 編譯 server 和 CLI
make build
```

**預期輸出**:
```
Building server...
Building CLI...
Build complete!
```

**驗證編譯結果**:
```bash
# 檢查執行檔
ls -lh bin/

# 測試 CLI 版本
./bin/webprism version
```

---

### 階段三：資料庫準備 (5 分鐘)

#### 3.1 選擇資料庫方案

**方案 A: 使用 Docker (推薦)**

```bash
# 啟動 PostgreSQL (使用 docker-compose)
make docker-up

# 檢查容器狀態
docker ps
```

**方案 B: 使用本機 PostgreSQL**

```bash
# 啟動 PostgreSQL
brew services start postgresql@16

# 驗證服務
brew services list | grep postgresql

# 創建資料庫
createdb webprism

# 執行 Migration
make migrate-up
```

#### 3.2 驗證資料庫連線

```bash
# 測試連線 (使用 psql)
psql -h localhost -U postgres -d webprism -c "SELECT version();"
```

**預期輸出**: PostgreSQL 版本資訊

#### 3.3 檢查資料庫結構

```bash
# 列出所有 table
psql -h localhost -U postgres -d webprism -c "\dt"
```

**預期看到的 tables**:
- `api_specs`
- `auth_configs`
- `health_checks`
- `schema_migrations`

---

### 階段四：啟動服務 (2 分鐘)

#### 4.1 載入環境變數

```bash
# 載入 .env 檔案
export $(cat .env | grep -v '^#' | xargs)

# 驗證環境變數
echo $DB_HOST
echo $DB_PORT
echo $SERVER_PORT
```

#### 4.2 啟動 WebPrism Server

**終端機 1 - 啟動 Server**:
```bash
# 方式 1: 使用 make (推薦)
make run-server

# 方式 2: 直接執行
./bin/server
```

**預期輸出**:
```
2026-01-01T10:00:00+0800 INFO Starting WebPrism Server
2026-01-01T10:00:00+0800 INFO HTTP server listening on :8080
2026-01-01T10:00:00+0800 INFO gRPC server listening on :9090
2026-01-01T10:00:00+0800 INFO Database connection established
```

**保持此終端機運行**

---

### 階段五：系統驗證 (3 分鐘)

#### 5.1 HTTP API 健康檢查

**開啟新的終端機 (終端機 2)**:

```bash
# 測試 HTTP server
curl http://localhost:8080/health

# 測試 API endpoint
curl http://localhost:8080/api/v1/specs
```

**預期輸出**:
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2026-01-01T10:00:00Z"
}
```

```json
{
  "specs": [],
  "total": 0,
  "limit": 10,
  "offset": 0
}
```

#### 5.2 gRPC 服務檢查

```bash
# 使用 grpcurl 測試 (需要安裝 grpcurl)
grpcurl -plaintext localhost:9090 list

# 或使用 CLI
./bin/webprism spec list
```

**預期輸出**: 空的規格列表（因為還沒上傳任何規格）

#### 5.3 CLI 工具驗證

```bash
# 檢查 CLI 版本
./bin/webprism version

# 檢查 CLI 說明
./bin/webprism --help

# 測試 CLI 連線
./bin/webprism spec list
```

---

### 階段六：第一個 API 代理測試 (10 分鐘)

#### 6.1 準備測試用 OpenAPI 規格

創建測試用的 OpenAPI 規格檔案：

```bash
# 創建測試目錄
mkdir -p testdata

# 創建簡單的 OpenAPI 規格
cat > testdata/httpbin.json << 'EOF'
{
  "openapi": "3.0.0",
  "info": {
    "title": "HTTPBin API",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://httpbin.org"
    }
  ],
  "paths": {
    "/get": {
      "get": {
        "operationId": "httpbinGet",
        "summary": "HTTP GET test",
        "responses": {
          "200": {
            "description": "Successful response"
          }
        }
      }
    },
    "/headers": {
      "get": {
        "operationId": "getHeaders",
        "summary": "Get request headers",
        "responses": {
          "200": {
            "description": "Returns request headers"
          }
        }
      }
    }
  }
}
EOF
```

#### 6.2 上傳 API 規格

```bash
# 使用 CLI 上傳
./bin/webprism spec upload testdata/httpbin.json \
  --name "HTTPBin API" \
  --version "1.0.0"

# 或使用 HTTP API
curl -X POST http://localhost:8080/api/v1/specs \
  -H "Content-Type: application/json" \
  -d @testdata/httpbin.json
```

**預期輸出**:
```json
{
  "id": "spec-abc123",
  "name": "HTTPBin API",
  "version": "1.0.0",
  "base_url": "https://httpbin.org",
  "created_at": "2026-01-01T10:00:00Z"
}
```

**記下 spec-id** (例如: `spec-abc123`)

#### 6.3 驗證規格已上傳

```bash
# 列出所有規格
./bin/webprism spec list

# 或使用 HTTP API
curl http://localhost:8080/api/v1/specs
```

**應該看到剛上傳的 HTTPBin API**

#### 6.4 執行代理請求

```bash
# 方式 1: 使用 CLI
./bin/webprism proxy call <spec-id> httpbinGet

# 方式 2: 使用 HTTP API
curl -X POST http://localhost:8080/api/v1/proxy/<spec-id>/httpbinGet \
  -H "Content-Type: application/json"
```

**預期輸出**:
```json
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Host": "httpbin.org",
    "User-Agent": "webprism/1.0"
  },
  "origin": "xxx.xxx.xxx.xxx",
  "url": "https://httpbin.org/get"
}
```

**✅ 恭喜！您已成功完成第一個 API 代理測試**

#### 6.5 測試帶認證的 API（選用）

```bash
# 設定 Bearer Token 認證
./bin/webprism auth set <spec-id> \
  --type bearer \
  --credentials token=test-token-123

# 執行帶認證的請求
./bin/webprism proxy call <spec-id> getHeaders

# 驗證請求中包含 Authorization header
```

---

## 🧪 執行單元測試

### 7.1 執行所有測試

```bash
# 執行所有單元測試
go test ./internal/domain/service/...

# 查看詳細輸出
go test -v ./internal/domain/service/...
```

**預期輸出**:
```
ok      github.com/ChangChunCheng/webprism/internal/domain/service  0.273s
coverage: 95.2% of statements
```

### 7.2 查看測試覆蓋率

```bash
# 生成覆蓋率報告
go test -coverprofile=coverage.out ./internal/domain/service/

# 查看覆蓋率摘要
go tool cover -func=coverage.out | tail -1

# 生成 HTML 報告
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

**預期覆蓋率**: 95.2%

---

## ✅ 系統啟動驗證檢查表

完成以下所有項目，確認系統正常運行：

### 環境檢查

- [ ] 所有工具已安裝並通過檢查
- [ ] Go 模組依賴已下載和驗證
- [ ] Protobuf 程式碼已生成
- [ ] 專案已成功編譯

### 資料庫檢查

- [ ] PostgreSQL 服務運行中
- [ ] 資料庫已創建
- [ ] Migration 已執行
- [ ] 資料庫連線測試成功

### 服務檢查

- [ ] WebPrism Server 啟動成功
- [ ] HTTP API (8080) 可訪問
- [ ] gRPC Service (9090) 可訪問
- [ ] Health endpoint 返回正常
- [ ] CLI 工具運行正常

### 功能檢查

- [ ] 可以上傳 OpenAPI 規格
- [ ] 可以列出已上傳的規格
- [ ] 可以執行代理請求
- [ ] 可以設定認證配置
- [ ] 單元測試全部通過

---

## 🐛 常見問題排除

### 問題 1: Port already in use

**症狀**:
```
Error: bind: address already in use
```

**解決方案**:
```bash
# 查找佔用 port 的程序
lsof -i :8080
lsof -i :9090

# 終止程序
kill -9 <PID>
```

### 問題 2: Database connection failed

**症狀**:
```
Error: failed to connect to database
```

**解決方案**:
```bash
# 檢查 PostgreSQL 狀態
brew services list | grep postgresql

# 重啟 PostgreSQL
brew services restart postgresql@16

# 檢查環境變數
echo $DB_HOST
echo $DB_PORT
echo $DB_USER
```

### 問題 3: Protobuf 生成失敗

**症狀**:
```
protoc-gen-go: program not found or is not executable
```

**解決方案**:
```bash
# 重新安裝 Protobuf 插件
make install-tools

# 確認 PATH 設定
echo $PATH | grep go/bin

# 如果需要，添加到 PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### 問題 4: CLI 找不到 server

**症狀**:
```
Error: failed to connect to server: connection refused
```

**解決方案**:
```bash
# 確認 server 正在運行
ps aux | grep webprism

# 檢查 server 監聽的 port
lsof -i :8080
lsof -i :9090

# 確認環境變數
echo $WEBPRISM_SERVER_ADDR
```

### 問題 5: 測試失敗

**症狀**:
```
FAIL    github.com/ChangChunCheng/webprism/internal/domain/service
```

**解決方案**:
```bash
# 查看詳細錯誤
go test -v ./internal/domain/service/

# 清理並重新測試
go clean -testcache
go test ./internal/domain/service/...

# 檢查是否有程式碼變更
git status
```

---

## 📊 系統監控

### 查看系統日誌

```bash
# Server 日誌 (如果使用 Docker)
docker logs webprism-server

# 或查看檔案日誌
tail -f server.log
```

### 查看資料庫狀態

```bash
# 連接到資料庫
psql -h localhost -U postgres -d webprism

# 查看資料
SELECT * FROM api_specs LIMIT 5;
SELECT * FROM auth_configs LIMIT 5;
SELECT * FROM health_checks LIMIT 5;
```

### 效能監控

```bash
# 查看 server 資源使用
top -pid $(pgrep webprism)

# 或使用 htop
htop -p $(pgrep webprism)
```

---

## 🔄 重啟系統流程

### 完整重啟

```bash
# 1. 停止所有服務
# Ctrl+C 停止 server

# 如果使用 Docker
make docker-down

# 2. 清理
make clean

# 3. 重新編譯
make build

# 4. 重啟資料庫 (如果使用 Docker)
make docker-up

# 5. 重新載入環境變數
export $(cat .env | grep -v '^#' | xargs)

# 6. 啟動 server
make run-server
```

### 快速重啟 (僅重啟 server)

```bash
# 1. 停止 server (Ctrl+C)

# 2. 重新啟動
make run-server
```

---

## 📚 下一步

完成系統啟動和測試後，您可以：

1. **深入學習 API 使用**
   - 查看 [V1 完整功能總結](../development/V1_COMPLETE_SUMMARY.md)
   - 研究更多 API 範例

2. **執行完整測試套件**
   - 閱讀 [測試快速指南](../testing/QUICK_START.md)
   - 查看 [測試完成報告](../testing/COMPLETION_SUMMARY.md)

3. **開始實際應用**
   - 上傳您自己的 OpenAPI 規格
   - 配置實際的 API 認證
   - 整合到您的應用程式中

4. **參與開發**
   - 查看 [當前狀態](../development/CURRENT_STATUS.md)
   - 了解未來規劃

---

## 📞 需要協助？

- 🐛 **遇到問題**: 查看上方的「常見問題排除」
- 📖 **查看文件**: 回到 [文件中心](../README.md)
- 🔍 **搜尋解決方案**: 查看 [已修復問題](../development/ISSUES_FIXED.md)

---

**祝您使用順利！** 🚀
