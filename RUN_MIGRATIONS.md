# How to Run Database Migrations on Railway

## Method 1: Using Railway's PostgreSQL Query Interface (Easiest)

1. **In Railway Dashboard:**
   - Go to your project
   - Click on your **PostgreSQL** service
   - Click the **"Data"** tab at the top
   - Click **"Query"** button

2. **Run First Migration:**
   - Open `backend/migrations/001_initial_schema.sql` in your code editor
   - Copy **ALL** the contents (Ctrl+A, Ctrl+C)
   - Paste into the Railway Query editor
   - Click **"Run"** or press Ctrl+Enter
   - Wait for "Success" message

3. **Run Second Migration:**
   - Open `backend/migrations/002_seed_data.sql` in your code editor
   - Copy **ALL** the contents
   - Paste into the Railway Query editor
   - Click **"Run"**
   - Wait for "Success" message

4. **Verify:**
   - In the Query editor, run: `SELECT COUNT(*) FROM movies;`
   - Should return `40` (you have 40 movies)
   - Run: `SELECT COUNT(*) FROM top4_sets;`
   - Should return `10` (you have 10 top 4 sets)

---

## Method 2: Using Railway's Terminal (Alternative)

1. **In Railway Dashboard:**
   - Go to your project
   - Click on your **PostgreSQL** service
   - Click the **"Connect"** tab
   - Copy the **"Connection String"** (looks like: `postgresql://postgres:password@host:port/railway`)

2. **Open Terminal in Railway:**
   - Click on your **PostgreSQL** service
   - Click **"Connect"** tab
   - Click **"Open in Terminal"** (if available)
   - OR use your local terminal:

3. **Run Migrations Locally:**
   ```bash
   # Install psql if you don't have it (macOS)
   # brew install postgresql
   
   # Run first migration
   psql "<PASTE_CONNECTION_STRING_HERE>" -f backend/migrations/001_initial_schema.sql
   
   # Run second migration
   psql "<PASTE_CONNECTION_STRING_HERE>" -f backend/migrations/002_seed_data.sql
   ```

---

## Method 3: Using Railway's PostgreSQL Service Terminal

1. **In Railway Dashboard:**
   - Click on **PostgreSQL** service
   - Click **"Connect"** tab
   - Find **"psql"** command (looks like: `psql $DATABASE_URL`)
   - Click **"Open in Terminal"** or copy the command

2. **In the Terminal:**
   ```sql
   -- You'll be in psql now, run SQL directly:
   
   -- First, run the schema (copy/paste contents of 001_initial_schema.sql)
   -- Then run the seed data (copy/paste contents of 002_seed_data.sql)
   ```

---

## Method 4: Using a Migration Service (Advanced)

Create a one-time migration service:

1. **In Railway:**
   - Click **"+ New"** → **"Empty Service"**
   - Add environment variable: `DATABASE_URL` (copy from PostgreSQL service)
   - Set **Start Command**: 
     ```bash
     psql $DATABASE_URL -f migrations/001_initial_schema.sql && psql $DATABASE_URL -f migrations/002_seed_data.sql
     ```
   - Deploy, wait for completion, then delete the service

---

## What the Migrations Do

**001_initial_schema.sql:**
- Creates all database tables (users, movies, top4_sets, comparisons, votes)
- Sets up indexes and foreign keys

**002_seed_data.sql:**
- Inserts 40 sample movies
- Creates 10 sample top 4 sets
- Creates an anonymous user

---

## Verification Queries

After running migrations, verify everything worked:

```sql
-- Check movies
SELECT COUNT(*) FROM movies;
-- Should return: 40

-- Check top 4 sets
SELECT COUNT(*) FROM top4_sets;
-- Should return: 10

-- Check a specific movie
SELECT title, year, director FROM movies WHERE title = 'The Godfather';
-- Should return: The Godfather, 1972, Francis Ford Coppola

-- Check a top 4 set
SELECT id, movie_ids FROM top4_sets LIMIT 1;
-- Should show an array of 4 movie IDs
```

---

## Troubleshooting

**Error: "relation already exists"**
- Tables already exist, that's okay
- Just run the seed data migration (002)

**Error: "permission denied"**
- Make sure you're using the connection string from Railway
- Don't use localhost credentials

**Error: "could not connect"**
- Check that PostgreSQL service is running (green status)
- Verify the DATABASE_URL is correct

**Can't find the Query button:**
- Make sure you're in the PostgreSQL service (not backend service)
- Click "Data" tab, then "Query"

---

## Quick Copy-Paste Method (Recommended)

1. Railway → PostgreSQL service → "Data" tab → "Query"
2. Copy entire contents of `backend/migrations/001_initial_schema.sql`
3. Paste and run
4. Copy entire contents of `backend/migrations/002_seed_data.sql`
5. Paste and run
6. Done! ✅

