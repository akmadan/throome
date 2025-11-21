package router

import (
	"context"
	"fmt"
	"sync"

	"github.com/akmadan/throome/pkg/adapters"
	"github.com/akmadan/throome/pkg/cluster"
)

// Router handles routing requests to appropriate adapters
type Router struct {
	config   cluster.Config
	adapters map[string]adapters.Adapter
	strategy Strategy
	mu       sync.RWMutex
}

// NewRouter creates a new router for a cluster
func NewRouter(config cluster.Config, adapterMap map[string]adapters.Adapter) *Router {
	router := &Router{
		config:   config,
		adapters: adapterMap,
	}

	// Initialize strategy based on config
	router.strategy = router.createStrategy(config.Routing.Strategy)

	return router
}

// GetAdapter returns an adapter for the given service name
func (r *Router) GetAdapter(serviceName string) (adapters.Adapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, exists := r.adapters[serviceName]
	if !exists {
		return nil, fmt.Errorf("adapter not found for service: %s", serviceName)
	}

	return adapter, nil
}

// Route routes a request to an appropriate adapter using the configured strategy
func (r *Router) Route(ctx context.Context, serviceName string, serviceType string) (adapters.Adapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get all adapters of the requested type
	var candidates []adapters.Adapter
	for _, adapter := range r.adapters {
		if adapter.GetType() == serviceType && adapter.IsConnected() {
			candidates = append(candidates, adapter)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no available adapters for service type: %s", serviceType)
	}

	// If specific service name requested, return it
	if serviceName != "" {
		if adapter, exists := r.adapters[serviceName]; exists && adapter.IsConnected() {
			return adapter, nil
		}
		return nil, fmt.Errorf("service not available: %s", serviceName)
	}

	// Use strategy to select adapter
	selected, err := r.strategy.Select(ctx, candidates)
	if err != nil {
		return nil, fmt.Errorf("failed to select adapter: %w", err)
	}

	return selected, nil
}

// AddAdapter adds a new adapter to the router
func (r *Router) AddAdapter(name string, adapter adapters.Adapter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.adapters[name] = adapter
}

// RemoveAdapter removes an adapter from the router
func (r *Router) RemoveAdapter(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.adapters, name)
}

// GetAllAdapters returns all adapters
func (r *Router) GetAllAdapters() map[string]adapters.Adapter {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy
	result := make(map[string]adapters.Adapter, len(r.adapters))
	for name, adapter := range r.adapters {
		result[name] = adapter
	}

	return result
}

// UpdateStrategy updates the routing strategy
func (r *Router) UpdateStrategy(strategyName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.strategy = r.createStrategy(strategyName)
}

// createStrategy creates a strategy based on the strategy name
func (r *Router) createStrategy(strategyName string) Strategy {
	switch strategyName {
	case "weighted":
		return NewWeightedStrategy()
	case "least_connections":
		return NewLeastConnectionsStrategy()
	case "ai":
		return NewAIStrategy()
	case "round_robin":
		fallthrough
	default:
		return NewRoundRobinStrategy()
	}
}

// HealthCheckAll performs health checks on all adapters
func (r *Router) HealthCheckAll(ctx context.Context) map[string]*adapters.HealthStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make(map[string]*adapters.HealthStatus)

	for name, adapter := range r.adapters {
		status, _ := adapter.HealthCheck(ctx)
		results[name] = status
	}

	return results
}
