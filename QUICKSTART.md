# Quick Start Guide

## Prerequisites

- PostgreSQL 15+ (local installation)
- Go 1.21+
- Node.js 18+ and npm
- Python 3.11+

## Setup Steps

### 1. Install and Setup PostgreSQL

**macOS** (using Homebrew):
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Linux** (Ubuntu/Debian):
```bash
sudo apt-get install postgresql-15
sudo systemctl start postgresql
```

**Windows**:
Download and install from [PostgreSQL Downloads](https://www.postgresql.org/download/windows/)

**Create Database**:
```bash
# Connect to PostgreSQL
psql postgres

# Create database and user
CREATE DATABASE proji;
CREATE USER moviemash WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE proji TO moviemash;
\q
```

### 2. Run Database Migrations

```bash
cd backend
go mod tidy
go run cmd/migrate/main.go
```

### 3. Start Backend API

```bash
cd backend
go run cmd/api/main.go
```

API will be available at `http://localhost:8080`

### 4. Start Frontend

```bash
cd frontend
npm install
ng serve
```

Frontend will be available at `http://localhost:4200`

### 5. Start ML Service (Optional)

```bash
cd ml-service
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
python main.py train  # One-time training
# OR
python main.py  # Continuous polling mode
```

### 6. Run Scraper (Optional)

```bash
cd scraper
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
# Edit main.py to add Letterboxd usernames
python main.py
```

## Testing the System

1. **Get a comparison**: `curl http://localhost:8080/api/v1/comparison`
2. **Submit a vote**: `curl -X POST http://localhost:8080/api/v1/votes -H "Content-Type: application/json" -d '{"comparison_id": 1, "winner_set_id": 1}'`
3. **View frontend**: Open `http://localhost:4200` in browser

## Environment Variables

Create `.env` files in each service directory:

**backend/.env**:
```
DATABASE_URL=postgres://moviemash:password@localhost:5432/proji?sslmode=disable
JWT_SECRET=your-secret-key
PORT=8080
```

**ml-service/.env**:
```
DATABASE_URL=postgres://moviemash:password@localhost:5432/proji?sslmode=disable
```

## Next Steps

1. **Populate data**: Run the scraper to get initial top 4 sets
2. **Implement auth**: Complete JWT authentication in backend
3. **Enhance ML model**: Tune hyperparameters, add features
4. **Add tests**: Unit and integration tests
5. **Deploy**: Set up production infrastructure

## Troubleshooting

- **Database connection errors**: 
  - Ensure PostgreSQL is running: `pg_isready` or `brew services list` (macOS)
  - Check connection string in `.env` file
  - Verify database exists: `psql -l`
- **Port conflicts**: PostgreSQL default port is 5432. Change if needed in connection string
- **Go module errors**: Run `go mod tidy` in backend directory
- **Permission errors**: Ensure database user has proper permissions

