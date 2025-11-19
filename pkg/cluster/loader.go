package cluster

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader handles loading and saving cluster configurations
type Loader struct {
	baseDir string
}

// NewLoader creates a new configuration loader
func NewLoader(baseDir string) *Loader {
	return &Loader{
		baseDir: baseDir,
	}
}

// Load loads a cluster configuration from disk
func (l *Loader) Load(clusterID string) (*Config, error) {
	configPath := l.getConfigPath(clusterID)

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cluster config not found: %s", clusterID)
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// Save saves a cluster configuration to disk
func (l *Loader) Save(config *Config) error {
	// Validate first
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Create cluster directory if it doesn't exist
	clusterDir := l.getClusterDir(config.ClusterID)
	if err := os.MkdirAll(clusterDir, 0755); err != nil {
		return fmt.Errorf("failed to create cluster directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	configPath := l.getConfigPath(config.ClusterID)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Delete deletes a cluster configuration from disk
func (l *Loader) Delete(clusterID string) error {
	clusterDir := l.getClusterDir(clusterID)

	// Check if directory exists
	if _, err := os.Stat(clusterDir); os.IsNotExist(err) {
		return fmt.Errorf("cluster not found: %s", clusterID)
	}

	// Remove directory and all contents
	if err := os.RemoveAll(clusterDir); err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}

	return nil
}

// List lists all cluster configurations
func (l *Loader) List() ([]string, error) {
	// Check if base directory exists
	if _, err := os.Stat(l.baseDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(l.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read clusters directory: %w", err)
	}

	var clusterIDs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if config.yaml exists
		configPath := filepath.Join(l.baseDir, entry.Name(), "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			clusterIDs = append(clusterIDs, entry.Name())
		}
	}

	return clusterIDs, nil
}

// Exists checks if a cluster configuration exists
func (l *Loader) Exists(clusterID string) bool {
	configPath := l.getConfigPath(clusterID)
	_, err := os.Stat(configPath)
	return err == nil
}

// getClusterDir returns the directory path for a cluster
func (l *Loader) getClusterDir(clusterID string) string {
	return filepath.Join(l.baseDir, clusterID)
}

// getConfigPath returns the config file path for a cluster
func (l *Loader) getConfigPath(clusterID string) string {
	return filepath.Join(l.getClusterDir(clusterID), "config.yaml")
}

// LoadAll loads all cluster configurations
func (l *Loader) LoadAll() (map[string]*Config, error) {
	clusterIDs, err := l.List()
	if err != nil {
		return nil, err
	}

	configs := make(map[string]*Config)
	for _, id := range clusterIDs {
		config, err := l.Load(id)
		if err != nil {
			// Log error but continue loading other clusters
			continue
		}
		configs[id] = config
	}

	return configs, nil
}
