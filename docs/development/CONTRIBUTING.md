# 🛠️ WEBPRISM 開發指南

本指南提供給想要參與 WEBPRISM 開發的貢獻者。

---

## 📋 目錄

- [開發環境設置](#開發環境設置)
- [專案結構](#專案結構)
- [開發工作流程](#開發工作流程)
- [程式碼規範](#程式碼規範)
- [測試要求](#測試要求)
- [提交規範](#提交規範)
- [Pull Request 流程](#pull-request-流程)

---

## 開發環境設置

### 前置需求

| 工具 | 版本要求 | 用途 |
|------|----------|------|
| **Go** | 1.21+ | 開發語言 |
| **protoc** | 3.x+ | Protobuf 編譯器 |
| **Docker** | 20.x+ | 容器化開發 |
| **Make** | 3.x+ | 建置自動化 |
| **Git** | 2.x+ | 版本控制 |

### 完整設置步驟

```bash
# 1. Clone 專案
git clone https://github.com/ChangChunCheng/webprism.git
cd webprism

# 2. 檢查必要工具
make check-tools

# 3. 初始化專案
make setup

# 4. 配置環境變數
cp .env.example .env
# 編輯 .env 填入開發環境設定

# 5. 啟動開發環境
make db-up          # 啟動 PostgreSQL
make run-server     # 執行 WEBPRISM

# 6. 執行測試驗證
make test
```

---

## 專案結構

```
webprism/
├── cmd/                        # 應用程式入口
│   ├── server/                 # HTTP/gRPC 伺服器
│   ├── mcp/                    # MCP 伺服器
│   └── cli/                    # CLI 工具
├── internal/                   # 內部程式碼（不對外開放）
│   ├── domain/                 # 領域層（核心業務邏輯）
│   │   ├── entity/             # 實體定義
│   │   ├── repository/         # Repository 介面
│   │   └── service/            # 業務邏輯服務
│   ├── ports/                  # 介面定義（六角形架構的 Ports）
│   │   ├── input/              # 輸入埠（Use Cases）
│   │   └── output/             # 輸出埠（Repositories）
│   ├── adapters/               # 適配器層
│   │   ├── api/                # HTTP/gRPC Handler
│   │   ├── cli/                # CLI 命令
│   │   └── mcp/                # MCP 介面
│   └── infrastructure/         # 基礎設施層
│       ├── database/           # 資料庫實作
│       ├── config/             # 配置管理
│       └── logger/             # 日誌系統
├── api/                        # API 定義
│   └── proto/                  # Protobuf 定義
├── gen/                        # 生成的程式碼（不提交）
├── deployments/                # 部署配置
│   ├── config.yaml             # 非機密配置
│   └── docker/                 # Docker 相關
├── docs/                       # 📚 完整文件
├── scripts/                    # 工具腳本
├── testdata/                   # 測試資料
└── third_party/                # 第三方定義
```

### 架構設計原則

#### 1. 六角形架構 (Hexagonal Architecture)

- **核心層 (Domain)**: 純業務邏輯，不依賴任何外部框架
- **介面層 (Ports)**: 定義輸入/輸出介面
- **適配器層 (Adapters)**: 實作各種傳輸協定
- **基礎設施層 (Infrastructure)**: 技術實作細節

#### 2. 依賴方向

```
Adapters → Ports → Domain
Infrastructure → Domain
```

**規則**: 依賴只能向內（Domain 不依賴任何外層）

#### 3. 介面導向設計

所有跨層通訊都透過介面，便於測試和替換實作。

---

## 開發工作流程

### 1. 建立開發分支

```bash
# 從 main 建立功能分支
git checkout main
git pull origin main
git checkout -b feature/your-feature-name

# 從 main 建立修復分支
git checkout -b fix/issue-number-description
```

### 2. 開發循環

```bash
# 1. 修改程式碼
# 2. 重新生成程式碼（如果修改了 proto 或介面）
make proto
make mocks

# 3. 執行測試
make test

# 4. 檢查程式碼品質
make lint

# 5. 建置驗證
make build
```

### 3. 提交變更

```bash
# 確保測試通過
make test

# 提交變更
git add .
git commit -m "feat: add new feature"

# 推送到遠端
git push origin feature/your-feature-name
```

---

## 程式碼規範

### Go 程式碼風格

遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

#### 命名規範

```go
// 好的命名
type SpecRepository interface {
    Save(spec *Spec) error
    FindByID(id string) (*Spec, error)
}

type specRepositoryImpl struct {
    db *sql.DB
}

// 避免的命名
type SpecRepo interface {}  // 不要縮寫介面名稱
type spec_repository struct {}  // 不要使用蛇形命名
```

#### 錯誤處理

```go
// 好的錯誤處理
result, err := service.DoSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// 避免忽略錯誤
result, _ := service.DoSomething()  // ❌
```

#### 註解規範

```go
// Package domain 提供核心業務邏輯
package domain

// Spec 代表一個 API 規格
type Spec struct {
    ID   string
    Name string
}

// Save 儲存規格到資料庫
// 如果規格已存在，將返回 ErrSpecAlreadyExists 錯誤
func (s *SpecService) Save(spec *Spec) error {
    // 實作...
}
```

### 程式碼檢查

```bash
# 執行 lint 檢查
make lint

# 自動格式化程式碼
go fmt ./...

# 整理 import
goimports -w .
```

---

## 測試要求

### 測試覆蓋率要求

- **最低要求**: 80%
- **目標**: 95%+
- **核心業務邏輯**: 100%

### 測試類型

#### 1. 單元測試

測試單一函數或方法的行為。

```go
func TestSpecService_Save(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewMockSpecRepository(t)
    service := NewSpecService(mockRepo)

    spec := &Spec{ID: "1", Name: "Test"}
    mockRepo.EXPECT().Save(spec).Return(nil)

    // Act
    err := service.Save(spec)

    // Assert
    assert.NoError(t, err)
}
```

#### 2. 整合測試

測試多個組件的整合。

```go
func TestSpecAPI_CreateSpec_Integration(t *testing.T) {
    // 設置測試資料庫
    db := setupTestDB(t)
    defer db.Close()

    // 執行整合測試
    // ...
}
```

### 執行測試

```bash
# 執行所有測試
make test

# 執行特定套件測試
go test ./internal/domain/service/

# 查看詳細輸出
go test -v ./internal/domain/service/

# 生成覆蓋率報告
make test-coverage

# 查看覆蓋率 HTML 報告
go tool cover -html=coverage.out
```

### Mock 生成

```bash
# 生成所有 Mock
make mocks

# 生成特定介面的 Mock
mockery --name=SpecRepository --dir=internal/domain/repository --output=internal/mocks
```

---

## 提交規範

### Commit Message 格式

使用 [Conventional Commits](https://www.conventionalcommits.org/) 規範：

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Type 類型

- `feat`: 新功能
- `fix`: 修復 Bug
- `docs`: 文件更新
- `style`: 程式碼格式調整（不影響功能）
- `refactor`: 重構（不影響功能）
- `test`: 測試相關
- `chore`: 建置或輔助工具相關

#### 範例

```bash
# 新功能
git commit -m "feat(spec): add OpenAPI 3.1 support"

# 修復 Bug
git commit -m "fix(auth): resolve bearer token validation issue"

# 文件更新
git commit -m "docs: update installation guide"

# 重構
git commit -m "refactor(service): simplify spec validation logic"
```

---

## Pull Request 流程

### 1. 建立 PR 前檢查清單

- [ ] 所有測試通過 (`make test`)
- [ ] 程式碼檢查通過 (`make lint`)
- [ ] 測試覆蓋率符合要求
- [ ] 更新相關文件
- [ ] Commit message 符合規範

### 2. PR 標題格式

```
<type>: <description>
```

範例：
```
feat: Add OAuth2 authentication support
fix: Resolve database connection leak
docs: Update API documentation
```

### 3. PR 描述範本

```markdown
## 變更說明
簡述這個 PR 的目的和變更內容

## 變更類型
- [ ] 新功能
- [ ] Bug 修復
- [ ] 重構
- [ ] 文件更新
- [ ] 測試更新

## 測試
說明如何測試這些變更

## 相關 Issue
Closes #123

## Checklist
- [ ] 測試通過
- [ ] 文件已更新
- [ ] Commit message 符合規範
```

### 4. Code Review 流程

1. 提交 PR 後，等待 CI/CD 檢查通過
2. 至少需要 1 位維護者審核
3. 解決所有審核意見
4. 獲得批准後合併到 main

---

## Makefile 快速參考

### 初始設定

```bash
make check-tools        # 檢查開發工具
make setup              # 完整設定專案
```

### 開發

```bash
make db-up              # 啟動 PostgreSQL
make run-server         # 執行 HTTP Server
make run-mcp            # 執行 MCP Server
make run-cli            # 執行 CLI
```

### 建置

```bash
make build              # 建置所有二進位檔
make proto              # 生成 Protobuf 程式碼
make mocks              # 生成測試 Mock
```

### 測試

```bash
make test               # 執行所有測試
make test-coverage      # 生成覆蓋率報告
make lint               # 執行程式碼檢查
```

### Docker

```bash
make docker-up-full     # 啟動所有服務
make docker-down        # 停止所有服務
make docker-logs        # 查看容器日誌
```

### 清理

```bash
make clean              # 清理建置檔案
make clean-all          # 清理所有生成檔案
```

完整指令說明執行 `make help`

---

## 配置管理

### 配置檔案分離

| 檔案 | 用途 | 是否提交 Git |
|------|------|--------------|
| `deployments/config.yaml` | 非機密配置 | ✅ 提交 |
| `.env` | 機密資訊 | ❌ 不提交 |

### 配置優先順序

```
環境變數 > .env > config.yaml > 程式預設值
```

### 範例配置

**deployments/config.yaml**（非機密）:
```yaml
server:
  port: 8080
  timeout: 30s

database:
  max_connections: 100
  connection_timeout: 10s
```

**.env**（機密）:
```bash
DB_PASSWORD=secret_password
ENCRYPTION_KEY=your-32-byte-key
API_TOKEN=secret_token
```

詳細說明請參考 [deployments/README.md](../../deployments/README.md)

---

## 常見開發任務

### 新增 API Endpoint

1. 在 `api/proto/` 定義 Protobuf
2. 執行 `make proto` 生成程式碼
3. 在 `internal/ports/input/` 定義 Use Case 介面
4. 在 `internal/domain/service/` 實作業務邏輯
5. 在 `internal/adapters/api/` 實作 Handler
6. 撰寫測試

### 新增 CLI 命令

1. 在 `internal/adapters/cli/` 建立命令檔案
2. 實作命令邏輯
3. 註冊到 `cmd/cli/main.go`
4. 撰寫測試
5. 更新文件

### 新增 Repository

1. 在 `internal/domain/repository/` 定義介面
2. 在 `internal/infrastructure/database/` 實作
3. 生成 Mock: `make mocks`
4. 撰寫測試

---

## 故障排除

### 常見問題

#### Q: Import 錯誤 - 找不到 gen/ 目錄

**解決方式**:
```bash
make proto
make mocks
```

#### Q: Mock 生成失敗

**解決方式**:
```bash
# 重新安裝 mockery
make install-tools

# 清理後重新生成
make clean-all
make setup
```

#### Q: 測試資料庫連線失敗

**解決方式**:
```bash
# 確保測試資料庫已啟動
make db-up

# 檢查 .env 設定
cat .env
```

---

## 取得協助

- **文件**: 查看 [docs/](../README.md)
- **Issue**: [GitHub Issues](https://github.com/ChangChunCheng/webprism/issues)
- **討論**: [GitHub Discussions](https://github.com/ChangChunCheng/webprism/discussions)

---

**感謝你的貢獻！** 🎉
