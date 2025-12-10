-- ================================================
-- Rollback: Drop ticket_event_log table and related types
-- PostgreSQL
-- ================================================

-- Drop indexes
DROP INDEX IF EXISTS idx_ticket_event_log_created_at;
DROP INDEX IF EXISTS idx_ticket_event_log_action;
DROP INDEX IF EXISTS idx_ticket_event_log_user_id;
DROP INDEX IF EXISTS idx_ticket_event_log_ticket_id;

-- Drop table
DROP TABLE IF EXISTS ticket_event_log;

-- Drop ENUM type
DROP TYPE IF EXISTS ticket_event_action;
