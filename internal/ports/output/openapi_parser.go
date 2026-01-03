package output

import (
	"context"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
)

// OpenAPIParser defines the output port for parsing OpenAPI specifications.
// This interface abstracts the OpenAPI parsing logic to support multiple versions (2.0, 3.0).
type OpenAPIParser interface {
	// Parse parses an OpenAPI specification and extracts structured information.
	// It supports both OpenAPI 2.0 (Swagger) and 3.0+ specifications.
	// For 2.0 specs, it automatically converts them to 3.0 internally.
	//
	// Returns:
	// - Parsed paths with operations, parameters, and schemas
	// - Base URL(s) extracted from the spec
	// - Version of the spec
	// - Any parsing errors
	Parse(ctx context.Context, specData map[string]interface{}) (*ParsedSpec, error)

	// Validate validates an OpenAPI specification for correctness.
	// Returns validation errors if the spec is invalid.
	Validate(ctx context.Context, specData map[string]interface{}) error

	// ExtractBaseURL extracts the base URL from the spec.
	// For OpenAPI 3.0, it uses the first server URL.
	// For Swagger 2.0, it constructs from schemes, host, and basePath.
	ExtractBaseURL(specData map[string]interface{}) (string, error)
}

// ParsedSpec represents a parsed OpenAPI specification.
type ParsedSpec struct {
	// Version is the OpenAPI version (e.g., "3.0.0", "2.0")
	Version string

	// BaseURL is the primary base URL for the API
	BaseURL string

	// Paths contains all API paths and their operations
	Paths map[string]*model.PathItem

	// Info contains API metadata (title, description, version)
	Info *SpecInfo
}

// SpecInfo contains metadata about the API specification.
type SpecInfo struct {
	Title       string
	Description string
	Version     string
}
