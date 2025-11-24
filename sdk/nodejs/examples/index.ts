import { ThroomClient, ServiceConfig } from 'throome-sdk';

async function main() {
  // Initialize the Throome client
  const client = new ThroomClient({ baseURL: 'http://localhost:9000' });

  try {
    // Example 1: Check gateway health
    console.log('=== Checking Gateway Health ===');
    const health = await client.health();
    console.log(`Gateway Status: ${health.status}\n`);

    // Example 2: List all clusters
    console.log('=== Listing Clusters ===');
    const clusters = await client.listClusters();
    console.log(`Found ${clusters.length} cluster(s)`);
    clusters.forEach((cluster) => {
      console.log(
        `- ${cluster.name} (${cluster.id}): ${cluster.services?.length || 0} services`
      );
    });
    console.log();

    // Example 3: Create a new cluster
    console.log('=== Creating a New Cluster ===');
    const services: Record<string, ServiceConfig> = {
      'redis-1': {
        type: 'redis',
        port: 6380,
      },
      'postgres-1': {
        type: 'postgres',
        port: 5434,
        username: 'postgres',
        password: 'password',
        database: 'demo_db',
      },
    };

    const createResponse = await client.createCluster({
      name: 'demo-cluster',
      services,
    });
    console.log(`Created cluster: demo-cluster (${createResponse.cluster_id})\n`);

    const clusterId = createResponse.cluster_id;

    // Wait for services to be ready
    console.log('Waiting for services to be ready...');
    await new Promise((resolve) => setTimeout(resolve, 10000));

    // Example 4: Work with a specific cluster
    console.log('=== Working with Cluster ===');
    const clusterClient = client.cluster(clusterId);

    // Check cluster health
    try {
      const clusterHealth = await clusterClient.health();
      console.log('Cluster Health:');
      Object.entries(clusterHealth.services).forEach(([serviceName, health]) => {
        const status = health.healthy ? 'healthy' : 'unhealthy';
        console.log(
          `- ${serviceName}: ${status} (response time: ${health.response_time}ms)`
        );
      });
    } catch (error) {
      console.log(`Failed to check cluster health: ${error}`);
    }
    console.log();

    // Example 5: Cache operations
    console.log('=== Cache Operations ===');
    const cache = clusterClient.cache();

    // Set a value
    await cache.set('user:123', 'John Doe', { expiration: 60 });
    console.log('Set cache: user:123 = John Doe (TTL: 60s)');

    // Get the value
    const value = await cache.get('user:123');
    console.log(`Get cache: user:123 = ${value}\n`);

    // Example 6: Database operations
    console.log('=== Database Operations ===');
    const db = clusterClient.db();

    // Create a table
    await db.execute(`
      CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100)
      )
    `);
    console.log('Created table: users');

    // Insert data
    await db.execute(
      'INSERT INTO users (name, email) VALUES ($1, $2)',
      'Alice',
      'alice@example.com'
    );
    console.log('Inserted user: Alice');

    // Query data
    const rows = await db.query('SELECT * FROM users');
    console.log(`Found ${rows.length} user(s):`);
    rows.forEach((row) => {
      console.log(`- ID: ${row.id}, Name: ${row.name}, Email: ${row.email}`);
    });
    console.log();

    // Example 7: Get activity logs
    console.log('=== Activity Logs ===');
    const activityLogs = await clusterClient.getActivity({ limit: 10 });
    console.log(`Recent activity (${activityLogs.length} logs):`);
    activityLogs.forEach((log) => {
      const time = new Date(log.timestamp).toLocaleTimeString();
      console.log(
        `- [${time}] ${log.service_name}.${log.operation}: ${log.command} (${log.status})`
      );
    });
    console.log();

    // Example 8: Get service logs
    console.log('=== Service Logs ===');
    const serviceClient = clusterClient.service('redis-1');
    const logs = await serviceClient.getLogs({ tail: 20 });
    console.log('Redis service logs (last 20 lines):');
    console.log(logs);
    console.log();

    // Example 9: Cleanup - Delete the cluster
    console.log('=== Cleanup ===');
    await client.deleteCluster(clusterId);
    console.log(`Deleted cluster: ${clusterId}`);
  } catch (error) {
    console.error('Error:', error);
    process.exit(1);
  }
}

main();

