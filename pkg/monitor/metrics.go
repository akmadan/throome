package monitor

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Collector collects and stores metrics
type Collector struct {
	// Prometheus metrics
	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	errorTotal      *prometheus.CounterVec
	activeConns     *prometheus.GaugeVec

	// Custom metrics storage
	clusterMetrics map[string]*ClusterMetrics
	mu             sync.RWMutex
}

// ClusterMetrics holds metrics for a cluster
type ClusterMetrics struct {
	ClusterID      string
	ServiceMetrics map[string]*ServiceMetrics
	TotalRequests  int64
	FailedRequests int64
	AverageLatency time.Duration
	LastUpdated    time.Time
}

// ServiceMetrics holds metrics for a service
type ServiceMetrics struct {
	ServiceName       string
	ServiceType       string
	TotalRequests     int64
	FailedRequests    int64
	SuccessRate       float64
	AverageLatency    time.Duration
	MinLatency        time.Duration
	MaxLatency        time.Duration
	P95Latency        time.Duration
	P99Latency        time.Duration
	ActiveConnections int
	HealthStatus      string
	LastRequestTime   time.Time
	Errors            []string
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{
		requestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "throome_requests_total",
				Help: "Total number of requests",
			},
			[]string{"cluster_id", "service", "type"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "throome_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"cluster_id", "service", "type"},
		),
		errorTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "throome_errors_total",
				Help: "Total number of errors",
			},
			[]string{"cluster_id", "service", "type", "error_type"},
		),
		activeConns: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "throome_active_connections",
				Help: "Number of active connections",
			},
			[]string{"cluster_id", "service", "type"},
		),
		clusterMetrics: make(map[string]*ClusterMetrics),
	}
}

// RecordRequest records a request metric
func (c *Collector) RecordRequest(clusterID, service, serviceType string, duration time.Duration, success bool) {
	c.requestTotal.WithLabelValues(clusterID, service, serviceType).Inc()
	c.requestDuration.WithLabelValues(clusterID, service, serviceType).Observe(duration.Seconds())

	if !success {
		c.errorTotal.WithLabelValues(clusterID, service, serviceType, "unknown").Inc()
	}

	// Update custom metrics
	c.updateServiceMetrics(clusterID, service, serviceType, duration, success)
}

// RecordError records an error metric
func (c *Collector) RecordError(clusterID, service, serviceType, errorType string) {
	c.errorTotal.WithLabelValues(clusterID, service, serviceType, errorType).Inc()
}

// SetActiveConnections sets the active connections gauge
func (c *Collector) SetActiveConnections(clusterID, service, serviceType string, count int) {
	c.activeConns.WithLabelValues(clusterID, service, serviceType).Set(float64(count))
}

// updateServiceMetrics updates custom service metrics
func (c *Collector) updateServiceMetrics(clusterID, service, serviceType string, duration time.Duration, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get or create cluster metrics
	cluster, exists := c.clusterMetrics[clusterID]
	if !exists {
		cluster = &ClusterMetrics{
			ClusterID:      clusterID,
			ServiceMetrics: make(map[string]*ServiceMetrics),
			LastUpdated:    time.Now(),
		}
		c.clusterMetrics[clusterID] = cluster
	}

	// Get or create service metrics
	svc, exists := cluster.ServiceMetrics[service]
	if !exists {
		svc = &ServiceMetrics{
			ServiceName:  service,
			ServiceType:  serviceType,
			MinLatency:   duration,
			MaxLatency:   duration,
			HealthStatus: "healthy",
		}
		cluster.ServiceMetrics[service] = svc
	}

	// Update metrics
	svc.TotalRequests++
	if !success {
		svc.FailedRequests++
	}

	// Update success rate
	svc.SuccessRate = float64(svc.TotalRequests-svc.FailedRequests) / float64(svc.TotalRequests) * 100

	// Update latency metrics
	if duration < svc.MinLatency {
		svc.MinLatency = duration
	}
	if duration > svc.MaxLatency {
		svc.MaxLatency = duration
	}

	// Calculate rolling average
	svc.AverageLatency = (svc.AverageLatency*time.Duration(svc.TotalRequests-1) + duration) / time.Duration(svc.TotalRequests)

	svc.LastRequestTime = time.Now()
	cluster.LastUpdated = time.Now()
}

// GetClusterMetrics returns metrics for a cluster
func (c *Collector) GetClusterMetrics(clusterID string) *ClusterMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.clusterMetrics[clusterID]
}

// GetAllMetrics returns all cluster metrics
func (c *Collector) GetAllMetrics() map[string]*ClusterMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy
	result := make(map[string]*ClusterMetrics, len(c.clusterMetrics))
	for id, metrics := range c.clusterMetrics {
		result[id] = metrics
	}

	return result
}

// GetServiceMetrics returns metrics for a specific service
func (c *Collector) GetServiceMetrics(clusterID, service string) *ServiceMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cluster, exists := c.clusterMetrics[clusterID]; exists {
		return cluster.ServiceMetrics[service]
	}

	return nil
}

// Clear clears all metrics
func (c *Collector) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.clusterMetrics = make(map[string]*ClusterMetrics)
}
