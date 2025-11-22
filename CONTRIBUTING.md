# Contributing to Throome

Thank you for your interest in contributing to Throome! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [UI Development Workflow](#ui-development-workflow)
- [Backend Development Workflow](#backend-development-workflow)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Code Style Guidelines](#code-style-guidelines)

---

## Development Environment Setup

### Prerequisites

- Go 1.24+
- Node.js 18+ and npm
- Docker Engine 20.10+
- Make
- Git

### Initial Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/throome.git
cd throome

# Add upstream remote
git remote add upstream https://github.com/akmadan/throome.git

# Install Go dependencies
go mod download

# Install UI dependencies
cd ui
npm install
cd ..
```

---

## UI Development Workflow

The UI is built with React + TypeScript + Vite and **embedded** into the Go binary. This requires a specific workflow for local testing and deployment.

### Understanding the UI Embedding

Throome uses Go's `embed` directive to include the built UI directly in the binary:

```go
// pkg/gateway/ui.go
//go:embed ui/dist
var UIFS embed.FS
```

This means the UI must be built and copied to `pkg/gateway/ui/dist/` before building the Go binary.

### Local UI Development (Hot Reload)

For rapid UI development with hot reload:

```bash
# Terminal 1: Run Throome backend
cd throome
make run

# Terminal 2: Run UI dev server with proxy
cd ui
npm run dev
```

The Vite dev server (port 3000) proxies API requests to the Throome backend (port 9000).

**Access**: http://localhost:3000

**Note**: This setup uses live UI files, not embedded ones. API calls are automatically proxied to the backend.

### Building UI for Embedding

When your UI changes are ready, build and embed them:

```bash
# Build UI
cd ui
npm run build
cd ..

# Copy built assets to Go embedding location
rm -rf pkg/gateway/ui/dist
mkdir -p pkg/gateway/ui
cp -r ui/dist pkg/gateway/ui/

# Verify files are copied
ls pkg/gateway/ui/dist/
```

**Important**: The `pkg/gateway/ui/dist/` directory must contain:
- `index.html`
- `assets/` directory with CSS and JS bundles
- Static assets (logos, icons)

### Testing UI Changes Locally

#### Option 1: Via Go Binary (Tests Embedding)

```bash
# 1. Build UI
cd ui && npm run build && cd ..

# 2. Copy to embedding location
rm -rf pkg/gateway/ui/dist
mkdir -p pkg/gateway/ui
cp -r ui/dist pkg/gateway/ui/

# 3. Build and run Go binary
make build
./bin/throome --port 9000

# 4. Test at http://localhost:9000
```

#### Option 2: Via Docker (Production-like)

```bash
# 1. Build UI
cd ui && npm run build && cd ..

# 2. Copy to embedding location
rm -rf pkg/gateway/ui/dist
mkdir -p pkg/gateway/ui
cp -r ui/dist pkg/gateway/ui/

# 3. Build Docker image
docker build -f deployments/docker/Dockerfile -t throome:test .

# 4. Run container
docker run --name throome-test \
  --user root \
  --add-host=host.docker.internal:host-gateway \
  -p 9000:9000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd)/clusters:/app/clusters \
  --rm \
  throome:test

# 5. Test at http://localhost:9000
```

### Common UI Pitfalls

1. **Forgot to copy UI assets**: Go binary serves old UI or crashes
   - **Solution**: Always run `cp -r ui/dist pkg/gateway/ui/` after building

2. **API calls failing in dev mode**: Proxy not configured
   - **Solution**: Check `ui/vite.config.ts` has correct proxy settings

3. **UI changes not showing**: Using old embedded files
   - **Solution**: Rebuild UI and copy to `pkg/gateway/ui/dist/` before rebuilding Go binary

4. **Docker build fails**: UI assets missing
   - **Solution**: Ensure `ui/dist` exists and is not in `.dockerignore`

### UI File Structure

```
ui/
├── src/
│   ├── api/
│   │   └── client.ts          # API client (axios)
│   ├── components/
│   │   ├── CanvasEditor.tsx   # Drag-and-drop service editor
│   │   ├── YamlEditor.tsx     # YAML configuration editor
│   │   ├── Layout.tsx         # Main layout wrapper
│   │   ├── Sidebar.tsx        # Navigation sidebar
│   │   └── Header.tsx         # Top header bar
│   ├── pages/
│   │   ├── Dashboard.tsx      # Main dashboard
│   │   ├── Clusters.tsx       # Cluster list
│   │   ├── CreateCluster.tsx  # Cluster creation page
│   │   ├── ViewCluster.tsx    # Cluster details page
│   │   └── Services.tsx       # Services overview
│   ├── lib/
│   │   └── utils.ts           # Utility functions
│   └── index.css              # Global styles (Tailwind)
├── public/                     # Static assets
├── dist/                       # Build output (gitignored)
└── package.json
```

### Making UI Changes: Complete Workflow

```bash
# 1. Create feature branch
git checkout -b feature/ui-improvement

# 2. Make changes in ui/src/

# 3. Test with hot reload
cd ui && npm run dev

# 4. Build UI
npm run build
cd ..

# 5. Copy to embedding location
rm -rf pkg/gateway/ui/dist
mkdir -p pkg/gateway/ui
cp -r ui/dist pkg/gateway/ui/

# 6. Test with Go binary
make build && ./bin/throome

# 7. Run validation script (optional)
./scripts/test-workflow.sh

# 8. Commit changes (include built files)
git add ui/ pkg/gateway/ui/dist/
git commit -m "feat(ui): your improvement description"

# 9. Push and create PR
git push origin feature/ui-improvement
```

**Important**: Do **NOT** gitignore `pkg/gateway/ui/dist/` as these files must be committed for the Go embed to work.

---

## Backend Development Workflow

### Project Layout

```
pkg/                    # Public packages
├── adapters/          # Service-specific adapters
│   ├── redis/
│   ├── postgres/
│   └── kafka/
├── cluster/           # Cluster management
├── gateway/           # HTTP server and core logic
├── provisioner/       # Docker container management
├── router/            # Routing strategies
├── monitor/           # Health checks and metrics
└── sdk/               # Go SDK for clients

internal/              # Private packages
├── config/            # Configuration structures
├── logger/            # Logging utilities
└── utils/             # Common utilities

cmd/                   # Application entrypoints
├── throome/           # Gateway service
└── throome-cli/       # CLI tool
```

### Making Backend Changes

```bash
# 1. Create feature branch
git checkout -b feature/backend-improvement

# 2. Make changes in pkg/ or internal/

# 3. Add/update tests
# Test files go next to the code: pkg/cluster/manager_test.go

# 4. Run tests
make test-unit

# 5. Run integration tests
make test-integration

# 6. Build and test locally
make build
./bin/throome

# 7. Commit changes
git add .
git commit -m "feat(gateway): your improvement description"

# 8. Push and create PR
git push origin feature/backend-improvement
```

### Adding a New Service Adapter

To add support for a new service type (e.g., MongoDB):

1. Create adapter package:
```
pkg/adapters/mongodb/
├── mongodb.go        # Adapter implementation
└── mongodb_test.go   # Unit tests
```

2. Implement `adapters.Adapter` interface:
```go
type Adapter struct {
    config *cluster.ServiceConfig
    client *mongo.Client
}

func (a *Adapter) Connect(ctx context.Context) error
func (a *Adapter) Disconnect(ctx context.Context) error
func (a *Adapter) HealthCheck(ctx context.Context) (*adapters.HealthStatus, error)
func (a *Adapter) GetMetrics() (*adapters.Metrics, error)
```

3. Register in factory (`pkg/adapters/adapter.go`):
```go
func (f *Factory) Create(config *cluster.ServiceConfig) (Adapter, error) {
    case "mongodb":
        return mongodb.NewAdapter(config)
```

4. Add provisioning support (`pkg/provisioner/docker.go`):
```go
case "mongodb":
    imageName = "mongo:7"
    // Configure environment and ports
```

5. Update UI (`ui/src/components/CanvasEditor.tsx`):
```typescript
const SERVICE_TYPES = [
  // ...
  { type: 'mongodb', label: 'MongoDB', icon: Database, color: 'green', defaultPort: 27017 },
]
```

---

## Testing

### Unit Tests

```bash
# Run all unit tests
make test-unit

# Run specific package tests
go test ./pkg/cluster/... -v

# Run with coverage
make test-coverage

# View coverage report
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Start test infrastructure (Redis, PostgreSQL, Kafka)
make test-setup

# Run integration tests
make test-integration

# Stop test infrastructure
make test-teardown
```

### Writing Tests

**Unit Test Example** (`pkg/cluster/manager_test.go`):

```go
func TestManager_Create(t *testing.T) {
    tmpDir := t.TempDir()
    manager := cluster.NewManager(tmpDir)
    
    config := &cluster.Config{
        Services: map[string]cluster.ServiceConfig{
            "redis-1": {
                Type: "redis",
                Host: "localhost",
                Port: 6379,
            },
        },
    }
    
    clusterID, err := manager.Create("test-cluster", config)
    assert.NoError(t, err)
    assert.NotEmpty(t, clusterID)
}
```

**Integration Test Example** (`test/integration/adapters_test.go`):

```go
func TestRedisAdapter_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    config := &cluster.ServiceConfig{
        Type: "redis",
        Host: "localhost",
        Port: 6379,
    }
    
    adapter := redis.NewAdapter(config)
    err := adapter.Connect(context.Background())
    require.NoError(t, err)
    defer adapter.Disconnect(context.Background())
    
    // Test operations
    err = adapter.Set(context.Background(), "key", "value")
    assert.NoError(t, err)
}
```

### Test Coverage Goals

- Unit tests: 80%+ coverage
- Integration tests: Critical paths covered
- All new features: Tests required

---

## Pull Request Process

### Before Submitting

1. **Run validation script**:
   ```bash
   ./scripts/test-workflow.sh
   ```
   This script runs all build and test steps.

2. **Ensure CI passes**:
   - All tests pass
   - No linting errors
   - Docker image builds successfully

3. **Update documentation**:
   - Update README if adding features
   - Add inline code comments
   - Update API documentation if relevant

### PR Guidelines

**Title Format**: `<type>(<scope>): <description>`

Examples:
- `feat(ui): add cluster duplication feature`
- `fix(gateway): resolve connection timeout issue`
- `docs(readme): update installation instructions`
- `test(adapters): add integration tests for PostgreSQL`

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Test additions or fixes
- `chore`: Build/tooling changes
- `perf`: Performance improvements

**PR Description Template**:

```markdown
## Description
Brief description of changes

## Motivation
Why is this change needed?

## Changes
- Change 1
- Change 2

## Testing
How was this tested?

## Screenshots (if UI changes)
Before/After screenshots

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] UI assets rebuilt (if UI changes)
- [ ] All tests passing
```

### Review Process

1. Automated checks run (tests, lint, build)
2. Maintainer review
3. Address feedback
4. Approval and merge

---

## Code Style Guidelines

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Maximum line length: 120 characters
- Use meaningful variable names
- Add comments for exported functions

**Example**:

```go
// CreateCluster provisions services and creates a new cluster configuration.
// It returns the cluster ID or an error if provisioning fails.
func (g *Gateway) CreateCluster(ctx context.Context, name string, config *cluster.Config) (string, error) {
    // Validate input
    if name == "" {
        return "", fmt.Errorf("cluster name cannot be empty")
    }
    
    // Implementation...
}
```

### TypeScript/React Code Style

- Use functional components with hooks
- TypeScript strict mode enabled
- Use Tailwind CSS for styling (no inline styles)
- Component file = PascalCase.tsx
- Utility file = camelCase.ts
- Export interfaces for props

**Example**:

```typescript
interface CanvasEditorProps {
  config: Cluster
  onChange: (newConfig: Cluster) => void
  readOnly?: boolean
}

export default function CanvasEditor({ config, onChange, readOnly = false }: CanvasEditorProps) {
  const [services, setServices] = useState<ServiceNode[]>([])
  
  // Implementation...
}
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Example**:

```
feat(gateway): add support for connection pooling

Implement per-adapter connection pools with configurable
min/max connections and idle timeout.

Closes #123
```

---

## Getting Help

- **Questions**: Open a [Discussion](https://github.com/akmadan/throome/discussions)
- **Bugs**: Open an [Issue](https://github.com/akmadan/throome/issues)
- **Chat**: Join our community (coming soon)

---

## License

By contributing to Throome, you agree that your contributions will be licensed under the Apache License 2.0.

Thank you for contributing!
