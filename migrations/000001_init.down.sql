-- Drop triggers
DROP TRIGGER IF EXISTS set_tickets_updated_at ON tickets;
DROP TRIGGER IF EXISTS set_users_updated_at ON users;

-- Drop tables
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();
