# Throome Node.js/TypeScript SDK

Official Node.js and TypeScript SDK for Throome - Universal Gateway for Modern Applications.

## Installation

```bash
npm install throome-sdk
# or
yarn add throome-sdk
# or
pnpm add throome-sdk
```

## Quick Start

### JavaScript

```javascript
const { ThroomClient } = require('throome-sdk');

const client = new ThroomClient({ baseURL: 'http://localhost:9000' });

async function main() {
  const health = await client.health();
  console.log('Gateway status:', health.status);

  const clusters = await client.listClusters();
  console.log(`Found ${clusters.length} clusters`);
}

main().catch(console.error);
```

### TypeScript

```typescript
import { ThroomClient, ServiceConfig } from 'throome-sdk';

const client = new ThroomClient({ baseURL: 'http://localhost:9000' });

async function main() {
  const health = await client.health();
  console.log('Gateway status:', health.status);

  const clusters = await client.listClusters();
  console.log(`Found ${clusters.length} clusters`);
}

main().catch(console.error);
```

## Features

- **Cluster Management**: Create, list, get, and delete clusters
- **Health Monitoring**: Check gateway and cluster health
- **Activity Logging**: View detailed activity logs
- **Service Operations**: Get service info and logs
- **Database Client**: Execute SQL queries through the gateway
- **Cache Client**: Redis operations (GET, SET, DELETE)
- **Queue Client**: Publish messages to Kafka topics
- **Full TypeScript Support**: Complete type definitions included

## Usage Examples

### Create a Cluster

```typescript
const services: Record<string, ServiceConfig> = {
  'redis-1': {
    type: 'redis',
    port: 6379,
  },
  'postgres-1': {
    type: 'postgres',
    port: 5432,
    username: 'postgres',
    password: 'password',
    database: 'mydb',
  },
};

const response = await client.createCluster({
  name: 'my-cluster',
  services,
});

console.log('Created cluster:', response.cluster_id);
```

### Cache Operations

```typescript
const cluster = client.cluster('cluster-id');
const cache = cluster.cache();

// Set value with TTL
await cache.set('user:123', 'John Doe', { expiration: 60 });

// Get value
const value = await cache.get('user:123');

// Delete value
await cache.delete('user:123');
```

### Database Operations

```typescript
const db = cluster.db();

// Execute statement
await db.execute('CREATE TABLE users (id SERIAL, name VARCHAR(100))');

// Query rows
const rows = await db.query('SELECT * FROM users WHERE id = $1', 123);

// Query single row
const row = await db.queryRow('SELECT * FROM users WHERE id = $1', 123);
```

### Get Service Logs

```typescript
const service = cluster.service('redis-1');

// Get last 100 lines
const logs = await service.getLogs({
  tail: 100,
  timestamps: true,
});
```

### Monitor Activity

```typescript
// Get cluster activity logs
const logs = await cluster.getActivity({ limit: 50 });

logs.forEach((log) => {
  console.log(
    `[${log.timestamp}] ${log.service_name}.${log.operation}: ${log.command} (${log.status})`
  );
});
```

## Complete Example

See [examples/index.ts](examples/index.ts) for a complete working example.

## API Reference

### ThroomClient

- `health()`: Check gateway health
- `listClusters()`: List all clusters
- `getCluster(id)`: Get cluster details
- `createCluster(req)`: Create new cluster
- `deleteCluster(id)`: Delete cluster
- `getActivity(filters?)`: Get global activity logs
- `cluster(id)`: Get cluster client

### ClusterClient

- `health()`: Check cluster health
- `metrics()`: Get cluster metrics
- `getActivity(filters?)`: Get cluster activity logs
- `service(name)`: Get service client
- `db()`: Get database client
- `cache()`: Get cache client
- `queue()`: Get queue client

### ServiceClient

- `getInfo()`: Get service information
- `getLogs(options?)`: Get Docker container logs
- `getActivity(filters?)`: Get service activity logs

## TypeScript Types

All TypeScript types are exported from the main package:

```typescript
import type {
  Cluster,
  Service,
  ServiceConfig,
  CreateClusterRequest,
  HealthResponse,
  ActivityLog,
  // ... and more
} from 'throome-sdk';
```

## License

MIT

