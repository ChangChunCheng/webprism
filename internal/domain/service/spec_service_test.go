package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/mocks"
	"github.com/ChangChunCheng/webprism/internal/ports/output"
	"github.com/ChangChunCheng/webprism/internal/testutil"
)

// setupSpecServiceTest creates test environment with all mocks
func setupSpecServiceTest(t *testing.T) (*SpecService, *mocks.MockSpecRepository, *mocks.MockOpenAPIParser) {
	mockSpecRepo := mocks.NewMockSpecRepository(t)
	mockParser := mocks.NewMockOpenAPIParser(t)
	log, _ := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
	})

	service := NewSpecService(mockSpecRepo, mockParser, log)
	return service, mockSpecRepo, mockParser
}

func TestSpecService_UploadSpec(t *testing.T) {
	tests := []struct {
		name       string
		specName   string
		version    string
		specData   map[string]interface{}
		setupMocks func(*mocks.MockSpecRepository, *mocks.MockOpenAPIParser)
		wantErr    bool
		validate   func(*testing.T, *model.APISpec)
	}{
		{
			name:     "successful upload with valid OpenAPI spec",
			specName: "Petstore API",
			version:  "1.0.0",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]interface{}{
					"title":   "Petstore API",
					"version": "1.0.0",
				},
				"servers": []interface{}{
					map[string]interface{}{"url": "https://api.petstore.com"},
				},
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, parser *mocks.MockOpenAPIParser) {
				parser.EXPECT().
					Validate(mock.Anything, mock.Anything).
					Return(nil)

				parsedSpec := &output.ParsedSpec{
					BaseURL: "https://api.petstore.com",
					Paths: map[string]*model.PathItem{
						"/pets": {
							Path: "/pets",
							Operations: map[string]*model.Operation{
								"GET": {
									OperationID: "listPets",
									Method:      "GET",
									Summary:     "List all pets",
								},
							},
						},
					},
				}
				parser.EXPECT().
					Parse(mock.Anything, mock.Anything).
					Return(parsedSpec, nil)

				specRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(spec *model.APISpec) bool {
						return spec.Name == "Petstore API" &&
							spec.Version == "1.0.0" &&
							spec.BaseURL == "https://api.petstore.com"
					})).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, spec *model.APISpec) {
				assert.NotEmpty(t, spec.ID)
				assert.Equal(t, "Petstore API", spec.Name)
				assert.Equal(t, "1.0.0", spec.Version)
				assert.Equal(t, "https://api.petstore.com", spec.BaseURL)
				assert.NotNil(t, spec.Paths)
				assert.Contains(t, spec.Paths, "/pets")
				assert.NotZero(t, spec.CreatedAt)
				assert.NotZero(t, spec.UpdatedAt)
			},
		},
		{
			name:     "validation failure - invalid OpenAPI format",
			specName: "Invalid API",
			version:  "1.0.0",
			specData: map[string]interface{}{
				"invalid": "data",
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, parser *mocks.MockOpenAPIParser) {
				parser.EXPECT().
					Validate(mock.Anything, mock.Anything).
					Return(errors.New("missing required field: openapi"))
			},
			wantErr: true,
		},
		{
			name:     "parse failure - malformed spec",
			specName: "Malformed API",
			version:  "1.0.0",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"info":    "invalid",
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, parser *mocks.MockOpenAPIParser) {
				parser.EXPECT().
					Validate(mock.Anything, mock.Anything).
					Return(nil)

				parser.EXPECT().
					Parse(mock.Anything, mock.Anything).
					Return(nil, errors.New("failed to parse info section"))
			},
			wantErr: true,
		},
		{
			name:     "repository failure - database error",
			specName: "Test API",
			version:  "1.0.0",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, parser *mocks.MockOpenAPIParser) {
				parser.EXPECT().
					Validate(mock.Anything, mock.Anything).
					Return(nil)

				parsedSpec := &output.ParsedSpec{
					BaseURL: "https://api.test.com",
					Paths:   make(map[string]*model.PathItem),
				}
				parser.EXPECT().
					Parse(mock.Anything, mock.Anything).
					Return(parsedSpec, nil)

				specRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("database connection failed"))
			},
			wantErr: true,
		},
		{
			name:     "upload with complex paths",
			specName: "Complex API",
			version:  "2.0.0",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, parser *mocks.MockOpenAPIParser) {
				parser.EXPECT().
					Validate(mock.Anything, mock.Anything).
					Return(nil)

				parsedSpec := &output.ParsedSpec{
					BaseURL: "https://api.complex.com",
					Paths: map[string]*model.PathItem{
						"/users/{id}": {
							Path: "/users/{id}",
							Operations: map[string]*model.Operation{
								"GET": {OperationID: "getUser", Method: "GET"},
								"PUT": {OperationID: "updateUser", Method: "PUT"},
							},
						},
						"/users": {
							Path: "/users",
							Operations: map[string]*model.Operation{
								"GET":  {OperationID: "listUsers", Method: "GET"},
								"POST": {OperationID: "createUser", Method: "POST"},
							},
						},
					},
				}
				parser.EXPECT().
					Parse(mock.Anything, mock.Anything).
					Return(parsedSpec, nil)

				specRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, spec *model.APISpec) {
				assert.Len(t, spec.Paths, 2)
				assert.Contains(t, spec.Paths, "/users/{id}")
				assert.Contains(t, spec.Paths, "/users")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockSpecRepo, mockParser := setupSpecServiceTest(t)
			tt.setupMocks(mockSpecRepo, mockParser)

			// Act
			spec, err := service.UploadSpec(context.Background(), tt.specName, tt.version, tt.specData)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, spec)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, spec)
				if tt.validate != nil {
					tt.validate(t, spec)
				}
			}
		})
	}
}

func TestSpecService_ListSpecs(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		setupMocks func(*mocks.MockSpecRepository)
		wantErr    bool
		validate   func(*testing.T, []*model.APISpec, int, int)
	}{
		{
			name:   "successful list with default limit",
			limit:  0, // Should default to 10
			offset: 0,
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specs := []*model.APISpec{
					testutil.CreateTestSpec("spec-1", "API 1"),
					testutil.CreateTestSpec("spec-2", "API 2"),
				}
				specRepo.EXPECT().
					FindAll(mock.Anything, 10, 0). // Default limit is 10
					Return(specs, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, specs []*model.APISpec, limit, offset int) {
				assert.Len(t, specs, 2)
				assert.Equal(t, "spec-1", specs[0].ID)
				assert.Equal(t, "spec-2", specs[1].ID)
			},
		},
		{
			name:   "list with custom limit",
			limit:  20,
			offset: 0,
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specs := make([]*model.APISpec, 0)
				for i := 1; i <= 20; i++ {
					specs = append(specs, testutil.CreateTestSpec("spec", "API"))
				}
				specRepo.EXPECT().
					FindAll(mock.Anything, 20, 0).
					Return(specs, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, specs []*model.APISpec, limit, offset int) {
				assert.Len(t, specs, 20)
			},
		},
		{
			name:   "list with maximum limit enforcement",
			limit:  200, // Exceeds max limit of 100
			offset: 0,
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specs := make([]*model.APISpec, 0)
				specRepo.EXPECT().
					FindAll(mock.Anything, 100, 0). // Should be capped at 100
					Return(specs, nil)
			},
			wantErr: false,
		},
		{
			name:   "list with pagination offset",
			limit:  10,
			offset: 20,
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specs := []*model.APISpec{
					testutil.CreateTestSpec("spec-21", "API 21"),
					testutil.CreateTestSpec("spec-22", "API 22"),
				}
				specRepo.EXPECT().
					FindAll(mock.Anything, 10, 20).
					Return(specs, nil)
			},
			wantErr: false,
		},
		{
			name:   "list with negative offset - should normalize to 0",
			limit:  10,
			offset: -5,
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specs := []*model.APISpec{}
				specRepo.EXPECT().
					FindAll(mock.Anything, 10, 0). // Negative offset normalized to 0
					Return(specs, nil)
			},
			wantErr: false,
		},
		{
			name:   "empty result set",
			limit:  10,
			offset: 0,
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specRepo.EXPECT().
					FindAll(mock.Anything, 10, 0).
					Return([]*model.APISpec{}, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, specs []*model.APISpec, limit, offset int) {
				assert.Empty(t, specs)
			},
		},
		{
			name:   "repository failure",
			limit:  10,
			offset: 0,
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specRepo.EXPECT().
					FindAll(mock.Anything, 10, 0).
					Return(nil, errors.New("database connection failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockSpecRepo, _ := setupSpecServiceTest(t)
			tt.setupMocks(mockSpecRepo)

			// Act
			specs, err := service.ListSpecs(context.Background(), tt.limit, tt.offset)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, specs)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, specs)
				if tt.validate != nil {
					tt.validate(t, specs, tt.limit, tt.offset)
				}
			}
		})
	}
}

func TestSpecService_GetSpec(t *testing.T) {
	tests := []struct {
		name       string
		specID     string
		setupMocks func(*mocks.MockSpecRepository)
		wantErr    bool
		wantErrIs  error
		validate   func(*testing.T, *model.APISpec)
	}{
		{
			name:   "successful get spec",
			specID: "existing-spec-id",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				spec := testutil.CreateTestSpecWithPaths("existing-spec-id", "My API")
				specRepo.EXPECT().
					FindByID(mock.Anything, "existing-spec-id").
					Return(spec, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, spec *model.APISpec) {
				assert.Equal(t, "existing-spec-id", spec.ID)
				assert.Equal(t, "My API", spec.Name)
				assert.NotNil(t, spec.Paths)
			},
		},
		{
			name:   "spec not found",
			specID: "non-existent-spec",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specRepo.EXPECT().
					FindByID(mock.Anything, "non-existent-spec").
					Return(nil, model.ErrSpecNotFound)
			},
			wantErr:   true,
			wantErrIs: model.ErrSpecNotFound,
		},
		{
			name:   "repository failure",
			specID: "some-spec-id",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specRepo.EXPECT().
					FindByID(mock.Anything, "some-spec-id").
					Return(nil, errors.New("database connection lost"))
			},
			wantErr: true,
		},
		{
			name:   "get spec with complete data",
			specID: "complete-spec-id",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				spec := testutil.CreateTestSpecWithPaths("complete-spec-id", "Complete API")
				spec.Version = "2.0.0"
				spec.BaseURL = "https://api.complete.com"
				specRepo.EXPECT().
					FindByID(mock.Anything, "complete-spec-id").
					Return(spec, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, spec *model.APISpec) {
				assert.Equal(t, "2.0.0", spec.Version)
				assert.Equal(t, "https://api.complete.com", spec.BaseURL)
				assert.NotEmpty(t, spec.Paths)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockSpecRepo, _ := setupSpecServiceTest(t)
			tt.setupMocks(mockSpecRepo)

			// Act
			spec, err := service.GetSpec(context.Background(), tt.specID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, spec)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, spec)
				if tt.validate != nil {
					tt.validate(t, spec)
				}
			}
		})
	}
}

func TestSpecService_DeleteSpec(t *testing.T) {
	tests := []struct {
		name       string
		specID     string
		setupMocks func(*mocks.MockSpecRepository)
		wantErr    bool
		wantErrIs  error
	}{
		{
			name:   "successful delete",
			specID: "spec-to-delete",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specRepo.EXPECT().
					Delete(mock.Anything, "spec-to-delete").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "delete non-existent spec",
			specID: "non-existent",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specRepo.EXPECT().
					Delete(mock.Anything, "non-existent").
					Return(model.ErrSpecNotFound)
			},
			wantErr:   true,
			wantErrIs: model.ErrSpecNotFound,
		},
		{
			name:   "repository failure",
			specID: "some-spec",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				specRepo.EXPECT().
					Delete(mock.Anything, "some-spec").
					Return(errors.New("foreign key constraint violation"))
			},
			wantErr: true,
		},
		{
			name:   "delete with cascade",
			specID: "spec-with-dependencies",
			setupMocks: func(specRepo *mocks.MockSpecRepository) {
				// Cascade should be handled by repository/database
				specRepo.EXPECT().
					Delete(mock.Anything, "spec-with-dependencies").
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockSpecRepo, _ := setupSpecServiceTest(t)
			tt.setupMocks(mockSpecRepo)

			// Act
			err := service.DeleteSpec(context.Background(), tt.specID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSpecService_NewSpecService(t *testing.T) {
	// Test constructor
	mockSpecRepo := mocks.NewMockSpecRepository(t)
	mockParser := mocks.NewMockOpenAPIParser(t)
	log, _ := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
	})

	service := NewSpecService(mockSpecRepo, mockParser, log)

	require.NotNil(t, service)
	assert.NotNil(t, service.specRepo)
	assert.NotNil(t, service.parser)
	assert.NotNil(t, service.logger)
}
