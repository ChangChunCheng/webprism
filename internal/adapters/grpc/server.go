package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authpb "github.com/ChangChunCheng/webprism/gen/go/v1"
	healthpb "github.com/ChangChunCheng/webprism/gen/go/v1"
	proxypb "github.com/ChangChunCheng/webprism/gen/go/v1"
	specpb "github.com/ChangChunCheng/webprism/gen/go/v1"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/input"
)

// Server wraps the gRPC server with all registered services.
type Server struct {
	grpcServer *grpc.Server
	port       int
	logger     *logger.Logger
}

// NewServer creates a new gRPC server with all services registered.
func NewServer(
	port int,
	specService input.SpecService,
	authService input.AuthService,
	proxyService input.ProxyService,
	healthService input.HealthService,
	log *logger.Logger,
) *Server {
	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(log),
			errorHandlingInterceptor(log),
		),
	)

	// Register services
	specpb.RegisterSpecServiceServer(grpcServer, NewSpecServiceServer(specService, log))
	authpb.RegisterAuthServiceServer(grpcServer, NewAuthServiceServer(authService, log))
	proxypb.RegisterProxyServiceServer(grpcServer, NewProxyServiceServer(proxyService, specService, log))
	healthpb.RegisterHealthServiceServer(grpcServer, NewHealthServiceServer(healthService, log))

	// Enable reflection for grpcurl and other tools
	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		port:       port,
		logger:     log,
	}
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.port, err)
	}

	s.logger.Info("gRPC server starting",
		logger.Int("port", s.port),
	)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("gRPC server failed: %w", err)
	}

	return nil
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("gRPC server stopping...")

	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.logger.Warn("gRPC server force stopping due to context cancellation")
		s.grpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	}
}

// loggingInterceptor logs all gRPC requests.
func loggingInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Debug("gRPC request received",
			logger.String("method", info.FullMethod),
		)

		resp, err := handler(ctx, req)

		if err != nil {
			log.Error("gRPC request failed",
				logger.String("method", info.FullMethod),
				logger.Any("error", err),
			)
		} else {
			log.Debug("gRPC request completed",
				logger.String("method", info.FullMethod),
			)
		}

		return resp, err
	}
}

// errorHandlingInterceptor handles errors and converts them to appropriate gRPC status codes.
func errorHandlingInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		// Error is already handled by individual handlers
		// This interceptor is a placeholder for future error handling logic

		return resp, err
	}
}
