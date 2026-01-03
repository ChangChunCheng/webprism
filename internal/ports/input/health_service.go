package input

import (
	"context"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
)

// HealthService defines the input port for health check operations.
// This interface handles checking target API health and retrieving health history.
type HealthService interface {
	// CheckHealth performs a health check on the target API.
	// It performs the following:
	// 1. Retrieves the API spec by specID
	// 2. Extracts the base URL from the spec
	// 3. Sends a HEAD or GET request to the health endpoint (or base URL)
	// 4. Records response time, status, and any errors
	// 5. Stores the health check result in the database
	//
	// Returns the HealthCheck result or an error.
	CheckHealth(ctx context.Context, specID string) (*model.HealthCheck, error)

	// GetHealthHistory retrieves the health check history for a spec.
	// Returns the most recent health checks, ordered by checked_at DESC.
	// The limit parameter controls how many records to return (default 10).
	//
	// Returns a list of HealthCheck records or an error.
	GetHealthHistory(ctx context.Context, specID string, limit int) ([]*model.HealthCheck, error)

	// GetLatestHealth retrieves the most recent health check for a spec.
	// Returns the latest HealthCheck or nil if no checks exist.
	GetLatestHealth(ctx context.Context, specID string) (*model.HealthCheck, error)
}
