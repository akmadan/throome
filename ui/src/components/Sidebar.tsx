import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Boxes,
  Database,
  Activity,
  GitBranch,
  Settings,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import ConnectionStatus from './ConnectionStatus'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Clusters', href: '/clusters', icon: Boxes },
  { name: 'Services', href: '/services', icon: Database },
  { name: 'Monitoring', href: '/monitoring', icon: Activity },
  { name: 'Routing', href: '/routing', icon: GitBranch },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export default function Sidebar() {
  return (
    <div className="flex flex-col w-56 bg-card border-r border-border">
      {/* Logo */}
      <div className="flex items-center h-14 px-4 border-b border-border">
        <div className="flex items-center space-x-3">
          <img
            src="/text_logo.svg"
            alt="throome"
            className="h-6 hidden dark:block"
          />
          <img
            src="/text_logo.svg"
            alt="throome"
            className="h-6 dark:hidden"
          />
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-2 py-2 space-y-0.5 overflow-y-auto">
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            end={item.href === '/'}
            className={({ isActive }) =>
              cn(
                'flex items-center px-3 py-2.5 text-[13px] font-normal rounded-md transition-all',
                isActive
                  ? 'bg-accent/80 text-foreground font-medium'
                  : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'
              )
            }
          >
            {() => (
              <>
                <div className="flex items-center space-x-3">
                  <item.icon className="w-4 h-4" />
                  <span>{item.name}</span>
                </div>
              </>
            )}
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div className="px-3 py-3 border-t border-border">
        <ConnectionStatus />
        <div className="mt-2 px-2 py-1.5 bg-muted/30 rounded-md">
          <p className="text-xs font-medium text-muted-foreground">
            throome v0.1 • :9000
          </p>
        </div>
      </div>
    </div>
  )
}

