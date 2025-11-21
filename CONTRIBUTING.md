# Contributing to Throome

Thank you for your interest in contributing to Throome! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates.

**When submitting a bug report, include:**
- Clear, descriptive title
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:
- Clear use case
- Detailed description of the proposed functionality
- Any examples or mockups
- Why this would be useful to most users

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **Make your changes**
   - Follow the code style guidelines
   - Add tests for new functionality
   - Update documentation as needed

4. **Commit your changes**
   ```bash
   git commit -m 'Add amazing feature'
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/amazing-feature
   ```

6. **Open a Pull Request**

## Development Setup

### Prerequisites

- Go 1.21+
- Make
- Docker (for integration tests)
- Git

### Local Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/throome.git
cd throome

# Add upstream remote
git remote add upstream https://github.com/akmadan/throome.git

# Install dependencies
make deps

# Build
make build

# Run tests
make test
```

## Code Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing

### Code Structure

```
throome/
├── cmd/          # Command-line applications
├── pkg/          # Public libraries
├── internal/     # Private application code
├── configs/      # Configuration examples
├── docs/         # Documentation
├── examples/     # Example applications
└── test/         # Tests
```

### Naming Conventions

- **Packages**: Short, lowercase, single-word names
- **Files**: Lowercase with underscores (e.g., `cluster_manager.go`)
- **Types**: PascalCase (e.g., `ClusterManager`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase

### Comments

- All exported types, functions, and constants must have comments
- Comments should be complete sentences
- Begin with the name of the element being described

Example:
```go
// ClusterManager manages the lifecycle of clusters.
// It provides methods to create, update, and delete clusters.
type ClusterManager struct {
    // ...
}

// Create creates a new cluster with the given configuration.
func (m *ClusterManager) Create(config *Config) error {
    // ...
}
```

## Testing

### Unit Tests

```bash
make test-unit
```

### Integration Tests

```bash
# Start test services with Docker
docker-compose -f test/docker-compose.yml up -d

# Run integration tests
make test-integration

# Stop test services
docker-compose -f test/docker-compose.yml down
```

### Writing Tests

- Place test files next to the code they test
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for >80% code coverage

Example:
```go
func TestClusterManager_Create(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: &Config{
                ClusterID: "test-01",
                Name:      "Test Cluster",
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := NewManager("./testdata")
            _, err := m.Create(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

- Update README.md for user-facing changes
- Add/update docs in `docs/` for detailed guides
- Include code examples where helpful
- Update CHANGELOG.md for notable changes

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(adapters): add MongoDB adapter

Implements MongoDB adapter with connection pooling and basic CRUD operations.

Closes #123
```

```
fix(router): prevent panic on nil adapter

Add nil check before accessing adapter methods.

Fixes #456
```

## Pull Request Process

1. **Update your fork**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Ensure all tests pass**
   ```bash
   make test
   make lint
   ```

3. **Update documentation** as needed

4. **Create a PR** with:
   - Clear description of changes
   - Link to related issues
   - Screenshots for UI changes
   - Breaking changes noted

5. **Respond to feedback** promptly

6. **PR will be merged** once approved by maintainers

## Building New Adapters

See [Adapter Development Guide](docs/adapter-development.md) for details on creating custom adapters.

## License

By contributing to Throome, you agree that your contributions will be licensed under the Apache License 2.0.

## Questions?

- Open a [GitHub Discussion](https://github.com/akmadan/throome/discussions)
- Comment on an existing issue
- Reach out to maintainers

Thank you for contributing! 🎉

