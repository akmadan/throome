import { Activity, Boxes, Database, TrendingUp } from 'lucide-react'

export default function Dashboard() {
  const stats = [
    { title: 'Total Clusters', value: '3', icon: Boxes, color: 'text-primary' },
    { title: 'Active Services', value: '12', icon: Database, color: 'text-blue-400' },
    { title: 'Health Status', value: '98.5%', icon: Activity, color: 'text-green-400' },
    { title: 'Requests/sec', value: '1,234', icon: TrendingUp, color: 'text-purple-400' },
  ]

  const recentActivity = [
    { action: 'Created cluster "production-us-east"', time: '2 min ago' },
    { action: 'Redis service health check passed', time: '5 min ago' },
    { action: 'Updated routing strategy to weighted', time: '12 min ago' },
    { action: 'New PostgreSQL service added', time: '1 hour ago' },
  ]

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl font-semibold text-foreground">Dashboard</h1>
        <p className="text-sm text-muted-foreground mt-0.5">
          Welcome to Throome Gateway Dashboard
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, index) => (
          <div
            key={index}
            className="bg-card border border-border rounded-lg p-4 hover:border-primary/30 transition-all"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs text-muted-foreground mb-1">{stat.title}</p>
                <p className="text-2xl font-semibold text-foreground">{stat.value}</p>
              </div>
              <div className={`w-10 h-10 bg-muted/30 rounded-md flex items-center justify-center ${stat.color}`}>
                <stat.icon className="w-5 h-5" />
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Quick Actions & Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {/* Quick Actions */}
        <div className="bg-card border border-border rounded-lg p-5">
          <h2 className="text-sm font-semibold text-foreground mb-4">
            Quick Actions
          </h2>
          <div className="space-y-2">
            <button className="w-full flex items-center justify-between px-3.5 py-2.5 bg-primary/10 text-primary rounded-md hover:bg-primary/15 transition-colors text-sm font-medium">
              <span>Create New Cluster</span>
              <span className="text-lg">+</span>
            </button>
            <button className="w-full flex items-center justify-between px-3.5 py-2.5 bg-muted/30 text-foreground rounded-md hover:bg-muted/50 transition-colors text-sm">
              <span>View All Services</span>
              <span className="text-muted-foreground">→</span>
            </button>
            <button className="w-full flex items-center justify-between px-3.5 py-2.5 bg-muted/30 text-foreground rounded-md hover:bg-muted/50 transition-colors text-sm">
              <span>Check System Health</span>
              <span className="text-muted-foreground">→</span>
            </button>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="bg-card border border-border rounded-lg p-5">
          <h2 className="text-sm font-semibold text-foreground mb-4">
            Recent Activity
          </h2>
          <div className="space-y-3">
            {recentActivity.map((activity, index) => (
              <div key={index} className="flex items-start space-x-3">
                <div className="w-1.5 h-1.5 mt-2 bg-primary rounded-full flex-shrink-0"></div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-foreground">{activity.action}</p>
                  <p className="text-xs text-muted-foreground mt-0.5">{activity.time}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
