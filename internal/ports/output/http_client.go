package output

import (
	"context"
	"net/http"
)

// HTTPClient defines the output port for executing HTTP requests to external APIs.
// This interface abstracts the HTTP client to allow for testing and different implementations.
type HTTPClient interface {
	// Do executes an HTTP request and returns the response.
	// This is a thin wrapper around http.Client.Do for dependency injection.
	// The caller is responsible for closing the response body.
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// HTTPRequestBuilder defines the output port for building HTTP requests.
// This interface helps construct HTTP requests with proper headers, auth, and body.
type HTTPRequestBuilder interface {
	// BuildRequest constructs an HTTP request from the given parameters.
	// It handles:
	// - URL construction (base URL + path + query parameters)
	// - Header injection (from spec and auth)
	// - Body serialization (JSON, form-data, etc.)
	// - Auth injection (API key, Bearer token, etc.)
	//
	// Returns a fully constructed *http.Request or an error.
	BuildRequest(ctx context.Context, params *HTTPRequestParams) (*http.Request, error)
}

// HTTPRequestParams contains all parameters needed to build an HTTP request.
type HTTPRequestParams struct {
	// Method is the HTTP method (GET, POST, PUT, DELETE, etc.)
	Method string

	// BaseURL is the target API's base URL
	BaseURL string

	// Path is the API endpoint path (e.g., "/users/{id}")
	Path string

	// PathParams are the path parameters (e.g., {"id": "123"})
	PathParams map[string]string

	// QueryParams are the query parameters (e.g., {"page": "1"})
	QueryParams map[string]string

	// Headers are the request headers
	Headers map[string]string

	// Body is the request body (will be serialized based on Content-Type)
	Body interface{}

	// AuthHeader is the authentication header (if any)
	// Format: "Authorization: Bearer <token>" or "X-API-Key: <key>"
	AuthHeader map[string]string
}
