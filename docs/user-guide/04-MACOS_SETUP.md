# macOS 本機開發設定指南

本指南說明如何在 macOS 上設定 WEBPRISM 開發環境，使用本機 PostgreSQL。

## 前置需求

### 1. 安裝 Homebrew（如果還沒有）

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. 安裝必要工具

```bash
# PostgreSQL
brew install postgresql@16

# Protobuf
brew install protobuf

# golangci-lint (程式碼檢查)
brew install golangci-lint

# golang-migrate (資料庫遷移工具)
brew install golang-migrate
```

### 3. 啟動 PostgreSQL

```bash
# 啟動 PostgreSQL 服務
brew services start postgresql@16

# 檢查狀態
brew services list | grep postgresql
```

## 專案設定

### 1. 清理專案結構

```bash
cd /Users/ChangChunCheng/projects/webprism
```

### 2. 設定環境變數

```bash
# 複製環境變數範例
cp .env.example .env

# 編輯 .env 檔案
nano .env
```

修改以下設定（如果使用 Homebrew 安裝的 PostgreSQL）：

```bash
WEBPRISM_DATABASE_HOST=localhost
WEBPRISM_DATABASE_PORT=5432
WEBPRISM_DATABASE_USER=your_username  # 通常是你的 macOS 使用者名稱
WEBPRISM_DATABASE_PASSWORD=           # Homebrew 預設無密碼
WEBPRISM_DATABASE_DATABASE=webprism
WEBPRISM_DATABASE_SSL_MODE=disable

# 生成安全的加密金鑰（32 bytes）
WEBPRISM_SECURITY_ENCRYPTION_KEY=$(openssl rand -base64 32 | head -c 32)
```

載入環境變數：

```bash
# 方法 1: 使用 export
export $(cat .env | grep -v '^#' | xargs)

# 方法 2: 使用 source（如果是 .env.sh）
source .env
```

### 3. 創建資料庫

```bash
# 連接到 PostgreSQL
psql postgres

# 創建資料庫和使用者
CREATE DATABASE webprism;
CREATE USER webprism WITH PASSWORD 'webprism';
GRANT ALL PRIVILEGES ON DATABASE webprism TO webprism;

# 如果使用你的使用者帳號，只需創建資料庫
CREATE DATABASE webprism;

# 離開 psql
\q
```

或使用命令列：

```bash
# 創建資料庫
createdb webprism

# 驗證
psql -l | grep webprism
```

### 4. 安裝 Go 工具

```bash
# 安裝 protoc-gen-go 和相關工具
make install-tools
```

### 5. 生成 Protobuf 程式碼

```bash
# 生成程式碼
make proto

# 驗證生成結果
ls -la gen/go/v1/
```

### 6. 下載 Go 依賴

```bash
# 下載依賴
go mod tidy

# 驗證
go mod verify
```

### 7. 執行資料庫遷移

```bash
# 檢查 migration 設定
cat Makefile | grep migrate-up

# 執行 migration
make migrate-up

# 驗證表結構
psql webprism -c "\dt"
```

應該看到以下表：
- api_specs
- auth_configs
- health_checks
- schema_migrations

### 8. 編譯專案

```bash
# 編譯所有二進位檔案
make build

# 驗證
ls -la bin/
```

應該看到：
- webprism-server
- webprism-mcp
- webprism

### 9. 啟動服務

#### 方法 A: 使用 Makefile（推薦）

```bash
# 終端 1: 啟動 HTTP + gRPC Server
make run-server

# 終端 2: 啟動 MCP Server（可選）
make run-mcp
```

#### 方法 B: 直接執行二進位檔案

```bash
# 啟動 server
./bin/webprism-server

# 使用自訂配置
./bin/webprism-server --config ./configs/config.local.yaml
```

### 10. 驗證服務

```bash
# 檢查 HTTP 服務
curl http://localhost:8080/api/v1/specs

# 檢查 gRPC 服務（需要 grpcurl）
brew install grpcurl
grpcurl -plaintext localhost:9090 list

# 使用 CLI 測試
./bin/webprism version
```

## 開發工作流程

### 修改 Protobuf 定義後

```bash
make proto
make build
```

### 修改程式碼後

```bash
# 格式化程式碼
make fmt

# 執行 lint
make lint

# 執行測試
make test

# 重新編譯
make build
```

### 新增 Migration

```bash
# 創建新的 migration 檔案
make migrate-create NAME=add_new_table

# 編輯生成的 .sql 檔案
# internal/adapters/storage/postgres/migrations/000002_add_new_table.up.sql
# internal/adapters/storage/postgres/migrations/000002_add_new_table.down.sql

# 執行 migration
make migrate-up

# 回滾（如需要）
make migrate-down
```

## 常見問題

### 問題 1: PostgreSQL 連線失敗

```bash
# 檢查 PostgreSQL 是否運行
brew services list | grep postgresql

# 重啟 PostgreSQL
brew services restart postgresql@16

# 檢查連線
psql postgres -c "SELECT version();"
```

### 問題 2: go mod tidy 錯誤

確認 `go.mod` 中有 replace 指令：

```go
replace github.com/ChangChunCheng/webprism => ./
```

### 問題 3: Protobuf 生成失敗

```bash
# 檢查 protoc 版本
protoc --version

# 重新安裝工具
make install-tools

# 清理後重新生成
rm -rf gen/go/v1/*
make proto
```

### 問題 4: Migration 失敗

```bash
# 檢查連線字串
echo "postgres://$WEBPRISM_DATABASE_USER:$WEBPRISM_DATABASE_PASSWORD@$WEBPRISM_DATABASE_HOST:$WEBPRISM_DATABASE_PORT/$WEBPRISM_DATABASE_DATABASE?sslmode=$WEBPRISM_DATABASE_SSL_MODE"

# 手動執行 migration
migrate -path internal/adapters/storage/postgres/migrations \
  -database "postgres://webprism:webprism@localhost:5432/webprism?sslmode=disable" \
  up
```

### 問題 5: 端口被占用

```bash
# 檢查端口使用情況
lsof -i :8080
lsof -i :9090

# 修改配置檔案使用不同端口
# configs/config.yaml
```

## 使用 Docker（可選）

如果想使用 Docker 搭配本機 PostgreSQL：

```bash
# 使用外部 PostgreSQL 的 docker-compose
docker-compose -f docker-compose.external-pg.yaml up

# 或設定環境變數後啟動
export WEBPRISM_DATABASE_HOST=host.docker.internal
docker-compose -f docker-compose.external-pg.yaml up
```

## 除錯技巧

### 啟用 Debug 日誌

```bash
export WEBPRISM_LOGGING_LEVEL=debug
make run-server
```

### 使用 VSCode 除錯

在 `.vscode/launch.json` 新增：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "env": {
        "WEBPRISM_DATABASE_HOST": "localhost",
        "WEBPRISM_DATABASE_USER": "webprism",
        "WEBPRISM_SECURITY_ENCRYPTION_KEY": "dev-key-12345678901234567890123"
      }
    }
  ]
}
```

## 下一步

完成設定後，請參考：
- [V1_COMPLETE_SUMMARY.md](./V1_COMPLETE_SUMMARY.md) - 完整功能說明
- [CLEANUP_AND_FIX.md](./CLEANUP_AND_FIX.md) - 清理指南
- [README.md](./README.md) - 專案簡介

開始開發！🚀
