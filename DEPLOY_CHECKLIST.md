# Deployment Checklist

## ✅ What I've Done For You

- [x] Updated CORS to use environment variable `ALLOWED_ORIGIN`
- [x] Created Railway configuration files
- [x] Created migration script
- [x] Updated Angular to use environment variables
- [x] Created comprehensive deployment guide

## 📋 Your Step-by-Step Checklist

### Railway Backend Setup

1. [ ] Sign up/login to Railway: https://railway.app
2. [ ] Create new project → "Deploy from GitHub repo"
3. [ ] Select your repository: `ledoit/MovieMash`
4. [ ] Add PostgreSQL database service
5. [ ] Copy `DATABASE_URL` from PostgreSQL service
6. [ ] Configure backend service:
   - Root Directory: `backend`
   - Build Command: `go build -o bin/api ./cmd/api`
   - Start Command: `./bin/api`
7. [ ] Add environment variables:
   - `DATABASE_URL` (from PostgreSQL)
   - `ALLOWED_ORIGIN` (set to your Vercel URL after deployment)
   - `TMDB_API_KEY=bcca2d98a5dce9840048ca4f0c13656c`
8. [ ] Run database migrations (see RAILWAY_DEPLOYMENT.md Step 5)
9. [ ] Generate Railway domain for backend
10. [ ] Copy backend URL (e.g., `https://xxx.up.railway.app`)

### Vercel Frontend Setup

11. [ ] Sign up/login to Vercel: https://vercel.com
12. [ ] Import GitHub repository: `ledoit/MovieMash`
13. [ ] Configure build:
    - Framework: Angular
    - Build Command: `npm run build -- --configuration=production`
    - Output Directory: `dist/moviemash/browser`
14. [ ] Add environment variable:
    - `API_URL` = `https://your-railway-backend-url.up.railway.app/api`
15. [ ] Deploy and copy Vercel URL

### Final Configuration

16. [ ] Update `src/environments/environment.prod.ts` with Railway backend URL
17. [ ] Update Railway `ALLOWED_ORIGIN` with Vercel URL
18. [ ] Push changes to trigger redeploy
19. [ ] Test the live site!

## 🔗 Quick Links

- **Railway Dashboard**: https://railway.app/dashboard
- **Vercel Dashboard**: https://vercel.com/dashboard
- **Full Guide**: See `RAILWAY_DEPLOYMENT.md`

## ⚠️ Important Notes

- Railway sets `PORT` automatically - don't override it
- `DATABASE_URL` is auto-provided by Railway PostgreSQL service
- Update `ALLOWED_ORIGIN` after you get your Vercel URL
- Both services will auto-redeploy on git push

