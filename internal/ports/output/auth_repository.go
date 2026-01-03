package output

import (
	"context"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
)

// AuthRepository defines the output port for authentication configuration persistence.
// This interface is implemented by storage adapters (PostgreSQL, Memory, etc.).
type AuthRepository interface {
	// Create stores a new authentication configuration.
	// The credentials should already be encrypted by the caller.
	// Returns an error if the auth config already exists or if storage fails.
	Create(ctx context.Context, auth *model.AuthConfig) error

	// FindBySpecID retrieves the authentication configuration for a spec.
	// Returns the AuthConfig or ErrAuthConfigNotFound if it doesn't exist.
	FindBySpecID(ctx context.Context, specID string) (*model.AuthConfig, error)

	// Update updates an existing authentication configuration.
	// The credentials should already be encrypted by the caller.
	// Returns ErrAuthConfigNotFound if the auth config doesn't exist.
	Update(ctx context.Context, auth *model.AuthConfig) error

	// Delete removes the authentication configuration for a spec.
	// Returns ErrAuthConfigNotFound if the auth config doesn't exist.
	Delete(ctx context.Context, specID string) error

	// Exists checks if an auth config exists for a spec.
	// Returns true if exists, false otherwise.
	Exists(ctx context.Context, specID string) (bool, error)
}
