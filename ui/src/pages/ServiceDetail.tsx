import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, RefreshCw, Download, Terminal, Activity as ActivityIcon, Info } from 'lucide-react'
import { toast } from 'sonner'
import api from '@/api/client'

interface ServiceInfo {
  cluster_id: string
  cluster_name: string
  service_name: string
  type: string
  host: string
  port: number
  container_id?: string
  container_status?: string
  container_running?: boolean
  container_started_at?: string
  container_image?: string
  database?: string
  username?: string
}

export default function ServiceDetail() {
  const { clusterId, serviceName } = useParams<{ clusterId: string; serviceName: string }>()
  const navigate = useNavigate()
  const [serviceInfo, setServiceInfo] = useState<ServiceInfo | null>(null)
  const [logs, setLogs] = useState<string>('')
  const [loading, setLoading] = useState(true)
  const [logsLoading, setLogsLoading] = useState(false)
  const [autoRefresh, setAutoRefresh] = useState(false)
  const [tailLines, setTailLines] = useState(100)

  const fetchServiceInfo = async () => {
    if (!clusterId || !serviceName) return

    try {
      const response = await api.get(`/clusters/${clusterId}/services/${serviceName}`)
      setServiceInfo(response.data)
    } catch (error) {
      toast.error('Failed to fetch service info', {
        description: error instanceof Error ? error.message : 'An unknown error occurred',
      })
    } finally {
      setLoading(false)
    }
  }

  const fetchLogs = async () => {
    if (!clusterId || !serviceName) return

    try {
      setLogsLoading(true)
      const response = await api.get(
        `/clusters/${clusterId}/services/${serviceName}/logs?tail=${tailLines}&timestamps=true`,
        {
          responseType: 'text',
        }
      )
      setLogs(response.data as string)
    } catch (error) {
      toast.error('Failed to fetch logs', {
        description: error instanceof Error ? error.message : 'An unknown error occurred',
      })
    } finally {
      setLogsLoading(false)
    }
  }

  useEffect(() => {
    fetchServiceInfo()
    fetchLogs()
  }, [clusterId, serviceName])

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(() => {
        fetchLogs()
      }, 5000)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, tailLines])

  const handleDownloadLogs = () => {
    const blob = new Blob([logs], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${serviceName}-logs-${new Date().toISOString()}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    toast.success('Logs downloaded!')
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  if (!serviceInfo) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">Service not found</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate('/services')}
            className="p-2 hover:bg-accent rounded-lg transition-colors"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <div>
            <h1 className="text-3xl font-bold text-foreground">{serviceName}</h1>
            <p className="text-muted-foreground mt-1">
              {serviceInfo.cluster_name} • {serviceInfo.type}
            </p>
          </div>
        </div>
      </div>

      {/* Service Info Card */}
      <div className="bg-card rounded-lg border border-border p-6">
        <div className="flex items-center space-x-2 mb-4">
          <Info className="w-5 h-5 text-primary" />
          <h2 className="text-lg font-semibold text-foreground">Service Information</h2>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p className="text-sm text-muted-foreground">Type</p>
            <p className="text-sm font-medium text-foreground mt-1">{serviceInfo.type}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Host</p>
            <p className="text-sm font-mono text-foreground mt-1">{serviceInfo.host}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Port</p>
            <p className="text-sm font-mono text-foreground mt-1">{serviceInfo.port}</p>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Status</p>
            <p className="text-sm font-medium mt-1">
              {serviceInfo.container_running ? (
                <span className="text-green-600 dark:text-green-400">Running</span>
              ) : (
                <span className="text-red-600 dark:text-red-400">Stopped</span>
              )}
            </p>
          </div>
          {serviceInfo.container_id && (
            <div className="col-span-2">
              <p className="text-sm text-muted-foreground">Container ID</p>
              <p className="text-xs font-mono text-foreground mt-1">
                {serviceInfo.container_id.substring(0, 12)}...
              </p>
            </div>
          )}
          {serviceInfo.container_image && (
            <div className="col-span-2">
              <p className="text-sm text-muted-foreground">Image</p>
              <p className="text-sm font-mono text-foreground mt-1">{serviceInfo.container_image}</p>
            </div>
          )}
          {serviceInfo.database && (
            <>
              <div>
                <p className="text-sm text-muted-foreground">Database</p>
                <p className="text-sm font-mono text-foreground mt-1">{serviceInfo.database}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Username</p>
                <p className="text-sm font-mono text-foreground mt-1">{serviceInfo.username}</p>
              </div>
            </>
          )}
        </div>
      </div>

      {/* Container Logs */}
      <div className="bg-card rounded-lg border border-border">
        {/* Logs Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-border">
          <div className="flex items-center space-x-2">
            <Terminal className="w-5 h-5 text-primary" />
            <h2 className="text-lg font-semibold text-foreground">Container Logs</h2>
          </div>
          <div className="flex items-center space-x-3">
            {/* Tail lines selector */}
            <select
              value={tailLines}
              onChange={(e) => setTailLines(Number(e.target.value))}
              className="px-3 py-1.5 text-sm border border-border rounded-lg bg-white dark:bg-gray-700 text-foreground"
            >
              <option value="50">Last 50 lines</option>
              <option value="100">Last 100 lines</option>
              <option value="250">Last 250 lines</option>
              <option value="500">Last 500 lines</option>
              <option value="1000">Last 1000 lines</option>
            </select>

            {/* Auto-refresh toggle */}
            <label className="flex items-center space-x-2 cursor-pointer">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
              />
              <span className="text-sm text-foreground">Auto-refresh</span>
            </label>

            {/* Download button */}
            <button
              onClick={handleDownloadLogs}
              disabled={!logs}
              className="flex items-center space-x-1 px-3 py-1.5 text-sm bg-muted text-foreground rounded-lg hover:bg-accent transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Download className="w-4 h-4" />
              <span>Download</span>
            </button>

            {/* Refresh button */}
            <button
              onClick={fetchLogs}
              disabled={logsLoading}
              className="flex items-center space-x-1 px-3 py-1.5 text-sm bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <RefreshCw className={`w-4 h-4 ${logsLoading ? 'animate-spin' : ''}`} />
              <span>Refresh</span>
            </button>
          </div>
        </div>

        {/* Logs Content */}
        <div className="p-4 bg-black">
          <pre className="text-xs text-green-400 font-mono overflow-x-auto whitespace-pre-wrap break-words min-h-[400px] max-h-[600px] overflow-y-auto">
            {logsLoading ? (
              <div className="flex items-center justify-center py-12">
                <RefreshCw className="w-6 h-6 animate-spin" />
                <span className="ml-2">Loading logs...</span>
              </div>
            ) : logs ? (
              logs
            ) : (
              <div className="flex items-center justify-center py-12 text-gray-500">
                No logs available
              </div>
            )}
          </pre>
        </div>
      </div>

      {/* Activity Link */}
      <div className="bg-primary/5 border border-primary rounded-lg p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <ActivityIcon className="w-5 h-5 text-primary" />
            <div>
              <h3 className="text-sm font-semibold text-primary">
                View Activity Logs
              </h3>
              <p className="text-sm text-primary mt-1">
                See all interactions with this service in the Monitoring page
              </p>
            </div>
          </div>
          <button
            onClick={() => navigate('/monitoring')}
            className="px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors text-sm"
          >
            View Activities
          </button>
        </div>
      </div>
    </div>
  )
}

