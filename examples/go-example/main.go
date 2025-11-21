package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/akmadan/throome/pkg/sdk"
)

func main() {
	// Create a new Throome client
	// Replace with your actual gateway URL and cluster ID
	client := sdk.NewClient("http://localhost:9000", "example-01")

	ctx := context.Background()

	// Example 1: Check cluster health
	fmt.Println("=== Health Check ===")
	health, err := client.Health(ctx)
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Printf("Cluster ID: %s\n", health.ClusterID)
		for service, status := range health.Services {
			fmt.Printf("  %s: healthy=%v, response_time=%dms\n",
				service, status.Healthy, status.ResponseTime)
		}
	}

	// Example 2: Cache operations
	fmt.Println("\n=== Cache Operations ===")
	cache := client.Cache()

	// Set a value
	if err := cache.Set(ctx, "user:123", "John Doe", 5*time.Minute); err != nil {
		log.Printf("Cache set failed: %v", err)
	} else {
		fmt.Println("✓ Set user:123 = John Doe")
	}

	// Get the value
	value, err := cache.Get(ctx, "user:123")
	if err != nil {
		log.Printf("Cache get failed: %v", err)
	} else {
		fmt.Printf("✓ Got user:123 = %s\n", value)
	}

	// Delete the key
	if err := cache.Delete(ctx, "user:123"); err != nil {
		log.Printf("Cache delete failed: %v", err)
	} else {
		fmt.Println("✓ Deleted user:123")
	}

	// Example 3: Database operations
	fmt.Println("\n=== Database Operations ===")
	db := client.DB()

	// Execute a query
	query := "INSERT INTO users (name, email) VALUES ($1, $2)"
	if err := db.Execute(ctx, query, "Alice", "alice@example.com"); err != nil {
		log.Printf("DB execute failed: %v", err)
	} else {
		fmt.Println("✓ Inserted user Alice")
	}

	// Query data
	selectQuery := "SELECT id, name, email FROM users WHERE name = $1"
	rows, err := db.Query(ctx, selectQuery, "Alice")
	if err != nil {
		log.Printf("DB query failed: %v", err)
	} else {
		fmt.Printf("✓ Query returned %d rows\n", len(rows))
		for _, row := range rows {
			fmt.Printf("  User: %v\n", row)
		}
	}

	// Example 4: Queue operations
	fmt.Println("\n=== Queue Operations ===")
	queue := client.Queue()

	// Publish a message
	message := []byte(`{"event": "user.created", "user_id": 123}`)
	if err := queue.Publish(ctx, "user-events", message); err != nil {
		log.Printf("Queue publish failed: %v", err)
	} else {
		fmt.Println("✓ Published message to user-events topic")
	}

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("This example demonstrates basic Throome SDK usage.")
	fmt.Println("Make sure the Throome gateway is running and configured properly.")
}
