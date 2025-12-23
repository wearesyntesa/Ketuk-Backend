-- ================================================
-- Migration: Drop form_links table
-- PostgreSQL
-- ================================================

-- Remove indexes
DROP INDEX IF EXISTS idx_tickets_form_link_id;
DROP INDEX IF EXISTS idx_form_links_is_active;
DROP INDEX IF EXISTS idx_form_links_expires_at;
DROP INDEX IF EXISTS idx_form_links_code;

-- Remove trigger
DROP TRIGGER IF EXISTS set_form_links_updated_at ON form_links;

-- Remove columns from tickets
ALTER TABLE tickets DROP COLUMN IF EXISTS submitter_phone;
ALTER TABLE tickets DROP COLUMN IF EXISTS submitter_email;
ALTER TABLE tickets DROP COLUMN IF EXISTS submitter_name;
ALTER TABLE tickets DROP COLUMN IF EXISTS form_link_id;

-- Drop table
DROP TABLE IF EXISTS form_links;
