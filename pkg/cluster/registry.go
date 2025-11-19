package cluster

import (
	"sync"
)

// Registry maintains an in-memory registry of cluster configurations
type Registry struct {
	clusters map[string]*Config
	mu       sync.RWMutex
}

// NewRegistry creates a new cluster registry
func NewRegistry() *Registry {
	return &Registry{
		clusters: make(map[string]*Config),
	}
}

// Register registers a cluster configuration
func (r *Registry) Register(clusterID string, config *Config) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clusters[clusterID] = config
}

// Unregister removes a cluster from the registry
func (r *Registry) Unregister(clusterID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.clusters, clusterID)
}

// Get retrieves a cluster configuration
func (r *Registry) Get(clusterID string) *Config {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.clusters[clusterID]
}

// GetAll returns all registered clusters
func (r *Registry) GetAll() map[string]*Config {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a copy to avoid concurrent modification
	result := make(map[string]*Config, len(r.clusters))
	for id, config := range r.clusters {
		result[id] = config
	}

	return result
}

// Exists checks if a cluster is registered
func (r *Registry) Exists(clusterID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.clusters[clusterID]
	return exists
}

// Count returns the number of registered clusters
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.clusters)
}

// Clear removes all clusters from the registry
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clusters = make(map[string]*Config)
}
