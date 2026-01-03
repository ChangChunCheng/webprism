package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/uptrace/bun"
)

// AuthRepository implements the output port for auth configuration persistence using PostgreSQL.
type AuthRepository struct {
	db     *bun.DB
	logger *logger.Logger
}

// NewAuthRepository creates a new PostgreSQL-backed AuthRepository.
func NewAuthRepository(db *bun.DB, log *logger.Logger) *AuthRepository {
	return &AuthRepository{
		db:     db,
		logger: log,
	}
}

// Create stores a new authentication configuration.
func (r *AuthRepository) Create(ctx context.Context, auth *model.AuthConfig) error {
	_, err := r.db.NewInsert().
		Model(auth).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to create auth config",
			logger.String("spec_id", auth.SpecID),
			logger.Any("error", err),
		)
		return fmt.Errorf("failed to create auth config: %w", err)
	}

	r.logger.Info("auth config created successfully",
		logger.String("spec_id", auth.SpecID),
		logger.String("auth_type", string(auth.AuthType)),
	)

	return nil
}

// FindBySpecID retrieves the authentication configuration for a spec.
func (r *AuthRepository) FindBySpecID(ctx context.Context, specID string) (*model.AuthConfig, error) {
	auth := &model.AuthConfig{}

	err := r.db.NewSelect().
		Model(auth).
		Where("spec_id = ?", specID).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrAuthConfigNotFound
		}
		r.logger.Error("failed to find auth config by spec ID",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to find auth config: %w", err)
	}

	return auth, nil
}

// Update updates an existing authentication configuration.
func (r *AuthRepository) Update(ctx context.Context, auth *model.AuthConfig) error {
	result, err := r.db.NewUpdate().
		Model(auth).
		Column("auth_type", "encrypted_credentials", "updated_at").
		Where("spec_id = ?", auth.SpecID).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to update auth config",
			logger.String("spec_id", auth.SpecID),
			logger.Any("error", err),
		)
		return fmt.Errorf("failed to update auth config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrAuthConfigNotFound
	}

	r.logger.Info("auth config updated successfully",
		logger.String("spec_id", auth.SpecID),
	)

	return nil
}

// Delete removes the authentication configuration for a spec.
func (r *AuthRepository) Delete(ctx context.Context, specID string) error {
	result, err := r.db.NewDelete().
		Model((*model.AuthConfig)(nil)).
		Where("spec_id = ?", specID).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to delete auth config",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return fmt.Errorf("failed to delete auth config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrAuthConfigNotFound
	}

	r.logger.Info("auth config deleted successfully",
		logger.String("spec_id", specID),
	)

	return nil
}

// Exists checks if an auth config exists for a spec.
func (r *AuthRepository) Exists(ctx context.Context, specID string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*model.AuthConfig)(nil)).
		Where("spec_id = ?", specID).
		Exists(ctx)

	if err != nil {
		r.logger.Error("failed to check auth config existence",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return false, fmt.Errorf("failed to check auth config existence: %w", err)
	}

	return exists, nil
}
