import { useState } from 'react'
import { Plus, Database, Server, MessageSquare, Trash2, Settings } from 'lucide-react'

interface CanvasEditorProps {
  config: any
  onChange: (config: any) => void
  readOnly?: boolean
}

interface ServiceNode {
  id: string
  name: string
  type: 'redis' | 'postgres' | 'kafka'
  host: string
  port: number
  username?: string
  password?: string
  database?: string
}

const SERVICE_TYPES = [
  { type: 'redis', label: 'Redis', icon: Database, color: '#ef4444', bgColor: '#fee2e2', darkBgColor: '#450a0a40', defaultPort: 6379 },
  { type: 'postgres', label: 'PostgreSQL', icon: Server, color: '#3b82f6', bgColor: '#dbeafe', darkBgColor: '#1e3a8a40', defaultPort: 5433 },
  { type: 'kafka', label: 'Kafka', icon: MessageSquare, color: '#a855f7', bgColor: '#f3e8ff', darkBgColor: '#581c8740', defaultPort: 9092 },
]

export default function CanvasEditor({ config, onChange, readOnly = false }: CanvasEditorProps) {
  const [services, setServices] = useState<ServiceNode[]>(() =>
    Object.entries(config.services || {}).map(([name, svc]: [string, any]) => ({
      id: name,
      name,
      type: svc.type,
      host: svc.host,
      port: svc.port,
      username: svc.username,
      password: svc.password,
      database: svc.database,
    }))
  )
  const [selectedService, setSelectedService] = useState<ServiceNode | null>(null)
  const [showAddMenu, setShowAddMenu] = useState(false)

  const addService = (type: 'redis' | 'postgres' | 'kafka') => {
    const serviceType = SERVICE_TYPES.find((s) => s.type === type)!
    const existingOfType = services.filter((s) => s.type === type)
    const count = existingOfType.length
    
    // Auto-increment port for multiple instances of same type
    const basePort = serviceType.defaultPort
    let port = basePort
    if (count > 0) {
      // Find the highest port number for this service type and add 1
      const maxPort = Math.max(...existingOfType.map(s => s.port))
      port = maxPort >= basePort ? maxPort + 1 : basePort + count
    }
    
    const newService: ServiceNode = {
      id: `${type}-${Date.now()}`,
      name: `${type}-${count + 1}`,
      type,
      host: 'localhost',
      port,
      username: type === 'postgres' ? 'postgres' : undefined,
      password: type === 'postgres' ? 'password' : undefined,
      database: type === 'postgres' ? `${type}_db_${count + 1}` : undefined,
    }

    const updated = [...services, newService]
    setServices(updated)
    updateConfig(updated)
    setSelectedService(newService)
    setShowAddMenu(false)
  }

  const updateService = (id: string, updates: Partial<ServiceNode>) => {
    const updated = services.map((s) => (s.id === id ? { ...s, ...updates } : s))
    setServices(updated)
    updateConfig(updated)
    if (selectedService?.id === id) {
      setSelectedService({ ...selectedService, ...updates })
    }
  }

  const deleteService = (id: string) => {
    const updated = services.filter((s) => s.id !== id)
    setServices(updated)
    updateConfig(updated)
    if (selectedService?.id === id) {
      setSelectedService(null)
    }
  }

  const updateConfig = (servicesList: ServiceNode[]) => {
    const servicesConfig = servicesList.reduce((acc, service) => {
      acc[service.name] = {
        type: service.type,
        host: service.host,
        port: service.port,
        ...(service.username && { username: service.username }),
        ...(service.password && { password: service.password }),
        ...(service.database && { database: service.database }),
      }
      return acc
    }, {} as Record<string, any>)

    onChange({ ...config, services: servicesConfig })
  }

  return (
    <div className="grid grid-cols-3 gap-0 h-full">
      {/* Canvas Area */}
      <div className="col-span-2 border-r border-border p-8 bg-background relative overflow-y-auto">
        <div className="space-y-4 mb-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-semibold text-foreground">
                Infrastructure Canvas
              </h3>
              <p className="text-sm text-muted-foreground mt-1">
                Drag and configure your services
              </p>
            </div>
            {!readOnly && (
              <button
                onClick={() => setShowAddMenu(!showAddMenu)}
                className="flex items-center space-x-2 px-4 py-2 bg-[#FF5050] text-white rounded-lg hover:bg-[#ed1515] transition-colors shadow-sm"
              >
                <Plus className="w-5 h-5" />
                <span>Add Service</span>
              </button>
            )}
          </div>
          
          {!readOnly && services.length === 0 && (
            <div className="p-4 bg-primary/5 border border-primary rounded-lg">
              <p className="text-sm text-primary">
                💡 <strong>Tip:</strong> You can add multiple instances of the same service type (e.g., 2 PostgreSQL databases, 3 Redis caches). 
                Ports will auto-increment to avoid conflicts.
              </p>
            </div>
          )}
        </div>

        {/* Add Service Menu */}
        {showAddMenu && (
          <div className="absolute top-20 right-8 bg-card border border-border rounded-lg shadow-xl p-2 z-10 min-w-[280px]">
            <div className="px-4 py-2 border-b border-border">
              <p className="text-xs text-muted-foreground">
                Add multiple instances of the same type
              </p>
            </div>
            {SERVICE_TYPES.map((service) => {
              const count = services.filter((s) => s.type === service.type).length
              return (
                <button
                  key={service.type}
                  onClick={() => addService(service.type as any)}
                  className="flex items-center justify-between w-full px-4 py-3 text-left hover:bg-accent rounded-lg transition-colors"
                >
                  <div className="flex items-center space-x-3">
                    <div 
                      className="w-8 h-8 rounded-lg flex items-center justify-center"
                      style={{ 
                        backgroundColor: document.documentElement.classList.contains('dark') ? service.darkBgColor : service.bgColor 
                      }}
                    >
                      <service.icon className="w-4 h-4" style={{ color: service.color }} />
                    </div>
                    <span className="text-foreground font-medium">{service.label}</span>
                  </div>
                  {count > 0 && (
                    <span className="px-2 py-1 text-xs font-medium rounded-full bg-muted text-foreground">
                      {count} added
                    </span>
                  )}
                </button>
              )
            })}
          </div>
        )}

        {/* Service Nodes */}
        <div className="grid grid-cols-2 xl:grid-cols-3 gap-6">
          {services.map((service) => {
            const serviceType = SERVICE_TYPES.find((s) => s.type === service.type)!
            const Icon = serviceType.icon
            const isSelected = selectedService?.id === service.id

            return (
              <div
                key={service.id}
                onClick={() => setSelectedService(service)}
                className={`p-5 bg-card border-2 rounded-xl cursor-pointer transition-all ${
                  isSelected
                    ? 'border-primary shadow-xl ring-4 ring-primary/20'
                    : 'border-border hover:border-secondary hover:shadow-lg'
                }`}
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center space-x-3">
                    <div
                      className="w-10 h-10 rounded-lg flex items-center justify-center"
                      style={{ 
                        backgroundColor: document.documentElement.classList.contains('dark') ? serviceType.darkBgColor : serviceType.bgColor 
                      }}
                    >
                      <Icon className="w-5 h-5" style={{ color: serviceType.color }} />
                    </div>
                    <div>
                      <div className="font-medium text-foreground">
                        {service.name}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {serviceType.label}
                      </div>
                    </div>
                  </div>
                  {!readOnly && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        deleteService(service.id)
                      }}
                      className="p-1 text-muted-foreground hover:text-error transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  )}
                </div>

                <div className="space-y-1 text-xs text-muted-foreground">
                  <div className="flex items-center justify-between">
                    <span>Host:</span>
                    <span className="font-mono">{service.host}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span>Port:</span>
                    <span className="font-mono">{service.port}</span>
                  </div>
                  {service.database && (
                    <div className="flex items-center justify-between">
                      <span>Database:</span>
                      <span className="font-mono">{service.database}</span>
                    </div>
                  )}
                </div>
              </div>
            )
          })}

          {services.length === 0 && (
            <div className="col-span-2 xl:col-span-3 flex flex-col items-center justify-center py-24 text-muted-foreground">
              <div className="w-24 h-24 bg-muted rounded-full flex items-center justify-center mb-6">
                <Database className="w-12 h-12 opacity-50" />
              </div>
              <h3 className="text-lg font-medium text-muted-foreground mb-2">
                No services yet
              </h3>
              <p className="text-sm text-center mb-6 max-w-md">
                Start building your cluster by adding Redis, PostgreSQL, or Kafka services
              </p>
              <button
                onClick={() => setShowAddMenu(true)}
                className="flex items-center space-x-2 px-6 py-3 bg-[#FF5050] text-white rounded-lg hover:bg-[#ed1515] transition-colors shadow-sm"
              >
                <Plus className="w-5 h-5" />
                <span>Add Your First Service</span>
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Properties Panel */}
      <div className="p-8 bg-secondary overflow-y-auto">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center space-x-2">
            <Settings className="w-6 h-6 text-muted-foreground" />
            <h3 className="text-lg font-semibold text-foreground">Properties</h3>
          </div>
          {selectedService && !readOnly && (
            <button
              onClick={() => {
                const existingOfType = services.filter((s) => s.type === selectedService.type)
                const count = existingOfType.length
                const maxPort = Math.max(...existingOfType.map(s => s.port))
                
                const duplicated: ServiceNode = {
                  id: `${selectedService.type}-${Date.now()}`,
                  name: `${selectedService.type}-${count + 1}`,
                  type: selectedService.type,
                  host: selectedService.host,
                  port: maxPort + 1,
                  username: selectedService.username,
                  password: selectedService.password,
                  database: selectedService.database ? `${selectedService.type}_db_${count + 1}` : undefined,
                }
                const updated = [...services, duplicated]
                setServices(updated)
                updateConfig(updated)
                setSelectedService(duplicated)
              }}
              className="px-3 py-1.5 text-xs bg-[#FF5050] text-white rounded-lg hover:bg-[#ed1515] transition-colors"
            >
              Duplicate
            </button>
          )}
        </div>

        {selectedService ? (
          <div className="space-y-5">
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">
                Service Name
              </label>
              <input
                type="text"
                value={selectedService.name}
                onChange={(e) => !readOnly && !readOnly && updateService(selectedService.id, { name: e.target.value })}
                readOnly={readOnly}
                className={`w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-[#FF5050] bg-input text-foreground text-sm ${readOnly ? 'cursor-not-allowed opacity-75' : ''}`}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-foreground mb-1">
                Host
              </label>
              <input
                type="text"
                value={selectedService.host}
                onChange={(e) => !readOnly && !readOnly && updateService(selectedService.id, { host: e.target.value })}
                className={`w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-[#FF5050] bg-input text-foreground text-sm ${readOnly ? "cursor-not-allowed opacity-75" : ""}`}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-foreground mb-1">
                Port
              </label>
              <input
                type="number"
                value={selectedService.port}
                onChange={(e) => !readOnly &&
                  updateService(selectedService.id, { port: parseInt(e.target.value) })
                }
                className={`w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-[#FF5050] bg-input text-foreground text-sm ${readOnly ? "cursor-not-allowed opacity-75" : ""}`}
              />
              {!readOnly && services.filter(s => s.id !== selectedService.id && s.port === selectedService.port).length > 0 && (
                <p className="text-xs text-error mt-1">
                  ⚠️ Port conflict detected! Another service uses this port.
                </p>
              )}
            </div>

            {selectedService.type === 'postgres' && (
              <>
                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Username
                  </label>
                  <input
                    type="text"
                    value={selectedService.username || ''}
                    onChange={(e) => !readOnly &&
                      updateService(selectedService.id, { username: e.target.value })
                    }
                    className={`w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-[#FF5050] bg-input text-foreground text-sm ${readOnly ? "cursor-not-allowed opacity-75" : ""}`}
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Password
                  </label>
                  <input
                    type="password"
                    value={selectedService.password || ''}
                    onChange={(e) => !readOnly &&
                      updateService(selectedService.id, { password: e.target.value })
                    }
                    className={`w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-[#FF5050] bg-input text-foreground text-sm ${readOnly ? "cursor-not-allowed opacity-75" : ""}`}
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-foreground mb-1">
                    Database
                  </label>
                  <input
                    type="text"
                    value={selectedService.database || ''}
                    onChange={(e) => !readOnly &&
                      updateService(selectedService.id, { database: e.target.value })
                    }
                    className={`w-full px-3 py-2 border border-border rounded-lg focus:ring-2 focus:ring-[#FF5050] bg-input text-foreground text-sm ${readOnly ? "cursor-not-allowed opacity-75" : ""}`}
                  />
                </div>
              </>
            )}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
            <Settings className="w-12 h-12 mb-3 opacity-50" />
            <p className="text-sm text-center">Select a service to edit its properties</p>
          </div>
        )}
      </div>
    </div>
  )
}

