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
      <div className="flex items-center justify-center h-screen bg-gray-50 dark:bg-gray-900">
        <div className="flex flex-col items-center space-y-4">
          <Loader2 className="w-12 h-12 text-[#FF5050] animate-spin" />
          <p className="text-gray-600 dark:text-gray-400">Loading cluster...</p>
        </div>
      </div>
    )
  }

  if (!cluster) {
    return null
  }

  // Convert cluster config to the format expected by editors
  const clusterConfig = {
    services: cluster.config?.services || {},
  }

  return (
    <div className="flex flex-col h-screen bg-gray-50 dark:bg-gray-900">
      {/* Top Bar */}
      <div className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <button
              onClick={() => navigate('/clusters')}
              className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
            >
              <ArrowLeft className="w-5 h-5 text-gray-600 dark:text-gray-400" />
            </button>
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                {cluster.name}
              </h1>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Cluster ID: {cluster.id}
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {/* View Mode Toggle */}
            <div className="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
              <button
                onClick={() => setViewMode('canvas')}
                className={`flex items-center space-x-2 px-4 py-2 rounded-md transition-all ${
                  viewMode === 'canvas'
                    ? 'bg-white dark:bg-gray-800 text-[#FF5050] shadow-sm'
                    : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
                }`}
              >
                <Workflow className="w-4 h-4" />
                <span className="font-medium text-sm">Canvas</span>
              </button>
              <button
                onClick={() => setViewMode('yaml')}
                className={`flex items-center space-x-2 px-4 py-2 rounded-md transition-all ${
                  viewMode === 'yaml'
                    ? 'bg-white dark:bg-gray-800 text-[#FF5050] shadow-sm'
                    : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
                }`}
              >
                <FileCode className="w-4 h-4" />
                <span className="font-medium text-sm">YAML</span>
              </button>
            </div>

            {/* Service Count Badge */}
            <div className="px-4 py-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
              <span className="text-sm font-medium text-green-600 dark:text-green-400">
                {cluster.services?.length || 0} service(s)
              </span>
            </div>

            {/* Health Check Button */}
            <button className="flex items-center space-x-2 px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors">
              <Activity className="w-4 h-4" />
              <span>Health Check</span>
            </button>

            {/* Delete Button */}
            <button
              onClick={handleDelete}
              disabled={deleting}
              className="flex items-center space-x-2 px-4 py-2 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 rounded-lg hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors disabled:opacity-50"
            >
              {deleting ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  <span>Deleting...</span>
                </>
              ) : (
                <>
                  <Trash2 className="w-4 h-4" />
                  <span>Delete</span>
                </>
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 overflow-hidden">
        {viewMode === 'canvas' ? (
          <CanvasEditor
            config={clusterConfig}
            onChange={() => {}} // Read-only mode
            readOnly={true}
          />
        ) : (
          <div className="h-full p-6">
            <YamlEditor
              config={clusterConfig}
              onChange={() => {}} // Read-only mode
              readOnly={true}
            />
          </div>
        )}
      </div>
    </div>
  )
}

