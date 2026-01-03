package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/uptrace/bun"
)

// SpecRepository implements the output port for API Spec persistence using PostgreSQL.
type SpecRepository struct {
	db     *bun.DB
	logger *logger.Logger
}

// NewSpecRepository creates a new PostgreSQL-backed SpecRepository.
func NewSpecRepository(db *bun.DB, log *logger.Logger) *SpecRepository {
	return &SpecRepository{
		db:     db,
		logger: log,
	}
}

// Create stores a new API specification.
func (r *SpecRepository) Create(ctx context.Context, spec *model.APISpec) error {
	_, err := r.db.NewInsert().
		Model(spec).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to create spec",
			logger.String("spec_id", spec.ID),
			logger.Any("error", err),
		)
		return fmt.Errorf("failed to create spec: %w", err)
	}

	r.logger.Info("spec created successfully",
		logger.String("spec_id", spec.ID),
		logger.String("name", spec.Name),
	)

	return nil
}

// FindByID retrieves an API specification by its ID.
func (r *SpecRepository) FindByID(ctx context.Context, id string) (*model.APISpec, error) {
	spec := &model.APISpec{}

	err := r.db.NewSelect().
		Model(spec).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrSpecNotFound
		}
		r.logger.Error("failed to find spec by ID",
			logger.String("spec_id", id),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to find spec: %w", err)
	}

	return spec, nil
}

// FindAll retrieves all API specifications with pagination.
// Note: Excludes spec_data to optimize memory usage in list views
func (r *SpecRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.APISpec, error) {
	var specs []*model.APISpec

	err := r.db.NewSelect().
		Model(&specs).
		ExcludeColumn("spec_data"). // Exclude large JSONB field for performance
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		r.logger.Error("failed to find all specs",
			logger.Int("limit", limit),
			logger.Int("offset", offset),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to find specs: %w", err)
	}

	return specs, nil
}

// Update updates an existing API specification.
func (r *SpecRepository) Update(ctx context.Context, spec *model.APISpec) error {
	result, err := r.db.NewUpdate().
		Model(spec).
		Column("name", "version", "base_url", "spec_data", "updated_at").
		Where("id = ?", spec.ID).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to update spec",
			logger.String("spec_id", spec.ID),
			logger.Any("error", err),
		)
		return fmt.Errorf("failed to update spec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrSpecNotFound
	}

	r.logger.Info("spec updated successfully",
		logger.String("spec_id", spec.ID),
	)

	return nil
}

// Delete removes an API specification by ID.
func (r *SpecRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.NewDelete().
		Model((*model.APISpec)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		r.logger.Error("failed to delete spec",
			logger.String("spec_id", id),
			logger.Any("error", err),
		)
		return fmt.Errorf("failed to delete spec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrSpecNotFound
	}

	r.logger.Info("spec deleted successfully",
		logger.String("spec_id", id),
	)

	return nil
}

// Exists checks if a spec exists by ID.
func (r *SpecRepository) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*model.APISpec)(nil)).
		Where("id = ?", id).
		Exists(ctx)

	if err != nil {
		r.logger.Error("failed to check spec existence",
			logger.String("spec_id", id),
			logger.Any("error", err),
		)
		return false, fmt.Errorf("failed to check spec existence: %w", err)
	}

	return exists, nil
}
