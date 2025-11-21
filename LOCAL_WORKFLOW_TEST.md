# 🧪 Local Workflow Testing Guide

## ✅ Issue Fixed

**Problem**: PostCSS config was using ES module syntax but package.json didn't specify module type.

**Solution**: Added `"type": "module"` to `ui/package.json`

---

## 📋 Test Workflow Locally

Follow these steps to simulate the GitHub Actions workflow on your local machine:

### **Step 1: Build UI**

```bash
cd /Users/akshitmadan/Documents/Akshit_Madan/throome

# Navigate to UI directory
cd ui

# Install dependencies (if not already done)
npm install --legacy-peer-deps

# Build UI
npm run build

# Go back to root
cd ..
```

**Expected output:**
```
✓ 1539 modules transformed.
dist/index.html                   0.47 kB │ gzip:  0.30 kB
dist/assets/index-BZOmT6NR.css   17.01 kB │ gzip:  3.87 kB
dist/assets/index-RQDycIvf.js   275.01 kB │ gzip: 88.46 kB
✓ built in 1.14s
```

### **Step 2: Copy UI to Go Package**

```bash
# Create directory
mkdir -p pkg/gateway/ui

# Copy built files
cp -r ui/dist pkg/gateway/ui/
```

### **Step 3: Build Go Code**

```bash
# Build all Go packages
go build ./...

# Build throome binary
go build -o throome ./cmd/throome

# Build CLI
go build -o throome-cli ./cmd/throome-cli
```

**Expected**: No errors, binaries created

### **Step 4: Run Tests**

```bash
# Run unit tests
make test-unit

# Or directly
go test -short ./...
```

### **Step 5: Test the Binary**

```bash
# Run throome (will fail without proper config, but that's OK)
./throome --version

# Or start it
./throome --port 9000 &

# Open browser
open http://localhost:9000

# You should see the dashboard UI!

# Stop throome
pkill throome
```

---

## 🔄 Complete Local CI/CD Simulation

This script simulates the entire GitHub Actions workflow:

```bash
#!/bin/bash
set -e

echo "🚀 Simulating GitHub Actions Workflow..."
echo ""

# Step 1: Build UI
echo "📦 Step 1/5: Building UI..."
cd ui
npm ci --legacy-peer-deps
npm run build
cd ..
echo "✅ UI built successfully"
echo ""

# Step 2: Copy UI
echo "📁 Step 2/5: Copying UI to Go package..."
mkdir -p pkg/gateway/ui
cp -r ui/dist pkg/gateway/ui/
echo "✅ UI copied"
echo ""

# Step 3: Download Go dependencies
echo "📥 Step 3/5: Downloading Go dependencies..."
go mod download
echo "✅ Dependencies downloaded"
echo ""

# Step 4: Run tests
echo "🧪 Step 4/5: Running tests..."
go test -short ./...
echo "✅ Tests passed"
echo ""

# Step 5: Build binaries
echo "🔨 Step 5/5: Building binaries..."
go build -o bin/throome ./cmd/throome
go build -o bin/throome-cli ./cmd/throome-cli
echo "✅ Binaries built"
echo ""

echo "🎉 All steps completed successfully!"
echo ""
echo "Binaries available:"
ls -lh bin/
```

Save this as `test-workflow.sh` and run:

```bash
chmod +x test-workflow.sh
./test-workflow.sh
```

---

## 🐳 Test Docker Build Locally

```bash
# Build Docker image (full multi-stage build)
docker build -f deployments/docker/Dockerfile -t throome:test .

# This will:
# 1. Build UI in Node.js container
# 2. Build Go with embedded UI
# 3. Create final runtime image

# Run the image
docker run -p 9000:9000 throome:test

# Open browser
open http://localhost:9000
```

---

## 🔍 Verify UI Embedding

Check that UI is properly embedded:

```bash
# Build throome
go build -o throome ./cmd/throome

# Check binary size (should be ~15-20MB with UI)
ls -lh throome

# Run and check logs
./throome --port 9000
# Should start without "UI not available" errors

# In another terminal, test UI endpoint
curl http://localhost:9000/
# Should return HTML

# Test API endpoint
curl http://localhost:9000/api/v1/health
# Should return JSON
```

---

## 🚨 Common Issues & Fixes

### Issue: "pattern ui/dist: no matching files found"

**Cause**: UI not built before Go build

**Fix**:
```bash
cd ui && npm run build && cd ..
mkdir -p pkg/gateway/ui
cp -r ui/dist pkg/gateway/ui/
go build ./...
```

### Issue: "PostCSS config error"

**Cause**: Missing `"type": "module"` in package.json

**Fix**: Already applied! ✅

### Issue: "npm install fails"

**Cause**: Peer dependency conflicts

**Fix**:
```bash
npm install --legacy-peer-deps
```

### Issue: "UI shows 404"

**Cause**: UI handler not catching all routes

**Fix**: Already configured in `pkg/gateway/server.go` ✅

---

## ✅ Checklist

Before pushing to GitHub, verify locally:

- [ ] UI builds without errors (`npm run build`)
- [ ] UI files exist at `ui/dist/`
- [ ] Files copied to `pkg/gateway/ui/dist/`
- [ ] Go code compiles (`go build ./...`)
- [ ] Unit tests pass (`make test-unit`)
- [ ] Binary runs (`./throome`)
- [ ] UI loads at `http://localhost:9000`
- [ ] API works at `http://localhost:9000/api/v1/health`
- [ ] Connection status shows green ✅

---

## 📊 Expected Results

### UI Build Output:
```
✓ 1539 modules transformed.
dist/index.html                   0.47 kB
dist/assets/index-*.css          17.01 kB
dist/assets/index-*.js          275.01 kB
✓ built in ~1s
```

### Go Build Output:
```
(no output = success)
```

### Test Output:
```
ok  	github.com/akmadan/throome/internal/utils	0.906s
ok  	github.com/akmadan/throome/pkg/cluster	    0.550s
ok  	github.com/akmadan/throome/test/integration	0.895s
```

### Binary Size:
```
-rwxr-xr-x  1 user  staff   18M Nov 21 10:00 throome
-rwxr-xr-x  1 user  staff   12M Nov 21 10:00 throome-cli
```

---

## 🎯 Next Steps

Once all local tests pass:

1. **Commit changes**:
   ```bash
   git add .
   git commit -m "fix: resolve PostCSS config ES module issue"
   ```

2. **Push to GitHub**:
   ```bash
   git push origin main
   ```

3. **Monitor workflows**:
   - Go to: https://github.com/akmadan/throome/actions
   - Watch workflows run
   - All should pass ✅

4. **Test Docker image**:
   ```bash
   # After CI completes
   docker pull akshitmadan/throome:latest
   docker run -p 9000:9000 akshitmadan/throome:latest
   open http://localhost:9000
   ```

---

**Status**: ✅ Ready to push!

**Last Updated**: 2025-11-21

