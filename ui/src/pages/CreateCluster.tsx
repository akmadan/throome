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
                Create New Cluster
              </h1>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                Design your infrastructure with visual canvas or YAML
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
                    ? 'bg-white dark:bg-gray-800 text-[#FF5050] dark:text-[#FF5050] shadow-sm'
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
                    ? 'bg-white dark:bg-gray-800 text-[#FF5050] dark:text-[#FF5050] shadow-sm'
                    : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-300'
                }`}
              >
                <FileCode className="w-4 h-4" />
                <span className="font-medium text-sm">YAML</span>
              </button>
            </div>

            {/* Cluster Name Input */}
            <input
              type="text"
              value={clusterName}
              onChange={(e) => setClusterName(e.target.value)}
              placeholder="Cluster name (e.g., production-us-east)"
              className="w-80 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-[#FF5050] focus:border-transparent bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
            />

            {/* Service Count Badge */}
            <div className="px-4 py-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
              <span className="text-sm font-medium text-[#FF5050] dark:text-[#FF5050]">
                {Object.keys(clusterConfig.services).length} service(s)
              </span>
            </div>

            {/* Create Button */}
            <button
              onClick={handleCreate}
              disabled={isCreating}
              className="flex items-center space-x-2 px-6 py-2 bg-[#FF5050] text-white rounded-lg hover:bg-[#ed1515] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
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
        {viewMode === 'canvas' ? (
          <CanvasEditor config={clusterConfig} onChange={setClusterConfig} />
        ) : (
          <div className="h-full p-6">
            <YamlEditor config={clusterConfig} onChange={setClusterConfig} />
          </div>
        )}
      </div>
    </div>
  )
}

