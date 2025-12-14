#!/bin/bash
# Migration script for Railway
# Usage: ./run-migrations.sh <DATABASE_URL>

set -e

DATABASE_URL="${1:-$DATABASE_URL}"

if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL not provided"
    echo "Usage: ./run-migrations.sh <DATABASE_URL>"
    echo "   OR: DATABASE_URL=... ./run-migrations.sh"
    exit 1
fi

echo "Running database migrations..."
echo "Database: $DATABASE_URL"

echo ""
echo "Step 1: Running initial schema..."
psql "$DATABASE_URL" -f migrations/001_initial_schema.sql

echo ""
echo "Step 2: Seeding initial data..."
psql "$DATABASE_URL" -f migrations/002_seed_data.sql

echo ""
echo "✅ Migrations completed successfully!"
echo ""
echo "Verifying..."
psql "$DATABASE_URL" -c "SELECT COUNT(*) as movie_count FROM movies;"
psql "$DATABASE_URL" -c "SELECT COUNT(*) as set_count FROM top4_sets;"

