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

// AuthServiceServer implements the gRPC AuthService server.
type AuthServiceServer struct {
	webprismv1.UnimplementedAuthServiceServer
	service input.AuthService
	logger  *logger.Logger
}

// NewAuthServiceServer creates a new AuthServiceServer.
func NewAuthServiceServer(service input.AuthService, log *logger.Logger) *AuthServiceServer {
	return &AuthServiceServer{
		service: service,
		logger:  log,
	}
}

// SetAuthConfig sets or updates the authentication configuration for a spec.
func (s *AuthServiceServer) SetAuthConfig(ctx context.Context, req *webprismv1.SetAuthConfigRequest) (*webprismv1.SetAuthConfigResponse, error) {
	s.logger.Info("gRPC SetAuthConfig called",
		logger.String("spec_id", req.SpecId),
		logger.String("auth_type", req.AuthType.String()),
	)

	// Convert protobuf auth type to domain auth type
	authType := s.toDomainAuthType(req.AuthType)

	// Convert credentials
	credentials := req.Credentials

	// Call domain service
	_, err := s.service.SetAuthConfig(ctx, req.SpecId, authType, credentials)
	if err != nil {
		if err == model.ErrSpecNotFound {
			return nil, status.Errorf(codes.NotFound, "spec not found: %s", req.SpecId)
		}
		s.logger.Error("failed to set auth config",
			logger.String("spec_id", req.SpecId),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to set auth config: %v", err)
	}

	// Convert to protobuf response
	return &webprismv1.SetAuthConfigResponse{
		Success: true,
		Message: "auth config set successfully",
	}, nil
}

// GetAuthConfig retrieves the authentication configuration for a spec.
func (s *AuthServiceServer) GetAuthConfig(ctx context.Context, req *webprismv1.GetAuthConfigRequest) (*webprismv1.GetAuthConfigResponse, error) {
	s.logger.Debug("gRPC GetAuthConfig called",
		logger.String("spec_id", req.SpecId),
	)

	// Call domain service
	authConfig, err := s.service.GetAuthConfig(ctx, req.SpecId)
	if err != nil {
		if err == model.ErrAuthConfigNotFound {
			return nil, status.Errorf(codes.NotFound, "auth config not found for spec: %s", req.SpecId)
		}
		s.logger.Error("failed to get auth config",
			logger.String("spec_id", req.SpecId),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to get auth config: %v", err)
	}

	// Convert to protobuf response
	return &webprismv1.GetAuthConfigResponse{
		SpecId:    authConfig.SpecID,
		AuthType:  s.toProtoAuthType(authConfig.AuthType),
		Position:  webprismv1.AuthPosition_AUTH_POSITION_HEADER, // Default position
		CreatedAt: timestamppb.New(authConfig.CreatedAt),
		UpdatedAt: timestamppb.New(authConfig.UpdatedAt),
	}, nil
}

// toDomainAuthType converts protobuf AuthType to domain AuthType.
func (s *AuthServiceServer) toDomainAuthType(protoType webprismv1.AuthType) model.AuthType {
	switch protoType {
	case webprismv1.AuthType_AUTH_TYPE_API_KEY:
		return model.AuthTypeAPIKey
	case webprismv1.AuthType_AUTH_TYPE_BEARER:
		return model.AuthTypeBearer
	case webprismv1.AuthType_AUTH_TYPE_BASIC:
		return model.AuthTypeBasic
	case webprismv1.AuthType_AUTH_TYPE_OAUTH2:
		return model.AuthTypeOAuth2
	default:
		return model.AuthTypeAPIKey
	}
}

// toProtoAuthType converts domain AuthType to protobuf AuthType.
func (s *AuthServiceServer) toProtoAuthType(domainType model.AuthType) webprismv1.AuthType {
	switch domainType {
	case model.AuthTypeAPIKey:
		return webprismv1.AuthType_AUTH_TYPE_API_KEY
	case model.AuthTypeBearer:
		return webprismv1.AuthType_AUTH_TYPE_BEARER
	case model.AuthTypeBasic:
		return webprismv1.AuthType_AUTH_TYPE_BASIC
	case model.AuthTypeOAuth2:
		return webprismv1.AuthType_AUTH_TYPE_OAUTH2
	default:
		return webprismv1.AuthType_AUTH_TYPE_API_KEY
	}
}
