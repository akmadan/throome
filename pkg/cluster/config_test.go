package cluster

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	clusterID := "test-01"
	name := "Test Cluster"

	config := DefaultConfig(clusterID, name)

	if config.ClusterID != clusterID {
		t.Errorf("Expected cluster ID %s, got %s", clusterID, config.ClusterID)
	}

	if config.Name != name {
		t.Errorf("Expected name %s, got %s", name, config.Name)
	}

	// Check defaults
	if config.Routing.Strategy != "round_robin" {
		t.Errorf("Expected default strategy round_robin, got %s", config.Routing.Strategy)
	}

	if !config.Routing.FailoverEnabled {
		t.Error("Expected failover to be enabled by default")
	}

	if !config.Health.Enabled {
		t.Error("Expected health checks to be enabled by default")
	}

	if config.AI.Enabled {
		t.Error("Expected AI to be disabled by default")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				ClusterID: "test-01",
				Name:      "Test",
				Services: map[string]ServiceConfig{
					"cache": {
						Type: "redis",
						Host: "localhost",
						Port: 6379,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing cluster ID",
			config: &Config{
				Name: "Test",
				Services: map[string]ServiceConfig{
					"cache": {
						Type: "redis",
						Host: "localhost",
						Port: 6379,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing name",
			config: &Config{
				ClusterID: "test-01",
				Services: map[string]ServiceConfig{
					"cache": {
						Type: "redis",
						Host: "localhost",
						Port: 6379,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no services",
			config: &Config{
				ClusterID: "test-01",
				Name:      "Test",
				Services:  map[string]ServiceConfig{},
			},
			wantErr: true,
		},
		{
			name: "invalid service type",
			config: &Config{
				ClusterID: "test-01",
				Name:      "Test",
				Services: map[string]ServiceConfig{
					"cache": {
						Type: "invalid",
						Host: "localhost",
						Port: 6379,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: &Config{
				ClusterID: "test-01",
				Name:      "Test",
				Services: map[string]ServiceConfig{
					"cache": {
						Type: "redis",
						Host: "localhost",
						Port: 99999,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		service ServiceConfig
		wantErr bool
	}{
		{
			name: "valid redis",
			service: ServiceConfig{
				Type: "redis",
				Host: "localhost",
				Port: 6379,
			},
			wantErr: false,
		},
		{
			name: "valid postgres",
			service: ServiceConfig{
				Type: "postgres",
				Host: "localhost",
				Port: 5432,
			},
			wantErr: false,
		},
		{
			name: "missing type",
			service: ServiceConfig{
				Host: "localhost",
				Port: 6379,
			},
			wantErr: true,
		},
		{
			name: "missing host",
			service: ServiceConfig{
				Type: "redis",
				Port: 6379,
			},
			wantErr: true,
		},
		{
			name: "invalid port low",
			service: ServiceConfig{
				Type: "redis",
				Host: "localhost",
				Port: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid port high",
			service: ServiceConfig{
				Type: "redis",
				Host: "localhost",
				Port: 70000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultPoolConfig(t *testing.T) {
	pool := DefaultPoolConfig()

	if pool.MinConnections <= 0 {
		t.Error("Expected min connections to be positive")
	}

	if pool.MaxConnections <= pool.MinConnections {
		t.Error("Expected max connections to be greater than min connections")
	}

	if pool.MaxIdleTime <= 0 {
		t.Error("Expected max idle time to be positive")
	}

	if pool.MaxLifetime <= 0 {
		t.Error("Expected max lifetime to be positive")
	}
}

