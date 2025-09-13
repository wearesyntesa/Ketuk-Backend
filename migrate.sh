#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
until pg_isready -h postgres -p 5432 -U user; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "PostgreSQL is ready. Running migrations..."

# Run migrations
psql -h postgres -U user -d mydb < /migrations/000001_init.up.sql

echo "Migrations completed successfully!"
