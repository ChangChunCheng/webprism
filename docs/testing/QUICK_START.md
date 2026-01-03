# 測試快速指南

快速參考文件，用於執行和維護 WEBPRISM 的測試。

---

## 🚀 快速開始

### 執行所有測試

```bash
go test -v ./internal/domain/service/
```

### 執行特定測試

```bash
# 只執行 ProxyService 測試
go test -v -run TestProxyService ./internal/domain/service/

# 執行特定測試案例
go test -v -run TestProxyService/successful_proxy_without_auth ./internal/domain/service/
```

---

## 📊 測試覆蓋率

### 生成覆蓋率報告

```bash
# 生成覆蓋率檔案
go test -coverprofile=coverage.out ./internal/domain/service/

# 查看覆蓋率摘要
go tool cover -func=coverage.out

# 生成 HTML 報告
go tool cover -html=coverage.out -o coverage.html

# 在瀏覽器中查看
open coverage.html  # macOS
```

### 查看特定服務的覆蓋率

```bash
# 只看 ProxyService 覆蓋率
go test -coverprofile=coverage.out ./internal/domain/service/
go tool cover -func=coverage.out | grep proxy_service
```

---

## 🔧 Mock 管理

### 重新生成所有 Mocks

```bash
mockery
```

### 生成特定 Mock

```bash
mockery --name=SpecRepository
mockery --name=AuthRepository
mockery --name=HTTPClient
```

### Mock 配置檔案

編輯 `.mockery.yaml` 來添加新的 interface:

```yaml
with-expecter: true
dir: "internal/mocks"
outpkg: mocks
packages:
  github.com/ChangChunCheng/webprism/internal/ports/output:
    interfaces:
      YourNewInterface:  # 添加這裡
```

---

## 🧪 測試開發

### 使用 testutil helpers

```go
// 創建測試 Spec
spec := testutil.CreateTestSpecWithPaths("spec-id", "My API")

// 創建測試請求
req := testutil.CreateTestProxyRequest("spec-id", "operationId")

// 創建認證配置
authConfig := testutil.CreateTestAuthConfig(model.AuthTypeBearer)
```

### Table-Driven Tests 模板

```go
func TestMyService_MyMethod(t *testing.T) {
    tests := []struct {
        name       string
        input      string
        setupMocks func(*mocks.MockDependency)
        wantErr    bool
        validate   func(*testing.T, *Result)
    }{
        {
            name:  "success case",
            input: "test-input",
            setupMocks: func(mock *mocks.MockDependency) {
                mock.EXPECT().
                    SomeMethod(mock.Anything, "test-input").
                    Return(&Response{}, nil)
            },
            wantErr: false,
            validate: func(t *testing.T, result *Result) {
                assert.NotNil(t, result)
            },
        },
        // 更多測試案例...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockDep := mocks.NewMockDependency(t)
            tt.setupMocks(mockDep)
            service := NewMyService(mockDep)

            // Act
            result, err := service.MyMethod(context.Background(), tt.input)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            if tt.validate != nil {
                tt.validate(t, result)
            }
        })
    }
}
```

---

## 📈 測試目標

### 覆蓋率目標

| Service | 目標覆蓋率 | 當前狀態 |
|---------|------------|----------|
| ProxyService | >85% | ✅ 93.1% |
| SpecService | >85% | ⏳ 待實作 |
| AuthService | >85% | ⏳ 待實作 |
| HealthService | >85% | ⏳ 待實作 |

### 測試品質檢查清單

每個測試應該:
- [ ] 獨立執行 (不依賴其他測試)
- [ ] 有清楚的測試名稱
- [ ] 遵循 AAA 模式 (Arrange-Act-Assert)
- [ ] 使用 testutil helpers (不硬編碼資料)
- [ ] 完整 mock 所有依賴
- [ ] 驗證錯誤情況
- [ ] 執行速度快 (< 100ms)

---

## 🐛 常見問題

### 測試失敗: "mock: I don't know what to return"

**原因**: Mock 沒有設定預期行為

**解決**:
```go
mockRepo.EXPECT().
    FindByID(mock.Anything, "some-id").
    Return(&Entity{}, nil)
```

### 測試失敗: Type mismatch

**原因**: 測試資料型別與實際不符

**解決**: 檢查 model 定義，確保使用正確的型別

### Mock 沒有更新

**原因**: 修改了 interface 但沒有重新生成 mock

**解決**:
```bash
mockery  # 重新生成所有 mocks
```

---

## 📚 相關文件

- [完整測試報告](./UNIT_TEST_REPORT.md) - 詳細的測試分析
- [測試策略](./TESTING_STRATEGY.md) - 整體測試規劃
- [工作總結](../development/WORK_SUMMARY_2026-01-01_UNIT_TESTS.md) - 開發過程記錄

---

## 🎯 下一步

1. **實作 SpecService 測試** (~10-12 個測試案例)
2. **實作 AuthService 測試** (~8-10 個測試案例)
3. **實作 HealthService 測試** (~6-8 個測試案例)
4. **整合到 CI/CD** (自動執行測試)
5. **Integration Tests** (資料庫、HTTP 層測試)

---

**最後更新**: 2026-01-01
