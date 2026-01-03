// TODO(Architecture): 考慮支援一個可擴充的「委派驗證」(Delegated Authentication) 模式。
//
// 該模式的設想如下：
// 1. 新增一個 AuthType，例如 `AUTH_TYPE_WEBHOOK`。
// 2. 當 AuthConfig 為此類型時，`Credentials` 中儲存的將是一個外部服務的 URL。
// 3. ProxyService 在代理請求時，會先呼叫此外部 URL（Call-out）。
// 4. 該外部服務可以選擇性地回呼（Call-back）WEBPRISM 的 GetAuthConfig API，
//    以取得儲存在系統中的原始設定（例如，一個靜態的 API Key 或其他元數據）。
// 5. 外部服務根據自己的邏輯和從 WEBPRISM 取得的資訊，動態決定驗證結果，
//    並將最終的驗證憑證（或拒絕原因）回傳給 WEBPRISM。
// 6. WEBPRISM 根據回傳結果繼續執行代理或拒絕請求。
//
// 實現此功能需要：
// - 保護 GetAuthConfig API，使其能被受信任的外部服務安全地呼叫。
// - 在儲存 WEBHOOK 類型的 AuthConfig 時，建立一個健康檢查機制，以驗證外部服務的可用性。

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/output"
)

// AuthService implements the AuthService input port.
type AuthService struct {
	authRepo output.AuthRepository
	specRepo output.SpecRepository
	crypto   output.CryptoService
	logger   *logger.Logger
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(
	authRepo output.AuthRepository,
	specRepo output.SpecRepository,
	crypto output.CryptoService,
	log *logger.Logger,
) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		specRepo: specRepo,
		crypto:   crypto,
		logger:   log,
	}
}

// SetAuthConfig sets or updates the authentication configuration for a spec.
func (s *AuthService) SetAuthConfig(ctx context.Context, specID string, authType model.AuthType, credentials map[string]string) (*model.AuthConfig, error) {
	s.logger.Info("setting auth config",
		logger.String("spec_id", specID),
		logger.String("auth_type", string(authType)),
	)

	// Verify that the spec exists
	exists, err := s.specRepo.Exists(ctx, specID)
	if err != nil {
		return nil, fmt.Errorf("failed to check spec existence: %w", err)
	}
	if !exists {
		return nil, model.ErrSpecNotFound
	}

	// Validate auth config
	authConfig := &model.AuthConfig{
		SpecID:      specID,
		AuthType:    authType,
		Credentials: credentials,
	}

	if err := authConfig.Validate(); err != nil {
		s.logger.Error("invalid auth config",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, err
	}

	// Encrypt credentials
	encrypted, err := s.encryptCredentials(credentials)
	if err != nil {
		s.logger.Error("failed to encrypt credentials",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to encrypt credentials: %w", err)
	}

	authConfig.EncryptedCredentials = encrypted

	// Check if auth config already exists
	authExists, err := s.authRepo.Exists(ctx, specID)
	if err != nil {
		return nil, fmt.Errorf("failed to check auth config existence: %w", err)
	}

	now := time.Now()
	authConfig.CreatedAt = now
	authConfig.UpdatedAt = now

	// Create or update
	if authExists {
		if err := s.authRepo.Update(ctx, authConfig); err != nil {
			s.logger.Error("failed to update auth config",
				logger.String("spec_id", specID),
				logger.Any("error", err),
			)
			return nil, fmt.Errorf("failed to update auth config: %w", err)
		}
		s.logger.Info("auth config updated successfully",
			logger.String("spec_id", specID),
		)
	} else {
		if err := s.authRepo.Create(ctx, authConfig); err != nil {
			s.logger.Error("failed to create auth config",
				logger.String("spec_id", specID),
				logger.Any("error", err),
			)
			return nil, fmt.Errorf("failed to create auth config: %w", err)
		}
		s.logger.Info("auth config created successfully",
			logger.String("spec_id", specID),
		)
	}

	return authConfig, nil
}

// GetAuthConfig retrieves the authentication configuration for a spec.
func (s *AuthService) GetAuthConfig(ctx context.Context, specID string) (*model.AuthConfig, error) {
	s.logger.Debug("getting auth config",
		logger.String("spec_id", specID),
	)

	// Retrieve auth config
	authConfig, err := s.authRepo.FindBySpecID(ctx, specID)
	if err != nil {
		s.logger.Error("failed to get auth config",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, err
	}

	// Decrypt credentials
	credentials, err := s.decryptCredentials(authConfig.EncryptedCredentials)
	if err != nil {
		s.logger.Error("failed to decrypt credentials",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return nil, fmt.Errorf("failed to decrypt credentials: %w", err)
	}

	authConfig.Credentials = credentials

	s.logger.Debug("auth config retrieved successfully",
		logger.String("spec_id", specID),
		logger.String("auth_type", string(authConfig.AuthType)),
	)

	return authConfig, nil
}

// DeleteAuthConfig removes the authentication configuration for a spec.
func (s *AuthService) DeleteAuthConfig(ctx context.Context, specID string) error {
	s.logger.Info("deleting auth config",
		logger.String("spec_id", specID),
	)

	if err := s.authRepo.Delete(ctx, specID); err != nil {
		s.logger.Error("failed to delete auth config",
			logger.String("spec_id", specID),
			logger.Any("error", err),
		)
		return err
	}

	s.logger.Info("auth config deleted successfully",
		logger.String("spec_id", specID),
	)

	return nil
}

// encryptCredentials encrypts the credentials map using the crypto service.
func (s *AuthService) encryptCredentials(credentials map[string]string) ([]byte, error) {
	// Marshal credentials to JSON
	data, err := json.Marshal(credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Encrypt
	encrypted, err := s.crypto.Encrypt(data)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// decryptCredentials decrypts the encrypted credentials.
func (s *AuthService) decryptCredentials(encrypted []byte) (map[string]string, error) {
	// Decrypt
	decrypted, err := s.crypto.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var credentials map[string]string
	if err := json.Unmarshal(decrypted, &credentials); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return credentials, nil
}
