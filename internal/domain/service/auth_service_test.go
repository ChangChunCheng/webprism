package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/mocks"
	"github.com/ChangChunCheng/webprism/internal/testutil"
)

// setupAuthServiceTest creates test environment with all mocks
func setupAuthServiceTest(t *testing.T) (*AuthService, *mocks.MockAuthRepository, *mocks.MockSpecRepository, *mocks.MockCryptoService) {
	mockAuthRepo := mocks.NewMockAuthRepository(t)
	mockSpecRepo := mocks.NewMockSpecRepository(t)
	mockCrypto := mocks.NewMockCryptoService(t)
	log, _ := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
	})

	service := NewAuthService(mockAuthRepo, mockSpecRepo, mockCrypto, log)
	return service, mockAuthRepo, mockSpecRepo, mockCrypto
}

func TestAuthService_SetAuthConfig(t *testing.T) {
	tests := []struct {
		name        string
		specID      string
		authType    model.AuthType
		credentials map[string]string
		setupMocks  func(*mocks.MockAuthRepository, *mocks.MockSpecRepository, *mocks.MockCryptoService)
		wantErr     bool
		wantErrIs   error
		validate    func(*testing.T, *model.AuthConfig)
	}{
		{
			name:     "create new bearer auth config",
			specID:   "spec-1",
			authType: model.AuthTypeBearer,
			credentials: map[string]string{
				"token": "secret-bearer-token",
			},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "spec-1").
					Return(true, nil)

				credJSON, _ := json.Marshal(map[string]string{"token": "secret-bearer-token"})
				crypto.EXPECT().
					Encrypt(credJSON).
					Return([]byte("encrypted-bearer-token"), nil)

				authRepo.EXPECT().
					Exists(mock.Anything, "spec-1").
					Return(false, nil)

				authRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(config *model.AuthConfig) bool {
						return config.SpecID == "spec-1" &&
							config.AuthType == model.AuthTypeBearer &&
							string(config.EncryptedCredentials) == "encrypted-bearer-token"
					})).
					Return(nil)
			},
			wantErr: false,
			validate: func(t *testing.T, config *model.AuthConfig) {
				assert.Equal(t, "spec-1", config.SpecID)
				assert.Equal(t, model.AuthTypeBearer, config.AuthType)
				assert.Equal(t, []byte("encrypted-bearer-token"), config.EncryptedCredentials)
				assert.NotZero(t, config.CreatedAt)
				assert.NotZero(t, config.UpdatedAt)
			},
		},
		{
			name:     "update existing basic auth config",
			specID:   "spec-2",
			authType: model.AuthTypeBasic,
			credentials: map[string]string{
				"username": "admin",
				"password": "newpassword",
			},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "spec-2").
					Return(true, nil)

				credJSON, _ := json.Marshal(map[string]string{"username": "admin", "password": "newpassword"})
				crypto.EXPECT().
					Encrypt(credJSON).
					Return([]byte("encrypted-basic-auth"), nil)

				authRepo.EXPECT().
					Exists(mock.Anything, "spec-2").
					Return(true, nil)

				authRepo.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "create api key auth config",
			specID:   "spec-3",
			authType: model.AuthTypeAPIKey,
			credentials: map[string]string{
				"key":   "X-API-Key",
				"value": "api-key-12345",
			},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "spec-3").
					Return(true, nil)

				credJSON, _ := json.Marshal(map[string]string{"key": "X-API-Key", "value": "api-key-12345"})
				crypto.EXPECT().
					Encrypt(credJSON).
					Return([]byte("encrypted-api-key"), nil)

				authRepo.EXPECT().
					Exists(mock.Anything, "spec-3").
					Return(false, nil)

				authRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "spec not found",
			specID:   "non-existent-spec",
			authType: model.AuthTypeBearer,
			credentials: map[string]string{
				"token": "some-token",
			},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "non-existent-spec").
					Return(false, nil)
			},
			wantErr:   true,
			wantErrIs: model.ErrSpecNotFound,
		},
		{
			name:        "invalid auth config - missing credentials",
			specID:      "spec-4",
			authType:    model.AuthTypeBearer,
			credentials: map[string]string{}, // Empty credentials
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "spec-4").
					Return(true, nil)
			},
			wantErr:   true,
			wantErrIs: model.ErrMissingCredentials,
		},
		{
			name:        "invalid auth config - empty spec ID",
			specID:      "",
			authType:    model.AuthTypeBearer,
			credentials: map[string]string{"token": "token"},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "").
					Return(false, nil)
			},
			wantErr:   true,
			wantErrIs: model.ErrSpecNotFound,
		},
		{
			name:     "encryption failure",
			specID:   "spec-5",
			authType: model.AuthTypeBearer,
			credentials: map[string]string{
				"token": "token",
			},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "spec-5").
					Return(true, nil)

				crypto.EXPECT().
					Encrypt(mock.Anything).
					Return(nil, errors.New("encryption failed"))
			},
			wantErr: true,
		},
		{
			name:     "repository creation failure",
			specID:   "spec-6",
			authType: model.AuthTypeBearer,
			credentials: map[string]string{
				"token": "token",
			},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "spec-6").
					Return(true, nil)

				crypto.EXPECT().
					Encrypt(mock.Anything).
					Return([]byte("encrypted"), nil)

				authRepo.EXPECT().
					Exists(mock.Anything, "spec-6").
					Return(false, nil)

				authRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:     "repository update failure",
			specID:   "spec-7",
			authType: model.AuthTypeBearer,
			credentials: map[string]string{
				"token": "token",
			},
			setupMocks: func(authRepo *mocks.MockAuthRepository, specRepo *mocks.MockSpecRepository, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					Exists(mock.Anything, "spec-7").
					Return(true, nil)

				crypto.EXPECT().
					Encrypt(mock.Anything).
					Return([]byte("encrypted"), nil)

				authRepo.EXPECT().
					Exists(mock.Anything, "spec-7").
					Return(true, nil)

				authRepo.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockAuthRepo, mockSpecRepo, mockCrypto := setupAuthServiceTest(t)
			tt.setupMocks(mockAuthRepo, mockSpecRepo, mockCrypto)

			// Act
			config, err := service.SetAuthConfig(context.Background(), tt.specID, tt.authType, tt.credentials)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				if tt.validate != nil {
					tt.validate(t, config)
				}
			}
		})
	}
}

func TestAuthService_GetAuthConfig(t *testing.T) {
	tests := []struct {
		name       string
		specID     string
		setupMocks func(*mocks.MockAuthRepository, *mocks.MockCryptoService)
		wantErr    bool
		wantErrIs  error
		validate   func(*testing.T, *model.AuthConfig)
	}{
		{
			name:   "successful get with bearer auth",
			specID: "spec-1",
			setupMocks: func(authRepo *mocks.MockAuthRepository, crypto *mocks.MockCryptoService) {
				authConfig := testutil.CreateTestAuthConfig(model.AuthTypeBearer)
				authConfig.SpecID = "spec-1"
				authConfig.EncryptedCredentials = []byte("encrypted-bearer")

				authRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-1").
					Return(authConfig, nil)

				credentials := map[string]string{"token": "decrypted-bearer-token"}
				credJSON, _ := json.Marshal(credentials)
				crypto.EXPECT().
					Decrypt([]byte("encrypted-bearer")).
					Return(credJSON, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, config *model.AuthConfig) {
				assert.Equal(t, "spec-1", config.SpecID)
				assert.Equal(t, model.AuthTypeBearer, config.AuthType)
				assert.NotNil(t, config.Credentials)
				assert.Equal(t, "decrypted-bearer-token", config.Credentials["token"])
			},
		},
		{
			name:   "successful get with basic auth",
			specID: "spec-2",
			setupMocks: func(authRepo *mocks.MockAuthRepository, crypto *mocks.MockCryptoService) {
				authConfig := testutil.CreateTestAuthConfig(model.AuthTypeBasic)
				authConfig.SpecID = "spec-2"
				authConfig.EncryptedCredentials = []byte("encrypted-basic")

				authRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-2").
					Return(authConfig, nil)

				credentials := map[string]string{"username": "user", "password": "pass"}
				credJSON, _ := json.Marshal(credentials)
				crypto.EXPECT().
					Decrypt([]byte("encrypted-basic")).
					Return(credJSON, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, config *model.AuthConfig) {
				assert.Equal(t, "user", config.Credentials["username"])
				assert.Equal(t, "pass", config.Credentials["password"])
			},
		},
		{
			name:   "auth config not found",
			specID: "non-existent-spec",
			setupMocks: func(authRepo *mocks.MockAuthRepository, crypto *mocks.MockCryptoService) {
				authRepo.EXPECT().
					FindBySpecID(mock.Anything, "non-existent-spec").
					Return(nil, model.ErrAuthConfigNotFound)
			},
			wantErr:   true,
			wantErrIs: model.ErrAuthConfigNotFound,
		},
		{
			name:   "decryption failure",
			specID: "spec-decrypt-fail",
			setupMocks: func(authRepo *mocks.MockAuthRepository, crypto *mocks.MockCryptoService) {
				authConfig := testutil.CreateTestAuthConfig(model.AuthTypeBearer)
				authConfig.SpecID = "spec-decrypt-fail"
				authConfig.EncryptedCredentials = []byte("corrupted-data")

				authRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-decrypt-fail").
					Return(authConfig, nil)

				crypto.EXPECT().
					Decrypt([]byte("corrupted-data")).
					Return(nil, errors.New("decryption failed"))
			},
			wantErr: true,
		},
		{
			name:   "repository failure",
			specID: "spec-repo-fail",
			setupMocks: func(authRepo *mocks.MockAuthRepository, crypto *mocks.MockCryptoService) {
				authRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-repo-fail").
					Return(nil, errors.New("database connection lost"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockAuthRepo, _, mockCrypto := setupAuthServiceTest(t)
			tt.setupMocks(mockAuthRepo, mockCrypto)

			// Act
			config, err := service.GetAuthConfig(context.Background(), tt.specID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				if tt.validate != nil {
					tt.validate(t, config)
				}
			}
		})
	}
}

func TestAuthService_DeleteAuthConfig(t *testing.T) {
	tests := []struct {
		name       string
		specID     string
		setupMocks func(*mocks.MockAuthRepository)
		wantErr    bool
		wantErrIs  error
	}{
		{
			name:   "successful delete",
			specID: "spec-1",
			setupMocks: func(authRepo *mocks.MockAuthRepository) {
				authRepo.EXPECT().
					Delete(mock.Anything, "spec-1").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "delete non-existent auth config",
			specID: "non-existent",
			setupMocks: func(authRepo *mocks.MockAuthRepository) {
				authRepo.EXPECT().
					Delete(mock.Anything, "non-existent").
					Return(model.ErrAuthConfigNotFound)
			},
			wantErr:   true,
			wantErrIs: model.ErrAuthConfigNotFound,
		},
		{
			name:   "repository failure",
			specID: "spec-error",
			setupMocks: func(authRepo *mocks.MockAuthRepository) {
				authRepo.EXPECT().
					Delete(mock.Anything, "spec-error").
					Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockAuthRepo, _, _ := setupAuthServiceTest(t)
			tt.setupMocks(mockAuthRepo)

			// Act
			err := service.DeleteAuthConfig(context.Background(), tt.specID)

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

func TestAuthService_NewAuthService(t *testing.T) {
	// Test constructor
	mockAuthRepo := mocks.NewMockAuthRepository(t)
	mockSpecRepo := mocks.NewMockSpecRepository(t)
	mockCrypto := mocks.NewMockCryptoService(t)
	log, _ := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
	})

	service := NewAuthService(mockAuthRepo, mockSpecRepo, mockCrypto, log)

	require.NotNil(t, service)
	assert.NotNil(t, service.authRepo)
	assert.NotNil(t, service.specRepo)
	assert.NotNil(t, service.crypto)
	assert.NotNil(t, service.logger)
}
