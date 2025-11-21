import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
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
  services: Service[]
}

export interface Service {
  name: string
  type: string
  host: string
  port: number
  healthy: boolean
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

export const createCluster = async (data: Partial<Cluster>): Promise<Cluster> => {
  const response = await api.post<Cluster>('/clusters', data)
  return response.data
}

export const deleteCluster = async (id: string): Promise<void> => {
  await api.delete(`/clusters/${id}`)
}

export default api

