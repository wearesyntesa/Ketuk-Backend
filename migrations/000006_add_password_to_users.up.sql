-- Add password column to users table for JWT authentication
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255);

-- Make google_sub nullable since we now support email/password auth
ALTER TABLE users ALTER COLUMN google_sub DROP NOT NULL;
