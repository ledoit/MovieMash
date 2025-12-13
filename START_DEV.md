# Development Server Startup Guide

## Prerequisites Setup (Already Done ✅)
- ✅ Go installed
- ✅ Node.js and npm installed
- ✅ Python 3 installed
- ✅ PostgreSQL installed and running
- ✅ Database created and migrated
- ✅ All dependencies installed

## Starting the Development Servers

### 1. Start Backend API

Open a terminal and run:
```bash
cd backend
export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"
go run cmd/api/main.go
```

Backend will be available at: `http://localhost:8080`

### 2. Start Frontend

Open a **new terminal** and run:
```bash
cd frontend
npm start
# or
npx ng serve
```

Frontend will be available at: `http://localhost:4200`

### 3. (Optional) Start ML Service

Open a **new terminal** and run:
```bash
cd ml-service
source venv/bin/activate
python main.py
```

## Quick Test

Once both backend and frontend are running:

1. Open browser to `http://localhost:4200`
2. You should see the MovieMash interface
3. Try the different tabs:
   - **Gladiators**: Compare top 4s (needs data first)
   - **Recommendations**: View recommendations (needs votes first)
   - **Leaderboard**: View global rankings (needs votes first)

## Adding PostgreSQL to PATH (Permanent)

To avoid typing `export PATH` every time, add to your `~/.zshrc`:
```bash
echo 'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Troubleshooting

- **Backend won't start**: Check PostgreSQL is running: `brew services list | grep postgresql`
- **Database connection error**: Verify `.env` file has correct DATABASE_URL
- **Frontend won't start**: Run `npm install` again in frontend directory
- **Port already in use**: Change PORT in backend/.env or kill process using port 8080/4200

