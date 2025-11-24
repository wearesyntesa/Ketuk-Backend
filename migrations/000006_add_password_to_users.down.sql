-- Remove password column from users table
ALTER TABLE users DROP COLUMN IF EXISTS password;

-- Restore google_sub NOT NULL constraint
ALTER TABLE users ALTER COLUMN google_sub SET NOT NULL;
