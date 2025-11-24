#!/usr/bin/env python3
"""Throome SDK Python Example"""

import time
from throome import ThroomClient, ServiceConfig


def main():
    # Initialize the Throome client
    client = ThroomClient(base_url="http://localhost:9000")

    # Example 1: Check gateway health
    print("=== Checking Gateway Health ===")
    health = client.health()
    print(f"Gateway Status: {health.status}\n")

    # Example 2: List all clusters
    print("=== Listing Clusters ===")
    clusters = client.list_clusters()
    print(f"Found {len(clusters)} cluster(s)")
    for cluster in clusters:
        print(f"- {cluster.name} ({cluster.id}): {len(cluster.services)} services")
    print()

    # Example 3: Create a new cluster
    print("=== Creating a New Cluster ===")
    services = {
        "redis-1": ServiceConfig(type="redis", port=6380),
        "postgres-1": ServiceConfig(
            type="postgres",
            port=5434,
            username="postgres",
            password="password",
            database="demo_db",
        ),
    }

    create_response = client.create_cluster(name="demo-cluster", services=services)
    print(f"Created cluster: demo-cluster ({create_response.cluster_id})\n")

    cluster_id = create_response.cluster_id

    # Wait for services to be ready
    print("Waiting for services to be ready...")
    time.sleep(10)

    # Example 4: Work with a specific cluster
    print("=== Working with Cluster ===")
    cluster_client = client.cluster(cluster_id)

    # Check cluster health
    try:
        cluster_health = cluster_client.health()
        print("Cluster Health:")
        for service_name, health in cluster_health.services.items():
            status = "healthy" if health.healthy else "unhealthy"
            print(f"- {service_name}: {status} (response time: {health.response_time}ms)")
    except Exception as e:
        print(f"Failed to check cluster health: {e}")
    print()

    # Example 5: Cache operations
    print("=== Cache Operations ===")
    cache = cluster_client.cache()

    # Set a value
    cache.set("user:123", "John Doe", expiration=60)
    print("Set cache: user:123 = John Doe (TTL: 60s)")

    # Get the value
    value = cache.get("user:123")
    print(f"Get cache: user:123 = {value}\n")

    # Example 6: Database operations
    print("=== Database Operations ===")
    db = cluster_client.db()

    # Create a table
    db.execute(
        """
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100),
            email VARCHAR(100)
        )
    """
    )
    print("Created table: users")

    # Insert data
    db.execute(
        "INSERT INTO users (name, email) VALUES ($1, $2)", "Alice", "alice@example.com"
    )
    print("Inserted user: Alice")

    # Query data
    rows = db.query("SELECT * FROM users")
    print(f"Found {len(rows)} user(s):")
    for row in rows:
        print(f"- ID: {row['id']}, Name: {row['name']}, Email: {row['email']}")
    print()

    # Example 7: Get activity logs
    print("=== Activity Logs ===")
    from throome import ActivityFilters

    activity_logs = cluster_client.get_activity(filters=ActivityFilters(limit=10))
    print(f"Recent activity ({len(activity_logs)} logs):")
    for log in activity_logs:
        timestamp = log.timestamp.strftime("%H:%M:%S")
        print(
            f"- [{timestamp}] {log.service_name}.{log.operation}: {log.command} ({log.status})"
        )
    print()

    # Example 8: Get service logs
    print("=== Service Logs ===")
    from throome import LogOptions

    service_client = cluster_client.service("redis-1")
    logs = service_client.get_logs(options=LogOptions(tail=20))
    print("Redis service logs (last 20 lines):")
    print(logs)
    print()

    # Example 9: Cleanup - Delete the cluster
    print("=== Cleanup ===")
    client.delete_cluster(cluster_id)
    print(f"Deleted cluster: {cluster_id}")


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"Error: {e}")
        exit(1)

