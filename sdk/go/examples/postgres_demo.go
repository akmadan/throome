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
		log.Fatal("No clusters found. Please create a cluster with PostgreSQL service first.")
	}

	cluster := clusters[0]
	clusterID := cluster.ID

	fmt.Printf("🚀 PostgreSQL Demo - Cluster: %s (%s)\n\n", cluster.Name, clusterID)

	// Get cluster client
	clusterClient := client.Cluster(clusterID)
	db := clusterClient.DB()

	// 1. Create a table
	fmt.Println("1️⃣  Creating 'users' table...")
	err = db.Execute(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			email VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		fmt.Printf("   ⚠️  Error: %v\n", err)
	} else {
		fmt.Println("   ✅ Table created")
	}
	time.Sleep(500 * time.Millisecond)

	// 2. Insert data
	fmt.Println("\n2️⃣  Inserting users...")
	users := []struct {
		name  string
		email string
	}{
		{"Alice Johnson", "alice@example.com"},
		{"Bob Smith", "bob@example.com"},
		{"Charlie Brown", "charlie@example.com"},
	}

	for _, user := range users {
		err = db.Execute(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", user.name, user.email)
		if err != nil {
			fmt.Printf("   ❌ Failed to insert %s: %v\n", user.name, err)
		} else {
			fmt.Printf("   ✅ Inserted: %s\n", user.name)
		}
		time.Sleep(300 * time.Millisecond)
	}

	// 3. Query all users
	fmt.Println("\n3️⃣  Querying all users...")
	rows, err := db.Query(ctx, "SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		fmt.Printf("   ❌ Query failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Found %d users:\n", len(rows))
		for _, row := range rows {
			fmt.Printf("      - ID: %v, Name: %v, Email: %v\n", row["id"], row["name"], row["email"])
		}
	}

	fmt.Println("\n✅ Demo complete! Check the monitoring page: http://localhost:9000/monitoring")
}

