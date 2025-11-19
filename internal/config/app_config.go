package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// AppConfig holds the application-level configuration
type AppConfig struct {
	Server     ServerConfig     `yaml:"server"`
	Gateway    GatewayConfig    `yaml:"gateway"`
	Dashboard  DashboardConfig  `yaml:"dashboard"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Logging    LoggingConfig    `yaml:"logging"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`  // seconds
	WriteTimeout int    `yaml:"write_timeout"` // seconds
}

// GatewayConfig holds gateway-specific configuration
type GatewayConfig struct {
	ClustersDir       string `yaml:"clusters_dir"`
	MaxConnections    int    `yaml:"max_connections"`
	ConnectionTimeout int    `yaml:"connection_timeout"` // seconds
	EnableAI          bool   `yaml:"enable_ai"`
}

// DashboardConfig holds dashboard configuration
type DashboardConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Enabled            bool   `yaml:"enabled"`
	MetricsPath        string `yaml:"metrics_path"`
	CollectionInterval int    `yaml:"collection_interval"` // seconds
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level       string `yaml:"level"` // debug, info, warn, error
	Development bool   `yaml:"development"`
	OutputPath  string `yaml:"output_path"`
}

// DefaultConfig returns the default application configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         9000,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Gateway: GatewayConfig{
			ClustersDir:       "./clusters",
			MaxConnections:    1000,
			ConnectionTimeout: 10,
			EnableAI:          false,
		},
		Dashboard: DashboardConfig{
			Enabled: true,
			Port:    9001,
			Path:    "/dashboard",
		},
		Monitoring: MonitoringConfig{
			Enabled:            true,
			MetricsPath:        "/metrics",
			CollectionInterval: 10,
		},
		Logging: LoggingConfig{
			Level:       "info",
			Development: false,
			OutputPath:  "stdout",
		},
	}
}

// LoadConfig loads application configuration from a YAML file
func LoadConfig(configPath string) (*AppConfig, error) {
	// Start with defaults
	config := DefaultConfig()

	// If no config file specified, return defaults
	if configPath == "" {
		return config, nil
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and normalize paths
	if !filepath.IsAbs(config.Gateway.ClustersDir) {
		absPath, err := filepath.Abs(config.Gateway.ClustersDir)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve clusters directory: %w", err)
		}
		config.Gateway.ClustersDir = absPath
	}

	return config, nil
}

// Validate validates the configuration
func (c *AppConfig) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Dashboard.Enabled && (c.Dashboard.Port < 1 || c.Dashboard.Port > 65535) {
		return fmt.Errorf("invalid dashboard port: %d", c.Dashboard.Port)
	}

	if c.Gateway.ClustersDir == "" {
		return fmt.Errorf("clusters directory cannot be empty")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}
