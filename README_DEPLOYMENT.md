# Deployment Guide

## Current Status
**The app will NOT work on Vercel as-is** because:
1. The Go backend needs to be deployed separately (Vercel doesn't support Go)
2. PostgreSQL database needs to be hosted
3. CORS needs to be configured for production

## Steps to Deploy

### 1. Deploy Go Backend
Choose one:
- **Railway**: Easy PostgreSQL + Go deployment
- **Render**: Free tier available
- **Fly.io**: Good for Go apps
- **VPS**: DigitalOcean, AWS, etc.

Update CORS in `backend/cmd/api/main.go` to allow your Vercel domain.

### 2. Host PostgreSQL Database
Options:
- **Supabase**: Free tier, easy setup
- **Neon**: Serverless Postgres
- **Railway**: Includes database
- **Render**: Includes database

Update `DATABASE_URL` in your backend environment.

### 3. Deploy Frontend to Vercel

1. Push code to GitHub
2. Connect repo to Vercel
3. Set environment variable:
   - `API_URL`: Your backend URL (e.g., `https://your-backend.railway.app/api`)
4. Build command: `npm run build`
5. Output directory: `dist/moviemash/browser`

### 4. Update Environment File
Edit `src/environments/environment.prod.ts` with your actual backend URL.

### Quick Deploy (Railway + Vercel)
1. Deploy backend + database on Railway
2. Get Railway URL
3. Deploy frontend on Vercel with `API_URL` env var
4. Update CORS on backend to allow Vercel domain

