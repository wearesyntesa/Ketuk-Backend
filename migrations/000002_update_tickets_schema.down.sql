-- ================================================
-- Rollback: Update tickets table schema
-- Remove added columns for RequestData fields
-- ================================================

-- Remove the added columns
ALTER TABLE tickets 
DROP COLUMN IF EXISTS category,
DROP COLUMN IF EXISTS start_date,
DROP COLUMN IF EXISTS end_date,
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS pic;

-- Drop the added indexes
DROP INDEX IF EXISTS idx_tickets_category;
DROP INDEX IF EXISTS idx_tickets_email;
DROP INDEX IF EXISTS idx_tickets_request_date;