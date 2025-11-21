package cluster

import (
	"time"
)

// Config represents a cluster configuration
type Config struct {
	ClusterID   string                   `yaml:"cluster_id" json:"cluster_id"`
	Name        string                   `yaml:"name" json:"name"`
	Description string                   `yaml:"description,omitempty" json:"description,omitempty"`
	Services    map[string]ServiceConfig `yaml:"services" json:"services"`
	Routing     RoutingConfig            `yaml:"routing,omitempty" json:"routing,omitempty"`
	Health      HealthConfig             `yaml:"health,omitempty" json:"health,omitempty"`
	AI          AIConfig                 `yaml:"ai,omitempty" json:"ai,omitempty"`
	CreatedAt   time.Time                `yaml:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time                `yaml:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// ServiceConfig represents configuration for a single infrastructure service
type ServiceConfig struct {
	Type        string                 `yaml:"type" json:"type"` // postgres, redis, kafka, etc.
	Host        string                 `yaml:"host" json:"host"`
	Port        int                    `yaml:"port" json:"port"`
	Username    string                 `yaml:"username,omitempty" json:"username,omitempty"`
	Password    string                 `yaml:"password,omitempty" json:"password,omitempty"`
	Database    string                 `yaml:"database,omitempty" json:"database,omitempty"`         // For databases
	ContainerID string                 `yaml:"container_id,omitempty" json:"container_id,omitempty"` // Docker container ID (if provisioned by Throome)
	Options     map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`           // Service-specific options
	Pool        PoolConfig             `yaml:"pool,omitempty" json:"pool,omitempty"`
	TLS         TLSConfig              `yaml:"tls,omitempty" json:"tls,omitempty"`
	Weight      int                    `yaml:"weight,omitempty" json:"weight,omitempty"` // For weighted routing
	Replicas    []ReplicaConfig        `yaml:"replicas,omitempty" json:"replicas,omitempty"`
}

// PoolConfig represents connection pool configuration
type PoolConfig struct {
	MinConnections int `yaml:"min_connections,omitempty" json:"min_connections,omitempty"`
	MaxConnections int `yaml:"max_connections,omitempty" json:"max_connections,omitempty"`
	MaxIdleTime    int `yaml:"max_idle_time,omitempty" json:"max_idle_time,omitempty"` // seconds
	MaxLifetime    int `yaml:"max_lifetime,omitempty" json:"max_lifetime,omitempty"`   // seconds
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled            bool   `yaml:"enabled" json:"enabled"`
	CertFile           string `yaml:"cert_file,omitempty" json:"cert_file,omitempty"`
	KeyFile            string `yaml:"key_file,omitempty" json:"key_file,omitempty"`
	CAFile             string `yaml:"ca_file,omitempty" json:"ca_file,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify,omitempty" json:"insecure_skip_verify,omitempty"`
}

// ReplicaConfig represents a replica of a service
type ReplicaConfig struct {
	Host   string `yaml:"host" json:"host"`
	Port   int    `yaml:"port" json:"port"`
	Role   string `yaml:"role,omitempty" json:"role,omitempty"` // primary, replica, readonly
	Weight int    `yaml:"weight,omitempty" json:"weight,omitempty"`
}

// RoutingConfig represents routing strategy configuration
type RoutingConfig struct {
	Strategy        string   `yaml:"strategy" json:"strategy"` // round_robin, weighted, least_connections, ai
	FailoverEnabled bool     `yaml:"failover_enabled" json:"failover_enabled"`
	TimeoutMS       int      `yaml:"timeout_ms,omitempty" json:"timeout_ms,omitempty"`
	RetryAttempts   int      `yaml:"retry_attempts,omitempty" json:"retry_attempts,omitempty"`
	CircuitBreaker  CBConfig `yaml:"circuit_breaker,omitempty" json:"circuit_breaker,omitempty"`
}

// CBConfig represents circuit breaker configuration
type CBConfig struct {
	Enabled          bool `yaml:"enabled" json:"enabled"`
	FailureThreshold int  `yaml:"failure_threshold,omitempty" json:"failure_threshold,omitempty"`
	ResetTimeout     int  `yaml:"reset_timeout,omitempty" json:"reset_timeout,omitempty"` // seconds
}

// HealthConfig represents health check configuration
type HealthConfig struct {
	Enabled   bool `yaml:"enabled" json:"enabled"`
	Interval  int  `yaml:"interval,omitempty" json:"interval,omitempty"`   // seconds
	Timeout   int  `yaml:"timeout,omitempty" json:"timeout,omitempty"`     // seconds
	Threshold int  `yaml:"threshold,omitempty" json:"threshold,omitempty"` // consecutive failures
}

// AIConfig represents AI optimization configuration
type AIConfig struct {
	Enabled        bool     `yaml:"enabled" json:"enabled"`
	Model          string   `yaml:"model,omitempty" json:"model,omitempty"`
	MinDataPoints  int      `yaml:"min_data_points,omitempty" json:"min_data_points,omitempty"`
	UpdateInterval int      `yaml:"update_interval,omitempty" json:"update_interval,omitempty"` // seconds
	Features       []string `yaml:"features,omitempty" json:"features,omitempty"`
}

// DefaultConfig returns a default cluster configuration
func DefaultConfig(clusterID, name string) *Config {
	return &Config{
		ClusterID:   clusterID,
		Name:        name,
		Description: "",
		Services:    make(map[string]ServiceConfig),
		Routing: RoutingConfig{
			Strategy:        "round_robin",
			FailoverEnabled: true,
			TimeoutMS:       5000,
			RetryAttempts:   3,
			CircuitBreaker: CBConfig{
				Enabled:          false,
				FailureThreshold: 5,
				ResetTimeout:     60,
			},
		},
		Health: HealthConfig{
			Enabled:   true,
			Interval:  10,
			Timeout:   5,
			Threshold: 3,
		},
		AI: AIConfig{
			Enabled:        false,
			Model:          "linear_regression",
			MinDataPoints:  100,
			UpdateInterval: 300,
			Features:       []string{"latency", "load", "error_rate"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// DefaultPoolConfig returns default pool configuration
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MinConnections: 2,
		MaxConnections: 10,
		MaxIdleTime:    300,  // 5 minutes
		MaxLifetime:    3600, // 1 hour
	}
}

// Validate validates the cluster configuration
func (c *Config) Validate() error {
	if c.ClusterID == "" {
		return ErrInvalidClusterConfig{Field: "cluster_id", Message: "cannot be empty"}
	}

	if c.Name == "" {
		return ErrInvalidClusterConfig{Field: "name", Message: "cannot be empty"}
	}

	if len(c.Services) == 0 {
		return ErrInvalidClusterConfig{Field: "services", Message: "at least one service is required"}
	}

	for name := range c.Services {
		svc := c.Services[name]
		if err := svc.Validate(); err != nil {
			return ErrInvalidClusterConfig{Field: "services." + name, Message: err.Error()}
		}
	}

	return nil
}

// Validate validates a service configuration
func (s *ServiceConfig) Validate() error {
	if s.Type == "" {
		return ErrInvalidClusterConfig{Field: "type", Message: "cannot be empty"}
	}

	validTypes := map[string]bool{
		"postgres": true,
		"redis":    true,
		"kafka":    true,
		"mongodb":  true,
		"mysql":    true,
		"rabbitmq": true,
	}

	if !validTypes[s.Type] {
		return ErrInvalidClusterConfig{Field: "type", Message: "unsupported service type: " + s.Type}
	}

	if s.Host == "" {
		return ErrInvalidClusterConfig{Field: "host", Message: "cannot be empty"}
	}

	if s.Port < 1 || s.Port > 65535 {
		return ErrInvalidClusterConfig{Field: "port", Message: "must be between 1 and 65535"}
	}

	return nil
}

// ErrInvalidClusterConfig represents a configuration validation error
type ErrInvalidClusterConfig struct {
	Field   string
	Message string
}

func (e ErrInvalidClusterConfig) Error() string {
	return "invalid cluster config [" + e.Field + "]: " + e.Message
}
