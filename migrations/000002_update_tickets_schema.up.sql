-- ================================================
-- Migration: Update tickets table to match current model
-- Add missing columns for RequestData fields
-- ================================================
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ticket_category') THEN
        CREATE TYPE ticket_category AS ENUM ('Kelas', 'Lainnya', 'Praktikum', 'Skripsi');
    END IF;
END $$;

-- Add new columns to tickets table with default values
ALTER TABLE tickets 
ADD COLUMN IF NOT EXISTS category ticket_category DEFAULT 'Lainnya',
ADD COLUMN IF NOT EXISTS starts_date timestamptz DEFAULT,
ADD COLUMN IF NOT EXISTS ends_date timestamptz DEFAULT,
ADD COLUMN IF NOT EXISTS email VARCHAR(100) DEFAULT '',
ADD COLUMN IF NOT EXISTS phone VARCHAR(15) DEFAULT '',
ADD COLUMN IF NOT EXISTS pic VARCHAR(100) DEFAULT '';

-- Now make the columns NOT NULL after adding defaults
ALTER TABLE tickets 
ALTER COLUMN category SET NOT NULL,
ALTER COLUMN start_date SET NOT NULL,
ALTER COLUMN end_date SET NOT NULL,
ALTER COLUMN email SET NOT NULL,
ALTER COLUMN phone SET NOT NULL,
ALTER COLUMN pic SET NOT NULL;

-- Columns title and description already exist, no changes needed

-- Add index on commonly queried fields
CREATE INDEX IF NOT EXISTS idx_tickets_category ON tickets(category);
CREATE INDEX IF NOT EXISTS idx_tickets_email ON tickets(email);
CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets(created_at);

-- Add comments for documentation
COMMENT ON COLUMN tickets.category IS 'Ticket category (Kelas, Lainnya, Praktikum, Skripsi)';
COMMENT ON COLUMN tickets.start_date IS 'Start date in timestamp format';
COMMENT ON COLUMN tickets.end_date IS 'End date in timestamp format';
COMMENT ON COLUMN tickets.email IS 'Contact email for the ticket';
COMMENT ON COLUMN tickets.phone IS 'Contact phone number';
COMMENT ON COLUMN tickets.pic IS 'Person in charge';