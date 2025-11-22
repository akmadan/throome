import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, FileCode, Workflow, Trash2, Activity, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import YamlEditor from '@/components/YamlEditor'
import CanvasEditor from '@/components/CanvasEditor'
import { getCluster, deleteCluster, type Cluster } from '@/api/client'

type ViewMode = 'yaml' | 'canvas'

export default function ViewCluster() {
  const navigate = useNavigate()
  const { clusterId } = useParams<{ clusterId: string }>()
  const [viewMode, setViewMode] = useState<ViewMode>('canvas')
  const [cluster, setCluster] = useState<Cluster | null>(null)
  const [loading, setLoading] = useState(true)
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    loadCluster()
  }, [clusterId])

  const loadCluster = async () => {
    if (!clusterId) return

    try {
      setLoading(true)
      const data = await getCluster(clusterId)
      setCluster(data)
    } catch (error) {
      toast.error('Failed to load cluster')
      navigate('/clusters')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!cluster) return

    if (!confirm(`Are you sure you want to delete cluster "${cluster.name}"?`)) {
      return
    }

    setDeleting(true)
    try {
      await deleteCluster(cluster.id)
      toast.success(`Cluster "${cluster.name}" deleted successfully`)
      navigate('/clusters')
    } catch (error) {
      toast.error('Failed to delete cluster')
    } finally {
      setDeleting(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="flex flex-col items-center space-y-4">
          <Loader2 className="w-10 h-10 text-primary animate-spin" />
          <p className="text-sm text-muted-foreground">Loading cluster...</p>
        </div>
      </div>
    )
  }

  if (!cluster) {
    return null
  }

  const clusterConfig = {
    services: cluster.config?.services || {},
  }

  const serviceCount = Object.keys(clusterConfig.services).length

  return (
    <div className="flex flex-col h-screen bg-background">
      {/* Top Bar */}
      <div className="bg-card border-b border-border px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <button
              onClick={() => navigate('/clusters')}
              className="p-1.5 hover:bg-muted/50 rounded-md transition-colors"
            >
              <ArrowLeft className="w-5 h-5 text-muted-foreground" />
            </button>
            <div>
              <h1 className="text-xl font-semibold text-foreground">
                {cluster.name}
              </h1>
              <p className="text-xs text-muted-foreground mt-0.5 font-mono">
                Cluster ID: {cluster.id}
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-3">
            {/* View Mode Toggle */}
            <div className="flex bg-muted/50 rounded-md p-0.5">
              <button
                onClick={() => setViewMode('canvas')}
                className={`flex items-center space-x-1.5 px-3 py-1.5 rounded-sm transition-all text-xs font-medium ${
                  viewMode === 'canvas'
                    ? 'bg-card text-primary shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                <Workflow className="w-3.5 h-3.5" />
                <span>Canvas</span>
              </button>
              <button
                onClick={() => setViewMode('yaml')}
                className={`flex items-center space-x-1.5 px-3 py-1.5 rounded-sm transition-all text-xs font-medium ${
                  viewMode === 'yaml'
                    ? 'bg-card text-primary shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                <FileCode className="w-3.5 h-3.5" />
                <span>YAML</span>
              </button>
            </div>

            {/* Service Count Badge */}
            <div className="px-3 py-1.5 bg-green-500/10 border border-green-500/20 rounded-md">
              <span className="text-xs font-medium text-green-500">
                {serviceCount} service{serviceCount !== 1 ? 's' : ''}
              </span>
            </div>

            {/* Health Check Button */}
            <button className="flex items-center space-x-1.5 px-3 py-1.5 bg-muted/50 text-foreground rounded-md hover:bg-muted transition-colors text-xs font-medium">
              <Activity className="w-3.5 h-3.5" />
              <span>Health Check</span>
            </button>

            {/* Delete Button */}
            <button
              onClick={handleDelete}
              disabled={deleting}
              className="flex items-center space-x-1.5 px-3 py-1.5 bg-destructive/10 text-destructive rounded-md hover:bg-destructive/20 transition-colors disabled:opacity-50 text-xs font-medium"
            >
              {deleting ? (
                <>
                  <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  <span>Deleting...</span>
                </>
              ) : (
                <>
                  <Trash2 className="w-3.5 h-3.5" />
                  <span>Delete</span>
                </>
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 overflow-hidden">
        {viewMode === 'yaml' ? (
          <YamlEditor config={clusterConfig} onChange={() => {}} readOnly />
        ) : (
          <CanvasEditor config={clusterConfig} onChange={() => {}} readOnly />
        )}
      </div>
    </div>
  )
}
