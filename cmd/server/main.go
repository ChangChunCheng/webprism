package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChangChunCheng/webprism/internal/adapters/grpc"
	"github.com/ChangChunCheng/webprism/internal/adapters/http"
	"github.com/ChangChunCheng/webprism/internal/adapters/httpclient"
	"github.com/ChangChunCheng/webprism/internal/adapters/parser"
	"github.com/ChangChunCheng/webprism/internal/adapters/storage/postgres"
	"github.com/ChangChunCheng/webprism/internal/domain/service"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/config"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/crypto"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/database"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		OutputPath: cfg.Logging.OutputPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Starting WEBPRISM server...")

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

	// Initialize gRPC server
	grpcServer := grpc.NewServer(
		cfg.Server.GRPC.Port,
		specService,
		authService,
		proxyService,
		healthService,
		log,
	)

	// Initialize HTTP server (grpc-gateway)
	httpServer, err := http.NewServer(cfg.Server.HTTP.Port, cfg.Server.GRPC.Port, log)
	if err != nil {
		log.Error("Failed to initialize HTTP server", logger.Any("error", err))
		os.Exit(1)
	}

	// Start servers in goroutines
	errCh := make(chan error, 2)

	// Start gRPC server
	go func() {
		log.Info("Starting gRPC server", logger.Int("port", cfg.Server.GRPC.Port))
		if err := grpcServer.Start(); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Info("Starting HTTP server", logger.Int("port", cfg.Server.HTTP.Port))
		if err := httpServer.Start(); err != nil {
			errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Wait for interrupt signal or error
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Error("Server error", logger.Any("error", err))
	case sig := <-sigCh:
		log.Info("Received shutdown signal", logger.String("signal", sig.String()))
	}

	// Graceful shutdown
	log.Info("Shutting down servers...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop servers
	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown error", logger.Any("error", err))
	}

	if err := grpcServer.Stop(shutdownCtx); err != nil {
		log.Error("gRPC server shutdown error", logger.Any("error", err))
	}

	log.Info("WEBPRISM server stopped")
}
