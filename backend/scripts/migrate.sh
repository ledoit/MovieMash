#!/bin/bash
# Migration script for Railway deployment

set -e

echo "Running database migrations..."

# Get database URL from environment
DATABASE_URL="${DATABASE_URL}"

if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL not set"
    exit 1
fi

# Run migrations
echo "Running initial schema..."
psql "$DATABASE_URL" -f migrations/001_initial_schema.sql

echo "Seeding initial data..."
psql "$DATABASE_URL" -f migrations/002_seed_data.sql

echo "Migrations completed successfully!"

