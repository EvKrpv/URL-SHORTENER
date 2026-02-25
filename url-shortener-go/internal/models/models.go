package models

import (
	"time"
)

// URL представляет собой основную сущность сокращенного URL, хранящуюся в PostgreSQL
type URL struct {
	ID          uint      `json:"id" db:"id"`                                             // Auto-incrementing primary key
	ShortCode   string    `json:"short_code" db:"short_code"`                             // Unique identifier (e.g., "abc123")
	OriginalURL string    `json:"original_url" db:"original_url" validate:"required,url"` // Must be a valid URL
	UserID      *string   `json:"user_id,omitempty" db:"user_id"`                         // Optional user association
	CreatedAt   time.Time `json:"created_at" db:"created_at"`                             // Auto-set on creation
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`                             // Default: created_at + 30 days
	ClickCount  int       `json:"click_count" db:"click_count"`                           // Auto-incremented on redirect
	IsActive    bool      `json:"is_active" db:"is_active"`                               // Can be deactivated
}

// ShortenRequest представляет собой API-данные для создания коротких URL
type ShortenRequest struct {
	URL         string         `json:"url" validate:"required,url"`                          // Required long URL
	CustomAlias *string        `json:"custom_alias,omitempty" validate:"omitempty,alphanum"` // Optional custom short code
	ExpiresIn   *time.Duration `json:"expires_in"`                                           // Optional expiry (e.g., "24h")
}

// ShortenResponse представляет собой API-ответ после сокращения
type ShortenResponse struct {
	ShortURL    string    `json:"short_url"`    // Full short URL (e.g., "https://short.ly/abc123")
	OriginalURL string    `json:"original_url"` // Echoes the submitted URL
	ExpiresAt   time.Time `json:"expires_at"`   // Calculated expiration timestamp
	ClickCount  int       `json:"click_count"`  // Always 0 initially
}

// URLStats представляет собой аналитические данные для короткого URL
type URLStats struct {
	ShortCode     string         `json:"short_code"`
	TotalClicks   int            `json:"total_clicks"`
	LastClickedAt time.Time      `json:"last_clicked_at"`
	BrowserStats  map[string]int `json:"browser_stats,omitempty"` // Browser breakdown
	CountryStats  map[string]int `json:"country_stats,omitempty"` // Geographic breakdown
}

// User представляет собой учетную запись (будущая реализация аутентификации)
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email" validate:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	APIKey    string    `json:"-" db:"api_key"` // Never serialized to JSON
}

// ErrorResponse стандартизированные сообщения об ошибках
type ErrorResponse struct {
	Error     string `json:"error"`
	Code      int    `json:"code"`
	Details   string `json:"details,omitempty"`
	RequestID string `json:"request_id,omitempty"` // For correlating logs
}

// Струкутра кэша для Redis.
type CachedURL struct {
	OriginalURL string `json:"original_url"`
	ExpiresAt   int64  `json:"expires_at"` // Unix timestamp
	Metadata    string `json:"metadata"`   // Reserved for future use
}

type RedisConfig struct {
	Address  string // Redis server address (e.g., "localhost:6379")
	Password string // Authentication password (empty if none)
	DB       int    // Redis database number (default: 0)
}

type Config struct {
	Postgres PostgresURL
	Redis    RedisURL
	Server   ServerConfig
}

type PostgresURL struct {
	URL string // postgres://user:pass@host:port/db?sslmode=disable
}

type RedisURL struct {
	URL string // redis://:pass@host:port/db
}

type ServerConfig struct {
	Port        string
	Environment string
	Timeout     time.Duration
}
