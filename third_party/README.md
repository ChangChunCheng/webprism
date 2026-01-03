# Third-Party Dependencies

本目錄包含 WEBPRISM 構建所需的第三方 proto 定義檔案。

## 📋 依賴清單

### 1. googleapis
- **來源**: https://github.com/googleapis/googleapis
- **用途**: Google API 標準定義（gRPC-Gateway 必需）
- **授權**: Apache License 2.0
- **需要的檔案**:
  - `google/api/annotations.proto` - HTTP API 註解
  - `google/api/http.proto` - HTTP 規則定義
  - `google/api/field_behavior.proto` - 欄位行為定義

## 🚀 設定方式

### 自動設定（推薦）

使用提供的設定腳本：

```bash
# 下載所有第三方依賴
./scripts/setup-third-party.sh

# 或使用 Makefile
make setup-third-party
```

這個腳本會：
1. 從官方 repository 下載指定版本的 googleapis
2. 使用 sparse checkout 只下載必要的檔案
3. 驗證所有必要的 proto 檔案是否存在

### 手動設定（進階）

如果您需要手動設定：

```bash
# 1. 進入 third_party 目錄
cd third_party

# 2. Clone googleapis（使用 sparse checkout）
git clone --filter=blob:none --no-checkout https://github.com/googleapis/googleapis.git
cd googleapis

# 3. Checkout 指定版本
git checkout 2f9af297c84c55c8b871ba4495e01ade42476c92

# 4. 只獲取需要的目錄
git sparse-checkout init --cone
git sparse-checkout set google/api google/rpc
```

## 🔍 版本資訊

當前使用的版本：

| 依賴 | Commit Hash | 日期 | 說明 |
|------|-------------|------|------|
| googleapis | `2f9af297c84c55c8b871ba4495e01ade42476c92` | 2024-01 | 穩定版本 |

## ⚠️ 重要提醒

### 為什麼不直接提交到 Git？

1. **檔案大小** - googleapis 完整 repository 超過 100MB
2. **版本控制** - 透過腳本管理更容易追蹤和更新版本
3. **授權清楚** - 明確這些是外部依賴，有各自的授權
4. **減少衝突** - 不會因為 proto 檔案變更導致 merge conflict

### 更新第三方依賴

如需更新到新版本的 googleapis：

1. 到 https://github.com/googleapis/googleapis 找到新的 commit hash
2. 更新 `scripts/setup-third-party.sh` 中的 `GOOGLEAPIS_COMMIT`
3. 更新本 README 中的版本資訊表格
4. 重新執行設定腳本：`./scripts/setup-third-party.sh`
5. 測試編譯：`make proto && make build`

## 🆘 常見問題

### Q: 為什麼 make proto 失敗說找不到 google/api/annotations.proto？

A: 您需要先設定第三方依賴：
```bash
./scripts/setup-third-party.sh
# 或
make setup-third-party
```

### Q: 可以使用其他版本的 googleapis 嗎？

A: 可以，但建議使用腳本指定的版本以確保相容性。如需使用其他版本：
1. 修改 `setup-third-party.sh` 中的 commit hash
2. 測試是否能正常編譯
3. 如有問題，請回報 issue

### Q: third_party/ 目錄可以刪除嗎？

A: 可以。這個目錄的內容會被 `.gitignore` 忽略，可以安全地刪除並重新下載：
```bash
rm -rf third_party/googleapis
./scripts/setup-third-party.sh
```

## 📚 延伸閱讀

- [googleapis Repository](https://github.com/googleapis/googleapis)
- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Protocol Buffers Guide](https://protobuf.dev/)

---

**注意**: 本目錄中的第三方依賴檔案不會被提交到 Git，會在專案設定時自動下載。
