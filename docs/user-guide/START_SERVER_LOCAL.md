# 🚀 本機啟動 WebPrism Server 指南

**日期**: 2026-01-01
**狀態**: ✅ 準備就緒

---

## ✅ 前置檢查

已完成以下準備：
- ✅ PostgreSQL 正在運行
- ✅ 資料庫 `webprism` 已創建
- ✅ 所有資料表已建立 (api_specs, auth_configs, health_checks)
- ✅ Server 已編譯完成 (`bin/webprism-server`)
- ✅ CLI 已編譯完成 (`bin/webprism`)
- ✅ 環境變數已設定 (`.env`)

---

## 🎯 啟動 Server

### 方式 1: 前台運行（推薦用於測試）

```bash
# 在專案根目錄執行
./bin/webprism-server
```

**優點**:
- 可以即時看到所有日誌
- 容易發現錯誤
- Ctrl+C 即可停止

**預期輸出**:
```json
{"level":"info","timestamp":"2026-01-01T22:31:55+0800","message":"Starting WEBPRISM server..."}
{"level":"info","timestamp":"2026-01-01T22:31:55+0800","message":"PostgreSQL connection established","host":"127.0.0.1","port":5432,"database":"webprism"}
{"level":"info","timestamp":"2026-01-01T22:31:55+0800","message":"Starting gRPC server","port":9090}
{"level":"info","timestamp":"2026-01-01T22:31:55+0800","message":"Starting HTTP server","port":8080}
{"level":"info","timestamp":"2026-01-01T22:31:55+0800","message":"HTTP server starting","addr":":8080","grpc_port":9090}
{"level":"info","timestamp":"2026-01-01T22:31:55+0800","message":"gRPC server starting","port":9090}
```

### 方式 2: 背景運行

```bash
# 背景運行並將日誌輸出到文件
./bin/webprism-server > webprism.log 2>&1 &

# 記錄 PID
echo $! > webprism.pid

# 查看日誌
tail -f webprism.log
```

**停止服務**:
```bash
# 使用 PID 停止
kill $(cat webprism.pid)

# 或直接 kill
pkill -f webprism-server
```

---

## 🔍 驗證 Server 啟動成功

### 1. 檢查進程

```bash
# 檢查 server 是否在運行
ps aux | grep webprism-server

# 預期看到類似這樣的輸出
# ChangChunCheng  54125  webprism-server
```

### 2. 檢查端口

```bash
# 檢查 HTTP 端口 (8080)
lsof -i :8080

# 檢查 gRPC 端口 (9090)
lsof -i :9090

# 預期輸出：兩個端口都被 webprism-server 佔用
```

### 3. 測試 API

```bash
# 測試 Spec API (應該回傳空陣列或已存在的規格)
curl http://localhost:8080/api/v1/specs

# 預期輸出
[]  # 或已存在的規格列表
```

---

## 🧪 第一個 API 測試

### 測試 1: 上傳 OpenAPI 規格

創建測試規格檔案：

```bash
cat > /tmp/test-api.json << 'EOF'
{
  "name": "JSONPlaceholder Test",
  "version": "1.0.0",
  "spec": {
    "openapi": "3.0.0",
    "info": {
      "title": "JSONPlaceholder",
      "version": "1.0.0"
    },
    "servers": [
      {
        "url": "https://jsonplaceholder.typicode.com"
      }
    ],
    "paths": {
      "/posts/1": {
        "get": {
          "operationId": "getPost",
          "summary": "Get a single post",
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
EOF
```

上傳規格：

```bash
curl -X POST http://localhost:8080/api/v1/specs \
  -H "Content-Type: application/json" \
  -d @/tmp/test-api.json

# 預期輸出（記錄 spec_id）
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "name": "JSONPlaceholder Test",
  "version": "1.0.0",
  ...
}
```

### 測試 2: 執行代理請求

```bash
# 使用上面得到的 spec_id
SPEC_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

# 執行代理請求
curl -X POST http://localhost:8080/api/v1/proxy/$SPEC_ID/getPost

# 預期輸出（JSONPlaceholder 的文章資料）
{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere...",
  "body": "quia et suscipit..."
}
```

✅ **如果看到文章資料，代表系統運行正常！**

---

## 🐛 常見問題排查

### 問題 1: Port 8080 已被佔用

```bash
# 查看誰在使用 port 8080
lsof -i :8080

# 停止佔用的進程
kill -9 <PID>

# 或修改 .env 中的端口
WEBPRISM_HTTP_PORT=8081
```

### 問題 2: 資料庫連線失敗

檢查錯誤日誌中是否有：
```
"error":"failed to connect to database"
```

解決方法：
```bash
# 確認 PostgreSQL 運行中
pg_isready

# 確認資料庫存在
psql -U webprism -d webprism -c "SELECT 1"

# 檢查 .env 設定
cat .env | grep DATABASE
```

### 問題 3: gRPC 連線錯誤

如果看到：
```
"error":"grpc: the client connection is closing"
```

可能是 HTTP gateway 到 gRPC 的連線問題。檢查：
```bash
# gRPC server 是否正常啟動
lsof -i :9090

# 查看完整錯誤日誌
# (前台模式運行時可直接看到)
```

### 問題 4: API 回傳 404

確認 API 路徑正確：
- ✅ `/api/v1/specs` （正確）
- ❌ `/specs` （錯誤）
- ❌ `/health` （可能未實現）

---

## 📝 測試檢查清單

啟動 server 後，請依序檢查：

- [ ] Server 進程正在運行
- [ ] HTTP port 8080 已監聽
- [ ] gRPC port 9090 已監聽
- [ ] 日誌沒有 ERROR
- [ ] 可以成功調用 `/api/v1/specs`
- [ ] 可以成功上傳規格
- [ ] 可以成功執行代理請求

---

## 🎯 完成後

當 server 正常運行後，請參考：

1. **[人機互動測試指南](docs/testing/MANUAL_TESTING_GUIDE.md)**
   - 完整的測試流程
   - 10 個測試階段
   - 詳細驗證步驟

2. **[測試執行總結](TEST_EXECUTION_SUMMARY.md)**
   - 單元測試結果
   - 測試覆蓋率分析
   - 測試計劃

---

## 🔧 Makefile 指令

專案提供了一些 Makefile 指令：

```bash
# 編譯所有組件
make build

# 只編譯 server
make build-server

# 執行測試
make test

# 查看測試覆蓋率
make test-coverage

# 格式化程式碼
make fmt

# Linting
make lint

# 清理建置產物
make clean
```

---

## 📊 監控服務狀態

### 查看日誌

如果使用背景模式運行：
```bash
# 實時查看日誌
tail -f webprism.log

# 過濾錯誤
tail -f webprism.log | grep -i error

# 過濾特定請求
tail -f webprism.log | grep -i "HTTP request"
```

### 查看資料庫

```bash
# 連線到資料庫
psql -U webprism -d webprism

# 查看已上傳的規格
SELECT id, name, version, created_at FROM api_specs;

# 查看認證設定
SELECT spec_id, auth_type, created_at FROM auth_configs;

# 退出
\q
```

---

## 🚀 準備就緒！

現在您可以：

1. **啟動 server**: `./bin/webprism-server`
2. **觀察日誌**: 確認沒有錯誤
3. **執行測試**: 按照測試指南進行
4. **回報問題**: 如果發現任何錯誤

---

**祝測試順利！** 🎉

*最後更新: 2026-01-01*
