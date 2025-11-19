# Getting Started with Throome

This guide will help you get Throome up and running in minutes.

## Prerequisites

- Go 1.21 or higher
- One or more infrastructure services (Redis, PostgreSQL, Kafka, etc.)
- Basic understanding of Go and REST APIs

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/akshitmadan/throome.git
cd throome

# Build binaries
make build

# Binaries will be in ./bin/
ls bin/
# Output: throome  throome-cli
```

### Option 2: Install with Go

```bash
go install github.com/akshitmadan/throome/cmd/throome@latest
go install github.com/akshitmadan/throome/cmd/throome-cli@latest
```

### Option 3: Download Pre-built Binaries

Download from [GitHub Releases](https://github.com/akshitmadan/throome/releases)

## Quick Start

### 1. Start the Gateway

```bash
# Start with defaults (port 9000, ./clusters directory)
./bin/throome

# Or with custom configuration
./bin/throome --config configs/throome.yaml --port 9000
```

You should see:

```
INFO  Starting Throome Gateway  version=0.1.0
INFO  Gateway initialized successfully
INFO  Throome Gateway is running  port=9000
```

### 2. Create Your First Cluster

```bash
# Create a cluster
./bin/throome-cli create-cluster --name my-first-app

# Output:
# ✓ Cluster created successfully!
#   Cluster ID: my-first-01
#   Name: my-first-app
#   Config: ./clusters/my-first-01/config.yaml
```

### 3. Configure Your Services

Edit the generated configuration file:

```bash
vim clusters/my-first-01/config.yaml
```

Add your services:

```yaml
cluster_id: "my-first-01"
name: "my-first-app"

services:
  cache:
    type: redis
    host: localhost
    port: 6379
  
  database:
    type: postgres
    host: localhost
    port: 5432
    username: myuser
    password: mypass
    database: mydb

routing:
  strategy: "round_robin"
  failover_enabled: true

health:
  enabled: true
  interval: 10
```

### 4. Restart the Gateway

```bash
# Stop the current gateway (Ctrl+C)
# Start it again
./bin/throome
```

The gateway will automatically load your cluster.

### 5. Use in Your Application

Create a new Go project:

```bash
mkdir my-app
cd my-app
go mod init my-app
```

Add Throome SDK:

```bash
go get github.com/akshitmadan/throome
```

Create `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/akshitmadan/throome/pkg/sdk"
)

func main() {
    // Connect to your cluster
    client := sdk.NewClient("http://localhost:9000", "my-first-01")
    ctx := context.Background()

    // Check health
    health, err := client.Health(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Cluster healthy: %+v\n", health)

    // Use cache
    cache := client.Cache()
    if err := cache.Set(ctx, "greeting", "Hello, Throome!", 0); err != nil {
        log.Fatal(err)
    }

    value, err := cache.Get(ctx, "greeting")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Cached value: %s\n", value)

    // Use database
    db := client.DB()
    if err := db.Execute(ctx, "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)"); err != nil {
        log.Fatal(err)
    }

    if err := db.Execute(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice"); err != nil {
        log.Fatal(err)
    }

    rows, err := db.Query(ctx, "SELECT * FROM users")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Users: %+v\n", rows)
}
```

Run your application:

```bash
go run main.go
```

## Verify Installation

### Check Gateway Status

```bash
curl http://localhost:9000/
```

Response:

```json
{
  "service": "Throome Gateway",
  "version": "0.1.0",
  "status": "running"
}
```

### List Clusters

```bash
./bin/throome-cli list-clusters
```

### Get Cluster Details

```bash
./bin/throome-cli get-cluster my-first-01
```

### Check Cluster Health

```bash
curl http://localhost:9000/api/v1/clusters/my-first-01/health
```

### View Metrics

```bash
curl http://localhost:9000/metrics
```

## Next Steps

- [Configure advanced routing strategies](cluster-configuration.md#routing)
- [Set up monitoring and alerting](deployment.md#monitoring)
- [Build custom adapters](adapter-development.md)
- [Deploy to production](deployment.md)

## Troubleshooting

### Gateway won't start

- Check if port 9000 is already in use: `lsof -i :9000`
- Verify clusters directory exists and is writable
- Check logs for error messages

### Can't connect to services

- Verify service credentials in cluster config
- Check if services are running: `netstat -an | grep <port>`
- Review health check results: `curl http://localhost:9000/api/v1/clusters/<id>/health`

### SDK connection errors

- Ensure gateway URL is correct
- Verify cluster ID exists: `./bin/throome-cli list-clusters`
- Check network connectivity to gateway

## Getting Help

- [GitHub Issues](https://github.com/akshitmadan/throome/issues)
- [GitHub Discussions](https://github.com/akshitmadan/throome/discussions)
- [Documentation](README.md)

