# 📚 WEBPRISM 文件中心

歡迎來到 WEBPRISM 文件中心！這裡是所有文件的導航頁面。

---

## 🎯 我應該看哪份文件？

### 👤 我是使用者

**從這裡開始** → [GETTING_STARTED.md](GETTING_STARTED.md)

1. **[快速開始](GETTING_STARTED.md)** - 完整的新手入門指南（推薦）
2. [安裝指南](user-guide/01-INSTALLATION.md) - 詳細安裝步驟
3. [使用教學](user-guide/02-QUICKSTART.md) - 第一次使用教學
4. [系統啟動指南](user-guide/05-SYSTEM_STARTUP.md) - 完整的啟動與測試流程

### 👨‍💻 我是開發者

**從這裡開始** → [development/CONTRIBUTING.md](development/CONTRIBUTING.md)

1. **[開發指南](development/CONTRIBUTING.md)** - 開發環境設置與工作流程（推薦）
2. **[版本管理指南](VERSION_MANAGEMENT.md)** - 版本類型、構建流程和 Release 發布（必讀）
3. [當前狀態](development/CURRENT_STATUS.md) - 專案當前進度
4. [測試指南](testing/QUICK_START.md) - 測試執行與覆蓋率

### 🧪 我要執行測試

**測試相關** → [testing/](testing/)

1. [測試快速指南](testing/QUICK_START.md) - 如何執行測試
2. [測試完成報告](testing/COMPLETION_SUMMARY.md) - 測試狀態與結果
3. [測試文件索引](testing/README.md) - 所有測試文件導航

---

## 📂 文件結構

```plaintext
docs/
├── README.md                          # 📍 你在這裡 - 文件導航中心
├── GETTING_STARTED.md                 # 🚀 使用者快速開始指南
├── VERSION_MANAGEMENT.md              # 🏷️  版本管理與構建指南
│
├── user-guide/                        # 👤 使用者指南
│   ├── 01-INSTALLATION.md             # 安裝指南
│   ├── 02-QUICKSTART.md               # 快速開始
│   ├── 03-INSTALL_TOOLS.md            # 工具安裝
│   ├── 04-MACOS_SETUP.md              # macOS 專用設定
│   ├── 05-SYSTEM_STARTUP.md           # 系統啟動與測試（完整流程）
│   └── START_SERVER_LOCAL.md          # 本機啟動指南
│
├── development/                       # 👨‍💻 開發者文件
│   ├── CONTRIBUTING.md                # 🆕 開發指南（開發者必讀）
│   ├── CURRENT_STATUS.md              # 當前狀態
│   ├── WORK_SUMMARY_2026-01-01.md     # 工作記錄
│   ├── WORK_SUMMARY_2026-01-01_UNIT_TESTS.md
│   ├── BUG_FIXES_2026-01-01.md
│   ├── BUILD_SUCCESS.md
│   ├── CLEANUP_AND_FIX.md
│   ├── ISSUES_FIXED.md
│   ├── V1_COMPLETE_SUMMARY.md
│   └── DOCUMENTATION_READY.md
│
├── testing/                           # 🧪 測試文件
│   ├── README.md                      # 測試文件索引
│   ├── QUICK_START.md                 # 測試快速指南
│   ├── COMPLETION_SUMMARY.md          # 測試完成總結
│   ├── TESTING_COMPLETE.md            # 單元測試完成狀態
│   ├── MANUAL_TESTING_GUIDE.md        # 手動測試指南
│   ├── MANUAL_TEST_RESULTS_2026-01-01.md
│   ├── TEST_EXECUTION_SUMMARY.md
│   ├── FINAL_TEST_REPORT_2026-01-01.md
│   ├── UNIT_TEST_REPORT.md
│   ├── TESTING_STRATEGY.md
│   └── UNIT_TEST_PROGRESS.md
│
└── archive/                           # 📦 歷史文件
    ├── ARCHITECTURE.md
    ├── DEVELOPMENT_PLAN.md
    ├── DEVELOPMENT_PROGRESS.md
    └── ...
```

---

## 🚀 快速導航

### 按任務分類

#### 🎬 第一次使用 WEBPRISM

**推薦**: 閱讀 **[快速開始指南](GETTING_STARTED.md)** (15-20 分鐘)

或分步閱讀：

1. **[安裝指南](user-guide/01-INSTALLATION.md)** (10 分鐘) - 系統需求與工具安裝
2. **[使用教學](user-guide/02-QUICKSTART.md)** (5 分鐘) - Docker 快速啟動
3. **[系統啟動指南](user-guide/05-SYSTEM_STARTUP.md)** (15 分鐘) - 完整測試流程

#### 👨‍💻 參與開發

**推薦**: 閱讀 **[開發指南](development/CONTRIBUTING.md)** (20-30 分鐘)

包含：

- 開發環境設置
- 專案結構說明
- 程式碼規範
- 測試要求
- 提交規範

#### 🔧 工具安裝問題

- [工具安裝指南](user-guide/03-INSTALL_TOOLS.md)
- [macOS 專用設定](user-guide/04-MACOS_SETUP.md)

#### 🧪 執行測試

- [測試快速指南](testing/QUICK_START.md) - 如何執行測試
- [測試完成報告](testing/COMPLETION_SUMMARY.md) - 當前測試狀態

#### 🐛 問題排除

- [快速開始 - 常見問題](GETTING_STARTED.md#常見問題)
- [開發指南 - 故障排除](development/CONTRIBUTING.md#故障排除)
- [已修復問題](development/ISSUES_FIXED.md)

#### 📊 查看專案狀態

- [當前狀態](development/CURRENT_STATUS.md)
- [測試完成報告](testing/COMPLETION_SUMMARY.md)
- [V1 完成總結](development/V1_COMPLETE_SUMMARY.md)

---

## 📖 文件索引

### 📚 核心文件

| 文件 | 說明 | 適合對象 |
| --- | --- | --- |
| **[GETTING_STARTED.md](GETTING_STARTED.md)** | **使用者快速開始指南** | 新使用者（推薦） |
| **[VERSION_MANAGEMENT.md](VERSION_MANAGEMENT.md)** | **版本管理與構建指南** | 開發者（必讀） |
| **[development/CONTRIBUTING.md](development/CONTRIBUTING.md)** | **開發者貢獻指南** | 開發者（推薦） |

### 👤 使用者指南 (User Guide)

| 文件 | 說明 | 適合對象 |
| --- | --- | --- |
| [01-INSTALLATION.md](user-guide/01-INSTALLATION.md) | 完整安裝指南 | 新使用者 |
| [02-QUICKSTART.md](user-guide/02-QUICKSTART.md) | 5 分鐘快速開始 | 所有人 |
| [03-INSTALL_TOOLS.md](user-guide/03-INSTALL_TOOLS.md) | 工具安裝說明 | 需要安裝工具者 |
| [04-MACOS_SETUP.md](user-guide/04-MACOS_SETUP.md) | macOS 特定設定 | macOS 使用者 |
| [05-SYSTEM_STARTUP.md](user-guide/05-SYSTEM_STARTUP.md) | 完整啟動與測試流程 | 第一次啟動者 |

### 👨‍💻 開發者文件 (Development)

| 文件 | 說明 | 適合對象 |
| --- | --- | --- |
| **[CONTRIBUTING.md](development/CONTRIBUTING.md)** | **開發指南（必讀）** | 開發者 |
| [CURRENT_STATUS.md](development/CURRENT_STATUS.md) | 當前專案狀態 | 所有人 |
| [WORK_SUMMARY_2026-01-01_UNIT_TESTS.md](development/WORK_SUMMARY_2026-01-01_UNIT_TESTS.md) | 單元測試工作記錄 | 開發者 |
| [BUG_FIXES_2026-01-01.md](development/BUG_FIXES_2026-01-01.md) | Bug 修復記錄 | 開發者 |
| [V1_COMPLETE_SUMMARY.md](development/V1_COMPLETE_SUMMARY.md) | V1 完成總結 | PM/Lead |

### 🧪 測試文件 (Testing)

| 文件 | 說明 | 適合對象 |
| --- | --- | --- |
| [README.md](testing/README.md) | 測試文件索引 | 所有人 |
| [QUICK_START.md](testing/QUICK_START.md) | 測試快速指南 | 開發者/QA |
| [COMPLETION_SUMMARY.md](testing/COMPLETION_SUMMARY.md) | 完整測試總結 | PM/Tech Lead |
| [FINAL_TEST_REPORT_2026-01-01.md](testing/FINAL_TEST_REPORT_2026-01-01.md) | 最終測試報告 | 開發者/QA |
| [TESTING_STRATEGY.md](testing/TESTING_STRATEGY.md) | 測試策略 | 架構師 |

---

## 🎯 常見任務快速指引

### 任務 1: 第一次使用 WEBPRISM

```bash
# 閱讀快速開始指南
open docs/GETTING_STARTED.md

# 按照指南步驟操作
make setup
make db-up
make run-server
```

### 任務 2: 參與開發

```bash
# 閱讀開發指南
open docs/development/CONTRIBUTING.md

# 閱讀版本管理指南
open docs/VERSION_MANAGEMENT.md

# 設置開發環境
make check-tools
make setup
make test
```

### 任務 3: 執行測試

```bash
# 閱讀測試指南
open docs/testing/QUICK_START.md

# 執行測試
make test
make test-coverage
```

### 任務 4: 排除問題

```bash
# 查看故障排除
open docs/GETTING_STARTED.md#常見問題
open docs/development/CONTRIBUTING.md#故障排除
open docs/development/ISSUES_FIXED.md
```

---

## 🔍 按關鍵字搜尋

- **新手入門**: [GETTING_STARTED.md](GETTING_STARTED.md)
- **版本管理**: [VERSION_MANAGEMENT.md](VERSION_MANAGEMENT.md)
- **構建發布**: [VERSION_MANAGEMENT.md](VERSION_MANAGEMENT.md)
- **開發環境**: [development/CONTRIBUTING.md](development/CONTRIBUTING.md)
- **安裝工具**: [user-guide/03-INSTALL_TOOLS.md](user-guide/03-INSTALL_TOOLS.md)
- **Docker**: [GETTING_STARTED.md](GETTING_STARTED.md#啟動模式)
- **測試相關**: [testing/QUICK_START.md](testing/QUICK_START.md)
- **故障排除**: [GETTING_STARTED.md](GETTING_STARTED.md#常見問題), [development/CONTRIBUTING.md](development/CONTRIBUTING.md#故障排除)
- **PostgreSQL**: [GETTING_STARTED.md](GETTING_STARTED.md#啟動模式)
- **程式碼規範**: [development/CONTRIBUTING.md](development/CONTRIBUTING.md#程式碼規範)
- **Commit 規範**: [development/CONTRIBUTING.md](development/CONTRIBUTING.md#提交規範)

---

## 📞 需要協助？

### 找不到想要的資訊？

1. **使用者**: 先看 [快速開始指南](GETTING_STARTED.md)
2. **開發者**: 先看 [開發指南](development/CONTRIBUTING.md)
3. **測試相關**: 查看 [測試文件索引](testing/README.md)
4. **專案概覽**: 查看 [專案 README](../README.md)

### 遇到問題？

1. **常見問題**:
   - [使用者常見問題](GETTING_STARTED.md#常見問題)
   - [開發者故障排除](development/CONTRIBUTING.md#故障排除)

2. **已知問題**:
   - [已修復問題](development/ISSUES_FIXED.md)
   - [Bug 修復記錄](development/BUG_FIXES_2026-01-01.md)

3. **提出問題**: [GitHub Issues](https://github.com/ChangChunCheng/webprism/issues)

---

## 📝 文件維護

### 文件更新記錄

- **2026-01-03**: 重構文件結構，新增快速開始指南和開發指南
- **2026-01-01**: 創建文件中心，整理文件結構
- **2026-01-01**: 完成單元測試，新增測試文件
- **2026-01-01**: 創建系統啟動指南

### 文件版本

所有文件與專案版本保持同步。當前版本: **v1.0**

### 文件分類原則

- **README.md**: 專案概覽、技術棧、快速開始
- **docs/GETTING_STARTED.md**: 使用者完整入門指南
- **docs/development/CONTRIBUTING.md**: 開發者貢獻指南
- **docs/user-guide/**: 詳細使用文件
- **docs/development/**: 開發相關文件
- **docs/testing/**: 測試相關文件
- **docs/archive/**: 歷史文件

---

## 🗂️ 歷史文件

開發過程中的規劃文件已移至 [archive/](archive/) 目錄:

- 架構設計文件
- 開發計劃
- 技術決策記錄
- 早期進度追蹤

這些文件保留作為歷史參考。

---

**祝您使用愉快！** 🚀

如有任何問題或建議，歡迎提出 Issue 或 Pull Request。
