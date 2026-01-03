package parser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ChangChunCheng/webprism/internal/domain/model"
	"github.com/ChangChunCheng/webprism/internal/infrastructure/logger"
	"github.com/ChangChunCheng/webprism/internal/ports/output"
)

// OpenAPIParser implements the OpenAPIParser port.
// It supports both OpenAPI 2.0 (Swagger) and 3.0+ specifications.
type OpenAPIParser struct {
	logger *logger.Logger
}

// NewOpenAPIParser creates a new OpenAPIParser instance.
func NewOpenAPIParser(log *logger.Logger) *OpenAPIParser {
	return &OpenAPIParser{
		logger: log,
	}
}

// Parse parses an OpenAPI specification and extracts structured information.
func (p *OpenAPIParser) Parse(ctx context.Context, specData map[string]interface{}) (*output.ParsedSpec, error) {
	// Determine the OpenAPI version
	version, err := p.detectVersion(specData)
	if err != nil {
		return nil, err
	}

	p.logger.Debug("detected OpenAPI version",
		logger.String("version", version),
	)

	// Parse based on version
	var doc *openapi3.T
	switch version {
	case "2.0":
		doc, err = p.parseV2(specData)
	case "3.0", "3.1":
		doc, err = p.parseV3(specData)
	default:
		return nil, fmt.Errorf("unsupported OpenAPI version: %s", version)
	}

	if err != nil {
		return nil, err
	}

	// Extract base URL
	baseURL, err := p.extractBaseURLFromV3(doc)
	if err != nil {
		return nil, err
	}

	// Extract paths
	paths, err := p.extractPaths(doc)
	if err != nil {
		return nil, err
	}

	// Extract info
	info := &output.SpecInfo{
		Title:       doc.Info.Title,
		Description: doc.Info.Description,
		Version:     doc.Info.Version,
	}

	return &output.ParsedSpec{
		Version: version,
		BaseURL: baseURL,
		Paths:   paths,
		Info:    info,
	}, nil
}

// Validate validates an OpenAPI specification for correctness.
func (p *OpenAPIParser) Validate(ctx context.Context, specData map[string]interface{}) error {
	version, err := p.detectVersion(specData)
	if err != nil {
		return err
	}

	switch version {
	case "2.0":
		_, err = p.parseV2(specData)
	case "3.0", "3.1":
		_, err = p.parseV3(specData)
	default:
		return fmt.Errorf("unsupported OpenAPI version: %s", version)
	}

	return err
}

// ExtractBaseURL extracts the base URL from the spec.
func (p *OpenAPIParser) ExtractBaseURL(specData map[string]interface{}) (string, error) {
	version, err := p.detectVersion(specData)
	if err != nil {
		return "", err
	}

	switch version {
	case "2.0":
		return p.extractBaseURLFromV2(specData)
	case "3.0", "3.1":
		doc, err := p.parseV3(specData)
		if err != nil {
			return "", err
		}
		return p.extractBaseURLFromV3(doc)
	default:
		return "", fmt.Errorf("unsupported OpenAPI version: %s", version)
	}
}

// detectVersion determines the OpenAPI version from the spec data.
func (p *OpenAPIParser) detectVersion(specData map[string]interface{}) (string, error) {
	// Check for OpenAPI 3.x
	if openapi, ok := specData["openapi"].(string); ok {
		if len(openapi) >= 3 && openapi[:3] == "3.0" {
			return "3.0", nil
		}
		if len(openapi) >= 3 && openapi[:3] == "3.1" {
			return "3.1", nil
		}
	}

	// Check for Swagger 2.0
	if swagger, ok := specData["swagger"].(string); ok {
		if swagger == "2.0" {
			return "2.0", nil
		}
	}

	return "", fmt.Errorf("unable to detect OpenAPI version from spec")
}

// parseV2 parses an OpenAPI 2.0 spec and converts it to 3.0.
func (p *OpenAPIParser) parseV2(specData map[string]interface{}) (*openapi3.T, error) {
	// Marshal spec data to JSON
	data, err := json.Marshal(specData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec data: %w", err)
	}

	// Unmarshal into OpenAPI 2.0 struct
	var v2Doc openapi2.T
	if err := json.Unmarshal(data, &v2Doc); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI 2.0 spec: %w", err)
	}

	// Convert to OpenAPI 3.0
	v3Doc, err := openapi2conv.ToV3(&v2Doc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert OpenAPI 2.0 to 3.0: %w", err)
	}

	p.logger.Debug("converted OpenAPI 2.0 to 3.0")

	return v3Doc, nil
}

// parseV3 parses an OpenAPI 3.0 spec.
func (p *OpenAPIParser) parseV3(specData map[string]interface{}) (*openapi3.T, error) {
	// Marshal spec data to JSON
	data, err := json.Marshal(specData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec data: %w", err)
	}

	// Load OpenAPI 3.0 spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI 3.0 spec: %w", err)
	}

	// Validate the spec
	if err := doc.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI 3.0 spec: %w", err)
	}

	return doc, nil
}

// extractBaseURLFromV2 extracts base URL from OpenAPI 2.0 spec.
func (p *OpenAPIParser) extractBaseURLFromV2(specData map[string]interface{}) (string, error) {
	// Get scheme (default to https)
	scheme := "https"
	if schemes, ok := specData["schemes"].([]interface{}); ok && len(schemes) > 0 {
		if s, ok := schemes[0].(string); ok {
			scheme = s
		}
	}

	// Get host
	host, ok := specData["host"].(string)
	if !ok || host == "" {
		return "", fmt.Errorf("missing 'host' field in OpenAPI 2.0 spec")
	}

	// Get basePath (default to /)
	basePath := ""
	if bp, ok := specData["basePath"].(string); ok {
		basePath = bp
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, basePath), nil
}

// extractBaseURLFromV3 extracts base URL from OpenAPI 3.0 spec.
func (p *OpenAPIParser) extractBaseURLFromV3(doc *openapi3.T) (string, error) {
	if len(doc.Servers) == 0 {
		return "", fmt.Errorf("no servers defined in OpenAPI 3.0 spec")
	}

	// Return the first server URL
	return doc.Servers[0].URL, nil
}

// extractPaths converts OpenAPI paths to domain PathItems.
func (p *OpenAPIParser) extractPaths(doc *openapi3.T) (map[string]*model.PathItem, error) {
	paths := make(map[string]*model.PathItem)

	for path, pathItem := range doc.Paths.Map() {
		domainPath := &model.PathItem{
			Path:       path,
			Operations: make(map[string]*model.Operation),
		}

		// Extract operations for each HTTP method
		for method, operation := range pathItem.Operations() {
			if operation == nil {
				continue
			}

			domainOp := &model.Operation{
				OperationID: operation.OperationID,
				Method:      method,
				Summary:     operation.Summary,
				Description: operation.Description,
				Parameters:  p.extractParameters(operation),
			}

			domainPath.Operations[method] = domainOp
		}

		paths[path] = domainPath
	}

	return paths, nil
}

// extractParameters converts OpenAPI parameters to domain Parameters.
func (p *OpenAPIParser) extractParameters(operation *openapi3.Operation) []*model.Parameter {
	var params []*model.Parameter

	for _, paramRef := range operation.Parameters {
		if paramRef.Value == nil {
			continue
		}

		param := paramRef.Value
		domainParam := &model.Parameter{
			Name:        param.Name,
			In:          param.In,
			Description: param.Description,
			Required:    param.Required,
			Type:        p.extractParameterType(param.Schema),
		}

		params = append(params, domainParam)
	}

	return params
}

// extractParameterType extracts the type from a parameter schema.
func (p *OpenAPIParser) extractParameterType(schema *openapi3.SchemaRef) string {
	if schema == nil || schema.Value == nil {
		return "string"
	}

	types := schema.Value.Type.Slice()
	if len(types) == 0 {
		return "string"
	}

	return types[0]
}
