import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Boxes,
  Database,
  Activity,
  GitBranch,
  Settings,
  ChevronRight,
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
    <div className="flex flex-col w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
      {/* Logo */}
      <div className="flex items-center justify-between h-16 px-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center space-x-3">
          <img src="/text_logo.svg" alt="throome" className="h-10" />
       
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            end={item.href === '/'}
            className={({ isActive }) =>
              cn(
                'flex items-center justify-between px-4 py-3 text-sm font-medium rounded-lg transition-colors group',
                isActive
                  ? 'bg-red-50 text-[#FF5050] dark:bg-red-900/20 dark:text-[#FF5050]'
                  : 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700'
              )
            }
          >
            {({ isActive }) => (
              <>
                <div className="flex items-center space-x-3">
                  <item.icon className="w-5 h-5" />
                  <span>{item.name}</span>
                </div>
                {isActive && (
                  <ChevronRight className="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity" />
                )}
              </>
            )}
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 space-y-3">
        <ConnectionStatus />
        <div className="flex items-center space-x-3">
          <div className="w-8 h-8 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center">
            <span className="text-xs font-medium text-gray-600 dark:text-gray-400">v0.1</span>
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
              throome gateway
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">Port :9000</p>
          </div>
        </div>
      </div>
    </div>
  )
}

