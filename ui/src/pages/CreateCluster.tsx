import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft, FileCode, Workflow, Save, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import YamlEditor from '@/components/YamlEditor'
import CanvasEditor from '@/components/CanvasEditor'
import { createCluster } from '@/api/client'

type ViewMode = 'yaml' | 'canvas'

export default function CreateCluster() {
  const navigate = useNavigate()
  const [viewMode, setViewMode] = useState<ViewMode>('canvas')
  const [clusterName, setClusterName] = useState('')
  const [clusterConfig, setClusterConfig] = useState({
    services: {} as Record<string, any>,
  })
  const [isCreating, setIsCreating] = useState(false)

  const handleCreate = async () => {
    if (!clusterName.trim()) {
      toast.error('Please enter a cluster name')
      return
    }

    if (Object.keys(clusterConfig.services).length === 0) {
      toast.error('Please add at least one service')
      return
    }

    setIsCreating(true)

    try {
      await createCluster({
        name: clusterName,
        config: clusterConfig,
      })

      toast.success(`Cluster "${clusterName}" created successfully!`)
      navigate('/clusters')
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Failed to create cluster')
    } finally {
      setIsCreating(false)
    }
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
                Create New Cluster
              </h1>
              <p className="text-xs text-muted-foreground mt-0.5">
                Design your infrastructure with visual canvas or YAML
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

            {/* Cluster Name Input */}
            <input
              type="text"
              value={clusterName}
              onChange={(e) => setClusterName(e.target.value)}
              placeholder="Cluster name (e.g., production-us-east)"
              className="w-72 px-3 py-1.5 border border-border rounded-md focus:ring-1 focus:ring-primary/50 focus:border-primary bg-background text-foreground text-sm placeholder:text-muted-foreground"
            />

            {/* Service Count Badge */}
            <div className="px-3 py-1.5 bg-muted/50 rounded-md">
              <span className="text-xs font-medium text-foreground">
                {serviceCount} service{serviceCount !== 1 ? 's' : ''}
              </span>
            </div>

            {/* Create Button */}
            <button
              onClick={handleCreate}
              disabled={isCreating}
              className="flex items-center space-x-2 px-4 py-1.5 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm font-medium"
            >
              {isCreating ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  <span>Creating...</span>
                </>
              ) : (
                <>
                  <Save className="w-4 h-4" />
                  <span>Create Cluster</span>
                </>
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 overflow-hidden">
        {viewMode === 'yaml' ? (
          <YamlEditor config={clusterConfig} onChange={setClusterConfig} />
        ) : (
          <CanvasEditor config={clusterConfig} onChange={setClusterConfig} />
        )}
      </div>
    </div>
  )
}
