package grpc

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v3"

	webprismv1 "github.com/ChangChunCheng/webprism/gen/go/v1"
	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/input"
)

// SpecServiceServer implements the gRPC SpecService server.
type SpecServiceServer struct {
	webprismv1.UnimplementedSpecServiceServer
	service input.SpecService
	logger  *logger.Logger
}

// NewSpecServiceServer creates a new SpecServiceServer.
func NewSpecServiceServer(service input.SpecService, log *logger.Logger) *SpecServiceServer {
	return &SpecServiceServer{
		service: service,
		logger:  log,
	}
}

// UploadSpec uploads a new OpenAPI specification.
func (s *SpecServiceServer) UploadSpec(ctx context.Context, req *webprismv1.UploadSpecRequest) (*webprismv1.UploadSpecResponse, error) {
	s.logger.Info("gRPC UploadSpec called",
		logger.String("name", req.Name),
		logger.String("format", req.Format.String()),
	)

	// Parse spec content (JSON or YAML) to map
	var specData map[string]interface{}
	var err error

	if req.Format == webprismv1.SpecFormat_SPEC_FORMAT_JSON {
		err = json.Unmarshal([]byte(req.SpecContent), &specData)
	} else {
		err = yaml.Unmarshal([]byte(req.SpecContent), &specData)
	}

	if err != nil {
		s.logger.Error("failed to parse spec content",
			logger.String("name", req.Name),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse spec content: %v", err)
	}

	// Extract version from spec data (OpenAPI info.version)
	version := ""
	if info, ok := specData["info"].(map[string]interface{}); ok {
		if v, ok := info["version"].(string); ok {
			version = v
		}
	}

	// Call domain service
	spec, err := s.service.UploadSpec(ctx, req.Name, version, specData)
	if err != nil {
		s.logger.Error("failed to upload spec",
			logger.String("name", req.Name),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to upload spec: %v", err)
	}

	// Convert to protobuf response
	return &webprismv1.UploadSpecResponse{
		Spec: s.toProtoSpec(spec),
	}, nil
}

// ListSpecs retrieves all API specifications with pagination.
func (s *SpecServiceServer) ListSpecs(ctx context.Context, req *webprismv1.ListSpecsRequest) (*webprismv1.ListSpecsResponse, error) {
	// Convert page-based pagination to limit/offset
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10 // default
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1 // default to first page
	}
	offset := (page - 1) * pageSize

	s.logger.Debug("gRPC ListSpecs called",
		logger.Int("page", page),
		logger.Int("page_size", pageSize),
	)

	// Call domain service
	specs, err := s.service.ListSpecs(ctx, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to list specs",
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to list specs: %v", err)
	}

	// Convert to protobuf response
	protoSpecs := make([]*webprismv1.APISpec, len(specs))
	for i, spec := range specs {
		protoSpecs[i] = s.toProtoSpec(spec)
	}

	return &webprismv1.ListSpecsResponse{
		Specs: protoSpecs,
	}, nil
}

// GetSpec retrieves a specific API specification by ID.
func (s *SpecServiceServer) GetSpec(ctx context.Context, req *webprismv1.GetSpecRequest) (*webprismv1.GetSpecResponse, error) {
	s.logger.Debug("gRPC GetSpec called",
		logger.String("spec_id", req.Id),
	)

	// Call domain service
	spec, err := s.service.GetSpec(ctx, req.Id)
	if err != nil {
		if err == model.ErrSpecNotFound {
			return nil, status.Errorf(codes.NotFound, "spec not found: %s", req.Id)
		}
		s.logger.Error("failed to get spec",
			logger.String("spec_id", req.Id),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to get spec: %v", err)
	}

	// Convert to protobuf response
	return &webprismv1.GetSpecResponse{
		Spec: s.toProtoSpec(spec),
	}, nil
}

// DeleteSpec removes an API specification by ID.
func (s *SpecServiceServer) DeleteSpec(ctx context.Context, req *webprismv1.DeleteSpecRequest) (*webprismv1.DeleteSpecResponse, error) {
	s.logger.Info("gRPC DeleteSpec called",
		logger.String("spec_id", req.Id),
	)

	// Call domain service
	if err := s.service.DeleteSpec(ctx, req.Id); err != nil {
		if err == model.ErrSpecNotFound {
			return nil, status.Errorf(codes.NotFound, "spec not found: %s", req.Id)
		}
		s.logger.Error("failed to delete spec",
			logger.String("spec_id", req.Id),
			logger.Any("error", err),
		)
		return nil, status.Errorf(codes.Internal, "failed to delete spec: %v", err)
	}

	return &webprismv1.DeleteSpecResponse{
		Success: true,
	}, nil
}

// toProtoSpec converts a domain APISpec to protobuf APISpec.
func (s *SpecServiceServer) toProtoSpec(spec *model.APISpec) *webprismv1.APISpec {
	return &webprismv1.APISpec{
		Id:        spec.ID,
		Name:      spec.Name,
		Version:   spec.Version,
		BaseUrl:   spec.BaseURL,
		CreatedAt: timestamppb.New(spec.CreatedAt),
		UpdatedAt: timestamppb.New(spec.UpdatedAt),
	}
}
