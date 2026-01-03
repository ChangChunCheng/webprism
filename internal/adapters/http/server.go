package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/ChangChunCheng/webprism/gen/go/v1"
	healthpb "github.com/ChangChunCheng/webprism/gen/go/v1"
	proxypb "github.com/ChangChunCheng/webprism/gen/go/v1"
	specpb "github.com/ChangChunCheng/webprism/gen/go/v1"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
)

// Server wraps the HTTP server with grpc-gateway.
type Server struct {
	httpServer *http.Server
	grpcPort   int
	logger     *logger.Logger
}

// NewServer creates a new HTTP server with grpc-gateway.
func NewServer(httpPort, grpcPort int, log *logger.Logger) (*Server, error) {
	// Create gRPC-gateway mux
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithErrorHandler(customErrorHandler(log)),
	)

	// Create gRPC client connection
	grpcAddr := fmt.Sprintf("localhost:%d", grpcPort)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Use background context for service registration
	// The handlers will maintain their own connections
	ctx := context.Background()

	// Register all service handlers
	if err := specpb.RegisterSpecServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register SpecService handler: %w", err)
	}

	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register AuthService handler: %w", err)
	}

	if err := proxypb.RegisterProxyServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register ProxyService handler: %w", err)
	}

	if err := healthpb.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register HealthService handler: %w", err)
	}

	// Create HTTP handler with middleware
	handler := corsMiddleware(
		loggingMiddleware(log,
			mux,
		),
	)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", httpPort),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		grpcPort:   grpcPort,
		logger:     log,
	}, nil
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	s.logger.Info("HTTP server starting",
		logger.String("addr", s.httpServer.Addr),
		logger.Int("grpc_port", s.grpcPort),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server stopping...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP server shutdown failed", logger.Any("error", err))
		return err
	}

	s.logger.Info("HTTP server stopped gracefully")
	return nil
}

// customErrorHandler handles errors from grpc-gateway.
func customErrorHandler(log *logger.Logger) runtime.ErrorHandlerFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		log.Error("HTTP request error",
			logger.String("method", r.Method),
			logger.String("path", r.URL.Path),
			logger.Any("error", err),
		)

		// Use default error handler
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
	}
}

// corsMiddleware adds CORS headers to responses.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs all HTTP requests.
func loggingMiddleware(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Debug("HTTP request received",
			logger.String("method", r.Method),
			logger.String("path", r.URL.Path),
			logger.String("remote_addr", r.RemoteAddr),
		)

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		log.Info("HTTP request completed",
			logger.String("method", r.Method),
			logger.String("path", r.URL.Path),
			logger.Int("status_code", wrapped.statusCode),
			logger.Duration("duration", duration),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
