package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/output"
)

// ProxyService implements the ProxyService input port.
type ProxyService struct {
	specRepo       output.SpecRepository
	authRepo       output.AuthRepository
	httpClient     output.HTTPClient
	requestBuilder output.HTTPRequestBuilder
	crypto         output.CryptoService
	parser         output.OpenAPIParser
	logger         *logger.Logger
}

// NewProxyService creates a new ProxyService instance.
func NewProxyService(
	specRepo output.SpecRepository,
	authRepo output.AuthRepository,
	httpClient output.HTTPClient,
	requestBuilder output.HTTPRequestBuilder,
	crypto output.CryptoService,
	parser output.OpenAPIParser,
	log *logger.Logger,
) *ProxyService {
	return &ProxyService{
		specRepo:       specRepo,
		authRepo:       authRepo,
		httpClient:     httpClient,
		requestBuilder: requestBuilder,
		crypto:         crypto,
		parser:         parser,
		logger:         log,
	}
}

// ExecuteProxy executes a proxied API request to the target API.
func (s *ProxyService) ExecuteProxy(ctx context.Context, req *model.ProxyRequest) (*model.ProxyResponse, error) {
	s.logger.Info("executing proxy request",
		logger.String("spec_id", req.SpecID),
		logger.String("operation_id", req.OperationID),
	)

	// 1. Retrieve the API spec
	spec, err := s.specRepo.FindByID(ctx, req.SpecID)
	if err != nil {
		return s.createErrorResponse(model.ErrorTypeSystem, "SPEC_NOT_FOUND", err.Error()), err
	}

	// Re-parse SpecData to populate Paths if needed
	if spec.SpecData != nil && spec.Paths == nil {
		parsed, err := s.parser.Parse(ctx, spec.SpecData)
		if err != nil {
			s.logger.Error("failed to parse spec data",
				logger.String("spec_id", req.SpecID),
				logger.Any("error", err),
			)
			return s.createErrorResponse(model.ErrorTypeSystem, "SPEC_PARSE_ERROR", err.Error()), err
		}
		spec.Paths = parsed.Paths
	}

	// 2. Find the operation
	operation, err := s.findOperation(spec, req.OperationID)
	if err != nil {
		return s.createErrorResponse(model.ErrorTypeClient, "OPERATION_NOT_FOUND", err.Error()), err
	}

	// 3. Retrieve auth config if exists
	authHeaders, err := s.getAuthHeaders(ctx, req.SpecID)
	if err != nil && err != model.ErrAuthConfigNotFound {
		s.logger.Error("failed to get auth headers",
			logger.String("spec_id", req.SpecID),
			logger.Any("error", err),
		)
		return s.createErrorResponse(model.ErrorTypeSystem, "AUTH_CONFIG_ERROR", err.Error()), err
	}

	// 4. Build the HTTP request
	httpReq, err := s.requestBuilder.BuildRequest(ctx, &output.HTTPRequestParams{
		Method:      operation.Method,
		BaseURL:     spec.BaseURL,
		Path:        operation.Path,
		PathParams:  req.Parameters.Path,
		QueryParams: req.Parameters.Query,
		Headers:     req.Parameters.Headers,
		Body:        req.Parameters.Body,
		AuthHeader:  authHeaders,
	})
	if err != nil {
		s.logger.Error("failed to build HTTP request",
			logger.String("spec_id", req.SpecID),
			logger.Any("error", err),
		)
		return s.createErrorResponse(model.ErrorTypeClient, "INVALID_REQUEST", err.Error()), err
	}

	// 5. Execute the request
	httpResp, err := s.httpClient.Do(ctx, httpReq)
	if err != nil {
		s.logger.Error("failed to execute HTTP request",
			logger.String("spec_id", req.SpecID),
			logger.String("url", httpReq.URL.String()),
			logger.Any("error", err),
		)
		return s.createErrorResponse(
			model.ErrorTypeExternal,
			"TARGET_API_UNREACHABLE",
			fmt.Sprintf("Failed to connect to target API: %v", err),
		), nil // Return nil error because we've classified it as external
	}
	defer httpResp.Body.Close()

	// 6. Parse response body
	body, err := s.parseResponseBody(httpResp)
	if err != nil {
		s.logger.Error("failed to parse response body",
			logger.String("spec_id", req.SpecID),
			logger.Any("error", err),
		)
		return s.createErrorResponse(model.ErrorTypeExternal, "INVALID_RESPONSE", err.Error()), nil
	}

	// 7. Build response
	response := &model.ProxyResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    s.extractHeaders(httpResp),
		Body:       body,
	}

	// 8. Check if target API returned an error
	if httpResp.StatusCode >= 400 {
		response.Error = &model.ProxyError{
			Code:    "TARGET_API_ERROR",
			Message: fmt.Sprintf("Target API returned status %d", httpResp.StatusCode),
			Type:    model.ErrorTypeExternal,
			Details: map[string]interface{}{
				"status_code":   httpResp.StatusCode,
				"response_body": body,
			},
		}
	}

	s.logger.Info("proxy request completed successfully",
		logger.String("spec_id", req.SpecID),
		logger.String("operation_id", req.OperationID),
		logger.Int("status_code", httpResp.StatusCode),
	)

	return response, nil
}

// findOperation finds an operation by operation ID in the spec.
func (s *ProxyService) findOperation(spec *model.APISpec, operationID string) (*model.OperationWithPath, error) {
	for path, pathItem := range spec.Paths {
		for method, operation := range pathItem.Operations {
			if operation.OperationID == operationID {
				return &model.OperationWithPath{
					Operation: *operation,
					Path:      path,
					Method:    method,
				}, nil
			}
		}
	}
	return nil, model.ErrOperationNotFound
}

// getAuthHeaders retrieves and builds authentication headers.
func (s *ProxyService) getAuthHeaders(ctx context.Context, specID string) (map[string]string, error) {
	authConfig, err := s.authRepo.FindBySpecID(ctx, specID)
	if err != nil {
		return nil, err
	}

	// Decrypt credentials
	decrypted, err := s.crypto.Decrypt(authConfig.EncryptedCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt credentials: %w", err)
	}

	var credentials map[string]string
	if err := json.Unmarshal(decrypted, &credentials); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	// Build auth headers based on auth type
	headers := make(map[string]string)
	switch authConfig.AuthType {
	case model.AuthTypeAPIKey:
		key := credentials["key"]
		value := credentials["value"]
		headers[key] = value

	case model.AuthTypeBearer:
		token := credentials["token"]
		headers["Authorization"] = fmt.Sprintf("Bearer %s", token)

	case model.AuthTypeBasic:
		username := credentials["username"]
		password := credentials["password"]
		auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		headers["Authorization"] = fmt.Sprintf("Basic %s", auth)

	default:
		return nil, fmt.Errorf("unsupported auth type: %s", authConfig.AuthType)
	}

	return headers, nil
}

// parseResponseBody parses the HTTP response body.
func (s *ProxyService) parseResponseBody(resp *http.Response) (map[string]interface{}, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// If empty body, return empty map
	if len(bodyBytes) == 0 {
		return make(map[string]interface{}), nil
	}

	// Try to parse as JSON
	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		// If not JSON, return as string
		return map[string]interface{}{
			"_raw": string(bodyBytes),
		}, nil
	}

	return body, nil
}

// extractHeaders extracts headers from HTTP response.
func (s *ProxyService) extractHeaders(resp *http.Response) map[string]string {
	headers := make(map[string]string)
	for key, values := range resp.Header {
		// Join multiple values with comma
		headers[key] = strings.Join(values, ", ")
	}
	return headers
}

// createErrorResponse creates a ProxyResponse with an error.
func (s *ProxyService) createErrorResponse(errorType model.ErrorType, code, message string) *model.ProxyResponse {
	return &model.ProxyResponse{
		StatusCode: 0,
		Headers:    make(map[string]string),
		Body:       make(map[string]interface{}),
		Error: &model.ProxyError{
			Code:    code,
			Message: message,
			Type:    errorType,
			Details: make(map[string]interface{}),
		},
	}
}
