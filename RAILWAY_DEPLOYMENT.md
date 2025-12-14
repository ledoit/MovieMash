# Railway Deployment Guide - Step by Step

## Prerequisites
- Railway account (sign up at https://railway.app)
- GitHub account with your code pushed

---

## Step 1: Create Railway Project

1. Go to https://railway.app
2. Click **"New Project"**
3. Select **"Deploy from GitHub repo"**
4. Authorize Railway to access your GitHub
5. Select your repository: `ledoit/MovieMash`
6. Click **"Deploy Now"**

---

## Step 2: Add PostgreSQL Database

1. In your Railway project dashboard, click **"+ New"**
2. Select **"Database"** → **"Add PostgreSQL"**
3. Wait for it to provision (takes ~30 seconds)
4. Click on the PostgreSQL service
5. Go to the **"Variables"** tab
6. Copy the **`DATABASE_URL`** value (you'll need this)

---

## Step 3: Configure Backend Service

1. In Railway dashboard, you should see a service (might be called "moviemash" or similar)
2. If not, click **"+ New"** → **"GitHub Repo"** → Select your repo
3. Click on the service
4. Go to **"Settings"** tab
5. Set the **Root Directory** to: `backend`
6. Go to **"Variables"** tab
7. Add these environment variables:

   ```
   PORT=8080
   DATABASE_URL=<paste the DATABASE_URL from PostgreSQL service>
   ALLOWED_ORIGIN=https://your-vercel-app.vercel.app
   TMDB_API_KEY=bcca2d98a5dce9840048ca4f0c13656c
   ```

   **Important**: Replace `your-vercel-app.vercel.app` with your actual Vercel domain (you'll get this after deploying frontend)

---

## Step 4: Configure Build Settings

1. Still in the backend service, go to **"Settings"** → **"Build & Deploy"**
2. Set **Build Command**: `go build -o bin/api ./cmd/api`
3. Set **Start Command**: `./bin/api`
4. Set **Watch Paths**: `backend/**`

---

## Step 5: Run Database Migrations

1. In Railway, go to your PostgreSQL service
2. Click **"Connect"** tab
3. Copy the connection details or use the **"psql"** command shown
4. Run migrations:

   ```bash
   # Option 1: Use Railway's built-in terminal
   # Click "Connect" → "Open in Terminal" in PostgreSQL service
   
   # Option 2: Use local psql with Railway connection string
   psql <DATABASE_URL> -f backend/migrations/001_initial_schema.sql
   psql <DATABASE_URL> -f backend/migrations/002_seed_data.sql
   ```

   **OR** use Railway's PostgreSQL service:
   1. Click PostgreSQL service → **"Data"** tab
   2. Click **"Query"** 
   3. Copy contents of `backend/migrations/001_initial_schema.sql` and run
   4. Copy contents of `backend/migrations/002_seed_data.sql` and run

---

## Step 6: Deploy Backend

1. Railway should auto-deploy when you push to main branch
2. Or manually trigger: Go to service → **"Deployments"** → **"Redeploy"**
3. Wait for deployment to complete
4. Go to **"Settings"** → **"Networking"**
5. Click **"Generate Domain"** to get your backend URL
6. Copy this URL (e.g., `https://moviemash-production.up.railway.app`)

---

## Step 7: Update CORS in Backend

1. Go back to backend service → **"Variables"**
2. Update `ALLOWED_ORIGIN` to include your Vercel domain:
   ```
   ALLOWED_ORIGIN=https://your-vercel-app.vercel.app
   ```
3. Railway will auto-redeploy

---

## Step 8: Deploy Frontend to Vercel

1. Go to https://vercel.com
2. Click **"Add New"** → **"Project"**
3. Import your GitHub repository: `ledoit/MovieMash`
4. Configure:
   - **Framework Preset**: Angular
   - **Root Directory**: `./` (leave as is)
   - **Build Command**: `npm run build -- --configuration=production`
   - **Output Directory**: `dist/moviemash/browser`
5. Add Environment Variable:
   - **Name**: `API_URL`
   - **Value**: `https://your-railway-backend-url.up.railway.app/api`
     (Use the URL from Step 6)
6. Click **"Deploy"**
7. Wait for deployment
8. Copy your Vercel URL (e.g., `https://moviemash.vercel.app`)

---

## Step 9: Update Production Environment File

1. Edit `src/environments/environment.prod.ts`
2. Replace the placeholder with your Railway backend URL:
   ```typescript
   apiUrl: 'https://your-railway-backend-url.up.railway.app/api'
   ```
3. Commit and push:
   ```bash
   git add src/environments/environment.prod.ts
   git commit -m "Update production API URL"
   git push origin main
   ```
4. Vercel will auto-redeploy

---

## Step 10: Final CORS Update

1. Go back to Railway → Backend service → **"Variables"**
2. Update `ALLOWED_ORIGIN` with your actual Vercel URL:
   ```
   ALLOWED_ORIGIN=https://moviemash.vercel.app
   ```
3. Railway will auto-redeploy

---

## Step 11: Test Everything

1. Visit your Vercel URL
2. Test the Gladiators page (should load comparisons)
3. Test voting (should work)
4. Test Leaderboard (should show data)
5. Check browser console for any CORS errors

---

## Troubleshooting

### Backend not starting:
- Check Railway logs: Service → **"Deployments"** → Click latest → **"View Logs"**
- Verify `DATABASE_URL` is set correctly
- Verify `PORT` is set (Railway sets this automatically, but check)

### CORS errors:
- Verify `ALLOWED_ORIGIN` matches your Vercel domain exactly
- Check browser console for exact error
- Make sure backend is using the updated CORS middleware

### Database connection errors:
- Verify `DATABASE_URL` is correct
- Check PostgreSQL service is running in Railway
- Verify migrations ran successfully

### Frontend can't reach backend:
- Verify `API_URL` in Vercel environment variables
- Check `environment.prod.ts` has correct URL
- Verify backend is deployed and accessible (try visiting backend URL directly)

---

## Quick Reference

**Railway Backend URL**: `https://your-service.up.railway.app`  
**Vercel Frontend URL**: `https://your-app.vercel.app`  
**Database**: Managed by Railway PostgreSQL service

**Environment Variables Needed:**
- Backend (Railway): `DATABASE_URL`, `ALLOWED_ORIGIN`, `TMDB_API_KEY`, `PORT`
- Frontend (Vercel): `API_URL`

---

## Next Steps After Deployment

1. Set up custom domains (optional)
2. Enable Railway/Vercel analytics
3. Set up monitoring/alerts
4. Consider adding rate limiting to backend
5. Set up automated backups for database

