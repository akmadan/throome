# Testing Guide for Throome

This guide covers how to test the Throome gateway effectively.

## 📋 Table of Contents

- [Testing Levels](#testing-levels)
- [Quick Start](#quick-start)
- [Unit Tests](#unit-tests)
- [Integration Tests](#integration-tests)
- [End-to-End Tests](#end-to-end-tests)
- [Benchmarks](#benchmarks)
- [Coverage](#coverage)
- [CI/CD Testing](#cicd-testing)

---

## 🎯 Testing Levels

### 1. **Unit Tests**
Test individual components in isolation.
- **Location**: `pkg/*/`
- **Pattern**: `*_test.go` files next to source
- **No external dependencies** (use mocks)

### 2. **Integration Tests**
Test components with real infrastructure services.
- **Location**: `test/integration/`
- **Requires**: Docker services (Redis, PostgreSQL, Kafka)
- **Tests**: Adapters, connections, data flow

### 3. **End-to-End Tests**
Test complete workflows through the gateway.
- **Location**: `test/e2e/`
- **Tests**: Full request cycles, SDK usage, multi-service operations

---

## 🚀 Quick Start

### Prerequisites

```bash
# Install dependencies
go mod download

# Install Docker (for integration tests)
# macOS: brew install docker
# Linux: apt-get install docker.io
```

### Run All Tests

```bash
# Unit tests only (fast, no dependencies)
make test-unit

# All tests including integration (requires Docker)
make test

# With coverage
make test-coverage
```

---

## 🧪 Unit Tests

Unit tests don't require external services.

### Run Unit Tests

```bash
# All unit tests
go test -short ./...

# Specific package
go test -short ./pkg/cluster

# Verbose output
go test -v -short ./pkg/cluster

# With race detection
go test -race -short ./...
```

### Writing Unit Tests

**Example: `pkg/cluster/config_test.go`**

```go
func TestConfigValidate(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: &Config{
                ClusterID: "test-01",
                Name:      "Test",
                Services: map[string]ServiceConfig{
                    "cache": {
                        Type: "redis",
                        Host: "localhost",
                        Port: 6379,
                    },
                },
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Current Unit Test Coverage

- ✅ `pkg/cluster/config_test.go` - Configuration validation
- ✅ `pkg/cluster/manager_test.go` - Cluster lifecycle
- ✅ `internal/utils/validation_test.go` - Input validation
- 🚧 More tests to be added...

---

## 🔗 Integration Tests

Integration tests require real services running.

### Setup Test Infrastructure

```bash
# Start all test services
cd test
docker-compose up -d

# Check services are running
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Run Integration Tests

```bash
# Set environment variable to enable integration tests
export INTEGRATION_TESTS=true

# Run integration tests
make test-integration

# Or directly
go test -v ./test/integration/...

# Stop services when done
cd test && docker-compose down
```

### Test Services

The `docker-compose.yml` provides:

| Service    | Port | Credentials           |
|------------|------|-----------------------|
| Redis      | 6379 | (no password)         |
| PostgreSQL | 5432 | user: test, pass: test|
| Kafka      | 9092 | (no auth)             |
| Zookeeper  | 2181 | (kafka dependency)    |

### Writing Integration Tests

**Example: Test Redis Adapter**

```go
// test/integration/redis_test.go
func TestRedisAdapter_Integration(t *testing.T) {
    if os.Getenv("INTEGRATION_TESTS") != "true" {
        t.Skip("Skipping integration test")
    }

    config := cluster.ServiceConfig{
        Type: "redis",
        Host: "localhost",
        Port: 6379,
    }

    adapter, err := redis.NewRedisAdapter(config)
    if err != nil {
        t.Fatalf("Failed to create adapter: %v", err)
    }

    ctx := context.Background()

    // Test connection
    if err := adapter.Connect(ctx); err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer adapter.Disconnect(ctx)

    // Test operations
    cacheAdapter := adapter.(adapters.CacheAdapter)
    
    err = cacheAdapter.Set(ctx, "test-key", "test-value", 0)
    if err != nil {
        t.Errorf("Set failed: %v", err)
    }

    value, err := cacheAdapter.Get(ctx, "test-key")
    if err != nil {
        t.Errorf("Get failed: %v", err)
    }

    if value != "test-value" {
        t.Errorf("Expected test-value, got %s", value)
    }
}
```

---

## 🌐 End-to-End Tests

E2E tests test complete workflows.

### Run E2E Tests

```bash
# Start gateway server
./bin/throome &
GATEWAY_PID=$!

# Start test services
cd test && docker-compose up -d

# Run E2E tests
go test -v ./test/e2e/...

# Cleanup
kill $GATEWAY_PID
cd test && docker-compose down
```

### E2E Test Scenarios

1. **Full Gateway Lifecycle**
   - Start gateway
   - Create cluster via CLI
   - Connect via SDK
   - Perform operations
   - Check metrics
   - Cleanup

2. **Multi-Service Operation**
   - Use cache, database, and queue together
   - Verify data consistency
   - Test failover

3. **Load Testing**
   - Concurrent requests
   - Connection pooling
   - Memory usage

---

## ⚡ Benchmarks

Measure performance characteristics.

### Run Benchmarks

```bash
# All benchmarks
go test -bench=. ./...

# Specific package
go test -bench=. ./pkg/router

# With memory stats
go test -bench=. -benchmem ./pkg/router

# Save results
go test -bench=. ./... > benchmark.txt
```

### Writing Benchmarks

```go
func BenchmarkRoundRobinStrategy(b *testing.B) {
    strategy := NewRoundRobinStrategy()
    adapters := make([]adapters.Adapter, 10)
    // ... setup adapters ...

    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := strategy.Select(ctx, adapters)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## 📊 Coverage

Track test coverage.

### Generate Coverage Report

```bash
# Generate coverage
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

### Coverage Goals

- **Overall**: > 80%
- **Critical paths**: > 90%
  - Cluster management
  - Adapter connections
  - Routing logic
- **Utilities**: > 95%

### Check Coverage by Package

```bash
# Coverage per package
go test -cover ./...

# Detailed coverage
go test -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out | grep -E "^total:"
```

---

## 🔄 CI/CD Testing

Automate testing in GitHub Actions.

### GitHub Actions Workflow

Create `.github/workflows/test.yml`:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -short -race -coverprofile=coverage.txt ./...
      - uses: codecov/codecov-action@v3
        with:
          file: ./coverage.txt

  integration-tests:
    runs-on: ubuntu-latest
    services:
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: throome_test
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: INTEGRATION_TESTS=true go test -v ./test/integration/...
```

---

## 🛠️ Useful Testing Commands

```bash
# Run tests with timeout
go test -timeout 30s ./...

# Run specific test
go test -run TestConfigValidate ./pkg/cluster

# Run tests matching pattern
go test -run "^TestManager" ./pkg/cluster

# List all tests
go test -list . ./...

# Run tests with CPU profiling
go test -cpuprofile cpu.prof ./pkg/router
go tool pprof cpu.prof

# Run tests with memory profiling
go test -memprofile mem.prof ./pkg/router
go tool pprof mem.prof

# Parallel execution (default)
go test -parallel 4 ./...

# Disable test caching
go test -count=1 ./...
```

---

## 🐛 Debugging Tests

### Debug Failing Test

```bash
# Run with verbose output
go test -v ./pkg/cluster -run TestManagerCreate

# Add debug prints
t.Logf("Debug: cluster ID = %s", clusterID)

# Use delve debugger
dlv test ./pkg/cluster -- -test.run TestManagerCreate
```

### Common Issues

**1. Port Already in Use**
```bash
# Find process using port
lsof -i :6379

# Kill process
kill -9 <PID>
```

**2. Docker Services Not Ready**
```bash
# Check service health
docker-compose ps

# View logs
docker-compose logs redis
```

**3. Permission Denied**
```bash
# Fix permissions
chmod -R 755 test/
```

---

## ✅ Testing Checklist

Before committing:

- [ ] All unit tests pass
- [ ] New code has tests
- [ ] Coverage >= 80%
- [ ] Integration tests pass (if applicable)
- [ ] No race conditions (`-race`)
- [ ] Linter passes (`make lint`)
- [ ] Documentation updated

---

## 📚 Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [testify](https://github.com/stretchr/testify) - Testing toolkit
- [gomock](https://github.com/golang/mock) - Mocking framework
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

---

**Happy Testing! 🧪**

