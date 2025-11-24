package main

import (
	"context"
	"fmt"
	"log"
	"time"

	throome "github.com/akmadan/throome/sdk/go"
)

func main() {
	// Initialize the Throome client
	client := throome.NewClient("http://localhost:9000")

	ctx := context.Background()

	// Example 1: Check gateway health
	fmt.Println("=== Checking Gateway Health ===")
	health, err := client.Health(ctx)
	if err != nil {
		log.Fatalf("Failed to check health: %v", err)
	}
	fmt.Printf("Gateway Status: %s\n\n", health.Status)

	// Example 2: List all clusters
	fmt.Println("=== Listing Clusters ===")
	clusters, err := client.ListClusters(ctx)
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}
	fmt.Printf("Found %d cluster(s)\n", len(clusters))
	for _, cluster := range clusters {
		fmt.Printf("- %s (%s): %d services\n", cluster.Name, cluster.ID, len(cluster.Services))
	}
	fmt.Println()

	// Example 3: Create a new cluster
	fmt.Println("=== Creating a New Cluster ===")
	createReq := throome.CreateClusterRequest{
		Name: "demo-cluster",
		Services: map[string]throome.ServiceConfig{
			"redis-1": {
				Type: "redis",
				Port: 6380,
			},
			"postgres-1": {
				Type:     "postgres",
				Port:     5434,
				Username: "postgres",
				Password: "password",
				Database: "demo_db",
			},
		},
	}

	createResp, err := client.CreateCluster(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}
	fmt.Printf("Created cluster: %s (%s)\n\n", createReq.Name, createResp.ClusterID)

	clusterID := createResp.ClusterID

	// Wait for services to be ready
	fmt.Println("Waiting for services to be ready...")
	time.Sleep(10 * time.Second)

	// Example 4: Work with a specific cluster
	fmt.Println("=== Working with Cluster ===")
	clusterClient := client.Cluster(clusterID)

	// Check cluster health
	clusterHealth, err := clusterClient.Health(ctx)
	if err != nil {
		log.Printf("Failed to check cluster health: %v", err)
	} else {
		fmt.Printf("Cluster Health:\n")
		for serviceName, health := range clusterHealth.Services {
			status := "unhealthy"
			if health.Healthy {
				status = "healthy"
			}
			fmt.Printf("- %s: %s (response time: %dms)\n", serviceName, status, health.ResponseTime)
		}
	}
	fmt.Println()

	// Example 5: Cache operations
	fmt.Println("=== Cache Operations ===")
	cache := clusterClient.Cache()

	// Set a value
	err = cache.Set(ctx, "user:123", "John Doe", 60*time.Second)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	} else {
		fmt.Println("Set cache: user:123 = John Doe (TTL: 60s)")
	}

	// Get the value
	value, err := cache.Get(ctx, "user:123")
	if err != nil {
		log.Printf("Failed to get cache: %v", err)
	} else {
		fmt.Printf("Get cache: user:123 = %s\n", value)
	}
	fmt.Println()

	// Example 6: Database operations
	fmt.Println("=== Database Operations ===")
	db := clusterClient.DB()

	// Create a table
	err = db.Execute(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			email VARCHAR(100)
		)
	`)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
	} else {
		fmt.Println("Created table: users")
	}

	// Insert data
	err = db.Execute(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "Alice", "alice@example.com")
	if err != nil {
		log.Printf("Failed to insert: %v", err)
	} else {
		fmt.Println("Inserted user: Alice")
	}

	// Query data
	rows, err := db.Query(ctx, "SELECT * FROM users")
	if err != nil {
		log.Printf("Failed to query: %v", err)
	} else {
		fmt.Printf("Found %d user(s):\n", len(rows))
		for _, row := range rows {
			fmt.Printf("- ID: %v, Name: %v, Email: %v\n", row["id"], row["name"], row["email"])
		}
	}
	fmt.Println()

	// Example 7: Get activity logs
	fmt.Println("=== Activity Logs ===")
	activityLogs, err := clusterClient.GetActivity(ctx, throome.ActivityFilters{Limit: 10})
	if err != nil {
		log.Printf("Failed to get activity logs: %v", err)
	} else {
		fmt.Printf("Recent activity (%d logs):\n", len(activityLogs))
		for _, log := range activityLogs {
			fmt.Printf("- [%s] %s.%s: %s (%s)\n",
				log.Timestamp.Format("15:04:05"),
				log.ServiceName,
				log.Operation,
				log.Command,
				log.Status,
			)
		}
	}
	fmt.Println()

	// Example 8: Get service logs
	fmt.Println("=== Service Logs ===")
	serviceClient := clusterClient.Service("redis-1")
	logs, err := serviceClient.GetLogs(ctx, throome.LogOptions{Tail: 20})
	if err != nil {
		log.Printf("Failed to get logs: %v", err)
	} else {
		fmt.Println("Redis service logs (last 20 lines):")
		fmt.Println(logs)
	}
	fmt.Println()

	// Example 9: Cleanup - Delete the cluster
	fmt.Println("=== Cleanup ===")
	err = client.DeleteCluster(ctx, clusterID)
	if err != nil {
		log.Printf("Failed to delete cluster: %v", err)
	} else {
		fmt.Printf("Deleted cluster: %s\n", clusterID)
	}
}
