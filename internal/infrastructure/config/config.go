package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	MCP      MCPConfig      `mapstructure:"mcp"`
}

// ServerConfig contains HTTP and gRPC server configuration.
type ServerConfig struct {
	HTTP HTTPServerConfig `mapstructure:"http"`
	GRPC GRPCServerConfig `mapstructure:"grpc"`
}

// HTTPServerConfig contains HTTP server settings.
type HTTPServerConfig struct {
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// GRPCServerConfig contains gRPC server settings.
type GRPCServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig contains PostgreSQL connection settings.
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// SecurityConfig contains security-related settings.
type SecurityConfig struct {
	// EncryptionKey is the AES-256 key for encrypting credentials (32 bytes).
	// This should be set via environment variable: WEBPRISM_ENCRYPTION_KEY
	EncryptionKey string `mapstructure:"encryption_key"`
}

// LoggingConfig contains logging settings.
type LoggingConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, console
	OutputPath string `mapstructure:"output_path"` // stdout, stderr, or file path
}

// MCPConfig contains MCP server settings.
type MCPConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// Load loads configuration from file and environment variables.
// It follows this precedence (highest to lowest):
// 1. Environment variables (WEBPRISM_*)
// 2. Config file (config.yaml)
// 3. Default values
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./deployments")
		v.AddConfigPath(".")
	}

	// Set default values
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found; use defaults and env vars
	}

	// Enable environment variable overrides
	v.SetEnvPrefix("WEBPRISM")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind environment variables for nested keys
	// Viper's AutomaticEnv() doesn't handle nested keys well
	v.BindEnv("security.encryption_key", "WEBPRISM_SECURITY_ENCRYPTION_KEY")
	v.BindEnv("database.password", "WEBPRISM_DATABASE_PASSWORD")

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values.
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.http.port", 8080)
	v.SetDefault("server.http.read_timeout", 30*time.Second)
	v.SetDefault("server.http.write_timeout", 30*time.Second)
	v.SetDefault("server.http.shutdown_timeout", 10*time.Second)
	v.SetDefault("server.grpc.port", 9090)

	// Database defaults (non-sensitive only)
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "webprism")
	// Note: database.password must be set via config.yaml or environment variable
	v.SetDefault("database.database", "webprism")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", 5*time.Minute)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output_path", "stdout")

	// MCP defaults
	v.SetDefault("mcp.enabled", true)
	v.SetDefault("mcp.name", "webprism")
	v.SetDefault("mcp.version", "1.0.0")
}

// Validate validates the configuration values.
func (c *Config) Validate() error {
	// Validate HTTP port
	if c.Server.HTTP.Port < 1 || c.Server.HTTP.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Server.HTTP.Port)
	}

	// Validate gRPC port
	if c.Server.GRPC.Port < 1 || c.Server.GRPC.Port > 65535 {
		return fmt.Errorf("invalid gRPC port: %d", c.Server.GRPC.Port)
	}

	// Validate database config
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate encryption key (must be 32 bytes for AES-256)
	if len(c.Security.EncryptionKey) != 32 {
		return fmt.Errorf("encryption key must be exactly 32 bytes, got %d", len(c.Security.EncryptionKey))
	}

	// Validate logging level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s (must be debug, info, warn, or error)", c.Logging.Level)
	}

	return nil
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}
