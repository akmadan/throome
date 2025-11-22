import { Bell, Search, Moon, Sun } from 'lucide-react'
import { useState, useEffect } from 'react'

export default function Header() {
  const [darkMode, setDarkMode] = useState(() => {
    // Initialize from localStorage or default to dark mode
    const saved = localStorage.getItem('throome-theme')
    if (saved) {
      return saved === 'dark'
    }
    // Default to dark mode (since we're using Docker Desktop dark theme)
    return true
  })

  // Apply theme on mount
  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [])

  const toggleDarkMode = () => {
    const newMode = !darkMode
    setDarkMode(newMode)
    
    // Update DOM
    if (newMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
    
    // Persist to localStorage
    localStorage.setItem('throome-theme', newMode ? 'dark' : 'light')
  }

  return (
    <header className="h-14 bg-card border-b border-border flex items-center justify-between px-4">
      {/* Search */}
      <div className="flex-1 max-w-md">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search..."
            className="w-full pl-9 pr-4 py-1.5 text-sm bg-muted/50 border border-transparent rounded-md focus:outline-none focus:ring-1 focus:ring-primary/50 focus:bg-background text-foreground placeholder:text-muted-foreground"
          />
        </div>
      </div>

      {/* Actions */}
      <div className="flex items-center space-x-2">
        {/* Dark Mode Toggle */}
        <button
          onClick={toggleDarkMode}
          className="p-2 text-muted-foreground hover:text-foreground rounded-md hover:bg-muted/50 transition-colors"
          aria-label="Toggle dark mode"
        >
          {darkMode ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
        </button>

        {/* Notifications */}
        <button
          className="relative p-2 text-muted-foreground hover:text-foreground rounded-md hover:bg-muted/50 transition-colors"
          aria-label="Notifications"
        >
          <Bell className="w-4 h-4" />
          <span className="absolute top-1.5 right-1.5 w-1.5 h-1.5 bg-primary rounded-full"></span>
        </button>

        {/* User Avatar */}
        <div className="w-7 h-7 bg-muted rounded-full flex items-center justify-center ml-2">
          <span className="text-xs font-medium text-foreground">AM</span>
        </div>
      </div>
    </header>
  )
}

