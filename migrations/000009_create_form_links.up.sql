-- ================================================
-- Migration: Create form_links table for public form feature
-- PostgreSQL
-- ================================================

-- Create table form_links
CREATE TABLE IF NOT EXISTS form_links (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pic_name VARCHAR(100) NOT NULL,
    pic_email VARCHAR(100) NOT NULL,
    pic_phone VARCHAR(20),
    expires_at TIMESTAMPTZ NOT NULL,
    max_submissions INT,
    submission_count INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add form_link_id to tickets table for tracking which form link a ticket came from
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS form_link_id INT REFERENCES form_links(id) ON DELETE SET NULL;

-- Add submitter info for public submissions (users without accounts)
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS submitter_name VARCHAR(100);
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS submitter_email VARCHAR(100);
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS submitter_phone VARCHAR(20);

-- Add trigger for auto-update timestamps
CREATE TRIGGER set_form_links_updated_at
BEFORE UPDATE ON form_links
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_form_links_code ON form_links(code);
CREATE INDEX IF NOT EXISTS idx_form_links_expires_at ON form_links(expires_at);
CREATE INDEX IF NOT EXISTS idx_form_links_is_active ON form_links(is_active);
CREATE INDEX IF NOT EXISTS idx_tickets_form_link_id ON tickets(form_link_id);
