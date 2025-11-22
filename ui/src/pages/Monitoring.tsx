import { useState, useEffect } from 'react'
import { getActivity, Activity, ActivityFilters as Filters } from '@/api/client'
import { toast } from 'sonner'
import { RefreshCw, Activity as ActivityIcon, TrendingUp, AlertCircle } from 'lucide-react'
import ActivityTable from '@/components/ActivityTable'
import ActivityFilters from '@/components/ActivityFilters'

export default function Monitoring() {
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState<Filters>({ limit: 100 })
  const [autoRefresh, setAutoRefresh] = useState(false)
  const [refreshInterval] = useState(5000) // 5 seconds

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

  // Calculate stats
  const totalActivities = activities.length
  const successCount = activities.filter((a) => a.status === 'success').length
  const errorCount = activities.filter((a) => a.status === 'error').length
  const avgDuration =
    activities.length > 0
      ? activities.reduce((sum, a) => sum + a.duration, 0) / activities.length
      : 0

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Activity Monitoring</h1>
          <p className="text-muted-foreground mt-1">
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
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 dark:peer-focus:ring-primary/30 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary"></div>
            </label>
            <span className="text-sm text-muted-foreground">
              Auto-refresh ({refreshInterval / 1000}s)
            </span>
          </div>

          {/* Manual refresh button */}
          <button
            onClick={handleRefresh}
            disabled={loading}
            className="flex items-center space-x-2 px-4 py-2 bg-white dark:bg-[#1a1a1a] border border-border rounded-lg hover:bg-gray-50 dark:hover:bg-[#262626] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            <span className="text-sm font-medium">Refresh</span>
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="rounded-lg border border-border p-6 shadow-sm">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-500">
                Total Activities
              </p>
              <p className="text-2xl font-bold text-foreground mt-1">
                {totalActivities}
              </p>
            </div>
            <div className="p-3 bg-primary/10 rounded-lg">
              <ActivityIcon className="w-6 h-6 text-primary" />
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-border p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-muted-foreground">
                Successful
              </p>
              <p className="text-2xl font-bold text-success mt-1">
                {successCount}
              </p>
            </div>
            <div className="p-3 bg-success-bg rounded-lg">
              <TrendingUp className="w-6 h-6 text-success" />
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-border p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-muted-foreground">Errors</p>
              <p className="text-2xl font-bold text-error mt-1">
                {errorCount}
              </p>
            </div>
            <div className="p-3 bg-error-bg rounded-lg">
              <AlertCircle className="w-6 h-6 text-error" />
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-border p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-muted-foreground">
                Avg Duration
              </p>
              <p className="text-2xl font-bold text-foreground mt-1">
                {avgDuration < 1000
                  ? `${avgDuration.toFixed(0)}ms`
                  : `${(avgDuration / 1000).toFixed(2)}s`}
              </p>
            </div>
            <div className="p-3 bg-purple-100 dark:bg-purple-900/20 rounded-lg">
              <ActivityIcon className="w-6 h-6 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      <ActivityFilters filters={filters} onFiltersChange={handleFiltersChange} />

      {/* Activity Table */}
      <ActivityTable activities={activities} loading={loading} />

      {/* Info Banner */}
      {activities.length === 0 && !loading && (
        <div className="bg-primary/5 border border-primary rounded-lg p-4">
          <div className="flex items-start space-x-3">
            <ActivityIcon className="w-5 h-5 text-primary mt-0.5" />
            <div>
              <h3 className="text-sm font-semibold text-primary">
                No Activities Yet
              </h3>
              <p className="text-sm text-primary mt-1">
                Activity logs capture every interaction with your services (Redis, PostgreSQL,
                Kafka). Operations will appear here once you start using your clusters.
              </p>
              <ul className="list-disc list-inside text-sm text-primary mt-2 space-y-1">
                <li>All GET, SET, DELETE operations on Redis</li>
                <li>Database queries on PostgreSQL (coming soon)</li>
                <li>Message publish/subscribe on Kafka (coming soon)</li>
              </ul>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
