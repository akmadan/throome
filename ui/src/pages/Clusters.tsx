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
    <div className="h-full flex flex-col">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-semibold text-foreground">Clusters</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            Manage your infrastructure clusters
          </p>
        </div>
        <button
          onClick={() => navigate('/clusters/create')}
          className="flex items-center space-x-2 px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm font-medium"
        >
          <Plus className="w-4 h-4" />
          <span>Create Cluster</span>
        </button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-16">
          <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
        </div>
      ) : clusters.length === 0 ? (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center">
            <Server className="w-16 h-16 mx-auto mb-4 text-muted-foreground opacity-50" />
            <h3 className="text-lg font-medium text-foreground mb-2">No clusters yet</h3>
            <p className="text-sm text-muted-foreground mb-6">
              Create your first cluster to get started
            </p>
            <button
              onClick={() => navigate('/clusters/create')}
              className="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm font-medium"
            >
              Create Your First Cluster
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
                    Name
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Cluster ID
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Services
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Created
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {clusters.map((cluster, index) => (
                  <tr
                    key={cluster.id}
                    className={`border-b border-border hover:bg-accent/30 transition-colors ${
                      index === clusters.length - 1 ? 'border-b-0' : ''
                    }`}
                  >
                    <td className="px-4 py-3">
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-muted/50 rounded-md flex items-center justify-center flex-shrink-0">
                          <Server className="w-4 h-4 text-primary" />
                        </div>
                        <div>
                          <div className="text-sm font-medium text-foreground">
                            {cluster.name}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <span className="text-sm font-mono text-muted-foreground">{cluster.id}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-muted/50 text-foreground">
                        {cluster.services?.length || 0} service{(cluster.services?.length || 0) !== 1 ? 's' : ''}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <span className="text-sm text-foreground">
                        {new Date(cluster.created_at).toLocaleDateString()}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={() => navigate(`/clusters/${cluster.id}`)}
                          className="flex items-center space-x-1 px-2.5 py-1.5 text-xs font-medium text-foreground hover:bg-muted/50 rounded-md transition-colors"
                        >
                          <Eye className="w-3.5 h-3.5" />
                          <span>View</span>
                        </button>
                        <button
                          onClick={() => handleDelete(cluster.id, cluster.name)}
                          className="p-1.5 text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md transition-colors"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </button>
                      </div>
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
