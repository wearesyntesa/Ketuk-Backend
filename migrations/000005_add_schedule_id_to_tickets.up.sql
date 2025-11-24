-- ================================================
-- Migration: Add schedule_id to tickets table
-- Add bidirectional foreign key references between tickets and schedule_ticket
-- ================================================

-- Add id_schedule column to tickets table
ALTER TABLE tickets 
ADD COLUMN IF NOT EXISTS id_schedule INTEGER;

-- Ensure ticket_id column exists in schedule_ticket (in case table was created before this column was added)
ALTER TABLE schedule_ticket 
ADD COLUMN IF NOT EXISTS ticket_id INTEGER;

-- Add foreign key constraint from tickets to schedule_ticket
ALTER TABLE tickets 
ADD CONSTRAINT fk_tickets_schedule 
FOREIGN KEY (id_schedule) 
REFERENCES schedule_ticket(id_schedule) 
ON DELETE SET NULL;

-- Add foreign key constraint from schedule_ticket to tickets
ALTER TABLE schedule_ticket 
ADD CONSTRAINT fk_schedule_ticket_tickets
FOREIGN KEY (ticket_id) 
REFERENCES tickets(id) 
ON DELETE CASCADE;

-- Add index for better performance on queries
CREATE INDEX IF NOT EXISTS idx_tickets_id_schedule ON tickets(id_schedule);

-- Add comments for documentation
COMMENT ON COLUMN tickets.id_schedule IS 'Foreign key reference to schedule_ticket table. NULL if ticket is not associated with a schedule.';
COMMENT ON COLUMN schedule_ticket.ticket_id IS 'Foreign key reference to tickets table. NULL if schedule was not created from a ticket.';
COMMENT ON TABLE schedule_ticket IS 'Table to store schedules created from tickets with bidirectional relationship';