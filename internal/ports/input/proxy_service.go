package input

import (
	"context"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
)

// ProxyService defines the input port for API proxy operations.
// This interface handles executing proxied requests to target APIs.
type ProxyService interface {
	// ExecuteProxy executes a proxied API request to the target API.
	// It performs the following:
	// 1. Retrieves the API spec by specID
	// 2. Finds the operation by operationID
	// 3. Retrieves and decrypts auth config if exists
	// 4. Builds the HTTP request with parameters and auth
	// 5. Executes the request to the target API
	// 6. Returns the response or error (classified as system/client/external)
	//
	// Returns ProxyResponse with the target API's response, or ProxyError with classified error type.
	ExecuteProxy(ctx context.Context, req *model.ProxyRequest) (*model.ProxyResponse, error)
}
