import { useState } from 'react'
import { Activity } from '@/api/client'
import { ChevronDown, ChevronRight, Clock, CheckCircle, XCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

interface ActivityTableProps {
  activities: Activity[]
  loading?: boolean
}

export default function ActivityTable({ activities, loading }: ActivityTableProps) {
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set())

  const toggleRow = (id: string) => {
    const newExpanded = new Set(expandedRows)
    if (newExpanded.has(id)) {
      newExpanded.delete(id)
    } else {
      newExpanded.add(id)
    }
    setExpandedRows(newExpanded)
  }

  const formatDuration = (ms: number) => {
    if (ms < 1) return '<1ms'
    if (ms < 1000) return `${ms}ms`
    return `${(ms / 1000).toFixed(2)}s`
  }

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp)
    return date.toLocaleString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      month: 'short',
      day: 'numeric',
    })
  }

  const truncateCommand = (command: string, maxLength: number = 60) => {
    if (command.length <= maxLength) return command
    return command.substring(0, maxLength) + '...'
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        <span className="ml-3 text-muted-foreground">Loading activities...</span>
      </div>
    )
  }

  if (activities.length === 0) {
    return (
      <div className="text-center py-12 rounded-lg border border-border">
        <Clock className="w-12 h-12 text-muted-foreground mx-auto mb-3" />
        <p className="text-muted-foreground">No activities recorded yet</p>
        <p className="text-sm text-muted-foreground mt-2">
          Activities will appear here when services are accessed
        </p>
      </div>
    )
  }

  return (
    <div className="rounded-lg border border-border overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-border">
          <thead className="bg-secondary">
            <tr>
              <th className="w-10 px-3 py-3"></th>
              <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Time
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Service
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-[#a6a6a6] uppercase tracking-wider">
                Operation
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Command
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-[#a6a6a6] uppercase tracking-wider">
                Duration
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-[#a6a6a6] uppercase tracking-wider">
                Status
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[#333333]">
            {activities.map((activity) => {
              const isExpanded = expandedRows.has(activity.id)
              return (
                <>
                  <tr
                    key={activity.id}
                    className="hover:bg-gray-50 dark:hover:bg-secondary cursor-pointer transition-colors"
                    onClick={() => toggleRow(activity.id)}
                  >
                    <td className="px-3 py-4 whitespace-nowrap text-sm">
                      {isExpanded ? (
                        <ChevronDown className="w-4 h-4 text-gray-500" />
                      ) : (
                        <ChevronRight className="w-4 h-4 text-gray-500" />
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-foreground">
                      {formatTimestamp(activity.timestamp)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <div className="flex flex-col">
                        <span className="font-medium text-foreground">
                          {activity.service_name}
                        </span>
                        <span className="text-xs text-gray-500 dark:text-[#a6a6a6]">
                          {activity.service_type}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={cn(
                          'px-2 py-1 text-xs font-semibold rounded-full',
                          activity.operation === 'GET' || activity.operation === 'SELECT'
                            ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400'
                            : activity.operation === 'SET' ||
                              activity.operation === 'INSERT' ||
                              activity.operation === 'UPDATE'
                            ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
                            : activity.operation === 'DELETE'
                            ? 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400'
                            : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
                        )}
                      >
                        {activity.operation}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-foreground">
                      <code className="font-mono text-xs">
                        {truncateCommand(activity.command)}
                      </code>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-foreground">
                      <span className="font-mono">{formatDuration(activity.duration)}</span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {activity.status === 'success' ? (
                        <div className="flex items-center space-x-1 text-green-600 dark:text-green-400">
                          <CheckCircle className="w-4 h-4" />
                          <span className="text-sm font-medium">Success</span>
                        </div>
                      ) : (
                        <div className="flex items-center space-x-1 text-red-600 dark:text-red-400">
                          <XCircle className="w-4 h-4" />
                          <span className="text-sm font-medium">Error</span>
                        </div>
                      )}
                    </td>
                  </tr>
                  {isExpanded && (
                    <tr className="bg-gray-50 bg-secondary">
                      <td colSpan={7} className="px-6 py-4">
                        <div className="space-y-3">
                          <div>
                            <h4 className="text-xs font-semibold text-gray-500 dark:text-[#a6a6a6] uppercase mb-1">
                              Full Command
                            </h4>
                            <pre className="bg-secondary border border-border rounded p-3 text-sm font-mono overflow-x-auto">
                              {activity.command}
                            </pre>
                          </div>
                          {activity.response && (
                            <div>
                              <h4 className="text-xs font-semibold text-gray-500 dark:text-[#a6a6a6] uppercase mb-1">
                                Response
                              </h4>
                              <pre className="bg-secondary border border-border rounded p-3 text-sm font-mono overflow-x-auto max-h-48 overflow-y-auto">
                                {activity.response}
                              </pre>
                            </div>
                          )}
                          {activity.error && (
                            <div>
                              <h4 className="text-xs font-semibold text-red-600 dark:text-red-400 uppercase mb-1">
                                Error
                              </h4>
                              <pre className="bg-red-50 dark:bg-red-900/10 border border-primary/20 rounded p-3 text-sm font-mono overflow-x-auto text-red-900 dark:text-red-300">
                                {activity.error}
                              </pre>
                            </div>
                          )}
                          <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                              <span className="text-gray-500 dark:text-[#a6a6a6]">
                                Cluster ID:{' '}
                              </span>
                              <span className="font-mono text-foreground">
                                {activity.cluster_id}
                              </span>
                            </div>
                            <div>
                              <span className="text-gray-500 dark:text-[#a6a6a6]">
                                Activity ID:{' '}
                              </span>
                              <span className="font-mono text-xs text-muted-foreground">
                                {activity.id.substring(0, 8)}...
                              </span>
                            </div>
                            {activity.rows_affected !== undefined && (
                              <div>
                                <span className="text-gray-500 dark:text-[#a6a6a6]">
                                  Rows Affected:{' '}
                                </span>
                                <span className="font-mono text-foreground">
                                  {activity.rows_affected}
                                </span>
                              </div>
                            )}
                          </div>
                        </div>
                      </td>
                    </tr>
                  )}
                </>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}

