package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/akmadan/throome/internal/logger"
	"github.com/akmadan/throome/pkg/adapters"
	"go.uber.org/zap"
)

// HealthChecker performs periodic health checks on adapters
type HealthChecker struct {
	interval  time.Duration
	timeout   time.Duration
	threshold int
	running   bool
	mu        sync.RWMutex
	stopChan  chan struct{}
	statuses  map[string]*HealthHistory
}

// HealthHistory tracks health check history for an adapter
type HealthHistory struct {
	ServiceName        string
	ConsecutiveFails   int
	ConsecutiveSuccess int
	LastHealthy        time.Time
	LastUnhealthy      time.Time
	TotalChecks        int64
	FailedChecks       int64
	History            []adapters.HealthStatus
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(interval time.Duration, timeout time.Duration, threshold int) *HealthChecker {
	return &HealthChecker{
		interval:  interval,
		timeout:   timeout,
		threshold: threshold,
		running:   false,
		stopChan:  make(chan struct{}),
		statuses:  make(map[string]*HealthHistory),
	}
}

// Start starts the health checker
func (h *HealthChecker) Start(ctx context.Context, adapterMap map[string]adapters.Adapter) {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return
	}
	h.running = true
	h.mu.Unlock()

	logger.Info("Health checker started", zap.Duration("interval", h.interval))

	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-h.stopChan:
			logger.Info("Health checker stopped")
			return
		case <-ctx.Done():
			logger.Info("Health checker context cancelled")
			return
		case <-ticker.C:
			h.performHealthChecks(ctx, adapterMap)
		}
	}
}

// Stop stops the health checker
func (h *HealthChecker) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.running {
		return
	}

	close(h.stopChan)
	h.running = false
}

// performHealthChecks performs health checks on all adapters
func (h *HealthChecker) performHealthChecks(ctx context.Context, adapterMap map[string]adapters.Adapter) {
	checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	var wg sync.WaitGroup

	for name, adapter := range adapterMap {
		wg.Add(1)
		go func(name string, adapter adapters.Adapter) {
			defer wg.Done()
			h.checkAdapter(checkCtx, name, adapter)
		}(name, adapter)
	}

	wg.Wait()
}

// checkAdapter performs a health check on a single adapter
func (h *HealthChecker) checkAdapter(ctx context.Context, name string, adapter adapters.Adapter) {
	status, err := adapter.HealthCheck(ctx)
	if err != nil {
		logger.Error("Health check failed",
			zap.String("service", name),
			zap.Error(err),
		)
		status = &adapters.HealthStatus{
			Healthy:      false,
			ErrorMessage: err.Error(),
			LastChecked:  time.Now(),
		}
	}

	h.recordHealthStatus(name, status)
}

// recordHealthStatus records a health status
func (h *HealthChecker) recordHealthStatus(name string, status *adapters.HealthStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()

	history, exists := h.statuses[name]
	if !exists {
		history = &HealthHistory{
			ServiceName: name,
			History:     make([]adapters.HealthStatus, 0, 100),
		}
		h.statuses[name] = history
	}

	history.TotalChecks++

	if status.Healthy {
		history.ConsecutiveSuccess++
		history.ConsecutiveFails = 0
		history.LastHealthy = status.LastChecked
	} else {
		history.ConsecutiveFails++
		history.ConsecutiveSuccess = 0
		history.FailedChecks++
		history.LastUnhealthy = status.LastChecked

		// Log if threshold exceeded
		if history.ConsecutiveFails >= h.threshold {
			logger.Warn("Service unhealthy threshold exceeded",
				zap.String("service", name),
				zap.Int("consecutive_fails", history.ConsecutiveFails),
			)
		}
	}

	// Keep last 100 statuses
	history.History = append(history.History, *status)
	if len(history.History) > 100 {
		history.History = history.History[1:]
	}
}

// GetHealthHistory returns health history for a service
func (h *HealthChecker) GetHealthHistory(name string) *HealthHistory {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.statuses[name]
}

// GetAllHealthHistories returns all health histories
func (h *HealthChecker) GetAllHealthHistories() map[string]*HealthHistory {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]*HealthHistory, len(h.statuses))
	for name, history := range h.statuses {
		result[name] = history
	}

	return result
}

// IsHealthy checks if a service is healthy
func (h *HealthChecker) IsHealthy(name string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	history, exists := h.statuses[name]
	if !exists {
		return true // Assume healthy if no history
	}

	return history.ConsecutiveFails < h.threshold
}
