# CI/CD Pipeline Setup Guide

This document explains how to set up and use the CI/CD pipeline for Throome.

## 🚀 Overview

Throome uses GitHub Actions for CI/CD with three main workflows:

1. **Tests** - Automated testing on every push/PR
2. **Docker** - Build and push Docker images
3. **Release** - Create GitHub releases with binaries

---

## 📋 Prerequisites

### 1. DockerHub Account

Create a DockerHub account and repository:

1. Go to [hub.docker.com](https://hub.docker.com)
2. Create account or sign in
3. Create a new repository: `akshitmadan/throome`
4. Generate access token:
   - Account Settings → Security → New Access Token
   - Name: `github-actions`
   - Permissions: Read & Write
   - Copy the token

### 2. GitHub Secrets

Add these secrets to your GitHub repository:

**Settings → Secrets and variables → Actions → New repository secret**

| Secret Name | Description | Value |
|-------------|-------------|-------|
| `DOCKERHUB_USERNAME` | Your DockerHub username | `your-username` |
| `DOCKERHUB_TOKEN` | DockerHub access token | Token from step 1 |

---

## 🔄 Workflows

### 1. Test Workflow (`.github/workflows/test.yml`)

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main` or `develop`

**Jobs:**
- **Unit Tests** - Fast tests without external dependencies
- **Integration Tests** - Tests with Redis, PostgreSQL
- **Coverage** - Code coverage report (uploaded to Codecov)
- **Lint** - Code quality checks
- **Build** - Multi-platform binary builds

**Status Badge:**
```markdown
![Tests](https://github.com/akmadan/throome/actions/workflows/test.yml/badge.svg)
```

### 2. Docker Workflow (`.github/workflows/docker.yml`)

**Triggers:**
- Push to `main` branch
- Push tags matching `v*.*.*`
- Pull requests to `main` (build only, no push)

**Jobs:**
- **Build** - Test Docker image builds
- **Push** - Push to DockerHub (main branch or tags only)

**Docker Tags Generated:**
- `throome/throome:latest` - Latest main branch
- `throome/throome:v1.2.3` - Specific version
- `throome/throome:v1.2` - Minor version
- `throome/throome:v1` - Major version
- `throome/throome:main-abc1234` - Commit SHA

**Multi-arch Support:**
- `linux/amd64`
- `linux/arm64`

### 3. Release Workflow (`.github/workflows/release.yml`)

**Triggers:**
- Push tags matching `v*.*.*`

**Artifacts:**
- Binaries for Linux, macOS, Windows (amd64 & arm64)
- Docker images on DockerHub and GitHub Container Registry
- GitHub Release with changelog

**Created Files:**
```
throome-linux-amd64.tar.gz
throome-linux-arm64.tar.gz
throome-darwin-amd64.tar.gz
throome-darwin-arm64.tar.gz
throome-windows-amd64.zip
throome-cli-linux-amd64.tar.gz
... (CLI for all platforms)
```

---

## 🎯 Usage Guide

### Running Tests Automatically

Tests run automatically on every push and PR:

```bash
# Push code
git add .
git commit -m "Add new feature"
git push origin main

# GitHub Actions will:
# 1. Run unit tests
# 2. Run integration tests with Redis/PostgreSQL
# 3. Generate coverage report
# 4. Run linter
# 5. Build binaries for all platforms
```

### Creating a Release

1. **Update version in code** (if needed):
```go
// cmd/throome/main.go
var Version = "0.2.0"  // Update this
```

2. **Update CHANGELOG.md**:
```markdown
## [0.2.0] - 2025-11-20
### Added
- New feature X
- Improvement Y

### Fixed
- Bug Z
```

3. **Create and push tag**:
```bash
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

4. **GitHub Actions will automatically**:
   - Build binaries for all platforms
   - Create GitHub release
   - Push Docker images
   - Update DockerHub description

### Building Docker Image Locally

```bash
# Build image
docker build -t throome:local -f deployments/docker/Dockerfile .

# Run container
docker run -p 9000:9000 throome:local

# Test
curl http://localhost:9000/health
```

### Using Docker Compose

```bash
# Start full stack (Gateway + Redis + PostgreSQL + Kafka)
cd deployments/docker
docker-compose up -d

# View logs
docker-compose logs -f throome

# Stop all services
docker-compose down
```

---

## 📊 CI/CD Status

Add these badges to your README:

```markdown
# Throome

![Tests](https://github.com/akmadan/throome/actions/workflows/test.yml/badge.svg)
![Docker](https://github.com/akmadan/throome/actions/workflows/docker.yml/badge.svg)
![Release](https://github.com/akmadan/throome/actions/workflows/release.yml/badge.svg)
[![codecov](https://codecov.io/gh/akshitmadan/throome/branch/main/graph/badge.svg)](https://codecov.io/gh/akshitmadan/throome)
[![Go Report Card](https://goreportcard.com/badge/github.com/akmadan/throome)](https://goreportcard.com/report/github.com/akmadan/throome)
[![Docker Pulls](https://img.shields.io/docker/pulls/throome/throome)](https://hub.docker.com/r/throome/throome)
[![License](https://img.shields.io/github/license/akshitmadan/throome)](LICENSE)
```

---

## 🐛 Troubleshooting

### Tests Failing in CI but Pass Locally

**Problem:** Integration tests fail because services aren't ready.

**Solution:** The workflow includes health checks. If still failing:
```yaml
# Add longer wait time in .github/workflows/test.yml
- name: Wait for services
  run: sleep 15  # Increase this
```

### Docker Push Fails

**Problem:** `unauthorized: authentication required`

**Solution:**
1. Check `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets
2. Verify token has Read & Write permissions
3. Ensure DockerHub repository exists: `throome/throome`

### Build Fails on ARM64

**Problem:** Some dependencies don't support ARM64.

**Solution:**
```yaml
# Temporarily remove ARM64
platforms: linux/amd64  # Remove linux/arm64
```

### Linter Failures

**Problem:** Code doesn't pass linting checks.

**Solution:**
```bash
# Run linter locally
make lint

# Auto-fix issues
golangci-lint run --fix

# Check specific file
golangci-lint run path/to/file.go
```

---

## 🔐 Security

### CodeQL Analysis

Automatic security scanning runs:
- On every push to main/develop
- On pull requests
- Weekly on Sunday

View results: **Security → Code scanning alerts**

### Dependabot

Enable Dependabot for automatic dependency updates:

**Settings → Security → Dependabot**
- ✅ Enable Dependabot alerts
- ✅ Enable Dependabot security updates
- ✅ Enable Dependabot version updates

---

## 📈 Monitoring

### GitHub Actions

- **Actions tab** - View all workflow runs
- **Insights → Actions** - Usage statistics
- **Failed runs** - Email notifications (configure in Settings)

### DockerHub

- **Tags** - View all pushed images
- **Pulls** - Download statistics
- **Security** - Vulnerability scanning

### Codecov

- Sign up at [codecov.io](https://codecov.io)
- Link GitHub repository
- View coverage trends and reports

---

## 🎨 Customization

### Change Docker Image Name

Edit `.github/workflows/docker.yml`:
```yaml
env:
  IMAGE_NAME: your-username/throome  # Change this
```

### Add New Test Services

Edit `.github/workflows/test.yml`:
```yaml
services:
  mongodb:  # Add new service
    image: mongo:7
    ports:
      - 27017:27017
```

### Modify Release Notes

Edit `.github/workflows/release.yml`:
```yaml
body: |
  # Your custom release notes template
```

---

## 📚 Best Practices

### 1. Branch Protection

Enable branch protection for `main`:
- **Settings → Branches → Add rule**
- Branch name pattern: `main`
- ✅ Require status checks to pass
  - Select: `Unit Tests`, `Integration Tests`, `Lint`
- ✅ Require branches to be up to date
- ✅ Require linear history

### 2. Semantic Versioning

Follow [semver.org](https://semver.org):
- `v1.0.0` - Major version (breaking changes)
- `v1.1.0` - Minor version (new features)
- `v1.1.1` - Patch version (bug fixes)

### 3. Changelog

Keep CHANGELOG.md updated:
```markdown
## [Unreleased]
### Added
- New feature

## [1.0.0] - 2025-11-20
### Added
- Initial release
```

### 4. Test Before Release

```bash
# Always test before creating a tag
make test
make test-integration
make build

# Then create release
git tag -a v1.0.0 -m "Release v1.0.0"
```

---

## 🚀 Quick Start Checklist

- [ ] Add DockerHub secrets to GitHub
- [ ] Push code to trigger first workflow
- [ ] Verify tests pass
- [ ] Create first release tag
- [ ] Verify Docker image on DockerHub
- [ ] Test downloading and running image
- [ ] Add status badges to README
- [ ] Enable branch protection
- [ ] Set up Codecov (optional)

---

## 📞 Support

Issues with CI/CD?
1. Check workflow logs in GitHub Actions
2. Review this guide
3. Open an issue: [github.com/akmadan/throome/issues](https://github.com/akmadan/throome/issues)

---

**Happy Shipping! 🚢**

