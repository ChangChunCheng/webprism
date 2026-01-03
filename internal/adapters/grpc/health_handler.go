package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	webprismv1 "github.com/ChangChunCheng/webprism/gen/go/v1"
	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/input"
)

// HealthServiceServer implements the gRPC HealthService server.
type HealthServiceServer struct {
	webprismv1.UnimplementedHealthServiceServer
	service input.HealthService
	logger  *logger.Logger
}

// NewHealthServiceServer creates a new HealthServiceServer.
func NewHealthServiceServer(service input.HealthService, log *logger.Logger) *HealthServiceServer {
	return &HealthServiceServer{
		service: service,
		logger:  log,
	}
}

// CheckHealth performs a health check on the target API.
func (s *HealthServiceServer) CheckHealth(ctx context.Context, req *webprismv1.CheckHealthRequest) (*webprismv1.CheckHealthResponse, error) {
	s.logger.Info("gRPC CheckHealth called",
		logger.String("spec_id", req.SpecId),
	)

	// Call domain service
	healthCheck, err := s.service.CheckHealth(ctx, req.SpecId)
	if err != nil {
		if err == model.ErrSpecNotFound {
			return nil, status.Errorf(codes.NotFound, "spec not found: %s", req.SpecId)
		}
		s.logger.Error("failed to check health",
			logger.String("spec_id", req.SpecId),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to check health: %v", err)
	}

	// Convert to protobuf response
	return &webprismv1.CheckHealthResponse{
		Result: s.toProtoHealthCheck(healthCheck),
	}, nil
}

// GetHealthHistory retrieves the health check history for a spec.
func (s *HealthServiceServer) GetHealthHistory(ctx context.Context, req *webprismv1.GetHealthHistoryRequest) (*webprismv1.GetHealthHistoryResponse, error) {
	s.logger.Debug("gRPC GetHealthHistory called",
		logger.String("spec_id", req.SpecId),
		logger.Int("limit", int(req.Limit)),
	)

	// Call domain service
	history, err := s.service.GetHealthHistory(ctx, req.SpecId, int(req.Limit))
	if err != nil {
		s.logger.Error("failed to get health history",
			logger.String("spec_id", req.SpecId),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to get health history: %v", err)
	}

	// Convert to protobuf response
	protoHistory := make([]*webprismv1.HealthCheckResult, len(history))
	for i, check := range history {
		protoHistory[i] = s.toProtoHealthCheck(check)
	}

	return &webprismv1.GetHealthHistoryResponse{
		Results: protoHistory,
	}, nil
}

// toProtoHealthCheck converts a domain HealthCheck to protobuf HealthCheckResult.
func (s *HealthServiceServer) toProtoHealthCheck(healthCheck *model.HealthCheck) *webprismv1.HealthCheckResult {
	return &webprismv1.HealthCheckResult{
		Id:             healthCheck.ID,
		SpecId:         healthCheck.SpecID,
		CheckUrl:       healthCheck.CheckURL,
		Status:         s.toProtoHealthStatus(healthCheck.Status),
		ResponseTimeMs: int32(healthCheck.ResponseTimeMs),
		CheckedAt:      timestamppb.New(healthCheck.CheckedAt),
		ErrorMessage:   healthCheck.ErrorMessage,
	}
}

// toProtoHealthStatus converts domain HealthStatus to protobuf HealthStatus.
func (s *HealthServiceServer) toProtoHealthStatus(domainStatus model.HealthStatus) webprismv1.HealthStatus {
	switch domainStatus {
	case model.HealthStatusUP:
		return webprismv1.HealthStatus_HEALTH_STATUS_UP
	case model.HealthStatusDown:
		return webprismv1.HealthStatus_HEALTH_STATUS_DOWN
	case model.HealthStatusDegraded:
		return webprismv1.HealthStatus_HEALTH_STATUS_DEGRADED
	default:
		return webprismv1.HealthStatus_HEALTH_STATUS_DOWN
	}
}
