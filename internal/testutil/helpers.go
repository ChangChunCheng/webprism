package testutil

import (
	"time"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
)

// CreateTestSpec 創建測試用的 APISpec
func CreateTestSpec(id, name string) *model.APISpec {
	return &model.APISpec{
		ID:        id,
		Name:      name,
		Version:   "1.0.0",
		BaseURL:   "https://api.example.com",
		SpecData:  make(map[string]interface{}),
		Paths:     make(map[string]*model.PathItem),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestSpecWithPaths 創建帶有 paths 的測試 spec
func CreateTestSpecWithPaths(id, name string) *model.APISpec {
	spec := CreateTestSpec(id, name)
	spec.Paths = map[string]*model.PathItem{
		"/users/{id}": {
			Path: "/users/{id}",
			Operations: map[string]*model.Operation{
				"GET": {
					OperationID: "getUser",
					Method:      "GET",
					Summary:     "Get user by ID",
					Parameters: []*model.Parameter{
						{
							Name:     "id",
							In:       "path",
							Required: true,
							Type:     "string",
						},
					},
				},
			},
		},
		"/users": {
			Path: "/users",
			Operations: map[string]*model.Operation{
				"GET": {
					OperationID: "listUsers",
					Method:      "GET",
					Summary:     "List all users",
					Parameters: []*model.Parameter{
						{
							Name:     "limit",
							In:       "query",
							Required: false,
							Type:     "integer",
						},
					},
				},
				"POST": {
					OperationID: "createUser",
					Method:      "POST",
					Summary:     "Create a new user",
					RequestBody: &model.RequestBody{
						Required: true,
						Content: map[string]*model.MediaType{
							"application/json": {
								Schema: &model.Schema{
									Type: "object",
								},
							},
						},
					},
				},
			},
		},
	}
	return spec
}

// CreateTestProxyRequest 創建測試用的代理請求
func CreateTestProxyRequest(specID, operationID string) *model.ProxyRequest {
	return &model.ProxyRequest{
		SpecID:      specID,
		OperationID: operationID,
		Parameters:  model.ProxyParameters{},
	}
}

// CreateTestAuthConfig 創建測試用的認證配置
func CreateTestAuthConfig(authType model.AuthType) *model.AuthConfig {
	credentials := make(map[string]string)

	switch authType {
	case model.AuthTypeAPIKey:
		credentials["key"] = "test-api-key"
		credentials["name"] = "X-API-Key"
	case model.AuthTypeBearer:
		credentials["token"] = "test-bearer-token"
	case model.AuthTypeBasic:
		credentials["username"] = "testuser"
		credentials["password"] = "testpass"
	case model.AuthTypeOAuth2:
		credentials["access_token"] = "test-oauth-token"
	}

	return &model.AuthConfig{
		AuthType:    authType,
		Credentials: credentials,
	}
}

// CreateTestHealthCheck 創建測試用的健康檢查記錄
func CreateTestHealthCheck(specID string, status model.HealthStatus) *model.HealthCheck {
	return &model.HealthCheck{
		ID:             "health-" + specID,
		SpecID:         specID,
		CheckURL:       "https://api.example.com/health",
		Status:         status,
		ResponseTimeMs: 100,
		CheckedAt:      time.Now(),
	}
}
