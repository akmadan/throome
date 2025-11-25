package main

import (
	"context"
	"fmt"
	"log"
	"time"

	throome "github.com/akmadan/throome/sdk/go"
)

func main() {
	// Create Throome client
	client := throome.NewClient("http://localhost:9000")
	ctx := context.Background()

	// Get cluster info
	clusters, err := client.ListClusters(ctx)
	if err != nil {
		log.Fatalf("Failed to list clusters: %v", err)
	}

	if len(clusters) == 0 {
		log.Fatal("No clusters found. Please create a cluster with Kafka service first.")
	}

	cluster := clusters[0]
	clusterID := cluster.ID

	// Check if cluster has Kafka
	hasKafka := false
	for _, svc := range cluster.Services {
		if svc.Type == "kafka" {
			hasKafka = true
			break
		}
	}

	if !hasKafka {
		log.Fatal("No Kafka service found in cluster. Please add a Kafka service to your cluster.")
	}

	fmt.Printf("🚀 Kafka Demo - Cluster: %s (%s)\n\n", cluster.Name, clusterID)

	// Get queue client
	clusterClient := client.Cluster(clusterID)
	queue := clusterClient.Queue()

	// Test 1: Publish messages
	fmt.Println("1️⃣  Publishing messages...")
	topics := []string{"user-events", "orders", "notifications"}

	for i, topic := range topics {
		message := []byte(fmt.Sprintf("Test message %d to %s at %s", i+1, topic, time.Now().Format(time.RFC3339)))

		err = queue.Publish(ctx, topic, message)
		if err != nil {
			fmt.Printf("   ❌ Failed to publish to %s: %v\n", topic, err)
		} else {
			fmt.Printf("   ✅ Published to %s: %s\n", topic, string(message))
		}
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("\n✅ Demo complete! Check the monitoring page: http://localhost:9000/monitoring")
	fmt.Println("\nYou should see Kafka operations logged:")
	fmt.Println("  • PUBLISH operations")
	fmt.Println("  • Topic names")
	fmt.Println("  • Message sizes")
	fmt.Println("  • Duration and status")
}
