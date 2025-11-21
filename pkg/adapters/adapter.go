package adapters

import (
	"context"
	"time"

	"github.com/akmadan/throome/pkg/cluster"
)

// Adapter is the base interface for all infrastructure adapters
type Adapter interface {
	// Connect establishes connection to the infrastructure service
	Connect(ctx context.Context) error

	// Disconnect closes the connection
	Disconnect(ctx context.Context) error

	// Ping checks if the connection is alive
	Ping(ctx context.Context) error

	// HealthCheck performs a health check
	HealthCheck(ctx context.Context) (*HealthStatus, error)

	// GetType returns the adapter type (redis, postgres, kafka, etc.)
	GetType() string

	// GetMetrics returns current adapter metrics
	GetMetrics() *Metrics

	// IsConnected returns connection status
	IsConnected() bool
}

// DatabaseAdapter extends Adapter for database operations
type DatabaseAdapter interface {
	Adapter

	// Execute executes a query/command
	Execute(ctx context.Context, query string, args ...interface{}) (Result, error)

	// Query performs a query and returns rows
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)

	// QueryRow performs a query that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) Row

	// Begin starts a transaction
	Begin(ctx context.Context) (Transaction, error)
}

// CacheAdapter extends Adapter for cache operations
type CacheAdapter interface {
	Adapter

	// Get retrieves a value
	Get(ctx context.Context, key string) (string, error)

	// Set sets a value
	Set(ctx context.Context, key string, value string, expiration time.Duration) error

	// Delete deletes a key
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)

	// Keys returns keys matching a pattern
	Keys(ctx context.Context, pattern string) ([]string, error)

	// TTL returns the time-to-live of a key
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Expire sets expiration on a key
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// QueueAdapter extends Adapter for message queue operations
type QueueAdapter interface {
	Adapter

	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, message []byte) error

	// Subscribe subscribes to a topic
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error

	// Unsubscribe unsubscribes from a topic
	Unsubscribe(ctx context.Context, topic string) error

	// CreateTopic creates a new topic
	CreateTopic(ctx context.Context, topic string, config map[string]interface{}) error

	// DeleteTopic deletes a topic
	DeleteTopic(ctx context.Context, topic string) error

	// ListTopics lists all topics
	ListTopics(ctx context.Context) ([]string, error)
}

// Result represents the result of a database operation
type Result interface {
	RowsAffected() int64
	LastInsertID() int64
}

// Rows represents query result rows
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// Row represents a single query result row
type Row interface {
	Scan(dest ...interface{}) error
}

// Transaction represents a database transaction
type Transaction interface {
	Commit() error
	Rollback() error
	Execute(ctx context.Context, query string, args ...interface{}) (Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
}

// MessageHandler is a function that handles incoming messages
type MessageHandler func(ctx context.Context, message *Message) error

// Message represents a message from a queue
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Timestamp time.Time
	Offset    int64
}

// HealthStatus represents the health status of an adapter
type HealthStatus struct {
	Healthy          bool
	ResponseTime     time.Duration
	ErrorMessage     string
	LastChecked      time.Time
	ConsecutiveFails int
}

// Metrics holds adapter performance metrics
type Metrics struct {
	TotalRequests     int64
	FailedRequests    int64
	SuccessRate       float64
	AverageLatency    time.Duration
	MinLatency        time.Duration
	MaxLatency        time.Duration
	ActiveConnections int
	TotalConnections  int64
	LastRequestTime   time.Time
}

// Factory creates adapters based on service configuration
type Factory struct {
	constructors map[string]AdapterConstructor
}

// AdapterConstructor is a function that creates an adapter
type AdapterConstructor func(config cluster.ServiceConfig) (Adapter, error)

// NewFactory creates a new adapter factory
func NewFactory() *Factory {
	return &Factory{
		constructors: make(map[string]AdapterConstructor),
	}
}

// Register registers an adapter constructor
func (f *Factory) Register(serviceType string, constructor AdapterConstructor) {
	f.constructors[serviceType] = constructor
}

// Create creates an adapter for the given service configuration
func (f *Factory) Create(config cluster.ServiceConfig) (Adapter, error) {
	constructor, exists := f.constructors[config.Type]
	if !exists {
		return nil, ErrAdapterNotFound{Type: config.Type}
	}

	return constructor(config)
}

// ErrAdapterNotFound is returned when an adapter type is not registered
type ErrAdapterNotFound struct {
	Type string
}

func (e ErrAdapterNotFound) Error() string {
	return "adapter not found: " + e.Type
}

// BaseAdapter provides common functionality for adapters
type BaseAdapter struct {
	config    cluster.ServiceConfig
	connected bool
	metrics   *Metrics
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(config cluster.ServiceConfig) *BaseAdapter {
	return &BaseAdapter{
		config:    config,
		connected: false,
		metrics: &Metrics{
			TotalRequests:     0,
			FailedRequests:    0,
			SuccessRate:       100.0,
			AverageLatency:    0,
			MinLatency:        time.Duration(0),
			MaxLatency:        time.Duration(0),
			ActiveConnections: 0,
			TotalConnections:  0,
		},
	}
}

// GetType returns the adapter type
func (b *BaseAdapter) GetType() string {
	return b.config.Type
}

// GetMetrics returns the adapter metrics
func (b *BaseAdapter) GetMetrics() *Metrics {
	return b.metrics
}

// IsConnected returns the connection status
func (b *BaseAdapter) IsConnected() bool {
	return b.connected
}

// SetConnected sets the connection status
func (b *BaseAdapter) SetConnected(connected bool) {
	b.connected = connected
}

// RecordRequest records a request in metrics
func (b *BaseAdapter) RecordRequest(latency time.Duration, success bool) {
	b.metrics.TotalRequests++
	b.metrics.LastRequestTime = time.Now()

	if !success {
		b.metrics.FailedRequests++
	}

	// Update success rate
	b.metrics.SuccessRate = float64(b.metrics.TotalRequests-b.metrics.FailedRequests) / float64(b.metrics.TotalRequests) * 100

	// Update latency metrics
	if b.metrics.MinLatency == 0 || latency < b.metrics.MinLatency {
		b.metrics.MinLatency = latency
	}
	if latency > b.metrics.MaxLatency {
		b.metrics.MaxLatency = latency
	}

	// Calculate rolling average
	b.metrics.AverageLatency = (b.metrics.AverageLatency*time.Duration(b.metrics.TotalRequests-1) + latency) / time.Duration(b.metrics.TotalRequests)
}
