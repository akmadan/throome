package gateway

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/akmadan/throome/internal/logger"
	"github.com/akmadan/throome/pkg/adapters"
	"github.com/akmadan/throome/pkg/adapters/kafka"
	"github.com/akmadan/throome/pkg/adapters/postgres"
	"github.com/akmadan/throome/pkg/adapters/redis"
	"github.com/akmadan/throome/pkg/cluster"
	"github.com/akmadan/throome/pkg/monitor"
	"github.com/akmadan/throome/pkg/router"
	"go.uber.org/zap"
)

// Gateway is the main Throome gateway service
type Gateway struct {
	clusterManager *cluster.Manager
	routers        map[string]*router.Router
	adapters       map[string]map[string]adapters.Adapter // clusterID -> serviceName -> adapter
	adapterFactory *adapters.Factory
	collector      *monitor.Collector
	healthChecker  *monitor.HealthChecker
	mu             sync.RWMutex
}

// NewGateway creates a new gateway instance
func NewGateway(clustersDir string) (*Gateway, error) {
	// Create cluster manager
	clusterManager := cluster.NewManager(clustersDir)

	// Create adapter factory
	factory := adapters.NewFactory()

	// Register adapter constructors
	factory.Register("redis", redis.NewRedisAdapter)
	factory.Register("postgres", postgres.NewPostgresAdapter)
	factory.Register("kafka", kafka.NewKafkaAdapter)

	// Create collector
	collector := monitor.NewCollector()

	// Create health checker (10s interval, 5s timeout, 3 failures threshold)
	healthChecker := monitor.NewHealthChecker(10*time.Second, 5*time.Second, 3)

	return &Gateway{
		clusterManager: clusterManager,
		routers:        make(map[string]*router.Router),
		adapters:       make(map[string]map[string]adapters.Adapter),
		adapterFactory: factory,
		collector:      collector,
		healthChecker:  healthChecker,
	}, nil
}

// Initialize initializes the gateway by loading all clusters
func (g *Gateway) Initialize(ctx context.Context) error {
	logger.Info("Initializing gateway...")

	// Load all clusters
	if err := g.clusterManager.LoadAll(); err != nil {
		return fmt.Errorf("failed to load clusters: %w", err)
	}

	configs := g.clusterManager.GetAllConfigs()
	logger.Info("Loaded clusters", zap.Int("count", len(configs)))

	// Initialize adapters for each cluster
	for clusterID, config := range configs {
		if err := g.initializeCluster(ctx, clusterID, config); err != nil {
			logger.Error("Failed to initialize cluster",
				zap.String("cluster_id", clusterID),
				zap.Error(err),
			)
			continue
		}
	}

	logger.Info("Gateway initialized successfully")
	return nil
}

// initializeCluster initializes a single cluster
func (g *Gateway) initializeCluster(ctx context.Context, clusterID string, config *cluster.Config) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	logger.Info("Initializing cluster",
		zap.String("cluster_id", clusterID),
		zap.String("name", config.Name),
	)

	// Create adapters for this cluster
	clusterAdapters := make(map[string]adapters.Adapter)

	for serviceName := range config.Services {
		serviceConfig := config.Services[serviceName]
		adapter, err := g.adapterFactory.Create(&serviceConfig)
		if err != nil {
			logger.Error("Failed to create adapter",
				zap.String("cluster_id", clusterID),
				zap.String("service", serviceName),
				zap.Error(err),
			)
			continue
		}

		// Connect to the service
		if err := adapter.Connect(ctx); err != nil {
			logger.Error("Failed to connect adapter",
				zap.String("cluster_id", clusterID),
				zap.String("service", serviceName),
				zap.Error(err),
			)
			continue
		}

		clusterAdapters[serviceName] = adapter
		logger.Info("Connected to service",
			zap.String("cluster_id", clusterID),
			zap.String("service", serviceName),
			zap.String("type", serviceConfig.Type),
		)
	}

	// Store adapters
	g.adapters[clusterID] = clusterAdapters

	// Create router for this cluster
	g.routers[clusterID] = router.NewRouter(config, clusterAdapters)

	return nil
}

// GetRouter returns the router for a cluster
func (g *Gateway) GetRouter(clusterID string) (*router.Router, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	r, exists := g.routers[clusterID]
	if !exists {
		return nil, fmt.Errorf("cluster not found: %s", clusterID)
	}

	return r, nil
}

// GetAdapter returns an adapter for a specific service in a cluster
func (g *Gateway) GetAdapter(clusterID, serviceName string) (adapters.Adapter, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	clusterAdapters, exists := g.adapters[clusterID]
	if !exists {
		return nil, fmt.Errorf("cluster not found: %s", clusterID)
	}

	adapter, exists := clusterAdapters[serviceName]
	if !exists {
		return nil, fmt.Errorf("service not found: %s", serviceName)
	}

	return adapter, nil
}

// GetCollector returns the metrics collector
func (g *Gateway) GetCollector() *monitor.Collector {
	return g.collector
}

// GetHealthChecker returns the health checker
func (g *Gateway) GetHealthChecker() *monitor.HealthChecker {
	return g.healthChecker
}

// GetClusterManager returns the cluster manager
func (g *Gateway) GetClusterManager() *cluster.Manager {
	return g.clusterManager
}

// CreateCluster creates a new cluster
func (g *Gateway) CreateCluster(ctx context.Context, name string, config *cluster.Config) (string, error) {
	// Create cluster
	clusterID, err := g.clusterManager.Create(name, config)
	if err != nil {
		return "", err
	}

	// Initialize the cluster
	loadedConfig, err := g.clusterManager.Get(clusterID)
	if err != nil {
		return "", err
	}

	if err := g.initializeCluster(ctx, clusterID, loadedConfig); err != nil {
		return "", err
	}

	logger.Info("Cluster created",
		zap.String("cluster_id", clusterID),
		zap.String("name", name),
	)

	return clusterID, nil
}

// DeleteCluster deletes a cluster
func (g *Gateway) DeleteCluster(ctx context.Context, clusterID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Disconnect all adapters
	if clusterAdapters, exists := g.adapters[clusterID]; exists {
		for _, adapter := range clusterAdapters {
			if err := adapter.Disconnect(ctx); err != nil {
				logger.Error("Failed to disconnect adapter",
					zap.String("cluster_id", clusterID),
					zap.Error(err),
				)
			}
		}
		delete(g.adapters, clusterID)
	}

	// Remove router
	delete(g.routers, clusterID)

	// Delete cluster
	if err := g.clusterManager.Delete(clusterID); err != nil {
		return err
	}

	logger.Info("Cluster deleted", zap.String("cluster_id", clusterID))
	return nil
}

// Shutdown gracefully shuts down the gateway
func (g *Gateway) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down gateway...")

	g.mu.Lock()
	defer g.mu.Unlock()

	// Stop health checker
	g.healthChecker.Stop()

	// Disconnect all adapters
	for clusterID, clusterAdapters := range g.adapters {
		for serviceName, adapter := range clusterAdapters {
			logger.Info("Disconnecting adapter",
				zap.String("cluster_id", clusterID),
				zap.String("service", serviceName),
			)
			if err := adapter.Disconnect(ctx); err != nil {
				logger.Error("Failed to disconnect adapter",
					zap.String("cluster_id", clusterID),
					zap.String("service", serviceName),
					zap.Error(err),
				)
			}
		}
	}

	logger.Info("Gateway shutdown complete")
	return nil
}

// ListClusters returns a list of all clusters
func (g *Gateway) ListClusters() ([]string, error) {
	return g.clusterManager.List()
}

// GetClusterConfig returns the configuration for a cluster
func (g *Gateway) GetClusterConfig(clusterID string) (*cluster.Config, error) {
	return g.clusterManager.Get(clusterID)
}
