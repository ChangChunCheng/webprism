# WEBPRISM 測試指南

**狀態**: ✅ 單元測試完成 (95.2% 覆蓋率) | ✅ 整合測試完成

---

## 🎯 快速開始

### 測試類型

WEBPRISM 提供多種測試類型：

- **[QUICK_START.md](QUICK_START.md)** - 單元測試快速指南
  - 如何執行單元測試
  - 測試覆蓋率報告
  - Mock 管理

- **[INTEGRATION_TEST_GUIDE.md](INTEGRATION_TEST_GUIDE.md)** - 整合測試指南 ⭐ 推薦
  - 完整的 API 工作流程測試
  - 與真實外部 API 整合
  - 自動化測試腳本

---

## 📊 測試狀態

### 單元測試

```
✅ 測試數量: 80 個測試
✅ 通過率: 100% (80/80)
✅ 覆蓋率: 95.2%
✅ 執行時間: 0.273s
```

### 整合測試

```
✅ 測試腳本: scripts/test-petstore-full.sh
✅ 測試步驟: 10 個步驟
✅ OpenAPI 支援: 3.0.4, Swagger 2.0
✅ 認證類型: API Key, Bearer Token
✅ 參數類型: Path, Query, Header, Body
```

### 測試覆蓋範圍

| 模組 | 單元測試覆蓋率 | 整合測試 | 狀態 |
| ------ | ---------------- | ---------- | ------ |
| SpecService | 100% | ✅ | ✅ |
| AuthService | 95.2% | ✅ | ✅ |
| ProxyService | 93.1% | ✅ | ✅ |
| HealthService | 93.8% | ✅ | ✅ |

---

## 🧪 執行測試

### 執行單元測試

```bash
# 執行所有單元測試
make test

# 執行特定服務測試
go test -v ./internal/domain/service/ -run TestSpecService
go test -v ./internal/domain/service/ -run TestAuthService
go test -v ./internal/domain/service/ -run TestProxyService
```

### 執行整合測試 ⭐ 推薦

```bash
# 1. 啟動服務
make db-up        # 啟動資料庫
make run-server   # 啟動 WEBPRISM (另一個終端)

# 2. 執行完整整合測試
./scripts/test-petstore-full.sh

# 3. 查看測試日誌
ls -lh /tmp/webprism-test-petstore-full-*.log
```

詳細說明請參考 [整合測試指南](INTEGRATION_TEST_GUIDE.md)

### 查看測試覆蓋率

```bash
make test-coverage
```

這會生成 `coverage.html` 文件,用瀏覽器打開查看詳細覆蓋率報告。

---

## 📁 測試結構

```
internal/
├── domain/
│   └── service/
│       ├── spec_service.go
│       ├── spec_service_test.go      # SpecService 測試
│       ├── auth_service.go
│       ├── auth_service_test.go      # AuthService 測試
│       ├── proxy_service.go
│       ├── proxy_service_test.go     # ProxyService 測試
│       ├── health_service.go
│       └── health_service_test.go    # HealthService 測試
└── mocks/                             # 自動生成的 Mock
    ├── MockSpecRepository.go
    ├── MockAuthRepository.go
    └── MockProxyClient.go
```

---

## 🔧 測試工具

### 必要工具

- **Go 1.21+** - 執行測試
- **mockery** - 生成 Mock (已包含在 `make install-tools` 中)

### 生成 Mock

```bash
make mocks
```

這會根據 `.mockery.yaml` 配置自動生成所有介面的 Mock 實現。

---

## 📚 測試最佳實踐

### 1. 使用 Table-Driven Tests

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {
        name:    "valid input",
        input:   "test",
        want:    "result",
        wantErr: false,
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test logic
    })
}
```

### 2. 使用 Mock 隔離依賴

```go
mockRepo := mocks.NewMockSpecRepository(t)
mockRepo.EXPECT().FindByID(ctx, specID).Return(spec, nil)

service := service.NewSpecService(mockRepo)
```

### 3. 測試錯誤情況

確保測試包含:
- ✅ 正常情況 (Happy Path)
- ✅ 錯誤處理 (Error Cases)
- ✅ 邊界條件 (Edge Cases)
- ✅ 空值處理 (Nil Cases)

---

## 🐛 故障排除

### 問題 1: Mock 生成失敗

```bash
# 清理舊的 Mock
make clean-all

# 重新生成
make mocks
```

### 問題 2: 測試找不到依賴

```bash
# 重新下載依賴
go mod download
go mod tidy
```

### 問題 3: 覆蓋率報告無法生成

```bash
# 確保有寫入權限
chmod +w .

# 重新執行
make test-coverage
```

---

## 🎯 下一步

- 查看 [QUICK_START.md](QUICK_START.md) 了解詳細測試指令
- 執行 `make test` 驗證所有測試通過
- 查看覆蓋率報告確保代碼品質

---

**測試是保證代碼品質的關鍵!** 🧪
