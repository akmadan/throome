package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestManagerCreate(t *testing.T) {
	// Create temp directory for tests
	tmpDir, err := os.MkdirTemp("", "throome-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)

	config := DefaultConfig("", "test-cluster")
	config.Services = map[string]ServiceConfig{
		"cache": {
			Type: "redis",
			Host: "localhost",
			Port: 6379,
		},
	}

	// Test creation
	clusterID, err := manager.Create("test-cluster", config)
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	if clusterID == "" {
		t.Error("Expected non-empty cluster ID")
	}

	// Verify cluster directory exists
	clusterDir := filepath.Join(tmpDir, clusterID)
	if _, err := os.Stat(clusterDir); os.IsNotExist(err) {
		t.Error("Cluster directory was not created")
	}

	// Verify config file exists
	configFile := filepath.Join(clusterDir, "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestManagerGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "throome-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)

	// Create a cluster
	config := DefaultConfig("", "test-cluster")
	config.Services = map[string]ServiceConfig{
		"cache": {
			Type: "redis",
			Host: "localhost",
			Port: 6379,
		},
	}

	clusterID, err := manager.Create("test-cluster", config)
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	// Get the cluster
	retrieved, err := manager.Get(clusterID)
	if err != nil {
		t.Fatalf("Failed to get cluster: %v", err)
	}

	if retrieved.ClusterID != clusterID {
		t.Errorf("Expected cluster ID %s, got %s", clusterID, retrieved.ClusterID)
	}

	if retrieved.Name != "test-cluster" {
		t.Errorf("Expected name test-cluster, got %s", retrieved.Name)
	}
}

func TestManagerDelete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "throome-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)

	// Create a cluster
	config := DefaultConfig("", "test-cluster")
	config.Services = map[string]ServiceConfig{
		"cache": {
			Type: "redis",
			Host: "localhost",
			Port: 6379,
		},
	}

	clusterID, err := manager.Create("test-cluster", config)
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	// Delete the cluster
	if err := manager.Delete(clusterID); err != nil {
		t.Fatalf("Failed to delete cluster: %v", err)
	}

	// Verify cluster directory is gone
	clusterDir := filepath.Join(tmpDir, clusterID)
	if _, err := os.Stat(clusterDir); !os.IsNotExist(err) {
		t.Error("Cluster directory still exists after deletion")
	}

	// Try to get deleted cluster
	_, err = manager.Get(clusterID)
	if err == nil {
		t.Error("Expected error when getting deleted cluster")
	}
}

func TestManagerList(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "throome-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)

	// Initially should be empty
	clusters, err := manager.List()
	if err != nil {
		t.Fatalf("Failed to list clusters: %v", err)
	}

	if len(clusters) != 0 {
		t.Errorf("Expected 0 clusters, got %d", len(clusters))
	}

	// Create multiple clusters with different names
	for i := 1; i <= 3; i++ {
		config := DefaultConfig("", fmt.Sprintf("test-cluster-%d", i))
		config.Services = map[string]ServiceConfig{
			"cache": {
				Type: "redis",
				Host: "localhost",
				Port: 6379,
			},
		}

		_, err := manager.Create(fmt.Sprintf("test-cluster-%d", i), config)
		if err != nil {
			t.Fatalf("Failed to create cluster %d: %v", i, err)
		}
	}

	// List again
	clusters, err = manager.List()
	if err != nil {
		t.Fatalf("Failed to list clusters: %v", err)
	}

	if len(clusters) != 3 {
		t.Errorf("Expected 3 clusters, got %d", len(clusters))
	}
}

func TestManagerExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "throome-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewManager(tmpDir)

	// Create a cluster
	config := DefaultConfig("", "test-cluster")
	config.Services = map[string]ServiceConfig{
		"cache": {
			Type: "redis",
			Host: "localhost",
			Port: 6379,
		},
	}

	clusterID, err := manager.Create("test-cluster", config)
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	// Check exists
	if !manager.Exists(clusterID) {
		t.Error("Expected cluster to exist")
	}

	// Check non-existent cluster
	if manager.Exists("non-existent") {
		t.Error("Expected cluster to not exist")
	}
}

