# Railway Dockerfile Not Found - Fix

## The Issue
Railway says: `Dockerfile 'Dockerfile' does not exist`

## The Solution

Since Railway's **Root Directory** is set to `backend`, Railway automatically looks for `Dockerfile` in that directory.

## What I Fixed

1. ✅ Removed `dockerfilePath` from railway.json and railway.toml
2. ✅ Railway will now auto-detect `backend/Dockerfile`

## Verify in Railway Dashboard

1. Go to your service → **Settings**
2. Check **Root Directory** is set to: `backend`
3. Railway should auto-detect the Dockerfile at `backend/Dockerfile`

## If Still Not Working

Make sure in Railway:
- **Root Directory**: `backend` (not empty, not `.`)
- **Dockerfile Path**: Leave empty (auto-detect) OR set to `Dockerfile`

The Dockerfile exists at `backend/Dockerfile` in your repo, so with root directory set to `backend`, Railway should find it automatically.

