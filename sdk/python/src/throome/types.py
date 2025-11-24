"""Throome SDK types"""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Any, Dict, List, Optional


@dataclass
class Service:
    """Service information"""

    name: str
    type: str
    host: str
    port: int
    healthy: bool
    username: Optional[str] = None
    database: Optional[str] = None
    container_id: Optional[str] = None


@dataclass
class Cluster:
    """Cluster information"""

    id: str
    name: str
    created_at: str
    services: List[Service] = field(default_factory=list)


@dataclass
class ServiceConfig:
    """Service configuration for cluster creation"""

    type: str
    provision: bool  # If True, Throome provisions a new Docker container; if False, connects to existing service
    port: int
    host: Optional[str] = None  # Required when provision is False
    username: Optional[str] = None  # Required for databases when provision is False
    password: Optional[str] = None  # Required for databases when provision is False
    database: Optional[str] = None  # Required for databases when provision is False


@dataclass
class CreateClusterRequest:
    """Request to create a cluster"""

    name: str
    services: Dict[str, ServiceConfig]


@dataclass
class CreateClusterResponse:
    """Response from creating a cluster"""

    cluster_id: str
    message: str


@dataclass
class HealthResponse:
    """Gateway health response"""

    status: str
    timestamp: int


@dataclass
class ServiceHealth:
    """Service health status"""

    healthy: bool
    response_time: int
    error_message: Optional[str] = None


@dataclass
class ClusterHealthResponse:
    """Cluster health response"""

    cluster_id: str
    services: Dict[str, ServiceHealth]


@dataclass
class MetricsResponse:
    """Cluster metrics"""

    requests: int
    errors: int
    avg_response_ms: float
    p95_response_ms: float
    active_services: int


@dataclass
class ServiceInfo:
    """Detailed service information"""

    name: str
    type: str
    host: str
    port: int
    healthy: bool
    container_id: Optional[str] = None
    status: Optional[str] = None


@dataclass
class ActivityLog:
    """Activity log entry"""

    id: str
    timestamp: datetime
    cluster_id: str
    service_name: str
    service_type: str
    operation: str
    command: str
    parameters: List[Any]
    duration: int
    status: str
    response: str
    error: Optional[str] = None
    client_info: Optional[Dict[str, str]] = None


@dataclass
class ActivityFilters:
    """Filters for activity logs"""

    limit: Optional[int] = None


@dataclass
class LogOptions:
    """Options for fetching service logs"""

    tail: Optional[int] = None
    timestamps: bool = False

