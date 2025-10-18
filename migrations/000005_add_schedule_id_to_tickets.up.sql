-- ================================================
-- Migration: Add schedule_id to tickets table
-- Add foreign key reference to schedule_ticket table
-- ================================================

-- Add id_schedule column to tickets table
ALTER TABLE tickets 
ADD COLUMN IF NOT EXISTS id_schedule INTEGER;

-- Add foreign key constraint
ALTER TABLE tickets 
ADD CONSTRAINT fk_tickets_schedule 
FOREIGN KEY (id_schedule) 
REFERENCES schedule_ticket(id_schedule) 
ON DELETE SET NULL;

-- Add index for better performance on queries
CREATE INDEX IF NOT EXISTS idx_tickets_id_schedule ON tickets(id_schedule);

-- Add comment for documentation
COMMENT ON COLUMN tickets.id_schedule IS 'Foreign key reference to schedule_ticket table. NULL if ticket is not associated with a schedule.';

-- Remove fk id from tickets table in schedule_ticket table
ALTER TABLE schedule_ticket 
DROP CONSTRAINT IF EXISTS fk_schedule_ticket_tickets;
ALTER TABLE schedule_ticket
DROP COLUMN IF EXISTS ticket_id;
-- Add comment to schedule_ticket table
COMMENT ON TABLE schedule_ticket IS 'Table to store schedules created from tickets. ticket_id column removed