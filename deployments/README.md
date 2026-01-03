# 部署配置說明

本目錄包含 WEBPRISM 的部署配置檔案。

## 📁 檔案結構

```
deployments/
├── config.yaml          # 應用程式配置（可提交）
└── docker-compose.yaml  # Docker Compose 配置（含 PostgreSQL）
```

## 🔐 配置管理策略

WEBPRISM 採用**分離式配置管理**，遵循 12-factor app 原則：

### 1. config.yaml - 非機密配置（可提交到 Git）

包含所有**非機密**的應用程式配置：
- 伺服器端口
- 連線超時設定
- 資料庫連線參數（不含密碼）
- 日誌設定
- 功能開關

**✅ 可以安全地提交到 Git**

### 2. .env - 機密資訊（不應提交）

只包含**機密資訊**：
- 資料庫密碼
- 加密金鑰
- API 金鑰（未來）

**❌ 不應提交到 Git**（已在 .gitignore 中）

### 3. 配置載入優先順序

1. **環境變數** (WEBPRISM_*) - 最高優先
2. **config.yaml** - 中等優先
3. **程式碼預設值** - 最低優先

## 🚀 使用方式

### 開發環境

```bash
# 1. 複製環境變數範例
cp .env.example .env

# 2. 編輯 .env，填入機密資訊
vim .env

# 3. 啟動服務（會自動載入 config.yaml 和 .env）
make run-server
```

### Docker 環境

```bash
# 使用內建 PostgreSQL（推薦）
docker compose up

# 使用外部 PostgreSQL
# 修改 docker-compose.yaml 中的 WEBPRISM_DATABASE_HOST
# 或通過環境變數覆蓋:
export WEBPRISM_DATABASE_HOST=your-postgres-host
export WEBPRISM_DATABASE_PASSWORD=your-password
docker compose up
```

### 生產環境

推薦使用環境變數覆蓋配置：

```bash
# 方式 1：直接設定環境變數
export WEBPRISM_DATABASE_PASSWORD="production-password"
export WEBPRISM_SECURITY_ENCRYPTION_KEY="production-key-32bytes-exactly"
./bin/webprism-server

# 方式 2：使用 systemd 環境檔
# /etc/webprism/environment
WEBPRISM_DATABASE_PASSWORD=xxx
WEBPRISM_SECURITY_ENCRYPTION_KEY=xxx

# /etc/systemd/system/webprism.service
[Service]
EnvironmentFile=/etc/webprism/environment
ExecStart=/usr/local/bin/webprism-server
```

## 🔧 配置覆蓋範例

### 覆蓋資料庫連線

```bash
# .env
WEBPRISM_DATABASE_HOST=production-db.example.com
WEBPRISM_DATABASE_PASSWORD=secure-password
```

### 覆蓋日誌等級

```bash
# .env
WEBPRISM_LOGGING_LEVEL=debug
WEBPRISM_LOGGING_FORMAT=console
```

### 覆蓋伺服器端口

```bash
# .env
WEBPRISM_SERVER_HTTP_PORT=8888
WEBPRISM_SERVER_GRPC_PORT=9999
```

## 🛡️ 安全最佳實踐

### ✅ DO（應該做的）

1. **config.yaml 只放非機密資訊**
   - 端口號、超時時間等
   - 可以安全地提交到 Git

2. **機密資訊使用環境變數**
   - 資料庫密碼
   - 加密金鑰
   - API 金鑰

3. **生產環境使用強加密金鑰**
   ```bash
   # 生成 32 bytes 隨機金鑰
   openssl rand -base64 32 | head -c 32
   ```

4. **使用 .env.example 作為範本**
   - 提交到 Git
   - 不包含實際機密值

### ❌ DON'T（不應該做的）

1. **不要把密碼提交到 Git**
   - 不要在 config.yaml 中放密碼
   - 不要提交 .env 檔案

2. **不要在生產環境使用開發金鑰**
   ```yaml
   # ❌ 不要這樣做
   encryption_key: "dev-key-12345..."
   ```

3. **不要在程式碼中硬編碼機密**
   ```go
   // ❌ 不要這樣做
   const password = "secret123"
   ```

## 🔍 驗證配置

```bash
# 檢查配置檔案語法
yamllint deployments/config.yaml

# 測試配置載入
go run cmd/server/main.go --help

# 檢查環境變數
env | grep WEBPRISM
```

## 📚 參考資料

- [12-Factor App - Config](https://12factor.net/config)
- [Viper Configuration](https://github.com/spf13/viper)
- [Docker Compose Environment Variables](https://docs.docker.com/compose/environment-variables/)
