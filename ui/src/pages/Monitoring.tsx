import { useState, useEffect } from 'react'
import { getActivity, Activity, ActivityFilters as Filters } from '@/api/client'
import { toast } from 'sonner'
import { RefreshCw, X, Clock, CheckCircle2, XCircle, Code, Copy, Check } from 'lucide-react'
import ActivityFilters from '@/components/ActivityFilters'

export default function Monitoring() {
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState<Filters>({ limit: 100 })
  const [autoRefresh, setAutoRefresh] = useState(false)
  const [refreshInterval] = useState(5000) // 5 seconds
  const [selectedActivity, setSelectedActivity] = useState<Activity | null>(null)
  const [copied, setCopied] = useState(false)

  const fetchActivities = async () => {
    try {
      setLoading(true)
      const response = await getActivity(filters)
      setActivities(response.activities)
    } catch (error) {
      toast.error('Failed to fetch activity logs', {
        description: error instanceof Error ? error.message : 'An unknown error occurred',
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchActivities()
  }, [filters])

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(() => {
        fetchActivities()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval, filters])

  const handleRefresh = () => {
    fetchActivities()
    toast.info('Activity logs refreshed')
  }

  const handleFiltersChange = (newFilters: Filters) => {
    setFilters(newFilters)
  }

  const handleCopyResponse = () => {
    if (selectedActivity?.response) {
      navigator.clipboard.writeText(selectedActivity.response)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
      toast.success('Response copied to clipboard')
    }
  }

  const isJsonResponse = (response: string): boolean => {
    try {
      JSON.parse(response)
      return true
    } catch {
      return false
    }
  }

  const formatJsonResponse = (response: string): string => {
    try {
      return JSON.stringify(JSON.parse(response), null, 2)
    } catch {
      return response
    }
  }

  // Calculate stats
  const totalActivities = activities.length
  const successCount = activities.filter((a) => a.status === 'success').length
  const errorCount = activities.filter((a) => a.status === 'error').length
  const avgDuration =
    activities.length > 0
      ? activities.reduce((sum, a) => sum + a.duration, 0) / activities.length
      : 0

  const formatDuration = (ns: number) => {
    const ms = ns / 1000000
    if (ms < 1) return '<1ms'
    if (ms < 1000) return `${Math.round(ms)}ms`
    return `${(ms / 1000).toFixed(2)}s`
  }

  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp)
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: true,
    })
  }

  return (
    <div className="flex h-screen overflow-hidden">
      {/* Main Content */}
      <div className={`flex-1 flex flex-col ${selectedActivity ? 'mr-[600px]' : ''} transition-all duration-300`}>
        <div className="p-6 space-y-6">
          {/* Header */}
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-foreground">Activity Monitoring</h1>
              <p className="text-sm text-muted-foreground mt-1">
                Real-time activity logs for all service interactions
              </p>
            </div>
            <div className="flex items-center space-x-3">
              {/* Auto-refresh toggle */}
              <div className="flex items-center space-x-2">
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={autoRefresh}
                    onChange={(e) => setAutoRefresh(e.target.checked)}
                    className="sr-only peer"
                  />
                  <div className="w-11 h-6 bg-border peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                </label>
                <span className="text-xs text-muted-foreground">
                  Auto-refresh ({refreshInterval / 1000}s)
                </span>
              </div>

              {/* Manual refresh button */}
              <button
                onClick={handleRefresh}
                disabled={loading}
                className="flex items-center space-x-2 px-3 py-2 bg-card border border-border rounded-lg hover:bg-accent transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                <span className="text-sm font-medium">Refresh</span>
              </button>
            </div>
          </div>

          {/* Simple Stats - Key Value */}
          <div className="flex items-center space-x-6 text-sm">
            <div className="flex items-center space-x-2">
              <span className="text-muted-foreground">Total:</span>
              <span className="font-semibold text-foreground">{totalActivities}</span>
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-muted-foreground">Success:</span>
              <span className="font-semibold text-success">{successCount}</span>
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-muted-foreground">Errors:</span>
              <span className="font-semibold text-error">{errorCount}</span>
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-muted-foreground">Avg Duration:</span>
              <span className="font-semibold text-foreground">{formatDuration(avgDuration)}</span>
            </div>
          </div>

          {/* Filters */}
          <ActivityFilters filters={filters} onFiltersChange={handleFiltersChange} />

          {/* Activity Table */}
          <div className="bg-card border border-border rounded-lg overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-secondary border-b border-border">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider w-32">
                      Service
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider w-28">
                      Operation
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                      Command
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider w-24">
                      Duration
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider w-24">
                      Status
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wider w-28">
                      Time
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {loading ? (
                    <tr>
                      <td colSpan={6} className="px-4 py-12 text-center text-sm text-muted-foreground">
                        <RefreshCw className="w-5 h-5 animate-spin mx-auto mb-2" />
                        Loading activities...
                      </td>
                    </tr>
                  ) : activities.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="px-4 py-12 text-center text-sm text-muted-foreground">
                        No activities found
                      </td>
                    </tr>
                  ) : (
                    activities.map((activity) => (
                      <tr
                        key={activity.id}
                        onClick={() => setSelectedActivity(activity)}
                        className="hover:bg-accent cursor-pointer transition-colors"
                      >
                        <td className="px-4 py-3">
                          <div>
                            <div className="font-medium text-sm text-foreground">{activity.service_name}</div>
                            <div className="text-xs text-muted-foreground">{activity.service_type}</div>
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <span className="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-secondary text-foreground">
                            {activity.operation}
                          </span>
                        </td>
                        <td className="px-4 py-3">
                          <div className="font-mono text-sm text-foreground truncate max-w-md" title={activity.command}>
                            {activity.command}
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <span className="text-sm text-muted-foreground">{formatDuration(activity.duration)}</span>
                        </td>
                        <td className="px-4 py-3">
                          {activity.status === 'success' ? (
                            <span className="inline-flex items-center space-x-1 text-success text-sm">
                              <CheckCircle2 className="w-4 h-4" />
                              <span>Success</span>
                            </span>
                          ) : (
                            <span className="inline-flex items-center space-x-1 text-error text-sm">
                              <XCircle className="w-4 h-4" />
                              <span>Error</span>
                            </span>
                          )}
                        </td>
                        <td className="px-4 py-3">
                          <span className="text-xs text-muted-foreground">{formatTime(activity.timestamp)}</span>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      {/* Side Panel */}
      {selectedActivity && (
        <div className="fixed right-0 top-0 h-screen w-[600px] bg-card border-l border-border shadow-2xl overflow-y-auto">
          <div className="sticky top-0 bg-card border-b border-border px-6 py-4 flex items-center justify-between z-10">
            <div>
              <h2 className="text-lg font-semibold text-foreground">Activity Details</h2>
              <p className="text-xs text-muted-foreground mt-0.5">
                {selectedActivity.service_name} • {selectedActivity.service_type}
              </p>
            </div>
            <button
              onClick={() => setSelectedActivity(null)}
              className="p-2 hover:bg-accent rounded-lg transition-colors"
            >
              <X className="w-5 h-5 text-muted-foreground" />
            </button>
          </div>

          <div className="p-6 space-y-6">
            {/* Status Banner */}
            <div
              className={`px-4 py-3 rounded-lg border ${
                selectedActivity.status === 'success'
                  ? 'bg-success/5 border-success/20'
                  : 'bg-error/5 border-error/20'
              }`}
            >
              <div className="flex items-center space-x-2">
                {selectedActivity.status === 'success' ? (
                  <>
                    <CheckCircle2 className="w-5 h-5 text-success" />
                    <span className="font-semibold text-success">Success</span>
                  </>
                ) : (
                  <>
                    <XCircle className="w-5 h-5 text-error" />
                    <span className="font-semibold text-error">Error</span>
                  </>
                )}
                <span className="text-muted-foreground">•</span>
                <Clock className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm text-muted-foreground">{formatDuration(selectedActivity.duration)}</span>
              </div>
            </div>

            {/* Operation */}
            <div>
              <label className="block text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                Operation
              </label>
              <div className="bg-secondary px-4 py-2 rounded-lg">
                <span className="font-mono text-sm text-foreground">{selectedActivity.operation}</span>
              </div>
            </div>

            {/* Command */}
            <div>
              <label className="block text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                Command
              </label>
              <div className="bg-secondary px-4 py-3 rounded-lg">
                <code className="font-mono text-sm text-foreground break-all">{selectedActivity.command}</code>
              </div>
            </div>

            {/* Parameters */}
            {selectedActivity.parameters && selectedActivity.parameters.length > 0 && (
              <div>
                <label className="block text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                  Parameters
                </label>
                <div className="bg-secondary px-4 py-3 rounded-lg">
                  <code className="font-mono text-sm text-foreground">
                    {JSON.stringify(selectedActivity.parameters, null, 2)}
                  </code>
                </div>
              </div>
            )}

            {/* Response/Output */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                  {selectedActivity.error ? 'Error' : 'Response'}
                </label>
                {selectedActivity.response && (
                  <button
                    onClick={handleCopyResponse}
                    className="flex items-center space-x-1 px-2 py-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {copied ? (
                      <>
                        <Check className="w-3 h-3" />
                        <span>Copied!</span>
                      </>
                    ) : (
                      <>
                        <Copy className="w-3 h-3" />
                        <span>Copy</span>
                      </>
                    )}
                  </button>
                )}
              </div>
              <div className="bg-[#1e1e1e] border border-border rounded-lg overflow-hidden">
                <div className="flex items-center justify-between px-4 py-2 bg-[#2d2d2d] border-b border-border">
                  <div className="flex items-center space-x-2">
                    <Code className="w-4 h-4 text-muted-foreground" />
                    <span className="text-xs font-medium text-muted-foreground">
                      {isJsonResponse(selectedActivity.response || selectedActivity.error || '')
                        ? 'JSON Output'
                        : 'Text Output'}
                    </span>
                  </div>
                </div>
                <div className="p-4 overflow-x-auto max-h-96">
                  <pre className="font-mono text-xs text-gray-300 leading-relaxed">
                    {selectedActivity.error ||
                      (isJsonResponse(selectedActivity.response || '')
                        ? formatJsonResponse(selectedActivity.response || '')
                        : selectedActivity.response || '(empty)')}
                  </pre>
                </div>
              </div>
            </div>

            {/* Metadata */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                  Cluster ID
                </label>
                <div className="bg-secondary px-3 py-2 rounded-lg">
                  <span className="font-mono text-xs text-foreground">{selectedActivity.cluster_id}</span>
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                  Timestamp
                </label>
                <div className="bg-secondary px-3 py-2 rounded-lg">
                  <span className="text-xs text-foreground">
                    {new Date(selectedActivity.timestamp).toLocaleString()}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
