package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/akmadan/throome/pkg/cluster"
)

// TestMain sets up and tears down integration test environment
func TestMain(m *testing.M) {
	// Check if we should skip integration tests
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		fmt.Println("Skipping integration tests. Set INTEGRATION_TESTS=true to run them.")
		os.Exit(0)
	}

	// Wait for services to be ready
	fmt.Println("Waiting for test services to be ready...")
	if err := waitForServices(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to test services: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All test services ready. Running integration tests...")

	// Run tests
	code := m.Run()

	// Cleanup
	fmt.Println("Integration tests complete. Cleaning up...")

	os.Exit(code)
}

// waitForServices waits for all test services to be ready
func waitForServices() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	services := map[string]func() error{
		"Redis":      checkRedis,
		"PostgreSQL": checkPostgres,
		"Kafka":      checkKafka,
	}

	for name, check := range services {
		fmt.Printf("Checking %s...\n", name)

		// Retry with backoff
		for i := 0; i < 12; i++ {
			if err := check(); err == nil {
				fmt.Printf("✓ %s is ready\n", name)
				break
			}

			if i == 11 {
				return fmt.Errorf("%s is not ready after 60 seconds", name)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
				// Continue
			}
		}
	}

	return nil
}

// checkRedis checks if Redis is ready
func checkRedis() error {
	config := cluster.ServiceConfig{
		Type: "redis",
		Host: "localhost",
		Port: 6379,
	}

	adapter, err := getRedisAdapter(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := adapter.Connect(ctx); err != nil {
		return err
	}
	defer adapter.Disconnect(ctx)

	return adapter.Ping(ctx)
}

// checkPostgres checks if PostgreSQL is ready
func checkPostgres() error {
	config := cluster.ServiceConfig{
		Type:     "postgres",
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	}

	adapter, err := getPostgresAdapter(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := adapter.Connect(ctx); err != nil {
		return err
	}
	defer adapter.Disconnect(ctx)

	return adapter.Ping(ctx)
}

// checkKafka checks if Kafka is ready
func checkKafka() error {
	config := cluster.ServiceConfig{
		Type: "kafka",
		Host: "localhost",
		Port: 9092,
	}

	adapter, err := getKafkaAdapter(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := adapter.Connect(ctx); err != nil {
		return err
	}
	defer adapter.Disconnect(ctx)

	return adapter.Ping(ctx)
}

// Helper functions to create adapters
func getRedisAdapter(config cluster.ServiceConfig) (interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Ping(context.Context) error
}, error) {
	// Import and create Redis adapter
	return nil, fmt.Errorf("not implemented")
}

func getPostgresAdapter(config cluster.ServiceConfig) (interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Ping(context.Context) error
}, error) {
	// Import and create Postgres adapter
	return nil, fmt.Errorf("not implemented")
}

func getKafkaAdapter(config cluster.ServiceConfig) (interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Ping(context.Context) error
}, error) {
	// Import and create Kafka adapter
	return nil, fmt.Errorf("not implemented")
}
