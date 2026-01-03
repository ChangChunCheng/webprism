# WEBPRISM 測試資料說明

本目錄包含用於測試 WEBPRISM 的 OpenAPI 規格檔案。

## 📋 包含的範例檔案

### petstore-simple.json
- **用途**: 簡單的 Petstore API 範例，用於基本功能測試
- **來源**: OpenAPI 官方範例
- **狀態**: ✅ 已驗證可用
- **大小**: 約 1.5KB (72 行)

## 🔧 添加自己的測試資料

您可以添加任何 OpenAPI 規格檔案到此目錄進行測試：

```bash
# 1. 將您的 OpenAPI 規格放入 testdata/
cp your-api-spec.json testdata/

# 2. 使用 WEBPRISM 上傳測試
./bin/webprism spec upload testdata/your-api-spec.json \
  --name "Your API" \
  --version "1.0.0"
```

## 📚 取得更多測試資料

### 台灣證券交易所 OpenAPI

如需測試大型、複雜的 OpenAPI 規格，可以使用台灣證券交易所的公開 API：

```bash
# 下載 TWSE OpenAPI 規格（約 653KB）
curl -o testdata/twse-openapi.json \
  https://openapi.twse.com.tw/v1/openapi/swagger.json

# 或使用瀏覽器訪問
# https://openapi.twse.com.tw/v1/openapi/swagger.json
```

### 其他公開 API 資源

- **OpenAPI Directory**: https://apis.guru/browse-apis/
- **Public APIs**: https://github.com/public-apis/public-apis
- **Swagger Petstore**: https://petstore.swagger.io/

## ⚠️ 注意事項

- **預設不提交**: testdata/ 中的檔案預設不會被 git 追蹤（除了 petstore-simple.json）
- **明確加入**: 如果您想提交某個測試檔案，請在 `.gitignore` 中明確添加例外規則：
  ```gitignore
  !testdata/your-verified-spec.json
  ```
- **檔案大小**: 避免提交超大的測試檔案（建議 < 1MB）
- **機密資訊**: 確保測試檔案中不包含任何機密資訊（API keys、tokens 等）

## 🧪 測試建議

1. **小型測試**: 使用 `petstore-simple.json` 進行快速功能驗證
2. **複雜測試**: 使用 TWSE OpenAPI 或其他大型規格測試效能和邊界情況
3. **自定義測試**: 使用您的實際 API 規格進行整合測試

---

**提示**: 所有測試資料都應該是公開的、非機密的 OpenAPI 規格檔案。
