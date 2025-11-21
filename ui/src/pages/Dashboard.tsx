import { Activity, Boxes, Database, TrendingUp } from 'lucide-react'
import StatsCard from '@/components/StatsCard'

export default function Dashboard() {
  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Welcome to Throome Gateway Dashboard
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard
          title="Total Clusters"
          value="3"
          change="+2 this month"
          trend="up"
          icon={Boxes}
          color="blue"
        />
        <StatsCard
          title="Active Services"
          value="12"
          change="3 Redis, 5 PostgreSQL, 4 Kafka"
          icon={Database}
          color="green"
        />
        <StatsCard
          title="Health Status"
          value="98.5%"
          change="All systems operational"
          trend="up"
          icon={Activity}
          color="emerald"
        />
        <StatsCard
          title="Requests/sec"
          value="1,234"
          change="+12% from last hour"
          trend="up"
          icon={TrendingUp}
          color="purple"
        />
      </div>

      {/* Quick Actions & Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Quick Actions */}
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Quick Actions
          </h2>
          <div className="space-y-3">
            <button className="w-full flex items-center justify-between px-4 py-3 bg-red-50 dark:bg-red-900/20 text-[#FF5050] dark:text-[#FF5050] rounded-lg hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors">
              <span className="font-medium">Create New Cluster</span>
              <span className="text-2xl">+</span>
            </button>
            <button className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors">
              <span className="font-medium">View All Services</span>
              <span className="text-xl">→</span>
            </button>
            <button className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors">
              <span className="font-medium">Check System Health</span>
              <span className="text-xl">→</span>
            </button>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Recent Activity
          </h2>
          <div className="space-y-4">
            {[
              { type: 'cluster', action: 'Created cluster "production-us-east"', time: '2 min ago' },
              { type: 'service', action: 'Redis service health check passed', time: '5 min ago' },
              { type: 'routing', action: 'Updated routing strategy to weighted', time: '12 min ago' },
              { type: 'service', action: 'New PostgreSQL service added', time: '1 hour ago' },
            ].map((activity, index) => (
              <div key={index} className="flex items-start space-x-3">
                <div className="w-2 h-2 mt-2 bg-red-500 rounded-full"></div>
                <div className="flex-1">
                  <p className="text-sm text-gray-900 dark:text-white">{activity.action}</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">{activity.time}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

