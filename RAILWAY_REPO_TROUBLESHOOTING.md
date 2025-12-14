# Railway Can't Find Repository - Fix Guide

Your repository is: `ledoit/MovieMash`

## Quick Fix Steps

### Option 1: Reconnect GitHub (Most Common Fix)

1. **In Railway Dashboard:**
   - Click your profile icon (top right)
   - Go to **"Account Settings"**
   - Find **"GitHub"** section
   - Click **"Disconnect"** or **"Revoke Access"**
   - Then click **"Connect GitHub"** again
   - Authorize Railway to access your repositories

2. **In GitHub:**
   - Go to: https://github.com/settings/applications
   - Click **"Authorized OAuth Apps"**
   - Find **"Railway"**
   - Click **"Revoke"**
   - Go back to Railway and reconnect

### Option 2: Install Railway GitHub App

1. Go to: https://github.com/apps/railway
2. Click **"Install"**
3. Choose:
   - **"Only select repositories"** → Select `MovieMash`
   - OR **"All repositories"** (if you're okay with that)
4. Click **"Install"**
5. Go back to Railway and try connecting again

### Option 3: Manual Deploy (If GitHub Still Doesn't Work)

1. In Railway, click **"+ New"**
2. Select **"Empty Project"**
3. Click **"+ New"** → **"GitHub Repo"**
4. If it still doesn't show, try:
   - **"Deploy from GitHub repo"** → Search for `MovieMash`
   - OR use **"Deploy from public Git repository"**:
     - URL: `https://github.com/ledoit/MovieMash.git`

### Option 4: Check Repository Visibility

1. Go to: https://github.com/ledoit/MovieMash/settings
2. Scroll to **"Danger Zone"**
3. If it's private, Railway needs access to private repos
4. Make sure Railway GitHub App has **"Private repository access"** enabled

### Option 5: Use Railway CLI (Alternative)

If web interface doesn't work:

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Link to your project
railway link

# Deploy
railway up
```

## Still Not Working?

1. **Check Railway Status**: https://status.railway.app
2. **Try Different Browser**: Sometimes browser extensions block OAuth
3. **Clear Browser Cache**: Clear cookies for railway.app
4. **Contact Railway Support**: support@railway.app

## Alternative: Deploy Backend Manually

If Railway GitHub connection is problematic, you can:

1. **Use Railway CLI** (see above)
2. **Use Render.com** instead (similar service, sometimes easier GitHub integration)
3. **Use Fly.io** (also supports Go + PostgreSQL)

---

## Quick Test

After reconnecting, try:
1. Railway Dashboard → **"+ New"** → **"Deploy from GitHub repo"**
2. Search for: `MovieMash`
3. Should show: `ledoit/MovieMash`

If it still doesn't appear, the GitHub app likely needs to be reinstalled.

