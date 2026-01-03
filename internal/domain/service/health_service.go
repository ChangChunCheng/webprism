package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/output"
)

// HealthService implements the HealthService input port.
type HealthService struct {
	specRepo   output.SpecRepository
	healthRepo output.HealthRepository
	httpClient output.HTTPClient
	logger     *logger.Logger
}

// NewHealthService creates a new HealthService instance.
func NewHealthService(
	specRepo output.SpecRepository,
	healthRepo output.HealthRepository,
	httpClient output.HTTPClient,
	log *logger.Logger,
) *HealthService {
	return &HealthService{
		specRepo:   specRepo,
		healthRepo: healthRepo,
		httpClient: httpClient,
		logger:     log,
	}
}

// CheckHealth performs a health check on the target API.
func (s *HealthService) CheckHealth(ctx context.Context, specID string) (*model.HealthCheck, error) {
	s.logger.Info("performing health check",
		logger.String("spec_id", specID),
	)

	// 1. Retrieve the API spec
	spec, err := s.specRepo.FindByID(ctx, specID)
	if err != nil {
		s.logger.Error("failed to retrieve spec for health check",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, err
	}

	// 2. Extract base URL
	checkURL := spec.BaseURL

	// 3. Perform health check
	status, responseTime, errorMsg := s.performHealthCheck(ctx, checkURL)

	// 4. Create health check record
	healthCheck := &model.HealthCheck{
		ID:             uuid.New().String(),
		SpecID:         specID,
		CheckURL:       checkURL,
		Status:         status,
		ResponseTimeMs: responseTime,
		CheckedAt:      time.Now(),
		ErrorMessage:   errorMsg,
	}

	// 5. Store health check record
	if err := s.healthRepo.Create(ctx, healthCheck); err != nil {
		s.logger.Error("failed to store health check result",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		// Continue even if storage fails
	}

	s.logger.Info("health check completed",
		logger.String("spec_id", specID),
		logger.String("status", string(status)),
		logger.Int("response_time_ms", responseTime),
	)

	return healthCheck, nil
}

// GetHealthHistory retrieves the health check history for a spec.
func (s *HealthService) GetHealthHistory(ctx context.Context, specID string, limit int) ([]*model.HealthCheck, error) {
	s.logger.Debug("getting health check history",
		logger.String("spec_id", specID),
		logger.Int("limit", limit),
	)

	// Set default limit if not specified
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	healthChecks, err := s.healthRepo.FindBySpecID(ctx, specID, limit, 0)
	if err != nil {
		s.logger.Error("failed to get health check history",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to get health check history: %w", err)
	}

	s.logger.Debug("health check history retrieved",
		logger.String("spec_id", specID),
		logger.Int("count", len(healthChecks)),
	)

	return healthChecks, nil
}

// GetLatestHealth retrieves the most recent health check for a spec.
func (s *HealthService) GetLatestHealth(ctx context.Context, specID string) (*model.HealthCheck, error) {
	s.logger.Debug("getting latest health check",
		logger.String("spec_id", specID),
	)

	healthCheck, err := s.healthRepo.FindLatestBySpecID(ctx, specID)
	if err != nil {
		s.logger.Error("failed to get latest health check",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to get latest health check: %w", err)
	}

	return healthCheck, nil
}

// performHealthCheck executes the actual health check request.
// Returns: status, responseTimeMs, errorMessage
func (s *HealthService) performHealthCheck(ctx context.Context, url string) (model.HealthStatus, int, string) {
	start := time.Now()

	// Create HEAD request (lightweight check)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		s.logger.Error("failed to create health check request",
			logger.String("url", url),
			logger.Any("error", err),
		)
		return model.HealthStatusDown, 0, fmt.Sprintf("Failed to create request: %v", err)
	}

	// Set timeout for health check
	checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Execute request
	resp, err := s.httpClient.Do(checkCtx, req)
	responseTime := int(time.Since(start).Milliseconds())

	if err != nil {
		s.logger.Warn("health check request failed",
			logger.String("url", url),
			logger.Int("response_time_ms", responseTime),
			logger.Any("error", err),
		)
		return model.HealthStatusDown, responseTime, fmt.Sprintf("Connection failed: %v", err)
	}
	defer resp.Body.Close()

	// Determine status based on HTTP status code
	var status model.HealthStatus
	var errorMsg string

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		status = model.HealthStatusUP
	case resp.StatusCode >= 500:
		status = model.HealthStatusDown
		errorMsg = fmt.Sprintf("Server error: HTTP %d", resp.StatusCode)
	default:
		status = model.HealthStatusDegraded
		errorMsg = fmt.Sprintf("Unexpected status: HTTP %d", resp.StatusCode)
	}

	s.logger.Debug("health check response received",
		logger.String("url", url),
		logger.Int("status_code", resp.StatusCode),
		logger.Int("response_time_ms", responseTime),
	)

	return status, responseTime, errorMsg
}
