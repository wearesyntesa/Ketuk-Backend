-- ================================================
-- Migration: Create ticket_event_log table for audit trails
-- PostgreSQL
-- ================================================

-- Create ENUM type for event actions
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ticket_event_action') THEN
        CREATE TYPE ticket_event_action AS ENUM (
            'created',
            'updated',
            'status_changed',
            'deleted',
            'assigned',
            'commented',
            'approved',
            'rejected'
        );
    END IF;
END $$;

-- Create ticket_event_log table
CREATE TABLE IF NOT EXISTS ticket_event_log (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    action ticket_event_action NOT NULL,
    old_value JSONB,
    new_value JSONB,
    changes JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_ticket_event_log_ticket_id ON ticket_event_log(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_event_log_user_id ON ticket_event_log(user_id);
CREATE INDEX IF NOT EXISTS idx_ticket_event_log_action ON ticket_event_log(action);
CREATE INDEX IF NOT EXISTS idx_ticket_event_log_created_at ON ticket_event_log(created_at DESC);

-- Add comment to table
COMMENT ON TABLE ticket_event_log IS 'Audit trail for all ticket-related events and changes';
COMMENT ON COLUMN ticket_event_log.action IS 'Type of action performed on the ticket';
COMMENT ON COLUMN ticket_event_log.old_value IS 'Previous state of the ticket (JSON)';
COMMENT ON COLUMN ticket_event_log.new_value IS 'New state of the ticket (JSON)';
COMMENT ON COLUMN ticket_event_log.changes IS 'Specific fields that were changed (JSON)';
COMMENT ON COLUMN ticket_event_log.ip_address IS 'IP address of the user who performed the action';
COMMENT ON COLUMN ticket_event_log.user_agent IS 'Browser/client information';
