import { useState, useEffect } from 'react'
import { ActivityFilters as Filters } from '@/api/client'
import { Filter, X } from 'lucide-react'

interface ActivityFiltersProps {
  filters: Filters
  onFiltersChange: (filters: Filters) => void
  showClusterFilter?: boolean
}

export default function ActivityFilters({
  filters,
  onFiltersChange,
  showClusterFilter = true,
}: ActivityFiltersProps) {
  const [localFilters, setLocalFilters] = useState<Filters>(filters)
  const [showFilters, setShowFilters] = useState(false)

  useEffect(() => {
    setLocalFilters(filters)
  }, [filters])

  const handleApplyFilters = () => {
    onFiltersChange(localFilters)
    setShowFilters(false)
  }

  const handleClearFilters = () => {
    const clearedFilters: Filters = { limit: 100 }
    setLocalFilters(clearedFilters)
    onFiltersChange(clearedFilters)
  }

  const activeFilterCount = Object.keys(localFilters).filter(
    (key) => key !== 'limit' && localFilters[key as keyof Filters]
  ).length

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <button
          onClick={() => setShowFilters(!showFilters)}
          className="flex items-center space-x-2 px-4 py-2 bg-white dark:bg-[#1a1a1a] border border-border rounded-lg hover:bg-gray-50 dark:hover:bg-[#262626] transition-colors"
        >
          <Filter className="w-4 h-4" />
          <span className="text-sm font-medium">Filters</span>
          {activeFilterCount > 0 && (
            <span className="px-2 py-0.5 text-xs font-semibold bg-primary text-white rounded-full">
              {activeFilterCount}
            </span>
          )}
        </button>

        {activeFilterCount > 0 && (
          <button
            onClick={handleClearFilters}
            className="flex items-center space-x-2 px-4 py-2 text-sm text-muted-foreground hover:text-gray-900 dark:hover:text-white transition-colors"
          >
            <X className="w-4 h-4" />
            <span>Clear Filters</span>
          </button>
        )}
      </div>

      {showFilters && (
        <div className="bg-white dark:bg-[#1a1a1a] border border-border rounded-lg p-4 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {showClusterFilter && (
              <div>
                <label
                  htmlFor="cluster-filter"
                  className="block text-sm font-medium text-foreground mb-1"
                >
                  Cluster ID
                </label>
                <input
                  id="cluster-filter"
                  type="text"
                  value={localFilters.cluster_id || ''}
                  onChange={(e) =>
                    setLocalFilters({ ...localFilters, cluster_id: e.target.value || undefined })
                  }
                  placeholder="Filter by cluster..."
                  className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent bg-white dark:bg-[#121212] text-foreground text-sm"
                />
              </div>
            )}

            <div>
              <label
                htmlFor="service-type-filter"
                className="block text-sm font-medium text-foreground mb-1"
              >
                Service Type
              </label>
              <select
                id="service-type-filter"
                value={localFilters.service_type || ''}
                onChange={(e) =>
                  setLocalFilters({ ...localFilters, service_type: e.target.value || undefined })
                }
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent bg-white dark:bg-[#121212] text-foreground text-sm"
              >
                <option value="">All Types</option>
                <option value="redis">Redis</option>
                <option value="postgres">PostgreSQL</option>
                <option value="kafka">Kafka</option>
              </select>
            </div>

            <div>
              <label
                htmlFor="operation-filter"
                className="block text-sm font-medium text-foreground mb-1"
              >
                Operation
              </label>
              <input
                id="operation-filter"
                type="text"
                value={localFilters.operation || ''}
                onChange={(e) =>
                  setLocalFilters({ ...localFilters, operation: e.target.value || undefined })
                }
                placeholder="e.g., GET, SET, SELECT..."
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent bg-white dark:bg-[#121212] text-foreground text-sm"
              />
            </div>

            <div>
              <label
                htmlFor="status-filter"
                className="block text-sm font-medium text-foreground mb-1"
              >
                Status
              </label>
              <select
                id="status-filter"
                value={localFilters.status || ''}
                onChange={(e) =>
                  setLocalFilters({
                    ...localFilters,
                    status: (e.target.value as 'success' | 'error') || undefined,
                  })
                }
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent bg-white dark:bg-[#121212] text-foreground text-sm"
              >
                <option value="">All Statuses</option>
                <option value="success">Success</option>
                <option value="error">Error</option>
              </select>
            </div>

            <div>
              <label
                htmlFor="limit-filter"
                className="block text-sm font-medium text-foreground mb-1"
              >
                Limit
              </label>
              <select
                id="limit-filter"
                value={localFilters.limit || 100}
                onChange={(e) =>
                  setLocalFilters({ ...localFilters, limit: parseInt(e.target.value) })
                }
                className="w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent bg-white dark:bg-[#121212] text-foreground text-sm"
              >
                <option value="25">25 activities</option>
                <option value="50">50 activities</option>
                <option value="100">100 activities</option>
                <option value="250">250 activities</option>
                <option value="500">500 activities</option>
              </select>
            </div>
          </div>

          <div className="flex items-center justify-end space-x-3 pt-3 border-t border-border">
            <button
              onClick={() => setShowFilters(false)}
              className="px-4 py-2 text-sm text-muted-foreground hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleApplyFilters}
              className="px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors text-sm font-medium"
            >
              Apply Filters
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

