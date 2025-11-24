# Throome Go SDK

Official Go SDK for Throome - Universal Gateway for Modern Applications.

## Installation

```bash
go get github.com/akmadan/throome/sdk/go
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    throome "github.com/akmadan/throome/sdk/go"
)

func main() {
    // Initialize client
    client := throome.NewClient("http://localhost:9000")
    ctx := context.Background()

    // Check health
    health, err := client.Health(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Gateway status: %s", health.Status)

    // List clusters
    clusters, err := client.ListClusters(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Found %d clusters", len(clusters))
}
```

## Features

- **Cluster Management**: Create, list, get, and delete clusters
- **Health Monitoring**: Check gateway and cluster health
- **Activity Logging**: View detailed activity logs
- **Service Operations**: Get service info and logs
- **Database Client**: Execute SQL queries through the gateway
- **Cache Client**: Redis operations (GET, SET, DELETE)
- **Queue Client**: Publish messages to Kafka topics

## Usage Examples

### Create a Cluster

```go
createReq := throome.CreateClusterRequest{
    Name: "my-cluster",
    Services: map[string]throome.ServiceConfig{
        "redis-1": {
            Type: "redis",
            Port: 6379,
        },
        "postgres-1": {
            Type:     "postgres",
            Port:     5432,
            Username: "postgres",
            Password: "password",
            Database: "mydb",
        },
    },
}

resp, err := client.CreateCluster(ctx, createReq)
if err != nil {
    log.Fatal(err)
}
log.Printf("Created cluster: %s", resp.ClusterID)
```

### Cache Operations

```go
cluster := client.Cluster("cluster-id")
cache := cluster.Cache()

// Set value with TTL
err := cache.Set(ctx, "user:123", "John Doe", 60*time.Second)

// Get value
value, err := cache.Get(ctx, "user:123")

// Delete value
err = cache.Delete(ctx, "user:123")
```

### Database Operations

```go
db := cluster.DB()

// Execute statement
err := db.Execute(ctx, "CREATE TABLE users (id SERIAL, name VARCHAR(100))")

// Query rows
rows, err := db.Query(ctx, "SELECT * FROM users WHERE id = $1", 123)

// Query single row
row, err := db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", 123)
```

### Get Service Logs

```go
service := cluster.Service("redis-1")

// Get last 100 lines
logs, err := service.GetLogs(ctx, throome.LogOptions{
    Tail:       100,
    Timestamps: true,
})
```

### Monitor Activity

```go
// Get cluster activity logs
logs, err := cluster.GetActivity(ctx, throome.ActivityFilters{
    Limit: 50,
})

for _, log := range logs {
    fmt.Printf("[%s] %s.%s: %s (%s)\n",
        log.Timestamp,
        log.ServiceName,
        log.Operation,
        log.Command,
        log.Status,
    )
}
```

## Complete Example

See [examples/main.go](examples/main.go) for a complete working example.

## API Reference

### ThroomClient

- `Health(ctx)`: Check gateway health
- `ListClusters(ctx)`: List all clusters
- `GetCluster(ctx, id)`: Get cluster details
- `CreateCluster(ctx, req)`: Create new cluster
- `DeleteCluster(ctx, id)`: Delete cluster
- `GetActivity(ctx, filters)`: Get global activity logs
- `Cluster(id)`: Get cluster client

### ClusterClient

- `Health(ctx)`: Check cluster health
- `Metrics(ctx)`: Get cluster metrics
- `GetActivity(ctx, filters)`: Get cluster activity logs
- `Service(name)`: Get service client
- `DB()`: Get database client
- `Cache()`: Get cache client
- `Queue()`: Get queue client

### ServiceClient

- `GetInfo(ctx)`: Get service information
- `GetLogs(ctx, options)`: Get Docker container logs
- `GetActivity(ctx, filters)`: Get service activity logs

## License

MIT

