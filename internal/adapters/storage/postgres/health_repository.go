package postgres

import (
	"context"
	"fmt"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/uptrace/bun"
)

// HealthRepository implements the output port for health check persistence using PostgreSQL.
type HealthRepository struct {
	db     *bun.DB
	logger *logger.Logger
}

// NewHealthRepository creates a new PostgreSQL-backed HealthRepository.
func NewHealthRepository(db *bun.DB, log *logger.Logger) *HealthRepository {
	return &HealthRepository{
		db:     db,
		logger: log,
	}
}

// Create stores a new health check record.
func (r *HealthRepository) Create(ctx context.Context, health *model.HealthCheck) error {
	_, err := r.db.NewInsert().
		Model(health).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to create health check",
			logger.String("spec_id", health.SpecID),
			logger.Any("error", err),
		)
		return fmt.Errorf("failed to create health check: %w", err)
	}

	r.logger.Debug("health check created successfully",
		logger.String("health_check_id", health.ID),
		logger.String("spec_id", health.SpecID),
		logger.String("status", string(health.Status)),
	)

	return nil
}

// FindBySpecID retrieves health check records for a spec with pagination.
func (r *HealthRepository) FindBySpecID(ctx context.Context, specID string, limit, offset int) ([]*model.HealthCheck, error) {
	var healthChecks []*model.HealthCheck

	err := r.db.NewSelect().
		Model(&healthChecks).
		Where("spec_id = ?", specID).
		Order("checked_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		r.logger.Error("failed to find health checks by spec ID",
			logger.String("spec_id", specID),
			logger.Int("limit", limit),
			logger.Int("offset", offset),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to find health checks: %w", err)
	}

	return healthChecks, nil
}

// FindLatestBySpecID retrieves the most recent health check for a spec.
func (r *HealthRepository) FindLatestBySpecID(ctx context.Context, specID string) (*model.HealthCheck, error) {
	health := &model.HealthCheck{}

	err := r.db.NewSelect().
		Model(health).
		Where("spec_id = ?", specID).
		Order("checked_at DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		r.logger.Error("failed to find latest health check",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to find latest health check: %w", err)
	}

	return health, nil
}

// DeleteBySpecID removes all health check records for a spec.
func (r *HealthRepository) DeleteBySpecID(ctx context.Context, specID string) (int, error) {
	result, err := r.db.NewDelete().
		Model((*model.HealthCheck)(nil)).
		Where("spec_id = ?", specID).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to delete health checks",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return 0, fmt.Errorf("failed to delete health checks: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	r.logger.Info("health checks deleted successfully",
		logger.String("spec_id", specID),
		logger.Int64("count", rowsAffected),
	)

	return int(rowsAffected), nil
}
