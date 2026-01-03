package model

import "time"

// APISpec 代表一個 API 規格
type APISpec struct {
	ID        string                 `json:"id" bun:"id,pk"`                       // UUID
	Name      string                 `json:"name" bun:"name,notnull"`              // API 名稱
	Version   string                 `json:"version" bun:"version"`                // 版本
	BaseURL   string                 `json:"base_url" bun:"base_url,notnull"`      // 基礎 URL
	SpecData  map[string]interface{} `json:"spec_data" bun:"spec_data,type:jsonb"` // 完整的 OpenAPI 規格（JSONB）
	CreatedAt time.Time              `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time              `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`

	// 關聯
	AuthConfig *AuthConfig `json:"auth_config,omitempty" bun:"rel:has-one,join:id=spec_id"`

	// 解析後的資料（僅用於記憶體，從 SpecData 解析而來）
	Paths map[string]*PathItem `json:"paths,omitempty" bun:"-"`
}

// PathItem represents an API path with multiple operations.
type PathItem struct {
	Path       string
	Operations map[string]*Operation // HTTP method -> Operation
}

// Operation represents a single HTTP operation on a path.
type Operation struct {
	OperationID string
	Method      string // GET, POST, PUT, DELETE, etc.
	Summary     string
	Description string
	Parameters  []*Parameter
	RequestBody *RequestBody
	Responses   map[string]*Response
}

// Parameter represents a parameter.
type Parameter struct {
	Name        string
	In          string // path, query, header, cookie
	Description string
	Required    bool
	Type        string // string, integer, boolean, etc.
	Schema      *Schema
}

// Schema 代表參數的 Schema
type Schema struct {
	Type       string
	Format     string
	Items      *Schema            // for array type
	Properties map[string]*Schema // for object type
}

// RequestBody 代表請求 Body
type RequestBody struct {
	Description string
	Required    bool
	Content     map[string]*MediaType // media type -> content
}

// MediaType 代表媒體類型
type MediaType struct {
	Schema *Schema
}

// Response 代表回應
type Response struct {
	Description string
	Content     map[string]*MediaType
}

// SpecFormat 規格格式
type SpecFormat string

const (
	SpecFormatJSON SpecFormat = "json"
	SpecFormatYAML SpecFormat = "yaml"
)

// SpecVersion OpenAPI 版本
type SpecVersion string

const (
	SpecVersion2 SpecVersion = "2.0" // Swagger 2.0
	SpecVersion3 SpecVersion = "3.0" // OpenAPI 3.0
)
