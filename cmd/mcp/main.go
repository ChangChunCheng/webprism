package main

import (
	"fmt"
	"os"

	"github.com/ChangChunCheng/webprism/internal/adapters/httpclient"
	"github.com/ChangChunCheng/webprism/internal/adapters/mcp"
	"github.com/ChangChunCheng/webprism/internal/adapters/parser"
	"github.com/ChangChunCheng/webprism/internal/adapters/storage/postgres"
	"github.com/ChangChunCheng/webprism/internal/domain/service"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/config"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/crypto"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/database"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/version"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger (stderr for MCP stdio mode)
	log, err := logger.New(logger.Config{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		OutputPath: "stderr",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()

	// 顯示版本資訊
	versionInfo := version.Get()
	log.Info("Starting WEBPRISM MCP server...",
		logger.String("version", versionInfo.FullVersion()),
		logger.String("git_commit", versionInfo.GitCommit),
		logger.String("git_branch", versionInfo.GitBranch),
		logger.String("build_time", versionInfo.BuildTime),
		logger.String("go_version", versionInfo.GoVersion),
		logger.String("platform", versionInfo.Platform),
		logger.Bool("development", versionInfo.IsDevelopment()),
	)

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database, log)
	if err != nil {
		log.Error("Failed to initialize database", logger.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// Initialize crypto service
	cryptoService, err := crypto.NewAESCrypto(cfg.Security.EncryptionKey)
	if err != nil {
		log.Error("Failed to initialize crypto service", logger.Any("error", err))
		os.Exit(1)
	}

	// Initialize repositories
	specRepo := postgres.NewSpecRepository(db.DB, log)
	authRepo := postgres.NewAuthRepository(db.DB, log)
	healthRepo := postgres.NewHealthRepository(db.DB, log)

	// Initialize adapters
	openAPIParser := parser.NewOpenAPIParser(log)
	httpClient := httpclient.NewClient(log)
	requestBuilder := httpclient.NewRequestBuilder(log)

	// Initialize domain services
	specService := service.NewSpecService(specRepo, openAPIParser, log)
	authService := service.NewAuthService(authRepo, specRepo, cryptoService, log)
	proxyService := service.NewProxyService(specRepo, authRepo, httpClient, requestBuilder, cryptoService, openAPIParser, log)
	healthService := service.NewHealthService(specRepo, healthRepo, httpClient, log)

	// Initialize MCP server
	mcpServer := mcp.NewServer(
		cfg.MCP.Name,
		cfg.MCP.Version,
		specService,
		authService,
		proxyService,
		healthService,
		log,
	)

	// Start MCP server (stdio mode)
	if err := mcpServer.Start(); err != nil {
		log.Error("MCP server failed", logger.Any("error", err))
		os.Exit(1)
	}
}
