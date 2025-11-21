# Throome Dashboard UI

Modern React-based dashboard for managing Throome Gateway.

## Features

- 🎨 Beautiful, responsive UI built with React + TypeScript
- 🌓 Dark mode support
- 📊 Real-time monitoring and metrics
- 🔧 Complete cluster management (CRUD operations)
- 🚀 Fast development with Vite
- 💅 Styled with Tailwind CSS

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn/pnpm
- Throome Gateway running on port 9000

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The UI will be available at `http://localhost:3000`

### Build for Production

```bash
# Build optimized bundle
npm run build

# Preview production build
npm run preview
```

## Project Structure

```
ui/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── Layout.tsx    # Main layout wrapper
│   │   ├── Sidebar.tsx   # Navigation sidebar
│   │   ├── Header.tsx    # Top header bar
│   │   └── StatsCard.tsx # Statistics card component
│   ├── pages/           # Page components
│   │   ├── Dashboard.tsx # Main dashboard
│   │   ├── Clusters.tsx  # Cluster management
│   │   ├── Services.tsx  # Service listing
│   │   ├── Monitoring.tsx # Monitoring view
│   │   ├── Routing.tsx   # Routing configuration
│   │   └── Settings.tsx  # Settings panel
│   ├── lib/             # Utilities and helpers
│   ├── App.tsx          # Main app component
│   ├── main.tsx         # Entry point
│   └── index.css        # Global styles
├── public/              # Static assets
└── package.json
```

## Sidebar Navigation

The sidebar includes the following sections:

### 📊 Dashboard
- Overview of all clusters and services
- Quick stats and metrics
- Recent activity feed
- Quick action buttons

### 📦 Clusters
- List all clusters
- Create new clusters
- Edit cluster configurations
- Delete clusters
- View cluster details and services

### 🗄️ Services
- View all services across clusters
- Service health status
- Service-specific metrics
- Filter by service type (Redis, PostgreSQL, Kafka)

### 📈 Monitoring
- Real-time metrics and charts
- System health status
- Performance graphs
- Alert history
- Custom dashboards

### 🔀 Routing
- Configure routing strategies
- View active routes
- Load balancing configuration
- AI-based routing settings

### ⚙️ Settings
- Gateway configuration
- User preferences
- API keys management
- Notification settings
- Theme preferences

## API Integration

The UI communicates with Throome Gateway via API proxy configured in `vite.config.ts`:

```typescript
proxy: {
  '/api': {
    target: 'http://localhost:9000',
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api/, ''),
  },
}
```

Example API usage:

```typescript
import axios from 'axios'

// Get all clusters
const clusters = await axios.get('/api/clusters')

// Create cluster
const newCluster = await axios.post('/api/clusters', {
  name: 'my-cluster',
  config: { /* ... */ }
})
```

## Customization

### Colors and Theme

Edit `tailwind.config.js` and `src/index.css` to customize colors and theme.

### Add New Pages

1. Create new component in `src/pages/`
2. Add route in `src/App.tsx`
3. Add navigation item in `src/components/Sidebar.tsx`

## Development

```bash
# Run linter
npm run lint

# Format code
npm run format

# Type checking
npm run build
```

## Deployment

The built files can be served by any static file server or integrated with the Throome Gateway binary.

### Serve with Throome

The UI can be embedded in the Go binary using `go:embed` or served from the `ui/dist` directory.

## License

Apache 2.0 - Same as Throome Gateway
