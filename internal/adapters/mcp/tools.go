package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/input"
)

// Tools implements MCP tools for WEBPRISM.
type Tools struct {
	specService   input.SpecService
	authService   input.AuthService
	proxyService  input.ProxyService
	healthService input.HealthService
	logger        *logger.Logger
}

// NewTools creates a new Tools instance.
func NewTools(
	specService input.SpecService,
	authService input.AuthService,
	proxyService input.ProxyService,
	healthService input.HealthService,
	log *logger.Logger,
) *Tools {
	return &Tools{
		specService:   specService,
		authService:   authService,
		proxyService:  proxyService,
		healthService: healthService,
		logger:        log,
	}
}

// RegisterTools registers all MCP tools with the server.
func (t *Tools) RegisterTools(s *server.MCPServer) {
	// Upload API Spec
	s.AddTool(mcp.Tool{
		Name:        "upload_api_spec",
		Description: "Upload an OpenAPI specification (Swagger 2.0 or OpenAPI 3.0)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the API",
				},
				"version": map[string]interface{}{
					"type":        "string",
					"description": "Version of the API",
				},
				"spec_data": map[string]interface{}{
					"type":        "object",
					"description": "The OpenAPI specification as JSON object",
				},
			},
			Required: []string{"name", "version", "spec_data"},
		},
	}, t.handleUploadSpec)

	// List API Specs
	s.AddTool(mcp.Tool{
		Name:        "list_api_specs",
		Description: "List all uploaded API specifications",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of specs to return (default: 10)",
				},
			},
		},
	}, t.handleListSpecs)

	// Get API Spec
	s.AddTool(mcp.Tool{
		Name:        "get_api_spec",
		Description: "Get detailed information about a specific API specification including all available operations",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"spec_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the API specification",
				},
			},
			Required: []string{"spec_id"},
		},
	}, t.handleGetSpec)

	// Set Auth
	s.AddTool(mcp.Tool{
		Name:        "set_auth",
		Description: "Set authentication configuration for an API",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"spec_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the API specification",
				},
				"auth_type": map[string]interface{}{
					"type":        "string",
					"description": "Authentication type (api_key, bearer, basic, oauth2)",
					"enum":        []string{"api_key", "bearer", "basic", "oauth2"},
				},
				"credentials": map[string]interface{}{
					"type":        "object",
					"description": "Credentials (format depends on auth_type)",
				},
			},
			Required: []string{"spec_id", "auth_type", "credentials"},
		},
	}, t.handleSetAuth)

	// Call API
	s.AddTool(mcp.Tool{
		Name:        "call_api",
		Description: "Execute an API call through WEBPRISM proxy",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"spec_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the API specification",
				},
				"operation_id": map[string]interface{}{
					"type":        "string",
					"description": "The operation ID to execute",
				},
				"parameters": map[string]interface{}{
					"type":        "object",
					"description": "Request parameters (path, query, headers, body)",
				},
			},
			Required: []string{"spec_id", "operation_id"},
		},
	}, t.handleCallAPI)

	// Health Check
	s.AddTool(mcp.Tool{
		Name:        "health_check",
		Description: "Check the health status of a target API",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"spec_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the API specification",
				},
			},
			Required: []string{"spec_id"},
		},
	}, t.handleHealthCheck)
}

// handleUploadSpec handles the upload_api_spec tool.
func (t *Tools) handleUploadSpec(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name, _ := args["name"].(string)
	version, _ := args["version"].(string)
	specData, _ := args["spec_data"].(map[string]interface{})

	spec, err := t.specService.UploadSpec(ctx, name, version, specData)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to upload spec: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(spec, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("API spec uploaded successfully:\n%s", string(result)),
			},
		},
	}, nil
}

// handleListSpecs handles the list_api_specs tool.
func (t *Tools) handleListSpecs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	specs, err := t.specService.ListSpecs(ctx, limit, 0)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to list specs: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(specs, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Found %d API specs:\n%s", len(specs), string(result)),
			},
		},
	}, nil
}

// handleGetSpec handles the get_api_spec tool.
func (t *Tools) handleGetSpec(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	specID, _ := args["spec_id"].(string)

	spec, err := t.specService.GetSpec(ctx, specID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to get spec: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(spec, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("API spec details:\n%s", string(result)),
			},
		},
	}, nil
}

// handleSetAuth handles the set_auth tool.
func (t *Tools) handleSetAuth(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	specID, _ := args["spec_id"].(string)
	authTypeStr, _ := args["auth_type"].(string)
	credentials, _ := args["credentials"].(map[string]interface{})

	// Convert credentials to map[string]string
	creds := make(map[string]string)
	for k, v := range credentials {
		creds[k] = fmt.Sprintf("%v", v)
	}

	authType := model.AuthType(authTypeStr)
	_, err := t.authService.SetAuthConfig(ctx, specID, authType, creds)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to set auth: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Authentication configured successfully for spec %s", specID),
			},
		},
	}, nil
}

// handleCallAPI handles the call_api tool.
func (t *Tools) handleCallAPI(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	specID, _ := args["spec_id"].(string)
	operationID, _ := args["operation_id"].(string)
	parameters, _ := args["parameters"].(map[string]interface{})

	// Parse parameters
	params := model.ProxyParameters{}
	if p, ok := parameters["path"].(map[string]interface{}); ok {
		params.Path = make(map[string]string)
		for k, v := range p {
			params.Path[k] = fmt.Sprintf("%v", v)
		}
	}
	if q, ok := parameters["query"].(map[string]interface{}); ok {
		params.Query = make(map[string]string)
		for k, v := range q {
			params.Query[k] = fmt.Sprintf("%v", v)
		}
	}
	if h, ok := parameters["headers"].(map[string]interface{}); ok {
		params.Headers = make(map[string]string)
		for k, v := range h {
			params.Headers[k] = fmt.Sprintf("%v", v)
		}
	}
	if b, ok := parameters["body"].(map[string]interface{}); ok {
		params.Body = b
	}

	req := &model.ProxyRequest{
		SpecID:      specID,
		OperationID: operationID,
		Parameters:  params,
	}

	resp, err := t.proxyService.ExecuteProxy(ctx, req)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to execute proxy: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(resp, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("API call completed:\n%s", string(result)),
			},
		},
	}, nil
}

// handleHealthCheck handles the health_check tool.
func (t *Tools) handleHealthCheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	specID, _ := args["spec_id"].(string)

	healthCheck, err := t.healthService.CheckHealth(ctx, specID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Failed to check health: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	result, _ := json.MarshalIndent(healthCheck, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Health check result:\n%s", string(result)),
			},
		},
	}, nil
}
