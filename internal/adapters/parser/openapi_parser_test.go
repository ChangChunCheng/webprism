package parser

import (
	"context"
	"testing"

	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenAPIParser(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)

	assert.NotNil(t, parser)
	assert.NotNil(t, parser.logger)
}

func TestDetectVersion(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)

	tests := []struct {
		name     string
		specData map[string]interface{}
		want     string
		wantErr  bool
	}{
		{
			name: "OpenAPI 3.0",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
			},
			want:    "3.0",
			wantErr: false,
		},
		{
			name: "OpenAPI 3.0.3",
			specData: map[string]interface{}{
				"openapi": "3.0.3",
			},
			want:    "3.0",
			wantErr: false,
		},
		{
			name: "OpenAPI 3.1",
			specData: map[string]interface{}{
				"openapi": "3.1.0",
			},
			want:    "3.1",
			wantErr: false,
		},
		{
			name: "Swagger 2.0",
			specData: map[string]interface{}{
				"swagger": "2.0",
			},
			want:    "2.0",
			wantErr: false,
		},
		{
			name:     "no version field",
			specData: map[string]interface{}{},
			want:     "",
			wantErr:  true,
		},
		{
			name: "invalid version",
			specData: map[string]interface{}{
				"openapi": "4.0.0",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.detectVersion(tt.specData)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParse_OpenAPI30(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)
	ctx := context.Background()

	specData := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Test API",
			"description": "Test API Description",
			"version":     "1.0.0",
		},
		"servers": []interface{}{
			map[string]interface{}{
				"url": "https://api.example.com",
			},
		},
		"paths": map[string]interface{}{
			"/users": map[string]interface{}{
				"get": map[string]interface{}{
					"operationId": "getUsers",
					"summary":     "Get all users",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
					},
				},
			},
		},
	}

	result, err := parser.Parse(ctx, specData)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "3.0", result.Version)
	assert.Equal(t, "https://api.example.com", result.BaseURL)
	assert.NotNil(t, result.Info)
	assert.Equal(t, "Test API", result.Info.Title)
	assert.Equal(t, "Test API Description", result.Info.Description)
	assert.Equal(t, "1.0.0", result.Info.Version)
	assert.NotEmpty(t, result.Paths)
	assert.Contains(t, result.Paths, "/users")
}

func TestParse_Swagger20(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)
	ctx := context.Background()

	specData := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       "Test API",
			"description": "Test API Description",
			"version":     "1.0.0",
		},
		"host":     "api.example.com",
		"basePath": "/v1",
		"schemes": []interface{}{
			"https",
		},
		"paths": map[string]interface{}{
			"/users": map[string]interface{}{
				"get": map[string]interface{}{
					"operationId": "getUsers",
					"summary":     "Get all users",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
					},
				},
			},
		},
	}

	result, err := parser.Parse(ctx, specData)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2.0", result.Version)
	assert.Equal(t, "https://api.example.com/v1", result.BaseURL)
	assert.NotNil(t, result.Info)
	assert.Equal(t, "Test API", result.Info.Title)
	assert.NotEmpty(t, result.Paths)
}

func TestParse_InvalidSpec(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)
	ctx := context.Background()

	tests := []struct {
		name     string
		specData map[string]interface{}
		errMsg   string
	}{
		{
			name:     "empty spec",
			specData: map[string]interface{}{},
			errMsg:   "unable to detect OpenAPI version",
		},
		{
			name: "unsupported version",
			specData: map[string]interface{}{
				"openapi": "4.0.0",
			},
			errMsg: "unable to detect OpenAPI version",
		},
		{
			name: "missing servers in 3.0",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]interface{}{
					"title":   "Test",
					"version": "1.0.0",
				},
				"paths": map[string]interface{}{},
			},
			errMsg: "no servers defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(ctx, tt.specData)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestValidate(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)
	ctx := context.Background()

	tests := []struct {
		name     string
		specData map[string]interface{}
		wantErr  bool
	}{
		{
			name: "valid OpenAPI 3.0",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]interface{}{
					"title":   "Test API",
					"version": "1.0.0",
				},
				"servers": []interface{}{
					map[string]interface{}{"url": "https://api.example.com"},
				},
				"paths": map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "valid Swagger 2.0",
			specData: map[string]interface{}{
				"swagger": "2.0",
				"info": map[string]interface{}{
					"title":   "Test API",
					"version": "1.0.0",
				},
				"host":  "api.example.com",
				"paths": map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "invalid spec - missing info",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"paths":   map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name:     "no version",
			specData: map[string]interface{}{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.Validate(ctx, tt.specData)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractBaseURL(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)

	tests := []struct {
		name     string
		specData map[string]interface{}
		want     string
		wantErr  bool
	}{
		{
			name: "OpenAPI 3.0 with single server",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]interface{}{
					"title":   "Test",
					"version": "1.0.0",
				},
				"servers": []interface{}{
					map[string]interface{}{
						"url": "https://api.example.com/v1",
					},
				},
				"paths": map[string]interface{}{},
			},
			want:    "https://api.example.com/v1",
			wantErr: false,
		},
		{
			name: "Swagger 2.0 with https",
			specData: map[string]interface{}{
				"swagger": "2.0",
				"host":    "api.example.com",
				"basePath": "/v1",
				"schemes":  []interface{}{"https"},
			},
			want:    "https://api.example.com/v1",
			wantErr: false,
		},
		{
			name: "Swagger 2.0 with http",
			specData: map[string]interface{}{
				"swagger": "2.0",
				"host":    "localhost:8080",
				"basePath": "/api",
				"schemes":  []interface{}{"http"},
			},
			want:    "http://localhost:8080/api",
			wantErr: false,
		},
		{
			name: "Swagger 2.0 without basePath",
			specData: map[string]interface{}{
				"swagger": "2.0",
				"host":    "api.example.com",
				"schemes": []interface{}{"https"},
			},
			want:    "https://api.example.com",
			wantErr: false,
		},
		{
			name: "Swagger 2.0 missing host",
			specData: map[string]interface{}{
				"swagger": "2.0",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "OpenAPI 3.0 missing servers",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]interface{}{
					"title":   "Test",
					"version": "1.0.0",
				},
				"paths": map[string]interface{}{},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ExtractBaseURL(tt.specData)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestExtractPaths(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)
	ctx := context.Background()

	specData := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":   "Test API",
			"version": "1.0.0",
		},
		"servers": []interface{}{
			map[string]interface{}{"url": "https://api.example.com"},
		},
		"paths": map[string]interface{}{
			"/users": map[string]interface{}{
				"get": map[string]interface{}{
					"operationId": "getUsers",
					"summary":     "Get all users",
					"description": "Retrieves a list of all users",
					"parameters": []interface{}{
						map[string]interface{}{
							"name":        "limit",
							"in":          "query",
							"description": "Max number of results",
							"required":    false,
							"schema": map[string]interface{}{
								"type": "integer",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
					},
				},
				"post": map[string]interface{}{
					"operationId": "createUser",
					"summary":     "Create a user",
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Created",
						},
					},
				},
			},
			"/users/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"operationId": "getUser",
					"summary":     "Get user by ID",
					"parameters": []interface{}{
						map[string]interface{}{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
					},
				},
			},
		},
	}

	result, err := parser.Parse(ctx, specData)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Check paths exist
	assert.Len(t, result.Paths, 2)
	assert.Contains(t, result.Paths, "/users")
	assert.Contains(t, result.Paths, "/users/{id}")

	// Check /users path
	usersPath := result.Paths["/users"]
	assert.NotNil(t, usersPath)
	assert.Equal(t, "/users", usersPath.Path)
	assert.Len(t, usersPath.Operations, 2)

	// Check GET operation
	getOp := usersPath.Operations["GET"]
	assert.NotNil(t, getOp)
	assert.Equal(t, "getUsers", getOp.OperationID)
	assert.Equal(t, "GET", getOp.Method)
	assert.Equal(t, "Get all users", getOp.Summary)
	assert.Equal(t, "Retrieves a list of all users", getOp.Description)
	assert.Len(t, getOp.Parameters, 1)

	// Check parameter
	param := getOp.Parameters[0]
	assert.Equal(t, "limit", param.Name)
	assert.Equal(t, "query", param.In)
	assert.Equal(t, "Max number of results", param.Description)
	assert.False(t, param.Required)
	assert.Equal(t, "integer", param.Type)

	// Check POST operation
	postOp := usersPath.Operations["POST"]
	assert.NotNil(t, postOp)
	assert.Equal(t, "createUser", postOp.OperationID)
	assert.Equal(t, "POST", postOp.Method)

	// Check /users/{id} path
	userPath := result.Paths["/users/{id}"]
	assert.NotNil(t, userPath)
	getUserOp := userPath.Operations["GET"]
	assert.NotNil(t, getUserOp)
	assert.Len(t, getUserOp.Parameters, 1)
	assert.True(t, getUserOp.Parameters[0].Required)
	assert.Equal(t, "path", getUserOp.Parameters[0].In)
}

func TestExtractParameterType(t *testing.T) {
	log := logger.NewNop()
	parser := NewOpenAPIParser(log)

	tests := []struct {
		name     string
		specData map[string]interface{}
		want     string
	}{
		{
			name: "parameter with explicit string schema",
			specData: map[string]interface{}{
				"openapi": "3.0.0",
				"info": map[string]interface{}{
					"title":   "Test",
					"version": "1.0.0",
				},
				"servers": []interface{}{
					map[string]interface{}{"url": "https://api.example.com"},
				},
				"paths": map[string]interface{}{
					"/test": map[string]interface{}{
						"get": map[string]interface{}{
							"operationId": "test",
							"parameters": []interface{}{
								map[string]interface{}{
									"name": "param",
									"in":   "query",
									"schema": map[string]interface{}{
										"type": "string",
									},
								},
							},
							"responses": map[string]interface{}{
								"200": map[string]interface{}{
									"description": "Success",
								},
							},
						},
					},
				},
			},
			want: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := parser.Parse(ctx, tt.specData)
			require.NoError(t, err)

			testPath := result.Paths["/test"]
			getOp := testPath.Operations["GET"]
			if len(getOp.Parameters) > 0 {
				assert.Equal(t, tt.want, getOp.Parameters[0].Type)
			}
		})
	}
}
