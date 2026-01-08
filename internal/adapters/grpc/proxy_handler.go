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
	service    input.ProxyService
	specRepo   input.SpecService // Need to access spec to categorize parameters
	logger     *logger.Logger
}

// NewProxyServiceServer creates a new ProxyServiceServer.
func NewProxyServiceServer(service input.ProxyService, specService input.SpecService, log *logger.Logger) *ProxyServiceServer {
	return &ProxyServiceServer{
		service:  service,
		specRepo: specService,
		logger:   log,
	}
}

// ExecuteProxy executes a proxied API request to the target API.
func (s *ProxyServiceServer) ExecuteProxy(ctx context.Context, req *webprismv1.ExecuteProxyRequest) (*webprismv1.ExecuteProxyResponse, error) {
	s.logger.Info("gRPC ExecuteProxy called",
		logger.String("spec_id", req.SpecId),
		logger.String("operation_id", req.OperationId),
	)

	// Convert protobuf request to domain ProxyRequest
	bodyMap := make(map[string]interface{})
	if req.Body != nil {
		bodyMap = req.Body.AsMap()
	}

	// Categorize parameters based on OpenAPI spec
	pathParams, queryParams, err := s.categorizeParameters(ctx, req.SpecId, req.OperationId, req.Parameters)
	if err != nil {
		s.logger.Warn("failed to categorize parameters, treating all as query params",
			logger.Any("error", err),
		)
		// Fallback: treat all as query params
		pathParams = make(map[string]string)
		queryParams = req.Parameters
	}

	domainReq := &model.ProxyRequest{
		SpecID:      req.SpecId,
		OperationID: req.OperationId,
		Parameters: model.ProxyParameters{
			Path:    pathParams,
			Query:   queryParams,
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

// categorizeParameters categorizes parameters into path and query parameters
// based on the OpenAPI specification.
func (s *ProxyServiceServer) categorizeParameters(ctx context.Context, specID, operationID string, params map[string]string) (path, query map[string]string, err error) {
	path = make(map[string]string)
	query = make(map[string]string)

	// Get the spec to understand parameter types
	spec, err := s.specRepo.GetSpec(ctx, specID)
	if err != nil {
		return nil, nil, err
	}

	// Find the operation to get parameter definitions
	var operation *model.Operation
	for _, pathItem := range spec.Paths {
		if pathItem.Operations != nil {
			for _, op := range pathItem.Operations {
				if op.OperationID == operationID {
					operation = op
					break
				}
			}
		}
		if operation != nil {
			break
		}
	}

	if operation == nil {
		// If we can't find the operation, treat all as query params
		return path, params, nil
	}

	// Categorize parameters based on their definition in the operation
	for paramName, paramValue := range params {
		// Check if this parameter is defined in the operation
		isPathParam := false
		if operation.Parameters != nil {
			for _, param := range operation.Parameters {
				if param.Name == paramName && param.In == "path" {
					isPathParam = true
					break
				}
			}
		}

		if isPathParam {
			path[paramName] = paramValue
		} else {
			query[paramName] = paramValue
		}
	}

	return path, query, nil
}
