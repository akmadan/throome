# 🛠️ Development Workflow Guide

## Quick Answer

### **For UI Changes ONLY:**
```bash
# 1. Make your UI changes in ui/src/...
# 2. Test locally (optional but recommended)
cd ui && npm run dev

# 3. Just commit and push!
git add .
git commit -m "feat: your UI changes"
git push origin main
```
**GitHub Actions will build the UI automatically** ✅

### **For Go Changes ONLY:**
```bash
# 1. Make your Go changes
# 2. Test locally (optional but recommended)
go test ./...

# 3. Just commit and push!
git add .
git commit -m "feat: your Go changes"
git push origin main
```
**No special commands needed** ✅

### **For Both UI + Go Changes:**
```bash
# 1. Make your changes
# 2. Test locally (recommended)
./test-workflow.sh

# 3. Just commit and push!
git add .
git commit -m "feat: your changes"
git push origin main
```

---

## 🎯 The Simple Truth

**You can just commit and push as usual!** 🎉

GitHub Actions will:
1. ✅ Build the UI
2. ✅ Copy UI files
3. ✅ Build Go with embedded UI
4. ✅ Run all tests
5. ✅ Create Docker image

**You don't NEED to run any special commands before pushing.**

---

## 🧪 But... Testing Locally is Recommended

While GitHub will handle everything, it's **good practice** to test locally first to catch errors early.

### Quick Local Test Commands

#### **Test UI Changes:**
```bash
cd ui
npm run build
cd ..
```
If this succeeds, your UI changes are good ✅

#### **Test Go Changes:**
```bash
go build ./...
go test ./...
```
If these succeed, your Go changes are good ✅

#### **Test Everything Together:**
```bash
./test-workflow.sh
```
This simulates the full CI/CD pipeline locally.

---

## 📋 Recommended Workflows

### **Scenario 1: Quick UI Fix (e.g., fix typo, change color)**

```bash
# 1. Edit files in ui/src/
vim ui/src/components/Dashboard.tsx

# 2. Preview (optional)
cd ui && npm run dev
# Check at http://localhost:3000

# 3. Commit and push
git add .
git commit -m "fix: typo in dashboard"
git push
```

**Time**: ~30 seconds
**Local build**: Not required

---

### **Scenario 2: Significant UI Feature**

```bash
# 1. Edit files in ui/src/
vim ui/src/pages/NewFeature.tsx

# 2. Test locally
cd ui && npm run build
# Make sure it builds without errors

# 3. Optional: Test with Go
cd ..
cp -r ui/dist pkg/gateway/ui/
go build -o throome ./cmd/throome
./throome --port 9000
# Check at http://localhost:9000

# 4. Commit and push
git add .
git commit -m "feat: add new feature page"
git push
```

**Time**: ~2-3 minutes
**Local build**: Recommended

---

### **Scenario 3: Go Backend Changes**

```bash
# 1. Edit Go files
vim pkg/gateway/server.go

# 2. Test locally
go build ./...
go test ./...

# 3. Commit and push
git add .
git commit -m "feat: add new API endpoint"
git push
```

**Time**: ~1 minute
**Local build**: Recommended

---

### **Scenario 4: Changes to Both UI and Go**

```bash
# 1. Make your changes
vim ui/src/components/NewComponent.tsx
vim pkg/gateway/handlers.go

# 2. Run full test
./test-workflow.sh

# 3. Commit and push
git add .
git commit -m "feat: add new feature with API integration"
git push
```

**Time**: ~3-4 minutes
**Local build**: Highly recommended

---

## 🚫 What You DON'T Need to Do

### ❌ Don't manually copy UI files to Git
```bash
# ❌ DON'T DO THIS before pushing:
cp -r ui/dist pkg/gateway/ui/
git add pkg/gateway/ui/dist
```

**Why?** 
- `pkg/gateway/ui/dist/` is in `.gitignore`
- GitHub Actions builds it fresh
- Keeps repo size small

### ❌ Don't commit node_modules
```bash
# ❌ DON'T DO THIS:
git add ui/node_modules
```

**Why?**
- Already in `.gitignore`
- ~300MB of dependencies
- npm installs them automatically

### ❌ Don't commit built binaries
```bash
# ❌ DON'T DO THIS:
git add bin/throome
git add throome
```

**Why?**
- Binaries are built by CI/CD
- Different for each platform
- Large file sizes

---

## ✅ What to Commit

### **Always commit:**
- ✅ Source code: `ui/src/**/*.tsx`, `pkg/**/*.go`, `cmd/**/*.go`
- ✅ Config files: `ui/package.json`, `go.mod`, `go.sum`
- ✅ Workflows: `.github/workflows/*.yml`
- ✅ Documentation: `*.md`

### **Never commit:**
- ❌ `ui/node_modules/`
- ❌ `ui/dist/`
- ❌ `pkg/gateway/ui/dist/`
- ❌ `bin/`
- ❌ Built binaries (`throome`, `throome-cli`)

---

## 🎨 Daily Development Workflow

### **Morning: Start working**
```bash
# Pull latest changes
git pull origin main

# Install any new dependencies (if package.json changed)
cd ui && npm install --legacy-peer-deps && cd ..

# Start UI dev server (optional, for live preview)
cd ui && npm run dev &
cd ..
```

### **During the day: Make changes**
```bash
# Edit files as needed
vim ui/src/...
vim pkg/...

# Preview changes at http://localhost:3000 (if dev server running)
# Dev server has hot reload, so changes appear instantly!
```

### **Before lunch/breaks: Quick commit**
```bash
git add .
git commit -m "wip: working on feature X"
git push
```

### **End of day: Final commit**
```bash
# Optional: Test everything
./test-workflow.sh

# Commit
git add .
git commit -m "feat: completed feature X"
git push
```

---

## 🐛 Pre-Push Checks (Optional but Recommended)

Create a pre-push checklist:

```bash
# Quick checks before pushing big changes:

# 1. Does UI build?
cd ui && npm run build && cd ..

# 2. Does Go compile?
go build ./...

# 3. Do tests pass?
go test ./...

# If all pass, you're good to push! ✅
```

You can create an alias for this:

```bash
# Add to your ~/.zshrc or ~/.bashrc
alias throome-check='cd ui && npm run build && cd .. && go build ./... && go test ./...'

# Then just run:
throome-check
```

---

## 🚀 The Lazy Developer's Workflow

**Minimum viable workflow:**

```bash
# Make changes
vim somefile.tsx

# Commit and push
git add .
git commit -m "feat: my changes"
git push
```

**GitHub Actions will:**
- Build UI ✅
- Test everything ✅
- Create Docker image ✅
- Notify you if something fails ✅

**If tests fail on GitHub:**
- You'll get an email
- Fix the issue
- Push again

**This works! But testing locally saves time.**

---

## ⚡ Pro Tips

### **1. Use Git Hooks (Optional)**

Create `.git/hooks/pre-push`:

```bash
#!/bin/bash
echo "🧪 Running pre-push checks..."

# Build UI
cd ui && npm run build && cd .. || exit 1

# Test Go
go test ./... || exit 1

echo "✅ Pre-push checks passed!"
exit 0
```

```bash
chmod +x .git/hooks/pre-push
```

Now tests run automatically before every push!

### **2. Use Makefile Shortcuts**

```bash
# Already available:
make build          # Build everything
make test           # Run all tests
make test-unit      # Quick unit tests
make clean          # Clean build artifacts

# Use these before pushing:
make test-unit && git push
```

### **3. Use Branch Protection**

On GitHub:
- Settings → Branches → Add rule
- Require tests to pass before merging
- Prevents broken code from reaching main

---

## 📊 What Happens in CI/CD

When you push to GitHub:

```
Your Push
    ↓
GitHub Actions Triggered
    ↓
1. Install Node.js
2. Build UI (npm run build)
3. Copy ui/dist → pkg/gateway/ui/dist
4. Install Go
5. Download Go dependencies
6. Run tests
7. Build binaries
8. Build Docker image
9. Push to DockerHub
    ↓
✅ Deploy Complete!
```

**Total time**: ~10-15 minutes

---

## 🎯 Summary

### **Bare Minimum (works fine):**
```bash
git add .
git commit -m "your changes"
git push
```

### **Recommended (catches errors early):**
```bash
./test-workflow.sh
git add .
git commit -m "your changes"
git push
```

### **Pro Level (automated):**
```bash
# Set up pre-push hooks
# Just make changes and push
# Everything tests automatically
```

---

## ❓ Common Questions

**Q: Do I need to rebuild UI every time I change Go code?**
**A:** No! Only if you're testing the embedded UI locally.

**Q: Do I need to build Go every time I change UI?**
**A:** No! Just test UI with `npm run dev`.

**Q: Can I skip local testing entirely?**
**A:** Yes! GitHub will test everything. But you'll wait 10-15 minutes to see results.

**Q: What if CI/CD fails after I push?**
**A:** Check GitHub Actions logs, fix the issue, push again. No big deal!

**Q: How do I test the full Docker build locally?**
**A:** `docker build -f deployments/docker/Dockerfile -t throome:test .`

---

## 🎉 TL;DR

**Just commit and push as normal!** 🎊

GitHub Actions handles everything automatically.

**Optional but nice:** Run `./test-workflow.sh` before pushing to catch errors early.

---

**Happy Coding!** 🚀

