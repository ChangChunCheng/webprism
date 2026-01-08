package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// SwaggerSpec represents a merged OpenAPI/Swagger specification
type SwaggerSpec struct {
	Swagger             string                       `json:"swagger"`
	Info                SwaggerInfo                  `json:"info"`
	Host                string                       `json:"host,omitempty"`
	BasePath            string                       `json:"basePath,omitempty"`
	Schemes             []string                     `json:"schemes"`
	Consumes            []string                     `json:"consumes"`
	Produces            []string                     `json:"produces"`
	Paths               map[string]interface{}       `json:"paths"`
	Definitions         map[string]interface{}       `json:"definitions"`
	Tags                []SwaggerTag                 `json:"tags"`
	SecurityDefinitions map[string]interface{}       `json:"securityDefinitions,omitempty"`
}

// SwaggerInfo represents API information
type SwaggerInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Contact     struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
		URL   string `json:"url,omitempty"`
	} `json:"contact,omitempty"`
}

// SwaggerTag represents an API tag
type SwaggerTag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// mergeSwaggerSpecs merges all swagger.json files into a single specification
func mergeSwaggerSpecs() (*SwaggerSpec, error) {
	merged := &SwaggerSpec{
		Swagger:  "2.0",
		Schemes:  []string{"http", "https"},
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Info: SwaggerInfo{
			Title:       "WEBPRISM API",
			Description: "Web Proxy Request Integration Service Manager - Unified API Gateway with automatic authentication injection",
			Version:     "1.0.0",
		},
		Paths:       make(map[string]interface{}),
		Definitions: make(map[string]interface{}),
		Tags:        []SwaggerTag{},
	}

	// Swagger files to merge (relative to project root)
	swaggerFiles := []string{
		"gen/openapi/v1/spec.swagger.json",
		"gen/openapi/v1/auth.swagger.json",
		"gen/openapi/v1/proxy.swagger.json",
		"gen/openapi/v1/health.swagger.json",
	}

	// Get project root (assuming we're in internal/adapters/http)
	projectRoot, err := getProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to get project root: %w", err)
	}

	for _, file := range swaggerFiles {
		fullPath := filepath.Join(projectRoot, file)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", file, err)
		}

		var spec map[string]interface{}
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %w", file, err)
		}

		// Merge paths
		if paths, ok := spec["paths"].(map[string]interface{}); ok {
			for path, methods := range paths {
				merged.Paths[path] = methods
			}
		}

		// Merge definitions
		if definitions, ok := spec["definitions"].(map[string]interface{}); ok {
			for name, def := range definitions {
				merged.Definitions[name] = def
			}
		}

		// Merge tags
		if tags, ok := spec["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if tagMap, ok := tag.(map[string]interface{}); ok {
					swaggerTag := SwaggerTag{
						Name: tagMap["name"].(string),
					}
					if desc, ok := tagMap["description"].(string); ok {
						swaggerTag.Description = desc
					}
					merged.Tags = append(merged.Tags, swaggerTag)
				}
			}
		}
	}

	return merged, nil
}

// getProjectRoot finds the project root directory
func getProjectRoot() (string, error) {
	// Try to find go.mod file
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("project root not found (no go.mod)")
}

// serveSwaggerUI serves a simple Swagger UI page using CDN
func serveSwaggerUI() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>WEBPRISM API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
    window.onload = function() {
        window.ui = SwaggerUIBundle({
            url: "/swagger.json",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        });
    };
    </script>
</body>
</html>`
		w.Write([]byte(html))
	})
}

// serveSwaggerJSON serves the merged swagger.json
func serveSwaggerJSON() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spec, err := mergeSwaggerSpecs()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate swagger spec: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(spec)
	})
}

// registerSwaggerHandlers registers Swagger UI and JSON handlers
func registerSwaggerHandlers(mux *http.ServeMux) {
	// Serve merged swagger.json
	mux.Handle("/swagger.json", serveSwaggerJSON())

	// Serve Swagger UI
	mux.Handle("/swagger-ui/", serveSwaggerUI())
	mux.Handle("/swagger-ui", serveSwaggerUI())
}
