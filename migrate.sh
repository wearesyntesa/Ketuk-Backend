#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
until pg_isready -h postgres -p 5432 -U user; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "PostgreSQL is ready. Running migrations..."

# Run migrations in order
echo "Running migration 000001_init.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000001_init.up.sql

echo "Running migration 000002_update_tickets_schema.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000002_update_tickets_schema.up.sql

echo "Running migration 000003_create_items_tables.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000003_create_items_tables.up.sql

echo "Running migration 000004_create_schedule_tables.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000004_create_schedule_tables.up.sql

echo "Running migration 000005_add_schedule_id_to_tickets.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000005_add_schedule_id_to_tickets.up.sql

echo "All migrations completed successfully!"
