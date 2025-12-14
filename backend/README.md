# MovieMash Backend

Go backend for MovieMash application with PostgreSQL database.

## Setup

### Prerequisites
- Go 1.21 or later
- PostgreSQL 12 or later

### Database Setup

1. Create PostgreSQL database:
```bash
createdb moviemash
```

2. Create user (optional, or use existing user):
```bash
psql -d moviemash -c "CREATE USER moviemash WITH PASSWORD 'password';"
psql -d moviemash -c "GRANT ALL PRIVILEGES ON DATABASE moviemash TO moviemash;"
```

3. Run migrations:
```bash
psql -d moviemash -f migrations/001_initial_schema.sql
```

4. Seed sample data:
```bash
cd cmd/seed
go run main.go
```

Or use the SQL file directly:
```bash
psql -d moviemash -f migrations/002_seed_data.sql
```

### Environment Variables

Create a `.env` file in the `backend` directory:

```
DATABASE_URL=postgres://moviemash:password@localhost:5432/moviemash?sslmode=disable
```

## Database Schema

- **users**: User accounts
- **movies**: Movie information
- **top4_sets**: Collections of 4 movies (Letterboxd top 4s)
- **comparisons**: Pairwise comparisons between top 4 sets
- **votes**: User votes on comparisons

## Sample Data

The seed script includes:
- 40 popular movies
- 10 sample top 4 sets
- Anonymous user for testing

