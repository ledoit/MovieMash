# Railway Build Fix

## Problem
Railway was failing with: `error: undefined variable 'go_1_23'`

## Solution Applied
1. ✅ Switched from Nixpacks to Dockerfile (more reliable)
2. ✅ Fixed Go version in nixpacks.toml (changed `go_1_23` to `go`)
3. ✅ Created proper Dockerfile for Go 1.23
4. ✅ Updated railway.json and railway.toml to use Dockerfile

## What You Need to Do in Railway

1. **Go to your service settings:**
   - Click on your backend service in Railway
   - Go to **"Settings"** tab
   - Scroll to **"Build & Deploy"** section

2. **Set Root Directory:**
   - **Root Directory**: `backend`
   - This tells Railway to build from the backend folder

3. **Verify Build Settings:**
   - Railway should auto-detect the Dockerfile
   - If not, set **Dockerfile Path**: `Dockerfile` (relative to backend folder)

4. **Redeploy:**
   - Railway should auto-redeploy after the git push
   - Or manually: **"Deployments"** → **"Redeploy"**

## Expected Build Process

The Dockerfile will:
1. Use Go 1.23 Alpine image
2. Copy go.mod and go.sum
3. Download dependencies
4. Copy source code
5. Build the binary: `bin/api`
6. Create minimal Alpine image with just the binary
7. Run on port 8080

## If Build Still Fails

Check the build logs for:
- Database connection errors (need DATABASE_URL set)
- Missing dependencies (should be fixed now)
- Port issues (Railway sets PORT automatically)

The build should work now! 🚀

