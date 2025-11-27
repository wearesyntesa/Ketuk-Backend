-- adding reason column to tickets table
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS reason TEXT NOT NULL DEFAULT 'No reason provided';