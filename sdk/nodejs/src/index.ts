import axios, { AxiosInstance, AxiosError } from 'axios';

// Types
export interface ThroomClientOptions {
  baseURL: string;
  timeout?: number;
}

export interface Cluster {
  id: string;
  name: string;
  services?: Service[];
  created_at: string;
}

export interface Service {
  name: string;
  type: string;
  host: string;
  port: number;
  username?: string;
  database?: string;
  healthy: boolean;
  container_id?: string;
}

export interface CreateClusterRequest {
  name: string;
  services: Record<string, ServiceConfig>;
}

export interface ServiceConfig {
  type: string;
  host?: string;
  port: number;
  username?: string;
  password?: string;
  database?: string;
}

export interface CreateClusterResponse {
  cluster_id: string;
  message: string;
}

export interface HealthResponse {
  status: string;
  timestamp: number;
}

export interface ClusterHealthResponse {
  cluster_id: string;
  services: Record<string, ServiceHealth>;
}

export interface ServiceHealth {
  healthy: boolean;
  response_time: number;
  error_message?: string;
}

export interface MetricsResponse {
  requests: number;
  errors: number;
  avg_response_ms: number;
  p95_response_ms: number;
  active_services: number;
}

export interface ServiceInfo {
  name: string;
  type: string;
  host: string;
  port: number;
  healthy: boolean;
  container_id?: string;
  status?: string;
}

export interface ActivityLog {
  id: string;
  timestamp: string;
  cluster_id: string;
  service_name: string;
  service_type: string;
  operation: string;
  command: string;
  parameters: any[];
  duration: number;
  status: string;
  response: string;
  error?: string;
  client_info?: Record<string, string>;
}

export interface ActivityFilters {
  limit?: number;
}

export interface LogOptions {
  tail?: number;
  timestamps?: boolean;
}

export interface DBQueryRequest {
  query: string;
  args?: any[];
}

export interface DBQueryResponse {
  rows: Record<string, any>[];
}

export interface CacheSetOptions {
  expiration?: number; // in seconds
}

// Main Client
export class ThroomClient {
  private client: AxiosInstance;
  private baseURL: string;

  constructor(options: ThroomClientOptions) {
    this.baseURL = options.baseURL;
    this.client = axios.create({
      baseURL: options.baseURL,
      timeout: options.timeout || 120000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response) {
          const data = error.response.data as any;
          throw new Error(
            `Throome API Error (${error.response.status}): ${data.message || data.error || error.message}`
          );
        }
        throw error;
      }
    );
  }

  /**
   * Get gateway health
   */
  async health(): Promise<HealthResponse> {
    const response = await this.client.get<HealthResponse>('/api/v1/health');
    return response.data;
  }

  /**
   * List all clusters
   */
  async listClusters(): Promise<Cluster[]> {
    const response = await this.client.get<Cluster[]>('/api/v1/clusters');
    return response.data;
  }

  /**
   * Get a specific cluster
   */
  async getCluster(clusterId: string): Promise<Cluster> {
    const response = await this.client.get<Cluster>(`/api/v1/clusters/${clusterId}`);
    return response.data;
  }

  /**
   * Create a new cluster
   */
  async createCluster(request: CreateClusterRequest): Promise<CreateClusterResponse> {
    const response = await this.client.post<CreateClusterResponse>('/api/v1/clusters', request);
    return response.data;
  }

  /**
   * Delete a cluster
   */
  async deleteCluster(clusterId: string): Promise<void> {
    await this.client.delete(`/api/v1/clusters/${clusterId}`);
  }

  /**
   * Get global activity logs
   */
  async getActivity(filters?: ActivityFilters): Promise<ActivityLog[]> {
    const params = filters?.limit ? { limit: filters.limit } : {};
    const response = await this.client.get<ActivityLog[]>('/api/v1/activity', { params });
    return response.data;
  }

  /**
   * Get a cluster client for cluster-specific operations
   */
  cluster(clusterId: string): ClusterClient {
    return new ClusterClient(this.client, clusterId);
  }
}

// Cluster Client
export class ClusterClient {
  constructor(
    private client: AxiosInstance,
    private clusterId: string
  ) {}

  /**
   * Get cluster health
   */
  async health(): Promise<ClusterHealthResponse> {
    const response = await this.client.get<ClusterHealthResponse>(
      `/api/v1/clusters/${this.clusterId}/health`
    );
    return response.data;
  }

  /**
   * Get cluster metrics
   */
  async metrics(): Promise<MetricsResponse> {
    const response = await this.client.get<MetricsResponse>(
      `/api/v1/clusters/${this.clusterId}/metrics`
    );
    return response.data;
  }

  /**
   * Get cluster activity logs
   */
  async getActivity(filters?: ActivityFilters): Promise<ActivityLog[]> {
    const params = filters?.limit ? { limit: filters.limit } : {};
    const response = await this.client.get<ActivityLog[]>(
      `/api/v1/clusters/${this.clusterId}/activity`,
      { params }
    );
    return response.data;
  }

  /**
   * Get a service client
   */
  service(serviceName: string): ServiceClient {
    return new ServiceClient(this.client, this.clusterId, serviceName);
  }

  /**
   * Get a database client
   */
  db(): DBClient {
    return new DBClient(this.client, this.clusterId);
  }

  /**
   * Get a cache client
   */
  cache(): CacheClient {
    return new CacheClient(this.client, this.clusterId);
  }

  /**
   * Get a queue client
   */
  queue(): QueueClient {
    return new QueueClient(this.client, this.clusterId);
  }
}

// Service Client
export class ServiceClient {
  constructor(
    private client: AxiosInstance,
    private clusterId: string,
    private serviceName: string
  ) {}

  /**
   * Get service information
   */
  async getInfo(): Promise<ServiceInfo> {
    const response = await this.client.get<ServiceInfo>(
      `/api/v1/clusters/${this.clusterId}/services/${this.serviceName}`
    );
    return response.data;
  }

  /**
   * Get service Docker container logs
   */
  async getLogs(options?: LogOptions): Promise<string> {
    const params: any = {};
    if (options?.tail) params.tail = options.tail;
    if (options?.timestamps) params.timestamps = true;

    const response = await this.client.get<string>(
      `/api/v1/clusters/${this.clusterId}/services/${this.serviceName}/logs`,
      { params }
    );
    return response.data;
  }

  /**
   * Get service activity logs
   */
  async getActivity(filters?: ActivityFilters): Promise<ActivityLog[]> {
    const params = filters?.limit ? { limit: filters.limit } : {};
    const response = await this.client.get<ActivityLog[]>(
      `/api/v1/clusters/${this.clusterId}/services/${this.serviceName}/activity`,
      { params }
    );
    return response.data;
  }
}

// Database Client
export class DBClient {
  constructor(
    private client: AxiosInstance,
    private clusterId: string
  ) {}

  /**
   * Execute a SQL statement without returning results
   */
  async execute(query: string, ...args: any[]): Promise<void> {
    await this.client.post(`/api/v1/clusters/${this.clusterId}/db/execute`, {
      query,
      args,
    });
  }

  /**
   * Execute a SQL query and return results
   */
  async query(query: string, ...args: any[]): Promise<Record<string, any>[]> {
    const response = await this.client.post<DBQueryResponse>(
      `/api/v1/clusters/${this.clusterId}/db/query`,
      { query, args }
    );
    return response.data.rows;
  }

  /**
   * Execute a query that returns a single row
   */
  async queryRow(query: string, ...args: any[]): Promise<Record<string, any>> {
    const rows = await this.query(query, ...args);
    if (rows.length === 0) {
      throw new Error('No rows returned');
    }
    return rows[0];
  }
}

// Cache Client
export class CacheClient {
  constructor(
    private client: AxiosInstance,
    private clusterId: string
  ) {}

  /**
   * Get a value from cache
   */
  async get(key: string): Promise<string> {
    const response = await this.client.post<{ value: string }>(
      `/api/v1/clusters/${this.clusterId}/cache/get`,
      { key }
    );
    return response.data.value;
  }

  /**
   * Set a value in cache
   */
  async set(key: string, value: string, options?: CacheSetOptions): Promise<void> {
    await this.client.post(`/api/v1/clusters/${this.clusterId}/cache/set`, {
      key,
      value,
      expiration: options?.expiration,
    });
  }

  /**
   * Delete a key from cache
   */
  async delete(key: string): Promise<void> {
    await this.client.post(`/api/v1/clusters/${this.clusterId}/cache/delete`, {
      key,
    });
  }
}

// Queue Client
export class QueueClient {
  constructor(
    private client: AxiosInstance,
    private clusterId: string
  ) {}

  /**
   * Publish a message to a topic
   */
  async publish(topic: string, message: Buffer | string): Promise<void> {
    const messageData = Buffer.isBuffer(message) ? message : Buffer.from(message);
    await this.client.post(`/api/v1/clusters/${this.clusterId}/queue/publish`, {
      topic,
      message: Array.from(messageData),
    });
  }

  /**
   * Subscribe to a topic (placeholder - use direct Kafka consumer)
   */
  async subscribe(topic: string, handler: (message: Buffer) => void | Promise<void>): Promise<void> {
    throw new Error('Subscribe not yet implemented in SDK - use direct Kafka consumer');
  }
}

// Default export
export default ThroomClient;

