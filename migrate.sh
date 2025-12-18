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

echo "Running migration 000006_add_password_to_users.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000006_add_password_to_users.up.sql

echo "Running migration 000007_adding_reason_in_ticket.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000007_adding_reason_in_ticket.up.sql

echo "Running migration 000008_create_ticket_event_log.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000008_create_ticket_event_log.up.sql

echo "Running migration 000009_add_lecturer_to_tables.up.sql..."
psql -h postgres -U user -d mydb < /migrations/000009_add_lecturer_to_tables.up.sql

echo "All migrations completed successfully!"
