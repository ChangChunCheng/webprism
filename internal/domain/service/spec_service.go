package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/output"
)

// SpecService implements the SpecService input port.
type SpecService struct {
	specRepo output.SpecRepository
	parser   output.OpenAPIParser
	logger   *logger.Logger
}

// NewSpecService creates a new SpecService instance.
func NewSpecService(
	specRepo output.SpecRepository,
	parser output.OpenAPIParser,
	log *logger.Logger,
) *SpecService {
	return &SpecService{
		specRepo: specRepo,
		parser:   parser,
		logger:   log,
	}
}

// UploadSpec uploads a new OpenAPI specification.
func (s *SpecService) UploadSpec(ctx context.Context, name, version string, specData map[string]interface{}) (*model.APISpec, error) {
	s.logger.Info("uploading API spec",
		logger.String("name", name),
		logger.String("version", version),
	)

	// Validate the spec
	if err := s.parser.Validate(ctx, specData); err != nil {
		s.logger.Error("invalid OpenAPI spec",
			logger.String("name", name),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("invalid OpenAPI spec: %w", err)
	}

	// Parse the spec
	parsed, err := s.parser.Parse(ctx, specData)
	if err != nil {
		s.logger.Error("failed to parse OpenAPI spec",
			logger.String("name", name),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}

	// Generate UUID for the spec (application-generated)
	specID := uuid.New().String()

	// Create APISpec domain model
	now := time.Now()
	spec := &model.APISpec{
		ID:        specID,
		Name:      name,
		Version:   version,
		BaseURL:   parsed.BaseURL,
		SpecData:  specData,
		Paths:     parsed.Paths,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Store the spec
	if err := s.specRepo.Create(ctx, spec); err != nil {
		s.logger.Error("failed to store spec",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to store spec: %w", err)
	}

	s.logger.Info("API spec uploaded successfully",
		logger.String("spec_id", specID),
		logger.String("name", name),
		logger.String("base_url", parsed.BaseURL),
	)

	return spec, nil
}

// ListSpecs retrieves all API specifications with pagination.
func (s *SpecService) ListSpecs(ctx context.Context, limit, offset int) ([]*model.APISpec, error) {
	s.logger.Debug("listing API specs",
		logger.Int("limit", limit),
		logger.Int("offset", offset),
	)

	// Set default limit if not specified
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	if offset < 0 {
		offset = 0
	}

	specs, err := s.specRepo.FindAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list specs",
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to list specs: %w", err)
	}

	s.logger.Debug("specs listed successfully",
		logger.Int("count", len(specs)),
	)

	return specs, nil
}

// GetSpec retrieves a specific API specification by ID.
func (s *SpecService) GetSpec(ctx context.Context, specID string) (*model.APISpec, error) {
	s.logger.Debug("getting API spec",
		logger.String("spec_id", specID),
	)

	// Retrieve the spec
	spec, err := s.specRepo.FindByID(ctx, specID)
	if err != nil {
		s.logger.Error("failed to get spec",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, err
	}

	// Re-parse SpecData to populate Paths (since Paths is not persisted)
	if spec.SpecData != nil && spec.Paths == nil {
		parsed, err := s.parser.Parse(ctx, spec.SpecData)
		if err != nil {
			s.logger.Error("failed to parse spec data",
				logger.String("spec_id", specID),
				logger.Any("error", err),
			)
			// Don't fail the entire request, just log the error
			// The spec data is still available even if parsing fails
		} else {
			spec.Paths = parsed.Paths
		}
	}

	s.logger.Debug("spec retrieved successfully",
		logger.String("spec_id", specID),
		logger.String("name", spec.Name),
	)

	return spec, nil
}

// DeleteSpec removes an API specification by ID.
func (s *SpecService) DeleteSpec(ctx context.Context, specID string) error {
	s.logger.Info("deleting API spec",
		logger.String("spec_id", specID),
	)

	// Delete the spec (cascade will handle auth configs and health checks)
	if err := s.specRepo.Delete(ctx, specID); err != nil {
		s.logger.Error("failed to delete spec",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return err
	}

	s.logger.Info("spec deleted successfully",
		logger.String("spec_id", specID),
	)

	return nil
}
