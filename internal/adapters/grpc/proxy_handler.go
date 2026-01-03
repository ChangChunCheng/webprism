package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	webprismv1 "github.com/ChangChunCheng/webprism/gen/go/v1"
	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/input"
)

// ProxyServiceServer implements the gRPC ProxyService server.
type ProxyServiceServer struct {
	webprismv1.UnimplementedProxyServiceServer
	service input.ProxyService
	logger  *logger.Logger
}

// NewProxyServiceServer creates a new ProxyServiceServer.
func NewProxyServiceServer(service input.ProxyService, log *logger.Logger) *ProxyServiceServer {
	return &ProxyServiceServer{
		service: service,
		logger:  log,
	}
}

// ExecuteProxy executes a proxied API request to the target API.
func (s *ProxyServiceServer) ExecuteProxy(ctx context.Context, req *webprismv1.ExecuteProxyRequest) (*webprismv1.ExecuteProxyResponse, error) {
	s.logger.Info("gRPC ExecuteProxy called",
		logger.String("spec_id", req.SpecId),
		logger.String("operation_id", req.OperationId),
	)

	// Convert protobuf request to domain ProxyRequest
	// Note: protobuf uses simple map<string,string> for parameters
	// We'll parse query params from parameters map (prefixed with "query_")
	// and headers from parameters map (prefixed with "header_")
	bodyMap := make(map[string]interface{})
	if req.Body != nil {
		bodyMap = req.Body.AsMap()
	}

	domainReq := &model.ProxyRequest{
		SpecID:      req.SpecId,
		OperationID: req.OperationId,
		Parameters: model.ProxyParameters{
			Query:   req.Parameters, // Simple pass-through for now
			Headers: make(map[string]string),
			Body:    bodyMap,
		},
	}

	// Validate request
	if err := domainReq.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Call domain service
	resp, err := s.service.ExecuteProxy(ctx, domainReq)
	if err != nil {
		// Check error type for proper gRPC status code
		if err == model.ErrSpecNotFound {
			return nil, status.Errorf(codes.NotFound, "spec not found: %s", req.SpecId)
		}
		if err == model.ErrOperationNotFound {
			return nil, status.Errorf(codes.NotFound, "operation not found: %s", req.OperationId)
		}
		s.logger.Error("failed to execute proxy",
			logger.String("spec_id", req.SpecId),
			logger.String("operation_id", req.OperationId),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to execute proxy: %v", err)
	}

	// Convert body to protobuf Struct
	bodyStruct, err := structpb.NewStruct(resp.Body)
	if err != nil {
		s.logger.Error("failed to convert response body to proto struct",
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to convert response body: %v", err)
	}

	// Build protobuf response
	protoResp := &webprismv1.ExecuteProxyResponse{
		StatusCode: int32(resp.StatusCode),
		Headers:    resp.Headers,
		Body:       bodyStruct,
	}

	// Add error if present
	if resp.Error != nil {
		protoResp.Error = &webprismv1.ProxyError{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
			Type:    s.toProtoErrorType(resp.Error.Type),
		}

		// Convert error details to struct
		if len(resp.Error.Details) > 0 {
			detailsStruct, err := structpb.NewStruct(resp.Error.Details)
			if err == nil {
				protoResp.Error.Details = detailsStruct
			}
		}
	}

	return protoResp, nil
}

// toProtoErrorType converts domain ErrorType to protobuf ErrorType.
func (s *ProxyServiceServer) toProtoErrorType(domainType model.ErrorType) webprismv1.ErrorType {
	switch domainType {
	case model.ErrorTypeSystem:
		return webprismv1.ErrorType_ERROR_TYPE_SYSTEM
	case model.ErrorTypeClient:
		return webprismv1.ErrorType_ERROR_TYPE_CLIENT
	case model.ErrorTypeExternal:
		return webprismv1.ErrorType_ERROR_TYPE_EXTERNAL
	default:
		return webprismv1.ErrorType_ERROR_TYPE_SYSTEM
	}
}
