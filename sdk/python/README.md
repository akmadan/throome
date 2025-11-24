# Throome Python SDK

Official Python SDK for Throome - Universal Gateway for Modern Applications.

## Installation

```bash
pip install throome-sdk
```

## Quick Start

```python
from throome import ThroomClient

client = ThroomClient(base_url="http://localhost:9000")

# Check health
health = client.health()
print(f"Gateway status: {health.status}")

# List clusters
clusters = client.list_clusters()
print(f"Found {len(clusters)} clusters")
```

## Features

- **Cluster Management**: Create, list, get, and delete clusters
- **Health Monitoring**: Check gateway and cluster health
- **Activity Logging**: View detailed activity logs
- **Service Operations**: Get service info and logs
- **Database Client**: Execute SQL queries through the gateway
- **Cache Client**: Redis operations (GET, SET, DELETE)
- **Queue Client**: Publish messages to Kafka topics
- **Type Hints**: Full type annotation support
- **Dataclasses**: Clean, Pythonic data structures

## Usage Examples

### Create a Cluster

```python
from throome import ServiceConfig

services = {
    "redis-1": ServiceConfig(type="redis", port=6379),
    "postgres-1": ServiceConfig(
        type="postgres",
        port=5432,
        username="postgres",
        password="password",
        database="mydb",
    ),
}

response = client.create_cluster(name="my-cluster", services=services)
print(f"Created cluster: {response.cluster_id}")
```

### Cache Operations

```python
cluster = client.cluster("cluster-id")
cache = cluster.cache()

# Set value with TTL
cache.set("user:123", "John Doe", expiration=60)

# Get value
value = cache.get("user:123")

# Delete value
cache.delete("user:123")
```

### Database Operations

```python
db = cluster.db()

# Execute statement
db.execute("CREATE TABLE users (id SERIAL, name VARCHAR(100))")

# Query rows
rows = db.query("SELECT * FROM users WHERE id = $1", 123)

# Query single row
row = db.query_row("SELECT * FROM users WHERE id = $1", 123)
```

### Get Service Logs

```python
from throome import LogOptions

service = cluster.service("redis-1")

# Get last 100 lines
logs = service.get_logs(options=LogOptions(tail=100, timestamps=True))
```

### Monitor Activity

```python
from throome import ActivityFilters

# Get cluster activity logs
logs = cluster.get_activity(filters=ActivityFilters(limit=50))

for log in logs:
    timestamp = log.timestamp.strftime("%H:%M:%S")
    print(f"[{timestamp}] {log.service_name}.{log.operation}: {log.command} ({log.status})")
```

## Complete Example

See [examples/main.py](examples/main.py) for a complete working example.

## API Reference

### ThroomClient

- `health()`: Check gateway health
- `list_clusters()`: List all clusters
- `get_cluster(cluster_id)`: Get cluster details
- `create_cluster(name, services)`: Create new cluster
- `delete_cluster(cluster_id)`: Delete cluster
- `get_activity(filters=None)`: Get global activity logs
- `cluster(cluster_id)`: Get cluster client

### ClusterClient

- `health()`: Check cluster health
- `metrics()`: Get cluster metrics
- `get_activity(filters=None)`: Get cluster activity logs
- `service(service_name)`: Get service client
- `db()`: Get database client
- `cache()`: Get cache client
- `queue()`: Get queue client

### ServiceClient

- `get_info()`: Get service information
- `get_logs(options=None)`: Get Docker container logs
- `get_activity(filters=None)`: Get service activity logs

## Error Handling

```python
from throome import ThroomAPIError, ThroomConnectionError

try:
    cluster = client.get_cluster("invalid-id")
except ThroomAPIError as e:
    print(f"API Error ({e.status_code}): {e}")
except ThroomConnectionError as e:
    print(f"Connection Error: {e}")
```

## Type Hints

The SDK includes full type annotations for better IDE support:

```python
from throome import (
    Cluster,
    Service,
    ServiceConfig,
    CreateClusterRequest,
    HealthResponse,
    ActivityLog,
    # ... and more
)
```

## Requirements

- Python 3.8 or higher
- requests >= 2.31.0

## License

MIT

