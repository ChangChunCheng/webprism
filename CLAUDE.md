# CLAUDE.md - WEBPRISM 開發歷程與需求文件
# CLAUDE.md - WEBPRISM Development History & Requirements Document

**專案**: WEBPRISM (Web Proxy Request Integration Service Manager)
**建立日期**: 2026-01-07
**維護者**: Jacky Chang
**目的**: 記錄所有過往及未來的開發要求，供 Claude AI 參考

---

## 📚 文件說明

本文件記錄 WEBPRISM 專案從初始到現在的所有開發歷程、設計決策、已完成功能及未來規劃。供 Claude AI 助手在協助開發時參考，確保開發方向一致且符合專案架構原則。

---

## 🎯 專案核心理念

### 核心目標
打造一個智能 API 整合中介層，讓開發者只需上傳 OpenAPI 規格文件，即可立即使用任何第三方 API，無需處理複雜的認證和請求細節。

### 設計哲學
1. **文件驅動** - 上傳 API 規格即可使用
2. **統一認證** - 集中管理所有 API 憑證
3. **安全代理** - 安全轉發請求，自動注入認證
4. **多介面支援** - REST API、gRPC、CLI、MCP Server
5. **企業級品質** - 完整的日誌、監控、測試

---

## 🏗️ 架構設計原則

### 六角形架構 (Hexagonal Architecture)
```
┌─────────────────────────────────────────┐
│         Adapters Layer                  │
│  (HTTP, gRPC, CLI, MCP)                │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │      Ports Layer                  │ │
│  │  (Input/Output Interfaces)        │ │
│  │                                   │ │
│  │  ┌─────────────────────────────┐ │ │
│  │  │    Domain Layer             │ │ │
│  │  │  (Business Logic)           │ │ │
│  │  └─────────────────────────────┘ │ │
│  └───────────────────────────────────┘ │
│                                         │
│      Infrastructure Layer               │
│  (Database, Config, Logger, Crypto)    │
└─────────────────────────────────────────┘
```

### 依賴規則
- **依賴方向**: Adapters → Ports → Domain
- **核心原則**: Domain Layer 不依賴任何外層
- **介面導向**: 所有跨層通訊透過介面

### 技術棧
- **語言**: Go 1.21+
- **資料庫**: PostgreSQL 16+
- **協定**: REST + gRPC
- **容器**: Docker + Docker Compose
- **API 文件**: OpenAPI 3.0
- **測試**: testify, mockery

---

## ✅ 已完成功能 (Completed Features)

### Phase 1: 基礎架構建立 (2025-12-30)

#### 1.1 Protocol Buffers 定義
- ✅ 定義 4 個核心服務的 gRPC 介面
  - `SpecService` - API 規格管理
  - `AuthService` - 認證配置管理
  - `ProxyService` - 請求代理執行
  - `HealthService` - 健康檢查
- ✅ 支援 gRPC-Gateway HTTP 轉換
- ✅ 完整的 protobuf 文件註解

**檔案位置**: `api/proto/v1/*.proto`

#### 1.2 Domain Layer 實作
- ✅ 定義核心 Domain Models
  - `APISpec` - API 規格實體
  - `AuthConfig` - 認證配置實體
  - `ProxyRequest/Response` - 代理請求/回應
  - `HealthStatus` - 健康狀態
- ✅ 實作 4 個 Domain Services
  - 完整的業務邏輯封裝
  - 錯誤處理機制
  - 參數驗證

**檔案位置**: `internal/domain/model/*.go`, `internal/domain/service/*.go`

#### 1.3 Ports Layer 定義
- ✅ 定義 Input Ports (Use Cases)
  - 清晰的業務用例介面
- ✅ 定義 Output Ports (Repositories)
  - 資料存取抽象介面
  - HTTP Client 介面
  - Crypto Service 介面
  - OpenAPI Parser 介面

**檔案位置**: `internal/ports/input/*.go`, `internal/ports/output/*.go`

#### 1.4 Adapters Layer 實作
- ✅ gRPC Handlers (4 個服務)
- ✅ HTTP Server (gRPC-Gateway)
- ✅ MCP Server 整合
- ✅ CLI Commands
- ✅ PostgreSQL Repositories
- ✅ HTTP Client Adapter
- ✅ OpenAPI Parser Adapter

**檔案位置**: `internal/adapters/*`

#### 1.5 Infrastructure Layer
- ✅ 配置管理 (Viper)
  - 支援環境變數
  - YAML 配置檔
  - 配置優先順序管理
- ✅ 日誌系統 (slog)
  - 結構化日誌
  - 多級別日誌
- ✅ 加密服務 (AES-256-GCM)
  - 敏感憑證加密
  - 安全金鑰管理
- ✅ 資料庫連接 (pgx)
  - 連接池管理
  - 健康檢查

**檔案位置**: `internal/infrastructure/*`

### Phase 2: 測試完善 (2026-01-01)

#### 2.1 單元測試實作
- ✅ 80 個測試案例，全部通過
- ✅ 95.2% 測試覆蓋率 (超出目標 85%)
- ✅ 0.273 秒執行時間
- ✅ 完整的 Mock 隔離
- ✅ Table-Driven Tests 模式
- ✅ AAA (Arrange-Act-Assert) 結構

**測試覆蓋率詳情**:
- ProxyService: 93.1%
- SpecService: 100.0%
- AuthService: 95.2%
- HealthService: 93.8%

**檔案位置**: `internal/domain/service/*_test.go`

#### 2.2 測試基礎設施
- ✅ Mock 生成配置 (.mockery.yaml)
- ✅ 7 個自動生成的 Mock interfaces
- ✅ 測試工具包 (testutil)
  - 動態測試資料生成
  - 無硬編碼測試值
  - 可重用測試輔助函數

**檔案位置**: `internal/mocks/*`, `internal/testutil/helpers.go`

### Phase 3: 文件體系建立 (2025-12-31 - 2026-01-03)

#### 3.1 使用者文件
- ✅ [README.md](README.md) - 專案概述
- ✅ [GETTING_STARTED.md](docs/GETTING_STARTED.md) - 快速開始
- ✅ [Installation Guide](docs/user-guide/01-INSTALLATION.md) - 安裝指南
- ✅ [Quickstart](docs/user-guide/02-QUICKSTART.md) - 使用教學

#### 3.2 開發者文件
- ✅ [CONTRIBUTING.md](docs/development/CONTRIBUTING.md) - 開發指南
- ✅ [CURRENT_STATUS.md](docs/development/CURRENT_STATUS.md) - 專案狀態
- ✅ [Testing Guide](docs/testing/QUICK_START.md) - 測試指南
- ✅ [Test Reports](docs/testing/) - 測試報告

#### 3.3 架構文件
- ✅ [VERSION_MANAGEMENT.md](docs/VERSION_MANAGEMENT.md) - 版本管理
- ✅ [DOCUMENTATION_READY.md](docs/DOCUMENTATION_READY.md) - 文件索引

### Phase 4: MCP 介面測試環境 (2026-01-03)

#### 4.1 MCP 測試環境建立
- ✅ React + TypeScript 測試介面
- ✅ MCP SDK 整合
- ✅ 互動式測試工具
- ✅ 完整的測試腳本

**檔案位置**: `tests/mcp/*`

**相關 Commit**: `799cf17 - feat: add comprehensive MCP interface testing environment (#1)`

### Phase 5: 版本管理系統 (2025-12-31)

#### 5.1 版本管理機制
- ✅ 自動版本號生成
- ✅ Git 版本資訊整合
- ✅ 建置資訊追蹤
- ✅ CLI 版本命令

**檔案位置**: `internal/version/version.go`, `scripts/version.sh`

**相關 Commit**: `498a79f - feat: implement comprehensive version management system`

### Phase 6: 開發工具與腳本

#### 6.1 Makefile 指令
- ✅ `make setup` - 完整專案設置
- ✅ `make build` - 建置所有二進位檔
- ✅ `make test` - 執行測試
- ✅ `make proto` - 生成 Protobuf 程式碼
- ✅ `make mocks` - 生成 Mock
- ✅ `make lint` - 程式碼檢查
- ✅ `make docker-up-full` - Docker 完整啟動
- ✅ `make run-server` - 執行服務器
- ✅ `make run-mcp` - 執行 MCP Server
- ✅ `make run-cli` - 執行 CLI

#### 6.2 開發腳本
- ✅ `scripts/check-tools.sh` - 工具檢查
- ✅ `scripts/cleanup.sh` - 清理腳本
- ✅ `scripts/version.sh` - 版本管理
- ✅ `scripts/detect-platform.sh` - 平台偵測

---

## 🔧 技術決策記錄 (Technical Decisions)

### TD-001: 選擇 Go 作為開發語言
**日期**: 2025-12-30
**決策**: 使用 Go 1.21+ 作為主要開發語言
**理由**:
- 高效能，適合 API 代理
- 優秀的並發支援
- 豐富的生態系統
- 容易部署（單一二進位檔）
- 強大的 gRPC 支援

### TD-002: 採用六角形架構
**日期**: 2025-12-30
**決策**: 使用六角形架構 (Hexagonal Architecture)
**理由**:
- 業務邏輯與技術細節分離
- 易於測試（可替換外部依賴）
- 支援多種介面（HTTP、gRPC、CLI、MCP）
- 長期可維護性高
- 符合 SOLID 原則

### TD-003: 使用 PostgreSQL 作為資料庫
**日期**: 2025-12-30
**決策**: 選擇 PostgreSQL 16+ 作為持久化儲存
**理由**:
- 成熟穩定的關聯式資料庫
- 優秀的 JSONB 支援（存儲 OpenAPI 規格）
- 強大的索引能力
- 良好的 Go 驅動支援 (pgx)
- 適合企業環境

### TD-004: gRPC + gRPC-Gateway 雙協定
**日期**: 2025-12-30
**決策**: 同時支援 gRPC 和 REST API
**理由**:
- gRPC 提供高效能
- REST API 提供廣泛相容性
- gRPC-Gateway 自動轉換
- 維護單一 Protobuf 定義即可

### TD-005: AES-256-GCM 加密敏感憑證
**日期**: 2025-12-30
**決策**: 使用 AES-256-GCM 加密 API 憑證
**理由**:
- 業界標準加密演算法
- 提供認證加密 (AEAD)
- 防止憑證洩漏
- Go 原生支援

### TD-006: 配置分離策略
**日期**: 2025-12-30
**決策**: 非機密配置用 YAML，機密資訊用環境變數
**理由**:
- 符合 12-Factor App 原則
- 安全管理機密資訊
- 環境變數優先級最高
- 適合容器化部署

### TD-007: Table-Driven Tests 測試模式
**日期**: 2026-01-01
**決策**: 使用 Table-Driven Tests 作為標準測試模式
**理由**:
- Go 社群最佳實踐
- 易於添加測試案例
- 減少重複程式碼
- 提高測試可讀性

### TD-008: 動態測試資料生成
**日期**: 2026-01-01
**決策**: 建立 testutil 包管理測試資料，禁止硬編碼
**理由**:
- 提高測試可維護性
- 避免測試資料污染
- 集中管理測試邏輯
- 可重用測試工具

---

## 🚧 已知問題與解決方案

### Issue-001: Base64 編碼 Basic Auth (已解決)
**發現日期**: 2026-01-01
**狀態**: ✅ 已解決
**問題**: Basic Auth 未正確 Base64 編碼
**解決方案**: 在 proxy_service.go 中正確實作 Base64 編碼
**相關測試**: TestProxyService_ExecuteProxy_BasicAuth_Base64Encoding

### Issue-002: 硬編碼測試資料 (已解決)
**發現日期**: 2026-01-01
**狀態**: ✅ 已解決
**問題**: 測試中出現硬編碼值
**解決方案**: 建立 testutil.helpers.go 提供動態資料生成
**影響檔案**: 所有 *_test.go 檔案

### Issue-003: Logger 建構返回值處理 (已解決)
**發現日期**: 2026-01-01
**狀態**: ✅ 已解決
**問題**: logger.New() 返回兩個值未正確處理
**解決方案**: 正確接收並處理兩個返回值

### Issue-004: Mock 設定不完整 (已解決)
**發現日期**: 2026-01-01
**狀態**: ✅ 已解決
**問題**: 某些測試的 Mock expectations 設定不完整
**解決方案**: 為所有依賴方法設定完整的 Mock

---

## 🎯 未來開發規劃

### Phase 7: 整合測試 (優先級: 高)
**預估時間**: 1-2 天
**目標覆蓋率**: 增加 10-15%

#### 任務清單
- [ ] Repository 層與 PostgreSQL 整合測試
- [ ] HTTP Handler 整合測試
- [ ] gRPC Handler 整合測試
- [ ] 完整服務棧測試
- [ ] 資料庫 Migration 測試
- [ ] 測試資料庫設置與清理自動化

#### 預期產出
- 15-20 個整合測試案例
- 整合測試文件
- CI/CD 整合測試配置

### Phase 8: CI/CD 流水線 (優先級: 高)
**預估時間**: 0.5-1 天

#### 任務清單
- [ ] GitHub Actions workflow 配置
- [ ] 自動化單元測試執行
- [ ] 自動化整合測試執行
- [ ] 測試覆蓋率報告生成
- [ ] 自動化建置與發布
- [ ] Docker 映像自動建置
- [ ] 測試失敗通知機制

#### 預期產出
- `.github/workflows/test.yml`
- `.github/workflows/build.yml`
- `.github/workflows/release.yml`
- Coverage badge 整合

### Phase 9: E2E 測試 (優先級: 中)
**預估時間**: 1-2 天

#### 任務清單
- [ ] 完整 API 工作流程測試
  - 上傳 OpenAPI 規格
  - 設定認證
  - 執行代理請求
  - 查詢健康狀態
- [ ] 真實第三方 API 整合測試
- [ ] MCP Server E2E 測試
- [ ] CLI E2E 測試
- [ ] 錯誤情境測試

#### 預期產出
- 8-10 個 E2E 測試場景
- E2E 測試文件
- 測試資料準備腳本

### Phase 10: 效能優化與測試 (優先級: 中低)
**預估時間**: 1-2 天

#### 任務清單
- [ ] 效能基準測試建立
- [ ] 負載測試 (1000+ RPS)
- [ ] 並發測試
- [ ] 資料庫查詢優化
- [ ] HTTP Client 連接池優化
- [ ] 記憶體使用分析
- [ ] CPU 使用分析

#### 預期產出
- 效能測試報告
- 優化建議文件
- 效能監控儀表板

### Phase 11: 功能擴展 (優先級: 中)

#### 11.1 OAuth2 完整支援
- [ ] OAuth2 Authorization Code Flow
- [ ] OAuth2 Client Credentials Flow
- [ ] Token 自動更新機制
- [ ] OAuth2 配置 UI

#### 11.2 API 速率限制
- [ ] 基於 Token Bucket 的速率限制
- [ ] Per-API 限制配置
- [ ] Per-User 限制配置
- [ ] 速率限制監控

#### 11.3 請求/回應轉換
- [ ] 請求參數映射
- [ ] 回應資料轉換
- [ ] 自定義轉換規則
- [ ] JSONPath 支援

#### 11.4 快取機制
- [ ] Redis 快取整合
- [ ] 可配置快取策略
- [ ] 快取失效機制
- [ ] 快取監控

### Phase 12: 可觀測性提升 (優先級: 中)

#### 12.1 日誌增強
- [ ] 分散式追蹤 (OpenTelemetry)
- [ ] 日誌聚合整合 (ELK/Loki)
- [ ] 結構化日誌標準化
- [ ] 日誌保留策略

#### 12.2 監控與告警
- [ ] Prometheus metrics 導出
- [ ] Grafana 儀表板
- [ ] 關鍵指標告警
- [ ] SLA 監控

#### 12.3 審計日誌
- [ ] API 呼叫審計
- [ ] 認證配置變更審計
- [ ] 管理操作審計
- [ ] 審計日誌查詢介面

### Phase 13: 安全性增強 (優先級: 高)

#### 13.1 認證與授權
- [ ] API Key 認證
- [ ] JWT Token 認證
- [ ] RBAC (角色基礎存取控制)
- [ ] 多租戶隔離

#### 13.2 安全掃描
- [ ] 依賴漏洞掃描
- [ ] 程式碼安全掃描
- [ ] Docker 映像掃描
- [ ] 滲透測試

#### 13.3 資料安全
- [ ] 敏感資料遮罩
- [ ] 加密金鑰輪換
- [ ] 資料備份與恢復
- [ ] GDPR 合規支援

### Phase 14: 部署與運維 (優先級: 高)

#### 14.1 容器化部署
- [ ] 多階段 Docker 建置優化
- [ ] Kubernetes Deployment
- [ ] Helm Charts
- [ ] Health Check 與 Liveness Probe

#### 14.2 高可用性
- [ ] 水平擴展支援
- [ ] 資料庫讀寫分離
- [ ] 故障轉移機制
- [ ] 負載均衡配置

#### 14.3 災難恢復
- [ ] 自動備份策略
- [ ] 備份驗證機制
- [ ] 災難恢復演練
- [ ] 恢復時間目標 (RTO) 優化

### Phase 15: 開發者體驗 (優先級: 中)

#### 15.1 API 文件
- [ ] OpenAPI 規格自動生成
- [ ] Swagger UI 整合
- [ ] API 使用範例
- [ ] SDK 生成（多語言）

#### 15.2 開發工具
- [ ] VSCode 擴充套件
- [ ] Postman Collection
- [ ] CLI 自動補全
- [ ] 開發環境快速啟動腳本

---

## 📋 開發規範與最佳實踐

### 程式碼風格
1. **遵循 Effective Go** - 標準 Go 程式碼風格
2. **使用 golangci-lint** - 自動化程式碼檢查
3. **清晰的命名** - 避免縮寫，使用描述性名稱
4. **適當的註解** - Package、Public 函數必須有註解
5. **錯誤處理** - 永遠處理錯誤，使用 fmt.Errorf 包裝

### 測試規範
1. **測試覆蓋率 ≥ 85%** - 當前已達 95.2%
2. **Table-Driven Tests** - 所有測試使用此模式
3. **AAA 結構** - Arrange, Act, Assert
4. **Mock 隔離** - 所有外部依賴必須 Mock
5. **無硬編碼** - 使用 testutil 動態生成測試資料
6. **清晰命名** - 測試名稱描述測試場景

### Git 工作流程
1. **Conventional Commits** - 使用標準化 commit 訊息
   - `feat:` - 新功能
   - `fix:` - 錯誤修復
   - `docs:` - 文件更新
   - `test:` - 測試相關
   - `refactor:` - 重構
   - `chore:` - 建置/工具相關
2. **分支策略**
   - `main` - 穩定版本
   - `dev` - 開發分支
   - `feature/*` - 功能分支
   - `fix/*` - 修復分支
3. **Pull Request** - 所有變更透過 PR 審查

### 文件規範
1. **README First** - 先寫文件，再寫程式碼
2. **保持更新** - 程式碼變更時同步更新文件
3. **範例優先** - 提供實際可執行的範例
4. **雙語支援** - 中英文並列（重要文件）

---

## 🔍 程式碼審查檢查清單

### 功能性
- [ ] 功能符合需求規格
- [ ] 錯誤處理完整
- [ ] 邊界條件考慮
- [ ] 並發安全性

### 程式碼品質
- [ ] 遵循專案架構原則
- [ ] 依賴方向正確（不違反六角形架構）
- [ ] 介面設計合理
- [ ] 命名清晰易懂
- [ ] 無重複程式碼
- [ ] 註解適當且準確

### 測試
- [ ] 單元測試覆蓋率 ≥ 85%
- [ ] 測試案例完整
- [ ] 使用 Table-Driven Tests
- [ ] Mock 隔離正確
- [ ] 無硬編碼測試資料

### 安全性
- [ ] 無敏感資訊洩漏
- [ ] 輸入驗證完整
- [ ] SQL Injection 防護
- [ ] XSS 防護
- [ ] CSRF 防護（如適用）

### 效能
- [ ] 無明顯效能問題
- [ ] 資料庫查詢優化
- [ ] 記憶體使用合理
- [ ] 無資源洩漏

### 文件
- [ ] Public API 有註解
- [ ] 複雜邏輯有說明
- [ ] README 已更新
- [ ] CHANGELOG 已更新（如適用）

---

## 📊 專案指標與目標

### 程式碼指標
- **總程式碼行數**: ~7,300 行 (不含生成程式碼)
- **測試覆蓋率**: 95.2% (目標 ≥ 85%)
- **Go 檔案數**: ~60 個
- **測試案例數**: 80 個
- **Mock 介面數**: 7 個

### 品質指標
- **單元測試通過率**: 100%
- **測試執行時間**: 0.273 秒
- **Lint 警告數**: 0
- **已知安全漏洞**: 0

### 文件指標
- **文件檔案數**: 20+
- **文件覆蓋率**: 100% (所有 Public API)
- **範例程式碼**: 完整

---

## 🎓 學習資源與參考

### Go 語言
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)

### 架構設計
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture (Robert C. Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

### 測試
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [testify Documentation](https://github.com/stretchr/testify)
- [mockery Documentation](https://vektra.github.io/mockery/)

### gRPC
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [gRPC-Gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Protocol Buffers](https://protobuf.dev/)

---

## 🤝 協作指南

### 與 Claude AI 協作
1. **提供足夠上下文** - 說明需求的背景和目的
2. **明確指定範圍** - 清楚說明要修改的範圍
3. **參考本文件** - 確保符合專案架構和規範
4. **測試驅動** - 要求同時撰寫測試
5. **文件同步** - 要求更新相關文件

### 範例請求格式
```
任務: [簡短描述]
目的: [為什麼需要這個功能]
範圍: [涉及的檔案/模組]
需求:
  - [具體需求 1]
  - [具體需求 2]
限制:
  - 遵循六角形架構
  - 測試覆蓋率 ≥ 85%
  - 使用 Table-Driven Tests
  - 更新相關文件
```

---

## 📞 聯絡與支援

### 專案資訊
- **Repository**: https://github.com/ChangChunCheng/webprism
- **Issue Tracker**: https://github.com/ChangChunCheng/webprism/issues
- **維護者**: Jacky Chang

### 取得協助
1. 查看 [文件目錄](docs/README.md)
2. 搜尋 [已知問題](https://github.com/ChangChunCheng/webprism/issues)
3. 提出新的 Issue
4. 參與 Discussions

---

## 📝 變更日誌

### 2026-01-07
- ✅ 建立 CLAUDE.md 文件
- 📝 記錄完整開發歷程
- 📋 規劃未來開發路線圖

### 2026-01-03
- ✅ 完成 MCP 測試環境建立
- ✅ 整理文件體系

### 2026-01-01
- ✅ 完成單元測試 (95.2% 覆蓋率)
- ✅ 建立測試基礎設施
- 📝 產出完整測試報告

### 2025-12-31
- ✅ 實作版本管理系統
- ✅ 完成核心功能開發

### 2025-12-30
- 🎉 專案初始化
- ✅ 建立六角形架構
- ✅ 實作 4 個核心服務

---

## 🎯 總結

WEBPRISM 專案已建立堅實的技術基礎，包含：
- ✅ 清晰的六角形架構
- ✅ 完整的單元測試（95.2% 覆蓋率）
- ✅ 豐富的文件體系
- ✅ 自動化建置與測試工具
- ✅ 4 個核心服務完整實作

下一階段重點：
1. 整合測試與 E2E 測試
2. CI/CD 流水線建立
3. 效能優化與監控
4. 安全性增強
5. 功能擴展

**目標**: 打造一個穩定、高效、易用的 API 整合中介層服務。

---

**最後更新**: 2026-01-07
**文件版本**: v1.0
**維護者**: Jacky Chang

---

*本文件將持續更新，記錄所有重要的開發決策與變更*
