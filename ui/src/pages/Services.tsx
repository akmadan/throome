import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Database, Server, MessageSquare, ExternalLink, CheckCircle2, XCircle, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { getClusters } from '@/api/client'

interface ServiceRow {
  clusterId: string
  clusterName: string
  serviceName: string
  type: string
  host: string
  port: number
  database?: string
  username?: string
  healthy?: boolean
}

export default function Services() {
  const navigate = useNavigate()
  const [services, setServices] = useState<ServiceRow[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadServices()
  }, [])

  const loadServices = async () => {
    try {
      setLoading(true)
      const clusters = await getClusters()
      
      // Flatten all services from all clusters
      const allServices: ServiceRow[] = []
      clusters.forEach((cluster) => {
        if (cluster.services) {
          cluster.services.forEach((service) => {
            allServices.push({
              clusterId: cluster.id,
              clusterName: cluster.name,
              serviceName: service.name,
              type: service.type,
              host: service.host,
              port: service.port,
              database: service.database,
              username: service.username,
              healthy: service.healthy,
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

  const getServiceIcon = (type: string) => {
    switch (type) {
      case 'redis':
        return <Database className="w-5 h-5 text-red-500" />
      case 'postgres':
        return <Server className="w-5 h-5 text-[#FF5050]" />
      case 'kafka':
        return <MessageSquare className="w-5 h-5 text-purple-500" />
      default:
        return <Server className="w-5 h-5 text-gray-500" />
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Services</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            View all services across clusters
          </p>
        </div>
        <div className="px-4 py-2 bg-gray-100 dark:bg-gray-700 rounded-lg">
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {services.length} total service(s)
          </span>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <div className="w-8 h-8 border-4 border-[#FF5050] border-t-transparent rounded-full animate-spin" />
        </div>
      ) : services.length === 0 ? (
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-12">
          <div className="flex flex-col items-center justify-center text-gray-400">
            <Server className="w-16 h-16 mb-4 opacity-50" />
            <h3 className="text-lg font-medium mb-2">No services found</h3>
            <p className="text-sm text-center mb-6">
              Create a cluster with services to see them here
            </p>
            <button
              onClick={() => navigate('/clusters/create')}
              className="px-4 py-2 bg-[#FF5050] text-white rounded-lg hover:bg-[#ed1515] transition-colors"
            >
              Create Cluster
            </button>
          </div>
        </div>
      ) : (
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-gray-900">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Service
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Cluster
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Host
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Port
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                {services.map((service, index) => (
                  <tr
                    key={`${service.clusterId}-${service.serviceName}-${index}`}
                    className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center space-x-3">
                        {getServiceIcon(service.type)}
                        <div>
                          <div className="text-sm font-medium text-gray-900 dark:text-white">
                            {service.serviceName}
                          </div>
                          {service.database && (
                            <div className="text-xs text-gray-500 dark:text-gray-400">
                              DB: {service.database}
                            </div>
                          )}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="px-2 py-1 text-xs font-medium rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300">
                        {service.type}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm text-gray-900 dark:text-white">
                          {service.clusterName}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400 font-mono">
                          {service.clusterId.slice(0, 8)}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-mono text-gray-900 dark:text-white">
                        {service.host}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-mono text-gray-900 dark:text-white">
                        {service.port}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {service.healthy === true ? (
                        <span className="flex items-center space-x-1 text-green-600 dark:text-green-400">
                          <CheckCircle2 className="w-4 h-4" />
                          <span className="text-sm">Healthy</span>
                        </span>
                      ) : service.healthy === false ? (
                        <span className="flex items-center space-x-1 text-red-600 dark:text-red-400">
                          <XCircle className="w-4 h-4" />
                          <span className="text-sm">Unhealthy</span>
                        </span>
                      ) : (
                        <span className="flex items-center space-x-1 text-gray-500 dark:text-gray-400">
                          <Loader2 className="w-4 h-4" />
                          <span className="text-sm">Unknown</span>
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <button
                        onClick={() => navigate(`/clusters/${service.clusterId}`)}
                        className="flex items-center space-x-1 text-[#FF5050] hover:text-[#ed1515] transition-colors"
                      >
                        <ExternalLink className="w-4 h-4" />
                        <span className="text-sm">View Cluster</span>
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
