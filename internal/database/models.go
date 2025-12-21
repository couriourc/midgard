package database

import (
	"time"

	"gorm.io/gorm"
)

// Collection represents a collection of endpoints
type Collection struct {
	ID            string         `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	Prefix        string         `gorm:"type:varchar(255);not null" json:"prefix"` // External gateway prefix
	BaseURL       string         `gorm:"type:varchar(500);not null" json:"base_url"` // Target base URL
	OpenAPIURL    string         `gorm:"type:varchar(500)" json:"openapi_url"`
	HealthPath    string         `gorm:"type:varchar(255)" json:"health_path"` // Health check path, e.g., /health
	HealthInterval int           `gorm:"default:30" json:"health_interval"` // Health check interval in seconds
	LogEnabled     bool          `gorm:"default:true" json:"log_enabled"`
	LogRolling     bool          `gorm:"default:true" json:"log_rolling"`
	LogMaxEntries  int           `gorm:"default:1000" json:"log_max_entries"`
	CacheEnabled    bool          `gorm:"default:false" json:"cache_enabled"`
	CacheTTL        int           `gorm:"default:300" json:"cache_ttl"` // Cache TTL in seconds
	CacheKeyStrategy string       `gorm:"type:varchar(50);default:'all'" json:"cache_key_strategy"` // "params", "body", "all"
	Active          bool          `gorm:"default:true" json:"active"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Endpoints []Endpoint `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE" json:"endpoints,omitempty"`
}

// Endpoint represents an API endpoint
type Endpoint struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CollectionID string    `gorm:"type:varchar(255);not null;index" json:"collection_id"`
	Path         string    `gorm:"type:varchar(500);not null" json:"path"`
	Method       string    `gorm:"type:varchar(10);not null" json:"method"`
	Summary      string    `gorm:"type:text" json:"summary"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RequestLog represents a request log entry
type RequestLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CollectionID   string    `gorm:"type:varchar(255);not null;index" json:"collection_id"`
	Path           string    `gorm:"type:varchar(500);not null" json:"path"`
	Method         string    `gorm:"type:varchar(10);not null" json:"method"`
	TargetURL      string    `gorm:"type:varchar(1000);not null" json:"target_url"`
	Status         int       `gorm:"not null" json:"status"`
	Duration       int64     `gorm:"not null" json:"duration"` // in milliseconds
	RequestSize    int       `json:"request_size"`
	ResponseSize   int       `json:"response_size"`
	ClientIP       string    `gorm:"type:varchar(50)" json:"client_ip"`
	RequestHeaders string    `gorm:"type:text" json:"request_headers"` // JSON string
	ResponseHeaders string   `gorm:"type:text" json:"response_headers"` // JSON string
	RequestBody    string    `gorm:"type:text" json:"request_body"` // Request body content
	RequestParams  string    `gorm:"type:text" json:"request_params"` // Query parameters (JSON string)
	FromCache      bool      `gorm:"default:false" json:"from_cache"` // Whether response came from cache
	Timestamp      time.Time `gorm:"index" json:"timestamp"`
}

