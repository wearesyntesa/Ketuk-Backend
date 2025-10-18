-- ================================================
-- Rollback Migration: Remove schedule_id from tickets table
-- ================================================

-- Drop the index
DROP INDEX IF EXISTS idx_tickets_id_schedule;

-- Drop the foreign key constraint
ALTER TABLE tickets 
DROP CONSTRAINT IF EXISTS fk_tickets_schedule;

-- Drop the column
ALTER TABLE tickets 
DROP COLUMN IF EXISTS id_schedule;
-- ================================================
-- Rollback Migration: Re-add ticket_id to schedule_ticket table
-- ================================================
ALTER TABLE schedule_ticket
ADD COLUMN IF NOT EXISTS ticket_id INTEGER NOT NULL;
ALTER TABLE schedule_ticket
ADD CONSTRAINT fk_schedule_ticket_tickets
FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE;
-- Remove comment from schedule_ticket table
COMMENT ON TABLE schedule_ticket IS NULL;
