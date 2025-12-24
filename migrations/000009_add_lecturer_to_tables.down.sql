-- ================================================
-- Rollback Migration: Remove lecturer column from tables
-- ================================================

-- Drop indexes first
DROP INDEX IF EXISTS idx_tickets_lecturer;
DROP INDEX IF EXISTS idx_schedule_ticket_lecturer;
DROP INDEX IF EXISTS idx_schedule_reguler_lecturer;

-- Remove lecturer column from tickets table
ALTER TABLE tickets
DROP COLUMN IF EXISTS lecturer;

-- Remove lecturer column from schedule_ticket table
ALTER TABLE schedule_ticket
DROP COLUMN IF EXISTS lecturer;

-- Remove lecturer column from schedule_reguler table
ALTER TABLE schedule_reguler
DROP COLUMN IF EXISTS lecturer;
