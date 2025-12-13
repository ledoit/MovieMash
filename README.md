# MovieMash

A FaceMash-style web application for comparing Letterboxd top 4 movie sets, with aggregated votes powering a collaborative filtering recommendation engine and global top movies leaderboard.

## Tech Stack

- **Frontend**: Angular
- **Backend API**: Go
- **Database**: PostgreSQL (local installation)
- **ML Service**: PyTorch

## Project Structure

```
proji/
├── frontend/          # Angular application
├── backend/           # Go API service
├── scraper/           # Letterboxd scraper (Python)
├── ml-service/        # PyTorch recommendation engine
└── docs/              # Documentation
```

## Quick Start

1. Install PostgreSQL locally (see [QUICKSTART.md](./QUICKSTART.md) for details)

2. Run database migrations:
   ```bash
   cd backend && go run cmd/migrate/main.go
   ```

3. Start backend:
   ```bash
   cd backend && go run cmd/api/main.go
   ```

4. Start frontend:
   ```bash
   cd frontend && npm install && ng serve
   ```

5. Start ML service:
   ```bash
   cd ml-service && python -m venv venv && source venv/bin/activate
   pip install -r requirements.txt
   python main.py
   ```

## Development

See [SYSTEM_DESIGN.md](./SYSTEM_DESIGN.md) for detailed architecture.

