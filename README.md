# 🚀 Throome

> **A lightweight, open-source Go gateway for unified backend infrastructure access.**

Throome provides a single gateway layer to access multiple infrastructure components (Redis, PostgreSQL, Kafka, etc.) via one cluster ID – eliminating direct integration complexity.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

---

## ✨ Features

- **🎯 Unified API**: Single SDK instead of multiple infrastructure SDKs
- **📦 Cluster Management**: Group services into logical clusters
- **🔄 Smart Routing**: Round-robin, weighted, least-connections, or AI-based routing
- **📊 Built-in Monitoring**: Prometheus metrics & health checks
- **⚡ High Performance**: Built in Go for low latency (<3ms overhead)
- **🔌 Extensible**: Plugin system for custom adapters
- **🛡️ Production Ready**: Connection pooling, circuit breakers, failover

---

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/akshitmadan/throome.git
cd throome

# Build binaries
make build

# Or install directly
make install
```

### Start the Gateway

```bash
# Start with default configuration
./bin/throome

# Or with custom config
./bin/throome --config throome.yaml --port 9000
```

### Create Your First Cluster

```bash
# Create a cluster
./bin/throome-cli create-cluster --name my-app

# Edit the generated config
vim clusters/<cluster-id>/config.yaml

# List clusters
./bin/throome-cli list-clusters
```

### Use in Your Application

```go
package main

import (
    "context"
    "github.com/akshitmadan/throome/pkg/sdk"
)

func main() {
    // Connect to Throome
    client := sdk.NewClient("http://localhost:9000", "your-cluster-id")
    ctx := context.Background()

    // Use cache (Redis)
    client.Cache().Set(ctx, "key", "value", 0)
    value, _ := client.Cache().Get(ctx, "key")

    // Use database (PostgreSQL)
    client.DB().Execute(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
    rows, _ := client.DB().Query(ctx, "SELECT * FROM users")

    // Use queue (Kafka)
    client.Queue().Publish(ctx, "events", []byte("message"))
}
```

---

## 📖 Documentation

- [Architecture Overview](docs/architecture.md)
- [Getting Started Guide](docs/getting-started.md)
- [Cluster Configuration](docs/cluster-configuration.md)
- [API Reference](docs/api-reference.md)
- [Building Custom Adapters](docs/adapter-development.md)
- [Deployment Guide](docs/deployment.md)

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│       Your Application                   │
│   (Uses Throome SDK - One Library)      │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         THROOME GATEWAY                  │
│  ┌────────────────────────────────┐    │
│  │  Cluster Manager                │    │
│  │  Router (with strategies)       │    │
│  │  Monitoring & Health Checks     │    │
│  │  (Optional) AI Engine           │    │
│  └────────────────────────────────┘    │
└──────┬──────────┬──────────┬───────────┘
       │          │          │
       ▼          ▼          ▼
   [Redis]   [Postgres]  [Kafka]
```

---

## 🧩 Supported Services

| Service      | Status | Adapter        |
|--------------|--------|----------------|
| Redis        | ✅     | `redis`        |
| PostgreSQL   | ✅     | `postgres`     |
| Kafka        | ✅     | `kafka`        |
| MongoDB      | 🚧     | Coming soon    |
| MySQL        | 🚧     | Coming soon    |
| RabbitMQ     | 🚧     | Coming soon    |
| Elasticsearch| 📋     | Planned        |

---

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Make
- Docker (for testing with real services)

### Building from Source

```bash
# Install dependencies
make deps

# Build all binaries
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

### Running Tests

```bash
# Unit tests
make test-unit

# Integration tests (requires Docker)
make test-integration

# Coverage report
make test-coverage
```

---

## 📊 Monitoring

Throome exposes Prometheus metrics at `/metrics`:

```bash
# View metrics
curl http://localhost:9000/metrics
```

**Key Metrics:**
- `throome_requests_total` - Total requests per cluster/service
- `throome_request_duration_seconds` - Request latency histogram
- `throome_errors_total` - Error count by type
- `throome_active_connections` - Active connections gauge

---

## 🔧 Configuration

### Gateway Configuration (`throome.yaml`)

```yaml
server:
  host: "0.0.0.0"
  port: 9000

gateway:
  clusters_dir: "./clusters"
  enable_ai: false

monitoring:
  enabled: true
  metrics_path: "/metrics"

logging:
  level: "info"
```

### Cluster Configuration (`clusters/<id>/config.yaml`)

```yaml
cluster_id: "my-cluster"
name: "My Application"

services:
  cache:
    type: redis
    host: localhost
    port: 6379
  
  database:
    type: postgres
    host: localhost
    port: 5432
    username: user
    password: pass
    database: mydb

routing:
  strategy: "round_robin"
  failover_enabled: true

health:
  enabled: true
  interval: 10
```

See [Configuration Guide](docs/cluster-configuration.md) for complete details.

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## 🗺️ Roadmap

- [x] Core gateway functionality
- [x] Redis, PostgreSQL, Kafka adapters
- [x] CLI tool
- [x] Go SDK
- [ ] Python SDK
- [ ] Node.js SDK
- [ ] Docker support
- [ ] Kubernetes operator
- [ ] AI-based routing
- [ ] gRPC support
- [ ] Dashboard UI

See [ROADMAP.md](ROADMAP.md) for detailed plans.

---

## 💬 Community

- **Issues**: [GitHub Issues](https://github.com/akshitmadan/throome/issues)
- **Discussions**: [GitHub Discussions](https://github.com/akshitmadan/throome/discussions)
- **Twitter**: [@throome_dev](https://twitter.com/throome_dev)

---

## 🙏 Acknowledgments

Built with ❤️ using:
- [Go](https://go.dev/)
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [go-redis](https://github.com/go-redis/redis) - Redis client
- [kafka-go](https://github.com/segmentio/kafka-go) - Kafka client
- [Prometheus](https://prometheus.io/) - Monitoring
- [Cobra](https://github.com/spf13/cobra) - CLI framework

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/akshitmadan">Akshit Madan</a>
</p>

<p align="center">
  <sub>If you find this project useful, please consider giving it a ⭐️!</sub>
</p>

