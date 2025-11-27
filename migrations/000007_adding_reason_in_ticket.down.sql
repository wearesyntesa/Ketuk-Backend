-- removing reason column from tickets table
ALTER TABLE tickets DROP COLUMN IF EXISTS reason;