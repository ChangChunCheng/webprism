package model

import "time"

// HealthCheck 代表一個健康檢查記錄
type HealthCheck struct {
	ID             string       `json:"id" bun:"id,pk"`                          // UUID
	SpecID         string       `json:"spec_id" bun:"spec_id,notnull"`           // 關聯的 API Spec ID
	CheckURL       string       `json:"check_url" bun:"check_url,notnull"`       // 檢查的 URL
	Status         HealthStatus `json:"status" bun:"status,notnull"`             // 健康狀態
	ResponseTimeMs int          `json:"response_time_ms" bun:"response_time_ms"` // 回應時間（毫秒）
	CheckedAt      time.Time    `json:"checked_at" bun:"checked_at,notnull,default:current_timestamp"`
	ErrorMessage   string       `json:"error_message,omitempty" bun:"error_message"` // 錯誤訊息（如果有）
}

// HealthStatus 健康狀態
type HealthStatus string

const (
	HealthStatusUnspecified HealthStatus = ""
	HealthStatusUP          HealthStatus = "UP"
	HealthStatusDown        HealthStatus = "DOWN"
	HealthStatusDegraded    HealthStatus = "DEGRADED"
)
