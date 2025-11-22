import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Server, Database, Box, ExternalLink } from 'lucide-react'
import { toast } from 'sonner'
import { getClusters, type Cluster } from '@/api/client'

interface ServiceRow {
  clusterId: string
  clusterName: string
  serviceName: string
  type: string
  host: string
  port: number
  database?: string
  healthy?: boolean
}

export default function Services() {
  const navigate = useNavigate()
  const [services, setServices] = useState<ServiceRow[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchClusters = async () => {
      try {
        const clusters = await getClusters()
        const allServices: ServiceRow[] = []
        
        clusters.forEach((cluster: Cluster) => {
          if (cluster.config?.services) {
            Object.entries(cluster.config.services).forEach(([serviceName, config]: [string, any]) => {
              allServices.push({
                clusterId: cluster.id,
                clusterName: cluster.name,
                serviceName,
                type: config.type,
                host: config.host || 'N/A',
                port: config.port || 0,
                database: config.database,
                healthy: config.healthy !== false,
              })
            })
          }
        })
        
        setServices(allServices)
      } catch (error) {
        toast.error('Failed to load services')
      } finally {
        setLoading(false)
      }
    }

    fetchClusters()
  }, [])

  const getServiceIcon = (type: string) => {
    switch (type.toLowerCase()) {
      case 'postgres':
      case 'postgresql':
        return <Database className="w-4 h-4 text-blue-400" />
      case 'redis':
        return <Database className="w-4 h-4 text-red-400" />
      case 'kafka':
        return <Box className="w-4 h-4 text-purple-400" />
      default:
        return <Server className="w-4 h-4 text-muted-foreground" />
    }
  }

  return (
    <div className="h-full flex flex-col">
      {/* Page Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-semibold text-foreground">Services</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            View all services across clusters
          </p>
        </div>
        <div className="px-3 py-1.5 bg-muted/50 rounded-md">
          <span className="text-sm text-muted-foreground">
            {services.length} total service{services.length !== 1 ? 's' : ''}
          </span>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-16">
          <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
        </div>
      ) : services.length === 0 ? (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center">
            <Server className="w-16 h-16 mx-auto mb-4 text-muted-foreground opacity-50" />
            <h3 className="text-lg font-medium text-foreground mb-2">No services found</h3>
            <p className="text-sm text-muted-foreground mb-6">
              Create a cluster with services to see them here
            </p>
            <button
              onClick={() => navigate('/clusters/create')}
              className="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm font-medium"
            >
              Create Cluster
            </button>
          </div>
        </div>
      ) : (
        <div className="bg-card rounded-lg border border-border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border">
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Service
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Type
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Cluster
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Host
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Port
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Status
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {services.map((service, index) => (
                  <tr
                    key={`${service.clusterId}-${service.serviceName}-${index}`}
                    className="border-b border-border hover:bg-accent/30 transition-colors"
                  >
                    <td className="px-4 py-3">
                      <div className="flex items-center space-x-2.5">
                        {getServiceIcon(service.type)}
                        <div>
                          <div className="text-sm font-medium text-foreground">
                            {service.serviceName}
                          </div>
                          {service.database && (
                            <div className="text-xs text-muted-foreground">
                              {service.database}
                            </div>
                          )}
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <span className="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-muted/50 text-foreground">
                        {service.type}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <div className="text-sm text-foreground">{service.clusterName}</div>
                      <div className="text-xs text-muted-foreground font-mono">{service.clusterId}</div>
                    </td>
                    <td className="px-4 py-3">
                      <span className="text-sm font-mono text-muted-foreground">{service.host}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className="text-sm font-mono text-foreground">{service.port}</span>
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex items-center space-x-2">
                        <div className={`w-2 h-2 rounded-full ${service.healthy ? 'bg-green-500' : 'bg-red-500'}`} />
                        <span className={`text-xs font-medium ${service.healthy ? 'text-green-500' : 'text-red-500'}`}>
                          {service.healthy ? 'Healthy' : 'Unhealthy'}
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <button
                        onClick={() => navigate(`/services/${service.clusterId}/${service.serviceName}`)}
                        className="flex items-center space-x-1 px-2.5 py-1.5 text-xs font-medium text-foreground hover:bg-muted/50 rounded-md transition-colors"
                      >
                        <ExternalLink className="w-3.5 h-3.5" />
                        <span>View</span>
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}
