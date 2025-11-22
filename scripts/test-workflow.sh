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
make test-unit
echo "✅ Tests passed"
echo ""

# Step 5: Build binaries
echo "🔨 Step 5/5: Building binaries..."
mkdir -p bin
go build -o bin/throome ./cmd/throome
go build -o bin/throome-cli ./cmd/throome-cli
echo "✅ Binaries built"
echo ""

echo "🎉 All steps completed successfully!"
echo ""
echo "📦 Binaries available:"
ls -lh bin/
echo ""
echo "🌐 To test the UI:"
echo "   ./bin/throome --port 9000"
echo "   open http://localhost:9000"

