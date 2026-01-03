package input

import (
	"context"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
)

// AuthService defines the input port for authentication configuration management.
// This interface handles setting and retrieving auth configs for API specs.
type AuthService interface {
	// SetAuthConfig sets or updates the authentication configuration for a spec.
	// The credentials will be encrypted before storage.
	// Returns the created AuthConfig or an error (e.g., ErrSpecNotFound if spec doesn't exist).
	SetAuthConfig(ctx context.Context, specID string, authType model.AuthType, credentials map[string]string) (*model.AuthConfig, error)

	// GetAuthConfig retrieves the authentication configuration for a spec.
	// The credentials will be decrypted before returning.
	// Returns the AuthConfig with decrypted credentials, or ErrAuthConfigNotFound.
	GetAuthConfig(ctx context.Context, specID string) (*model.AuthConfig, error)

	// DeleteAuthConfig removes the authentication configuration for a spec.
	// Returns ErrAuthConfigNotFound if the config doesn't exist.
	DeleteAuthConfig(ctx context.Context, specID string) error
}
