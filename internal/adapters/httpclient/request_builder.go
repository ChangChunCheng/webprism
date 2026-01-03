package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/output"
)

// RequestBuilder implements the HTTPRequestBuilder port.
type RequestBuilder struct {
	logger *logger.Logger
}

// NewRequestBuilder creates a new HTTP request builder.
func NewRequestBuilder(log *logger.Logger) *RequestBuilder {
	return &RequestBuilder{
		logger: log,
	}
}

// BuildRequest constructs an HTTP request from the given parameters.
func (b *RequestBuilder) BuildRequest(ctx context.Context, params *output.HTTPRequestParams) (*http.Request, error) {
	// Build URL
	requestURL, err := b.buildURL(params)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Build body
	var body *bytes.Buffer
	if params.Body != nil {
		body, err = b.buildBody(params.Body, params.Headers)
		if err != nil {
			return nil, fmt.Errorf("failed to build body: %w", err)
		}
	}

	// Create request
	var req *http.Request
	if body != nil {
		req, err = http.NewRequestWithContext(ctx, params.Method, requestURL, body)
	} else {
		req, err = http.NewRequestWithContext(ctx, params.Method, requestURL, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	// Add auth headers
	for key, value := range params.AuthHeader {
		req.Header.Set(key, value)
	}

	// Set default Content-Type if not specified and body exists
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	b.logger.Debug("built HTTP request",
		logger.String("method", params.Method),
		logger.String("url", requestURL),
	)

	return req, nil
}

// buildURL constructs the full URL with path and query parameters.
func (b *RequestBuilder) buildURL(params *output.HTTPRequestParams) (string, error) {
	// Replace path parameters
	path := params.Path
	for key, value := range params.PathParams {
		placeholder := fmt.Sprintf("{%s}", key)
		path = strings.ReplaceAll(path, placeholder, value)
	}

	// Parse base URL
	baseURL, err := url.Parse(params.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Join base URL with path
	fullURL := baseURL.JoinPath(path)

	// Add query parameters
	if len(params.QueryParams) > 0 {
		q := fullURL.Query()
		for key, value := range params.QueryParams {
			q.Add(key, value)
		}
		fullURL.RawQuery = q.Encode()
	}

	return fullURL.String(), nil
}

// buildBody serializes the request body based on Content-Type.
func (b *RequestBuilder) buildBody(body interface{}, headers map[string]string) (*bytes.Buffer, error) {
	// Determine Content-Type
	contentType := "application/json"
	if ct, ok := headers["Content-Type"]; ok {
		contentType = ct
	}

	var buf bytes.Buffer

	// Serialize based on Content-Type
	switch {
	case strings.Contains(contentType, "application/json"):
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode JSON body: %w", err)
		}
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		// Convert body to form values
		formData, ok := body.(map[string]string)
		if !ok {
			return nil, fmt.Errorf("form-encoded body must be map[string]string")
		}
		values := url.Values{}
		for key, value := range formData {
			values.Add(key, value)
		}
		buf.WriteString(values.Encode())
	default:
		// For other content types, try JSON encoding
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode body: %w", err)
		}
	}

	return &buf, nil
}
