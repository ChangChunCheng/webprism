package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
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

// setupProxyServiceTest creates test environment with all mocks
func setupProxyServiceTest(t *testing.T) (*ProxyService, *mocks.MockSpecRepository, *mocks.MockAuthRepository, *mocks.MockHTTPClient, *mocks.MockHTTPRequestBuilder, *mocks.MockCryptoService) {
	mockSpecRepo := mocks.NewMockSpecRepository(t)
	mockAuthRepo := mocks.NewMockAuthRepository(t)
	mockHTTPClient := mocks.NewMockHTTPClient(t)
	mockRequestBuilder := mocks.NewMockHTTPRequestBuilder(t)
	mockCrypto := mocks.NewMockCryptoService(t)
	log, _ := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
	})

	mockParser := &mocks.MockOpenAPIParser{}
	service := NewProxyService(mockSpecRepo, mockAuthRepo, mockHTTPClient, mockRequestBuilder, mockCrypto, mockParser, log)
	return service, mockSpecRepo, mockAuthRepo, mockHTTPClient, mockRequestBuilder, mockCrypto
}

func TestProxyService_ExecuteProxy(t *testing.T) {
	tests := []struct {
		name           string
		specID         string
		operationID    string
		params         model.ProxyParameters
		setupMocks     func(*mocks.MockSpecRepository, *mocks.MockAuthRepository, *mocks.MockHTTPClient, *mocks.MockHTTPRequestBuilder, *mocks.MockCryptoService)
		wantStatusCode int
		wantErr        bool
		wantErrType    model.ErrorType
		wantErrCode    string
		validate       func(*testing.T, *model.ProxyResponse)
	}{
		{
			name:        "successful proxy without auth",
			specID:      "spec-no-auth",
			operationID: "getUser",
			params: model.ProxyParameters{
				Path: map[string]string{"id": "123"},
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-no-auth", "Test API")
				specRepo.EXPECT().
					FindByID(mock.Anything, "spec-no-auth").
					Return(spec, nil)

				authRepo.EXPECT().
					FindBySpecID(mock.Anything, "spec-no-auth").
					Return(nil, model.ErrAuthConfigNotFound)

				reqBuilder.EXPECT().
					BuildRequest(mock.Anything, mock.MatchedBy(func(params *output.HTTPRequestParams) bool {
						return params.Method == "GET" && params.PathParams["id"] == "123"
					})).
					Return(&http.Request{
						Method: "GET",
						URL:    mustParseURL("https://api.example.com/users/123"),
					}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"id": "123", "name": "John Doe"}`)),
					}, nil)
			},
			wantStatusCode: 200,
			wantErr:        false,
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.Equal(t, "John Doe", resp.Body["name"])
				assert.Equal(t, "123", resp.Body["id"])
				assert.Nil(t, resp.Error)
				assert.Contains(t, resp.Headers, "Content-Type")
			},
		},
		{
			name:        "successful proxy with bearer auth",
			specID:      "spec-bearer",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "456"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-bearer", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-bearer").Return(spec, nil)

				authConfig := &model.AuthConfig{
					SpecID:               "spec-bearer",
					AuthType:             model.AuthTypeBearer,
					EncryptedCredentials: []byte("encrypted-bearer-token"),
				}
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-bearer").Return(authConfig, nil)

				credentials := map[string]string{"token": "test-bearer-token"}
				credJSON, _ := json.Marshal(credentials)
				crypto.EXPECT().Decrypt([]byte("encrypted-bearer-token")).Return(credJSON, nil)

				var capturedAuthHeader map[string]string
				reqBuilder.EXPECT().
					BuildRequest(mock.Anything, mock.Anything).
					Run(func(_ context.Context, params *output.HTTPRequestParams) {
						capturedAuthHeader = params.AuthHeader
					}).
					Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/456")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"authenticated": true}`)),
					}, nil)

				// Validate Bearer token after test
				t.Cleanup(func() {
					assert.Equal(t, "Bearer test-bearer-token", capturedAuthHeader["Authorization"])
				})
			},
			wantStatusCode: 200,
			wantErr:        false,
		},
		{
			name:        "successful proxy with basic auth and correct base64 encoding",
			specID:      "spec-basic",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "789"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-basic", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-basic").Return(spec, nil)

				authConfig := &model.AuthConfig{
					SpecID:               "spec-basic",
					AuthType:             model.AuthTypeBasic,
					EncryptedCredentials: []byte("encrypted-basic-creds"),
				}
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-basic").Return(authConfig, nil)

				credentials := map[string]string{"username": "testuser", "password": "testpass"}
				credJSON, _ := json.Marshal(credentials)
				crypto.EXPECT().Decrypt([]byte("encrypted-basic-creds")).Return(credJSON, nil)

				var capturedAuthHeader map[string]string
				reqBuilder.EXPECT().
					BuildRequest(mock.Anything, mock.Anything).
					Run(func(_ context.Context, params *output.HTTPRequestParams) {
						capturedAuthHeader = params.AuthHeader
					}).
					Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/789")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"authenticated": true}`)),
					}, nil)

				// Validate Base64 encoding after test
				t.Cleanup(func() {
					authHeader := capturedAuthHeader["Authorization"]
					assert.True(t, strings.HasPrefix(authHeader, "Basic "))
					encoded := strings.TrimPrefix(authHeader, "Basic ")
					decoded, err := base64.StdEncoding.DecodeString(encoded)
					require.NoError(t, err)
					assert.Equal(t, "testuser:testpass", string(decoded))
				})
			},
			wantStatusCode: 200,
			wantErr:        false,
		},
		{
			name:        "successful proxy with api key auth",
			specID:      "spec-apikey",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "999"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-apikey", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-apikey").Return(spec, nil)

				authConfig := &model.AuthConfig{
					SpecID:               "spec-apikey",
					AuthType:             model.AuthTypeAPIKey,
					EncryptedCredentials: []byte("encrypted-apikey"),
				}
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-apikey").Return(authConfig, nil)

				credentials := map[string]string{"key": "X-API-Key", "value": "secret-api-key-123"}
				credJSON, _ := json.Marshal(credentials)
				crypto.EXPECT().Decrypt([]byte("encrypted-apikey")).Return(credJSON, nil)

				var capturedAuthHeader map[string]string
				reqBuilder.EXPECT().
					BuildRequest(mock.Anything, mock.Anything).
					Run(func(_ context.Context, params *output.HTTPRequestParams) {
						capturedAuthHeader = params.AuthHeader
					}).
					Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/999")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"authenticated": true}`)),
					}, nil)

				// Validate API Key header
				t.Cleanup(func() {
					assert.Equal(t, "secret-api-key-123", capturedAuthHeader["X-API-Key"])
				})
			},
			wantStatusCode: 200,
			wantErr:        false,
		},
		{
			name:        "target API returns 4xx error",
			specID:      "spec-4xx",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "404"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-4xx", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-4xx").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-4xx").Return(nil, model.ErrAuthConfigNotFound)
				reqBuilder.EXPECT().BuildRequest(mock.Anything, mock.Anything).Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/404")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 404,
						Body:       io.NopCloser(strings.NewReader(`{"error": "User not found"}`)),
					}, nil)
			},
			wantStatusCode: 404,
			wantErr:        false,
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.NotNil(t, resp.Error)
				assert.Equal(t, "TARGET_API_ERROR", resp.Error.Code)
				assert.Equal(t, model.ErrorTypeExternal, resp.Error.Type)
			},
		},
		{
			name:        "target API returns 5xx error",
			specID:      "spec-5xx",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "500"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-5xx", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-5xx").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-5xx").Return(nil, model.ErrAuthConfigNotFound)
				reqBuilder.EXPECT().BuildRequest(mock.Anything, mock.Anything).Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/500")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 500,
						Body:       io.NopCloser(strings.NewReader(`{"error": "Internal server error"}`)),
					}, nil)
			},
			wantStatusCode: 500,
			wantErr:        false,
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.NotNil(t, resp.Error)
				assert.Equal(t, "TARGET_API_ERROR", resp.Error.Code)
			},
		},
		{
			name:        "spec not found",
			specID:      "non-existent-spec",
			operationID: "getUser",
			params:      model.ProxyParameters{},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				specRepo.EXPECT().
					FindByID(mock.Anything, "non-existent-spec").
					Return(nil, model.ErrSpecNotFound)
			},
			wantErr:     true,
			wantErrType: model.ErrorTypeSystem,
			wantErrCode: "SPEC_NOT_FOUND",
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, model.ErrorTypeSystem, resp.Error.Type)
			},
		},
		{
			name:        "operation not found",
			specID:      "spec-no-op",
			operationID: "nonExistentOperation",
			params:      model.ProxyParameters{},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-no-op", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-no-op").Return(spec, nil)
			},
			wantErr:     true,
			wantErrType: model.ErrorTypeClient,
			wantErrCode: "OPERATION_NOT_FOUND",
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Error)
			},
		},
		{
			name:        "auth decryption failure",
			specID:      "spec-decrypt-fail",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "123"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-decrypt-fail", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-decrypt-fail").Return(spec, nil)

				authConfig := &model.AuthConfig{
					SpecID:               "spec-decrypt-fail",
					AuthType:             model.AuthTypeBearer,
					EncryptedCredentials: []byte("corrupted-data"),
				}
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-decrypt-fail").Return(authConfig, nil)
				crypto.EXPECT().Decrypt([]byte("corrupted-data")).Return(nil, errors.New("decryption failed"))
			},
			wantErr:     true,
			wantErrType: model.ErrorTypeSystem,
			wantErrCode: "AUTH_CONFIG_ERROR",
		},
		{
			name:        "http client connection failure",
			specID:      "spec-network-fail",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "123"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-network-fail", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-network-fail").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-network-fail").Return(nil, model.ErrAuthConfigNotFound)
				reqBuilder.EXPECT().BuildRequest(mock.Anything, mock.Anything).Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/123")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(nil, errors.New("connection refused"))
			},
			wantErr: false, // Service returns ProxyResponse with error, not Go error
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Error)
				assert.Equal(t, "TARGET_API_UNREACHABLE", resp.Error.Code)
				assert.Equal(t, model.ErrorTypeExternal, resp.Error.Type)
			},
		},
		{
			name:        "build request failure",
			specID:      "spec-build-fail",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "123"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-build-fail", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-build-fail").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-build-fail").Return(nil, model.ErrAuthConfigNotFound)
				reqBuilder.EXPECT().BuildRequest(mock.Anything, mock.Anything).Return(nil, errors.New("invalid URL"))
			},
			wantErr:     true,
			wantErrType: model.ErrorTypeClient,
			wantErrCode: "INVALID_REQUEST",
		},
		{
			name:        "empty response body",
			specID:      "spec-empty-body",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "123"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-empty-body", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-empty-body").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-empty-body").Return(nil, model.ErrAuthConfigNotFound)
				reqBuilder.EXPECT().BuildRequest(mock.Anything, mock.Anything).Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/123")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 204,
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil)
			},
			wantStatusCode: 204,
			wantErr:        false,
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.NotNil(t, resp.Body)
				assert.Empty(t, resp.Body)
			},
		},
		{
			name:        "non-json response body",
			specID:      "spec-text-body",
			operationID: "getUser",
			params:      model.ProxyParameters{Path: map[string]string{"id": "123"}},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-text-body", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-text-body").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-text-body").Return(nil, model.ErrAuthConfigNotFound)
				reqBuilder.EXPECT().BuildRequest(mock.Anything, mock.Anything).Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/123")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("plain text response")),
					}, nil)
			},
			wantStatusCode: 200,
			wantErr:        false,
			validate: func(t *testing.T, resp *model.ProxyResponse) {
				assert.NotNil(t, resp.Body)
				assert.Equal(t, "plain text response", resp.Body["_raw"])
			},
		},
		{
			name:        "request with query parameters",
			specID:      "spec-query-params",
			operationID: "listUsers",
			params: model.ProxyParameters{
				Query: map[string]string{"limit": "10", "offset": "20"},
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-query-params", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-query-params").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-query-params").Return(nil, model.ErrAuthConfigNotFound)

				reqBuilder.EXPECT().
					BuildRequest(mock.Anything, mock.MatchedBy(func(params *output.HTTPRequestParams) bool {
						return params.QueryParams["limit"] == "10" && params.QueryParams["offset"] == "20"
					})).
					Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users?limit=10&offset=20")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"users": [], "total": 0}`)),
					}, nil)
			},
			wantStatusCode: 200,
			wantErr:        false,
		},
		{
			name:        "request with custom headers",
			specID:      "spec-custom-headers",
			operationID: "getUser",
			params: model.ProxyParameters{
				Path:    map[string]string{"id": "123"},
				Headers: map[string]string{"X-Request-ID": "req-123", "X-Client-Version": "1.0"},
			},
			setupMocks: func(specRepo *mocks.MockSpecRepository, authRepo *mocks.MockAuthRepository, httpClient *mocks.MockHTTPClient, reqBuilder *mocks.MockHTTPRequestBuilder, crypto *mocks.MockCryptoService) {
				spec := testutil.CreateTestSpecWithPaths("spec-custom-headers", "Test API")
				specRepo.EXPECT().FindByID(mock.Anything, "spec-custom-headers").Return(spec, nil)
				authRepo.EXPECT().FindBySpecID(mock.Anything, "spec-custom-headers").Return(nil, model.ErrAuthConfigNotFound)

				reqBuilder.EXPECT().
					BuildRequest(mock.Anything, mock.MatchedBy(func(params *output.HTTPRequestParams) bool {
						return params.Headers["X-Request-ID"] == "req-123" && params.Headers["X-Client-Version"] == "1.0"
					})).
					Return(&http.Request{Method: "GET", URL: mustParseURL("https://api.example.com/users/123")}, nil)

				httpClient.EXPECT().
					Do(mock.Anything, mock.Anything).
					Return(&http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"id": "123"}`)),
					}, nil)
			},
			wantStatusCode: 200,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockSpecRepo, mockAuthRepo, mockHTTPClient, mockRequestBuilder, mockCrypto := setupProxyServiceTest(t)
			tt.setupMocks(mockSpecRepo, mockAuthRepo, mockHTTPClient, mockRequestBuilder, mockCrypto)

			req := testutil.CreateTestProxyRequest(tt.specID, tt.operationID)
			req.Parameters = tt.params

			// Act
			resp, err := service.ExecuteProxy(context.Background(), req)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.NotNil(t, resp)
				if tt.wantErrCode != "" {
					assert.Equal(t, tt.wantErrCode, resp.Error.Code)
				}
				if tt.wantErrType != "" {
					assert.Equal(t, tt.wantErrType, resp.Error.Type)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.wantStatusCode > 0 {
					assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
				}
			}

			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}

// mustParseURL is a helper to parse URLs in tests
func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
