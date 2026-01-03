# 🌐 WEBPRISM

**Web Proxy Request Integration Service Manager**

> 智能 API 整合中介層 - 讓任何人都能輕鬆使用第三方 API

[![Tests](https://img.shields.io/badge/tests-80%20passed-brightgreen)](docs/testing/TESTING_COMPLETE.md)
[![Coverage](https://img.shields.io/badge/coverage-95.2%25-brightgreen)](docs/testing/TESTING_COMPLETE.md)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](go.mod)
[![License](https://img.shields.io/badge/license-TBD-lightgrey)](LICENSE)

---

## 🎯 專案目的

WEBPRISM 是一個智能 API 整合中介層，旨在簡化第三方 API 的使用和管理。

**核心理念**: 上傳 API 文件 → 立即可用

只需上傳 Swagger/OpenAPI 規格文件，WEBPRISM 自動建立代理介面，讓開發者、系統和 AI 模型都能輕鬆調用任何第三方服務。

### 主要功能

- 📄 **文件驅動**: 上傳 OpenAPI 規格即可使用
- 🔐 **統一認證**: 集中管理所有 API 憑證（Bearer、Basic、API Key、OAuth2）
- 🔄 **請求代理**: 安全轉發請求，自動注入認證
- 🚀 **多介面**: REST API、gRPC、CLI、MCP Server
- 📊 **企業級**: 日誌、監控、速率限制、錯誤處理

---

## 🛠️ 技術棧

### 語言與框架

- **語言**: Go 1.21+
- **資料庫**: PostgreSQL 16+
- **API 協定**: REST + gRPC
- **API 文件**: OpenAPI 3.0
- **容器化**: Docker + Docker Compose

### 架構設計

- **六角形架構** (Hexagonal Architecture)
- **介面導向設計** (Interface-Driven Design)
- **依賴反轉原則** (Dependency Inversion Principle)

---

## 📦 相依套件

### 核心依賴

- **[gRPC](https://grpc.io/)** - 高效能 RPC 框架
- **[Protocol Buffers](https://protobuf.dev/)** - 資料序列化
- **[pgx](https://github.com/jackc/pgx)** - PostgreSQL 驅動
- **[Viper](https://github.com/spf13/viper)** - 配置管理
- **[Cobra](https://github.com/spf13/cobra)** - CLI 框架

### 開發工具

- **[mockery](https://github.com/vektra/mockery)** - Mock 生成
- **[golangci-lint](https://github.com/golangci/golangci-lint)** - 程式碼檢查

完整依賴清單請參考 [go.mod](go.mod) 和 [DEPENDENCIES.md](DEPENDENCIES.md)

---

## 🚀 快速開始

### 安裝與啟動

```bash
# Clone 專案
git clone https://github.com/ChangChunCheng/webprism.git
cd webprism

# 初始化專案
make setup

# 配置環境變數
cp .env.example .env

# 啟動服務（推薦混合模式）
make db-up        # 啟動 PostgreSQL (Docker)
make run-server   # 執行 WEBPRISM (本機)
```

### 驗證服務

```bash
curl http://localhost:8080/api/v1/specs
```

### 詳細文件

- **使用者**: 參考 [快速開始指南](docs/GETTING_STARTED.md)
- **開發者**: 參考 [開發指南](docs/development/CONTRIBUTING.md)
- **文件中心**: [docs/README.md](docs/README.md)

---

## 📖 文件導航

所有文件已分類整理至 [docs/](docs/) 目錄：

### 👤 使用者文件

| 文件 | 說明 |
|------|------|
| [快速開始](docs/GETTING_STARTED.md) | 新手入門完整指南 |
| [安裝指南](docs/user-guide/01-INSTALLATION.md) | 詳細安裝步驟 |
| [使用教學](docs/user-guide/02-QUICKSTART.md) | 第一次使用教學 |

### 👨‍💻 開發者文件

| 文件 | 說明 |
|------|------|
| [開發指南](docs/development/CONTRIBUTING.md) | 開發環境設置與工作流程 |
| [專案狀態](docs/development/CURRENT_STATUS.md) | 當前開發進度 |
| [測試指南](docs/testing/QUICK_START.md) | 測試執行與覆蓋率 |

---

## 📄 授權

MIT License - 無限制開放使用

查看 [LICENSE](LICENSE) 了解詳情。

---

## 👥 開發者

**專案維護者**: Jacky Chang

---

## 📚 References

- [Go Programming Language](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [gRPC](https://grpc.io/)
- [Protocol Buffers](https://protobuf.dev/)

完整清單請參考 [go.mod](go.mod)

---

**有問題或建議？** 歡迎提出 [Issue](https://github.com/ChangChunCheng/webprism/issues) 或 [Pull Request](https://github.com/ChangChunCheng/webprism/pulls)
