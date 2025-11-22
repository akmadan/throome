import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 120000, // 2 minutes - enough time for Docker image pulls
})

export interface HealthResponse {
  status: string
  timestamp: string
  version?: string
}

export interface Cluster {
  id: string
  name: string
  created_at: string
  services?: Service[]
  config?: any
}

export interface Service {
  name: string
  type: string
  host: string
  port: number
  healthy?: boolean
  status?: string
  container_id?: string
  username?: string
  password?: string
  database?: string
}

export interface Activity {
  id: string
  timestamp: string
  cluster_id: string
  service_name: string
  service_type: string
  operation: string
  command: string
  parameters?: any[]
  duration: number // milliseconds
  status: 'success' | 'error'
  response?: string
  error?: string
  rows_affected?: number
  client_info?: Record<string, string>
}

export interface ActivityFilters {
  cluster_id?: string
  service_name?: string
  service_type?: string
  operation?: string
  status?: 'success' | 'error'
  since?: string
  limit?: number
}

export interface ActivityResponse {
  activities: Activity[]
  count: number
  cluster_id?: string
  service_name?: string
  filters?: ActivityFilters
}

// Health check
export const checkHealth = async (): Promise<HealthResponse> => {
  const response = await api.get<HealthResponse>('/health')
  return response.data
}

// Clusters
export const getClusters = async (): Promise<Cluster[]> => {
  const response = await api.get<Cluster[]>('/clusters')
  return response.data
}

export const getCluster = async (id: string): Promise<Cluster> => {
  const response = await api.get<Cluster>(`/clusters/${id}`)
  return response.data
}

export const createCluster = async (data: {
  name: string
  config: any
}): Promise<Cluster> => {
  const response = await api.post<Cluster>('/clusters', data)
  return response.data
}

export const deleteCluster = async (id: string): Promise<void> => {
  await api.delete(`/clusters/${id}`)
}

// Activity Logs
export const getActivity = async (filters?: ActivityFilters): Promise<ActivityResponse> => {
  const params = new URLSearchParams()
  if (filters?.cluster_id) params.append('cluster_id', filters.cluster_id)
  if (filters?.service_name) params.append('service_name', filters.service_name)
  if (filters?.service_type) params.append('service_type', filters.service_type)
  if (filters?.operation) params.append('operation', filters.operation)
  if (filters?.status) params.append('status', filters.status)
  if (filters?.since) params.append('since', filters.since)
  if (filters?.limit) params.append('limit', filters.limit.toString())

  const response = await api.get<ActivityResponse>(`/activity?${params.toString()}`)
  return response.data
}

export const getClusterActivity = async (
  clusterId: string,
  limit?: number
): Promise<ActivityResponse> => {
  const params = new URLSearchParams()
  if (limit) params.append('limit', limit.toString())

  const response = await api.get<ActivityResponse>(
    `/clusters/${clusterId}/activity?${params.toString()}`
  )
  return response.data
}

export const getServiceActivity = async (
  clusterId: string,
  serviceName: string,
  limit?: number
): Promise<ActivityResponse> => {
  const params = new URLSearchParams()
  if (limit) params.append('limit', limit.toString())

  const response = await api.get<ActivityResponse>(
    `/clusters/${clusterId}/services/${serviceName}/activity?${params.toString()}`
  )
  return response.data
}

export default api

