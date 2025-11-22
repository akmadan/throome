import { useEffect, useState } from 'react'
import { Activity, AlertCircle } from 'lucide-react'
import { checkHealth } from '@/api/client'

export default function ConnectionStatus() {
  const [isConnected, setIsConnected] = useState<boolean | null>(null)
  const [lastCheck, setLastCheck] = useState<Date | null>(null)

  const checkConnection = async () => {
    try {
      await checkHealth()
      setIsConnected(true)
      setLastCheck(new Date())
    } catch (error) {
      setIsConnected(false)
      setLastCheck(new Date())
    }
  }

  useEffect(() => {
    // Initial check
    checkConnection()

    // Check every 10 seconds
    const interval = setInterval(checkConnection, 10000)

    return () => clearInterval(interval)
  }, [])

  if (isConnected === null) {
    return (
      <div className="flex items-center space-x-2 text-muted-foreground">
        <div className="w-2 h-2 bg-gray-400 rounded-full animate-pulse"></div>
        <span className="text-xs">Connecting...</span>
      </div>
    )
  }

  if (!isConnected) {
    return (
      <div className="flex items-center space-x-2 text-red-600 dark:text-red-400">
        <AlertCircle className="w-4 h-4" />
        <span className="text-xs">Disconnected</span>
      </div>
    )
  }

  return (
    <div className="flex items-center space-x-2 text-green-600 dark:text-green-400">
      <div className="relative">
        <div className="w-2 h-2 bg-green-500 rounded-full"></div>
        <div className="absolute inset-0 w-2 h-2 bg-green-500 rounded-full animate-ping opacity-75"></div>
      </div>
      <Activity className="w-4 h-4" />
      <span className="text-xs">
        Connected
        {lastCheck && (
          <span className="text-muted-foreground ml-1">
            ({lastCheck.toLocaleTimeString()})
          </span>
        )}
      </span>
    </div>
  )
}

