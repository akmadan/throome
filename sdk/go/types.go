package throome

import "time"

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// Cluster represents a Throome cluster
type Cluster struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Services  []Service `json:"services,omitempty"`
	CreatedAt string    `json:"created_at"`
}

// Service represents a service in a cluster
type Service struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username,omitempty"`
	Database    string `json:"database,omitempty"`
	Healthy     bool   `json:"healthy"`
	ContainerID string `json:"container_id,omitempty"`
}

// CreateClusterRequest represents a request to create a cluster
type CreateClusterRequest struct {
	Name     string                   `json:"name"`
	Services map[string]ServiceConfig `json:"services"`
}

// ServiceConfig represents service configuration
type ServiceConfig struct {
	Type      string `json:"type"`
	Provision bool   `json:"provision"`          // If true, Throome provisions a new Docker container; if false, connects to existing service
	Host      string `json:"host,omitempty"`     // Required when Provision is false
	Port      int    `json:"port"`               // Required when Provision is false
	Username  string `json:"username,omitempty"` // Required for databases when Provision is false
	Password  string `json:"password,omitempty"` // Required for databases when Provision is false
	Database  string `json:"database,omitempty"` // Required for databases when Provision is false
}

// CreateClusterResponse represents the response from creating a cluster
type CreateClusterResponse struct {
	ClusterID string `json:"cluster_id"`
	Message   string `json:"message"`
}

// ClusterHealthResponse represents cluster health status
type ClusterHealthResponse struct {
	ClusterID string                   `json:"cluster_id"`
	Services  map[string]ServiceHealth `json:"services"`
}

// ServiceHealth represents service health status
type ServiceHealth struct {
	Healthy      bool   `json:"healthy"`
	ResponseTime int64  `json:"response_time"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// MetricsResponse represents cluster metrics
type MetricsResponse struct {
	Requests       int64   `json:"requests"`
	Errors         int64   `json:"errors"`
	AvgResponseMs  float64 `json:"avg_response_ms"`
	P95ResponseMs  float64 `json:"p95_response_ms"`
	ActiveServices int     `json:"active_services"`
}

// ServiceInfo represents detailed service information
type ServiceInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Healthy     bool   `json:"healthy"`
	ContainerID string `json:"container_id,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ActivityLog represents an activity log entry
type ActivityLog struct {
	ID          string            `json:"id"`
	Timestamp   time.Time         `json:"timestamp"`
	ClusterID   string            `json:"cluster_id"`
	ServiceName string            `json:"service_name"`
	ServiceType string            `json:"service_type"`
	Operation   string            `json:"operation"`
	Command     string            `json:"command"`
	Parameters  []interface{}     `json:"parameters"`
	Duration    time.Duration     `json:"duration"`
	Status      string            `json:"status"`
	Response    string            `json:"response"`
	Error       string            `json:"error,omitempty"`
	ClientInfo  map[string]string `json:"client_info,omitempty"`
}

// ActivityFilters represents filters for activity logs
type ActivityFilters struct {
	Limit int
}

// LogOptions represents options for fetching service logs
type LogOptions struct {
	Tail       int
	Timestamps bool
}

// DBQueryRequest represents a database query request
type DBQueryRequest struct {
	Query string        `json:"query"`
	Args  []interface{} `json:"args,omitempty"`
}

// DBQueryResponse represents a database query response
type DBQueryResponse struct {
	Rows []map[string]interface{} `json:"rows"`
}

// CacheGetRequest represents a cache get request
type CacheGetRequest struct {
	Key string `json:"key"`
}

// CacheGetResponse represents a cache get response
type CacheGetResponse struct {
	Value string `json:"value"`
}

// CacheSetRequest represents a cache set request
type CacheSetRequest struct {
	Key        string  `json:"key"`
	Value      string  `json:"value"`
	Expiration float64 `json:"expiration,omitempty"`
}

// CacheDeleteRequest represents a cache delete request
type CacheDeleteRequest struct {
	Key string `json:"key"`
}

// QueuePublishRequest represents a queue publish request
type QueuePublishRequest struct {
	Topic   string `json:"topic"`
	Message []byte `json:"message"`
}
