package cluster

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Manager manages the lifecycle of clusters
type Manager struct {
	loader   *Loader
	registry *Registry
	mu       sync.RWMutex
}

// NewManager creates a new cluster manager
func NewManager(baseDir string) *Manager {
	return &Manager{
		loader:   NewLoader(baseDir),
		registry: NewRegistry(),
	}
}

// Create creates a new cluster
func (m *Manager) Create(name string, config *Config) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate cluster ID if not provided
	if config.ClusterID == "" {
		config.ClusterID = generateClusterID()
	}

	// Check if cluster already exists
	if m.loader.Exists(config.ClusterID) {
		return "", fmt.Errorf("cluster already exists: %s", config.ClusterID)
	}

	// Set metadata
	config.Name = name
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	// Validate configuration
	if err := config.Validate(); err != nil {
		return "", fmt.Errorf("invalid configuration: %w", err)
	}

	// Save to disk
	if err := m.loader.Save(config); err != nil {
		return "", fmt.Errorf("failed to save cluster: %w", err)
	}

	// Register in memory
	m.registry.Register(config.ClusterID, config)

	return config.ClusterID, nil
}

// Get retrieves a cluster configuration
func (m *Manager) Get(clusterID string) (*Config, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check registry first
	if config := m.registry.Get(clusterID); config != nil {
		return config, nil
	}

	// Load from disk
	config, err := m.loader.Load(clusterID)
	if err != nil {
		return nil, err
	}

	// Register in memory
	m.registry.Register(clusterID, config)

	return config, nil
}

// Update updates a cluster configuration
func (m *Manager) Update(clusterID string, config *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if cluster exists
	if !m.loader.Exists(clusterID) {
		return fmt.Errorf("cluster not found: %s", clusterID)
	}

	// Ensure cluster ID matches
	config.ClusterID = clusterID
	config.UpdatedAt = time.Now()

	// Validate
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Save to disk
	if err := m.loader.Save(config); err != nil {
		return fmt.Errorf("failed to save cluster: %w", err)
	}

	// Update registry
	m.registry.Register(clusterID, config)

	return nil
}

// Delete deletes a cluster
func (m *Manager) Delete(clusterID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Delete from disk
	if err := m.loader.Delete(clusterID); err != nil {
		return err
	}

	// Unregister from memory
	m.registry.Unregister(clusterID)

	return nil
}

// List lists all clusters
func (m *Manager) List() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.loader.List()
}

// LoadAll loads all clusters into memory
func (m *Manager) LoadAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	configs, err := m.loader.LoadAll()
	if err != nil {
		return err
	}

	for id, config := range configs {
		m.registry.Register(id, config)
	}

	return nil
}

// GetAllConfigs returns all loaded cluster configurations
func (m *Manager) GetAllConfigs() map[string]*Config {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.registry.GetAll()
}

// Exists checks if a cluster exists
func (m *Manager) Exists(clusterID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.loader.Exists(clusterID)
}

// Reload reloads a cluster configuration from disk
func (m *Manager) Reload(clusterID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	config, err := m.loader.Load(clusterID)
	if err != nil {
		return err
	}

	m.registry.Register(clusterID, config)

	return nil
}

// generateClusterID generates a unique cluster ID
func generateClusterID() string {
	// Generate a UUID and take the first 8 characters
	id := uuid.New().String()
	return id[:8]
}
