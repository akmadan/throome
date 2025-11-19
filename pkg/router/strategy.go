package router

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/akshitmadan/throome/pkg/adapters"
)

// Strategy defines the interface for routing strategies
type Strategy interface {
	Select(ctx context.Context, candidates []adapters.Adapter) (adapters.Adapter, error)
	Name() string
}

// RoundRobinStrategy implements round-robin routing
type RoundRobinStrategy struct {
	counter uint64
}

// NewRoundRobinStrategy creates a new round-robin strategy
func NewRoundRobinStrategy() Strategy {
	return &RoundRobinStrategy{counter: 0}
}

// Select selects an adapter using round-robin
func (s *RoundRobinStrategy) Select(ctx context.Context, candidates []adapters.Adapter) (adapters.Adapter, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no adapters available")
	}

	index := atomic.AddUint64(&s.counter, 1) % uint64(len(candidates))
	return candidates[index], nil
}

// Name returns the strategy name
func (s *RoundRobinStrategy) Name() string {
	return "round_robin"
}

// WeightedStrategy implements weighted routing
type WeightedStrategy struct {
	counter uint64
}

// NewWeightedStrategy creates a new weighted strategy
func NewWeightedStrategy() Strategy {
	return &WeightedStrategy{counter: 0}
}

// Select selects an adapter based on weights
func (s *WeightedStrategy) Select(ctx context.Context, candidates []adapters.Adapter) (adapters.Adapter, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no adapters available")
	}

	// For now, fall back to round-robin
	// TODO: Implement actual weighted selection based on adapter metrics
	index := atomic.AddUint64(&s.counter, 1) % uint64(len(candidates))
	return candidates[index], nil
}

// Name returns the strategy name
func (s *WeightedStrategy) Name() string {
	return "weighted"
}

// LeastConnectionsStrategy implements least-connections routing
type LeastConnectionsStrategy struct{}

// NewLeastConnectionsStrategy creates a new least-connections strategy
func NewLeastConnectionsStrategy() Strategy {
	return &LeastConnectionsStrategy{}
}

// Select selects an adapter with the least active connections
func (s *LeastConnectionsStrategy) Select(ctx context.Context, candidates []adapters.Adapter) (adapters.Adapter, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no adapters available")
	}

	var selected adapters.Adapter
	minConnections := int(^uint(0) >> 1) // Max int

	for _, adapter := range candidates {
		metrics := adapter.GetMetrics()
		if metrics.ActiveConnections < minConnections {
			minConnections = metrics.ActiveConnections
			selected = adapter
		}
	}

	if selected == nil {
		return candidates[0], nil
	}

	return selected, nil
}

// Name returns the strategy name
func (s *LeastConnectionsStrategy) Name() string {
	return "least_connections"
}

// AIStrategy implements AI-based routing
type AIStrategy struct {
	counter uint64
}

// NewAIStrategy creates a new AI strategy
func NewAIStrategy() Strategy {
	return &AIStrategy{counter: 0}
}

// Select selects an adapter using AI predictions
func (s *AIStrategy) Select(ctx context.Context, candidates []adapters.Adapter) (adapters.Adapter, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no adapters available")
	}

	// For now, select based on lowest average latency
	// TODO: Integrate with AI engine for predictions
	var selected adapters.Adapter
	minLatency := int64(^uint64(0) >> 1) // Max int64

	for _, adapter := range candidates {
		metrics := adapter.GetMetrics()
		if int64(metrics.AverageLatency) < minLatency && metrics.AverageLatency > 0 {
			minLatency = int64(metrics.AverageLatency)
			selected = adapter
		}
	}

	if selected == nil {
		// Fall back to first adapter if no metrics available
		return candidates[0], nil
	}

	return selected, nil
}

// Name returns the strategy name
func (s *AIStrategy) Name() string {
	return "ai"
}
