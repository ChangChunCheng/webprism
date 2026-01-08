# WEBPRISM Scripts 說明

本目錄包含 WEBPRISM 專案的各種腳本工具。

---

## 📋 目錄

- [測試腳本](#測試腳本)
- [設置腳本](#設置腳本)
- [工具腳本](#工具腳本)
- [常用工作流程](#常用工作流程)

---

## 🧪 測試腳本

### test-petstore-full.sh ⭐ 推薦

**功能**: Petstore API 完整整合測試

**狀態**: ✅ 所有測試通過

**描述**:
使用真實的 Swagger Petstore API，測試 WEBPRISM 的完整工作流程，包含 OpenAPI 3.0.4 和 Swagger 2.0 支援。

**測試步驟** (10 steps):
1. ✅ 健康檢查
2. ✅ 上傳 OpenAPI 3.0 規格
3. ✅ 驗證 Spec 上傳成功
4. ✅ 配置 API Key 認證
5. ✅ 執行代理請求 #1 (`addPet` - POST with Body)
6. ✅ 執行代理請求 #2 (`findPetsByStatus` - GET with Query)
7. ✅ 執行代理請求 #3 (`getPetById` - GET with Path Param)
8. ✅ 上傳 Swagger 2.0 規格
9. ✅ 配置 Bearer Token 認證
10. ✅ 執行代理請求 #4 (Swagger 2.0 測試)

**使用方法**:
```bash
# 1. 啟動資料庫
make db-up

# 2. 啟動服務器 (另一個終端)
make run-server

# 3. 執行測試
./scripts/test-petstore-full.sh
```

**輸出**:
- 彩色終端輸出，清楚標示 ✅ 通過 / ❌ 失敗
- 自動保存日誌到 `/tmp/webprism-test-petstore-full-YYYYMMDD_HHMMSS.log`
- 顯示每個步驟的詳細回應
- 測試完成後提供清理指令

**測試涵蓋範圍**:

- ✅ OpenAPI 3.0.4 支援
- ✅ Swagger 2.0 支援（自動轉換為 OpenAPI 3.0）
- ✅ API Key 認證（Header）
- ✅ Bearer Token 認證
- ✅ Path Parameters
- ✅ Query Parameters
- ✅ Request Body
- ✅ GET, POST 方法

**需求**:

- WEBPRISM 服務器正在運行 (`http://localhost:8080`)
- PostgreSQL 資料庫運行中
- `curl` 和 `jq` 命令行工具已安裝

**相關文件**:

- 詳細指南: [docs/testing/INTEGRATION_TEST_GUIDE.md](../docs/testing/INTEGRATION_TEST_GUIDE.md)
- OpenAPI 3.0 規格: [testdata/petstore-openapi3.json](../testdata/petstore-openapi3.json)
- Swagger 2.0 規格: [testdata/petstore-full.json](../testdata/petstore-full.json)

---

### test-swagger.sh (已棄用)

**功能**: Swagger UI 整合測試

**狀態**: ⚠️ 已被 `test-petstore-full.sh` 取代

**描述**:
簡單的測試腳本，驗證 Swagger UI 和 OpenAPI 文件端點是否可訪問。

**使用方法**:
```bash
./scripts/test-swagger.sh
```

**建議**: 使用 `test-petstore-full.sh` 進行更完整的整合測試

---

## 🛠️ 設置腳本

### setup-docker.sh

**功能**: Docker 環境設置

**描述**:
自動設置 Docker 和 Docker Compose 環境。

**使用方法**:
```bash
./scripts/setup-docker.sh
```

---

### setup-macos.sh

**功能**: macOS 開發環境設置

**描述**:
在 macOS 上安裝所有必要的開發工具和依賴。

**使用方法**:
```bash
./scripts/setup-macos.sh
```

**安裝內容**:

- Homebrew (如果未安裝)
- Go
- PostgreSQL
- Protocol Buffers 編譯器
- 其他開發工具

---

### setup-ubuntu.sh

**功能**: Ubuntu 開發環境設置

**描述**:
在 Ubuntu/Debian 系統上安裝所有必要的開發工具和依賴。

**使用方法**:
```bash
./scripts/setup-ubuntu.sh
```

---

### setup-third-party.sh

**功能**: 第三方工具安裝

**描述**:
安裝 Go 相關的第三方工具，如 protoc-gen-go, grpc-gateway 等。

**使用方法**:
```bash
./scripts/setup-third-party.sh
```

---

## 🔧 工具腳本

### check-tools.sh

**功能**: 檢查開發工具

**描述**:
檢查所有必要的開發工具是否已安裝。

**使用方法**:
```bash
./scripts/check-tools.sh
```

**檢查項目**:

- Go
- Protocol Buffers 編譯器
- PostgreSQL
- Docker
- Make
- jq

---

### detect-platform.sh

**功能**: 平台偵測

**描述**:
偵測當前作業系統和架構。

**使用方法**:
```bash
./scripts/detect-platform.sh
```

**輸出範例**:
```
OS: darwin
ARCH: arm64
```

---

## 🚀 常用工作流程

### 初次設置專案

```bash
# 1. 檢查工具
./scripts/check-tools.sh

# 2. 完整設置專案
make setup
```

### 開發前準備

```bash
# 清理舊產物
make clean-all

# 建置專案
make build

# 啟動資料庫
make db-up

# 啟動服務
make run-server
```

### 執行完整測試

```bash
# 1. 單元測試
make test

# 2. 整合測試
./scripts/test-petstore-full.sh

# 3. 查看測試日誌
cat /tmp/webprism-test-petstore-full-*.log
```

### 清理與重置

```bash
# 清理所有產物
make clean-all

# 重新建置
make build
```

---

## 🐛 常見問題

### Q: 測試失敗，提示 "Server is not running"

**A**: 確保 WEBPRISM 服務器正在運行：
```bash
make run-server
```

### Q: 找不到 jq 命令

**A**: 安裝 jq：
```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt-get install jq
```

### Q: 測試日誌在哪裡？

**A**: 日誌自動保存在 `/tmp/` 目錄中，文件名格式為：

```
/tmp/webprism-test-petstore-full-YYYYMMDD_HHMMSS.log
```

### Q: 如何清理測試資料？

**A**: 測試完成後會顯示清理指令，例如：

```bash
curl -X DELETE http://localhost:8080/api/v1/specs/<spec-id>
```

### Q: 腳本無法執行怎麼辦？

**A**: 添加執行權限：

```bash
chmod +x scripts/*.sh
```

---

## 📚 相關文件

- [整合測試指南](../docs/testing/INTEGRATION_TEST_GUIDE.md) - 詳細測試文件
- [測試概述](../docs/testing/README.md) - 測試系統總覽
- [開發指南](../docs/development/CONTRIBUTING.md) - 開發環境設置
- [快速開始](../docs/GETTING_STARTED.md) - 使用教學
- [README](../README.md) - 專案首頁

---

## 🎯 未來計劃

- [x] 支援 Swagger 2.0 格式的測試
- [x] 增加 Bearer Token 認證測試
- [x] POST/PUT/DELETE 請求測試
- [ ] 增加 Basic Auth 認證測試
- [ ] 增加 OAuth 2.0 認證測試
- [ ] 錯誤處理測試
- [ ] 效能測試腳本
- [ ] CI/CD 整合測試
- [ ] 自動化 E2E 測試

---

**最後更新**: 2026-01-08
