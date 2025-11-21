import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, Server, Trash2, Eye } from 'lucide-react'
import { toast } from 'sonner'
import { getClusters, deleteCluster, type Cluster } from '@/api/client'

export default function Clusters() {
  const navigate = useNavigate()
  const [clusters, setClusters] = useState<Cluster[]>([])
  const [loading, setLoading] = useState(true)

  const loadClusters = async () => {
    try {
      setLoading(true)
      const data = await getClusters()
      setClusters(data)
    } catch (error) {
      toast.error('Failed to load clusters')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadClusters()
  }, [])

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`Are you sure you want to delete cluster "${name}"?`)) {
      return
    }

    try {
      await deleteCluster(id)
      toast.success(`Cluster "${name}" deleted successfully`)
      loadClusters()
    } catch (error) {
      toast.error('Failed to delete cluster')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Clusters</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage your infrastructure clusters
          </p>
        </div>
        <button
          onClick={() => navigate('/clusters/create')}
          className="flex items-center space-x-2 px-4 py-2 bg-[#FF5050] text-white rounded-lg hover:bg-[#ed1515] transition-colors"
        >
          <Plus className="w-5 h-5" />
          <span>Create Cluster</span>
        </button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <div className="w-8 h-8 border-4 border-[#FF5050] border-t-transparent rounded-full animate-spin" />
        </div>
      ) : clusters.length === 0 ? (
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-12">
          <div className="flex flex-col items-center justify-center text-gray-400">
            <Server className="w-16 h-16 mb-4 opacity-50" />
            <h3 className="text-lg font-medium mb-2">No clusters yet</h3>
            <p className="text-sm text-center mb-6">
              Create your first cluster to get started
            </p>
            <button
              onClick={() => navigate('/clusters/create')}
              className="px-4 py-2 bg-[#FF5050] text-white rounded-lg hover:bg-[#ed1515] transition-colors"
            >
              Create Your First Cluster
            </button>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {clusters.map((cluster) => (
            <div
              key={cluster.id}
              className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg transition-shadow"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 bg-red-100 dark:bg-red-900/20 rounded-lg flex items-center justify-center">
                    <Server className="w-6 h-6 text-[#FF5050] dark:text-[#FF5050]" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900 dark:text-white">
                      {cluster.name}
                    </h3>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {cluster.id}
                    </p>
                  </div>
                </div>
              </div>

              <div className="space-y-2 mb-4">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Services</span>
                  <span className="font-medium text-gray-900 dark:text-white">
                    {cluster.services?.length || 0}
                  </span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">Created</span>
                  <span className="font-medium text-gray-900 dark:text-white">
                    {new Date(cluster.created_at).toLocaleDateString()}
                  </span>
                </div>
              </div>

              <div className="flex space-x-2">
                <button 
                  onClick={() => navigate(`/clusters/${cluster.id}`)}
                  className="flex-1 flex items-center justify-center space-x-2 px-3 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors text-sm"
                >
                  <Eye className="w-4 h-4" />
                  <span>View</span>
                </button>
                <button
                  onClick={() => handleDelete(cluster.id, cluster.name)}
                  className="px-3 py-2 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

