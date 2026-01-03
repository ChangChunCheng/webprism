# WEBPRISM 版本管理指南

本文件說明 WEBPRISM 的版本管理系統、構建流程和 Release 發布流程。

## 📋 目錄

- [版本類型概述](#版本類型概述)
- [開發版本構建](#開發版本構建)
- [測試版本構建](#測試版本構建)
- [正式版本構建](#正式版本構建)
- [版本查詢方式](#版本查詢方式)
- [常見問題](#常見問題)

---

## 版本類型概述

WEBPRISM 使用三種版本類型來區分不同階段的構建：

| 版本類型 | 格式 | 標記 | 使用場景 |
|---------|------|------|---------|
| **開發版本** | `branch+commit` | `[DEVELOPMENT]` | 日常開發、功能分支 |
| **測試版本** | `alpha/beta+tag+commit` | `[ALPHA/BETA]` | 內部測試、QA 測試 |
| **正式版本** | `release+tag+commit` | `[RELEASE]` | 生產環境、正式發布 |

### 版本資訊內容

每個構建的二進位檔案都包含以下資訊：

- **版本號**：根據類型自動格式化
- **Git Commit**：完整的 commit SHA hash
- **Git Branch**：分支名稱
- **Build Time**：構建時間（UTC RFC3339 格式）
- **Go Version**：Go 編譯器版本
- **Platform**：目標平台（OS/架構）

---

## 開發版本構建

開發版本用於日常開發和功能分支測試，**不需要創建 Git tag**。

### 使用場景

- 功能開發
- Bug 修復
- 本地測試
- 功能分支構建

### 構建步驟

#### 1. 在功能分支工作

```bash
# 創建功能分支
git checkout -b feature/new-feature

# 或使用現有分支
git checkout feature/existing-feature
```

#### 2. 開發和提交代碼

```bash
# 正常開發...
git add .
git commit -m "Add new feature"
```

#### 3. 構建開發版本

```bash
make build
```

#### 4. 查看版本

```bash
./bin/webprism version --short
```

**輸出範例：**
```
WEBPRISM feature/new-feature+abc1234 (commit: abc1234)
```

**完整版本資訊：**
```bash
./bin/webprism version
```

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  WEBPRISM CLI [DEVELOPMENT]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Version:     feature/new-feature+abc1234
  Git Commit:  abc1234567890abcdef1234567890abcdef123456
  Git Branch:  feature/new-feature
  Build Time:  2026-01-03T12:00:00Z
  Go Version:  go1.24.11
  Platform:    darwin/arm64
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 版本格式說明

- **格式**：`branch+commit`
- **範例**：`build-foundation+a24b17f`
- **說明**：自動使用當前分支名稱加上短 commit hash

---

## 測試版本構建

測試版本用於內部測試、QA 測試和預發布階段，**需要創建帶有 `-alpha`、`-beta` 或 `-rc` 後綴的 Git tag**。

### 使用場景

- Alpha 測試（內部測試）
- Beta 測試（擴大測試）
- Release Candidate（發布候選）
- QA 環境部署

### Tag 命名規範

```
v主版本.次版本.修訂版本-階段.迭代號

範例：
  v1.0.0-alpha.1    # 第一個 alpha 版本
  v1.0.0-alpha.2    # 第二個 alpha 版本
  v1.0.0-beta.1     # 第一個 beta 版本
  v1.0.0-beta.2     # 第二個 beta 版本
  v1.0.0-rc.1       # Release Candidate 1
  v1.0.0-rc.2       # Release Candidate 2
```

### 構建步驟

#### 1. 確保代碼已提交

```bash
# 查看狀態
git status

# 提交所有變更
git add .
git commit -m "Prepare alpha release v1.0.0-alpha.1"
```

#### 2. 創建測試版本 Tag

```bash
# 使用 make 命令（推薦）
make release-tag VERSION=v1.0.0-alpha.1

# 或手動創建
git tag -a v1.0.0-alpha.1 -m "Alpha release 1.0.0-alpha.1"
```

#### 3. 構建測試版本

```bash
make build
```

#### 4. 驗證版本資訊

```bash
./bin/webprism version --short
```

**輸出範例：**
```
WEBPRISM alpha/beta+v1.0.0-alpha.1+abc1234 (commit: abc1234)
```

#### 5. 推送 Tag 到遠端（可選）

```bash
# 推送 tag
git push origin v1.0.0-alpha.1

# 或推送所有 tags
git push origin --tags
```

### 版本格式說明

- **格式**：`alpha/beta+tag+commit`
- **範例**：`alpha/beta+v1.0.0-alpha.1+a24b17f`
- **說明**：包含完整 tag 名稱和 commit hash

### 測試版本迭代

如果需要發布新的測試版本：

```bash
# Alpha 迭代
make release-tag VERSION=v1.0.0-alpha.2
make build

# 從 Alpha 升級到 Beta
make release-tag VERSION=v1.0.0-beta.1
make build

# Beta 到 RC
make release-tag VERSION=v1.0.0-rc.1
make build
```

---

## 正式版本構建

正式版本用於生產環境部署和正式發布，**需要創建正式的語義化版本 Git tag**。

### 使用場景

- 生產環境部署
- 正式對外發布
- 客戶交付版本

### Tag 命名規範

```
v主版本.次版本.修訂版本

範例：
  v1.0.0      # 首次正式發布
  v1.0.1      # Bug 修復版本
  v1.1.0      # 新功能版本（向後相容）
  v2.0.0      # 重大變更版本（不向後相容）
```

**語義化版本說明：**

- **主版本**（Major）：不向後相容的 API 變更
- **次版本**（Minor）：向後相容的新功能
- **修訂版本**（Patch）：向後相容的 bug 修復

### 構建步驟（推薦流程）

#### 1. 切換到主分支

```bash
git checkout main
git pull origin main
```

#### 2. 合併功能分支（如有需要）

```bash
# 合併功能分支
git merge feature/new-feature

# 解決衝突（如有）
git add .
git commit -m "Merge feature/new-feature"
```

#### 3. 確保所有變更已提交

```bash
# 查看狀態
make release-info

# 如有未提交變更
git add .
git commit -m "Prepare release v1.0.0"
```

#### 4. 創建 Release Tag

```bash
make release-tag VERSION=v1.0.0
```

這個命令會：
- ✅ 檢查工作目錄是否乾淨
- ✅ 創建帶註解的 Git tag
- ✅ 顯示後續步驟提示

#### 5. 構建 Release 版本

```bash
make release-build
```

這個命令會：
- ✅ 再次檢查工作目錄是否乾淨
- ✅ 執行構建
- ✅ 顯示構建完成資訊

**或者使用一般構建：**
```bash
make build
```

#### 6. 驗證版本資訊

```bash
./bin/webprism version
```

**輸出範例：**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  WEBPRISM CLI [RELEASE]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Version:     release+v1.0.0+abc1234
  Git Commit:  abc1234567890abcdef1234567890abcdef123456
  Git Branch:  main
  Build Time:  2026-01-03T12:00:00Z
  Go Version:  go1.24.11
  Platform:    darwin/arm64
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### 7. 推送到遠端倉庫

```bash
# 推送代碼
git push origin main

# 推送 tag
git push origin v1.0.0
```

#### 8. 打包發布（可選）

```bash
# 為不同平台打包
tar -czf webprism-v1.0.0-darwin-arm64.tar.gz -C bin .
tar -czf webprism-v1.0.0-linux-amd64.tar.gz -C bin .
```

### 版本格式說明

- **格式**：`release+tag+commit`
- **範例**：`release+v1.0.0+a24b17f`
- **說明**：清楚標示為正式發布版本

---

## 版本查詢方式

### 1. CLI 命令查詢

#### 完整版本資訊
```bash
./bin/webprism version
```

#### 簡短版本
```bash
./bin/webprism version --short
```

輸出範例：
```
WEBPRISM build-foundation+a24b17f (commit: a24b17f)
```

#### JSON 格式
```bash
./bin/webprism version -j
```

輸出範例：
```json
{
  "version": "build-foundation+a24b17f",
  "git_commit": "a24b17f38bcac2479228987846c28a805cb683bb",
  "git_branch": "build-foundation",
  "build_time": "2026-01-03T12:00:00Z",
  "go_version": "go1.24.11",
  "platform": "darwin/arm64"
}
```

### 2. Makefile 命令查詢

```bash
make release-info
```

輸出範例：
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Current Version Information
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Version:     build-foundation+a24b17f
  Git Commit:  a24b17f38bcac2479228987846c28a805cb683bb
  Git Branch:  build-foundation
  Clean:       No (has uncommitted changes)

  All Tags:
    v1.0.0
    v1.0.0-alpha.1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 3. API 端點查詢

#### HTTP API
```bash
curl http://localhost:8080/api/v1/version
```

#### gRPC API
```bash
grpcurl -plaintext localhost:9090 webprism.v1.HealthService/GetVersion
```

### 4. 伺服器日誌

伺服器啟動時會自動記錄版本資訊：

```
INFO  Starting WEBPRISM server...
  version: build-foundation+a24b17f
  git_commit: a24b17f38bcac2479228987846c28a805cb683bb
  git_branch: build-foundation
  build_time: 2026-01-03T12:00:00Z
  go_version: go1.24.11
  platform: darwin/arm64
  development: true
```

---

## 常見問題

### Q1: 如何刪除錯誤的 Tag？

```bash
# 刪除本地 tag
git tag -d v1.0.0

# 刪除遠端 tag
git push origin --delete v1.0.0
```

### Q2: 版本顯示 `-dirty` 是什麼意思？

`-dirty` 表示工作目錄有未提交的變更。要獲得乾淨的版本：

```bash
git add .
git commit -m "Clean up changes"
make build
```

### Q3: 如何在構建時自定義版本號？

```bash
# 自定義版本號
VERSION=v2.0.0 make build

# 自定義所有構建變數
VERSION=v2.0.0 GIT_BRANCH=custom make build
```

### Q4: 如何查看所有可用的 Tags？

```bash
# 列出所有 tags
git tag -l

# 列出特定模式的 tags
git tag -l "v1.*"

# 查看 tag 的詳細資訊
git show v1.0.0
```

### Q5: 構建失敗提示工作目錄不乾淨？

使用 `release-build` 時會檢查工作目錄。解決方法：

```bash
# 查看狀態
git status

# 提交變更
git add .
git commit -m "Commit changes"

# 或使用一般構建（不檢查）
make build
```

### Q6: 如何為多個平台構建？

```bash
# macOS ARM64
GOOS=darwin GOARCH=arm64 make build

# macOS AMD64
GOOS=darwin GOARCH=amd64 make build

# Linux AMD64
GOOS=linux GOARCH=amd64 make build

# Windows AMD64
GOOS=windows GOARCH=amd64 make build
```

### Q7: 版本號太長怎麼辦？

使用短版本：

```bash
./bin/webprism version --short
```

或在代碼中使用 `version.Short()` 方法。

---

## 快速參考

### 常用命令

| 操作 | 命令 |
|-----|------|
| 查看版本資訊 | `make release-info` |
| 開發版本構建 | `make build` |
| 創建測試版本 | `make release-tag VERSION=v1.0.0-alpha.1` |
| 創建正式版本 | `make release-tag VERSION=v1.0.0` |
| Release 構建 | `make release-build` |
| 查看版本（完整） | `./bin/webprism version` |
| 查看版本（簡短） | `./bin/webprism version --short` |
| 查看版本（JSON） | `./bin/webprism version -j` |

### 版本類型對照表

| 你在 | 需要 Tag？ | 版本格式 | 標記 |
|-----|-----------|---------|------|
| 功能分支開發 | ❌ 不需要 | `branch+commit` | `[DEVELOPMENT]` |
| 測試階段 | ✅ `-alpha/-beta/-rc` | `alpha/beta+tag+commit` | `[ALPHA/BETA]` |
| 正式發布 | ✅ `vX.Y.Z` | `release+tag+commit` | `[RELEASE]` |

---

## 相關文檔

- [語義化版本規範](https://semver.org/)
- [Git Tag 文檔](https://git-scm.com/book/en/v2/Git-Basics-Tagging)
- [Go Build 文檔](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)

---

**文檔版本：** 1.0.0
**最後更新：** 2026-01-03
**維護者：** WEBPRISM Team
