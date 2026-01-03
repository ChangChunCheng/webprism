package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/input"
)

// Server wraps the MCP server.
type Server struct {
	mcpServer *server.MCPServer
	tools     *Tools
	logger    *logger.Logger
}

// NewServer creates a new MCP server with all tools registered.
func NewServer(
	name, version string,
	specService input.SpecService,
	authService input.AuthService,
	proxyService input.ProxyService,
	healthService input.HealthService,
	log *logger.Logger,
) *Server {
	// Create MCP server
	mcpServer := server.NewMCPServer(name, version)

	// Create tools
	tools := NewTools(specService, authService, proxyService, healthService, log)

	// Register all tools
	tools.RegisterTools(mcpServer)

	log.Info("MCP server created",
		logger.String("name", name),
		logger.String("version", version),
		logger.Int("tools_count", 6),
	)

	return &Server{
		mcpServer: mcpServer,
		tools:     tools,
		logger:    log,
	}
}

// Start starts the MCP server (stdio mode).
func (s *Server) Start() error {
	s.logger.Info("MCP server starting (stdio mode)")

	if err := server.ServeStdio(s.mcpServer); err != nil {
		s.logger.Error("MCP server failed", logger.Any("error", err))
		return err
	}

	return nil
}

// Stop gracefully stops the MCP server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("MCP server stopping...")

	// MCP stdio server doesn't need explicit cleanup
	s.logger.Info("MCP server stopped gracefully")
	return nil
}
