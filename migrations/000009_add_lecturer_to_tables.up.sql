-- ================================================
-- Migration: Add lecturer column to tables
-- Add lecturer/dosen name field to tickets and schedule tables
-- ================================================

-- Add lecturer column to tickets table
ALTER TABLE tickets
ADD COLUMN IF NOT EXISTS lecturer VARCHAR(255);

-- Add lecturer column to schedule_ticket table
ALTER TABLE schedule_ticket
ADD COLUMN IF NOT EXISTS lecturer VARCHAR(255);

-- Add lecturer column to schedule_reguler table
ALTER TABLE schedule_reguler
ADD COLUMN IF NOT EXISTS lecturer VARCHAR(255);

-- Add index for better performance on queries
CREATE INDEX IF NOT EXISTS idx_tickets_lecturer ON tickets(lecturer);
CREATE INDEX IF NOT EXISTS idx_schedule_ticket_lecturer ON schedule_ticket(lecturer);
CREATE INDEX IF NOT EXISTS idx_schedule_reguler_lecturer ON schedule_reguler(lecturer);

-- Add comments for documentation
COMMENT ON COLUMN tickets.lecturer IS 'Name of the lecturer/dosen for the event';
COMMENT ON COLUMN schedule_ticket.lecturer IS 'Name of the lecturer/dosen for the scheduled event';
COMMENT ON COLUMN schedule_reguler.lecturer IS 'Name of the lecturer/dosen for the regular schedule';
