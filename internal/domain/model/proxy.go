package model

// ProxyRequest represents a proxy request.
type ProxyRequest struct {
	SpecID      string          `json:"spec_id"`      // API spec ID
	OperationID string          `json:"operation_id"` // Operation ID
	Parameters  ProxyParameters `json:"parameters"`   // Request parameters
}

// ProxyParameters contains all types of parameters for a proxy request.
type ProxyParameters struct {
	Path    map[string]string      `json:"path,omitempty"`    // Path parameters
	Query   map[string]string      `json:"query,omitempty"`   // Query parameters
	Headers map[string]string      `json:"headers,omitempty"` // Headers
	Body    map[string]interface{} `json:"body,omitempty"`    // Request body
}

// ProxyResponse represents a proxy response.
type ProxyResponse struct {
	StatusCode int                    `json:"status_code"`     // HTTP status code
	Headers    map[string]string      `json:"headers"`         // Response headers
	Body       map[string]interface{} `json:"body"`            // Response body
	Error      *ProxyError            `json:"error,omitempty"` // Error information (if any)
}

// ProxyError represents a proxy error.
type ProxyError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Type    ErrorType              `json:"type"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorType represents the error classification.
type ErrorType string

const (
	ErrorTypeSystem   ErrorType = "system"   // WEBPRISM internal error
	ErrorTypeClient   ErrorType = "client"   // Client request error
	ErrorTypeExternal ErrorType = "external" // External API error
)

// OperationWithPath combines an Operation with its path and method.
type OperationWithPath struct {
	Operation
	Path   string
	Method string
}

// Validate validates the proxy request.
func (p *ProxyRequest) Validate() error {
	if p.SpecID == "" {
		return ErrMissingSpecID
	}
	if p.OperationID == "" {
		return ErrMissingOperationID
	}
	return nil
}
