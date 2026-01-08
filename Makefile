.PHONY: help proto build test run clean lint migrate docker-build docker-up docker-down

# 載入 .env 檔案（如果存在）
-include .env
export

# 預設目標
.DEFAULT_GOAL := help

# 幫助資訊
help: ## 顯示幫助資訊
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  WEBPRISM Makefile Commands"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "\033[1;36m📦 初始設定\033[0m"
	@grep -hE '^(setup|install-tools|setup-third-party|check-tools|deps):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;32m💻 本機開發模式（推薦）\033[0m"
	@echo "  \033[2m資料庫用 Docker，應用程式在本機執行\033[0m"
	@echo ""
	@echo "  \033[1;33m▸ 資料庫管理\033[0m"
	@grep -hE '^(db-up|db-down|db-logs):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "    \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "  \033[1;33m▸ 應用程式執行（go run - 開發用）\033[0m"
	@grep -hE '^(run-server|run-mcp|run-cli|dev):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "    \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "  \033[1;33m▸ 應用程式執行（bin - 含完整版本資訊）\033[0m"
	@grep -hE '^(run-server-bin|run-mcp-bin|run-cli-bin):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "    \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "  \033[1;33m▸ 資料庫遷移\033[0m"
	@grep -hE '^(migrate-up|migrate-down|migrate-create):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "    \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;35m🐳 Docker 容器模式\033[0m"
	@echo "  \033[2m完整容器化執行\033[0m"
	@echo ""
	@grep -hE '^(docker-build|docker-up-full|docker-up|docker-down|docker-logs|docker-clean):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;34m🔧 完全本機模式\033[0m"
	@echo "  \033[2m自行配置 PostgreSQL\033[0m"
	@echo ""
	@grep -hE '^(run-local):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;33m🔨 建置與生成\033[0m"
	@grep -hE '^(build|proto|mocks):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;33m🧪 測試與品質\033[0m"
	@grep -hE '^(test|test-coverage|lint|security|fmt):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;33m🧹 清理\033[0m"
	@grep -hE '^(clean|clean-all):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;35m🚀 Release 管理\033[0m"
	@grep -hE '^(release-check|release-tag|release-build|release-info):.*?## .*$$' $(MAKEFILE_LIST) | \
		sed 's/^\([^:]*\):.*## \(.*\)/\1|\2/' | \
		awk -F'|' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  \033[1;32m推薦工作流程:\033[0m"
	@echo "  1️⃣  初次設定: \033[36mmake setup\033[0m"
	@echo "  2️⃣  啟動資料庫: \033[36mmake db-up\033[0m"
	@echo "  3️⃣  執行應用: \033[36mmake run-server\033[0m (或 run-mcp, run-cli)"
	@echo ""
	@echo "  \033[1;35m發布版本:\033[0m"
	@echo "  1️⃣  查看狀態: \033[36mmake release-info\033[0m"
	@echo "  2️⃣  創建 tag: \033[36mmake release-tag VERSION=v1.0.0\033[0m"
	@echo "  3️⃣  構建發布: \033[36mmake release-build\033[0m"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 生成 Protobuf 程式碼
proto: ## 生成 Protobuf 程式碼
	@echo "Generating protobuf code..."
	@mkdir -p gen/go/v1
	@mkdir -p gen/openapi/v1
	@protoc -I api/proto \
		-I third_party/googleapis \
		--plugin=protoc-gen-go=$(shell go env GOPATH)/bin/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(shell go env GOPATH)/bin/protoc-gen-go-grpc \
		--plugin=protoc-gen-grpc-gateway=$(shell go env GOPATH)/bin/protoc-gen-grpc-gateway \
		--plugin=protoc-gen-openapiv2=$(shell go env GOPATH)/bin/protoc-gen-openapiv2 \
		--go_out=gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen/go \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=gen/go \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		--openapiv2_out=gen/openapi \
		api/proto/v1/*.proto
	@echo "Protobuf code generated successfully!"

# 安裝依賴
deps: ## 安裝 Go 依賴
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

# 構建資訊變數
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_VERSION ?= $(shell go version | awk '{print $$3}')
PLATFORM ?= $(shell go env GOOS)/$(shell go env GOARCH)

# 構建標記（ldflags）
LDFLAGS := -X 'github.com/ChangChunCheng/webprism/internal/version.Version=$(VERSION)' \
	-X 'github.com/ChangChunCheng/webprism/internal/version.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/ChangChunCheng/webprism/internal/version.GitBranch=$(GIT_BRANCH)' \
	-X 'github.com/ChangChunCheng/webprism/internal/version.BuildTime=$(BUILD_TIME)'

# 建置專案
build: ## 建置所有二進位檔案
	@echo "Building with version info:"
	@echo "  Version:    $(VERSION)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Git Branch: $(GIT_BRANCH)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(GO_VERSION)"
	@echo "  Platform:   $(PLATFORM)"
	@echo ""
	@mkdir -p bin
	@echo "Building webprism-server..."
	@go build -ldflags "$(LDFLAGS)" -o bin/webprism-server cmd/server/main.go
	@echo "Building webprism-mcp..."
	@go build -ldflags "$(LDFLAGS)" -o bin/webprism-mcp cmd/mcp/main.go
	@echo "Building webprism-cli..."
	@go build -ldflags "$(LDFLAGS)" -o bin/webprism cmd/cli/main.go
	@echo "Build complete!"

# 執行測試
test: ## 執行所有測試
	@echo "Running tests..."
	@go test -v -cover ./...

# 執行測試並生成覆蓋率報告
test-coverage: ## 執行測試並生成覆蓋率報告
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 執行 Lint
lint: ## 執行程式碼檢查
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

# 執行安全性掃描
security: ## 執行安全性掃描
	@echo "Running gosec..."
	@gosec ./...

# 執行 HTTP Server
run-server: ## 執行 HTTP Server（讀取 .env）
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found"; \
		echo "Please run: cp .env.example .env"; \
		echo "Then edit .env to configure secrets"; \
		exit 1; \
	fi
	@echo "Starting HTTP Server..."
	@echo "Environment loaded from .env"
	@go run cmd/server/main.go

# 執行 MCP Server
run-mcp: ## 執行 MCP Server（讀取 .env）
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found"; \
		echo "Please run: cp .env.example .env"; \
		echo "Then edit .env to configure secrets"; \
		exit 1; \
	fi
	@echo "Starting MCP Server..."
	@echo "Environment loaded from .env"
	@go run cmd/mcp/main.go

# 執行 CLI
run-cli: ## 執行 CLI（讀取 .env）
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found"; \
		echo "Please run: cp .env.example .env"; \
		echo "Then edit .env to configure secrets"; \
		exit 1; \
	fi
	@go run cmd/cli/main.go

# ============================================
# 執行已建置的 bin 檔案（含完整版本資訊）
# ============================================

run-server-bin: ## 執行已建置的 webprism-server（含版本資訊，讀取 .env）
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found"; \
		echo "Please run: cp .env.example .env"; \
		echo "Then edit .env to configure secrets"; \
		exit 1; \
	fi
	@if [ ! -f bin/webprism-server ]; then \
		echo "❌ Error: bin/webprism-server not found"; \
		echo "Please run: make build"; \
		exit 1; \
	fi
	@echo "Starting webprism-server from bin/ (with version info)..."
	@echo "Environment loaded from .env"
	@set -a && . ./.env && set +a && ./bin/webprism-server

run-mcp-bin: ## 執行已建置的 webprism-mcp（含版本資訊，讀取 .env）
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found"; \
		echo "Please run: cp .env.example .env"; \
		echo "Then edit .env to configure secrets"; \
		exit 1; \
	fi
	@if [ ! -f bin/webprism-mcp ]; then \
		echo "❌ Error: bin/webprism-mcp not found"; \
		echo "Please run: make build"; \
		exit 1; \
	fi
	@echo "Starting webprism-mcp from bin/ (with version info)..."
	@echo "Environment loaded from .env"
	@set -a && . ./.env && set +a && ./bin/webprism-mcp

run-cli-bin: ## 執行已建置的 webprism CLI（含版本資訊，讀取 .env）
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found"; \
		echo "Please run: cp .env.example .env"; \
		echo "Then edit .env to configure secrets"; \
		exit 1; \
	fi
	@if [ ! -f bin/webprism ]; then \
		echo "❌ Error: bin/webprism not found"; \
		echo "Please run: make build"; \
		exit 1; \
	fi
	@set -a && . ./.env && set +a && ./bin/webprism $(ARGS)

# 資料庫遷移（向上）
migrate-up: ## 執行資料庫遷移（向上）
	@echo "Running database migrations up..."
	@migrate -path internal/adapters/storage/postgres/migrations \
		-database "postgres://webprism:webprism@localhost:5432/webprism?sslmode=disable" up

# 資料庫遷移（向下）
migrate-down: ## 執行資料庫遷移（向下）
	@echo "Running database migrations down..."
	@migrate -path internal/adapters/storage/postgres/migrations \
		-database "postgres://webprism:webprism@localhost:5432/webprism?sslmode=disable" down

# 建立新的遷移檔案
migrate-create: ## 建立新的遷移檔案 (usage: make migrate-create NAME=create_users_table)
	@migrate create -ext sql -dir internal/adapters/storage/postgres/migrations -seq $(NAME)

# 生成 Mock 程式碼
mocks: ## 生成測試用 Mock
	@echo "Generating mocks..."
	@mockery --config .mockery.yaml
	@echo "Mocks generated successfully!"

# 清理建置檔案
clean: ## 清理建置檔案
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# 清理所有生成的檔案
clean-all: clean ## 清理所有生成的檔案（包含 gen/ 和 mocks/）
	@echo "Cleaning all generated files..."
	@rm -rf gen/
	@rm -rf internal/mocks/
	@echo "All generated files cleaned!"

# ============================================
# 本機運行（Native Development）
# ============================================

run-local: ## 本機運行服務（需要自行配置 PostgreSQL）
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env file not found"; \
		echo "Please run: cp .env.example .env"; \
		echo "Then edit .env to configure secrets"; \
		exit 1; \
	fi
	@echo "Starting WEBPRISM locally..."
	@echo "Environment loaded from .env"
	@echo "Make sure PostgreSQL is configured and running"
	@go run cmd/server/main.go

# ============================================
# Docker 運行（混合模式 - 推薦）
# ============================================

db-up: ## 只啟動 PostgreSQL（Docker）- 推薦開發模式
	@echo "Starting PostgreSQL in Docker..."
	@docker compose up -d postgres
	@echo "PostgreSQL started on localhost:5432"
	@echo "Run 'make run-server' to start WEBPRISM locally"

db-down: ## 停止 PostgreSQL（Docker）
	@echo "Stopping PostgreSQL..."
	@docker compose down

db-logs: ## 查看 PostgreSQL 日誌
	@docker compose logs -f postgres

# ============================================
# Docker 運行（完整容器化）
# ============================================

docker-build: ## 建置 Docker 映像
	@echo "Building Docker image..."
	@docker build -t webprism:latest .
	@echo "Docker image built successfully!"

docker-up-full: ## 啟動所有服務（PostgreSQL + WEBPRISM）在 Docker 中
	@echo "Starting all services in Docker..."
	@docker compose --profile full up -d
	@echo "All services started!"
	@echo "WEBPRISM: http://localhost:8080"
	@echo "View logs: make docker-logs"

docker-up: db-up ## 啟動 PostgreSQL（預設混合模式）

docker-down: ## 停止所有 Docker 服務
	@echo "Stopping Docker Compose services..."
	@docker compose --profile full down
	@echo "Services stopped!"

docker-logs: ## 查看所有 Docker 日誌
	@docker compose --profile full logs -f

docker-clean: ## 清理 Docker 資源（包含 volumes）
	@echo "Cleaning Docker resources..."
	@docker compose down -v
	@echo "Docker resources cleaned!"

# 格式化程式碼
fmt: ## 格式化程式碼
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

# ============================================
# Release 管理
# ============================================

.PHONY: release-check
release-check: ## 檢查是否可以發布（無未提交變更）
	@echo "Checking release readiness..."
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "❌ Error: Working directory is not clean"; \
		echo "Please commit all changes before creating a release"; \
		exit 1; \
	fi
	@echo "✅ Working directory is clean"
	@echo "Current version: $$(git describe --tags --always 2>/dev/null || echo 'no tags')"

.PHONY: release-tag
release-tag: release-check ## 創建 release tag (usage: make release-tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Error: VERSION is required"; \
		echo "Usage: make release-tag VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "✅ Tag $(VERSION) created successfully"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Push tag: git push origin $(VERSION)"
	@echo "  2. Build release: make release-build"

.PHONY: release-build
release-build: release-check build ## 構建 release 版本（確保乾淨狀態）
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  Release Build Complete"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  Version:    $(VERSION)"
	@echo "  Binaries:   bin/"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

.PHONY: release-info
release-info: ## 顯示當前版本資訊
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  Current Version Information"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  Version:     $(VERSION)"
	@echo "  Git Commit:  $(GIT_COMMIT)"
	@echo "  Git Branch:  $(GIT_BRANCH)"
	@if [ -z "$$(git status --porcelain)" ]; then \
		echo "  Clean:       Yes"; \
	else \
		echo "  Clean:       No (has uncommitted changes)"; \
	fi
	@echo ""
	@echo "  All Tags:"
	@git tag -l | sed 's/^/    /' || echo "    (no tags)"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 安裝開發工具到專案本地
install-tools: ## 安裝開發工具到專案本地 GOPATH
	@echo "Installing development tools to GOPATH..."
	@echo "GOPATH: $(shell go env GOPATH)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@go install github.com/vektra/mockery/v2@latest
	@echo "Development tools installed to $(shell go env GOPATH)/bin"
	@echo ""
	@echo "Note: Make sure $(shell go env GOPATH)/bin is in your PATH"
	@echo "Add this to your ~/.zshrc or ~/.bashrc:"
	@echo '  export PATH="$$PATH:$(shell go env GOPATH)/bin"'

# 開發模式（自動重載）
dev: ## 開發模式
	@echo "Starting development mode..."
	@air

# 檢查工具安裝
check-tools: ## 檢查開發工具是否已安裝
	@chmod +x scripts/check-tools.sh
	@./scripts/check-tools.sh

# 設定第三方依賴
setup-third-party: ## 下載第三方 proto 依賴（googleapis）
	@chmod +x scripts/setup-third-party.sh
	@./scripts/setup-third-party.sh

# 完整設定（第一次使用）
setup: install-tools setup-third-party proto mocks deps ## 完整設定專案
	@echo "Project setup complete!"
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Next steps - Choose your development mode:"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "Option 1: 混合模式 (推薦) - PostgreSQL 用 Docker，WEBPRISM 本機運行"
	@echo "  1. Copy .env: cp .env.example .env"
	@echo "  2. Start DB: make db-up"
	@echo "  3. Run app:  make run-server"
	@echo ""
	@echo "Option 2: 完全本機 - 使用自訂的 PostgreSQL"
	@echo "  1. Copy .env: cp .env.example .env"
	@echo "  2. Configure PostgreSQL connection in .env"
	@echo "  3. Run app:  make run-local"
	@echo ""
	@echo "Option 3: 完全 Docker - 所有服務都在容器中"
	@echo "  1. Copy .env: cp .env.example .env"
	@echo "  2. Start all: make docker-up-full"
	@echo ""
	@echo "Run 'make test' to verify setup"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
