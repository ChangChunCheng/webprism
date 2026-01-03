package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/mocks"
	"github.com/ChangChunCheng/webprism/internal/testutil"
)

// setupHealthServiceTest creates test environment with all mocks
func setupHealthServiceTest(t *testing.T) (*HealthService, *mocks.MockSpecRepository, *mocks.MockHealthRepository, *mocks.MockHTTPClient) {
	mockSpecRepo := mocks.NewMockSpecRepository(t)
	mockHealthRepo := mocks.NewMockHealthRepository(t)
	mockHTTPClient := mocks.NewMockHTTPClient(t)
	log, _ := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
	})

	service := NewHealthService(mockSpecRepo, mockHealthRepo, mockHTTPClient, log)
	return service, mockSpecRepo, mockHealthRepo, mockHTTPClient
}

func TestHealthService_CheckHealth(t *testing.T) {
	tests := []struct {
		name       string
		specID     string
		setupMocks func(*mocks.MockSpecRepository, *mocks.MockHealthRepository, *mocks.MockHTTPClient)
		wantErr    bool
		validate   func(*testing.T, *model.HealthCheck)
	}{
		{
			name:   "successful health check - status UP",
			specID: "spec-healthy",
			setupMocks: func(specRepo *mocks.MockSpecRepository, healthRepo *mocks.MockHealthRepository, httpClient *mocks.MockHTTPClient) {
				spec := testutil.CreateTestSpec("spec-healthy", "Healthy API")
				spec.BaseURL = "https://api.healthy.com"
				specRepo.EXPECT().
					FindByID(mock.Anything, "spec-healthy").
					Return(spec, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.MatchedBy(func(req *http.Request) bool {
						return req.Method == http.MethodHead && req.URL.String() == "https://api.healthy.com"
					})).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil)

				healthRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(hc *model.HealthCheck) bool {
						return hc.SpecID == "spec-healthy" && hc.Status == model.HealthStatusUP
					})).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, hc *model.HealthCheck) {
				assert.NotEmpty(t, hc.ID)
				assert.Equal(t, "spec-healthy", hc.SpecID)
				assert.Equal(t, model.HealthStatusUP, hc.Status)
				assert.Equal(t, "https://api.healthy.com", hc.CheckURL)
				assert.GreaterOrEqual(t, hc.ResponseTimeMs, 0)
				assert.Empty(t, hc.ErrorMessage)
				assert.NotZero(t, hc.CheckedAt)
			},
		},
		{
			name:   "health check - status DOWN (connection failed)",
			specID: "spec-down",
			setupMocks: func(specRepo *mocks.MockSpecRepository, healthRepo *mocks.MockHealthRepository, httpClient *mocks.MockHTTPClient) {
				spec := testutil.CreateTestSpec("spec-down", "Down API")
				spec.BaseURL = "https://api.down.com"
				specRepo.EXPECT().
					FindByID(mock.Anything, "spec-down").
					Return(spec, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(nil, errors.New("connection refused"))

				healthRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(hc *model.HealthCheck) bool {
						return hc.Status == model.HealthStatusDown
					})).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, hc *model.HealthCheck) {
				assert.Equal(t, model.HealthStatusDown, hc.Status)
				assert.Contains(t, hc.ErrorMessage, "Connection failed")
			},
		},
		{
			name:   "health check - status DOWN (500 error)",
			specID: "spec-500",
			setupMocks: func(specRepo *mocks.MockSpecRepository, healthRepo *mocks.MockHealthRepository, httpClient *mocks.MockHTTPClient) {
				spec := testutil.CreateTestSpec("spec-500", "API with 500")
				spec.BaseURL = "https://api.error.com"
				specRepo.EXPECT().
					FindByID(mock.Anything, "spec-500").
					Return(spec, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 500,
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil)

				healthRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, hc *model.HealthCheck) {
				assert.Equal(t, model.HealthStatusDown, hc.Status)
				assert.Contains(t, hc.ErrorMessage, "Server error: HTTP 500")
			},
		},
		{
			name:   "health check - status DEGRADED (4xx error)",
			specID: "spec-degraded",
			setupMocks: func(specRepo *mocks.MockSpecRepository, healthRepo *mocks.MockHealthRepository, httpClient *mocks.MockHTTPClient) {
				spec := testutil.CreateTestSpec("spec-degraded", "Degraded API")
				spec.BaseURL = "https://api.degraded.com"
				specRepo.EXPECT().
					FindByID(mock.Anything, "spec-degraded").
					Return(spec, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 404,
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil)

				healthRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, hc *model.HealthCheck) {
				assert.Equal(t, model.HealthStatusDegraded, hc.Status)
				assert.Contains(t, hc.ErrorMessage, "Unexpected status: HTTP 404")
			},
		},
		{
			name:   "spec not found",
			specID: "non-existent-spec",
			setupMocks: func(specRepo *mocks.MockSpecRepository, healthRepo *mocks.MockHealthRepository, httpClient *mocks.MockHTTPClient) {
				specRepo.EXPECT().
					FindByID(mock.Anything, "non-existent-spec").
					Return(nil, model.ErrSpecNotFound)
			},
			wantErr: true,
		},
		{
			name:   "health check with storage failure - should still return result",
			specID: "spec-storage-fail",
			setupMocks: func(specRepo *mocks.MockSpecRepository, healthRepo *mocks.MockHealthRepository, httpClient *mocks.MockHTTPClient) {
				spec := testutil.CreateTestSpec("spec-storage-fail", "API")
				spec.BaseURL = "https://api.example.com"
				specRepo.EXPECT().
					FindByID(mock.Anything, "spec-storage-fail").
					Return(spec, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil)

				healthRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("database error"))
			},
			wantErr: false, // Should continue even if storage fails
			validate: func(t *testing.T, hc *model.HealthCheck) {
				assert.Equal(t, model.HealthStatusUP, hc.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockSpecRepo, mockHealthRepo, mockHTTPClient := setupHealthServiceTest(t)
			tt.setupMocks(mockSpecRepo, mockHealthRepo, mockHTTPClient)

			// Act
			healthCheck, err := service.CheckHealth(context.Background(), tt.specID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, healthCheck)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, healthCheck)
				if tt.validate != nil {
					tt.validate(t, healthCheck)
				}
			}
		})
	}
}

func TestHealthService_GetHealthHistory(t *testing.T) {
	tests := []struct {
		name       string
		specID     string
		limit      int
		setupMocks func(*mocks.MockHealthRepository)
		wantErr    bool
		validate   func(*testing.T, []*model.HealthCheck, int)
	}{
		{
			name:   "successful get history with default limit",
			specID: "spec-1",
			limit:  0, // Should default to 10
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthChecks := []*model.HealthCheck{
					testutil.CreateTestHealthCheck("spec-1", model.HealthStatusUP),
					testutil.CreateTestHealthCheck("spec-1", model.HealthStatusUP),
					testutil.CreateTestHealthCheck("spec-1", model.HealthStatusDegraded),
				}
				healthRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-1", 10, 0).
					Return(healthChecks, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, healthChecks []*model.HealthCheck, limit int) {
				assert.Len(t, healthChecks, 3)
				for _, hc := range healthChecks {
					assert.Equal(t, "spec-1", hc.SpecID)
				}
			},
		},
		{
			name:   "get history with custom limit",
			specID: "spec-2",
			limit:  20,
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthChecks := make([]*model.HealthCheck, 20)
				for i := range healthChecks {
					healthChecks[i] = testutil.CreateTestHealthCheck("spec-2", model.HealthStatusUP)
				}
				healthRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-2", 20, 0).
					Return(healthChecks, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, healthChecks []*model.HealthCheck, limit int) {
				assert.Len(t, healthChecks, 20)
			},
		},
		{
			name:   "get history with maximum limit enforcement",
			specID: "spec-3",
			limit:  200, // Should be capped at 100
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthChecks := []*model.HealthCheck{}
				healthRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-3", 100, 0).
					Return(healthChecks, nil)
			},
			wantErr: false,
		},
		{
			name:   "empty history",
			specID: "spec-no-history",
			limit:  10,
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-no-history", 10, 0).
					Return([]*model.HealthCheck{}, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, healthChecks []*model.HealthCheck, limit int) {
				assert.Empty(t, healthChecks)
			},
		},
		{
			name:   "repository failure",
			specID: "spec-error",
			limit:  10,
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-error", 10, 0).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, _, mockHealthRepo, _ := setupHealthServiceTest(t)
			tt.setupMocks(mockHealthRepo)

			// Act
			healthChecks, err := service.GetHealthHistory(context.Background(), tt.specID, tt.limit)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, healthChecks)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, healthChecks)
				if tt.validate != nil {
					tt.validate(t, healthChecks, tt.limit)
				}
			}
		})
	}
}

func TestHealthService_GetLatestHealth(t *testing.T) {
	tests := []struct {
		name       string
		specID     string
		setupMocks func(*mocks.MockHealthRepository)
		wantErr    bool
		validate   func(*testing.T, *model.HealthCheck)
	}{
		{
			name:   "successful get latest health",
			specID: "spec-1",
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthCheck := testutil.CreateTestHealthCheck("spec-1", model.HealthStatusUP)
				healthRepo.EXPECT().
					FindLatestBySpecID(mock.Anything, "spec-1").
					Return(healthCheck, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, hc *model.HealthCheck) {
				assert.Equal(t, "spec-1", hc.SpecID)
				assert.Equal(t, model.HealthStatusUP, hc.Status)
			},
		},
		{
			name:   "no health check found",
			specID: "spec-no-health",
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthRepo.EXPECT().
					FindLatestBySpecID(mock.Anything, "spec-no-health").
					Return(nil, errors.New("no health check found"))
			},
			wantErr: true,
		},
		{
			name:   "repository failure",
			specID: "spec-error",
			setupMocks: func(healthRepo *mocks.MockHealthRepository) {
				healthRepo.EXPECT().
					FindLatestBySpecID(mock.Anything, "spec-error").
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, _, mockHealthRepo, _ := setupHealthServiceTest(t)
			tt.setupMocks(mockHealthRepo)

			// Act
			healthCheck, err := service.GetLatestHealth(context.Background(), tt.specID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, healthCheck)
				if tt.validate != nil {
					tt.validate(t, healthCheck)
				}
			}
		})
	}
}

func TestHealthService_NewHealthService(t *testing.T) {
	// Test constructor
	mockSpecRepo := mocks.NewMockSpecRepository(t)
	mockHealthRepo := mocks.NewMockHealthRepository(t)
	mockHTTPClient := mocks.NewMockHTTPClient(t)
	log, _ := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
	})

	service := NewHealthService(mockSpecRepo, mockHealthRepo, mockHTTPClient, log)

	require.NotNil(t, service)
	assert.NotNil(t, service.specRepo)
	assert.NotNil(t, service.healthRepo)
	assert.NotNil(t, service.httpClient)
	assert.NotNil(t, service.logger)
}
