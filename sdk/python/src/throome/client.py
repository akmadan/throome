"""Throome SDK client"""

from typing import Any, Dict, List, Optional
import requests
from requests.exceptions import RequestException, Timeout

from .types import (
    Cluster,
    Service,
    CreateClusterRequest,
    ServiceConfig,
    CreateClusterResponse,
    HealthResponse,
    ClusterHealthResponse,
    ServiceHealth,
    MetricsResponse,
    ServiceInfo,
    ActivityLog,
    ActivityFilters,
    LogOptions,
)
from .exceptions import ThroomAPIError, ThroomConnectionError


class ThroomClient:
    """Main Throome SDK client"""

    def __init__(self, base_url: str, timeout: int = 120):
        """
        Initialize Throome client

        Args:
            base_url: Base URL of the Throome Gateway
            timeout: Request timeout in seconds (default: 120)
        """
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})

    def _request(
        self, method: str, path: str, data: Optional[Dict[str, Any]] = None, params: Optional[Dict[str, Any]] = None
    ) -> Any:
        """Make an HTTP request to the gateway"""
        url = f"{self.base_url}{path}"

        try:
            response = self.session.request(
                method=method, url=url, json=data, params=params, timeout=self.timeout
            )

            if response.status_code >= 400:
                error_data = response.json() if response.content else {}
                message = error_data.get("message") or error_data.get("error") or response.text
                raise ThroomAPIError(
                    f"Throome API Error ({response.status_code}): {message}",
                    status_code=response.status_code,
                )

            if response.content:
                return response.json()
            return None

        except Timeout as e:
            raise ThroomConnectionError(f"Request timed out: {e}")
        except RequestException as e:
            raise ThroomConnectionError(f"Connection error: {e}")

    def health(self) -> HealthResponse:
        """Get gateway health"""
        data = self._request("GET", "/api/v1/health")
        return HealthResponse(**data)

    def list_clusters(self) -> List[Cluster]:
        """List all clusters"""
        data = self._request("GET", "/api/v1/clusters")
        return [
            Cluster(
                id=c["id"],
                name=c["name"],
                created_at=c["created_at"],
                services=[Service(**s) for s in c.get("services", [])],
            )
            for c in data
        ]

    def get_cluster(self, cluster_id: str) -> Cluster:
        """Get a specific cluster"""
        data = self._request("GET", f"/api/v1/clusters/{cluster_id}")
        return Cluster(
            id=data["id"],
            name=data["name"],
            created_at=data["created_at"],
            services=[Service(**s) for s in data.get("services", [])],
        )

    def create_cluster(
        self, name: str, services: Dict[str, ServiceConfig]
    ) -> CreateClusterResponse:
        """Create a new cluster"""
        services_dict = {
            name: {
                "type": config.type,
                "port": config.port,
                "host": config.host,
                "username": config.username,
                "password": config.password,
                "database": config.database,
            }
            for name, config in services.items()
        }
        data = self._request("POST", "/api/v1/clusters", {"name": name, "services": services_dict})
        return CreateClusterResponse(**data)

    def delete_cluster(self, cluster_id: str) -> None:
        """Delete a cluster"""
        self._request("DELETE", f"/api/v1/clusters/{cluster_id}")

    def get_activity(self, filters: Optional[ActivityFilters] = None) -> List[ActivityLog]:
        """Get global activity logs"""
        params = {"limit": filters.limit} if filters and filters.limit else None
        data = self._request("GET", "/api/v1/activity", params=params)
        return [ActivityLog(**log) for log in data]

    def cluster(self, cluster_id: str) -> "ClusterClient":
        """Get a cluster client for cluster-specific operations"""
        return ClusterClient(self, cluster_id)


class ClusterClient:
    """Client for cluster-specific operations"""

    def __init__(self, client: ThroomClient, cluster_id: str):
        self._client = client
        self.cluster_id = cluster_id

    def health(self) -> ClusterHealthResponse:
        """Get cluster health"""
        data = self._client._request("GET", f"/api/v1/clusters/{self.cluster_id}/health")
        return ClusterHealthResponse(
            cluster_id=data["cluster_id"],
            services={k: ServiceHealth(**v) for k, v in data["services"].items()},
        )

    def metrics(self) -> MetricsResponse:
        """Get cluster metrics"""
        data = self._client._request("GET", f"/api/v1/clusters/{self.cluster_id}/metrics")
        return MetricsResponse(**data)

    def get_activity(self, filters: Optional[ActivityFilters] = None) -> List[ActivityLog]:
        """Get cluster activity logs"""
        params = {"limit": filters.limit} if filters and filters.limit else None
        data = self._client._request(
            "GET", f"/api/v1/clusters/{self.cluster_id}/activity", params=params
        )
        return [ActivityLog(**log) for log in data]

    def service(self, service_name: str) -> "ServiceClient":
        """Get a service client"""
        return ServiceClient(self._client, self.cluster_id, service_name)

    def db(self) -> "DBClient":
        """Get a database client"""
        return DBClient(self._client, self.cluster_id)

    def cache(self) -> "CacheClient":
        """Get a cache client"""
        return CacheClient(self._client, self.cluster_id)

    def queue(self) -> "QueueClient":
        """Get a queue client"""
        return QueueClient(self._client, self.cluster_id)


class ServiceClient:
    """Client for service-specific operations"""

    def __init__(self, client: ThroomClient, cluster_id: str, service_name: str):
        self._client = client
        self.cluster_id = cluster_id
        self.service_name = service_name

    def get_info(self) -> ServiceInfo:
        """Get service information"""
        data = self._client._request(
            "GET", f"/api/v1/clusters/{self.cluster_id}/services/{self.service_name}"
        )
        return ServiceInfo(**data)

    def get_logs(self, options: Optional[LogOptions] = None) -> str:
        """Get service Docker container logs"""
        params = {}
        if options:
            if options.tail:
                params["tail"] = options.tail
            if options.timestamps:
                params["timestamps"] = "true"

        # Get logs as text
        url = f"{self._client.base_url}/api/v1/clusters/{self.cluster_id}/services/{self.service_name}/logs"
        try:
            response = self._client.session.get(url, params=params, timeout=self._client.timeout)
            response.raise_for_status()
            return response.text
        except RequestException as e:
            raise ThroomConnectionError(f"Failed to get logs: {e}")

    def get_activity(self, filters: Optional[ActivityFilters] = None) -> List[ActivityLog]:
        """Get service activity logs"""
        params = {"limit": filters.limit} if filters and filters.limit else None
        data = self._client._request(
            "GET",
            f"/api/v1/clusters/{self.cluster_id}/services/{self.service_name}/activity",
            params=params,
        )
        return [ActivityLog(**log) for log in data]


class DBClient:
    """Client for database operations"""

    def __init__(self, client: ThroomClient, cluster_id: str):
        self._client = client
        self.cluster_id = cluster_id

    def execute(self, query: str, *args: Any) -> None:
        """Execute a SQL statement without returning results"""
        self._client._request(
            "POST",
            f"/api/v1/clusters/{self.cluster_id}/db/execute",
            {"query": query, "args": list(args)},
        )

    def query(self, query: str, *args: Any) -> List[Dict[str, Any]]:
        """Execute a SQL query and return results"""
        data = self._client._request(
            "POST",
            f"/api/v1/clusters/{self.cluster_id}/db/query",
            {"query": query, "args": list(args)},
        )
        return data["rows"]

    def query_row(self, query: str, *args: Any) -> Dict[str, Any]:
        """Execute a query that returns a single row"""
        rows = self.query(query, *args)
        if not rows:
            raise ValueError("No rows returned")
        return rows[0]


class CacheClient:
    """Client for cache operations"""

    def __init__(self, client: ThroomClient, cluster_id: str):
        self._client = client
        self.cluster_id = cluster_id

    def get(self, key: str) -> str:
        """Get a value from cache"""
        data = self._client._request(
            "POST", f"/api/v1/clusters/{self.cluster_id}/cache/get", {"key": key}
        )
        return data["value"]

    def set(self, key: str, value: str, expiration: Optional[int] = None) -> None:
        """
        Set a value in cache

        Args:
            key: Cache key
            value: Cache value
            expiration: Expiration time in seconds
        """
        payload = {"key": key, "value": value}
        if expiration is not None:
            payload["expiration"] = expiration
        self._client._request("POST", f"/api/v1/clusters/{self.cluster_id}/cache/set", payload)

    def delete(self, key: str) -> None:
        """Delete a key from cache"""
        self._client._request(
            "POST", f"/api/v1/clusters/{self.cluster_id}/cache/delete", {"key": key}
        )


class QueueClient:
    """Client for queue/message broker operations"""

    def __init__(self, client: ThroomClient, cluster_id: str):
        self._client = client
        self.cluster_id = cluster_id

    def publish(self, topic: str, message: bytes) -> None:
        """Publish a message to a topic"""
        self._client._request(
            "POST",
            f"/api/v1/clusters/{self.cluster_id}/queue/publish",
            {"topic": topic, "message": list(message)},
        )

    def subscribe(self, topic: str, handler: Any) -> None:
        """Subscribe to a topic (not yet implemented - use direct Kafka consumer)"""
        raise NotImplementedError("Subscribe not yet implemented in SDK - use direct Kafka consumer")

