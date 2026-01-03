package input

import (
	"context"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
)

// SpecService defines the input port for API Spec management operations.
// This interface represents the use cases available to external adapters (HTTP, gRPC, CLI, MCP).
type SpecService interface {
	// UploadSpec uploads a new OpenAPI specification.
	// It parses the spec, validates it, generates a UUID, and stores it.
	// Returns the created APISpec or an error.
	UploadSpec(ctx context.Context, name, version string, specData map[string]interface{}) (*model.APISpec, error)

	// ListSpecs retrieves all API specifications with pagination.
	// Returns a list of APISpecs or an error.
	ListSpecs(ctx context.Context, limit, offset int) ([]*model.APISpec, error)

	// GetSpec retrieves a specific API specification by ID.
	// Returns the APISpec with full details including parsed paths, or ErrSpecNotFound.
	GetSpec(ctx context.Context, specID string) (*model.APISpec, error)

	// DeleteSpec removes an API specification by ID.
	// This will cascade delete auth configs and health check records.
	// Returns ErrSpecNotFound if the spec doesn't exist.
	DeleteSpec(ctx context.Context, specID string) error
}
