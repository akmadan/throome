import { useState, useEffect } from 'react'
import { Copy, Check } from 'lucide-react'

interface YamlEditorProps {
  config: any
  onChange: (config: any) => void
}

export default function YamlEditor({ config, onChange }: YamlEditorProps) {
  const [yamlText, setYamlText] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  // Convert config to YAML string
  useEffect(() => {
    const yaml = configToYaml(config)
    setYamlText(yaml)
  }, [config])

  const handleYamlChange = (text: string) => {
    setYamlText(text)
    
    try {
      const parsed = yamlToConfig(text)
      onChange(parsed)
      setError(null)
    } catch (err: any) {
      setError(err.message)
    }
  }

  const handleCopy = () => {
    navigator.clipboard.writeText(yamlText)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const exampleYaml = `# Example cluster configuration
services:
  redis-cache:
    type: redis
    host: localhost
    port: 6379

  main-db:
    type: postgres
    host: localhost
    port: 5432
    username: postgres
    password: password
    database: myapp

  events-queue:
    type: kafka
    host: localhost
    port: 9092`

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-600 dark:text-gray-400">
          Edit your cluster configuration in YAML format
        </div>
        <button
          onClick={handleCopy}
          className="flex items-center space-x-2 px-3 py-1.5 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
        >
          {copied ? (
            <>
              <Check className="w-4 h-4 text-green-500" />
              <span>Copied!</span>
            </>
          ) : (
            <>
              <Copy className="w-4 h-4" />
              <span>Copy</span>
            </>
          )}
        </button>
      </div>

      {/* Editor */}
      <div className="relative">
        <textarea
          value={yamlText}
          onChange={(e) => handleYamlChange(e.target.value)}
          placeholder={exampleYaml}
          className={`w-full h-96 px-4 py-3 font-mono text-sm border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-white ${
            error
              ? 'border-red-500 dark:border-red-500'
              : 'border-gray-300 dark:border-gray-600'
          }`}
          spellCheck={false}
        />
        
        {/* Line numbers overlay (simple version) */}
        <div className="absolute top-0 left-0 px-2 py-3 text-xs text-gray-400 pointer-events-none select-none font-mono">
          {yamlText.split('\n').map((_, i) => (
            <div key={i} className="h-5">
              {i + 1}
            </div>
          ))}
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-sm text-red-600 dark:text-red-400">
            <strong>YAML Parse Error:</strong> {error}
          </p>
        </div>
      )}

      {/* Help */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <h4 className="text-sm font-medium text-blue-900 dark:text-blue-300 mb-2">
          Configuration Format
        </h4>
        <ul className="text-xs text-blue-800 dark:text-blue-400 space-y-1">
          <li>• Define services under the <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">services:</code> key</li>
          <li>• Each service needs: <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">type</code>, <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">host</code>, <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">port</code></li>
          <li>• Supported types: <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">redis</code>, <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">postgres</code>, <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">kafka</code></li>
          <li>• PostgreSQL services can include: <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">username</code>, <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">password</code>, <code className="bg-blue-100 dark:bg-blue-900 px-1 py-0.5 rounded">database</code></li>
        </ul>
      </div>
    </div>
  )
}

// Simple YAML to config converter
function yamlToConfig(yaml: string): any {
  const lines = yaml.split('\n').filter(line => line.trim() && !line.trim().startsWith('#'))
  const config: any = { services: {} }
  let currentService: string | null = null
  let currentServiceData: any = {}

  for (const line of lines) {
    const trimmed = line.trim()
    const indent = line.search(/\S/)

    if (indent === 0 && trimmed === 'services:') {
      continue
    }

    if (indent === 2 && trimmed.endsWith(':')) {
      // Save previous service
      if (currentService) {
        config.services[currentService] = currentServiceData
      }
      // Start new service
      currentService = trimmed.slice(0, -1)
      currentServiceData = {}
    } else if (indent === 4 && currentService) {
      const [key, ...valueParts] = trimmed.split(':')
      const value = valueParts.join(':').trim()
      if (key && value) {
        currentServiceData[key] = isNaN(Number(value)) ? value : Number(value)
      }
    }
  }

  // Save last service
  if (currentService) {
    config.services[currentService] = currentServiceData
  }

  return config
}

// Simple config to YAML converter
function configToYaml(config: any): string {
  if (!config.services || Object.keys(config.services).length === 0) {
    return ''
  }

  let yaml = 'services:\n'
  
  for (const [name, service] of Object.entries(config.services as Record<string, any>)) {
    yaml += `  ${name}:\n`
    yaml += `    type: ${service.type}\n`
    yaml += `    host: ${service.host}\n`
    yaml += `    port: ${service.port}\n`
    if (service.username) yaml += `    username: ${service.username}\n`
    if (service.password) yaml += `    password: ${service.password}\n`
    if (service.database) yaml += `    database: ${service.database}\n`
  }

  return yaml
}

