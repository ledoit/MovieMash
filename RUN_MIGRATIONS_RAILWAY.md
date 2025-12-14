# Running Migrations on Railway - Updated Guide

Since Railway's PostgreSQL service doesn't show a Data/Query tab, here are alternative methods:

## Method 1: Railway PostgreSQL Connect Tab (Easiest)

1. **In Railway Dashboard:**
   - Click on your **PostgreSQL** service
   - Click the **"Connect"** tab
   - You'll see connection details

2. **Option A: Use Railway's Built-in Terminal (if available)**
   - Look for **"Open in Terminal"** or **"Connect"** button
   - This opens a psql session directly

3. **Option B: Use Local psql with Railway Connection String**
   - Copy the **"Connection String"** from the Connect tab
   - It looks like: `postgresql://postgres:password@host:port/railway`
   - Open your local terminal and run:

```bash
# Run first migration
psql "postgresql://postgres:password@host:port/railway" -f backend/migrations/001_initial_schema.sql

# Run second migration  
psql "postgresql://postgres:password@host:port/railway" -f backend/migrations/002_seed_data.sql
```

---

## Method 2: Use Railway CLI

1. **Install Railway CLI:**
   ```bash
   npm i -g @railway/cli
   ```

2. **Login:**
   ```bash
   railway login
   ```

3. **Link to your project:**
   ```bash
   railway link
   ```

4. **Connect to PostgreSQL:**
   ```bash
   railway connect postgres
   ```
   This opens a psql session

5. **In psql, run migrations:**
   ```sql
   -- Copy and paste contents of backend/migrations/001_initial_schema.sql
   -- Then copy and paste contents of backend/migrations/002_seed_data.sql
   ```

---

## Method 3: Create a One-Time Migration Service

1. **In Railway Dashboard:**
   - Click **"+ New"** → **"Empty Service"**
   - Name it: `migrations` (or anything)

2. **Set Environment Variables:**
   - Add: `DATABASE_URL` (copy from PostgreSQL service → Variables tab)

3. **Set Start Command:**
   ```bash
   psql $DATABASE_URL -f migrations/001_initial_schema.sql && psql $DATABASE_URL -f migrations/002_seed_data.sql && echo "Migrations complete!"
   ```

4. **Set Root Directory:** `backend`

5. **Deploy:**
   - Railway will run the migrations
   - Check logs to confirm success
   - **Delete the service** after migrations complete

---

## Method 4: Use a Database GUI Tool

1. **Get Connection Details from Railway:**
   - PostgreSQL service → **"Connect"** tab
   - Copy: Host, Port, Database, User, Password

2. **Use a GUI Tool:**
   - **pgAdmin** (free): https://www.pgadmin.org/
   - **TablePlus** (Mac, free trial): https://tableplus.com/
   - **DBeaver** (free): https://dbeaver.io/
   - **Postico** (Mac, paid): https://eggerapps.at/postico/

3. **Connect using Railway credentials:**
   - Host: (from Railway Connect tab)
   - Port: (from Railway Connect tab)
   - Database: `railway` (usually)
   - User: `postgres` (usually)
   - Password: (from Railway)

4. **Run SQL:**
   - Open SQL editor
   - Copy/paste contents of `001_initial_schema.sql`
   - Execute
   - Copy/paste contents of `002_seed_data.sql`
   - Execute

---

## Method 5: Use Docker with Railway Connection

If you have Docker installed:

```bash
# Get connection string from Railway PostgreSQL → Connect tab
# Then run:

docker run -it --rm postgres:alpine psql "YOUR_RAILWAY_CONNECTION_STRING" -f - < backend/migrations/001_initial_schema.sql

docker run -it --rm postgres:alpine psql "YOUR_RAILWAY_CONNECTION_STRING" -f - < backend/migrations/002_seed_data.sql
```

---

## Quick Verification

After running migrations, verify they worked:

```sql
SELECT COUNT(*) FROM movies;
-- Should return: 40

SELECT COUNT(*) FROM top4_sets;
-- Should return: 10

SELECT title FROM movies LIMIT 5;
-- Should show 5 movie titles
```

---

## Recommended: Method 1 (Local psql)

This is usually the easiest:
1. Copy connection string from Railway PostgreSQL → Connect tab
2. Run migrations from your local terminal
3. Done!

If you don't have `psql` installed locally:
- **macOS**: `brew install postgresql`
- **Linux**: `sudo apt-get install postgresql-client` (or similar)
- **Windows**: Download from https://www.postgresql.org/download/windows/

