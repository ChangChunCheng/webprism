package model

import "time"

// AuthConfig 代表認證配置
type AuthConfig struct {
	// TODO: 為提高可讀性，此欄位應考慮重構成 APISpecID，以更明確地表示與 APISpec.ID 的關聯。
	SpecID               string       `json:"spec_id" bun:"spec_id,pk"`                         // 關聯的 API Spec ID，其值等於 APISpec.ID
	AuthType             AuthType     `json:"auth_type" bun:"auth_type,notnull"`                // 認證類型
	EncryptedCredentials []byte       `json:"-" bun:"encrypted_credentials,type:bytea,notnull"` // 加密後的憑證
	Position             AuthPosition `json:"position" bun:"-"`                                 // 認證位置（僅用於記憶體，未持久化）
	CreatedAt            time.Time    `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt            time.Time    `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`

	// 解密後的憑證（不儲存到資料庫）
	Credentials map[string]string `json:"credentials,omitempty" bun:"-"`
}

// AuthType 認證類型
type AuthType string

const (
	AuthTypeUnspecified AuthType = ""
	AuthTypeAPIKey      AuthType = "API_KEY"
	AuthTypeBearer      AuthType = "BEARER"
	AuthTypeBasic       AuthType = "BASIC"
	AuthTypeOAuth2      AuthType = "OAUTH2"
)

// AuthPosition 認證位置
type AuthPosition string

const (
	AuthPositionUnspecified AuthPosition = ""
	AuthPositionHeader      AuthPosition = "HEADER"
	AuthPositionQuery       AuthPosition = "QUERY"
	AuthPositionBody        AuthPosition = "BODY"
)

// Validate 驗證認證配置
func (a *AuthConfig) Validate() error {
	if a.SpecID == "" {
		return ErrMissingSpecID
	}
	if a.AuthType == AuthTypeUnspecified {
		return ErrInvalidAuthType
	}
	if len(a.Credentials) == 0 {
		return ErrMissingCredentials
	}
	return nil
}
